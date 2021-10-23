package network_test

import (
	"gochain/blockchain"
	"gochain/network"
	"log"
	"testing"
)

func TestConnect(t *testing.T){
	nw := network.Connect()
	chain := blockchain.ContinueBlockChain("../../tmp/blocks")

	bm := network.BlockMessage{}
	block,err := chain.GetBlock(chain.LastHash)
	if err != nil{
		log.Panic(err)
	}
	
		
	bm.Init(block)
	ne := network.NetworkEnvelope{bm.GetCommand(),bm.Serialize()}
	print("\n enveloped\n")
	nw.Publish(ne)
}