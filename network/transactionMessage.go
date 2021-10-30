package network

import (
	"bytes"
	"encoding/gob"
	"gochain/blockchain"
	
)


type TransactionMessage struct {
	Command     []byte
	Transaction *blockchain.Transaction
}

func (tm *TransactionMessage) Init(tx *blockchain.Transaction){
	tm.Command = []byte("transaction")
	
	tm.Transaction = tx
}

func (tm TransactionMessage) Serialize() []byte{
	b := bytes.Buffer{}
    e := gob.NewEncoder(&b)
    err := e.Encode(tm)
	Handle(err)
    return b.Bytes()
}

func (tm TransactionMessage) GetCommand() []byte{
	return tm.Command
}

func (tm *TransactionMessage) Parse(data []byte) *blockchain.Transaction {
	var tmsg TransactionMessage
    b := bytes.Buffer{}
    b.Write(data)
    d := gob.NewDecoder(&b)
    err := d.Decode(&tmsg)
	Handle(err)
	txn := tmsg.Transaction
    return txn
}