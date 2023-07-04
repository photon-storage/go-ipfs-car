package car

import (
	"context"
	"io"

	ipldcar "github.com/ipld/go-car/v2"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/storage/bsadapter"
	selectorparse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
)

type Builder struct {
	di *DataImporter
	wt io.WriterTo
}

func NewBuilder() *Builder {
	return &Builder{
		di: NewDataImporter(),
	}
}

func (b *Builder) Build(
	ctx context.Context,
	input any,
	opts ...ImportOption,
) (io.WriterTo, error) {
	root, err := b.di.Import(ctx, input, opts...)
	if err != nil {
		return nil, err
	}

	ls := cidlink.DefaultLinkSystem()
	ls.SetReadStorage(&bsadapter.Adapter{
		Wrapped: b.di.Blockstore(),
	})

	return ipldcar.NewSelectiveWriter(
		ctx,
		&ls,
		root,
		selectorparse.CommonSelector_ExploreAllRecursively,
	)
}
