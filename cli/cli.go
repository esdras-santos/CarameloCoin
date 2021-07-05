package cli

import (
	"bufio"
	"flag"
	"fmt"
	"gochain/blockchain"
	"gochain/network"

	"gochain/wallet"
	"log"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct{}
func (cli *CommandLine) printUsage(){
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get the balance for the account")
	fmt.Println(" createblockchain -address ADDRESS - creates a blockchain")
	fmt.Println(" printchain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT -mine - Send amount")
	fmt.Println(" createwallet - Creates a new Wallet")
	fmt.Println(" listaddresses - Lists the addresses in our wallet file")
	fmt.Println(" reindexutxo - Rebuilds the UTXO set")
	fmt.Println(" startnode -miner ADDRESS - Start a nodde with I specified in NODE_ID env, var, -miner enables mining")
}
func (cli *CommandLine) validateArgs(){
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}

}

func (cli *CommandLine) StartNode(){
	fmt.Printf("Starting Node \n")

	if len(network.NODEIP) > 0{
		if wallet.ValidateAddress(network.NODEIP){
			fmt.Println("Mining is on. Address to receive rewards: ",network.NODEIP)
		}else{
			log.Panic("Wrong miner address")
		}
	}
	network.StartServer()
}

func (cli *CommandLine) reindexUTXO(){
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n",count)
}

// func (cli *CommandLine) listAddresses(nodeID string){
// 	wallets,_ := wallet.CreateWallets(nodeID)
// 	addresses := wallets.GetAllAddresses()

// 	for _,address := range addresses{
// 		fmt.Println(address)
// 	}
// }

func (cli *CommandLine) createWallet(){
	w := wallet.MakeWallet()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Password: ")
	password, err := reader.ReadString('\n')
	if err != nil{
		log.Panic(err)
	}
	reader1 := bufio.NewReader(os.Stdin)
	fmt.Print("confirm Password: ")
	password1, err := reader1.ReadString('\n')
	if err != nil{
		log.Panic(err)
	}
	if password != password1{
		log.Panic("wrong passoword!")
	}
	w.SaveFile(password)
	
	fmt.Printf("New address is: %s \n",w.Address())
}

func (cli *CommandLine) printChain(){
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Prev. Hash: %x\n",block.BH.PrevBlock)
		fmt.Printf("Hash: %x\n",block.BH.Hash())

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW:%s \n",strconv.FormatBool(pow.Validate()))
		for _,tx := range block.Transactions{
			fmt.Println(tx)
		}
		fmt.Println()

		if len(block.BH.PrevBlock) == 0{
			break
		}
	}
}
func (cli *CommandLine) createblockchain(){
	w := wallet.MakeWallet()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Password: ")
	password, err := reader.ReadString('\n')
	if err != nil{
		log.Panic(err)
	}
	err = w.LoadFile(password)
	if err != nil{
		log.Panic(err)
	}
	chain := blockchain.InitBlockChain(w)
	chain.Database.Close()

	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()
	fmt.Println("Finished!")
}
func (cli *CommandLine) getBalance(address string){
	if !wallet.ValidateAddress(address){
		log.Panic("Address is not Valid")
	}
	chain := blockchain.ContinueBlockChain()
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Database.Close()

	balance := 0
	pubKeyHash := wallet.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUnspentTransactions(pubKeyHash)

	for _,out := range UTXOs{
		balance += int(out.Amount)
	}
	fmt.Printf("Balance of %s: %d \n",address,balance)
}
func (cli *CommandLine) send(to string, amount int, mineNow bool){
	if !wallet.ValidateAddress(to){
		log.Panic("Address is not Valid")
	}
	
	chain := blockchain.ContinueBlockChain()
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Database.Close()

	
	wFrom := wallet.MakeWallet()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Password: ")
	password, err := reader.ReadString('\n')
	if err != nil{
		log.Panic(err)
	}
	err = wFrom.LoadFile(password)
	if err != nil{
		log.Panic(err)
	}


	tx := blockchain.NewTransaction(wFrom, to, amount, &UTXOSet)
	if mineNow{
		cbTx := blockchain.CoinbaseTx(wFrom)
		
		txs := []*blockchain.Transaction{cbTx, tx}

		for _,t := range network.MEMPOOL{
			txs = append(txs, &t)

		}
		block := chain.MineBlock(txs)
		UTXOSet.Update(block)
		fmt.Println("Success!")
	}else{
		var nc network.NodeCommand
		//send the transaction to the mempool of all your known nodes
		go func(){
			for i,_ := range network.KNOWNNODES{
				nc.Init(network.KNOWNNODES[i])
				nc.SendTransaction(*tx)
			}
			fmt.Println("Success!")
		}()
	}

	
}

func (cli *CommandLine) Run(){
	cli.validateArgs()

	

	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet",flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses",flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo",flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address","","The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address","","The address to send genesis block reward to")
	sendFrom := sendCmd.String("from","","Source wallet address")
	sendTo := sendCmd.String("to","","Destination wallet address")
	sendAmount := sendCmd.Int("amount",0,"Amount to send")
	sendMine := sendCmd.Bool("mine",false,"Mine immediately on the same node")
	//startNodeMiner := startNodeCmd.String("miner", "","Enable mining node and send reward to")

	switch os.Args[1]{
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed(){
		if *getBalanceAddress == ""{
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}
	if createBlockchainCmd.Parsed(){
		if *createBlockchainAddress == ""{
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createblockchain()
	}
	if printChainCmd.Parsed(){
		cli.printChain()
	}
	if createWalletCmd.Parsed(){
		cli.createWallet()
	}
	// if listAddressesCmd.Parsed(){
	// 	cli.listAddresses()
	// }
	if reindexUTXOCmd.Parsed(){
		cli.reindexUTXO()
	}
	if sendCmd.Parsed(){
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0{
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendTo,*sendAmount,*sendMine)
	}

	if startNodeCmd.Parsed(){
		
		cli.StartNode()
	}
}