package executor

import (
	"strings"

	pt "github.com/assetcloud/plugin/plugin/dapp/paracross/types"
	"github.com/assetcloud/chain/common"
	dbm "github.com/assetcloud/chain/common/db"
	"github.com/assetcloud/chain/types"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func makeParaSupervisionNodeGroupReceipt(title string, prev, current *types.ConfigItem) *types.Receipt {
	key := calcParaSupervisionNodeGroupAddrsKey(title)
	log := &types.ReceiptConfig{Prev: prev, Current: current}
	return &types.Receipt{
		Ty: types.ExecOk,
		KV: []*types.KeyValue{
			{Key: key, Value: types.Encode(current)},
		},
		Logs: []*types.ReceiptLog{
			{
				Ty:  pt.TyLogParaSupervisionNodeGroupAddrsUpdate,
				Log: types.Encode(log),
			},
		},
	}
}

func makeSupervisionNodeConfigReceipt(fromAddr string, config *pt.ParaNodeGroupConfig, prev, current *pt.ParaNodeGroupStatus) *types.Receipt {
	log := &pt.ReceiptParaNodeGroupConfig{
		Addr:    fromAddr,
		Config:  config,
		Prev:    prev,
		Current: current,
	}
	return &types.Receipt{
		Ty: types.ExecOk,
		KV: []*types.KeyValue{
			{Key: []byte(current.Id), Value: types.Encode(current)},
		},
		Logs: []*types.ReceiptLog{
			{
				Ty:  pt.TyLogParaSupervisionNodeConfig,
				Log: types.Encode(log),
			},
		},
	}
}

func makeParaSupervisionNodeStatusReceipt(fromAddr string, prev, current *pt.ParaNodeAddrIdStatus) *types.Receipt {
	key := calcParaNodeAddrKey(current.Title, current.Addr)
	log := &pt.ReceiptParaNodeAddrStatUpdate{
		FromAddr: fromAddr,
		Prev:     prev,
		Current:  current,
	}
	return &types.Receipt{
		Ty: types.ExecOk,
		KV: []*types.KeyValue{
			{Key: key, Value: types.Encode(current)},
		},
		Logs: []*types.ReceiptLog{
			{
				Ty:  pt.TyLogParaSupervisionNodeStatusUpdate,
				Log: types.Encode(log),
			},
		},
	}
}

func getSupervisionNodeID(db dbm.KV, id string) (*pt.ParaNodeGroupStatus, error) {
	val, err := getDb(db, []byte(id))
	if err != nil {
		return nil, err
	}

	var status pt.ParaNodeGroupStatus
	err = types.Decode(val, &status)
	return &status, err
}

func (a *action) updateSupervisionNodeGroup(title, addr string, add bool) (*types.Receipt, error) {
	var item types.ConfigItem
	key := calcParaSupervisionNodeGroupAddrsKey(title)

	value, err := a.db.Get(key)
	if err != nil {
		return nil, err
	}

	if value != nil {
		err = types.Decode(value, &item)
		if err != nil {
			clog.Error("updateSupervisionNodeGroup", "decode db key", key)
			return nil, err
		}
	}

	copyValue := *item.GetArr()
	copyItem := item
	copyItem.Value = &types.ConfigItem_Arr{Arr: &copyValue}

	receipt := &types.Receipt{Ty: types.ExecOk}
	item.Addr = addr
	if add {
		item.GetArr().Value = append(item.GetArr().Value, addr)
		clog.Info("updateSupervisionNodeGroup add", "addr", addr)
	} else {
		item.GetArr().Value = make([]string, 0)
		for _, value := range copyItem.GetArr().Value {
			if value != addr {
				item.GetArr().Value = append(item.GetArr().Value, value)
			}
		}
		clog.Info("updateSupervisionNodeGroup delete", "addr", addr)
	}
	err = a.db.Set(key, types.Encode(&item))
	if err != nil {
		return nil, errors.Wrapf(err, "updateNodeGroup set dbkey=%s", key)
	}
	r := makeParaSupervisionNodeGroupReceipt(title, &copyItem, &item)
	receipt = mergeReceipt(receipt, r)
	return receipt, nil
}

func (a *action) checkValidSupervisionNode(config *pt.ParaNodeGroupConfig) (bool, error) {
	key := calcParaSupervisionNodeGroupAddrsKey(config.Title)
	nodes, _, err := getNodes(a.db, key)
	if err != nil && !(isNotFound(err) || errors.Cause(err) == pt.ErrTitleNotExist) {
		return false, errors.Wrapf(err, "getNodes for title:%s", config.Title)
	}

	if validNode(config.Addrs, nodes) {
		return true, nil
	}
	return false, nil
}

func (a *action) checkSupervisionNodeGroupExist(title string) (bool, error) {
	key := calcParaSupervisionNodeGroupAddrsKey(title)
	value, err := a.db.Get(key)
	if err != nil && !isNotFound(err) {
		return false, err
	}

	if value != nil {
		var item types.ConfigItem
		err = types.Decode(value, &item)
		if err != nil {
			clog.Error("updateSupervisionNodeGroup", "decode db key", key)
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (a *action) supervisionNodeGroupCreate(title, targetAddrs string) (*types.Receipt, error) {
	var item types.ConfigItem
	key := calcParaSupervisionNodeGroupAddrsKey(title)
	item.Key = string(key)
	emptyValue := &types.ArrayConfig{Value: make([]string, 0)}
	arr := types.ConfigItem_Arr{Arr: emptyValue}
	item.Value = &arr
	item.GetArr().Value = append(item.GetArr().Value, targetAddrs)
	item.Addr = a.fromaddr

	receipt := makeParaSupervisionNodeGroupReceipt(title, nil, &item)
	return receipt, nil
}

//??????propasal id ???quit id?????????quit id???????????????addr???proposal id???coinfrozen?????????????????????????????????addr????????????????????????
func (a *action) updateSupervisionNodeAddrStatus(stat *pt.ParaNodeGroupStatus) (*types.Receipt, error) {
	addrStat, err := getNodeAddr(a.db, stat.Title, stat.TargetAddrs)
	if err != nil {
		if !isNotFound(err) {
			return nil, errors.Wrapf(err, "nodeAddr:%s get error", stat.TargetAddrs)
		}
		addrStat = &pt.ParaNodeAddrIdStatus{}
		addrStat.Title = stat.Title
		addrStat.Addr = stat.TargetAddrs
		addrStat.BlsPubKey = stat.BlsPubKeys
		addrStat.Status = pt.ParacrossSupervisionNodeApprove
		addrStat.ProposalId = stat.Id
		addrStat.QuitId = ""
		return makeParaSupervisionNodeStatusReceipt(a.fromaddr, nil, addrStat), nil
	}

	preStat := *addrStat
	if stat.Status == pt.ParacrossSupervisionNodeQuit {
		proposalStat, err := getSupervisionNodeID(a.db, addrStat.ProposalId)
		if err != nil {
			return nil, errors.Wrapf(err, "nodeAddr:%s quiting wrong proposeid:%s", stat.TargetAddrs, addrStat.ProposalId)
		}

		addrStat.Status = stat.Status
		addrStat.QuitId = stat.Id
		receipt := makeParaSupervisionNodeStatusReceipt(a.fromaddr, &preStat, addrStat)

		cfg := a.api.GetConfig()
		if !cfg.IsPara() {
			r, err := a.nodeGroupCoinsActive(proposalStat.FromAddr, proposalStat.CoinsFrozen, 1)
			if err != nil {
				return nil, err
			}
			receipt = mergeReceipt(receipt, r)
		}
		return receipt, nil
	} else {
		addrStat.Status = stat.Status
		addrStat.ProposalId = stat.Id
		addrStat.QuitId = ""
		return makeParaSupervisionNodeStatusReceipt(a.fromaddr, &preStat, addrStat), nil
	}
}

func (a *action) supervisionNodeApply(config *pt.ParaNodeGroupConfig) (*types.Receipt, error) {
	// ???????????????????????? ???????????????????????? ??????????????????????????????
	if strings.Contains(config.Addrs, ",") {
		return nil, errors.Wrapf(types.ErrInvalidParam, "not support multi addr currently,addrs=%s", config.Addrs)
	}
	nodeCfg := &pt.ParaNodeAddrConfig{Title: config.Title, Addr: config.Addrs}
	addrExist, err := a.checkValidNode(nodeCfg)
	if err != nil {
		return nil, err
	}

	// ???????????????????????????
	if addrExist {
		return nil, errors.Wrapf(pt.ErrParaNodeAddrExisted, "supervisionNodeGroup Apply Addr existed:%s in super group", config.Addrs)
	}

	// ?????? node ??????????????????
	addrExist, err = a.checkValidSupervisionNode(config)
	if err != nil {
		return nil, err
	}
	if addrExist {
		return nil, errors.Wrapf(pt.ErrParaSupervisionNodeAddrExisted, "supervisionNodeGroup Apply Addr existed:%s", config.Addrs)
	}

	// ????????????????????????
	receipt := &types.Receipt{Ty: types.ExecOk}
	cfg := a.api.GetConfig()
	if !cfg.IsPara() {
		r, err := a.nodeGroupCoinsFrozen(a.fromaddr, config.CoinsFrozen, 1)
		if err != nil {
			return nil, err
		}
		receipt = mergeReceipt(receipt, r)
	}

	stat := &pt.ParaNodeGroupStatus{
		Id:          calcParaSupervisionNodeIDKey(config.Title, common.ToHex(a.txhash)),
		Status:      pt.ParacrossSupervisionNodeApply,
		Title:       config.Title,
		TargetAddrs: config.Addrs,
		BlsPubKeys:  config.BlsPubKeys,
		CoinsFrozen: config.CoinsFrozen,
		FromAddr:    a.fromaddr,
		Height:      a.height,
	}
	r := makeSupervisionNodeConfigReceipt(a.fromaddr, config, nil, stat)
	receipt = mergeReceipt(receipt, r)

	return receipt, nil
}

func (a *action) supervisionNodeApprove(config *pt.ParaNodeGroupConfig) (*types.Receipt, error) {
	cfg := a.api.GetConfig()

	apply, err := getSupervisionNodeID(a.db, calcParaSupervisionNodeIDKey(config.Title, config.Id))
	if err != nil {
		return nil, err
	}
	if config.Title != apply.Title {
		return nil, errors.Wrapf(pt.ErrNodeNotForTheTitle, "config title:%s,id title:%s", config.Title, apply.Title)
	}
	if apply.CoinsFrozen < config.CoinsFrozen {
		return nil, errors.Wrapf(pt.ErrParaNodeGroupFrozenCoinsNotEnough, "id not enough coins apply:%d,config:%d", apply.CoinsFrozen, config.CoinsFrozen)
	}

	//????????????????????? ?????????????????????????????????????????????????????????????????????????????????
	if !cfg.IsPara() {
		err := a.checkApproveOp(config, apply)
		if err != nil {
			return nil, err
		}
	}

	// ???????????????????????????????????????
	exist, err := a.checkSupervisionNodeGroupExist(config.Title)
	if err != nil {
		return nil, err
	}

	receipt := &types.Receipt{Ty: types.ExecOk}
	if !exist {
		// ????????????????????????
		r, err := a.supervisionNodeGroupCreate(apply.Title, apply.TargetAddrs)
		if err != nil {
			return nil, errors.Wrapf(err, "nodegroup create:title:%s,addrs:%s", config.Title, apply.TargetAddrs)
		}
		receipt = mergeReceipt(receipt, r)
	} else {
		// ???????????????????????????
		r, err := a.updateSupervisionNodeGroup(config.Title, apply.TargetAddrs, true)
		if err != nil {
			return nil, err
		}
		receipt = mergeReceipt(receipt, r)
	}

	copyStat := proto.Clone(apply).(*pt.ParaNodeGroupStatus)
	apply.Status = pt.ParacrossSupervisionNodeApprove
	apply.Height = a.height

	r, err := a.updateSupervisionNodeAddrStatus(apply)
	if err != nil {
		return nil, err
	}
	receipt = mergeReceipt(receipt, r)

	r = makeSupervisionNodeConfigReceipt(a.fromaddr, config, copyStat, apply)
	receipt = mergeReceipt(receipt, r)
	return receipt, nil
}

func (a *action) supervisionNodeQuit(config *pt.ParaNodeGroupConfig) (*types.Receipt, error) {
	addrExist, err := a.checkValidSupervisionNode(config)
	if err != nil {
		return nil, err
	}
	if !addrExist {
		return nil, errors.Wrapf(pt.ErrParaSupervisionNodeAddrNotExisted, "nodeAddr not existed:%s", config.Addrs)
	}

	status, err := getNodeAddr(a.db, config.Title, config.Addrs)
	if err != nil {
		return nil, errors.Wrapf(err, "nodeAddr:%s get error", config.Addrs)
	}
	if status.Status != pt.ParacrossSupervisionNodeApprove {
		return nil, errors.Wrapf(pt.ErrParaSupervisionNodeAddrNotExisted, "nodeAddr:%s status:%d", config.Addrs, status.Status)
	}

	cfg := a.api.GetConfig()
	stat := &pt.ParaNodeGroupStatus{
		Id:          calcParaSupervisionNodeIDKey(config.Title, common.ToHex(a.txhash)),
		Status:      pt.ParacrossSupervisionNodeQuit,
		Title:       config.Title,
		TargetAddrs: config.Addrs,
		FromAddr:    a.fromaddr,
		Height:      a.height,
	}

	//????????????????????????????????????????????????
	if a.fromaddr != status.Addr && !cfg.IsPara() && !isSuperManager(cfg, a.fromaddr) {
		return nil, errors.Wrapf(types.ErrNotAllow, "id create by:%s,not by:%s", status.Addr, a.fromaddr)
	}

	if config.Title != status.Title {
		return nil, errors.Wrapf(pt.ErrNodeNotForTheTitle, "config title:%s,id title:%s", config.Title, status.Title)
	}

	receipt := &types.Receipt{Ty: types.ExecOk}

	// node quit????????????committx??????2/3?????????????????????commitDone ????????????
	r, err := a.loopCommitTxDone(config.Title)
	if err != nil {
		clog.Error("updateSupervisionNodeGroup.loopCommitTxDone", "title", config.Title, "err", err.Error())
	}
	receipt = mergeReceipt(receipt, r)

	r, err = a.updateSupervisionNodeGroup(config.Title, stat.TargetAddrs, false)
	if err != nil {
		return nil, err
	}
	receipt = mergeReceipt(receipt, r)

	r, err = a.updateSupervisionNodeAddrStatus(stat)
	if err != nil {
		return nil, err
	}
	receipt = mergeReceipt(receipt, r)

	r = makeSupervisionNodeConfigReceipt(a.fromaddr, config, nil, stat)
	receipt = mergeReceipt(receipt, r)

	return receipt, nil
}

func (a *action) supervisionNodeCancel(config *pt.ParaNodeGroupConfig) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	status, err := getSupervisionNodeID(a.db, calcParaSupervisionNodeIDKey(config.Title, config.Id))
	if err != nil {
		return nil, errors.Wrapf(err, "getSupervisionNodeID=%s", config.Id)
	}

	//?????????????????????????????????
	if a.fromaddr != status.FromAddr {
		return nil, errors.Wrapf(types.ErrNotAllow, "id create by:%s,not by:%s", status.FromAddr, a.fromaddr)
	}
	if config.Title != status.Title {
		return nil, errors.Wrapf(pt.ErrNodeNotForTheTitle, "config title:%s,id title:%s", config.Title, status.Title)
	}
	if status.Status != pt.ParacrossSupervisionNodeApply {
		return nil, errors.Wrapf(pt.ErrParaNodeOpStatusWrong, "config id:%s,status:%d", config.Id, status.Status)
	}

	receipt := &types.Receipt{Ty: types.ExecOk}
	//main chain
	if !cfg.IsPara() {
		r, err := a.nodeGroupCoinsActive(status.FromAddr, status.CoinsFrozen, 1)
		if err != nil {
			return nil, err
		}
		receipt = mergeReceipt(receipt, r)
	}

	copyStat := proto.Clone(status).(*pt.ParaNodeGroupStatus)
	status.Status = pt.ParacrossSupervisionNodeCancel
	status.Height = a.height

	r := makeSupervisionNodeConfigReceipt(a.fromaddr, config, copyStat, status)
	receipt = mergeReceipt(receipt, r)

	return receipt, nil
}

func (a *action) supervisionNodeModify(config *pt.ParaNodeGroupConfig) (*types.Receipt, error) {
	addrStat, err := getNodeAddr(a.db, config.Title, config.Addrs)
	if err != nil {
		return nil, errors.Wrapf(err, "nodeAddr:%s get error", config.Addrs)
	}

	// ?????????????????????
	if a.fromaddr != config.Addrs {
		return nil, errors.Wrapf(types.ErrNotAllow, "addr create by:%s,not by:%s", config.Addrs, a.fromaddr)
	}

	preStat := *addrStat
	addrStat.BlsPubKey = config.BlsPubKeys

	return makeParaSupervisionNodeStatusReceipt(a.fromaddr, &preStat, addrStat), nil
}

func (a *action) SupervisionNodeConfig(config *pt.ParaNodeGroupConfig) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	if !validTitle(cfg, config.Title) {
		return nil, pt.ErrInvalidTitle
	}
	if !types.IsParaExecName(string(a.tx.Execer)) {
		return nil, errors.Wrapf(types.ErrInvalidParam, "exec=%s,should prefix with user.p.", string(a.tx.Execer))
	}
	if (config.Op == pt.ParacrossSupervisionNodeApprove || config.Op == pt.ParacrossSupervisionNodeCancel) && config.Id == "" {
		return nil, types.ErrInvalidParam
	}

	if config.Op == pt.ParacrossSupervisionNodeApply {
		return a.supervisionNodeApply(config)
	} else if config.Op == pt.ParacrossSupervisionNodeApprove {
		return a.supervisionNodeApprove(config)
	} else if config.Op == pt.ParacrossSupervisionNodeQuit {
		// ?????? group
		return a.supervisionNodeQuit(config)
	} else if config.Op == pt.ParacrossSupervisionNodeCancel {
		// ????????????????????????
		return a.supervisionNodeCancel(config)
	} else if config.Op == pt.ParacrossSupervisionNodeModify {
		return a.supervisionNodeModify(config)
	}

	return nil, pt.ErrParaUnSupportNodeOper
}
