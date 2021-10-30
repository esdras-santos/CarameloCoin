package network

import (
	"bytes"
	"encoding/gob"
	"gochain/blockchain"
)

//this message is a response to the "getblock" command
type GenBlockMessage struct {
	
	MinerAddr []byte
	Block *blockchain.Block
}

func (gbm *GenBlockMessage) Init(mineraddr []byte, block blockchain.Block){
	gbm.MinerAddr = mineraddr
	
	gbm.Block = &block
}

func (gbm GenBlockMessage) GetCommand() []byte{
	return []byte("genblock")
}

func (gbm GenBlockMessage) Serialize() []byte{
	by := bytes.Buffer{}
    e := gob.NewEncoder(&by)
    err := e.Encode(gbm)
	Handle(err)
    return by.Bytes()
}

func (gbm *GenBlockMessage) Parse(data []byte) (*blockchain.Block, []byte){
	var genblockmsg GenBlockMessage
    by := bytes.Buffer{}
    by.Write(data)
    d := gob.NewDecoder(&by)
    err := d.Decode(&genblockmsg)
	Handle(err)

    return genblockmsg.Block, genblockmsg.MinerAddr
}