package car

import (
	"context"
	"errors"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
)

var (
	ErrNoopExchgNotFound = errors.New("not found in noop exchange")
)

var _ exchange.Interface = (*noopExchg)(nil)

type noopExchg struct {
}

func newNoopExchg() *noopExchg {
	return &noopExchg{}
}

func (e *noopExchg) GetBlock(context.Context, cid.Cid) (blocks.Block, error) {
	return nil, ErrNoopExchgNotFound
}

func (e *noopExchg) GetBlocks(
	_ context.Context,
	_ []cid.Cid,
) (<-chan blocks.Block, error) {
	return nil, ErrNoopExchgNotFound
}

func (e *noopExchg) NotifyNewBlocks(
	_ context.Context,
	_ ...blocks.Block,
) error {
	return nil
}

func (e *noopExchg) Close() error {
	return nil
}
