package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"gochain/utils"
	"gochain/wallet"
)

type Transaction struct {
	Version *big.Int
	Locktime *big.Int
	//ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
	Testnet bool
}

func (tx *Transaction) hash() []byte {
	var hash [32]byte

	txCopy := *tx
	//txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func (tx *Transaction) Id() []byte{
	return tx.hash()
}




func (tx Transaction) Serialize() []byte {
	result := utils.ToLittleEndian(tx.Version.Bytes(),4)
	
	lenIns := big.NewInt(int64(len(tx.Inputs)))
	var lenInsEnc []byte
	utils.EncodeVarint(*lenIns,&lenInsEnc)
	result = append(result, lenInsEnc...)
	for i := 0; i < len(tx.Inputs);i++{
		result = append(result, tx.Inputs[i].Serialize()...)
	}

	lenOuts := big.NewInt(int64(len(tx.Outputs)))
	var lenOutsEnc []byte
	utils.EncodeVarint(*lenOuts,&lenOutsEnc)
	result = append(result, lenOutsEnc...)
	for i := 0; i < len(tx.Outputs);i++{
		result = append(result, tx.Outputs[i].Serialize()...)
	}
	result = append(result, utils.ToLittleEndian(tx.Locktime.Bytes(),4)...)
	return result
}

func DeserializeTransaction(data []byte, testnet bool) Transaction {
	var txn Transaction
	var lenIn uint
	utils.ReadVarint([]byte{data[5]},&lenIn)
	var startIn int
	if lenIn <= 253{
		startIn = 6
	}else if lenIn <= 254{
		startIn = 7
	}else if lenIn <= 255{
		startIn = 8
	}
	
	txn.Version.SetBytes(utils.ToLittleEndian(data[:5],4))
	
	for i := 0;i<int(lenIn);i++{
		data, len := DeserializeInput(data[startIn:])
		txn.Inputs = append(txn.Inputs, data)
		startIn += len
	}

	var lenOut uint
	utils.ReadVarint([]byte{data[startIn+1]},&lenOut)
	var startOut int
	if lenOut <= 253{
		startOut = startIn + 2
	}else if lenOut <= 254{
		startOut = startIn + 3
	}else if lenOut <= 255{
		startOut = startIn + 4
	}
	for i := 0;i<int(lenOut);i++{
		data,len := DeserializeOutput(data[startOut:])
		txn.Outputs = append(txn.Outputs, data)
		startOut += len
	}

	txn.Locktime.SetBytes(utils.ToLittleEndian(data[startOut:],4))

	txn.Testnet = testnet

	return txn
}

func (tx Transaction) Fee(testnet bool) *big.Int{
	inputSum , outputSum, aux := big.NewInt(0), big.NewInt(0), big.NewInt(0)

	for _,txin := range tx.Inputs{
		inputSum = aux.Add(inputSum,txin.Value(false))
	}
	for _,txout := range tx.Outputs{
		outputSum = aux.Add(outputSum,txout.Amount)
	}
	return aux.Sub(inputSum,outputSum)
}



func NewTransaction(w *wallet.Wallet, to string, amount int, UTXO *UTXOSet) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)
	acc, validOutputs := UTXO.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("Error: not enough funds")
	}

	// out and txId must be in little endian 32 bytes and 4 bytes respectivily
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	from := fmt.Sprintf("%s", w.Address())

	outputs = append(outputs, *NewTXOutput(amount, to))

	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXO.Blockchain.SignTransaction(&tx, w.PrivateKey)

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTX.Outputs[in.Out].PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
		Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inId].Signature = signature
		txCopy.Inputs[inId].PubKey = nil
	}
}


func (tx *Transaction) VerifyTransaction(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, in := range tx.Inputs {
		prevtxs := prevTXs[hex.EncodeToString(in.ID)] 
		if prevtxs.Id() == nil {
			log.Panic("Previous transaction not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range tx.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].ScriptSig = nil
		scriptpubKey := prevTx.Outputs[binary.BigEndian.Uint64(in.Out)].ScriptPubKey

		r := big.Int{}
		s := big.Int{}

		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Inputs[inId].PubKey = nil
	}

	return true
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}



func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:     %x", input.ID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Out))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}