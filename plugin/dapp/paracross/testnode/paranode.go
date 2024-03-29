package testnode

import (
	"github.com/assetcloud/chain/types"
	"github.com/assetcloud/chain/util/testnode"
	_ "github.com/assetcloud/plugin/plugin/mempool/init"
)

/*
1. solo 模式，后台启动一个 主节点
2. 启动一个平行链节点：注意，这个要测试的话，会依赖平行链插件
*/

//ParaNode 平行链节点由两个节点组成
type ParaNode struct {
	Main *testnode.ChainMock
	Para *testnode.ChainMock
}

//NewParaNode 创建一个平行链节点
func NewParaNode(main *testnode.ChainMock, para *testnode.ChainMock) *ParaNode {
	if main == nil {
		main = testnode.New("", nil)
		main.Listen()
	}
	if para == nil {
		cfg := types.NewChainConfig(DefaultConfig)
		cfg.GetModuleConfig().RPC.ParaChain.MainChainGrpcAddr = main.GetCfg().RPC.GrpcBindAddr
		para = testnode.NewWithConfig(cfg, nil)
		para.Listen()
	}
	return &ParaNode{Main: main, Para: para}
}

//Close 关闭系统
func (node *ParaNode) Close() {
	node.Para.Close()
	node.Main.Close()
}
