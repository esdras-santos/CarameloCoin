package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"gochain/wallet"
)

type TxOutput struct{
	Value int
	scriptPubKey string
}

type TxOutputs struct{
	Outputs []TxOutput
}

type TxInput struct {
	ID []byte
	Out int
	PubKey []byte
	Signature []byte
}

type TxInputs struct{
	Inputs []TxInput
}

func NewTXOutput(value int,address string) *TxOutput{
	txo := &TxOutput{value, ""}
	txo.LockScript([]byte(address))

	return txo
}

func NewInput(id []byte, out int, sequence []byte) *TxInput{
	txi := &TxInput{id,out, "", sequence}
	return txi
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