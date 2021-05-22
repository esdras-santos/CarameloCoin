package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"gochain/wallet"
)

type TxOutput struct{
	Amount int
	ScriptPubKey string
}

type TxInput struct {
	PrevTxID []byte // 32 bytes little-endian
	Out []byte // 4 bytes little-endian
	ScriptSig []byte
	Sequence []byte // 4 bytes little-endian
}
func (in *TxInput) NewInput(prevTx,prevIndex,scriptSig,sequence []byte) {
	in.PrevTxID = prevTx
	in.Out = prevIndex
	if scriptSig == nil{
		in.ScriptSig = Script()
	}else{
		in.ScriptSig = scriptSig
	}
	in.Sequence = sequence
}


func Script()[]byte{
	return nil
}

func NewTXOutput(value int,address string) *TxOutput{
	txo := &TxOutput{value, ""}
	txo.LockScript([]byte(address))

	return txo
}



func (outs TxOutputs) Serialize()[]byte{
	var buffer bytes.Buffer
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	Handle(err)
	return buffer.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs{
	var outputs TxOutputs
	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	Handle(err)
	return outputs
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool{
	lockingHash := wallet.PublicKeyHash(in.PubKey)
	return bytes.Compare(lockingHash,pubKeyHash) == 0
}

func (out *TxOutput) LockScript(address []byte){
	addr := wallet.Base58Decode(address)
	addr = addr[1:len(addr)-4]
	script := fmt.Sprintf("OP_DUP OP_HASH160 %s EQUALVERIFY CHECKSIG",string(addr[:]))
	out.scriptPubKey = script 
}

func (in *TxInput) UnlockScript(sig, pubK []byte){
	script := fmt.Sprintf("%s %s",string(sig[:]), string(pubK[:]))
	in.scripSig = script
}

func (out *TxOutput) IsLockedWithKey(scriptPubKey string) bool{
	return out.scriptPubKey == scriptPubKey	
}