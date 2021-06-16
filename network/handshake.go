package network

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
)


type HandShakeNode struct {
	//the ip and port of the host
	HostAddress string
	Conn        net.Conn
	Stream      []byte
}

//this need to started in a goroutine because of infity loop
func (hsn *HandShakeNode) Init(minerAdrress, host string) {
	node := fmt.Sprintf("%s%s", minerAddress, PORT)

	hsn.HostAddress = fmt.Sprintf("%s%s", host, PORT)
	var err error
	hsn.Conn, err = net.Dial(PROTOCOL, hsn.HostAddress)
	Handle(err)
	ln, err := net.Listen(PROTOCOL, node)
	Handle(err)
	for {
		conn, err := ln.Accept()
		defer conn.Close()
		Handle(err)
		hsn.Stream, err = ioutil.ReadAll(conn)
		Handle(err)
		break
	}
}

//send handshake
func (hsn *HandShakeNode) Send(message VersionMessage) {
	envelope := NetworkEnvelope{NETWORK_MAGIC, []byte("version"), message.Serialize()}
	_, err := io.Copy(hsn.Conn, bytes.NewReader(envelope.Serialize()))
	Handle(err)
}

//read the response of the handshake
func (hsn *HandShakeNode) Read() NetworkEnvelope {
	var ne NetworkEnvelope
	envelope := ne.Parse(hsn.Stream)
	return *envelope
}

func (hsn *HandShakeNode) WaitFor() {

}