/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package rpc_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/assetcloud/chain/rpc/jsonclient"
	rpctypes "github.com/assetcloud/chain/rpc/types"
	"github.com/assetcloud/chain/types"
	"github.com/assetcloud/chain/util/testnode"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	commonlog "github.com/assetcloud/chain/common/log"
	_ "github.com/assetcloud/chain/system"
	_ "github.com/assetcloud/plugin/plugin"
	oty "github.com/assetcloud/plugin/plugin/dapp/oracle/types"
)

func init() {
	commonlog.SetLogLevel("error")
}

func TestJRPCChannel(t *testing.T) {
	// 启动RPCmocker
	mocker := testnode.New("--notset--", nil)
	cfg := mocker.GetAPI().GetConfig()
	defer func() {
		mocker.Close()
	}()
	mocker.Listen()

	jrpcClient := mocker.GetJSONC()
	assert.NotNil(t, jrpcClient)

	testCases := []struct {
		fn func(*testing.T, *types.ChainConfig, *jsonclient.JSONClient) error
	}{
		{fn: testPublishEventRawCmd},
		{fn: testAbortEventRawTxCmd},
		{fn: testPrePublishResultRawTxCmd},
		{fn: testAbortPrePubResultRawTxCmd},
		{fn: testPublishResultRawTxCmd},
	}
	for index, testCase := range testCases {
		err := testCase.fn(t, cfg, jrpcClient)
		if err == nil {
			continue
		}
		assert.NotEqualf(t, err, types.ErrActionNotSupport, "test index %d", index)
		if strings.Contains(err.Error(), "rpc: can't find") {
			assert.FailNowf(t, err.Error(), "test index %d", index)
		}
	}

	testCases1 := []struct {
		fn func(*testing.T, *jsonclient.JSONClient) error
	}{
		{fn: testQueryOracleListByIDsRawTxCmd},
		{fn: testQueryEventIDByAddrAndStatusRawTxCmd},
		{fn: testQueryEventIDByTypeAndStatusRawTxCmd},
		{fn: testQueryEventIDByStatusRawTxCmd},
	}
	result := []error{
		oty.ErrParamNeedIDs,
		oty.ErrParamStatusInvalid,
		types.ErrNotFound,
		types.ErrNotFound,
	}
	for index, testCase := range testCases1 {
		err := testCase.fn(t, jrpcClient)
		assert.Equal(t, result[index], err, fmt.Sprint(index))
	}
}

func testPublishEventRawCmd(t *testing.T, cfg *types.ChainConfig, jrpc *jsonclient.JSONClient) error {
	timeStr := "2019-01-21 15:30:00"
	layout := "2006-01-02 15:04:05"
	ti, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Printf("time error:%v\n", err.Error())
		return errors.Errorf("time error:%v\n", err.Error())
	}
	payload := &oty.EventPublish{
		Type:         "football",
		SubType:      "Premier League",
		Time:         ti.Unix(),
		Content:      "{\"team1\":\"ChelSea\", \"team2\":\"Manchester\",\"resultType\":\"score\"}",
		Introduction: "guess the sore result of football game between ChelSea and Manchester in 2019-01-21 14:00:00",
	}
	params := &rpctypes.CreateTxIn{
		Execer:     cfg.ExecName(oty.OracleX),
		ActionName: oty.CreateEventPublishTx,
		Payload:    types.MustPBToJSON(payload),
	}
	var res string
	return jrpc.Call("Chain.CreateTransaction", params, &res)
}

func testAbortEventRawTxCmd(t *testing.T, cfg *types.ChainConfig, jrpc *jsonclient.JSONClient) error {
	payload := &oty.EventAbort{EventID: "123"}
	params := &rpctypes.CreateTxIn{
		Execer:     cfg.ExecName(oty.OracleX),
		ActionName: oty.CreateAbortEventPublishTx,
		Payload:    types.MustPBToJSON(payload),
	}
	var res string
	return jrpc.Call("Chain.CreateTransaction", params, &res)
}

func testPrePublishResultRawTxCmd(t *testing.T, cfg *types.ChainConfig, jrpc *jsonclient.JSONClient) error {
	payload := &oty.ResultPrePublish{
		EventID: "123",
		Source:  "新浪体育",
		Result:  "{\"team1\":3, \"team2\":2}",
	}
	params := &rpctypes.CreateTxIn{
		Execer:     cfg.ExecName(oty.OracleX),
		ActionName: oty.CreatePrePublishResultTx,
		Payload:    types.MustPBToJSON(payload),
	}
	var res string
	return jrpc.Call("Chain.CreateTransaction", params, &res)
}

func testAbortPrePubResultRawTxCmd(t *testing.T, cfg *types.ChainConfig, jrpc *jsonclient.JSONClient) error {
	payload := &oty.EventAbort{EventID: "123"}
	params := &rpctypes.CreateTxIn{
		Execer:     cfg.ExecName(oty.OracleX),
		ActionName: oty.CreateAbortResultPrePublishTx,
		Payload:    types.MustPBToJSON(payload),
	}
	var res string
	return jrpc.Call("Chain.CreateTransaction", params, &res)
}

func testPublishResultRawTxCmd(t *testing.T, cfg *types.ChainConfig, jrpc *jsonclient.JSONClient) error {
	payload := &oty.ResultPrePublish{
		EventID: "123",
		Source:  "新浪体育",
		Result:  "{\"team1\":3, \"team2\":2}",
	}
	params := &rpctypes.CreateTxIn{
		Execer:     cfg.ExecName(oty.OracleX),
		ActionName: oty.CreateResultPublishTx,
		Payload:    types.MustPBToJSON(payload),
	}
	var res string
	return jrpc.Call("Chain.CreateTransaction", params, &res)
}

func testQueryOracleListByIDsRawTxCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &oty.QueryOracleInfos{}
	params.Execer = oty.OracleX
	params.FuncName = oty.FuncNameQueryOracleListByIDs
	params.Payload = types.MustPBToJSON(req)
	rep = &oty.ReplyOracleStatusList{}
	return jrpc.Call("Chain.Query", params, rep)
}

func testQueryEventIDByAddrAndStatusRawTxCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &oty.QueryEventID{}
	params.Execer = oty.OracleX
	params.FuncName = oty.FuncNameQueryEventIDByAddrAndStatus
	params.Payload = types.MustPBToJSON(req)
	rep = &oty.ReplyEventIDs{}
	return jrpc.Call("Chain.Query", params, rep)
}

func testQueryEventIDByTypeAndStatusRawTxCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &oty.QueryEventID{
		Type:   "football",
		Status: 1,
		Addr:   "",
	}
	params.Execer = oty.OracleX
	params.FuncName = oty.FuncNameQueryEventIDByTypeAndStatus
	params.Payload = types.MustPBToJSON(req)
	rep = &oty.ReplyEventIDs{}
	return jrpc.Call("Chain.Query", params, rep)
}

func testQueryEventIDByStatusRawTxCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &oty.QueryEventID{
		Status: 1,
		Type:   "",
		Addr:   "",
	}
	params.Execer = oty.OracleX
	params.FuncName = oty.FuncNameQueryEventIDByStatus
	params.Payload = types.MustPBToJSON(req)
	rep = &oty.ReplyEventIDs{}
	return jrpc.Call("Chain.Query", params, rep)
}
