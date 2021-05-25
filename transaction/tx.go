package transaction

import (
	"bytes"
	"encoding/binary"
	"gochain/script"
	"gochain/wallet"
	"math/big"
)



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
func (in TxInput) Serialize() []byte{
	result := toLittleEndian(in.PrevTxID,32)
	result = append(result, toLittleEndian(in.Out,4)...)
	result = append(result, in.ScriptSig...)
	result = append(result, toLittleEndian(in.Sequence,4)...)
	
	return result
}
func (in TxInput) FetchTx(testnet bool) *Transaction{
	fet := TxFetcher{}
	return fet.Fetch(in.PrevTxID,testnet,false)
}
func (in TxInput) Value(testnet bool) *big.Int{
	tx := in.FetchTx(false)
	return tx.Outputs[binary.BigEndian.Uint64(in.Out)].Amount
}
func (in TxInput) ScriptpubKey(testnet bool) []byte{
	tx := in.FetchTx(testnet)
	return tx.Outputs[binary.BigEndian.Uint64(in.Out)].ScriptPubKey
}

func DeserializeInput(data []byte) (TxInput,int){
	var txin TxInput
	var lensc int
	txin.PrevTxID = toLittleEndian(data[:33],32)
	txin.Out = toLittleEndian(data[33:37],4)
	txin.ScriptSig,lensc = script.ScriptParser(data[37:])
	txin.Sequence = toLittleEndian(data[lensc+33 : lensc+37],4)
	return txin,len(data)
}


type TxOutput struct{
	Amount *big.Int
	ScriptPubKey []byte
}
func (out TxOutput) Serialize()[]byte{
	amount := out.Amount.Bytes()
	result := toLittleEndian(amount,8)
	result = append(result, byte(len(out.ScriptPubKey)))
	result = append(result, out.ScriptPubKey...)

	return result
}

func DeserializeOutput(data []byte) (TxOutput,int){
	var txout TxOutput
	var lensc int
	txout.Amount.SetBytes(data[:8])
	txout.ScriptPubKey,lensc = script.ScriptParser(data[8:])
	return txout,lensc+8
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool{
	lockingHash := wallet.PublicKeyHash(in.ScriptSig)
	return bytes.Compare(lockingHash,pubKeyHash) == 0
}

func (out *TxOutput) IsLockedWithKey(scriptPubKey []byte) bool{
	return bytes.Equal(out.ScriptPubKey,scriptPubKey)
}

func toLittleEndian(bytes []byte, length int) []byte{
	le := make([]byte,length)
	for i := len(le)-1;i >= 0;i--{
		if bytes[i] != 0x00{
			le = append(le, bytes[i])
		}
		le = append(le, 0x00)
	}
	return le
}

func Script() []byte{
	return nil
}