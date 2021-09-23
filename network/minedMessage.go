package network

import (
	"bytes"
	"encoding/gob"
	"gochain/blockchain"
)

type MinedMessage struct {
	Command     []byte
	FromIp      []byte
	Transaction *blockchain.Transaction
}

func (tm *MinedMessage) Init(fromIp []byte, tx *blockchain.Transaction) {
	tm.Command = []byte("mined")
	tm.FromIp = fromIp
	tm.Transaction = tx
}

func (tm MinedMessage) Serialize() []byte {
	b := bytes.Buffer{}
    e := gob.NewEncoder(&b)
    err := e.Encode(tm)
	Handle(err)
    return b.Bytes()
}

func (tm MinedMessage) GetCommand() []byte {
	return tm.Command
}

func (tm *MinedMessage) Parse(data []byte) *blockchain.Transaction {
	
	var mm MinedMessage
    b := bytes.Buffer{}
    b.Write(data)
    d := gob.NewDecoder(&b)
    err := d.Decode(&mm)
	Handle(err)
	txn := mm.Transaction
    return txn
}