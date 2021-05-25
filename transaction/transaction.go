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

func ReadVarint(s []byte, buf *uint){
	i := s[0]
	if i == 0xfd{
		a := binary.LittleEndian.Uint16(s[1:3])
		*buf = uint(a)
	}else if i == 0xfe{
		a := binary.LittleEndian.Uint32(s[1:5])
		*buf = uint(a)
	}else if i == 0xff{
		a := binary.LittleEndian.Uint64(s[1:9])
		*buf = uint(a)
	}else{
		*buf = uint(i)
	}
}

func encodeVarint(i big.Int, buf *[]byte) {
	var bignum, ok = new(big.Int).SetString("0x10000000000000000", 0)
	ibytes := i.Bytes()
	lebytes := toLittleEndian(ibytes,4) 
	if !ok {
		log.Panic("fails to create the big number")
	}
	if cmp := i.Cmp(big.NewInt(0xfd));cmp < 0 {
		*buf = ibytes
	}else if cmp := i.Cmp(big.NewInt(0x10000));cmp < 0{
		*buf = lebytes
		*buf = append([]byte{0xfd},*buf...) 
	}else if cmp := i.Cmp(big.NewInt(0x100000000));cmp < 0{
		*buf = lebytes
		*buf = append([]byte{0xfe},*buf...)
	}else if cmp := i.Cmp(bignum);cmp < 0{
		*buf = lebytes
		*buf = append([]byte{0xff},*buf...)
	}else{
		log.Panic("integer too large")
	}
}

func (tx Transaction) Serialize() []byte {
	result := toLittleEndian(tx.Version.Bytes(),4)
	result = append(result, toLittleEndian(tx.Locktime.Bytes(),4)...)
	lenIns := big.NewInt(int64(len(tx.Inputs)))
	var lenInsEnc []byte
	encodeVarint(*lenIns,&lenInsEnc)
	result = append(result, lenInsEnc...)
	for i := 0; i < len(tx.Inputs);i++{
		result = append(result, tx.Inputs[i].Serialize()...)
	}

	lenOuts := big.NewInt(int64(len(tx.Outputs)))
	var lenOutsEnc []byte
	encodeVarint(*lenOuts,&lenOutsEnc)
	result = append(result, lenOutsEnc...)
	for i := 0; i < len(tx.Outputs);i++{
		result = append(result, tx.Outputs[i].Serialize()...)
	}
	

	return result
}

func DeserializeTransaction(data []byte, testnet bool) Transaction {
	var txn Transaction
	var lenIn uint
	ReadVarint([]byte{data[9]},&lenIn)
	var startIn int
	if lenIn <= 253{
		startIn = 10
	}else if lenIn <= 254{
		startIn = 11
	}else if lenIn <= 255{
		startIn = 12
	}
	
	txn.Version.SetBytes(toLittleEndian(data[:5],4))
	txn.Locktime.SetBytes(toLittleEndian(data[5:9],4))
	for i := 0;i<int(lenIn);i++{
		data, len := DeserializeInput(data[startIn:])
		txn.Inputs = append(txn.Inputs, data)
		startIn += len
	}

	var lenOut uint
	ReadVarint([]byte{data[startIn+1]},&lenOut)
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

	txn.Testnet = testnet

	return txn
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 24)
		_, err := rand.Read(randData)
		Handle(err)
		data = fmt.Sprintf("%x", randData)
	}

	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(20, to)

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	//tx.ID = tx.Hash()

	return &tx
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

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("Previous transaction not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range tx.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTx.Outputs[in.Out].PubKeyHash

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