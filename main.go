package main

import (
	"bufio"
	"encoding/hex"
	"time"

	"fmt"
	"gochain/blockchain"
	"strings"


	"gochain/network"
	"gochain/wallet"
	"log"
	"strconv"

	"os"
)




func main(){

	
	if len(os.Args) == 1{
		fmt.Println("gen node")
	 	var w *wallet.Wallet
	 	var chain *blockchain.BlockChain
	 	if !walletExists("./tmp/wallet.data"){
	 		w = createWallet()
	 	} else {
	 	 	w = enterWallet()
	 	}

	 	if !DBexists("./tmp/blocks") {
	 	 	chain = blockchain.InitBlockChain(w,"./tmp/blocks")	
	 		chain.Database.Close()
	 		chain.Acc.AccDatabase.Close()
	 	} 
	 	chain = blockchain.ContinueBlockChain("./tmp/blocks")
		

	 	network.MyAddress = string(w.Address())

	 	reader := bufio.NewReader(os.Stdin)
	 	go network.Listen()
	 	time.Sleep(time.Second * 3)
	 	fmt.Println("\naddress: ",string(w.Address()))
	 	fmt.Println("\ntype \"command\" to see the list of commands")
	 	for{
	 		fmt.Print("\n>>> ")
	 		com, err := reader.ReadString('\n')
	 		com = strings.Replace(com, "\n", "", -1)
	 		com = strings.Replace(com, "\r", "", -1)

	 		Handler(err)
	 		switch com{
	 		case "send":
	 			send(w,chain)
	 		case "mine":
	 		 	mine(w,chain)
	 		case "balance":
	 		 	balance(chain)
	 		case "command":
	 		 	command()
	 		default:
	 		 	command()
	 		}
	 	}

		
	 } else{
	 	var w *wallet.Wallet
	 	var chain *blockchain.BlockChain
	 	if !walletExists("./tmp/wallet.data"){
	 		w = createWallet()
	 	} else {
	 		w = enterWallet()
	 	}

	 	network.MyAddress = string(w.Address())
	 	network.Connect(os.Args[1])

	 	time.Sleep(time.Second * 10)

	 	chain = blockchain.ContinueBlockChain("./tmp/blocks")

	 	reader := bufio.NewReader(os.Stdin)
	 	fmt.Println("\naddress: ",network.MyAddress)
	 	fmt.Println("\ntype \"command\" to see the list of commands")
	 	for{
	 		fmt.Print("\n>>> ")
	 		com, err := reader.ReadString('\n')
	 		com = strings.Replace(com, "\n", "", -1)
	 		com = strings.Replace(com, "\r", "", -1)
	 		Handler(err)
	 		switch com{
	 		case "send":
	 			send(w,chain)
	 		case "mine":
	 			mine(w,chain)
	 		case "balance":
	 			balance(chain)
	 		case "command":
	 			command()
	 		}
	 	}
	}
}	

func command(){
	fmt.Println("-------commands list---------")
	fmt.Println("send - send a transation to")
	fmt.Println("mine - mine a block")
	fmt.Println("balance - see the balance of")
}

func balance(chain *blockchain.BlockChain){
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n	Of: ")
	of, err := reader.ReadString('\n')
	of = strings.Replace(of, "\n", "", -1)
	of = strings.Replace(of, "\r", "", -1)
	Handler(err)
	balance, _ := chain.Acc.BalanceNonce(of)
	fmt.Printf("\n	balance: %d\n", balance)
}

func mine(w *wallet.Wallet, chain *blockchain.BlockChain){
	cbTx := blockchain.CoinbaseTx(w)
	for range network.MEMPOOL{
			//collecting fees
		cbTx.Value += 1
	}
		
	txs := []*blockchain.Transaction{cbTx}

	for _,t := range network.MEMPOOL{
		txs = append(txs, &t)
	}

	bm := network.BlockMessage{}

	block := chain.MineBlock(txs)
	bm.Init(*block)
	ne := network.NetworkEnvelope{[]byte(network.MyId),bm.GetCommand(),bm.Serialize()}
	network.Message <- ne
}

func send(w *wallet.Wallet, chain *blockchain.BlockChain){
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n	To: ")
	to, err := reader.ReadString('\n')
	to = strings.Replace(to, "\n", "", -1)
	to = strings.Replace(to, "\r", "", -1)
	Handler(err)
	fmt.Print("\n	Amount: ")
	samount, err := reader.ReadString('\n')
	samount = strings.Replace(samount, "\n", "", -1)
	samount = strings.Replace(samount, "\r", "", -1)
	Handler(err)
	amount, err := strconv.Atoi(samount)
	
	tx := blockchain.NewTransaction(w, to, uint64(amount), chain)
	network.MEMPOOL[hex.EncodeToString(tx.Id())] = *tx
	
	txm := network.TransactionMessage{}
	txm.Init(tx)
	
	ne := network.NetworkEnvelope{[]byte(network.MyId),txm.GetCommand(),txm.Serialize()}

	network.Message <- ne
}

func DBexists(path string) bool {
	if _, err := os.Stat(path+"/LOCK"); os.IsNotExist(err) {
		return false
	}
	return true
}


func walletExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func enterWallet() *wallet.Wallet{
	var w wallet.Wallet 
	
	err := w.LoadFile("./tmp/wallet.data")
	Handler(err)
	return &w
}

func createWallet() *wallet.Wallet{
	w := wallet.MakeWallet()
	
	w.SaveFile("./tmp/wallet.data")
	
	fmt.Printf("New address is: %s \n",w.Address())
	return w
}


func Handler(err error){
	if err != nil{
		log.Panic(err)
	}
}