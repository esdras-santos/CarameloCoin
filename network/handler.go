package network

import (
	"encoding/hex"
	"fmt"
	"gochain/blockchain"
	
	"log"
	
)

// func HandleAddr(request []byte) {
// 	var buff bytes.Buffer
// 	var payload Addr

// 	buff.Write(request[COMMANDLENGTH+4:])
// 	dec := gob.NewDecoder(&buff)
// 	err := dec.Decode(&payload)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	KNOWNNODES = append(KNOWNNODES, payload.AddrList...)
// 	fmt.Printf("there are %d known nodes\n", len(KNOWNNODES))
// 	RequestBlocks()
// }

// func NodeIsKnown(address string)bool{
// 	for _, item := range KNOWNNODES {
//         if item == address {
//             return true
//         }
//     }
//     return false
// }

// func HandleInv(request []byte, chain *blockchain.BlockChain) {
// 	var buff bytes.Buffer
// 	var payload Inv

// 	buff.Write(request[commandLength:])
// 	dec := gob.NewDecoder(&buff)
// 	err := dec.Decode(&payload)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	fmt.Printf("Recevied inventory with %d %s \n", len(payload.Items), payload.Type)

// 	if payload.Type == "block" {
// 		blocksInTransit = payload.Items
// 		blockHash := payload.Items[0]

// 		SendGetData(payload.AddrFrom, "block", blockHash)

// 		newInTransit := [][]byte{}
// 		for _, b := range blocksInTransit {
// 			if bytes.Compare(b, blockHash) != 0 {
// 				newInTransit = append(newInTransit, b)
// 			}
// 		}
// 		blocksInTransit = newInTransit
// 	}

// 	if payload.Type == "tx" {
// 		txID := payload.Items[0]

// 		if memoryPool[hex.EncodeToString(txID)].ID == nil {
// 			SendGetData(payload.AddrFrom, "tx", txID)
// 		}
// 	}
// }

// func HandleBlock(request []byte, chain *blockchain.BlockChain) {
// 	var buff bytes.Buffer
// 	var payload Block

// 	buff.Write(request[commandLength:])
// 	dec := gob.NewDecoder(&buff)
// 	err := dec.Decode(&payload)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	blockData := payload.Block
// 	block := blockchain.Deserialize(blockData)

// 	fmt.Println("Recevied a new block!")
// 	chain.AddBlock(block)

// 	fmt.Printf("added block %x\n", block.Hash)

// 	if len(blocksInTransit) > 0 {
// 		blockHash := blocksInTransit[0]
// 		SendGetData(payload.AddrFrom, "block", blockHash)
// 		blocksInTransit = blocksInTransit[1:]
// 	} else {
// 		UTXOSet := blockchain.UTXOSet{chain}
// 		UTXOSet.Reindex()
// 	}
// }

func HandleBlock(block *blockchain.Block){
	chain := blockchain.BlockchainInstance
	for _,tx := range block.Transactions{
		if chain.VerifyTransaction(tx) != true{
			log.Panic("Invalid Transaction")
		}
		_, ok := MEMPOOL[hex.EncodeToString(tx.Id())];
    	if ok {
        	delete(MEMPOOL, hex.EncodeToString(tx.Id()));
    	}
	}
	chain.AddBlock(block)
}


//this function need to be reajusted to get blocks in distributed way
// func HandleGetBlock(request []byte, chain *blockchain.BlockChain) {
//  	var payload GetBlockMessage

//  	payload.Parse(request[COMMANDLENGTH+4:])
// 	blockhashs := chain.GetBlockHashes()
// 	for _,h := range blockhashs{
// 		bm := BlockMessage{}
// 		block, err := chain.GetBlock(h)
// 		Handle(err)
// 		bm.Init(&block)

//  		SendData(string(payload.SenderIp), bm)
// 	}	
// }


//request for headers
//get all the hashs in the DB from the startBlock to the endBlock
//put all the blocks as argument in HeadersMessage struct
//send the HeadersMessage
// func HandleGetHeaders(request []byte, chain *blockchain.BlockChain) {
// 	var payload GetHeadersMessage

// 	payload.Parse(request[COMMANDLENGTH+4:])
	
// 	blockHeaders := chain.GetBlockHeaders(payload.StartingBlock,payload.EndingBlock)
// 	hm := HeadersMessage{}
// 	hm.Init(blockHeaders)
// 	SendData(string(payload.SenderIp),hm)
// }

//response for the getheaders command
//receive the headers and add to the database 
// func HandleHeaders(request []byte, chain *blockchain.BlockChain) {
// 	var payload HeadersMessage

// 	payload.Parse(request[COMMANDLENGTH+4:])

	
// }


func HandleTx(tx *blockchain.Transaction) {

	MEMPOOL[hex.EncodeToString(tx.Id())] = *tx
	print("\ntransaction added")

	fmt.Printf("127.0.0.1, %d\n", len(MEMPOOL))
}

func HandleMined(tx *blockchain.Transaction){

	_, ok := MEMPOOL[hex.EncodeToString(tx.Id())];
    if ok {
        delete(MEMPOOL, hex.EncodeToString(tx.Id()));
    }
	// this have to be the ip of the node
	fmt.Printf("127.0.0.1, %d\n", len(MEMPOOL))
}

// func HandleConnection(conn net.Conn, chain *blockchain.BlockChain) {
// 	req, err := ioutil.ReadAll(conn)
// 	defer conn.Close()
// 	connectedNode := conn.RemoteAddr().String()
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	command := string(req[4:COMMANDLENGTH+4])
// 	fmt.Printf("Received %s command\n", command)

// 	switch command {
// 	// case "addr":
// 	// 	if VERACKRECEIVED[connectedNode]{
// 	// 		HandleAddr(req)
// 	// 	}else{
// 	// 		log.Panic("you don't made the handshake")
// 	// 	}
// 	//"block" is a response to "getblock" command
// 	case "block":
// 	 	if VERACKRECEIVED[connectedNode]{
// 	 		HandleBlock(req, chain)
// 	 	}else{
// 	 		log.Panic("you don't made the handshake")
// 	 	}		
// 	// case "inv":
// 	// 	if VERACKRECEIVED[connectedNode]{
// 	// 		HandleInv(req, chain)
// 	// 	}else{
// 	// 		log.Panic("you don't made the handshake")
// 	// 	}
// 	//with this command you will receive a "block" command
// 	case "getblock":
// 	 	if VERACKRECEIVED[connectedNode]{
// 	 		HandleGetBlock(req, chain)
// 	 	}else{
// 	 		log.Panic("you don't made the handshake")
// 	 	}	
// 	//request headers
// 	case "getheaders":
// 		if VERACKRECEIVED[connectedNode]{
// 			//this need return all the block headers asked with a headers command
// 			HandleGetHeaders(req, chain)
// 		}else{
// 			log.Panic("you don't made the handshake")
// 		}
// 	//response of getheaders command
// 	case "headers":
// 		if VERACKRECEIVED[connectedNode]{
// 			//this need return all the block headers asked with a headers command
// 			HandleHeaders(req, chain)
// 		}else{
// 			log.Panic("you don't made the handshake")
// 		}
// 	// case "getdata":
// 	// 	if VERACKRECEIVED[connectedNode]{
// 	// 		HandleGetData(req, chain)
// 	// 	}else{
// 	// 		log.Panic("you don't made the handshake")
// 	// 	}	
// 	case "transaction":
// 	 	if VERACKRECEIVED[connectedNode]{
// 	 		HandleTx(req, chain)
// 	 	}else{
// 	 		log.Panic("you don't made the handshake")
// 	 	}
// 	case "mined":
// 		if VERACKRECEIVED[connectedNode]{
// 			HandleMined(req, chain)
// 		}else{
// 			log.Panic("you don't made the handshake")
// 		}	 		
// 	case "version":
// 		HandleVersion(req, chain)
// 	case "verack":
// 		HandleVerAck(req)
// 	default:
// 		fmt.Println("Unknown command")
// 	}
// }

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}