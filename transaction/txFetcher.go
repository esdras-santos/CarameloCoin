package transaction

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
)

// ALL THIS FILE MUST BE MODIFIED
type TxFetcher struct {
	cache map[*big.Int]*Transaction
}

func (txf TxFetcher) GetUrl(testnet bool) string {
	if testnet {
		return "http://testnet.programmingbitcoin.com"
	} else {
		return "http://mainnet.programmingbitcoin.com"
	}
}

func (txf TxFetcher) Fetch(txid []byte, testnet, fresh bool) *Transaction {
	var tx Transaction
	id := big.Int{}
	_,exists := txf.cache[id.SetBytes(txid)]
	if fresh || !exists {
		url := fmt.Sprintf("%s/tx/%d.hex",txf.GetUrl(false),txid)
		response, err := http.Get(url)
		if err != nil{
			print(err)
		} 
		defer response.Body.Close()
		raw, err := ioutil.ReadAll(response.Body)
		if err != nil{
			print(err)
		}
		
		if raw[4] == 0{
			raw = append(raw[:4], raw[6:]...)
			tx = DeserializeTransaction(raw,false)
			tx.Locktime.SetBytes(toLittleEndian(raw[4:],4))
		}else{
			tx = DeserializeTransaction(raw,false)
		}
		
		if bytes.Equal(tx.Id(),txid){
			log.Panic("not the same ids")
		}
		txf.cache[id.SetBytes(txid)] = &tx
	}
	txf.cache[id.SetBytes(txid)].Testnet = testnet 
	return txf.cache[id.SetBytes(txid)]
}