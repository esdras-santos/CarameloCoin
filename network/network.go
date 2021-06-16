package network

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"gochain/blockchain"
	"log"
	"net"
	"os"
	"runtime"
	"syscall"

	"gopkg.in/vrecan/death.v3"
)

const (
	PORT = ":9333"
	PROTOCOL = "tcp"
	VERSION  = 1
	COMMANDLENGTH = 12
)
//hex equivalenty to cmlc
var NETWORK_MAGIC = []byte{0x63,0x6d,0x6c,0x63}

var(
	nodeAddress string
	minerAddress string
	KnownNodes = []string{"localhost:9333"}
	blocksInTransit = [][]byte{}
	memoryPool = make(map[string]blockchain.Transaction)
)

type Addr struct{
	AddrList []string
}

type Block struct{
	AddrFrom string
	Block []byte
}

type GetBlocks struct{
	AddrFrom string
}

type GetData struct{
	AddrFrom string
	Type string
	ID []byte
}

type Inv struct{
	AddrFrom string
	Type string
	Items [][]byte
}

type Tx struct{
	AddrFrom string
	Transaction []byte
}

func StartServer(nodeID, minerAddress string){
	nodeAddress = fmt.Sprintf("%s%s",minerAddress,PORT)
	ln, err := net.Listen(PROTOCOL, nodeAddress)
	if err != nil{
		log.Panic(err)
	}
	defer ln.Close()

	chain := blockchain.ContinueBlockChain(nodeID)
	defer chain.Database.Close()
	go CloseDB(chain)

	if nodeAddress != KnownNodes[0]{
		SendVersion(KnownNodes[0], chain)
	}
	for{
		conn, err := ln.Accept()
		if err != nil{
			log.Panic(err)
		}
		go HandleConnection(conn, chain)
	}
}

//this function need to be reviewed later
func CmdToBytes(cmd string) []byte{
	var bytes [commandLength]byte
	for i,c := range cmd{
		bytes[i] = byte(c)
	}
	return bytes[:]
}

//this function need to be reveewed later
func BytesToCmd(bytes []byte) string{
	var cmd []byte

	for _,b := range bytes{
		if b != 0x0{
			cmd = append(cmd,b)
		}
	}

	return fmt.Sprintf("%s",cmd)
}

func RequestBlocks(){
	for _,node := range KnownNodes{
		SendGetBlocks(node)
	}
}

func ExtractCmd(request []byte) []byte{
	return request[:commandLength]
}

func MineTx(chain *blockchain.BlockChain){
	var txs []*blockchain.Transaction

	for id := range memoryPool{
		fmt.Printf("tx: %s\n",memoryPool[id].ID)
		tx := memoryPool[id]
		if chain.VerifyTransaction(&tx){
			txs = append(txs, &tx)
		}
	}

	if len(txs) == 0{
		fmt.Println("All Transactions are invalid")
		return
	}

	cbTx := blockchain.CoinbaseTx(minerAddress, "")
	txs = append(txs, cbTx)

	newBlock := chain.MineBlock(txs)
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	fmt.Println("New Block mined")
	 for _,tx := range txs{
		 txID := hex.EncodeToString(tx.ID)
		 delete(memoryPool,txID)
	 }

	 for _,node := range KnownNodes{
		 if node != nodeAddress{
			 SendInv(node, "block", [][]byte{newBlock.Hash})
		 }
	 }

	 if len(memoryPool) > 0{
		 MineTx(chain)
	 }
}

func GobEncode(data interface{}) []byte{
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil{
		log.Panic(err)
	}

	return buff.Bytes()
}

func NodeIsKnown(addr string) bool{
	for _,node := range KnownNodes{
		if node == addr{
			return true
		}
	}
	return false
}

func CloseDB(chain *blockchain.BlockChain){
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM,os.Interrupt)

	d.WaitForDeathWithFunc(func(){
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}


