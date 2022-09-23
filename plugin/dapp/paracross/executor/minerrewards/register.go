package minerrewards

import (
	"fmt"

	pt "github.com/assetcloud/plugin/plugin/dapp/paracross/types"
	"github.com/assetcloud/chain/types"
)

type RewardPolicy interface {
	GetConfigReward(cfg *types.Chain33Config, height int64) (int64, int64, int64)
	RewardMiners(cfg *types.Chain33Config, coinReward int64, miners []string, height int64) ([]*pt.ParaMinerReward, int64)
}

var MinerRewards = make(map[string]RewardPolicy)

func register(ty string, policy RewardPolicy) {
	if _, ok := MinerRewards[ty]; ok {
		panic(fmt.Sprintf("paracross minerreward ty=%s registered", ty))
	}
	MinerRewards[ty] = policy
}
