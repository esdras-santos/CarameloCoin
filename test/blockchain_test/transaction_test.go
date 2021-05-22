package blockchain_test

import(
	"testing"
	"gochain/blockchain"
)

func TestReadVarint(t *testing.T){
	var buf uint
	s := []byte{0xfd,0xff,0x00}
	blockchain.ReadVarint(s,&buf)
	if buf != 255{
		t.Error("fail")
	}
}