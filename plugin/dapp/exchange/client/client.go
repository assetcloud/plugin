package client

import (
	"github.com/assetcloud/chain/types"
	"github.com/golang/protobuf/proto"
)

type Cli interface {
	Query(fn string, msg proto.Message) ([]byte, error)
	Send(tx *types.Transaction, hexKey string) ([]*types.ReceiptLog, error)
}
