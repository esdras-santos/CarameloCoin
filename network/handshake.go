package network

import (
	"bytes"
	"fmt"
	"gochain/utils"
	"io"
	"io/ioutil"
	"net"
)

const VERACKMESSAGE = "verack"

type HandShakeNode struct {
	//the ip and port of the host
	HostAddress string
	Conn        net.Conn
	Stream      []byte
}

var bestHeight = chain.GetBestHeight()
var message = VersionMessage{}
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



//this need to started in a goroutine because of infity loop
func (hsn *HandShakeNode) Init(minerAdrress, host string) {
	hsn.HostAddress = fmt.Sprintf("%s%s", host, PORT)
	var err error
	hsn.Conn, err = net.Dial(PROTOCOL, hsn.HostAddress)
	Handle(err)
	
}

//send handshake
func (hsn *HandShakeNode) Send(message VersionMessage, handmess string) {
	envelope := NetworkEnvelope{NETWORK_MAGIC, []byte(handmess), message.Serialize()}
	_, err := io.Copy(hsn.Conn, bytes.NewReader(envelope.Serialize()))
	Handle(err)
}

func (hs *HandShakeNode) HandShake(){
	message.Init(version[:],services[:],nil,recServ[:],recIp,recPort,sendServ[:],sendIp,sendPort,nonce[:],userAge,lateBlock,rel)
		
	hs.Send(message,"version")
}