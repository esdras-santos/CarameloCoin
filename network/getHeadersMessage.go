package network

import (
	"gochain/utils"
	"log")


type GetHeadersMessage struct {
	Command       []byte
	Version       []byte
	NumberOfHashs []byte
	StartingBlock []byte
	EndingBlock   []byte
}

func (gh *GetHeadersMessage) Init(version, numHashs, startBlock, endBlock []byte) {
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

func (gh *GetHeadersMessage) Serialize() []byte {
	result := utils.ToLittleEndian(gh.Version,4)
	result = append(result, gh.NumberOfHashs...)
	result = append(result, utils.ToLittleEndian(gh.StartingBlock,32)...)
	result = append(result, utils.ToLittleEndian(gh.EndingBlock,32)...)
	return result
}

func (gh *GetHeadersMessage) Parse(data []byte){
	version := utils.ToLittleEndian(data[:4],4)
	numHash := data[5:6]
	sb := utils.ToLittleEndian(data[6:38],32)
	eb := utils.ToLittleEndian(data[38:70],32)
	gh.Init(version,numHash,sb,eb)
}