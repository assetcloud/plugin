package executor

import (
	"github.com/assetcloud/chain/types"
	storagetypes "github.com/assetcloud/plugin/plugin/dapp/storage/types"
)

//从statedb 读取原始数据
func (s *storage) Query_QueryStorage(in *storagetypes.QueryStorage) (types.Message, error) {
	return QueryStorage(s.GetStateDB(), s.GetLocalDB(), in.TxHash)
}

//通过状态查询ids
func (s *storage) Query_BatchQueryStorage(in *storagetypes.BatchQueryStorage) (types.Message, error) {
	return BatchQueryStorage(s.GetStateDB(), s.GetLocalDB(), in)
}
