package init

import (
	_ "github.com/assetcloud/plugin/plugin/consensus/dpos" //auto gen
	_ "github.com/assetcloud/plugin/plugin/consensus/para" //auto gen
	_ "github.com/assetcloud/plugin/plugin/consensus/pbft" //auto gen
	_ "github.com/assetcloud/plugin/plugin/consensus/qbft" //auto gen
	_ "github.com/assetcloud/plugin/plugin/consensus/raft" //auto gen
	_ "github.com/assetcloud/plugin/plugin/consensus/rollup"
	_ "github.com/assetcloud/plugin/plugin/consensus/tendermint" //auto gen
	_ "github.com/assetcloud/plugin/plugin/consensus/ticket"     //auto gen
)
