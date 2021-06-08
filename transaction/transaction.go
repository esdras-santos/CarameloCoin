package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"gochain/script"
	"gochain/utils"
	"gochain/wallet"
)

type Transaction struct {
	Version uint
	Locktime uint
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
	result := utils.ToLittleEndian([]byte{byte(tx.Version)},4)
	
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
	result = append(result, utils.ToLittleEndian([]byte{byte(tx.Locktime)},4)...)
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
	
	txn.Version = uint(binary.BigEndian.Uint64(utils.ToLittleEndian(data[:5],4)))
	
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
		startOut += int(len)
	}

	txn.Locktime = uint(binary.BigEndian.Uint64(utils.ToLittleEndian(data[startOut:],4)))

	txn.Testnet = testnet

	return txn
}

func (tx Transaction) Fee(testnet bool) uint{
	var inputSum, outputSum uint

	for _,txin := range tx.Inputs{
		inputSum += txin.Value(false)
	}
	for _,txout := range tx.Outputs{
		outputSum += txout.Amount
	}
	fee := inputSum - outputSum
	if outputSum + fee != inputSum{
		log.Panic("fee overflow")
	}
	return fee
}



func NewTransaction(w *wallet.Wallet, scriptPubKey []byte, amount int, UTXO *UTXOSet) *Transaction {
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
			input := TxInput{txID, uint(out), nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{uint(amount), scriptPubKey})

	if acc > amount {
		outputs = append(outputs, TxOutput{uint(acc-amount), w.PublicKey})
	}

	tx := Transaction{1,4294967295, inputs, outputs,false}
	UTXO.Blockchain.SignTransaction(&tx, w.PrivateKey)

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].PrevTxID) == 0 && tx.Inputs[0].Out == 0
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, in := range tx.Inputs {
		ptx := prevTXs[hex.EncodeToString(in.PrevTxID)]
		if ptx.Id == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	

	for inId, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.PrevTxID)]
		txCopy.Inputs[inId].ScriptSig = nil
		for outId, out := range prevTX.Outputs{
			txCopy.Outputs[outId].ScriptPubKey = out.ScriptPubKey
		}

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
		Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inId].ScriptSig = signature
		for outId, _ := range prevTX.Outputs{
			txCopy.Outputs[outId].ScriptPubKey = nil
		}
	}
}


func (tx *Transaction) VerifyTransaction(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, in := range tx.Inputs {
		prevtxs := prevTXs[hex.EncodeToString(in.PrevTxID)] 
		if prevtxs.Id() == nil {
			log.Panic("Previous transaction not correct")
		}
	}


	for inId, in := range tx.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.PrevTxID)]
		scriptsig := tx.Inputs[inId].ScriptSig
		scriptpubKey := prevTx.Outputs[in.Out].ScriptPubKey

		scriptin := script.Script{}
		scriptout := script.Script{}
		scriptin.Cmd,_ = scriptin.ScriptParser(scriptsig)
		scriptout.Cmd,_ = scriptin.ScriptParser(scriptpubKey)
		script := scriptin.Add(scriptout)
		dataToVerify := fmt.Sprintf("%x\n", tx)
		if !script.Evaluate([]byte(dataToVerify)){
			return false
		}
	}

	return true
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.PrevTxID, in.Out, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Amount, out.ScriptPubKey})
	}

	txCopy := Transaction{1,4294967295, inputs, outputs, false}

	return txCopy
}



func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.Id()))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:     %x", input.PrevTxID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Out))
		lines = append(lines, fmt.Sprintf("       Sript Signature: %x", input.ScriptSig))
		lines = append(lines, fmt.Sprintf("       Sequence:    %x", input.Sequence))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Amount))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.ScriptPubKey))
	}

	return strings.Join(lines, "\n")
}

func Handle(err error){
	if err != nil{
		log.Panic(err)
	}
}