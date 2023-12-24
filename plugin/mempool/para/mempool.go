package para

import (
	"github.com/assetcloud/chain/queue"
	drivers "github.com/assetcloud/chain/system/mempool"
	"github.com/assetcloud/chain/types"
)

//--------------------------------------------------------------------------------
// Module Mempool

func init() {
	drivers.Reg("para", New)
}

//New 创建price cache 结构的 mempool
func New(cfg *types.Mempool, sub []byte) queue.Module {
	return NewMempool(cfg)
}
