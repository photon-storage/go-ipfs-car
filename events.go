package car

import "github.com/ipfs/go-cid"

type ImportEvent struct {
	Name  string
	CID   cid.Cid
	Bytes int64
	Size  string
}
