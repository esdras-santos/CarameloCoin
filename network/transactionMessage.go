package network

import "gochain/blockchain"

type TransactionMessage struct {
	Command     []byte
	NodeAddress []byte
	Transaction *blockchain.Transaction
}

func (tm *TransactionMessage) Init(nodeAddress string, tx *blockchain.Transaction){
	tm.Command = []byte("transaction")
	tm.NodeAddress = []byte(nodeAddress)
	tm.Transaction = tx
}

func (tm TransactionMessage) Serialize() []byte{
	result := tm.NodeAddress
	result = append(result, tm.Transaction.Serialize()...)
	return result
}

func (tm TransactionMessage) GetCommand() []byte{
	return tm.Command
}

func (tm *TransactionMessage) Parse(data []byte) {
	tm.NodeAddress =  data[]
}