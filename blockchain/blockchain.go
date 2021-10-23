package blockchain

import (
	"bytes"

	"encoding/binary"
	"errors"
	"fmt"
	"gochain/wallet"
	"log"
	"os"
	
	"runtime"
	
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	
)

const (
	genesisData = "First Transaction from Genesis"
)

var (
	once sync.Once
	BlockchainInstance BlockChain
)

type BlockChain struct {
	LastHash []byte
	Database *leveldb.DB
	Acc *AccDB
	
}

func GetBlockChainInstance(lastHash []byte, db *leveldb.DB, acc *AccDB) BlockChain{
	once.Do(func(){
			BlockchainInstance = BlockChain{lastHash,db,acc}
	})
	return BlockchainInstance
}


func (chain *BlockChain) AddBlock(block *Block){
	var lastBlock Block
	
	//return if block already exists
	if _, err := chain.Database.Get(block.Hash(),nil); err == nil{
		return 
	}
	

	blockData := block.Serialize()
	err := chain.Database.Put(block.Hash(), blockData,nil)
	Handle(err)
	lastHash, err := chain.Database.Get([]byte("lh"),nil)
	Handle(err)
	data, err := chain.Database.Get(lastHash,nil)
	Handle(err)
	lastBlock = *lastBlock.Parse(data)
	
	if binary.LittleEndian.Uint64(block.Height) > binary.LittleEndian.Uint64(lastBlock.Height){
		
		err = chain.Database.Put([]byte("lh"), block.Hash(),nil)
		Handle(err)
		chain.LastHash = block.Hash()
	}
	
	chain.Acc.UpdateBalances(*block)
	

}



func (chain *BlockChain) GetBlock(blockHash []byte) (Block, error){
	var block *Block
	
	blockData, err := chain.Database.Get(blockHash,nil);
	if  err != nil{
		log.Panic(err)
	}else{
		block = block.Parse(blockData)
	}
	
	return *block, err
}


//return the block headers 
// func (chain *BlockChain) GetBlockHeaders(startBlock, endBlock []byte) []BlockHeader{
// 	var blockHeaders []BlockHeader
// 	if bytes.Equal(startBlock,endBlock){
// 		b,err := chain.GetBlock(endBlock)
// 		Handle(err)
// 		blockHeaders = append(blockHeaders, b.BH)
// 		return blockHeaders
// 	}
// 	iter := &BlockChainIterator{endBlock,chain.Database}

// 	for{
// 		block := *iter.Next()
// 		if bytes.Equal(block.BH.Hash(),startBlock){
// 			blockHeaders = append(blockHeaders, block.BH)
// 			break
// 		}
// 		blockHeaders = append(blockHeaders, block.BH)
// 	}
// 	return blockHeaders
// }

func (chain *BlockChain) GetBlockHashes() [][]byte{
	var blocks [][]byte

	iter := chain.Iterator()

	for {
		block := iter.Next()

		blocks = append(blocks, block.Hash())

		if len(block.PrevBlock) == 0{
			break
		}
	}

	return blocks
}

func (chain *BlockChain) GetBestHeight() uint64{
	var lastBlock Block
	
	lastHash, err := chain.Database.Get([]byte("lh"),nil);
	Handle(err)
	blockData, err := chain.Database.Get(lastHash,nil);
	Handle(err)
	lastBlock = *lastBlock.Parse(blockData)
	
	return binary.BigEndian.Uint64(lastBlock.Height)
}

func (chain *BlockChain) GetLastHash() []byte{
	var lastHash []byte
	var lastBlock Block
	
	lastHash, err := chain.Database.Get([]byte("lh"),nil);
	Handle(err)
	blockData, err := chain.Database.Get(lastHash,nil);
	Handle(err)
	lastBlock = *lastBlock.Parse(blockData)

	
	return lastBlock.Hash()
}

func (chain *BlockChain) MineBlock(transactions []*Transaction) *Block{
	var lastHash []byte
	var lastHeight uint64
	var lastBlock Block

	for _,tx := range transactions{
		if chain.VerifyTransaction(tx) != true{
			log.Panic("Invalid Transaction")
		}
	}
	
	lastHash, err := chain.Database.Get([]byte("lh"),nil);
	Handle(err)
	blockData, err := chain.Database.Get(lastHash,nil);
	Handle(err)
	lastBlock = *lastBlock.Parse(blockData)
	lastHeight = binary.LittleEndian.Uint64(lastBlock.Height)
	
	newBlock := CreateBlock(transactions, lastHash, lastHeight+1)
	
	err = chain.Database.Put(newBlock.Hash(),newBlock.Serialize(),nil)
	Handle(err)
	err = chain.Database.Put([]byte("lh"),newBlock.Hash(),nil)
	chain.LastHash = newBlock.Hash()


	chain.Acc.UpdateBalances(*newBlock)

	return newBlock	
}

func DBexists(path string) bool {
	if _, err := os.Stat(path+"/LOCK"); os.IsNotExist(err) {
		return false
	}
	return true
}

func InitBlockChain(w *wallet.Wallet,dbPath string) *BlockChain {
	

	if DBexists(dbPath) {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	var lastHash []byte
	
	db, err := leveldb.OpenFile(dbPath,nil)
	Handle(err)
	//defer db.Close()
	cbtx := CoinbaseTx(w)
	genesis := Genesis(cbtx)
	fmt.Println("Genesis created")

	err = db.Put(genesis.Hash(), genesis.Serialize(),nil)
	
	Handle(err)
	err = db.Put([]byte("lh"),genesis.Hash(),nil)
	Handle(err)
	lastHash = genesis.Hash()
	
	acc := InitAccounts(string(w.Address()))
	
	blockchain := BlockChain{lastHash, db, acc}
	
	return &blockchain
}

func ContinueBlockChain(dbPath string) *BlockChain {
	if DBexists(dbPath) == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := leveldb.OpenFile(dbPath,nil)
	Handle(err)
	lastHash, err = db.Get([]byte("lh"),nil)
	Handle(err)
	acc := GetAccounts()
	chain := GetBlockChainInstance(lastHash, db, acc)
	return &chain
}



// func (chain *BlockChain) FindUTXO() map[string][]TxOutput{
// 	UTXO := make(map[string][]TxOutput)
// 	spentTXOs := make(map[string][]int)

// 	iter := chain.Iterator()

// 	for {
// 		block := iter.Next()

// 		for _,tx := range block.Transactions{
// 			txID := hex.EncodeToString(tx.Id())

// 		Outputs:
// 			for outIdx,out := range tx.Outputs{
// 				if spentTXOs[txID] != nil{
// 					for _,spentOut := range spentTXOs[txID]{
// 						if spentOut == outIdx{
// 							continue Outputs
// 						}
// 					}
// 				}
// 				outs := UTXO[txID]
// 				outs = append(outs,out)
// 				UTXO[txID] = outs
// 			}
// 			if tx.IsCoinbase() == false{
// 				for _,in := range tx.Inputs{
// 					inTxID := hex.EncodeToString(in.PrevTxID)
// 					spentTXOs[inTxID] = append(spentTXOs[inTxID],int(in.Out))
// 				}
// 			}
// 		}
// 		if len(block.BH.PrevBlock) == 0{
// 			break
// 		}
// 	}
	
// 	return UTXO
// }

func (bc *BlockChain) FindTransaction(ID []byte) (*Transaction,error){
	iter := bc.Iterator()	

	for {
		block := iter.Next()
		for _,tx := range block.Transactions{
			if bytes.Compare(tx.Id(),ID) == 0{
				return tx,nil
			}
		}

		if len(block.PrevBlock) == 0{
			break
		}
	}
	return &Transaction{},errors.New("Transaction does not exist")
}



func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool{
	if tx.IsCoinbase(){

		return true
	}

	balance,_ := bc.Acc.BalanceNonce(wallet.PKHtoAddress(wallet.PktoPKH(tx.Pubkey)))
	if !wallet.VerifySignature(tx.Id(), tx.Pubkey, tx.Sig) {
		
		return false
	} else if balance < tx.Value{
		
		return false
	}

	return true
}

// func retry(dir string, originalOpts badger.Options) (*badger.DB,error){
// 	lockPath := filepath.Join(dir, "LOCK")
// 	if err := os.Remove(lockPath); err != nil{
// 		return nil, fmt.Errorf(`removing "LOCK": %s`,err)
// 	}
// 	retryOpts := originalOpts
// 	retryOpts.Truncate = true
// 	db,err := badger.Open(retryOpts)
// 	return db, err
// }

// func openDB(dir string, opts badger.Options) (*badger.DB,error){
// 	if db, err := badger.Open(opts); err != nil{
// 		if strings.Contains(err.Error(), "LOCK"){
// 			if db, err := retry(dir, opts); err == nil{
// 				log.Println("database unlocked, value log truncated")
// 				return db, nil
// 			}
// 			log.Println("could not unlock database:", err)
// 		}
// 		return nil, err
// 	}else{
// 		return db,nil
// 	}
// }