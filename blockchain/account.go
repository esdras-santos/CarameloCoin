package blockchain

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	
	"gochain/wallet"

	"github.com/syndtr/goleveldb/leveldb"
)




type Account struct{
 	Balance uint64
 	Nonce uint64
}

type AccDB struct{
	AccDatabase *leveldb.DB
}
const ACCPATH = "./tmp/accounts"


func (acc *Account) Parse(s []byte) *Account {
	var acco Account
    by := bytes.Buffer{}
    by.Write(s)
    d := gob.NewDecoder(&by)
    err := d.Decode(&acco)
	Handle(err)

    return &acco
}

func (acc *Account) Serialize() []byte{
	by := bytes.Buffer{}
    e := gob.NewEncoder(&by)
    err := e.Encode(acc)
	Handle(err)
    return by.Bytes()
}

func (acc *AccDB) BalanceNonce(address string) (uint64,uint64) {
	var account Account

	accdata, err := acc.AccDatabase.Get(wallet.AddressToPKH(address),nil)
	Handle(err)
	account = *account.Parse(accdata)	
	
	return account.Balance, account.Nonce
}

func (acc *AccDB) UpdateBalances(b Block){
	var ar Account
	var as Account

	for _,tx := range b.Transactions{
		
		if (tx.IsCoinbase() == true){
			rdata,err := acc.AccDatabase.Get(tx.Recipient,nil)
			if err == nil {
				ar = *ar.Parse(rdata)
			}
			ar.Balance = ar.Balance + tx.Value
			err = acc.AccDatabase.Put(tx.Recipient, ar.Serialize(), nil)
			Handle(err)
		} else {
			if _, err := acc.AccDatabase.Get(tx.Recipient,nil); err != nil{
				a := Account{Balance: tx.Value - 1,Nonce: 0}
	
				err = acc.AccDatabase.Put(tx.Recipient,a.Serialize(),nil)
				Handle(err)
			} else{
				
				rdata,err := acc.AccDatabase.Get(tx.Recipient,nil)
				if err == nil {
					ar = *ar.Parse(rdata)
				}
				
				ar.Balance = ar.Balance + tx.Value - 1
				
				err = acc.AccDatabase.Put(tx.Recipient, ar.Serialize(), nil)
				Handle(err)
			}
			
			

			sdata,err := acc.AccDatabase.Get(wallet.PktoPKH(tx.Pubkey),nil)
			Handle(err)
			as = *as.Parse(sdata)
			as.Balance = as.Balance - tx.Value
			as.Nonce = as.Nonce + 1
			err = acc.AccDatabase.Put(wallet.PktoPKH(tx.Pubkey), as.Serialize(), nil)
			Handle(err)

		}
		
	}
}


func GetAccounts() *AccDB{
	db, err := leveldb.OpenFile(ACCPATH,nil)
	Handle(err)
	return &AccDB{db}

}

func InitAccounts(address string) *AccDB{
	
	db, err := leveldb.OpenFile(ACCPATH,nil)
	Handle(err)

	acc := Account{Balance: 50,Nonce: 0}
	
	err = db.Put(wallet.AddressToPKH(address),acc.Serialize(),nil)
	Handle(err)
	

	accounts := AccDB{db}
	return &accounts
}

func ToBytes(n uint64) []byte{
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}

