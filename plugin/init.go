package plugin

import (
	_ "github.com/assetcloud/plugin/plugin/consensus/init" //consensus init
	_ "github.com/assetcloud/plugin/plugin/crypto/init"    //crypto init
	_ "github.com/assetcloud/plugin/plugin/dapp/init"      //dapp init
	_ "github.com/assetcloud/plugin/plugin/mempool/init"   //mempool init
	_ "github.com/assetcloud/plugin/plugin/p2p/init"       //p2p init
	_ "github.com/assetcloud/plugin/plugin/store/init"     //store init
)
