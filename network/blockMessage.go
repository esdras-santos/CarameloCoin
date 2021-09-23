package network

import "gochain/blockchain"


//this message is a response to the "getblock" command
type BlockMessage struct {
	Command []byte
	Block *blockchain.Block
}

func (bm *BlockMessage) Init(block *blockchain.Block){
	bm.Command = []byte("block")
	bm.Block = block
}

func (bm BlockMessage) GetCommand() []byte{
	return bm.Command
}

func (bm BlockMessage) Serialize() []byte{
	return bm.Block.Serialize()
}

func (bm *BlockMessage) Parse(data []byte) *blockchain.Block{
	block := &blockchain.Block{}
	block.Parse(data)
	bm.Init(block)
	return bm
}