// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tendermint

import (
	"fmt"

	dbm "github.com/assetcloud/chain/common/db"
	"github.com/assetcloud/chain/types"
	tmtypes "github.com/assetcloud/plugin/plugin/dapp/valnode/types"
	"github.com/golang/protobuf/proto"
)

var (
	stateKey = []byte("stateKey")
)

//ConsensusStore ...
type ConsensusStore struct {
	db dbm.DB
}

// NewConsensusStore returns a new ConsensusStore with the given DB
func NewConsensusStore() *ConsensusStore {
	db := DefaultDBProvider("state")
	db.SetCacheSize(100)
	return &ConsensusStore{
		db: db,
	}
}

// LoadStateFromStore ...
func (cs *ConsensusStore) LoadStateFromStore() *tmtypes.State {
	buf, err := cs.db.Get(stateKey)
	if err != nil {
		tendermintlog.Error("LoadStateFromStore", "err", err)
		return nil
	}
	state := &tmtypes.State{}
	err = types.Decode(buf, state)
	if err != nil {
		panic(err)
	}
	return state
}

// LoadStateHeight ...
func (cs *ConsensusStore) LoadStateHeight() int64 {
	state := cs.LoadStateFromStore()
	if state == nil {
		return int64(0)
	}
	return state.LastBlockHeight
}

// LoadSeenCommit by height
func (cs *ConsensusStore) LoadSeenCommit(height int64) *tmtypes.TendermintCommit {
	buf, err := cs.db.Get(calcSeenCommitKey(height))
	if err != nil {
		tendermintlog.Error("LoadSeenCommit", "err", err)
		return nil
	}
	commit := &tmtypes.TendermintCommit{}
	err = types.Decode(buf, commit)
	if err != nil {
		panic(err)
	}
	return commit
}

// SaveConsensusState save state and seenCommit
func (cs *ConsensusStore) SaveConsensusState(height int64, state *tmtypes.State, sc proto.Message) error {
	seenCommitBytes := types.Encode(sc)
	stateBytes := types.Encode(state)
	batch := cs.db.NewBatch(true)
	batch.Set(calcSeenCommitKey(height), seenCommitBytes)
	batch.Set(stateKey, stateBytes)
	err := batch.Write()
	if err != nil {
		tendermintlog.Error("SaveConsensusState batch.Write", "err", err)
		return err
	}
	return nil
}

func calcSeenCommitKey(height int64) []byte {
	return []byte(fmt.Sprintf("SC:%v", height))
}
