package car

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/exchange/offline"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/ipld/merkledag"
	unixfile "github.com/ipfs/boxo/ipld/unixfs/file"
	"github.com/ipfs/boxo/tar"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/sync"
	format "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log"
	carv1 "github.com/ipld/go-car"
	selectorparse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
	"io"
	"os"
	"path"
)

var log = logging.Logger("car")

type identityWriteCloser struct {
	w io.Writer
}

func (i *identityWriteCloser) Write(p []byte) (int, error) {
	return i.w.Write(p)
}

func (i *identityWriteCloser) Close() error {
	return nil
}

type ImportEvent struct {
	Name  string
	CID   cid.Cid
	Bytes int64
	Size  string
}

type Builder struct {
	di *DataImporter
	wt io.WriterTo
}

func NewBuilder() *Builder {
	return &Builder{
		di: NewDataImporter(),
	}
}

type CarV1 struct {
	root cid.Cid
	car  *carv1.SelectiveCar
}

func (c *CarV1) Write(w io.Writer) error {
	return c.car.Write(w)
}

func (c *CarV1) Root() cid.Cid {
	return c.root
}

func (b *Builder) BuildCar(
	ctx context.Context,
	input any,
	opts ...ImportOption,
) (*CarV1, error) {
	root, err := b.di.Import(ctx, input, opts...)
	if err != nil {
		return nil, err
	}

	car := carv1.NewSelectiveCar(
		ctx,
		b.di.Blockstore(),
		[]carv1.Dag{
			carv1.Dag{
				Root:     root,
				Selector: selectorparse.CommonSelector_ExploreAllRecursively,
			},
		},
		carv1.TraverseLinksOnlyOnce(),
	)

	return &CarV1{
		root: root,
		car:  &car,
	}, nil
}

func PackCarFormat(input, output string) (string, error) {
	rootCid := ""
	_, err := os.Stat(input)
	if err != nil {
		log.Fatalf("Failed to stat input path, err: %v", err)
		return rootCid, err
	}

	iFd, err := os.Open(input)
	if err != nil {
		log.Fatalf("Failed to open input path, err: %v", err)
		return rootCid, err
	}
	defer iFd.Close()

	oFd, err := os.Create(output)
	if err != nil {
		log.Fatalf("Failed to create output CAR path, err: %v", err)
		return rootCid, err
	}
	defer oFd.Close()

	b := NewBuilder()
	v1car, err := b.BuildCar(
		context.TODO(),
		input,
		ImportOpts.CIDv0(),
	)
	if err != nil {
		log.Fatalf("Error creating car v1 builder: %v", err)
		return rootCid, err
	}

	if err := v1car.Write(oFd); err != nil {
		log.Fatalf("Error writing out car v0 format: %v", err)
		return rootCid, err
	}

	rootCid = v1car.Root().String()
	log.Infof("Car v1 generated, CID =%v", v1car.Root())
	return rootCid, nil
}

func UnpackCarFormat(input, output string) error {
	iFd, err := os.Open(input)
	if err != nil {
		log.Fatalf("Failed to open input file, err: %v", err)
		return err
	}
	defer iFd.Close()

	stat, err := os.Stat(output)
	if err != nil {
		log.Fatalf("Failed to stat output file, err: %v", err)
		return err
	}
	if !stat.IsDir() {
		log.Fatalf("output path is not directory")
		return err
	}

	ds := sync.MutexWrap(datastore.NewMapDatastore())
	bs := blockstore.NewBlockstore(ds)
	bsvc := blockservice.New(bs, offline.Exchange(bs))
	roots, err := carv1.LoadCar(context.Background(), bs, iFd)
	if err != nil {
		log.Fatalf("Failed to load CAR file: %v", err)
		return err
	}

	dagService := merkledag.NewDAGService(bsvc)
	rootCid := roots.Roots[0]
	outputDir := path.Join(output, rootCid.String())
	err = restoreFilesFromDag(dagService, rootCid, outputDir)
	if err != nil {
		log.Fatalf("Failed to restore files from DAG: %v", err)
		return err
	}

	fmt.Println("successfully restored files to:", outputDir)
	return nil
}

func restoreFilesFromDag(dagService format.DAGService, rootCid cid.Cid, outputDir string) error {
	cleaned := path.Clean(outputDir)
	_, filename := path.Split(cleaned)

	rootNode, err := dagService.Get(context.Background(), rootCid)
	if err != nil {
		log.Fatalf("Failed to get root node: %v", err)
		return fmt.Errorf("failed to get root node: %v", err)
	}

	file, err := unixfile.NewUnixfsFile(context.Background(), dagService, rootNode)
	if err != nil {
		log.Fatalf("Failed to get unixfs file: %v", err)
		return err
	}

	piper, pipew := io.Pipe()
	checkErrAndClosePipe := func(err error) bool {
		if err != nil {
			_ = pipew.CloseWithError(err)
			return true
		}
		return false
	}

	bufw := bufio.NewWriterSize(pipew, 1048576)
	maybeGzw := &identityWriteCloser{bufw}
	w, err := files.NewTarWriter(maybeGzw)

	closeGzwAndPipe := func() {
		if err := maybeGzw.Close(); checkErrAndClosePipe(err) {
			return
		}
		if err := bufw.Flush(); checkErrAndClosePipe(err) {
			return
		}
		_ = pipew.Close()
	}

	go func() {
		if err := w.WriteFile(file, filename); checkErrAndClosePipe(err) {
			return
		}
		_ = w.Close()
		closeGzwAndPipe()
	}()

	extractor := &tar.Extractor{Path: outputDir}
	return extractor.Extract(piper)
}
