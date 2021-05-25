package transaction

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
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
	result = append(result, byte(len(in.ScriptSig)))
	result = append(result, in.ScriptSig...)
	result = append(result, toLittleEndian(in.Sequence,4)...)
	
	return result
}
func (in TxInput) fetchTx(testnet bool) *Transaction{
	fet := TxFetcher{}
	return fet.Fetch(in.PrevTxID,testnet,false)
}
func (in TxInput) value(testnet bool) big.Int{
	tx := in.fetchTx(false)
	return *tx.Outputs[binary.BigEndian.Uint64(in.Out)].Amount
}
func (in TxInput) ScriptpubKey(testnet bool) []byte{
	tx := in.fetchTx(testnet)
	return tx.Outputs[binary.BigEndian.Uint64(in.Out)].ScriptPubKey
}

func DeserializeInput(data []byte) (TxInput,int){
	var txin TxInput
	txin.PrevTxID = toLittleEndian(data[:33],32)
	txin.Out = toLittleEndian(data[33:37],4)
	scriptLen := binary.BigEndian.Uint64(data[37:38])
	txin.ScriptSig = data[38:int(scriptLen)+39]
	txin.Sequence = toLittleEndian(data[int(scriptLen)+39:int(scriptLen)+39+4],4)
	return txin,(41+int(scriptLen))
}

func Script()[]byte{
	return nil
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
	txout.Amount.SetBytes(data[:8])
	spkLen := binary.BigEndian.Uint64(data[9:10])
	txout.ScriptPubKey = data[10:spkLen+10]
	return txout,(9+int(spkLen))
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool{
	lockingHash := wallet.PublicKeyHash(in.PubKey)
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