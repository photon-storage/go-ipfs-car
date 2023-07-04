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

func (importScope) CIDv0() ImportOption {
	return func(opts *importOptions) error {
		opts.cidVersion = 0
		return nil
	}
}

func (importScope) CIDv1() ImportOption {
	return func(opts *importOptions) error {
		opts.cidVersion = 1
		return nil
	}
}

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

func (importScope) RawLeaves(enabled bool) ImportOption {
	return func(opts *importOptions) error {
		opts.rawLeaves = enabled
		opts.rawLeavesSet = true
		return nil
	}
}

func (importScope) InlineBlock() ImportOption {
	return func(opts *importOptions) error {
		opts.inline = true
		return nil
	}
}

func (importScope) InlineBlockLimit(limit int) ImportOption {
	return func(opts *importOptions) error {
		opts.inlineLimit = limit
		return nil
	}
}

func (importScope) Chunker(chunker string) ImportOption {
	return func(opts *importOptions) error {
		opts.chunker = chunker
		return nil
	}
}

func (importScope) BalancedLayout() ImportOption {
	return func(opts *importOptions) error {
		opts.layout = options.BalancedLayout
		return nil
	}
}

func (importScope) TrickleLayoutLayout() ImportOption {
	return func(opts *importOptions) error {
		opts.layout = options.TrickleLayout
		return nil
	}
}

func (importScope) Events(ch chan *ImportEvent) ImportOption {
	return func(opts *importOptions) error {
		opts.out = ch
		return nil
	}
}
