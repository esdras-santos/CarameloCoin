package cli

// import (
// 	"bufio"
// 	"flag"
// 	"fmt"
// 	"gochain/blockchain"
// 	"gochain/network"


// 	"gochain/wallet"
// 	"log"
// 	"os"
// 	"runtime"
// 	"strconv"
// )

// type CommandLine struct{}
// func (cli *CommandLine) printUsage(){
// 	fmt.Println("Usage:")
// 	fmt.Println(" getbalance -address ADDRESS - get the balance for the account")
// 	fmt.Println(" createblockchain -address ADDRESS - creates a blockchain")
// 	fmt.Println(" printchain - Prints the blocks in the chain")
// 	fmt.Println(" send -from FROM -to TO -amount AMOUNT -mine - Send amount")
// 	fmt.Println(" createwallet - Creates a new Wallet")
// 	fmt.Println(" listaddresses - Lists the addresses in our wallet file")
// 	fmt.Println(" reindexutxo - Rebuilds the UTXO set")
// 	fmt.Println(" startnode -miner ADDRESS - Start a nodde with I specified in NODE_ID env, var, -miner enables mining")
// }
// func (cli *CommandLine) validateArgs(){
// 	if len(os.Args) < 2 {
// 		cli.printUsage()
// 		runtime.Goexit()
// 	}

// }

// func (cli *CommandLine) StartNode(){
// 	fmt.Printf("Starting Node \n")
// 	w := wallet.MakeWallet()
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter Password: ")
// 	password, err := reader.ReadString('\n')
// 	if err != nil{
// 		log.Panic(err)
// 	}
// 	err = w.LoadFile(password,"./tmp/wallet.data")
// 	if err != nil{
// 		log.Panic(err)
// 	}
// 	// if len(network.NODEIP) > 0{
// 	//  	if wallet.ValidateAddress(string(w.Address())){
// 	//  		fmt.Println("Mining is on. Address to receive rewards: ",string(w.Address()))
// 	//  	}else{
// 	//  		log.Panic("Wrong miner address")
// 	//  	}
// 	// }
// 	//this have to wait until you manually finish the connection
// 	network.Connect()
// 	fmt.Printf("Running at %s...",w.Address())
// 	for{

// 	}

// }

// func Handle(err error){
// 	if err != nil {
// 		panic(err)
// 	}
// }

// // func (cli *CommandLine) reindexUTXO(){
// // 	chain := blockchain.ContinueBlockChain("./tmp/blocks")
// // 	defer chain.Database.Close()
// // 	UTXOSet := blockchain.UTXOSet{chain}
// // 	UTXOSet.Reindex()

// // 	count := UTXOSet.CountTransactions()
// // 	fmt.Printf("Done! There are %d transactions in the UTXO set.\n",count)
// // }

// // func (cli *CommandLine) listAddresses(nodeID string){
// // 	wallets,_ := wallet.CreateWallets(nodeID)
// // 	addresses := wallets.GetAllAddresses()

// // 	for _,address := range addresses{
// // 		fmt.Println(address)
// // 	}
// // }

// func walletExists(path string) bool {
// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		return false
// 	}
// 	return true
// }

// func (cli *CommandLine) createWallet(){
// 	if walletExists("./tmp/wallet.data"){
// 		fmt.Println("Wallet already exists")
// 		runtime.Goexit()
// 	}
// 	w := wallet.MakeWallet()
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter Password: ")
// 	password, err := reader.ReadString('\n')
// 	if err != nil{
// 		log.Panic(err)
// 	}
// 	reader1 := bufio.NewReader(os.Stdin)
// 	fmt.Print("confirm Password: ")
// 	password1, err := reader1.ReadString('\n')
// 	if err != nil{
// 		log.Panic(err)
// 	}
// 	if password != password1{
// 		log.Panic("wrong passoword!")
// 	}
// 	w.SaveFile(password,"./tmp/wallet.data")
	
// 	fmt.Printf("New address is: %s \n",w.Address())
// }

// func (cli *CommandLine) printChain(){
// 	chain := blockchain.ContinueBlockChain("./tmp/blocks")
// 	defer chain.Database.Close()
// 	iter := chain.Iterator()

// 	for {
// 		block := iter.Next()
// 		fmt.Printf("Prev. Hash: %x\n",block.PrevBlock)
// 		fmt.Printf("Hash: %x\n",block.Hash())

// 		pow := blockchain.NewProof(block)
// 		fmt.Printf("PoW:%s \n",strconv.FormatBool(pow.Validate()))
// 		for _,tx := range block.Transactions{
// 			fmt.Println(tx)
// 		}
// 		fmt.Println()

// 		if len(block.PrevBlock) == 0{
// 			break
// 		}
// 	}
// }
// func (cli *CommandLine) createblockchain(){
// 	w := wallet.MakeWallet()
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter Password: ")
// 	password, err := reader.ReadString('\n')
// 	if err != nil{
// 		log.Panic(err)
// 	}
// 	err = w.LoadFile(password,"./tmp/wallet.data")
// 	if err != nil{
// 		log.Panic(err)
// 	}
	
	


// 	chain := blockchain.InitBlockChain(w,"./tmp/blocks")

// 	bm := network.BlockMessage{}
// 	block,err := chain.GetBlock(chain.LastHash)
// 	Handle(err)
	
		
// 	bm.Init(block)
// 	//ne := network.NetworkEnvelope{[]byte(network.MyId),bm.GetCommand(),bm.Serialize()}
// 	print("\n enveloped\n")
	
// 	print("\n published \n")
	

// 	defer chain.Database.Close()


// 	// UTXOSet := blockchain.UTXOSet{chain}
// 	// UTXOSet.Reindex()
// 	fmt.Println("Finished!")
// }
// func (cli *CommandLine) getBalance(){
// 	chain := blockchain.ContinueBlockChain("./tmp/blocks")
// 	defer chain.Database.Close()
// 	w := wallet.MakeWallet()
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter Password: ")
// 	password, err := reader.ReadString('\n')
// 	if err != nil{
// 		log.Panic(err)
// 	}
// 	err = w.LoadFile(password,"./tmp/wallet.data")
// 	if err != nil{
// 		log.Panic(err)
// 	}

// 	balance, _ := chain.Acc.BalanceNonce(string(w.Address()))
// 	fmt.Printf("Balance of %s: %d \n",w.Address() ,balance)
// }

// func (cli *CommandLine) Mine(){
	
// 	chain := blockchain.ContinueBlockChain("./tmp/blocks")
// 	defer chain.Database.Close()
// 	wFrom := wallet.MakeWallet()
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter Password: ")
// 	password, err := reader.ReadString('\n')
// 	if err != nil{
// 		log.Panic(err)
// 	}
// 	err = wFrom.LoadFile(password,"./tmp/wallet.data")
// 	if err != nil{
// 		log.Panic(err)
// 	}


// 	cbTx := blockchain.CoinbaseTx(wFrom)
// 	for range network.MEMPOOL{
// 			//collecting fees
// 		cbTx.Value += 1
// 	}
		
// 	txs := []*blockchain.Transaction{cbTx}

// 	for _,t := range network.MEMPOOL{
// 		txs = append(txs, &t)

// 	}
	
// 	network.Connect()
// 	bm := network.BlockMessage{}

// 	block := chain.MineBlock(txs)
	
// 	fmt.Println(block.ToString())
// 	fmt.Printf("\nh: %x\n",block.Hash())
// 	bm.Init(*block)
// 	ne := network.NetworkEnvelope{[]byte(network.MyId),bm.GetCommand(),bm.Serialize()}
// 	network.Message <- ne
// 	fmt.Println("Success!")
// 	select{}
// }

// func (cli *CommandLine) send(to string, amount uint64, mineNow bool){
// 	if !wallet.ValidateAddress(to){
// 		log.Panic("Address is not Valid")
// 	}
	
// 	chain := blockchain.ContinueBlockChain("./tmp/blocks")
// 	defer chain.Database.Close()

	
// 	wFrom := wallet.MakeWallet()
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter Password: ")
// 	password, err := reader.ReadString('\n')
// 	if err != nil{
// 		log.Panic(err)
// 	}
// 	err = wFrom.LoadFile(password,"./tmp/wallet.data")
// 	if err != nil{
// 		log.Panic(err)
// 	}


// 	tx := blockchain.NewTransaction(wFrom, to, amount, chain)
// 	network.Connect()
// 	txm := network.TransactionMessage{}
// 	txm.Init(tx)
// 	ne := network.NetworkEnvelope{[]byte(network.MyId),txm.GetCommand(),txm.Serialize()}
// 	network.Message <- ne

// 	if mineNow{
// 		cbTx := blockchain.CoinbaseTx(wFrom)
// 		for range network.MEMPOOL{
// 			//collecting fees
// 			cbTx.Value += 1
// 		}
		
// 		txs := []*blockchain.Transaction{cbTx, tx}

// 		for _,t := range network.MEMPOOL{
// 			txs = append(txs, &t)

// 		}
// 		bm := network.BlockMessage{}

// 		block := chain.MineBlock(txs)
// 		Handle(err)
// 		bm.Init(*block)
// 		ne := network.NetworkEnvelope{[]byte(network.MyId),bm.GetCommand(),bm.Serialize()}
// 		network.Message <- ne
// 		//send message to all nodes for exclude the tx from your mempool because is already mined
		
// 		fmt.Println("Success!")
// 	}

	
// }

// func (cli *CommandLine) Run(){
// 	cli.validateArgs()

	

// 	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
// 	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
// 	sendCmd := flag.NewFlagSet("send",flag.ExitOnError)
// 	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
// 	createWalletCmd := flag.NewFlagSet("createwallet",flag.ExitOnError)
// 	listAddressesCmd := flag.NewFlagSet("listaddresses",flag.ExitOnError)
// 	reindexUTXOCmd := flag.NewFlagSet("reindexutxo",flag.ExitOnError)
// 	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
// 	mineCmd := flag.NewFlagSet("mine", flag.ExitOnError)

// 	//getBalanceAddress := getBalanceCmd.String("address","","The address to get balance for")
// 	//createBlockchainAddress := createBlockchainCmd.String("address","","The address to send genesis block reward to")
// 	sendFrom := sendCmd.String("from","","Source wallet address")
// 	sendTo := sendCmd.String("to","","Destination wallet address")
// 	sendAmount := sendCmd.Uint64("amount",0,"Amount to send")
// 	sendMine := sendCmd.Bool("mine",false,"Mine immediately on the same node")
// 	//startNodeMiner := startNodeCmd.String("miner", "","Enable mining node and send reward to")

// 	switch os.Args[1]{
// 	case "mine":
// 		err := mineCmd.Parse(os.Args[2:])
// 		blockchain.Handle(err)
// 	case "reindexutxo":
// 		err := reindexUTXOCmd.Parse(os.Args[2:])
// 		blockchain.Handle(err)
// 	case "getbalance":
// 		err := getBalanceCmd.Parse(os.Args[2:])
// 		blockchain.Handle(err)
// 	case "startnode":
// 		err := startNodeCmd.Parse(os.Args[2:])
// 		blockchain.Handle(err)
// 	case "createblockchain":
// 		err := createBlockchainCmd.Parse(os.Args[2:])
// 		blockchain.Handle(err)
// 	case "createwallet":
// 		err := createWalletCmd.Parse(os.Args[2:])
// 		blockchain.Handle(err)
// 	case "listaddresses":
// 		err := listAddressesCmd.Parse(os.Args[2:])
// 		blockchain.Handle(err)
// 	case "printchain":
// 		err := printChainCmd.Parse(os.Args[2:])
// 		blockchain.Handle(err)
// 	case "send":
// 		err := sendCmd.Parse(os.Args[2:])
// 		blockchain.Handle(err)
// 	default:
// 		cli.printUsage()
// 		runtime.Goexit()
// 	}

// 	if getBalanceCmd.Parsed(){
		
// 		cli.getBalance()
// 	}
// 	if createBlockchainCmd.Parsed(){
		
// 		cli.createblockchain()
// 	}
// 	if printChainCmd.Parsed(){
// 		cli.printChain()
// 	}
// 	if createWalletCmd.Parsed(){
// 		cli.createWallet()
// 	}
// 	// if listAddressesCmd.Parsed(){
// 	// 	cli.listAddresses()
// 	// }
// 	// if reindexUTXOCmd.Parsed(){
// 	// 	cli.reindexUTXO()
// 	// }
// 	if sendCmd.Parsed(){
// 		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0{
// 			sendCmd.Usage()
// 			runtime.Goexit()
// 		}
// 		cli.send(*sendTo,*sendAmount,*sendMine)
// 	}

// 	if startNodeCmd.Parsed(){
		
// 		cli.StartNode()
// 	}

// 	if mineCmd.Parsed(){
// 		cli.Mine()
// 	}
// }