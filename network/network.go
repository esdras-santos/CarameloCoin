package network

import (
	//"bytes"

	//"encoding/hex"
	"context"
	"fmt"
	"gochain/blockchain"
	"gochain/utils"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"syscall"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"gopkg.in/vrecan/death.v3"
)

const (
	PORT = "9333"
	PROTOCOL = "tcp"
	VERSION  = 1
	COMMANDLENGTH = 12
	DiscoveryServiceTag = "caramelocoinnetwork"
)

//the ip of your node
var NODEIP = utils.GetIp()

//hex equivalenty to cmlc
var NETWORK_MAGIC = []byte{0x63,0x6d,0x6c,0x63}

//will check if the handshack was maded
var VERACKRECEIVED map[string]bool

var(
	minerAddress string
	//hard-coded first node ip
	KNOWNNODES = []string{"45.167.55.3:9333"}
	blocksInTransit = [][]byte{}
	MEMPOOL = make(map[string]blockchain.Transaction)
)



type Addr struct{
	AddrList []string
}

// type Block struct{
// 	AddrFrom string
// 	Block []byte
// }

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

func StartServer(){
	nodeAddress := fmt.Sprintf("127.0.0.1:%s",PORT)
	
	ln, err := net.Listen(PROTOCOL, nodeAddress)
	if err != nil{
		log.Panic(err)
	}
	defer ln.Close()
	
	println("breaking point")
	chain := blockchain.ContinueBlockChain("./tmp/blocks")
	
	defer chain.Database.Close()
	go CloseDB(chain)

	nc := NodeCommand{}
	if nodeAddress != KNOWNNODES[0]{
		nc.Init(KNOWNNODES[0])
		nc.HandShake()
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
// func CmdToBytes(cmd string) []byte{
// 	var bytes [commandLength]byte
// 	for i,c := range cmd{
// 		bytes[i] = byte(c)
// 	}
// 	return bytes[:]
// }

//this function need to be reveewed later
// func BytesToCmd(bytes []byte) string{
// 	var cmd []byte

// 	for _,b := range bytes{
// 		if b != 0x0{
// 			cmd = append(cmd,b)
// 		}
// 	}

// 	return fmt.Sprintf("%s",cmd)
// }

// func RequestBlocks(){
// 	for _,node := range KnownNodes{
// 		SendGetBlocks(node)
// 	}
// }

// func ExtractCmd(request []byte) []byte{
// 	return request[:commandLength]
// }

// func MineTx(chain *blockchain.BlockChain){
// 	var txs []*blockchain.Transaction

// 	for id := range MEMPOOL{
// 		fmt.Printf("tx: %s\n",MEMPOOL[id].ID)
// 		tx := MEMPOOL[id]
// 		if chain.VerifyTransaction(&tx){
// 			txs = append(txs, &tx)
// 		}
// 	}

// 	if len(txs) == 0{
// 		fmt.Println("All Transactions are invalid")
// 		return
// 	}

// 	cbTx := blockchain.CoinbaseTx(minerAddress, "")
// 	txs = append(txs, cbTx)

// 	newBlock := chain.MineBlock(txs)
// 	UTXOSet := blockchain.UTXOSet{chain}
// 	UTXOSet.Reindex()

// 	fmt.Println("New Block mined")
// 	 for _,tx := range txs{
// 		 txID := hex.EncodeToString(tx.ID)
// 		 delete(MEMPOOL,txID)
// 	 }

// 	 for _,node := range KnownNodes{
// 		 if node != nodeAddress{
// 			 SendInv(node, "block", [][]byte{newBlock.Hash})
// 		 }
// 	 }

// 	 if len(MEMPOOL) > 0{
// 		 MineTx(chain)
// 	 }
// }


// func NodeIsKnown(addr string) bool{
// 	for _,node := range KnownNodes{
// 		if node == addr{
// 			return true
// 		}
// 	}
// 	return false
// }

func CloseDB(chain *blockchain.BlockChain){
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM,os.Interrupt)

	d.WaitForDeathWithFunc(func(){
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}



// new network with libp2p-pubsub

type Network struct{
	messages chan *NetworkEnvelope

	ctx 	context.Context
	ps 		*pubsub.PubSub
	topic 	*pubsub.Topic
	sub		*pubsub.Subscription

	self	peer.ID
}

func JoinNetwork(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID) (*Network, error){
	topic, err := ps.Join("caramelocoin")
	if err != nil{
		return nil, err
	}
	sub, err := topic.Subscribe()
	if err != nil{
		return nil, err
	}
	nw := &Network{
		messages: make(chan *NetworkEnvelope, 128),
		ctx: ctx,
		ps: ps,
		topic: topic,
		sub: sub,
	}

	go nw.readLoop()
	return nw, nil
}

func (nw *Network) Publish(message NetworkEnvelope) error{
	msgBytes := message.Serialize()
	return nw.topic.Publish(nw.ctx, msgBytes)
}

func (nw *Network) ListPeers() []peer.ID{
	return nw.ps.ListPeers("caramelocoin")
}

func (nw *Network) readLoop(){
	go nw.HandleMesssages()
	for{
		msg, err := nw.sub.Next(nw.ctx)
		if err != nil {
			close(nw.messages)
			return
		}

		// only forward messages delivered by other nodes
		if msg.ReceivedFrom == nw.self{
			continue
		}

		ne := new(NetworkEnvelope)
		ne = ne.Parse(msg.Data)
		
		nw.messages <- ne
	}
}

func (nw *Network) HandleMesssages(){

	for{
		m := <- nw.messages
		switch string(m.Command){
		case "transaction":
			//just forward if is a valid transaction
			tm := TransactionMessage{} 
			tx := tm.Parse(m.Payload)
			if !blockchain.BlockchainInstance.VerifyTransaction(tx){
				continue
			}
			HandleTx(tx)
		case "mined":
			tm := MinedMessage{} 
			tx := tm.Parse(m.Payload)
			HandleMined(tx)
		case "block":
			bm := BlockMessage{}
			block := bm.Parse(m.Payload)
			HandleBlock(block)
		}
	
	}

}

type discoveryNotifee struct{
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo){
	fmt.Printf("discovered new peer %s \n",pi.ID.Pretty())
	err := n.h.Connect(context.Background(),pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}

func SetupDiscovery(ctx context.Context, h host.Host) error{
	disc := mdns.NewMdnsService(h, DiscoveryServiceTag)

	n := discoveryNotifee{h}
	disc.RegisterNotifee(&n)
	return nil
}
