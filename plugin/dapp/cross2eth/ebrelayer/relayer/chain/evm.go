package chain

import (
	"strings"

	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/events"

	chainEvm "github.com/assetcloud/plugin/plugin/dapp/cross2eth/contracts/contracts4chain/generated"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func (relayer *Relayer4Chain) prePareSubscribeEvent() {
	var eventName string
	contractABI, err := abi.JSON(strings.NewReader(chainEvm.BridgeBankABI))
	if err != nil {
		panic(err)
	}

	eventName = events.ChainEventLogLock.String()
	relayer.bridgeBankEventLockSig = contractABI.Events[eventName].ID.Hex()
	eventName = events.ChainEventLogBurn.String()
	relayer.bridgeBankEventBurnSig = contractABI.Events[eventName].ID.Hex()
	eventName = events.ChainEventLogWithdraw.String()
	relayer.bridgeBankEventWithdrawSig = contractABI.Events[eventName].ID.Hex()

	relayer.bridgeBankAbi = contractABI

	relayerLog.Info("prePareSubscribeEvent", "bridgeBankEventLockSig", relayer.bridgeBankEventLockSig,
		"bridgeBankEventBurnSig", relayer.bridgeBankEventBurnSig, "bridgeBankEventWithdrawSig", relayer.bridgeBankEventWithdrawSig)
}
