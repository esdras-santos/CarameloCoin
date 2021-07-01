package network

import (
	"gochain/blockchain"	
	"gochain/utils"
)


type TransactionMessage struct {
	Command     []byte
	NodeAddress []byte
	Transaction *blockchain.Transaction
}

func (tm *TransactionMessage) Init(nodeAddress string, tx *blockchain.Transaction){
	tm.Command = []byte("transaction")
	tm.NodeAddress = AddressToBytes(nodeAddress) //IPv4
	tm.Transaction = tx
}

func (tm TransactionMessage) Serialize() []byte{
	result := utils.ToLittleEndian(tm.NodeAddress,4)
	result = append(result, tm.Transaction.Serialize()...)
	return result
}

func (tm TransactionMessage) GetCommand() []byte{
	return tm.Command
}

func (tm *TransactionMessage) Parse(data []byte) {
	var tx blockchain.Transaction
	tm.NodeAddress = utils.ToLittleEndian(data[:4],4)
	tm.Transaction = tx.Parse(data[4:])
}