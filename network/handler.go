package network

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gochain/blockchain"
	"io/ioutil"
	"log"
	"net"
)


func HandleAddr(request []byte) {
	var buff bytes.Buffer
	var payload Addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	KnownNodes = append(KnownNodes, payload.AddrList...)
	fmt.Printf("there are %d known nodes\n", len(KnownNodes))
	RequestBlocks()
}

func HandleInv(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s \n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items
		blockHash := payload.Items[0]

		SendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if memoryPool[hex.EncodeToString(txID)].ID == nil {
			SendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func HandleBlock(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := blockchain.Deserialize(blockData)

	fmt.Println("Recevied a new block!")
	chain.AddBlock(block)

	fmt.Printf("added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		SendGetData(payload.AddrFrom, "block", blockHash)
		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{chain}
		UTXOSet.Reindex()
	}
}

func HandleGetBlocks(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := chain.GetBlockHashes()
	SendInv(payload.AddrFrom, "block", blocks)
}

func HandleGetData(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload GetData

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := chain.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		SendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := memoryPool[txID]

		SendTx(payload.AddrFrom, &tx)
	}
}

func HandleVersion(request []byte, chain *blockchain.BlockChain) {
	var payload VersionMessage

	payload.Parse(request[COMMANDLENGTH+4:])
		
	SendVerAck(string(payload.SenderIp))
	SendVersion(string(payload.SenderIp), chain)
	

	if !NodeIsKnown(string(payload.SenderIp)) {
		KnownNodes = append(KnownNodes, string(payload.SenderIp))
	}
}

func HandleVerAck(request []byte){
	var payload VersionMessage

	payload.Parse(request[COMMANDLENGTH+4:])

	//handshack maded
	VERACKRECEIVED[string(payload.SenderIp)] = true
}

func HandleTx(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Tx

	buff.Write(request[command:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := blockchain.DeserializeTransaction(txData)
	memoryPool[hex.EncodeToString(tx.ID)] = tx

	fmt.Printf("%s, %d\n", nodeAddress, len(memoryPool))

	if nodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != nodeAddress && node != payload.AddrFrom {
				SendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(memoryPool) >= 2 && len(minerAddress) > 0 {
			MineTx(chain)
		}
	}
}

func HandleConnection(conn net.Conn, chain *blockchain.BlockChain) {
	req, err := ioutil.ReadAll(conn)
	defer conn.Close()
	connectedNode := conn.RemoteAddr().String()
	if err != nil {
		log.Panic(err)
	}
	command := BytesToCmd(req[4:COMMANDLENGTH+4])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		if VERACKRECEIVED[connectedNode]{
			HandleAddr(req)
		}else{
			log.Panic("you don't made the handshake")
		}
	case "block":
		if VERACKRECEIVED[connectedNode]{
			HandleBlock(req, chain)
		}else{
			log.Panic("you don't made the handshake")
		}		
	case "inv":
		if VERACKRECEIVED[connectedNode]{
			HandleInv(req, chain)
		}else{
			log.Panic("you don't made the handshake")
		}
	case "getblocks":
		if VERACKRECEIVED[connectedNode]{
			HandleGetBlocks(req, chain)
		}else{
			log.Panic("you don't made the handshake")
		}	
	case "getheaders":
		if VERACKRECEIVED[connectedNode]{
			//this need return all the block headers asked with a headers command
			HandleGetHeaders(req, chain)
		}else{
			log.Panic("you don't made the handshake")
		}
	case "getdata":
		if VERACKRECEIVED[connectedNode]{
			HandleGetData(req, chain)
		}else{
			log.Panic("you don't made the handshake")
		}	
	case "tx":
		if VERACKRECEIVED[connectedNode]{
			HandleTx(req, chain)
		}else{
			log.Panic("you don't made the handshake")
		}		
	case "version":
		HandleVersion(req, chain)
	case "verack":
		HandleVerAck(req)
	default:
		fmt.Println("Unknown command")
	}
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}