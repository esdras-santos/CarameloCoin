package blockchain

// import (
// 	"bytes"
// 	"encoding/hex"
// 	"log"

// 	"github.com/dgraph-io/badger"
// )

// var (
// 	utxoPrefix   = []byte("utxo-")
// 	prefixLength = len(utxoPrefix)
// )

// type UTXOSet struct {
// 	Blockchain *BlockChain
// }

// func (u UTXOSet) Reindex() {
// 	db := u.Blockchain.Database
// 	u.DeleteByPrefix(utxoPrefix)

// 	UTXO := u.Blockchain.FindUTXO()

// 	err := db.Update(func(txn *badger.Txn) error {
// 		for txId, outs := range UTXO {
// 			key, err := hex.DecodeString(txId)
// 			if err != nil {
// 				return err
// 			}
// 			key = append(utxoPrefix, key...)
// 			err = txn.Set(key, SerializeOutputs(outs))
// 			Handle(err)
// 		}
// 		return nil
// 	})
// 	Handle(err)
// }

// func (u *UTXOSet) Update(block *Block) {
// 	db := u.Blockchain.Database

// 	err := db.Update(func(txn *badger.Txn) error {
// 		for _, tx := range block.Transactions {
// 			if tx.IsCoinbase() == false {
// 				for _, in := range tx.Inputs {
// 					updatedOuts := []TxOutput{}
// 					inID := append(utxoPrefix, in.PrevTxID...)
// 					item, err := txn.Get(inID)
// 					Handle(err)
// 					v, err := item.Value()
// 					Handle(err)

// 					outs := ParseOutputs(v)

// 					for outIdx, out := range outs {
// 						if outIdx != int(in.Out) {
// 							updatedOuts = append(updatedOuts, out)
// 						}
// 					}

// 					if len(updatedOuts) == 0 {
// 						if err := txn.Delete(inID); err != nil {
// 							log.Panic(err)
// 						}
// 					} else {
// 						if err := txn.Set(inID, SerializeOutputs(updatedOuts)); err != nil {
// 							log.Panic(err)
// 						}
// 					}
// 				}
// 			}
// 			newOutputs := []TxOutput{}
// 			for _, out := range tx.Outputs {
// 				newOutputs = append(newOutputs, out)
// 			}

// 			txID := append(utxoPrefix, tx.Id()...)
// 			if err := txn.Set(txID, SerializeOutputs(newOutputs)); err != nil {
// 				log.Panic(err)
// 			}
// 		}

// 		return nil
// 	})
// 	Handle(err)
// }

// // func (u UTXOSet) FindUnspentTransactions(pubKeyHash []byte) []TxOutput {
// // 	var UTXOs []TxOutput

// // 	db := u.Blockchain.Database

// // 	err := db.View(func(txn *badger.Txn) error {
// // 		opts := badger.DefaultIteratorOptions

// // 		it := txn.NewIterator(opts)
// // 		defer it.Close()

// // 		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
// // 			item := it.Item()
// // 			v, err := item.Value()
// // 			Handle(err)
// // 			outs := ParseOutputs(v)
// // 			for _, out := range outs {
// // 				if out.IsLockedWithKey(pubKeyHash) {
// // 					UTXOs = append(UTXOs, out)
// // 				}
// // 			}
// // 		}
// // 		return nil
// // 	})
// // 	Handle(err)
// // 	return UTXOs
// // }

// // func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
// // 	unspentOuts := make(map[string][]int)
// // 	db := u.Blockchain.Database
// // 	accumulated := 0

// // 	err := db.View(func(txn *badger.Txn)error{
// // 		opts := badger.DefaultIteratorOptions

// // 		it := txn.NewIterator(opts)
// // 		defer it.Close()

// // 		for it.Seek(utxoPrefix);it.ValidForPrefix(utxoPrefix);it.Next(){
// // 			item := it.Item()
// // 			k := item.Key()
// // 			v, err := item.Value()
// // 			Handle(err)
// // 			k = bytes.TrimPrefix(k,utxoPrefix)
// // 			txID := hex.EncodeToString(k)
// // 			outs := ParseOutputs(v)
			
// // 			for outIdx, out := range outs{
// // 				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount{
// // 					accumulated += int(out.Amount)
// // 					unspentOuts[txID] = append(unspentOuts[txID],outIdx)
// // 				}
// // 			}
// // 		}
// // 		return nil
// // 	})
// // 	Handle(err)
	

// // 	return accumulated, unspentOuts
// // }

// func (u UTXOSet) CountTransactions() int {
// 	db := u.Blockchain.Database
// 	counter := 0

// 	err := db.View(func(txn *badger.Txn) error {
// 		opts := badger.DefaultIteratorOptions
// 		it := txn.NewIterator(opts)
// 		defer it.Close()
// 		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
// 			counter++
// 		}
// 		return nil
// 	})
// 	Handle(err)
// 	return counter
// }

// func (u *UTXOSet) DeleteByPrefix(prefix []byte) {
// 	deleteKeys := func(keysForDelete [][]byte) error {
// 		if err := u.Blockchain.Database.Update(func(txn *badger.Txn) error {
// 			for _, key := range keysForDelete {
// 				if err := txn.Delete(key); err != nil {
// 					return err
// 				}
// 			}
// 			return nil
// 		}); err != nil {
// 			return err
// 		}
// 		return nil
// 	}

// 	collectSize := 100000
// 	u.Blockchain.Database.View(func(txn *badger.Txn) error {
// 		opts := badger.DefaultIteratorOptions
// 		opts.PrefetchValues = false
// 		it := txn.NewIterator(opts)
// 		defer it.Close()

// 		keysForDelete := make([][]byte, 0, collectSize)
// 		keysCollected := 0
// 		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
// 			key := it.Item().KeyCopy(nil)
// 			keysForDelete = append(keysForDelete, key)
// 			keysCollected++
// 			if keysCollected == collectSize {
// 				if err := deleteKeys(keysForDelete); err != nil {
// 					log.Panic(err)
// 				}
// 				keysForDelete = make([][]byte, 0, collectSize)
// 				keysCollected = 0
// 			}
// 		}
// 		if keysCollected > 0 {
// 			if err := deleteKeys(keysForDelete); err != nil {
// 				log.Panic(err)
// 			}
// 		}
// 		return nil
// 	})

// }