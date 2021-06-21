package network

import (
	"encoding/binary"
	"gochain/blockchain"
	"gochain/utils"
	"log"
)


type HeadersMessage struct {
	Command []byte
	Blocks  []blockchain.BlockHeader
}

func (hm *HeadersMessage) Init(blocks []blockchain.BlockHeader){
	hm.Command = []byte("headers")
	hm.Blocks = blocks
}

func (hm *HeadersMessage) Parse(data []byte){
	var numberBlocks uint 
	utils.ReadVarint(data,&numberBlocks)
	var i uint
	var startIn int
	if numberBlocks <= 253{
		startIn = 1
	}else if numberBlocks <= 254{
		startIn = 2
	}else if numberBlocks <= 255{
		startIn = 3
	}
	var blocks []blockchain.BlockHeader
	var block blockchain.BlockHeader 
	for i < numberBlocks{
		var numTxs uint16
		blocks = append(blocks, block.Parse(data[startIn:startIn+80]))
		numTxs = binary.LittleEndian.Uint16(data[startIn+81:startIn+82])
		if numTxs != 0{
			log.Panic("number of txs not 0")
		}
		startIn += 82
		i++
	}
	hm.Init(blocks) 
}
