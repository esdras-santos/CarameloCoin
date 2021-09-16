package blockchain_test

import (
	"gochain/blockchain"
	"gochain/wallet"
	"testing"
)

func TestInitBlockchain(t *testing.T){
	w := wallet.MakeWallet()
	blockchain.InitBlockChain(w, "./tmp/blocks_test")
}