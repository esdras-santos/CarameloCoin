package network

import (
	"gochain/blockchain"
	"gochain/utils"
)


type TransactionMessage struct {
	Command     []byte
	FromIp	[]byte
	Transaction *blockchain.Transaction
}

func (tm *TransactionMessage) Init(fromIp []byte,tx *blockchain.Transaction){
	tm.Command = []byte("transaction")
	tm.FromIp = fromIp
	tm.Transaction = tx
}

func (tm TransactionMessage) Serialize() []byte{
	result := utils.ToLittleEndian(tm.FromIp,4)
	result = append(result, tm.Transaction.Serialize()...)
	return result
}

func (tm TransactionMessage) GetCommand() []byte{
	return tm.Command
}

func (tm *TransactionMessage) Parse(data []byte) {
	var tx blockchain.Transaction
	tm.FromIp = utils.ToLittleEndian(data[:4],4)
	tm.Transaction = tx.Parse(data[4:])
}