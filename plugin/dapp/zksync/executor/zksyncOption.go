package executor

import (
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/log/log15"
	"math/big"
	"strconv"
	"strings"

	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/client"
	"github.com/33cn/chain33/common/address"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	zt "github.com/33cn/plugin/plugin/dapp/zksync/types"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/pkg/errors"
)

var (
	zklog = log15.New("module", "exec.zksync")
)

// Action action struct
type Action struct {
	statedb   dbm.KV
	txhash    []byte
	fromaddr  string
	blocktime int64
	height    int64
	execaddr  string
	localDB   dbm.KVDB
	index     int
	api       client.QueueProtocolAPI
}

//NewAction ...
func NewAction(z *zksync, tx *types.Transaction, index int) *Action {
	hash := tx.Hash()
	fromaddr := tx.From()
	return &Action{
		statedb:   z.GetStateDB(),
		txhash:    hash,
		fromaddr:  fromaddr,
		blocktime: z.GetBlockTime(),
		height:    z.GetHeight(),
		execaddr:  dapp.ExecAddress(string(tx.Execer)),
		localDB:   z.GetLocalDB(),
		index:     index,
		api:       z.GetAPI(),
	}
}

//GetIndex get index
func (a *Action) GetIndex() int64 {
	return a.height*types.MaxTxsPerBlock + int64(a.index)
}

func (a *Action) Deposit(payload *zt.ZkDeposit) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	var localKvs []*types.KeyValue
	var err error

	err = checkParam(payload.Amount)
	if err != nil {
		return nil, errors.Wrapf(err, "checkParam")
	}

	zklog.Info("start zksync deposit", "eth", payload.EthAddress, "chain33", payload.Chain33Addr)
	//只有管理员能操作
	cfg := a.api.GetConfig()
	if !isSuperManager(cfg, a.fromaddr) && !isVerifier(a.statedb, a.fromaddr) {
		return nil, errors.Wrapf(types.ErrNotAllow, "from addr is not manager")
	}

	info, err := generateTreeUpdateInfo(a.statedb)
	if err != nil {
		return nil, errors.Wrapf(err, "db.generateTreeUpdateInfo")
	}

	leaf, err := GetLeafByChain33AndEthAddress(a.statedb, payload.GetChain33Addr(), payload.GetEthAddress(), info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetLeafByChain33AndEthAddress")
	}

	tree, err := getAccountTree(a.statedb, info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.getAccountTree")
	}
	zklog.Info("zksync deposit", "tree", tree)

	operationInfo := &zt.OperationInfo{
		BlockHeight: uint64(a.height),
		TxIndex:     uint32(a.index),
		TxType:      zt.TyDepositAction,
		TokenID:     payload.TokenId,
		Amount:      payload.Amount,
		SigData:     payload.Signature,
	}

	//leaf不存在就添加
	if leaf == nil {
		zklog.Info("zksync deposit add leaf")
		operationInfo.AccountID = tree.GetTotalIndex() + 1
		//添加之前先计算证明
		receipt, err := calProof(a.statedb, info, operationInfo.AccountID, payload.TokenId)
		if err != nil {
			return nil, errors.Wrapf(err, "calProof")
		}

		before := getBranchByReceipt(receipt, operationInfo, payload.EthAddress, payload.Chain33Addr, nil, "0")

		kvs, localKvs, err = AddNewLeaf(a.statedb, a.localDB, info, payload.GetEthAddress(), payload.GetTokenId(), payload.GetAmount(), payload.GetChain33Addr())
		if err != nil {
			return nil, errors.Wrapf(err, "db.AddNewLeaf")
		}
		receipt, err = calProof(a.statedb, info, operationInfo.AccountID, payload.TokenId)
		if err != nil {
			return nil, errors.Wrapf(err, "calProof")
		}

		after := getBranchByReceipt(receipt, operationInfo, payload.EthAddress, payload.Chain33Addr, nil, receipt.Token.Balance)
		rootHash, err := common.FromHex(receipt.TreeProof.RootHash)
		if err != nil {
			return nil, errors.Wrapf(err, "hex.DecodeString")
		}
		kv := &types.KeyValue{
			Key:   getHeightKey(a.height),
			Value: rootHash,
		}
		kvs = append(kvs, kv)

		branch := &zt.OperationPairBranch{
			Before: before,
			After:  after,
		}
		operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)
		zklog := &zt.ZkReceiptLog{
			OperationInfo: operationInfo,
			LocalKvs:      localKvs,
		}
		receiptLog := &types.ReceiptLog{Ty: zt.TyDepositLog, Log: types.Encode(zklog)}
		logs = append(logs, receiptLog)
	} else {
		operationInfo.AccountID = leaf.GetAccountId()

		receipt, err := calProof(a.statedb, info, leaf.AccountId, payload.TokenId)
		if err != nil {
			return nil, errors.Wrapf(err, "calProof")
		}

		var balance string
		if receipt.Token == nil {
			balance = "0"
		} else {
			balance = receipt.Token.Balance
		}
		before := getBranchByReceipt(receipt, operationInfo, payload.EthAddress, payload.Chain33Addr, nil, balance)

		kvs, localKvs, err = UpdateLeaf(a.statedb, a.localDB, info, leaf.GetAccountId(), payload.GetTokenId(), payload.GetAmount(), zt.Add)
		if err != nil {
			return nil, errors.Wrapf(err, "db.UpdateLeaf")
		}
		receipt, err = calProof(a.statedb, info, leaf.AccountId, payload.TokenId)
		if err != nil {
			return nil, errors.Wrapf(err, "calProof")
		}
		after := getBranchByReceipt(receipt, operationInfo, payload.EthAddress, payload.Chain33Addr, nil, receipt.Token.Balance)
		rootHash, err := common.FromHex(receipt.TreeProof.RootHash)
		if err != nil {
			return nil, errors.Wrapf(err, "hex.DecodeString")
		}
		kv := &types.KeyValue{
			Key:   getHeightKey(a.height),
			Value: rootHash,
		}
		kvs = append(kvs, kv)

		branch := &zt.OperationPairBranch{
			Before: before,
			After:  after,
		}
		operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)
		zklog := &zt.ZkReceiptLog{
			OperationInfo: operationInfo,
			LocalKvs:      localKvs,
		}
		receiptLog := &types.ReceiptLog{Ty: zt.TyDepositLog, Log: types.Encode(zklog)}
		logs = append(logs, receiptLog)
	}
	receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
	return receipts, nil
}

func getBranchByReceipt(receipt *zt.ZkReceiptLeaf, info *zt.OperationInfo, ethAddr string, chain33Addr string, pubKey *zt.ZkPubKey, balance string) *zt.OperationMetaBranch {
	info.Roots = append(info.Roots, receipt.TreeProof.RootHash)

	treePath := &zt.SiblingPath{
		Path:   receipt.TreeProof.ProofSet,
		Helper: receipt.TreeProof.GetHelpers(),
	}
	accountW := &zt.AccountWitness{
		ID:          info.AccountID,
		EthAddr:     ethAddr,
		Chain33Addr: chain33Addr,
		PubKey:      pubKey,
		Sibling:     treePath,
	}

	//token不存在不用生成TokenWitness
	if receipt.GetTokenProof() == nil {
		accountW.TokenTreeRoot = ""
		return &zt.OperationMetaBranch{
			AccountWitness: accountW,
		}
	}
	accountW.TokenTreeRoot = receipt.GetTokenProof().RootHash

	tokenPath := &zt.SiblingPath{
		Path:   receipt.TokenProof.ProofSet,
		Helper: receipt.TokenProof.GetHelpers(),
	}
	tokenW := &zt.TokenWitness{
		ID:      info.TokenID,
		Balance: balance,
		Sibling: tokenPath,
	}

	branch := &zt.OperationMetaBranch{
		AccountWitness: accountW,
		TokenWitness:   tokenW,
	}
	return branch
}

func generateTreeUpdateInfo(db dbm.KV) (*TreeUpdateInfo, error) {
	updateMap := make(map[string][]byte)
	val, err := db.Get(GetAccountTreeKey())
	if err != nil {
		//没查到就先初始化
		if err == types.ErrNotFound {
			tree := NewAccountTree(db)
			updateMap[string(GetAccountTreeKey())] = types.Encode(tree)
			return &TreeUpdateInfo{updateMap: updateMap}, nil
		} else {
			return nil, err
		}
	}
	var tree zt.AccountTree
	err = types.Decode(val, &tree)
	if err != nil {
		return nil, err
	}
	updateMap[string(GetAccountTreeKey())] = types.Encode(&tree)
	return &TreeUpdateInfo{updateMap: updateMap}, nil
}

func (a *Action) Withdraw(payload *zt.ZkWithdraw) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	var localKvs []*types.KeyValue
	err := checkParam(payload.Amount)
	if err != nil {
		return nil, errors.Wrapf(err, "checkParam")
	}

	info, err := generateTreeUpdateInfo(a.statedb)
	if err != nil {
		return nil, errors.Wrapf(err, "db.generateTreeUpdateInfo")
	}

	leaf, err := GetLeafByAccountId(a.statedb, payload.GetAccountId(), info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetLeafByAccountId")
	}
	if leaf == nil {
		return nil, errors.New("account not exist")
	}
	err = authVerification(payload.GetSignature().PubKey, leaf.GetPubKey())
	if err != nil {
		return nil, errors.Wrapf(err, "authVerification")
	}

	token, err := GetTokenByAccountIdAndTokenId(a.statedb, payload.AccountId, payload.TokenId, info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetTokenByAccountIdAndTokenId")
	}
	err = checkAmount(token, payload.GetAmount())
	if err != nil {
		return nil, errors.Wrapf(err, "db.checkAmount")
	}

	operationInfo := &zt.OperationInfo{
		BlockHeight: uint64(a.height),
		TxIndex:     uint32(a.index),
		TxType:      zt.TyWithdrawAction,
		TokenID:     payload.TokenId,
		Amount:      payload.Amount,
		SigData:     payload.Signature,
		AccountID:   payload.AccountId,
	}

	//更新之前先计算证明
	receipt, err := calProof(a.statedb, info, payload.AccountId, payload.TokenId)
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}
	before := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, receipt.Token.Balance)

	kvs, localKvs, err = UpdateLeaf(a.statedb, a.localDB, info, leaf.GetAccountId(), payload.GetTokenId(), payload.GetAmount(), zt.Sub)
	if err != nil {
		return nil, errors.Wrapf(err, "db.UpdateLeaf")
	}
	//取款之后计算证明
	receipt, err = calProof(a.statedb, info, payload.AccountId, payload.TokenId)
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}
	after := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, receipt.Token.Balance)

	rootHash, err := common.FromHex(receipt.TreeProof.RootHash)
	if err != nil {
		return nil, errors.Wrapf(err, "hex.DecodeString")
	}
	kv := &types.KeyValue{
		Key:   getHeightKey(a.height),
		Value: rootHash,
	}
	kvs = append(kvs, kv)

	branch := &zt.OperationPairBranch{
		Before: before,
		After:  after,
	}
	operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)
	zklog := &zt.ZkReceiptLog{
		OperationInfo: operationInfo,
		LocalKvs:      localKvs,
	}
	receiptLog := &types.ReceiptLog{Ty: zt.TyWithdrawLog, Log: types.Encode(zklog)}
	logs = append(logs, receiptLog)
	receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}

	return receipts, nil
}

func checkAmount(token *zt.TokenBalance, amount string) error {
	if token != nil {
		balance, _ := new(big.Int).SetString(token.Balance, 10)
		need, _ := new(big.Int).SetString(amount, 10)
		if balance.Cmp(need) >= 0 {
			return nil
		} else {
			return errors.New("balance not enough")
		}
	}
	//token为nil
	return errors.New("balance not enough")
}

func (a *Action) ContractToTree(payload *zt.ZkContractToTree) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	var localKvs []*types.KeyValue

	err := checkParam(payload.Amount)
	if err != nil {
		return nil, errors.Wrapf(err, "checkParam")
	}
	info, err := generateTreeUpdateInfo(a.statedb)
	if err != nil {
		return nil, errors.Wrapf(err, "db.generateTreeUpdateInfo")
	}

	leaf, err := GetLeafByAccountId(a.statedb, payload.GetAccountId(), info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetLeafByAccountId")
	}
	operationInfo := &zt.OperationInfo{
		BlockHeight: uint64(a.height),
		TxIndex:     uint32(a.index),
		TxType:      zt.TyContractToTreeAction,
		TokenID:     payload.TokenId,
		Amount:      payload.Amount,
		SigData:     payload.Signature,
		AccountID:   payload.AccountId,
	}

	//leaf不存在，需要新增
	if leaf == nil {
		leaf, err = GetLeafByChain33AndEthAddress(a.statedb, payload.GetChain33Addr(), payload.GetEthAddress(), info)
		if err != nil {
			return nil, errors.Wrapf(err, "db.GetLeafByChain33AndEthAddress")
		}
		//如果通过accountId没找到，但是通过给的地址找到了，说明参数出错了
		if leaf != nil {
			return nil, errors.New("account already exist")
		}
		//更新之前先计算证明
		receipt, err := calProof(a.statedb, info, payload.AccountId, payload.TokenId)
		if err != nil {
			return nil, errors.Wrapf(err, "calProof")
		}

		before := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, "0")

		kvs, localKvs, err = AddNewLeaf(a.statedb, a.localDB, info, payload.GetEthAddress(), payload.GetTokenId(), payload.GetAmount(), payload.GetChain33Addr())
		if err != nil {
			return nil, errors.Wrapf(err, "db.UpdateLeaf")
		}
		//更新合约账户
		accountKvs, err := a.UpdateContractAccount(a.fromaddr, payload.GetAmount(), payload.GetTokenId(), zt.Sub)
		if err != nil {
			return nil, errors.Wrapf(err, "db.UpdateContractAccount")
		}
		kvs = append(kvs, accountKvs...)
		//存款到叶子之后计算证明
		receipt, err = calProof(a.statedb, info, payload.AccountId, payload.TokenId)
		if err != nil {
			return nil, errors.Wrapf(err, "calProof")
		}

		after := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, receipt.Token.Balance)
		rootHash, err := common.FromHex(receipt.TreeProof.RootHash)
		if err != nil {
			return nil, errors.Wrapf(err, "hex.DecodeString")
		}
		kv := &types.KeyValue{
			Key:   getHeightKey(a.height),
			Value: rootHash,
		}
		kvs = append(kvs, kv)

		branch := &zt.OperationPairBranch{
			Before: before,
			After:  after,
		}
		operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)
	} else {
		//如果叶子存在，校验地址和accountId是否一致，防止提错账户
		if payload.GetEthAddress() != leaf.EthAddress || payload.GetEthAddress() != leaf.EthAddress {
			return nil, errors.New("check address failed")
		}
		//更新之前先计算证明
		receipt, err := calProof(a.statedb, info, payload.AccountId, payload.TokenId)
		if err != nil {
			return nil, errors.Wrapf(err, "calProof")
		}
		var balance string
		if receipt.Token == nil {
			balance = "0"
		} else {
			balance = receipt.Token.Balance
		}
		before := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, balance)

		kvs, localKvs, err = UpdateLeaf(a.statedb, a.localDB, info, leaf.GetAccountId(), payload.GetTokenId(), payload.GetAmount(), zt.Add)
		if err != nil {
			return nil, errors.Wrapf(err, "db.UpdateLeaf")
		}
		//更新合约账户
		accountKvs, err := a.UpdateContractAccount(a.fromaddr, payload.GetAmount(), payload.GetTokenId(), zt.Sub)
		if err != nil {
			return nil, errors.Wrapf(err, "db.UpdateContractAccount")
		}
		kvs = append(kvs, accountKvs...)
		//存款到叶子之后计算证明
		receipt, err = calProof(a.statedb, info, payload.AccountId, payload.TokenId)
		if err != nil {
			return nil, errors.Wrapf(err, "calProof")
		}

		after := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, receipt.Token.Balance)
		rootHash, err := common.FromHex(receipt.TreeProof.RootHash)
		if err != nil {
			return nil, errors.Wrapf(err, "hex.DecodeString")
		}
		kv := &types.KeyValue{
			Key:   getHeightKey(a.height),
			Value: rootHash,
		}
		kvs = append(kvs, kv)

		branch := &zt.OperationPairBranch{
			Before: before,
			After:  after,
		}
		operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)

	}
	zksynclog := &zt.ZkReceiptLog{
		OperationInfo: operationInfo,
		LocalKvs:      localKvs,
	}
	receiptLog := &types.ReceiptLog{Ty: zt.TyContractToTreeLog, Log: types.Encode(zksynclog)}
	logs = append(logs, receiptLog)
	receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
	return receipts, nil
}

func (a *Action) TreeToContract(payload *zt.ZkTreeToContract) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	var localKvs []*types.KeyValue

	err := checkParam(payload.Amount)
	if err != nil {
		return nil, errors.Wrapf(err, "checkParam")
	}
	info, err := generateTreeUpdateInfo(a.statedb)
	if err != nil {
		return nil, errors.Wrapf(err, "db.generateTreeUpdateInfo")
	}
	leaf, err := GetLeafByAccountId(a.statedb, payload.GetAccountId(), info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetLeafByAccountId")
	}
	if leaf == nil {
		return nil, errors.New("account not exist")
	}
	err = authVerification(payload.Signature.PubKey, leaf.GetPubKey())
	if err != nil {
		return nil, errors.Wrapf(err, "authVerification")
	}

	token, err := GetTokenByAccountIdAndTokenId(a.statedb, payload.AccountId, payload.TokenId, info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetTokenByAccountIdAndTokenId")
	}
	err = checkAmount(token, payload.GetAmount())
	if err != nil {
		return nil, errors.Wrapf(err, "db.checkAmount")
	}

	operationInfo := &zt.OperationInfo{
		BlockHeight: uint64(a.height),
		TxIndex:     uint32(a.index),
		TxType:      zt.TyTreeToContractAction,
		TokenID:     payload.TokenId,
		Amount:      payload.Amount,
		SigData:     payload.Signature,
		AccountID:   payload.AccountId,
	}

	//更新之前先计算证明
	receipt, err := calProof(a.statedb, info, payload.AccountId, payload.TokenId)
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}
	before := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, receipt.Token.Balance)

	kvs, localKvs, err = UpdateLeaf(a.statedb, a.localDB, info, leaf.GetAccountId(), payload.GetTokenId(), payload.GetAmount(), zt.Sub)
	if err != nil {
		return nil, errors.Wrapf(err, "db.UpdateLeaf")
	}
	//更新合约账户
	accountKvs, err := a.UpdateContractAccount(a.fromaddr, payload.GetAmount(), payload.GetTokenId(), zt.Add)
	if err != nil {
		return nil, errors.Wrapf(err, "db.UpdateContractAccount")
	}
	kvs = append(kvs, accountKvs...)
	//从叶子取款之后计算证明
	receipt, err = calProof(a.statedb, info, payload.AccountId, payload.TokenId)
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}

	after := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, receipt.Token.Balance)
	rootHash, err := common.FromHex(receipt.TreeProof.RootHash)
	if err != nil {
		return nil, errors.Wrapf(err, "hex.DecodeString")
	}
	kv := &types.KeyValue{
		Key:   getHeightKey(a.height),
		Value: rootHash,
	}
	kvs = append(kvs, kv)

	branch := &zt.OperationPairBranch{
		Before: before,
		After:  after,
	}
	operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)
	zklog := &zt.ZkReceiptLog{
		OperationInfo: operationInfo,
		LocalKvs:      localKvs,
	}
	receiptLog := &types.ReceiptLog{Ty: zt.TyTreeToContractLog, Log: types.Encode(zklog)}
	logs = append(logs, receiptLog)
	receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
	return receipts, nil
}

func (a *Action) UpdateContractAccount(addr string, amount string, tokenId uint64, option int32) ([]*types.KeyValue ,error) {
	execAddr := address.ExecAddress(zt.Zksync)

	accountdb, _ := account.NewAccountDB(a.api.GetConfig(), zt.Zksync, strconv.Itoa(int(tokenId)), a.statedb)
	contractAccount := accountdb.LoadExecAccount(addr, execAddr)
	change, _ := new(big.Int).SetString(amount, 10)
	//accountdb去除末尾8位小数
	shortChange := new(big.Int).Div(change, big.NewInt(100000000)).Int64()
	if option == zt.Sub {
		if contractAccount.Balance < shortChange {
			return nil, errors.New("balance not enough")
		}
		contractAccount.Balance -= shortChange
	} else {
		contractAccount.Balance += shortChange
	}

	kvs := accountdb.GetExecKVSet(execAddr, contractAccount)
	zlog.Info("zksync UpdateContractAccount", "key", string(kvs[0].GetKey()), "account", contractAccount)
	return kvs, nil
}

func (a *Action) Transfer(payload *zt.ZkTransfer) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	var localKvs []*types.KeyValue

	err := checkParam(payload.Amount)
	if err != nil {
		return nil, errors.Wrapf(err, "checkParam")
	}
	info, err := generateTreeUpdateInfo(a.statedb)
	if err != nil {
		return nil, errors.Wrapf(err, "db.generateTreeUpdateInfo")
	}
	fromLeaf, err := GetLeafByAccountId(a.statedb, payload.GetFromAccountId(), info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetLeafByAccountId")
	}
	if fromLeaf == nil {
		return nil, errors.New("account not exist")
	}
	err = authVerification(payload.Signature.PubKey, fromLeaf.PubKey)
	if err != nil {
		return nil, errors.Wrapf(err, "authVerification")
	}
	fromToken, err := GetTokenByAccountIdAndTokenId(a.statedb, payload.FromAccountId, payload.TokenId, info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetTokenByAccountIdAndTokenId")
	}
	err = checkAmount(fromToken, payload.GetAmount())
	if err != nil {
		return nil, errors.Wrapf(err, "db.checkAmount")
	}

	operationInfo := &zt.OperationInfo{
		BlockHeight: uint64(a.height),
		TxIndex:     uint32(a.index),
		TxType:      zt.TyTransferAction,
		TokenID:     payload.TokenId,
		Amount:      payload.Amount,
		SigData:     payload.Signature,
		AccountID:   payload.FromAccountId,
	}

	//更新之前先计算证明
	receipt, err := calProof(a.statedb, info, payload.FromAccountId, payload.TokenId)
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}

	before := getBranchByReceipt(receipt, operationInfo, fromLeaf.EthAddress, fromLeaf.Chain33Addr, fromLeaf.PubKey, receipt.Token.Balance)

	//更新fromLeaf
	fromKvs, fromLocal, err := UpdateLeaf(a.statedb, a.localDB, info, fromLeaf.GetAccountId(), payload.GetTokenId(), payload.GetAmount(), zt.Sub)
	if err != nil {
		return nil, errors.Wrapf(err, "db.UpdateLeaf")
	}
	kvs = append(kvs, fromKvs...)
	localKvs = append(localKvs, fromLocal...)
	//更新之后计算证明
	receipt, err = calProof(a.statedb, info, payload.FromAccountId, payload.TokenId)
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}

	after := getBranchByReceipt(receipt, operationInfo, fromLeaf.EthAddress, fromLeaf.Chain33Addr, fromLeaf.PubKey, receipt.Token.Balance)

	branch := &zt.OperationPairBranch{
		Before: before,
		After:  after,
	}
	operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)

	toLeaf, err := GetLeafByAccountId(a.statedb, payload.ToAccountId, info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetLeafByAccountId")
	}
	if toLeaf == nil {
		return nil, errors.New("account not exist")
	}

	//更新之前先计算证明
	receipt, err = calProof(a.statedb, info, payload.ToAccountId, payload.TokenId)
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}
	var balance string
	if receipt.Token == nil {
		balance = "0"
	} else {
		balance = receipt.Token.Balance
	}
	before = getBranchByReceipt(receipt, operationInfo, toLeaf.EthAddress, toLeaf.Chain33Addr, toLeaf.PubKey, balance)

	//更新toLeaf
	tokvs, toLocal, err := UpdateLeaf(a.statedb, a.localDB, info, toLeaf.GetAccountId(), payload.GetTokenId(), payload.GetAmount(), zt.Add)
	if err != nil {
		return nil, errors.Wrapf(err, "db.UpdateLeaf")
	}
	kvs = append(kvs, tokvs...)
	localKvs = append(localKvs, toLocal...)
	//更新之后计算证明
	receipt, err = calProof(a.statedb, info, payload.GetToAccountId(), payload.GetTokenId())
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}
	after = getBranchByReceipt(receipt, operationInfo, toLeaf.EthAddress, toLeaf.Chain33Addr, toLeaf.PubKey, receipt.Token.Balance)
	rootHash, err := common.FromHex(receipt.TreeProof.RootHash)
	if err != nil {
		return nil, errors.Wrapf(err, "hex.DecodeString")
	}
	kv := &types.KeyValue{
		Key:   getHeightKey(a.height),
		Value: rootHash,
	}
	kvs = append(kvs, kv)

	branch = &zt.OperationPairBranch{
		Before: before,
		After:  after,
	}
	operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)
	zklog := &zt.ZkReceiptLog{
		OperationInfo: operationInfo,
		LocalKvs:      localKvs,
	}
	receiptLog := &types.ReceiptLog{Ty: zt.TyTransferLog, Log: types.Encode(zklog)}
	logs = append(logs, receiptLog)

	receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
	return receipts, nil
}

func (a *Action) TransferToNew(payload *zt.ZkTransferToNew) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	var localKvs []*types.KeyValue
	err := checkParam(payload.Amount)
	if err != nil {
		return nil, errors.Wrapf(err, "checkParam")
	}
	info, err := generateTreeUpdateInfo(a.statedb)
	if err != nil {
		return nil, errors.Wrapf(err, "db.generateTreeUpdateInfo")
	}
	fromLeaf, err := GetLeafByAccountId(a.statedb, payload.GetFromAccountId(), info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetLeafByAccountId")
	}
	if fromLeaf == nil {
		return nil, errors.New("account not exist")
	}
	err = authVerification(payload.Signature.PubKey, fromLeaf.PubKey)
	if err != nil {
		return nil, errors.Wrapf(err, "authVerification")
	}

	fromToken, err := GetTokenByAccountIdAndTokenId(a.statedb, payload.FromAccountId, payload.TokenId, info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetTokenByAccountIdAndTokenId")
	}
	err = checkAmount(fromToken, payload.GetAmount())
	if err != nil {
		return nil, errors.Wrapf(err, "db.checkAmount")
	}

	operationInfo := &zt.OperationInfo{
		BlockHeight: uint64(a.height),
		TxIndex:     uint32(a.index),
		TxType:      zt.TyTransferToNewAction,
		TokenID:     payload.TokenId,
		Amount:      payload.Amount,
		SigData:     payload.Signature,
		AccountID:   payload.FromAccountId,
	}

	toLeaf, err := GetLeafByChain33AndEthAddress(a.statedb, payload.GetToChain33Address(), payload.GetToEthAddress(), info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetLeafByChain33AndEthAddress")
	}
	if toLeaf != nil {
		return nil, errors.New("to account already exist")
	}
	//更新之前先计算证明
	receipt, err := calProof(a.statedb, info, payload.GetFromAccountId(), payload.GetTokenId())
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}

	before := getBranchByReceipt(receipt, operationInfo, fromLeaf.EthAddress, fromLeaf.Chain33Addr, fromLeaf.PubKey, receipt.Token.Balance)

	//更新fromLeaf
	fromkvs, fromLocal, err := UpdateLeaf(a.statedb, a.localDB, info, fromLeaf.GetAccountId(), payload.GetTokenId(), payload.GetAmount(), zt.Sub)
	if err != nil {
		return nil, errors.Wrapf(err, "db.UpdateLeaf")
	}
	kvs = append(kvs, fromkvs...)
	localKvs = append(localKvs, fromLocal...)
	//更新之后计算证明
	receipt, err = calProof(a.statedb, info, payload.GetFromAccountId(), payload.GetTokenId())
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}

	after := getBranchByReceipt(receipt, operationInfo, fromLeaf.EthAddress, fromLeaf.Chain33Addr, fromLeaf.PubKey, receipt.Token.Balance)

	branch := &zt.OperationPairBranch{
		Before: before,
		After:  after,
	}
	operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)

	tree, err := getAccountTree(a.statedb, info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.getAccountTree")
	}
	accountId := tree.GetTotalIndex() + 1
	//更新之前先计算证明
	receipt, err = calProof(a.statedb, info, accountId, payload.GetTokenId())
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}

	before = getBranchByReceipt(receipt, operationInfo, payload.ToEthAddress, payload.ToChain33Address, nil, "0")

	//新增toLeaf
	tokvs, toLocal, err := AddNewLeaf(a.statedb, a.localDB, info, payload.GetToEthAddress(), payload.GetTokenId(), payload.GetAmount(), payload.GetToChain33Address())
	if err != nil {
		return nil, errors.Wrapf(err, "db.AddNewLeaf")
	}
	kvs = append(kvs, tokvs...)
	localKvs = append(localKvs, toLocal...)
	//新增之后计算证明
	receipt, err = calProof(a.statedb, info, accountId, payload.GetTokenId())
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}

	after = getBranchByReceipt(receipt, operationInfo, payload.ToEthAddress, payload.ToChain33Address, nil, receipt.Token.Balance)
	rootHash, err := common.FromHex(receipt.TreeProof.RootHash)
	if err != nil {
		return nil, errors.Wrapf(err, "hex.DecodeString")
	}
	kv := &types.KeyValue{
		Key:   getHeightKey(a.height),
		Value: rootHash,
	}
	kvs = append(kvs, kv)

	branch = &zt.OperationPairBranch{
		Before: before,
		After:  after,
	}
	operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)
	zklog := &zt.ZkReceiptLog{
		OperationInfo: operationInfo,
		LocalKvs:      localKvs,
	}
	receiptLog := &types.ReceiptLog{Ty: zt.TyTransferToNewLog, Log: types.Encode(zklog)}
	logs = append(logs, receiptLog)

	receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
	return receipts, nil
}

func (a *Action) ForceExit(payload *zt.ZkForceExit) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	var localKvs []*types.KeyValue
	info, err := generateTreeUpdateInfo(a.statedb)
	if err != nil {
		return nil, errors.Wrapf(err, "db.generateTreeUpdateInfo")
	}
	leaf, err := GetLeafByAccountId(a.statedb, payload.AccountId, info)
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}

	if leaf == nil {
		return nil, errors.New("account not exist")
	}

	token, err := GetTokenByAccountIdAndTokenId(a.statedb, payload.AccountId, payload.TokenId, info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetTokenByAccountIdAndTokenId")
	}

	//token不存在时，不需要取
	if token == nil {
		return nil, errors.New("token not find")
	}

	operationInfo := &zt.OperationInfo{
		BlockHeight: uint64(a.height),
		TxIndex:     uint32(a.index),
		TxType:      zt.TyForceExitAction,
		TokenID:     payload.TokenId,
		Amount:      token.Balance,
		SigData:     payload.Signature,
		AccountID:   payload.AccountId,
	}

	//更新之前先计算证明
	receipt, err := calProof(a.statedb, info, payload.GetAccountId(), payload.GetTokenId())
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}
	before := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, receipt.Token.Balance)

	//更新fromLeaf
	kvs, localKvs, err = UpdateLeaf(a.statedb, a.localDB, info, leaf.GetAccountId(), payload.GetTokenId(), token.Balance, zt.Sub)
	if err != nil {
		return nil, errors.Wrapf(err, "db.UpdateLeaf")
	}
	//更新之后计算证明
	receipt, err = calProof(a.statedb, info, payload.GetAccountId(), payload.GetTokenId())
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}

	after := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, receipt.Token.Balance)
	rootHash, err := common.FromHex(receipt.TreeProof.RootHash)
	if err != nil {
		return nil, errors.Wrapf(err, "hex.DecodeString")
	}
	kv := &types.KeyValue{
		Key:   getHeightKey(a.height),
		Value: rootHash,
	}
	kvs = append(kvs, kv)

	branch := &zt.OperationPairBranch{
		Before: before,
		After:  after,
	}
	operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)
	zklog := &zt.ZkReceiptLog{
		OperationInfo: operationInfo,
		LocalKvs:      localKvs,
	}
	receiptLog := &types.ReceiptLog{Ty: zt.TyForceExitLog, Log: types.Encode(zklog)}
	logs = append(logs, receiptLog)
	receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
	return receipts, nil

}

func calProof(statedb dbm.KV, info *TreeUpdateInfo, accountId uint64, tokenId uint64) (*zt.ZkReceiptLeaf, error) {
	receipt := &zt.ZkReceiptLeaf{}

	leaf, err := GetLeafByAccountId(statedb, accountId, info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetLeafByAccountId")
	}
	receipt.Leaf = leaf

	token, err := GetTokenByAccountIdAndTokenId(statedb, accountId, tokenId, info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetTokenByAccountIdAndTokenId")
	}
	receipt.Token = token

	leafProof, err := CalLeafProof(statedb, leaf, info)
	if err != nil {
		return nil, errors.Wrapf(err, "CalLeafProof")
	}
	receipt.TreeProof = leafProof

	tokenProof, err := CalTokenProof(statedb, leaf, token, info)
	if err != nil {
		return nil, errors.Wrapf(err, "CalTokenProof")
	}
	receipt.TokenProof = tokenProof

	return receipt, nil
}

func (a *Action) SetPubKey(payload *zt.ZkSetPubKey) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	var localKvs []*types.KeyValue
	info, err := generateTreeUpdateInfo(a.statedb)
	if err != nil {
		return nil, errors.Wrapf(err, "db.generateTreeUpdateInfo")
	}

	leaf, err := GetLeafByAccountId(a.statedb, payload.GetAccountId(), info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetLeafByEthAddress")
	}
	if leaf == nil {
		return nil, errors.New("account not exist")
	}

	//校验预存的地址是否和公钥匹配
	pubKey := &eddsa.PublicKey{}
	pubKey.A.X.SetString(payload.PubKey.X)
	pubKey.A.Y.SetString(payload.PubKey.Y)
	hash := mimc.NewMiMC(zt.ZkMimcHashSeed)
	hash.Write(pubKey.Bytes())
	if common.ToHex(hash.Sum(nil)) != leaf.Chain33Addr {
		return nil, errors.New("not your account")
	}

	if err != nil {
		return nil, errors.Wrapf(err, "authVerification")
	}

	operationInfo := &zt.OperationInfo{
		BlockHeight: uint64(a.height),
		TxIndex:     uint32(a.index),
		TxType:      zt.TySetPubKeyAction,
		TokenID:     leaf.TokenIds[0],
		Amount:      "0",
		SigData:     payload.Signature,
		AccountID:   payload.AccountId,
	}

	//更新之前先计算证明
	receipt, err := calProof(a.statedb, info, payload.AccountId, leaf.TokenIds[0])
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}
	before := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, nil, receipt.Token.Balance)

	kvs, localKvs, err = UpdatePubKey(a.statedb, a.localDB, info, payload.GetPubKey(), payload.AccountId)
	if err != nil {
		return nil, errors.Wrapf(err, "db.UpdateLeaf")
	}
	//更新之后计算证明
	receipt, err = calProof(a.statedb, info, payload.AccountId, leaf.TokenIds[0])
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}
	after := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, payload.PubKey, receipt.Token.Balance)
	rootHash, err := common.FromHex(receipt.TreeProof.RootHash)
	if err != nil {
		return nil, errors.Wrapf(err, "hex.DecodeString")
	}
	kv := &types.KeyValue{
		Key:   getHeightKey(a.height),
		Value: rootHash,
	}
	kvs = append(kvs, kv)

	branch := &zt.OperationPairBranch{
		Before: before,
		After:  after,
	}
	operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)
	zklog := &zt.ZkReceiptLog{
		OperationInfo: operationInfo,
		LocalKvs:      localKvs,
	}
	receiptLog := &types.ReceiptLog{Ty: zt.TySetPubKeyLog, Log: types.Encode(zklog)}
	logs = append(logs, receiptLog)
	receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
	return receipts, nil
}


func (a *Action) FullExit(payload *zt.ZkFullExit) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	var localKvs []*types.KeyValue

	//只有管理员能操作
	cfg := a.api.GetConfig()
	if !isSuperManager(cfg, a.fromaddr) && !isVerifier(a.statedb, a.fromaddr) {
		return nil, errors.Wrapf(types.ErrNotAllow, "from addr is not manager")
	}

	info, err := generateTreeUpdateInfo(a.statedb)
	if err != nil {
		return nil, errors.Wrapf(err, "db.generateTreeUpdateInfo")
	}
	leaf, err := GetLeafByAccountId(a.statedb, payload.AccountId, info)
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}

	if leaf == nil {
		return nil, errors.New("account not exist")
	}

	token, err := GetTokenByAccountIdAndTokenId(a.statedb, payload.AccountId, payload.TokenId, info)
	if err != nil {
		return nil, errors.Wrapf(err, "db.GetTokenByAccountIdAndTokenId")
	}

	//token不存在时，不需要取
	if token == nil {
		return nil, errors.New("token not find")
	}

	operationInfo := &zt.OperationInfo{
		BlockHeight: uint64(a.height),
		TxIndex:     uint32(a.index),
		TxType:      zt.TyForceExitAction,
		TokenID:     payload.TokenId,
		Amount:      token.Balance,
		SigData:     payload.Signature,
		AccountID:   payload.AccountId,
	}

	//更新之前先计算证明
	receipt, err := calProof(a.statedb, info, payload.GetAccountId(), payload.GetTokenId())
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}
	before := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, receipt.Token.Balance)

	//更新fromLeaf
	kvs, localKvs, err = UpdateLeaf(a.statedb, a.localDB, info, leaf.GetAccountId(), payload.GetTokenId(), token.Balance, zt.Sub)
	if err != nil {
		return nil, errors.Wrapf(err, "db.UpdateLeaf")
	}
	//更新之后计算证明
	receipt, err = calProof(a.statedb, info, payload.GetAccountId(), payload.GetTokenId())
	if err != nil {
		return nil, errors.Wrapf(err, "calProof")
	}

	after := getBranchByReceipt(receipt, operationInfo, leaf.EthAddress, leaf.Chain33Addr, leaf.PubKey, receipt.Token.Balance)
	rootHash, err := common.FromHex(receipt.TreeProof.RootHash)
	if err != nil {
		return nil, errors.Wrapf(err, "hex.DecodeString")
	}
	kv := &types.KeyValue{
		Key:   getHeightKey(a.height),
		Value: rootHash,
	}
	kvs = append(kvs, kv)

	branch := &zt.OperationPairBranch{
		Before: before,
		After:  after,
	}
	operationInfo.OperationBranches = append(operationInfo.GetOperationBranches(), branch)
	zklog := &zt.ZkReceiptLog{
		OperationInfo: operationInfo,
		LocalKvs:      localKvs,
	}
	receiptLog := &types.ReceiptLog{Ty: zt.TyForceExitLog, Log: types.Encode(zklog)}
	logs = append(logs, receiptLog)
	receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
	return receipts, nil

}

//验证身份
func authVerification(signPubKey *zt.ZkPubKey, leafPubKey *zt.ZkPubKey) error {
	if signPubKey == nil || leafPubKey == nil {
		return errors.New("set your pubKey")
	}
	if signPubKey.GetX() != leafPubKey.GetX() || signPubKey.GetY() != leafPubKey.GetY() {
		return errors.New("not your account")
	}
	return nil
}

//检查参数
func checkParam(amount string) error {
	if amount == "" || amount == "0" || strings.HasPrefix(amount, "-") {
		return types.ErrAmount
	}
	return nil
}
