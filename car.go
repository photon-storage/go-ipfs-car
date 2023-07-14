package car

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
	carv1 "github.com/ipld/go-car"
	carv2 "github.com/ipld/go-car/v2"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/storage/bsadapter"
	selectorparse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
)

// Builder builds a WriterTo from the given input source.
// The WriterTo can be used to output car format data to a io.Writer.
type Builder struct {
	di *DataImporter
	wt io.WriterTo
}

// NewBuilder creates a new car builder.
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

// Buildv1 builds a CarV1 for outputing car v1 format data.
func (b *Builder) Buildv1(
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

type CarV2 struct {
	root cid.Cid
	io.WriterTo
}

func (c *CarV2) Root() cid.Cid {
	return c.root
}

// Buildv2 builds a CarV2 for outputing car v2 format data.
func (b *Builder) Buildv2(
	ctx context.Context,
	input any,
	opts ...ImportOption,
) (*CarV2, error) {
	root, err := b.di.Import(ctx, input, opts...)
	if err != nil {
		return nil, err
	}

	ls := cidlink.DefaultLinkSystem()
	ls.SetReadStorage(&bsadapter.Adapter{
		Wrapped: b.di.Blockstore(),
	})

	w, err := carv2.NewSelectiveWriter(
		ctx,
		&ls,
		root,
		selectorparse.CommonSelector_ExploreAllRecursively,
	)
	if err != nil {
		return nil, err
	}

	return &CarV2{
		root:     root,
		WriterTo: w,
	}, nil
}
