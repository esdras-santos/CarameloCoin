package network

import (
	"gochain/utils"
)

//when you start a new node this message will be sended
//when that message is sended your node will receive a "blockchain" message
type GetBlockMessage struct{
	Command []byte
	SenderIp []byte
}

func (gb *GetBlockMessage) Init(senderIp []byte){
	gb.Command = []byte("getblock")
	gb.SenderIp = senderIp
}

func (gb GetBlockMessage) GetCommand() []byte{
	return gb.Command
}

func (gb GetBlockMessage) Serialize() []byte{
	result := utils.ToLittleEndian(gb.SenderIp,4)
	return result
}

func (gb *GetBlockMessage) Parse(data []byte){
	sip := utils.ToLittleEndian(data,4)
	gb.Init(sip)
}

