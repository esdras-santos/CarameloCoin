package network

import (
	"fmt"
	"gochain/blockchain"
)



type NodeCommand struct {
	//the ip and port of the host
	hostAddress *string
}

//this need to be started in a loop(with the range of your known nodes) as a goroutine for share the command through out all the network 
func (nm *NodeCommand) Init(host string) {
	*nm.hostAddress = fmt.Sprintf("%s%s", host, PORT)
}

func (nm *NodeCommand) HandShake(){
	var message VersionMessage
	message.Init(nil, AddressToBytes(*nm.hostAddress),nil, []byte{0x01})
		
	SendData(*nm.hostAddress,message)
}

func (nm *NodeCommand) GetBlock(height int64){
	var message GetBlockMessage
	message.Init(AddressToBytes(*nm.hostAddress))
	SendData(*nm.hostAddress,message)
}

//transaction will be sended to the mempool to the given ip address node
func (nm *NodeCommand) SendTransaction(tx blockchain.Transaction){
	var message TransactionMessage
	message.Init(AddressToBytes(NODEIP),&tx)
	SendData(*nm.hostAddress,message)
}

func (nm *NodeCommand) Mined(tx blockchain.Transaction){
	var message MinedMessage
	message.Init(AddressToBytes(NODEIP),&tx)
	SendData(*nm.hostAddress,message)
}

//you can split the return of bestHeight function in different nodes with goroutines and pass that in numhash parameter
func (nm *NodeCommand) GetHeaders(numhash, sb, eb []byte){
	var message GetHeadersMessage
	message.Init(AddressToBytes(*nm.hostAddress),[]byte{0x00,0x00,0x00,0x01},numhash,sb,eb)
	SendData(*nm.hostAddress,message)
}