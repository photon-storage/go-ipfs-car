package car

import (
	"context"
	"fmt"
	"testing"

	"github.com/photon-storage/go-common/testing/require"
)

func TestImport(t *testing.T) {
	ctx := context.Background()
	pi := NewDataImporter()

	// CID v0
	cid, err := pi.Import(
		ctx,
		"./data",
		ImportOpts.CIDv0(),
	)
	require.NoError(t, err)
	require.Equal(t,
		"QmVyN4KhsSEgQ21WhnrNVQondEBAXYSEVSNJpTqRsyMKkg",
		cid.String(),
	)

	// CID v1
	cid, err = pi.Import(
		ctx,
		"./data",
		ImportOpts.CIDv1(),
	)
	require.NoError(t, err)
	require.Equal(t,
		"bafybeid4ij3cn74tlwbnnwscmsjxz2h5n6j7xtafbp77xkekik6e42xjk4",
		cid.String(),
	)

	// CID v0, inline
	cid, err = pi.Import(
		ctx,
		"./data",
		ImportOpts.CIDv0(),
		ImportOpts.InlineBlock(),
	)
	require.NoError(t, err)
	require.Equal(t,
		"QmfFmGQrgrn2HiW44dfJLuV1aNGnjfEFY9qJ8SWCbVDZzw",
		cid.String(),
	)

	// CID v1, inline
	cid, err = pi.Import(
		ctx,
		"./data",
		ImportOpts.CIDv1(),
		ImportOpts.InlineBlock(),
	)
	require.NoError(t, err)
	require.Equal(t,
		"bafybeihtxqlpsentx42rarm44fz4d7brmzsq5yozk3cpzmtsomo4embx6m",
		cid.String(),
	)

	// CID v1, inline, exclude bob.txt
	cid, err = pi.Import(
		ctx,
		"./data",
		ImportOpts.CIDv1(),
		ImportOpts.InlineBlock(),
		ImportOpts.Ignores("bob.txt"),
	)
	require.NoError(t, err)
	require.Equal(t,
		"bafybeigszi3podzseg3m64vvosm4sgyfrk7reukhifpuf627iwjkebochu",
		cid.String(),
	)

	// CID v0, with events
	ch := make(chan *ImportEvent, 8)
	cid, err = pi.Import(
		ctx,
		"./data",
		ImportOpts.CIDv0(),
		ImportOpts.Events(ch),
	)
	require.NoError(t, err)
	require.Equal(t,
		"QmVyN4KhsSEgQ21WhnrNVQondEBAXYSEVSNJpTqRsyMKkg",
		cid.String(),
	)
	for ev := range ch {
		fmt.Printf("%v %v\n", ev.Name, ev.CID.String())
	}
}
