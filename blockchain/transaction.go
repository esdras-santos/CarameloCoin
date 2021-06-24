package blockchain

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

func (tx *Transaction) Parse(data []byte) *Transaction {
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
	bc := BlockChain{}

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
	prevTXs := make(map[string]Transaction)
	
	for _,in := range tx.Inputs{
		prevTX, err := bc.FindTransaction(in.PrevTxID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.Id())] = prevTX
	}
	tx.VerifyTransaction(prevTXs)
	
	UTXO.Blockchain.SignTransaction(&tx, w.PrivateKey)

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	if len(tx.Inputs) == 1{
		if len(tx.Inputs[0].PrevTxID) == 0{
			if tx.Inputs[0].Out == 0xffffffff{
				return true
			}
		}
	}
	return false
}

func (tx *Transaction) Sign(wallet wallet.Wallet, prevTxs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.Inputs {
		index := 0
		txCopy.Inputs[inId].ScriptSig = nil
		for index == len(tx.Outputs){
			txCopy.Outputs[index].ScriptPubKey = prevTxs[hex.EncodeToString(in.PrevTxID)].Outputs[index].ScriptPubKey
			index++
		}

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &wallet.PrivateKey, []byte(dataToSign))
		Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)

		scriptsig := []byte{byte(len(signature))}
		scriptsig = append(scriptsig, signature...)
		scriptsig = append(scriptsig, []byte{byte(len(wallet.PublicKey))}...)
		scriptsig = append(scriptsig, wallet.PublicKey...)
		//p2pkh script
		tx.Inputs[inId].ScriptSig = scriptsig 
		for outId, _ := range txCopy.Outputs{
			txCopy.Outputs[outId].ScriptPubKey = nil
		}
	}
}


func (tx *Transaction) VerifyTransaction(UTXOs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}


	for inId, in := range tx.Inputs {
	
		scriptsig := tx.Inputs[inId].ScriptSig
		scriptpubKey := UTXOs[hex.EncodeToString(in.PrevTxID)].Outputs[in.Out].ScriptPubKey

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

func P2pkhScript(w wallet.Wallet) []byte{
	script := []byte{0x76}
	script = append(script, 0xa9)
	hash := wallet.PublicKeyHash(w.PublicKey)
	script = append(script, []byte{byte(len(hash))}...)
	script = append(script, hash...)
	script = append(script, 0x88)
	script = append(script, 0xac)
	return script
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
