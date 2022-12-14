package executor

import (
	"testing"

	"github.com/assetcloud/chain/system/dapp"
	pty "github.com/assetcloud/plugin/plugin/dapp/trade/types"

	//"github.com/assetcloud/chain/common/db"
	//"github.com/assetcloud/chain/common/db/table"
	"github.com/assetcloud/chain/util"
	"github.com/stretchr/testify/assert"
)

var order1 = &pty.LocalOrder{
	AssetSymbol:       "bty",
	Owner:             "O1",
	AmountPerBoardlot: 1,
	MinBoardlot:       1,
	PricePerBoardlot:  1,
	TotalBoardlot:     10,
	TradedBoardlot:    0,
	BuyID:             "B1",
	Status:            pty.TradeOrderStatusOnBuy,
	SellID:            "",
	TxHash:            nil,
	Height:            1,
	Key:               "B1",
	BlockTime:         1,
	IsSellOrder:       false,
	AssetExec:         "coins",
	TxIndex:           dapp.HeightIndexStr(1, 1),
	IsFinished:        false,
	PriceExec:         "token",
	PriceSymbol:       "CCNY",
}

var order2 = &pty.LocalOrder{
	AssetSymbol:       "bty",
	Owner:             "O1",
	AmountPerBoardlot: 1,
	MinBoardlot:       1,
	PricePerBoardlot:  1,
	TotalBoardlot:     10,
	TradedBoardlot:    0,
	BuyID:             "B2",
	Status:            pty.TradeOrderStatusOnBuy,
	SellID:            "",
	TxHash:            nil,
	Height:            2,
	Key:               "B2",
	BlockTime:         2,
	IsSellOrder:       false,
	AssetExec:         "coins",
	TxIndex:           dapp.HeightIndexStr(2, 1),
	IsFinished:        false,
	PriceExec:         "token",
	PriceSymbol:       "CCNY",
}

// 两个fork前的数据
var order3 = &pty.LocalOrder{
	AssetSymbol:       "CCNY",
	Owner:             "O1",
	AmountPerBoardlot: 1,
	MinBoardlot:       1,
	PricePerBoardlot:  1,
	TotalBoardlot:     10,
	TradedBoardlot:    0,
	BuyID:             "B2",
	Status:            pty.TradeOrderStatusOnBuy,
	SellID:            "",
	TxHash:            nil,
	Height:            3,
	Key:               "B2",
	BlockTime:         3,
	IsSellOrder:       false,
	TxIndex:           dapp.HeightIndexStr(3, 1),
	IsFinished:        false,
}

func TestListAll(t *testing.T) {
	dir, ldb, tdb := util.CreateTestDB()
	t.Log(dir, ldb, tdb)
	odb := NewOrderTable(tdb)
	odb.Add(order1)
	odb.Add(order2)
	kvs, err := odb.Save()
	assert.Nil(t, err)
	t.Log(kvs)
	ldb.Close()
}

func TestListV2All(t *testing.T) {
	dir, ldb, tdb := util.CreateTestDB()
	t.Log(dir, ldb, tdb)
	odb := NewOrderTableV2(tdb)
	odb.Add(order1)
	odb.Add(order2)
	kvs, err := odb.Save()
	assert.Nil(t, err)
	t.Log(kvs)
	ldb.Close()
}
