package executor

import (
	"fmt"

	"github.com/assetcloud/chain/common/address"
	zt "github.com/assetcloud/plugin/plugin/dapp/zksync/types"
)

func GetAccountIdPrimaryKeyPrefix() string {
	return fmt.Sprintf("%s", KeyPrefixStateDB+"accountId-")
}

func GetAccountIdPrimaryKey(accountId uint64) []byte {
	return []byte(fmt.Sprintf("%s%022d", KeyPrefixStateDB+"accountId-", accountId))
}

func GetLocalChainEthPrimaryKey(chainAddr string, ethAddr string) []byte {
	return []byte(fmt.Sprintf("%s-%s", address.FormatAddrKey(chainAddr), address.FormatAddrKey(ethAddr)))
}

func GetChainEthPrimaryKey(chainAddr string, ethAddr string) []byte {
	return []byte(fmt.Sprintf("%s%s-%s", KeyPrefixStateDB, address.FormatAddrKey(chainAddr),
		address.FormatAddrKey(ethAddr)))
}

func GetTokenPrimaryKey(accountId uint64, tokenId uint64) []byte {
	return []byte(fmt.Sprintf("%s%022d%s%022d", KeyPrefixStateDB+"token-", accountId, "-", tokenId))
}

func GetTokenPrimaryKeyPrefix() string {
	return fmt.Sprintf("%s", KeyPrefixStateDB+"token-")
}

func GetNFTIdPrimaryKey(nftTokenId uint64) []byte {
	return []byte(fmt.Sprintf("%s%022d", KeyPrefixStateDB+"nftTokenId-", nftTokenId))
}

func GetNFTHashPrimaryKey(nftHash string) []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+"nftHash-"+nftHash))
}

func GetRootIndexPrimaryKey(rootIndex uint64) []byte {
	return []byte(fmt.Sprintf("%s%016d", KeyPrefixStateDB+"rootIndex-", rootIndex))
}

func GetAccountTreeKey() []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+"accountTree"))
}

func getHeightKey(height int64) []byte {
	return []byte(fmt.Sprintf("%s%022d", KeyPrefixStateDB+"treeHeightRoot", height))
}

func getVerifyKey(chainTitleId string) []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+chainTitleId+"-verifyKey"))
}

func getVerifier(chainTitleId string) []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+chainTitleId+"-"+zt.ZkVerifierKey))
}

func getLastProofKey() []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+"lastProof"))
}

func getLastOnChainProofIdKey(chainTitleId string) []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+chainTitleId+"-lastOnChainProofId"))
}

func getValidatorsKey() []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+"validators"))
}

func getEthPriorityQueueKey(chainID uint32) []byte {
	return []byte(fmt.Sprintf("%s-%d", KeyPrefixStateDB+"priorityQueue", chainID))
}

//特意把title放后面，方便按id=1搜索所有的chain
func getProofIdCommitProofKey(chainTitleId string, proofId uint64) []byte {
	return []byte(fmt.Sprintf("%016d-%s", proofId, chainTitleId))
}

func getRootCommitProofKey(chainTitleId, root string) []byte {
	return []byte(fmt.Sprintf("%s-%s", chainTitleId, root))
}

func getHistoryAccountTreeKey(proofId, accountId uint64) []byte {
	return []byte(fmt.Sprintf("%016d.%16d", proofId, accountId))
}

func getZkFeeKey(actionTy int32, tokenId uint64) []byte {
	return []byte(fmt.Sprintf("%s%02d-%03d", KeyPrefixStateDB+"fee-", actionTy, tokenId))
}

func CalcLatestAccountIDKey() []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+"latestAccountID"))
}

func getExodusModeKey() []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+"exodusMode"))
}

//GetTokenSymbolKey tokenId 对应symbol
func GetTokenSymbolKey(tokenId string) []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+"tokenId-"+tokenId))
}

//GetTokenSymbolIdKey token symbol 对应id
func GetTokenSymbolIdKey(symbol string) []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+"tokenSym-"+symbol))
}

func getLastProofIdKey(chainTitleId string) []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+chainTitleId+"-lastProofId"))
}

func getMaxRecordProofIdKey(chainTitleId string) []byte {
	return []byte(fmt.Sprintf("%s", KeyPrefixStateDB+chainTitleId+"-maxRecordProofId"))
}

func getProofIdKey(chainTitleId string, id uint64) []byte {
	return []byte(fmt.Sprintf("%s%022d", KeyPrefixStateDB+chainTitleId+"-ProofId", id))
}
