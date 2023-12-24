package main

import (
	"crypto/ecdsa"
	crand "crypto/rand"
	"fmt"

	chainAddress "github.com/assetcloud/chain/common/address"

	"github.com/ethereum/go-ethereum/common/math"

	"github.com/ethereum/go-ethereum/crypto"

	chainCommon "github.com/assetcloud/chain/common"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func KeyManageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "generate secp256k1 private key, show public key and address",
	}
	cmd.AddCommand(
		chainCmd(),
		ethereumCmd(),
	)
	return cmd
}

func chainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain",
		Short: "generate secp256k1 private key or show info for chain",
	}
	cmd.AddCommand(
		generareChainKeyCmd(),
		showChainpubCmd(),
	)
	return cmd
}

func generareChainKeyCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "new",
		Short: "create a private key for chain",
		Run:   generareChainKey,
	}
	return cmd
}

func generareChainKey(cmd *cobra.Command, _ []string) {

	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if nil != err {
		fmt.Println("Failed to generate private key for chain" + err.Error())
		return
	}

	fmt.Println(common.Bytes2Hex(privateKey.Serialize()))
}

func showChainpubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show public key(inluding uncompressed and compressed, chain address)",
		Run:   showChainKey,
	}
	addShowKeyFlags(cmd)
	return cmd
}

func showChainKey(cmd *cobra.Command, _ []string) {
	key, _ := cmd.Flags().GetString("key")
	privateKeySlice, err := chainCommon.FromHex(key)
	if nil != err {
		fmt.Println("convert string error due to:" + err.Error())
		return
	}

	if len(privateKeySlice) != 32 {
		fmt.Println("invalid priv key length", len(privateKeySlice))
		return
	}
	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privateKeySlice)

	uncompressedKey := pubKey.SerializeUncompressed()
	uncompressedKey = uncompressedKey[1:]
	compressedKey := pubKey.SerializeCompressed()
	fmt.Println("The uncompressed pub key = "+common.Bytes2Hex(uncompressedKey), "with length=", len(uncompressedKey))
	fmt.Println("The compressed pub key = "+common.Bytes2Hex(compressedKey), "with length=", len(compressedKey))
	fmt.Println("Chain address = " + chainAddress.PubKeyToAddr(chainAddress.DefaultID, compressedKey))
}

////////////////////////////////////////////////////////

func ethereumCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ethereum",
		Short: "generate secp256k1 private key or show info for ethereum",
	}
	cmd.AddCommand(
		generareEthereumKeyCmd(),
		showEtheremKeyCmd(),
	)
	return cmd
}

func showEtheremKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show public key(inluding uncompressed and compressed, ethereum address)",
		Run:   showEtheremKey,
	}
	addShowKeyFlags(cmd)
	return cmd
}

func addShowKeyFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "private key")
	_ = cmd.MarkFlagRequired("key")
}

// Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
func showEtheremKey(cmd *cobra.Command, _ []string) {
	key, _ := cmd.Flags().GetString("key")
	privateKeySlice, err := chainCommon.FromHex(key)
	if nil != err {
		fmt.Println("convert string error due to:" + err.Error())
		return
	}

	if len(privateKeySlice) != 32 {
		fmt.Println("invalid priv key length", len(privateKeySlice))
		return
	}

	privateKey, err := crypto.ToECDSA(privateKeySlice)
	if nil != err {
		fmt.Println("Failed ToECDSA due to " + err.Error())
		return
	}

	_, pubKey := btcec.PrivKeyFromBytes(crypto.S256(), privateKeySlice)
	uncompressedKey := pubKey.SerializeUncompressed()
	uncompressedKey = uncompressedKey[1:]
	compressedKey := pubKey.SerializeCompressed()
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("The uncompressed pub key = "+common.Bytes2Hex(uncompressedKey), "with length=", len(uncompressedKey))
	fmt.Println("The compressed pub key = "+common.Bytes2Hex(compressedKey), "with length=", len(compressedKey))
	fmt.Println("Address = " + address.String())
}

func generareEthereumKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "create a private key for ethereum",
		Run:   generareEthereumKey,
	}
	return cmd
}

func generareEthereumKey(cmd *cobra.Command, _ []string) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), crand.Reader)
	if err != nil {
		fmt.Println("Failed to generate private key for ethereum due to:" + err.Error())
		return
	}
	privateKeyBytes := math.PaddedBigBytes(privateKeyECDSA.D, 32)
	fmt.Println(common.Bytes2Hex(privateKeyBytes))
}
