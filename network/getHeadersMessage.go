package network

import (
	"gochain/utils"
	"log"
)


type GetHeadersMessage struct {
	Command       []byte
	SenderIp	  []byte
	Version       []byte
	NumberOfHashs []byte
	StartingBlock []byte
	EndingBlock   []byte
}

func (gh *GetHeadersMessage) Init(senderIp, version, numHashs, startBlock, endBlock []byte) {
	gh.SenderIp = senderIp
	gh.Command = []byte("getheaders")
	gh.Version = version
	gh.NumberOfHashs = numHashs
	if startBlock == nil {
		log.Panic("a start block is required")
	}
	gh.StartingBlock = startBlock
	if endBlock == nil{
		var h []byte
		for i:=0;i<32;i++{
			h = append(h, 0x00)
		}
		gh.EndingBlock = h
	}else{
		gh.EndingBlock = endBlock
	}
}

func (gh GetHeadersMessage) GetCommand() []byte{
	return gh.Command
}

func (gh GetHeadersMessage) Serialize() []byte {
	result := utils.ToLittleEndian(gh.SenderIp) 
	result = append(result, utils.ToLittleEndian(gh.Version)...)
	result = append(result, gh.NumberOfHashs...)
	result = append(result, utils.ToLittleEndian(gh.StartingBlock)...)
	result = append(result, utils.ToLittleEndian(gh.EndingBlock)...)
	return result
}

func (gh *GetHeadersMessage) Parse(data []byte){
	senderIp := utils.ToLittleEndian(data[:4])
	version := utils.ToLittleEndian(data[4:8])
	numHash := data[8:9]
	sb := utils.ToLittleEndian(data[9:41])
	eb := utils.ToLittleEndian(data[41:73])
	gh.Init(senderIp,version,numHash,sb,eb)
}