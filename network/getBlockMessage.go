package network

import (
	"gochain/utils"
)

//when you start a new node this message will be sended
//when that message is sended your node will receive a "blockchain" message
type GetBlockMessage struct{
	Command []byte
	BlockHash []byte
	SenderIp []byte
}

func (gb *GetBlockMessage) Init(hash []byte, senderIp []byte){
	gb.Command = []byte("getblock")
	gb.BlockHash = hash
	gb.SenderIp = senderIp
}

func (gb GetBlockMessage) GetCommand() []byte{
	return gb.Command
}

func (gb GetBlockMessage) Serialize() []byte{
	result := utils.ToLittleEndian(gb.BlockHash,len(gb.BlockHash))
	result = append(result, utils.ToLittleEndian(gb.SenderIp,len(gb.SenderIp))...)
	return result
}

func (gb *GetBlockMessage) Parse(data []byte){
	bh := utils.ToLittleEndian(data,len(data[:32]))
	sip := utils.ToLittleEndian(data[32:],len(data[32:]))
	gb.Init(bh,sip)
}

