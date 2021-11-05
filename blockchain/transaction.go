package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"math/big"

	"fmt"
	"log"

	"strings"

	"gochain/wallet"
)

type Transaction struct {
	Sig []byte
	Nonce uint64
	Pubkey []byte
	Recipient []byte
	Value uint64
}

func (tx *Transaction) hash() []byte {
	var hash [32]byte
	var data []byte

	
	data = ToBytes(tx.Nonce)
	data = append(data, tx.Pubkey...)
	data = append(data, tx.Recipient...)
	data = append(data, ToBytes(tx.Value)...)

	hash = sha256.Sum256(data)

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
	tx.Nonce = 0
	tx.Pubkey = []byte{0x0000000000000000000000000000000000000001}
	tx.Sig = []byte{0x00000001}
	tx.Recipient = wallet.AddressToPKH(string(w.Address()))
	//miner prize
	tx.Value = 50

	return &tx
}

func NewTransaction(from *wallet.Wallet, to string, amount uint64, chain *BlockChain) *Transaction {
	balance, Nonce := chain.Acc.BalanceNonce(string(from.Address()))
	if balance < amount+1{
		log.Panic("not enough balance!")
	}
	tx := Transaction{}
	tx.Nonce = Nonce+1
	tx.Pubkey = from.PublicKey
	tx.Value = amount + 1//fee
	tx.Recipient = wallet.AddressToPKH(to)

	tx.Sign(from,chain)

	return &tx
}

func (tx *Transaction) IsCoinbase() bool { 
	if tx.Nonce == 0{
		
		if bytes.Equal(tx.Pubkey,[]byte{0x0000000000000000000000000000000000000001}){
			if bytes.Equal(tx.Sig, []byte{0x00000001}){
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
	balance, _ := chain.Acc.BalanceNonce(string(wallet.Address()))
	if (balance >= uint64(tx.Value)){

		r,s, err := ecdsa.Sign(rand.Reader, &wallet.PrivateKey, tx.Id())
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Sig = signature
		Handle(err)
	} else{
		log.Panic("not enough founds!")
	}

}
func VerifySignature(txid ,pubkey, sig []byte) bool{
	curve := elliptic.P256()

	r := big.Int{}
	s := big.Int{}

	sigLen := len(sig)
	r.SetBytes(sig[:(sigLen / 2)])
	s.SetBytes(sig[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubkey)
	x.SetBytes(pubkey[:(keyLen / 2)])
	y.SetBytes(pubkey[(keyLen / 2):])

	

	rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	if ecdsa.Verify(&rawPubKey, txid, &r, &s) == false {
		
		return false
	}
	return true
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
	lines = append(lines, fmt.Sprintf("       Nonce:     %d", tx.Nonce))
	lines = append(lines, fmt.Sprintf("       Public Key:       %x", tx.Pubkey))
	lines = append(lines, fmt.Sprintf("       Signature: %x", tx.Sig))
	lines = append(lines, fmt.Sprintf("       Recipient:    %s", tx.Recipient))
	lines = append(lines, fmt.Sprintf("       Value:    %d", tx.Value))

	return strings.Join(lines, "\n")
}
