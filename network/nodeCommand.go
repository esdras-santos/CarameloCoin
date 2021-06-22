package network

import (
	"fmt"
)

var bestHeight = chain.GetBestHeight()

var version = [4]byte{0x00000001}
var services = [8]byte{0x00000000000000}
var recServ = [8]byte{0x00000000000000}
var recIp = []byte(addr)
var recPort = []byte(PORT)
var sendServ = [8]byte{0x00000000000000}
// getIp() function should return the current Ip of the node in []byte
var sendIp = getIp()
var sendPort = []byte(PORT)
var nonce = [8]byte{0x6e,0x6f,0x74,0x20,0x6d,0x65,0x21,0x21}
var userAge = []byte("/CarameloCoin:0.1/")
var lateBlock = utils.ToHex(int64(bestHeight))
var rel = []byte{0x01}

type NodeCommand struct {
	//the ip and port of the host
	HostAddress string
}

//this need to started in a goroutine because of infity loop
func (nm *NodeCommand) Init(host string) {
	nm.HostAddress = fmt.Sprintf("%s%s", host, PORT)
}

func (nm *NodeCommand) HandShake(){
	var message VersionMessage
	message.Init(version[:],services[:],nil,recServ[:],recIp,recPort,sendServ[:],sendIp,sendPort,nonce[:],userAge,lateBlock,rel)
		
	SendData(nm.HostAddress,message)
}

func (nm *NodeCommand) GetHeaders(numhash, sb, eb []byte){
	var message GetHeadersMessage
	message.Init([]byte{0x00,0x00,0x00,0x01},numhash,sb,eb)
	SendData(nm.HostAddress,message)
}