package network

import (
	//"bytes"

	//"encoding/hex"
	"bufio"
	"encoding/hex"

	//"strings"

	"context"
	"crypto/rand"

	"fmt"
	"gochain/blockchain"
	"io"
	"log"
	"sync"

	// "os"
	// "runtime"

	// "syscall"

	"github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-core/host"
	"github.com/syndtr/goleveldb/leveldb"

	//"github.com/libp2p/go-libp2p-core/peer"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	net "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"

	ma "github.com/multiformats/go-multiaddr"
)


var MEMPOOL = make(map[string]blockchain.Transaction)
var mutex = &sync.Mutex{}
var Host host.Host
var MyId string




var Peerids = make(map[string]*bufio.ReadWriter)
var Message = make(chan NetworkEnvelope,120)

var MyAddress string

// new network with libp2p-pubsub

type Network struct{
	RW *bufio.ReadWriter
}

func Connect(genesisnode string) error{
	
	 
	Host, err := MakeHost()
	
	
	Host.SetStreamHandler("/p2p/1.0.0", HandleStream)

	ipfsaddr, err := ma.NewMultiaddr(genesisnode)
	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	Handle(err)
	peerid, err := peer.IDB58Decode(pid)
	Handle(err)
	tpa := fmt.Sprintf("/ipfs/%s", peerid.Pretty())
	var targetPeerAddr, _ = ma.NewMultiaddr(tpa)
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)
	
	Host.Peerstore().AddAddr(peerid,targetAddr,pstore.PermanentAddrTTL)
		
	s, err := Host.NewStream(context.Background(),peerid,"/p2p/1.0.0")
	Handle(err)
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	Peerids[pid] = rw
	MyId = pid
	go HandleMesssages(rw)
	go Publish()
	gb := GetBlockMessage{}
	gb.Init()
	ne := NetworkEnvelope{[]byte(MyId),gb.Command,nil}
	Message <- ne


	return err
}

func Listen(){
	log.Println("listening for connections")
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
	Host, err := MakeHost()
	Handle(err)
	Host.SetStreamHandler("/p2p/1.0.0", HandleStream)

	
}

func HandleStream(s net.Stream){
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	

	go HandleMesssages(rw)
	go Publish()
}

func MakeHost() (host.Host, error){
	var r io.Reader
	r = rand.Reader

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA,2048,r)
	Handle(err)

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
		libp2p.Identity(priv),
	}

	h, err := libp2p.New(context.Background(), opts...)
	Handle(err)
	
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", h.ID().Pretty()))
	

	pid, err := hostAddr.ValueForProtocol(ma.P_IPFS)
	Handle(err)
	
	

	peerid, err := peer.IDB58Decode(pid)
	MyId = string(peerid)
	Handle(err)
	
	addr := h.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	log.Printf("peerid: %s\n", peerid.Pretty())

	return h, nil
}


func Publish() {
	for{
		msg :=<- Message
	
		for _, v := range Peerids{
			mutex.Lock()
			bytes := msg.Serialize()

			mutex.Unlock()

			mutex.Lock()
		
			xbytes := fmt.Sprintf("%x\n",string(bytes))	
			
			v.WriteString(xbytes)

			v.Flush()
		
		
			mutex.Unlock()
		}
	
	}
	
	
}

func PublishToTarget(msg NetworkEnvelope, rw *bufio.ReadWriter) {
	mutex.Lock()
	bytes := msg.Serialize()

	mutex.Unlock()

	mutex.Lock()
		
	xbytes := fmt.Sprintf("%x\n",string(bytes))	
			
	_, err := rw.WriteString(xbytes)
	Handle(err)
	

	err = rw.Flush()
	Handle(err)
		
	mutex.Unlock()
	
}


func HandleMesssages(rw *bufio.ReadWriter){
	var msg NetworkEnvelope
	for{
		str, err := rw.ReadString('\n')
		data, _ := hex.DecodeString(str)
		Handle(err)
		if str == ""{
			return
		}
		if str != "\n"{
			m := msg.Parse(data)
			if _, e := Peerids[string(m.Peerid)]; !e{
				Peerids[string(m.Peerid)] = rw
			}
			mutex.Lock()
			switch string(m.Command){
			case "transaction":
				//just forward if is a valid transaction
				
				tm := TransactionMessage{} 
				tx := tm.Parse(m.Payload)
				if !blockchain.BlockchainInstance.VerifyTransaction(tx){
					fmt.Println("invalid transaction received")
					mutex.Unlock()
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
			case "getblock":
				bm := GenBlockMessage{}
				lh := blockchain.BlockchainInstance.GetLastHash()
				gb, err := blockchain.BlockchainInstance.GetBlock(lh)
				Handle(err)
				bm.Init([]byte(MyAddress),gb)
				ne := NetworkEnvelope{[]byte(MyId),[]byte("genblock"),bm.Serialize()}
				Message <- ne
			case "genblock":
				gb := GenBlockMessage{}
				block, mineraddr := gb.Parse(m.Payload)
				db, err := leveldb.OpenFile("./tmp/blocks",nil)
				Handle(err)
				err = db.Put(block.Hash(), block.Serialize(),nil)
				Handle(err)
				err = db.Put([]byte("lh"), block.Hash(),nil)
				Handle(err)
				
				accdb := blockchain.InitAccounts(string(mineraddr))
				
				db.Close()
				accdb.AccDatabase.Close()

			}
		
			
			mutex.Unlock()

		}
	}

}


//this function is just for tests forwhile
// func ConnectWithPeer() *bufio.ReadWriter{
// 	var r io.Reader
// 	r = rand.Reader

// 	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA,2048,r)
// 	Handle(err)

// 	opts := []libp2p.Option{
// 		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
// 		libp2p.Identity(priv),
// 	}

// 	host, err := libp2p.New(context.Background(), opts...)
// 	Handle(err)

// 	ipfsaddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/56915/ipfs/QmS4x9RtzLhqYDVbFRDzNAYVdGXtFhGPsmHVDbtGQbChdA")
// 	tpa := fmt.Sprintf("/ipfs/%s", GenesisPeerID)
// 	var targetPeerAddr, _ = ma.NewMultiaddr(tpa)
// 	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)
// 	peid, err := peer.IDB58Decode(GenesisPeerID)
// 	Handle(err)
// 	host.Peerstore().AddAddr(peid,targetAddr,pstore.PermanentAddrTTL)
		
// 	s, err := host.NewStream(context.Background(),peid,"/p2p/1.0.0")
// 	Handle(err)
// 	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
// 	return rw
// }

// func (nw *Network) ListPeers() []peer.ID{
// 	return nw.ps.ListPeers("caramelocoin")
// }

// func JoinNetwork(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID) (*Network, error){
// 	topic, err := ps.Join("caramelocoin")
// 	if err != nil{
// 		return nil, err
// 	}
// 	sub, err := topic.Subscribe()
// 	if err != nil{
// 		return nil, err
// 	}
// 	nw := &Network{
// 		messages: make(chan *NetworkEnvelope, 128),
// 		ctx: ctx,
// 		ps: ps,
// 		topic: topic,
// 		sub: sub,
// 	}

// 	go nw.readLoop()
// 	return nw, nil
// }


// func (nw *Network) readLoop(){
// 	go nw.HandleMesssages()
// 	for{
// 		msg, err := nw.sub.Next(nw.ctx)
// 		if err != nil {
// 			close(nw.messages)
// 			return
// 		}

// 		// only forward messages delivered by other nodes
// 		if msg.ReceivedFrom == nw.self{
// 			continue
// 		}

// 		ne := new(NetworkEnvelope)
// 		ne = ne.Parse(msg.Data)
		
// 		nw.messages <- ne
// 	}
// }

 
// type discoveryNotifee struct{
// 	h host.Host
// }

// func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo){
// 	fmt.Printf("\ndiscovered new peer %s \n",pi.ID.Pretty())
// 	err := n.h.Connect(context.Background(),pi)
// 	if err != nil {
// 		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
// 	}
// }

// func SetupDiscovery(ctx context.Context, h host.Host) error{
// 	disc := mdns.NewMdnsService(h,DiscoveryServiceTag)

// 	n := discoveryNotifee{h}
// 	disc.RegisterNotifee(&n)
// 	return nil
// }



// const (
// 	PORT = "9333"
// 	PROTOCOL = "tcp"
// 	VERSION  = 1
// 	COMMANDLENGTH = 12
// 	DiscoveryServiceTag = "caramelocoinnetwork"
// )

// //the ip of your node
// var NODEIP = utils.GetIp()

// //hex equivalenty to cmlc
// var NETWORK_MAGIC = []byte{0x63,0x6d,0x6c,0x63}

// //will check if the handshack was maded
// var VERACKRECEIVED map[string]bool

// var(
// 	minerAddress string
// 	//hard-coded first node ip
// 	KNOWNNODES = []string{"45.167.55.3:9333"}
// 	blocksInTransit = [][]byte{}
// 	MEMPOOL = make(map[string]blockchain.Transaction)
// )



// type Addr struct{
// 	AddrList []string
// }

// // type Block struct{
// // 	AddrFrom string
// // 	Block []byte
// // }

// type GetBlocks struct{
// 	AddrFrom string
// }

// type GetData struct{
// 	AddrFrom string
// 	Type string
// 	ID []byte
// }

// type Inv struct{
// 	AddrFrom string
// 	Type string
// 	Items [][]byte
// }

// type Tx struct{
// 	AddrFrom string
// 	Transaction []byte
// }

// func StartServer(){
// 	nodeAddress := fmt.Sprintf("127.0.0.1:%s",PORT)
	
// 	ln, err := net.Listen(PROTOCOL, nodeAddress)
// 	if err != nil{
// 		log.Panic(err)
// 	}
// 	defer ln.Close()
	
// 	println("breaking point")
// 	chain := blockchain.ContinueBlockChain("./tmp/blocks")
	
// 	defer chain.Database.Close()
// 	go CloseDB(chain)

// 	nc := NodeCommand{}
// 	if nodeAddress != KNOWNNODES[0]{
// 		nc.Init(KNOWNNODES[0])
// 		nc.HandShake()
// 	}
// 	for{
// 		conn, err := ln.Accept()
// 		if err != nil{
// 			log.Panic(err)
// 		}
// 		go HandleConnection(conn, chain)
// 	}
// }

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

// func CloseDB(chain *blockchain.BlockChain){
// 	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM,os.Interrupt)

// 	d.WaitForDeathWithFunc(func(){
// 		defer os.Exit(1)
// 		defer runtime.Goexit()
// 		chain.Database.Close()
// 	})
// }



