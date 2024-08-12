package ipfs_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/urchinfs/go-ipfs-sdk/ipfs_api/options"
	"io"
	"os"
	"strings"

	"github.com/ipfs/boxo/files"
)

type DagImportRoot struct {
	Root struct {
		Cid struct {
			Value string `json:"/"`
		}
	}
	Stats *DagImportStats `json:"Stats,omitempty"`
}

type DagImportStats struct {
	BlockBytesCount uint64
	BlockCount      uint64
}

type DagImportOutput struct {
	Roots []DagImportRoot
	Stats *DagImportStats
}

func (h *HttpClient) DagGet(ref string, out interface{}) error {
	return h.Request("dag/get", ref).Exec(context.Background(), out)
}
func (h *HttpClient) DagImport(input string, silent, stats bool) (*DagImportOutput, error) {
	iFd, err := os.Open(input)
	if err != nil {
		return nil, err
	}

	rc := io.ReadCloser(iFd)
	return h.DagImportWithOpts(rc, options.Dag.Silent(silent), options.Dag.Stats(stats))
}

func (h *HttpClient) DagImportWithOpts(data interface{}, opts ...options.DagImportOption) (*DagImportOutput, error) {
	cfg, err := options.DagImportOptions(opts...)
	if err != nil {
		return nil, err
	}

	fileReader, err := h.dagToFilesReader(data)
	if err != nil {
		return nil, err
	}

	res, err := h.Request("dag/import").
		Option("pin-roots", cfg.PinRoots).
		Option("silent", cfg.Silent).
		Option("stats", cfg.Stats).
		Body(fileReader).
		Send(context.Background())
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.Error != nil {
		return nil, res.Error
	}

	if cfg.Silent {
		return nil, nil
	}

	out := DagImportOutput{
		Roots: []DagImportRoot{},
	}

	dec := json.NewDecoder(res.Output)

	for {
		var root DagImportRoot
		err := dec.Decode(&root)
		if err == io.EOF {
			break
		}

		if root.Stats != nil {
			out.Stats = root.Stats
			break
		}

		out.Roots = append(out.Roots, root)
	}

	return &out, err
}

func (h *HttpClient) dagToFilesReader(data interface{}) (*files.MultiFileReader, error) {
	var r io.Reader
	switch data := data.(type) {
	case *files.MultiFileReader:
		return data, nil
	case string:
		r = strings.NewReader(data)
	case []byte:
		r = bytes.NewReader(data)
	case io.Reader:
		r = data
	default:
		return nil, fmt.Errorf("values of type %T cannot be handled as DAG input", data)
	}

	fr := files.NewReaderFile(r)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
	return h.newMultiFileReader(slf)
}

func (h *HttpClient) DagExport(hash, outputFile string) error {
	resp, err := h.Request("dag/export", hash).Send(context.Background())
	if err != nil {
		return err
	}
	defer resp.Close()

	if resp.Error != nil {
		return resp.Error
	}
	defer resp.Output.Close()

	oFd, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output CAR file: %v", err)
	}
	defer oFd.Close()

	written, err := io.Copy(oFd, resp.Output)
	if err != nil {
		log.Fatalf("Failed to write output CAR file: %v", err)
		_ = os.Remove(outputFile)
		return err
	}

	fmt.Printf("io.Copy write:%v\n", written)
	return nil
}
