package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	
	"fmt"
	"log"
	
	"strings"

	
	
	"gochain/wallet"
)

type Transaction struct {
	sig []byte
	nonce uint64
	pubkey []byte
	receipent []byte
	Value uint64
}

func (tx *Transaction) hash() []byte {
	var hash [32]byte

	txCopy := tx
	//txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func (tx *Transaction) Id() []byte{
	return tx.hash()
}

func (tx Transaction) Serialize() []byte {

	b := bytes.Buffer{}
    e := gob.NewEncoder(&b)
    err := e.Encode(tx)
	Handle(err)
    return b.Bytes()

}

func (tx *Transaction) Parse(data []byte) *Transaction {
	var txn Transaction
    b := bytes.Buffer{}
    b.Write(data)
    d := gob.NewDecoder(&b)
    err := d.Decode(&txn)
	Handle(err)
    return &txn
}

// func (tx Transaction) Fee() uint{
// 	var inputSum, outputSum uint

// 	for _,txin := range tx.Inputs{
// 		inputSum += txin.Value()
// 	}
// 	for _,txout := range tx.Outputs{
// 		outputSum += txout.Amount
// 	}
// 	fee := inputSum - outputSum
// 	if outputSum + fee != inputSum{
// 		log.Panic("fee overflow")
// 	}
// 	return fee
// }

func CoinbaseTx(w *wallet.Wallet) *Transaction{
	var tx Transaction
	tx.nonce = 0
	tx.pubkey = []byte{0x0000000000000000000000000000000000000001}
	tx.sig = []byte{0x00000001}
	tx.receipent = w.Address()
	//miner prize
	tx.Value = 50

	return &tx
}

func NewTransaction(from *wallet.Wallet, to string, amount uint64, chain *BlockChain) *Transaction {
	if(chain.Acc.BalanceOf(string(from.Address())) < amount){
		log.Panic("not enough balance!")
	}
	tx := Transaction{}
	tx.nonce = chain.Acc.NonceOf(string(from.Address())) + 1
	tx.pubkey = from.PublicKey
	tx.Value = uint64(amount) + 1//fee
	tx.receipent = []byte(to)

	tx.Sign(from,chain)

	return &tx
}

func (tx *Transaction) IsCoinbase() bool { 
	if tx.nonce == 0{
		if tx.pubkey == nil{
			if bytes.Equal(tx.sig, []byte{0x00000001}){
				return true
			}
		}
	}
	return false
}

func (tx *Transaction) Sign(wallet *wallet.Wallet, chain *BlockChain) {
	if tx.IsCoinbase() {
		return
	}

	if (chain.Acc.BalanceOf(string(wallet.Address())) >= uint64(tx.Value)){
		r, s, err := ecdsa.Sign(rand.Reader, &wallet.PrivateKey, tx.Id())
		signature := append(r.Bytes(), s.Bytes()...)
		tx.sig = signature
		Handle(err)
	} else{
		log.Panic("not enough founds!")
	}

}


// func (tx *Transaction) VerifyTransaction(UTXOs map[string]Transaction) bool {
// 	if tx.IsCoinbase() {
// 		return true
// 	}


// 	for inId, in := range tx.Inputs {
	
// 		scriptsig := tx.Inputs[inId].ScriptSig
// 		scriptpubKey := UTXOs[hex.EncodeToString(in.PrevTxID)].Outputs[in.Out].ScriptPubKey

// 		scriptin := script.Script{}
// 		scriptout := script.Script{}
// 		scriptin.Cmd,_ = scriptin.ScriptParser(scriptsig)
// 		scriptout.Cmd,_ = scriptin.ScriptParser(scriptpubKey)
// 		script := scriptin.Add(scriptout)
// 		dataToVerify := fmt.Sprintf("%x\n", tx)
// 		if !script.Evaluate([]byte(dataToVerify)){
// 			return false
// 		}
// 	}

// 	return true
// }



// func (tx *Transaction) TrimmedCopy() Transaction {
// 	var inputs []TxInput
// 	var outputs []TxOutput

// 	for _, in := range tx.Inputs {
// 		inputs = append(inputs, TxInput{in.PrevTxID, in.Out, nil, nil})
// 	}

// 	for _, out := range tx.Outputs {
// 		outputs = append(outputs, TxOutput{out.Amount, out.ScriptPubKey})
// 	}

// 	txCopy := Transaction{[]byte{0x00000001},utils.ToHex(4294967295), inputs, outputs}

// 	return txCopy
// }



func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.Id()))
	lines = append(lines, fmt.Sprintf("       Nonce:     %d", tx.nonce))
	lines = append(lines, fmt.Sprintf("       Public Key:       %x", tx.pubkey))
	lines = append(lines, fmt.Sprintf("       Signature: %x", tx.sig))
	lines = append(lines, fmt.Sprintf("       Receipent:    %s", tx.receipent))
	lines = append(lines, fmt.Sprintf("       Value:    %d", tx.Value))

	return strings.Join(lines, "\n")
}
