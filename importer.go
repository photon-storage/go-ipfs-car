package car

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/ipfs/boxo/blockservice"
	blockstore "github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-cidutil"
	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	ipld "github.com/ipfs/go-ipld-format"
	coreiface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/core/coreiface/options"
	"github.com/ipfs/kubo/core/coreunix"
)

// DataImporter creates a new importer that imports data (can be byte slice,
// io.Reader or path from local file system) into in-memory dag service.
type DataImporter struct {
	bstore  blockstore.Blockstore
	dagServ ipld.DAGService
}

// NewDataImporter creates a new DataImporter.
func NewDataImporter() *DataImporter {
	bstore := blockstore.NewBlockstore(dssync.MutexWrap(ds.NewMapDatastore()))
	return &DataImporter{
		bstore: bstore,
		dagServ: merkledag.NewDAGService(
			blockservice.New(bstore, newNoopExchg()),
		),
	}
}

// Import imports the given input.
func (di *DataImporter) Import(
	ctx context.Context,
	input any,
	opts ...ImportOption,
) (cid.Cid, error) {
	// Build options from defaults.
	ioptions, err := buildImportOptions(opts...)
	if err != nil {
		return cid.Undef, err
	}

	// Build CID builder.
	prefix, err := merkledag.PrefixForCidVersion(ioptions.cidVersion)
	if err != nil {
		return cid.Undef, err
	}
	prefix.MhType = ioptions.mhType
	prefix.MhLength = -1

	var target files.Node
	var path string
	switch v := input.(type) {
	case string:
		var err error
		if target, err = newFsPath(
			v,
			ioptions.ignoreFile,
			ioptions.ignoreRules,
			ioptions.includeHiddenFiles,
		); err != nil {
			return cid.Undef, err
		}
		path = v
	case []byte:
		target = files.NewBytesFile(v)
	case io.Reader:
		target = files.NewReaderFile(v)
	}

	adder, err := coreunix.NewAdder(ctx, nil, nil, di.dagServ)
	if err != nil {
		return cid.Undef, err
	}
	adder.CidBuilder = prefix
	if ioptions.inline && ioptions.inlineLimit > 0 {
		adder.CidBuilder = cidutil.InlineBuilder{
			Builder: prefix,
			Limit:   ioptions.inlineLimit,
		}
	}
	adder.RawLeaves = ioptions.rawLeaves
	adder.Chunker = ioptions.chunker
	if ioptions.layout == options.TrickleLayout {
		adder.Trickle = true
	}
	adder.Pin = false
	if ioptions.out != nil {
		ch := make(chan interface{}, 8)
		adder.Progress = true
		adder.Out = ch
		defer close(ch)

		go func() {
			defer close(ioptions.out)

			_, isDir := target.(files.Directory)
			for v := range ch {
				ev, ok := v.(*coreiface.AddEvent)
				if !ok {
					continue
				}
				if ev.Path.String() == "" {
					continue
				}

				name := ev.Name
				if !isDir && path != "" {
					name = path
				} else {
					name = filepath.Join(path, ev.Name)
				}
				ioptions.out <- &ImportEvent{
					Name:  name,
					CID:   ev.Path.RootCid(),
					Bytes: ev.Bytes,
					Size:  ev.Size,
				}
			}
		}()
	}

	nd, err := adder.AddAllAndPin(ctx, target)
	if err != nil {
		return cid.Undef, err
	}

	return nd.Cid(), nil
}

func (di *DataImporter) Blockstore() blockstore.Blockstore {
	return di.bstore
}

func newFsPath(
	path string,
	ignoreFile string,
	ignoreRules []string,
	includeHiddenFiles bool,
) (files.Node, error) {
	stat, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	filter, err := files.NewFilter(ignoreFile, ignoreRules, includeHiddenFiles)
	if err != nil {
		return nil, err
	}

	return files.NewSerialFileWithFilter(path, filter, stat)
}
