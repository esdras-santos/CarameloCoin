package transaction

import (
	"bytes"
	"encoding/binary"
	"gochain/script"
	"gochain/wallet"
	"math/big"
	"gochain/utils"
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
	result := utils.ToLittleEndian(in.PrevTxID,32)
	result = append(result, utils.ToLittleEndian(in.Out,4)...)
	result = append(result, in.ScriptSig...)
	result = append(result, utils.ToLittleEndian(in.Sequence,4)...)
	
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
	var lensc uint
	txin.PrevTxID = utils.ToLittleEndian(data[:33],32)
	txin.Out = utils.ToLittleEndian(data[33:37],4)
	sc := script.Script{}
	txin.ScriptSig,lensc = sc.ScriptParser(data[37:])
	txin.Sequence = utils.ToLittleEndian(data[lensc+33 : lensc+37],4)
	return txin,len(data)
}


type TxOutput struct{
	Amount *big.Int
	ScriptPubKey []byte
}
func (out TxOutput) Serialize()[]byte{
	amount := out.Amount.Bytes()
	result := utils.ToLittleEndian(amount,8)
	result = append(result, byte(len(out.ScriptPubKey)))
	result = append(result, out.ScriptPubKey...)

	return result
}

func DeserializeOutput(data []byte) (TxOutput,uint){
	var txout TxOutput
	var lensc uint
	txout.Amount.SetBytes(data[:8])
	sc := script.Script{}
	txout.ScriptPubKey,lensc = sc.ScriptParser(data[8:])
	return txout,lensc+8
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool{
	lockingHash := wallet.PublicKeyHash(in.ScriptSig)
	return bytes.Compare(lockingHash,pubKeyHash) == 0
}

func (out *TxOutput) IsLockedWithKey(scriptPubKey []byte) bool{
	return bytes.Equal(out.ScriptPubKey,scriptPubKey)
}



func Script() []byte{
	return nil
}