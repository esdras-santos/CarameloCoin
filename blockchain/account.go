package blockchain

import (
	"encoding/binary"
	"gochain/wallet"

	"dgraph-io/badger"
)


type Account struct{
	BalanceDatabase *badger.DB
	NonceDatabase *badger.DB
}

var BALANCEPATH = "./tmp/balances"
var NONCEPATH = "./tmp/nonces"

func (acc *Account) BalanceOf(address string) uint64 {
	var balance uint64
	err := acc.BalanceDatabase.View(func(txn *badger.Txn) error {
		item, err := txn.Get(wallet.AddressToPKH(address))
		Handle(err)
		value, err := item.Value()
		Handle(err)
		balance = binary.LittleEndian.Uint64(value)	
			
		return err
	})
	Handle(err)
	return balance
}

func (acc *Account) UpdateBalances(b Block){

	var err error
	var nonce uint64
	var rbalance uint64
	var sbalance uint64
	for _,tx := range b.Transactions{

		err = acc.BalanceDatabase.Update(func(txn *badger.Txn) error {
			if (tx.IsCoinbase()){
				rbalance = acc.BalanceOf(string(tx.receipent))
				err = txn.Set(wallet.AddressToPKH(string(tx.receipent)), ToBytes(rbalance + tx.Value))
				Handle(err)
			} else {
				rbalance = acc.BalanceOf(string(tx.receipent))
				err = txn.Set(wallet.AddressToPKH(string(tx.receipent)), ToBytes(rbalance + tx.Value))
				Handle(err)

				sbalance = acc.BalanceOf(wallet.PKHtoAddress(wallet.PktoPKH(tx.pubkey)))
				err = txn.Set(wallet.PktoPKH(tx.pubkey), ToBytes(sbalance - tx.Value))
				Handle(err)
			}
			return err
		})
		Handle(err)
		
		if (!tx.IsCoinbase()){
			err = acc.NonceDatabase.Update(func(txn *badger.Txn) error {
				nonce = acc.NonceOf(wallet.PKHtoAddress(wallet.PktoPKH(tx.pubkey)))
				
				err = txn.Set(wallet.PktoPKH(tx.pubkey), ToBytes(nonce + 1))
				Handle(err)
				return err
			})
		}
		
		Handle(err)
		
	}
	acc.BalanceDatabase.Close()
	acc.NonceDatabase.Close()
}

func (acc *Account) NonceOf(address string) uint64{
	var nonce uint64
	err := acc.NonceDatabase.View(func(txn *badger.Txn) error {
		item, err := txn.Get(wallet.AddressToPKH(address))
		Handle(err)
		value, err := item.Value()
		Handle(err)
		nonce = binary.LittleEndian.Uint64(value)	
			
		return err
	})
	Handle(err)
	return nonce
}

func GetAccounts() *Account{
	optsb := badger.DefaultOptions
	optsb.Dir = BALANCEPATH
	optsb.ValueDir = BALANCEPATH

	bdb, err := openDB(BALANCEPATH,optsb)
	Handle(err)

	optsn := badger.DefaultOptions
	optsn.Dir = NONCEPATH
	optsn.ValueDir = NONCEPATH

	ndb, err := openDB(NONCEPATH,optsn)
	Handle(err)

	return &Account{bdb, ndb}

}

func InitAccounts(address string) *Account{
	
	optsb := badger.DefaultOptions
	optsb.Dir = BALANCEPATH
	optsb.ValueDir = BALANCEPATH

	bdb, err := openDB(BALANCEPATH,optsb)
	Handle(err)

	err = bdb.Update(func(txn *badger.Txn) error {
		err = txn.Set(wallet.AddressToPKH(address), ToBytes(50))
		Handle(err)
		return err
	})
	Handle(err)

	optsn := badger.DefaultOptions
	optsn.Dir = NONCEPATH
	optsn.ValueDir = NONCEPATH

	ndb, err := openDB(NONCEPATH,optsn)
	Handle(err)

	err = ndb.Update(func(txn *badger.Txn) error {
		err = txn.Set(wallet.AddressToPKH(address), ToBytes(1))
		Handle(err)
		return err
	})
	Handle(err)

	accounts := Account{bdb,ndb}
	return &accounts
}

func ToBytes(n uint64) []byte{
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}

