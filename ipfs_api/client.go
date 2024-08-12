package ipfs_api

import (
	"context"
	"encoding/json"
	"fmt"
	logging "github.com/ipfs/go-log"
	"github.com/urchinfs/go-urchin2-sdk/utils"
	"io"
	"net"
	gohttp "net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/blang/semver/v4"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/tar"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

const (
	DefaultApiFile = "api/v0"
)

var log = logging.Logger("api")

type HttpClient struct {
	url     string
	httpCli gohttp.Client

	versionMu sync.Mutex
	version   *semver.Version
}

func NewClient(url string) *HttpClient {
	c := &gohttp.Client{
		Transport: &gohttp.Transport{
			Proxy:             gohttp.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}

	var client HttpClient
	client.url = url
	client.httpCli = *c
	client.httpCli.CheckRedirect = func(_ *gohttp.Request, _ []*gohttp.Request) error {
		return fmt.Errorf("unexpected redirect")
	}

	maddr, err := ma.NewMultiaddr(url)
	if err != nil {
		return &client
	}

	network, host, err := manet.DialArgs(maddr)
	if err != nil {
		return &client
	}

	if network == "unix" {
		client.url = network

		var tptCopy *gohttp.Transport
		if tpt, ok := client.httpCli.Transport.(*gohttp.Transport); ok && tpt.DialContext == nil {
			tptCopy = tpt.Clone()
		} else if client.httpCli.Transport == nil {
			tptCopy = &gohttp.Transport{
				Proxy:             gohttp.ProxyFromEnvironment,
				DisableKeepAlives: true,
			}
		} else {
			return &client
		}

		tptCopy.DialContext = func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", host)
		}
		client.httpCli.Transport = tptCopy
	} else {
		client.url = host
	}

	return &client
}

var encodedAbsolutePathVersion = semver.MustParse("0.28.0-dev")

func (h *HttpClient) Version() (string, string, error) {
	ver := struct {
		Version string
		Commit  string
	}{}

	if err := h.Request("version").Exec(context.Background(), &ver); err != nil {
		return "", "", err
	}

	return ver.Version, ver.Commit, nil
}

func (h *HttpClient) loadRemoteVersion() (*semver.Version, error) {
	h.versionMu.Lock()
	defer h.versionMu.Unlock()

	if h.version == nil {
		version, _, err := h.Version()
		if err != nil {
			return nil, err
		}

		remoteVersion, err := semver.New(version)
		if err != nil {
			return nil, err
		}

		h.version = remoteVersion
	}

	return h.version, nil
}

func (h *HttpClient) newMultiFileReader(dir files.Directory) (*files.MultiFileReader, error) {
	version, err := h.loadRemoteVersion()
	if err != nil {
		return nil, err
	}

	return files.NewMultiFileReader(dir, true, version.LT(encodedAbsolutePathVersion)), nil
}

func (h *HttpClient) SetTimeout(d time.Duration) {
	h.httpCli.Timeout = d
}

func (h *HttpClient) Request(command string, args ...string) *RequestBuilder {
	return &RequestBuilder{
		command: command,
		args:    args,
		client:  h,
	}
}

func (h *HttpClient) Cat(path string) (io.ReadCloser, error) {
	resp, err := h.Request("cat", path).Send(context.Background())
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Output, nil
}

type LsLink struct {
	Hash string
	Name string
	Size uint64
	Type int
}

type LsObject struct {
	Links []*LsLink
	LsLink
}

func (h *HttpClient) List(path string) ([]*LsLink, error) {
	var out struct{ Objects []LsObject }
	err := h.Request("ls", path).Exec(context.Background(), &out)
	if err != nil {
		return nil, err
	}
	if len(out.Objects) != 1 {
		return nil, utils.ErrBadResponse
	}
	return out.Objects[0].Links, nil
}

func (h *HttpClient) Get(hash, outDir string) error {
	stat, err := os.Stat(outDir)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		outDir = path.Join(outDir, hash)
	}

	resp, err := h.Request("get", hash).Option("create", true).Send(context.Background())
	if err != nil {
		return err
	}
	defer resp.Close()

	if resp.Error != nil {
		return resp.Error
	}

	extractor := &tar.Extractor{Path: outDir}
	return extractor.Extract(resp.Output)
}

type SwarmStreamInfo struct {
	Protocol string
}

type SwarmConnInfo struct {
	Addr    string
	Peer    string
	Latency string
	Muxer   string
	Streams []SwarmStreamInfo
}

type SwarmConnInfos struct {
	Peers []SwarmConnInfo
}

func (h *HttpClient) SwarmPeers(ctx context.Context) (*SwarmConnInfos, error) {
	v := &SwarmConnInfos{}
	err := h.Request("swarm/peers").Exec(ctx, &v)
	return v, err
}

type swarmConnection struct {
	Strings []string
}

func (h *HttpClient) SwarmConnect(ctx context.Context, addr ...string) error {
	var conn *swarmConnection
	err := h.Request("swarm/connect").
		Arguments(addr...).
		Exec(ctx, &conn)
	return err
}

type object struct {
	Hash string
}

type AddOpts = func(*RequestBuilder) error

func Pin(enabled bool) AddOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("pin", enabled)
		return nil
	}
}

func (h *HttpClient) Add(inputFile string, options ...AddOpts) (string, error) {
	stat, err := os.Stat(inputFile)
	if err != nil {
		return "", err
	}
	if stat.IsDir() {
		return "", utils.ErrNotFile
	}

	wrapDataDir, err := utils.WarpPath(inputFile)
	if err != nil {
		return "", err
	}

	fileReader, err := h.newMultiFileReader(wrapDataDir)
	if err != nil {
		log.Errorf("new multi file reader err:%v", err)
		return "", err
	}

	rb := h.Request("add").Option("wrap-with-directory", true)
	for _, option := range options {
		option(rb)
	}

	resp, err := rb.Body(fileReader).Send(context.Background())
	if err != nil {
		log.Errorf("send http err:%v", err)
		return "", err
	}
	defer resp.Close()

	if resp.Error != nil {
		return "", resp.Error
	}

	dec := json.NewDecoder(resp.Output)
	var final string
	for {
		var out object
		err = dec.Decode(&out)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		final = out.Hash
	}

	if final == "" {
		return "", utils.ErrNotReceiveRet
	}

	log.Debugf("received warp hash: %s", final)
	return final, nil
}

func (h *HttpClient) AddNoPin(inputFile string) (string, error) {
	return h.Add(inputFile, Pin(false))
}

func (h *HttpClient) AddDir(dir string, options ...AddOpts) (string, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		return "", utils.ErrNotDir
	}

	wrapDataDir, err := utils.WarpPath(dir)
	if err != nil {
		return "", err
	}

	reader, err := h.newMultiFileReader(wrapDataDir)
	if err != nil {
		return "", err
	}

	rb := h.Request("add").Option("recursive", true).Option("wrap-with-directory", true)
	for _, option := range options {
		option(rb)
	}

	resp, err := rb.Body(reader).Send(context.Background())
	if err != nil {
		return "", err
	}
	defer resp.Close()

	if resp.Error != nil {
		return "", resp.Error
	}

	dec := json.NewDecoder(resp.Output)
	var final string
	for {
		var out object
		err = dec.Decode(&out)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		final = out.Hash
	}

	if final == "" {
		log.Warnf("no results received from ipfs peer")
		return "", utils.ErrNotReceiveRet
	}

	return final, nil
}
