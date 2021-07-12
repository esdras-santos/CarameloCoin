package blockchain_test

import (
	"gochain/wallet"
	"testing"
)

func TestId(t *testing.T){
	w := wallet.MakeWallet()
	scriptPubKey := []byte{byte(len(w.PublicKey))}
	scriptPubKey = append(scriptPubKey, w.PublicKey...)
	scriptPubKey = append(scriptPubKey, 0xac)
	
}

