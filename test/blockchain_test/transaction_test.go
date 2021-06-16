package blockchain_test

import (
	"gochain/blockchain"
	"gochain/utils"
	"gochain/wallet"
	"testing"
)

func TestId(t *testing.T){
	w := wallet.MakeWallet()
	scriptPubKey := []byte{byte(len(w.PublicKey))}
	scriptPubKey = append(scriptPubKey, w.PublicKey...)
	scriptPubKey = append(scriptPubKey, 0xac)
	
}

func TestReadVarint(t *testing.T){
	var buf uint
	s := []byte{0xfd,0xff,0x00}
	utils.ReadVarint(s,&buf)
	if buf != 255{
		t.Error("fail")
	}
}