package blockchain

import (
	"bytes"

	"encoding/binary"
	"errors"
	"fmt"
	"gochain/wallet"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/dgraph-io/badger"
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
	Database *badger.DB
	Acc *Account
	
}

func GetBlockChainInstance(lastHash []byte, db *badger.DB, acc *Account) BlockChain{
	once.Do(func(){
			BlockchainInstance = BlockChain{lastHash,db,acc}
	})
	return BlockchainInstance
}


func (chain *BlockChain) AddBlock(block *Block){
	var lastBlock Block
	err := chain.Database.Update(func(txn *badger.Txn) error{
		if _,err := txn.Get(block.BH.Hash()); err == nil{
			return nil
		}

		blockData := block.Serialize()
		err := txn.Set(block.BH.Hash(), blockData)
		Handle(err)

		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, _ := item.Value()

		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData,_ := item.Value()

		lastBlock.Parse(lastBlockData)

		if binary.BigEndian.Uint64(block.Height) > binary.BigEndian.Uint64(lastBlock.Height){
			err = txn.Set([]byte("lh"), block.BH.Hash())
			Handle(err)
			chain.LastHash = block.BH.Hash()
		}
		chain.Acc.UpdateBalances(*block)
		return nil
	})
	Handle(err)
}



func (chain *BlockChain) GetBlock(blockHash []byte) (Block, error){
	var block *Block
	
	err := chain.Database.View(func(txn *badger.Txn) error{
		item, err := txn.Get(blockHash);
		
		if  err != nil{
			return errors.New("Block is not found")
		}else{
			blockData,err := item.Value()
			Handle(err)
			block = block.Parse(blockData)
		}
	
		return nil
	})

	if err != nil{
		return *block,err
	}
	print("returned")
	return *block, nil
}


//return the block headers 
func (chain *BlockChain) GetBlockHeaders(startBlock, endBlock []byte) []BlockHeader{
	var blockHeaders []BlockHeader
	if bytes.Equal(startBlock,endBlock){
		b,err := chain.GetBlock(endBlock)
		Handle(err)
		blockHeaders = append(blockHeaders, b.BH)
		return blockHeaders
	}
	iter := &BlockChainIterator{endBlock,chain.Database}

	for{
		block := *iter.Next()
		if bytes.Equal(block.BH.Hash(),startBlock){
			blockHeaders = append(blockHeaders, block.BH)
			break
		}
		blockHeaders = append(blockHeaders, block.BH)
	}
	return blockHeaders
}

func (chain *BlockChain) GetBlockHashes() [][]byte{
	var blocks [][]byte

	iter := chain.Iterator()

	for {
		block := iter.Next()

		blocks = append(blocks, block.BH.Hash())

		if len(block.BH.PrevBlock) == 0{
			break
		}
	}

	return blocks
}

func (chain *BlockChain) GetBestHeight() uint64{
	var lastBlock Block

	err := chain.Database.View(func(txn *badger.Txn) error{
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash,_ := item.Value()

		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData,_ := item.Value()

		lastBlock.Parse(lastBlockData)

		return nil
	})
	Handle(err)
	return binary.BigEndian.Uint64(lastBlock.Height)
}

func (chain *BlockChain) GetLastHash() []byte{
	var lastHash []byte
	var lastBlock *Block
	
	err := chain.Database.View(func(txn *badger.Txn) error{
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		lastHash,err = item.Value()
		Handle(err)

		
		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData,_ := item.Value()

		lastBlock = lastBlock.Parse(lastBlockData)

		return nil
	})
	Handle(err)
	
	return lastBlock.BH.Hash()
}

func (chain *BlockChain) MineBlock(transactions []*Transaction) *Block{
	var lastHash []byte
	var lastHeight uint64

	for _,tx := range transactions{
		if chain.VerifyTransaction(tx) != true{
			log.Panic("Invalid Transaction")
		}
	}
	
	err := chain.Database.View(func(txn *badger.Txn) error {
		var lastBlock *Block
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()
		Handle(err)

		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData, err := item.Value()
		Handle(err)
		
		lastBlock = lastBlock.Parse(lastBlockData)
		
		lastHeight = binary.LittleEndian.Uint64(lastBlock.Height)
		
		return err
	})
	Handle(err)
	
	newBlock := CreateBlock(transactions, lastHash, lastHeight+1)
	

	err = chain.Database.Update(func(txn *badger.Txn) error {
		//the blockheader hash will be linked to the block
		err := txn.Set(newBlock.BH.Hash(), newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.BH.Hash())
		chain.LastHash = newBlock.BH.Hash()

		return err
	})
	Handle(err)

	chain.Acc.UpdateBalances(*newBlock)

	return newBlock	
}

func DBexists(path string) bool {
	if _, err := os.Stat(path+"/MANIFEST"); os.IsNotExist(err) {
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
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := openDB(dbPath,opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(w)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		err = txn.Set(genesis.BH.Hash(), genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.BH.Hash())

		lastHash = genesis.BH.Hash()

		return err
	})
	Handle(err)
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

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := openDB(dbPath,opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()
		return err
	})
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

		if len(block.BH.PrevBlock) == 0{
			break
		}
	}
	return &Transaction{},errors.New("Transaction does not exist")
}



func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool{
	if tx.IsCoinbase(){

		return true
	}

	if !wallet.VerifySignature(tx.Id(), tx.pubkey, tx.sig) {
		
		return false
	} else if bc.Acc.BalanceOf(wallet.PKHtoAddress(tx.pubkey)) < tx.Value{
		
		return false
	}

	return true
}

func retry(dir string, originalOpts badger.Options) (*badger.DB,error){
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil{
		return nil, fmt.Errorf(`removing "LOCK": %s`,err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	db,err := badger.Open(retryOpts)
	return db, err
}

func openDB(dir string, opts badger.Options) (*badger.DB,error){
	if db, err := badger.Open(opts); err != nil{
		if strings.Contains(err.Error(), "LOCK"){
			if db, err := retry(dir, opts); err == nil{
				log.Println("database unlocked, value log truncated")
				return db, nil
			}
			log.Println("could not unlock database:", err)
		}
		return nil, err
	}else{
		return db,nil
	}
}