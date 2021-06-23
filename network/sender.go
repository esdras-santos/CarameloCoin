package network

import (
	"bytes"
	"fmt"
	"gochain/blockchain"
	"gochain/utils"
	"io"
	"log"
	"net"
)



type Message interface{
	GetCommand() []byte
	Serialize() []byte
}

func SendAddr(address string) {
	nodes := Addr{KnownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := GobEncode(nodes)
	request := append(CmdToBytes("addr"), payload...)

	SendData(address, request)
}

func SendBlock(addr string, b *blockchain.Block) {
	data := Block{nodeAddress, b.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("block"), payload...)

	SendData(addr, request)
}

func SendInv(address, kind string, items [][]byte) {
	inventory := Inv{nodeAddress, kind, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes("inv"), payload...)

	SendData(address, request)
}

func SendTx(addr string, tnx *blockchain.Transaction) {
	data := Tx{nodeAddress, tnx.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	SendData(addr, request)
}

// func SendVersion(addr string, chain *blockchain.BlockChain) {
// 	bestHeight := chain.GetBestHeight()
// 	message := VersionMessage{}
// 	version := [4]byte{0x00000001}
// 	services := [8]byte{0x00000000000000}
// 	recServ := [8]byte{0x00000000000000}
// 	recIp := []byte(addr)
// 	recPort := []byte(PORT)
// 	sendServ := [8]byte{0x00000000000000}
// 	// getIp() function should return the current Ip of the node in []byte
// 	sendIp := getIp()
// 	sendPort := []byte(PORT)
// 	nonce := [8]byte{0x6e,0x6f,0x74,0x20,0x6d,0x65,0x21,0x21}
// 	userAge := []byte("/CarameloCoin:0.1/")
// 	lateBlock := utils.ToHex(int64(bestHeight))
// 	rel := []byte{0x01}

// 	message.Init(version[:],services[:],nil,recServ[:],recIp,recPort,sendServ[:],sendIp,sendPort,nonce[:],userAge,lateBlock,rel)
 
// 	ne := NetworkEnvelope{NETWORK_MAGIC,[]byte("version"),message.Serialize()}
// 	SendData(addr, ne.Serialize())
// }

// func SendVerAck(address string){
// 	ne := NetworkEnvelope{NETWORK_MAGIC,[]byte("verack"),nil}
// 	SendData(address,ne.Serialize())
// }

// func SendGetBlocks(address string) {
// 	payload := GobEncode(GetBlocks{nodeAddress})
// 	request := append(CmdToBytes("getblocks"), payload...)

// 	SendData(address, request)
// }

// func SendGetData(address, kind string, id []byte) {
// 	payload := GobEncode(GetData{nodeAddress, kind, id})
// 	request := append(CmdToBytes("getdata"), payload...)

// 	SendData(address, request)
// }


func SendData(hostAddr string, message Message) {
	envelope := NetworkEnvelope{NETWORK_MAGIC, message.GetCommand(), message.Serialize()}
	host := fmt.Sprintf("%s:%s",hostAddr,PORT)
	conn, err := net.Dial(PROTOCOL, host)

	if err != nil {
		fmt.Printf("%s is not available\n", host)
		var updatedNodes []string

		for _, node := range KnownNodes {
			if node != host {
				updatedNodes = append(updatedNodes, node)
			}
		}

		KnownNodes = updatedNodes

		return
	}

	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(envelope.Serialize()))
	if err != nil {
		log.Panic(err)
	}
}
