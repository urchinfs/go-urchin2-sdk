package options

import (
	"errors"
	"fmt"

	dag "github.com/ipfs/boxo/ipld/merkledag"
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

type Layout int

const (
	BalancedLayout Layout = iota
	TrickleLayout
)

type UnixfsAddSettings struct {
	CidVersion int
	MhType     uint64

	Inline       bool
	InlineLimit  int
	RawLeaves    bool
	RawLeavesSet bool

	Chunker string
	Layout  Layout

	Events chan<- interface{}
	Silent bool
}

type UnixfsLsSettings struct {
	ResolveChildren   bool
	UseCumulativeSize bool
}

type (
	UnixfsAddOption func(*UnixfsAddSettings) error
	UnixfsLsOption  func(*UnixfsLsSettings) error
)

func UnixfsAddOptions(opts ...UnixfsAddOption) (*UnixfsAddSettings, cid.Prefix, error) {
	options := &UnixfsAddSettings{
		CidVersion: -1,
		MhType:     mh.SHA2_256,

		Inline:       false,
		InlineLimit:  32,
		RawLeaves:    false,
		RawLeavesSet: false,

		Chunker: "size-262144",
		Layout:  BalancedLayout,
		Events:  nil,
	}

	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, cid.Prefix{}, err
		}
	}

	if options.MhType != mh.SHA2_256 {
		switch options.CidVersion {
		case 0:
			return nil, cid.Prefix{}, errors.New("CIDv0 only supports sha2-256")
		case 1, -1:
			options.CidVersion = 1
		default:
			return nil, cid.Prefix{}, fmt.Errorf("unknown CID version: %d", options.CidVersion)
		}
	} else {
		if options.CidVersion < 0 {
			// Default to CIDv0
			options.CidVersion = 0
		}
	}

	if options.CidVersion > 0 && !options.RawLeavesSet {
		options.RawLeaves = true
	}

	prefix, err := dag.PrefixForCidVersion(options.CidVersion)
	if err != nil {
		return nil, cid.Prefix{}, err
	}

	prefix.MhType = options.MhType
	prefix.MhLength = -1

	return options, prefix, nil
}

func UnixfsLsOptions(opts ...UnixfsLsOption) (*UnixfsLsSettings, error) {
	options := &UnixfsLsSettings{
		ResolveChildren: true,
	}

	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}

	return options, nil
}

type unixfsOpts struct{}

var Unixfs unixfsOpts

func (unixfsOpts) CidVersion(version int) UnixfsAddOption {
	return func(settings *UnixfsAddSettings) error {
		settings.CidVersion = version
		return nil
	}
}

func (unixfsOpts) Hash(mhtype uint64) UnixfsAddOption {
	return func(settings *UnixfsAddSettings) error {
		settings.MhType = mhtype
		return nil
	}
}

func (unixfsOpts) RawLeaves(enable bool) UnixfsAddOption {
	return func(settings *UnixfsAddSettings) error {
		settings.RawLeaves = enable
		settings.RawLeavesSet = true
		return nil
	}
}

func (unixfsOpts) Inline(enable bool) UnixfsAddOption {
	return func(settings *UnixfsAddSettings) error {
		settings.Inline = enable
		return nil
	}
}

func (unixfsOpts) InlineLimit(limit int) UnixfsAddOption {
	return func(settings *UnixfsAddSettings) error {
		settings.InlineLimit = limit
		return nil
	}
}

func (unixfsOpts) Chunker(chunker string) UnixfsAddOption {
	return func(settings *UnixfsAddSettings) error {
		settings.Chunker = chunker
		return nil
	}
}

func (unixfsOpts) Layout(layout Layout) UnixfsAddOption {
	return func(settings *UnixfsAddSettings) error {
		settings.Layout = layout
		return nil
	}
}

func (unixfsOpts) Events(sink chan<- interface{}) UnixfsAddOption {
	return func(settings *UnixfsAddSettings) error {
		settings.Events = sink
		return nil
	}
}

func (unixfsOpts) ResolveChildren(resolve bool) UnixfsLsOption {
	return func(settings *UnixfsLsSettings) error {
		settings.ResolveChildren = resolve
		return nil
	}
}

func (unixfsOpts) UseCumulativeSize(use bool) UnixfsLsOption {
	return func(settings *UnixfsLsSettings) error {
		settings.UseCumulativeSize = use
		return nil
	}
}
