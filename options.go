package car

import (
	"errors"

	"github.com/ipfs/boxo/coreiface/options"
	mh "github.com/multiformats/go-multihash"
)

var (
	ErrInvalidMhType           = errors.New("invalid multihash type")
	ErrIncompactibleCidVersion = errors.New("incompactible CID version")
)

type importOptions struct {
	cidVersion   int
	mhType       uint64
	rawLeaves    bool
	rawLeavesSet bool
	inline       bool
	inlineLimit  int
	chunker      string
	layout       options.Layout
	out          chan *ImportEvent
}

func buildImportOptions(opts ...ImportOption) (*importOptions, error) {
	ioptions := &importOptions{
		cidVersion:   1,
		mhType:       mh.SHA2_256,
		rawLeaves:    false,
		rawLeavesSet: false,
		inline:       false,
		inlineLimit:  32,
		chunker:      "size-262144",
		layout:       options.BalancedLayout,
	}

	for _, opt := range opts {
		if err := opt(ioptions); err != nil {
			return nil, err
		}
	}

	if ioptions.mhType != mh.SHA2_256 && ioptions.cidVersion != 1 {
		return nil, ErrIncompactibleCidVersion
	}

	if ioptions.cidVersion == 1 && !ioptions.rawLeavesSet {
		ioptions.rawLeaves = true
	}

	return ioptions, nil
}

type ImportOption func(*importOptions) error

type importScope struct{}

var ImportOpts importScope

// CIDv0 uses CID v0.
func (importScope) CIDv0() ImportOption {
	return func(opts *importOptions) error {
		opts.cidVersion = 0
		return nil
	}
}

// CIDv1 uses CID v1 (default).
func (importScope) CIDv1() ImportOption {
	return func(opts *importOptions) error {
		opts.cidVersion = 1
		return nil
	}
}

// MhType sets multihash type to use (default: mh.SHA2_256).
func (importScope) MhType(code uint64) ImportOption {
	return func(opts *importOptions) error {
		_, found := mh.Codes[code]
		if !found {
			return ErrInvalidMhType
		}

		opts.mhType = code
		return nil
	}
}

// RawLeaves enables raw leaves in the DAG tree (default to false for CIDv0,
// true for CIDv1).
func (importScope) RawLeaves(enabled bool) ImportOption {
	return func(opts *importOptions) error {
		opts.rawLeaves = enabled
		opts.rawLeavesSet = true
		return nil
	}
}

// InlineBlock enables inline small blocks into CID (default false).
func (importScope) InlineBlock() ImportOption {
	return func(opts *importOptions) error {
		opts.inline = true
		return nil
	}
}

// InlineBlockLimit sets the threshold for triggering inline (default 32).
func (importScope) InlineBlockLimit(limit int) ImportOption {
	return func(opts *importOptions) error {
		opts.inlineLimit = limit
		return nil
	}
}

// Chunker sets chunker configuration (default size-262144).
func (importScope) Chunker(chunker string) ImportOption {
	return func(opts *importOptions) error {
		opts.chunker = chunker
		return nil
	}
}

// BalancedLayout uses balanced DAG layout.
func (importScope) BalancedLayout() ImportOption {
	return func(opts *importOptions) error {
		opts.layout = options.BalancedLayout
		return nil
	}
}

// TrickleLayout uses trickle DAG layout.
func (importScope) TrickleLayout() ImportOption {
	return func(opts *importOptions) error {
		opts.layout = options.TrickleLayout
		return nil
	}
}

// Events sets event channel to receive import progress.
func (importScope) Events(ch chan *ImportEvent) ImportOption {
	return func(opts *importOptions) error {
		opts.out = ch
		return nil
	}
}
