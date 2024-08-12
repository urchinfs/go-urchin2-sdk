package cid

import (
	"context"
	"github.com/ipfs/boxo/blockservice"
	bstore "github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/ipld/merkledag"
	dagtest "github.com/ipfs/boxo/ipld/merkledag/test"
	ft "github.com/ipfs/boxo/ipld/unixfs"
	"github.com/ipfs/boxo/mfs"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/urchinfs/go-urchin2-sdk/cid/options"
	"github.com/urchinfs/go-urchin2-sdk/utils"
	"os"
)

type AddEvent struct {
	Name  string
	Path  path.ImmutablePath `json:",omitempty"`
	Bytes int64              `json:",omitempty"`
	Size  string             `json:",omitempty"`
}

// FileType is an enum of possible UnixFS file types.
type FileType int32

const (
	TUnknown FileType = iota
	TFile
	TDirectory
	TSymlink
)

func (t FileType) String() string {
	switch t {
	case TUnknown:
		return "unknown"
	case TFile:
		return "file"
	case TDirectory:
		return "directory"
	case TSymlink:
		return "symlink"
	default:
		return "<unknown file type>"
	}
}

type DirEntry struct {
	Name string
	Cid  cid.Cid

	Size   uint64
	Type   FileType
	Target string

	Err error
}

func AddAndBuildCid(ctx context.Context, files files.Node, opts ...options.UnixfsAddOption) (cid.Cid, error) {
	settings, prefix, err := options.UnixfsAddOptions(opts...)
	if err != nil {
		return cid.Cid{}, err
	}

	dstore := dssync.MutexWrap(ds.NewNullDatastore())
	bs := bstore.NewBlockstore(dstore, bstore.WriteThrough())
	addblockstore := bstore.NewGCBlockstore(bs, nil)

	bserv := blockservice.New(addblockstore, nil)
	dserv := merkledag.NewDAGService(bserv)

	syncDserv := &SyncDagService{
		DAGService: dserv,
		syncFn:     func() error { return nil },
	}

	fileAdder, err := NewAdder(ctx, addblockstore, syncDserv)
	if err != nil {
		return cid.Cid{}, err
	}

	fileAdder.Chunker = settings.Chunker
	if settings.Events != nil {
		fileAdder.Out = settings.Events
	}
	fileAdder.Silent = settings.Silent
	fileAdder.RawLeaves = settings.RawLeaves
	fileAdder.CidBuilder = prefix

	md := dagtest.Mock()
	emptyDirNode := ft.EmptyDirNode()
	// Use the same prefix for the "empty" MFS root as for the file adder.
	err = emptyDirNode.SetCidBuilder(fileAdder.CidBuilder)
	if err != nil {
		return cid.Cid{}, err
	}
	mr, err := mfs.NewRoot(ctx, md, emptyDirNode, nil)
	if err != nil {
		return cid.Cid{}, err
	}

	fileAdder.SetMfsRoot(mr)

	nd, err := fileAdder.AddAll(ctx, files)
	if err != nil {
		return cid.Cid{}, err
	}

	return nd.Cid(), nil
}

type SyncDagService struct {
	ipld.DAGService
	syncFn func() error
}

func (s *SyncDagService) Sync() error {
	return s.syncFn()
}

func AppendFile(fpath string, recursive bool, filter *files.Filter) (files.Node, error) {
	stat, err := os.Lstat(fpath)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		if !recursive {
			return nil, utils.ErrNotSupported
		}
	} else if (stat.Mode() & os.ModeNamedPipe) != 0 {
		file, err := os.Open(fpath)
		if err != nil {
			return nil, err
		}

		return files.NewReaderFile(file), nil
	}
	return files.NewSerialFileWithFilter(fpath, filter, stat)
}
