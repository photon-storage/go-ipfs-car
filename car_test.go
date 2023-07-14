package car

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
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

	v1car, err := b.Buildv1(
		ctx,
		"./data",
		ImportOpts.CIDv1(),
	)
	require.NoError(t, err)

	v1buf := bytes.Buffer{}
	require.NoError(t, v1car.Write(&v1buf))
	require.Equal(t,
		"bafybeid4ij3cn74tlwbnnwscmsjxz2h5n6j7xtafbp77xkekik6e42xjk4",
		v1car.Root().String(),
	)

	// CID v1
	ch := make(chan *ImportEvent, 8)
	v2car, err := b.Buildv2(
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

	v2buf := bytes.Buffer{}
	n, err := v2car.WriteTo(&v2buf)
	require.NoError(t, err)
	require.Equal(t, int64(len(v2buf.Bytes())), n)
	require.Equal(t,
		"bafybeid4ij3cn74tlwbnnwscmsjxz2h5n6j7xtafbp77xkekik6e42xjk4",
		v2car.Root().String(),
	)

	v2r, err := ipldcar.NewReader(bytes.NewReader(v2buf.Bytes()))
	require.NoError(t, err)
	cids, err := v2r.Roots()
	require.NoError(t, err)
	require.Equal(t, 1, len(cids))
	require.Equal(t,
		"bafybeid4ij3cn74tlwbnnwscmsjxz2h5n6j7xtafbp77xkekik6e42xjk4",
		cids[0].String(),
	)

	dr, err := v2r.DataReader()
	require.NoError(t, err)
	v1data, err := ioutil.ReadAll(dr)
	require.NoError(t, err)
	require.DeepEqual(t, v1buf.Bytes(), v1data)

	v1r, err := ipldcarv1.NewCarReader(&v1buf)
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
