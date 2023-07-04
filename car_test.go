package car

import (
	"bytes"
	"context"
	"io"
	"sort"
	"testing"

	"github.com/ipfs/go-cid"
	ipldcarv1 "github.com/ipld/go-car"
	ipldcar "github.com/ipld/go-car/v2"

	"github.com/photon-storage/go-common/testing/require"
)

func TestBuilder(t *testing.T) {
	ctx := context.Background()
	b := NewBuilder()

	// CID v1
	ch := make(chan *ImportEvent, 8)
	wt, err := b.Build(
		ctx,
		"./data",
		ImportOpts.CIDv1(),
		ImportOpts.Events(ch),
	)
	require.NoError(t, err)
	var builderCids []cid.Cid
	for v := range ch {
		builderCids = append(builderCids, v.CID)
	}

	buf := bytes.Buffer{}
	n, err := wt.WriteTo(&buf)
	require.NoError(t, err)
	require.Equal(t, int64(len(buf.Bytes())), n)

	v2r, err := ipldcar.NewReader(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	cids, err := v2r.Roots()
	require.NoError(t, err)
	require.Equal(t, 1, len(cids))
	require.Equal(t,
		"bafybeid4ij3cn74tlwbnnwscmsjxz2h5n6j7xtafbp77xkekik6e42xjk4",
		cids[0].String(),
	)

	sr, err := v2r.DataReader()
	require.NoError(t, err)
	v1r, err := ipldcarv1.NewCarReader(sr)
	require.NoError(t, err)

	// Compare block CIDs from v1 reader match CIDs reported by builder.
	var readerCids []cid.Cid
	for {
		b, err := v1r.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		readerCids = append(readerCids, b.Cid())
	}
	sort.Slice(builderCids, func(i, j int) bool {
		return builderCids[i].String() < builderCids[j].String()
	})
	sort.Slice(readerCids, func(i, j int) bool {
		return readerCids[i].String() < readerCids[j].String()
	})
	require.DeepEqual(t, builderCids, readerCids)
}
