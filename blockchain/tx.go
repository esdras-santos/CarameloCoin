package blockchain

import (
	"bytes"
	"encoding/binary"
	"gochain/wallet"
	
	"gochain/utils"
)



type TxInput struct {
	PrevTxID []byte // 32 bytes little-endian
	Out uint // 4 bytes little-endian
	ScriptSig []byte
	Sequence []byte // 4 bytes little-endian
}
func (in *TxInput) NewInput(prevTx,prevIndex,scriptSig,sequence []byte) {
	in.PrevTxID = prevTx
	in.Out = uint(binary.BigEndian.Uint64(prevIndex)) 
	in.ScriptSig = scriptSig
	in.Sequence = sequence
}
func (in TxInput) Serialize() []byte{
	result := utils.ToLittleEndian(in.PrevTxID,32)
	result = append(result, utils.ToLittleEndian([]byte{byte(in.Out)},4)...)
	result = append(result, in.ScriptSig...)
	result = append(result, utils.ToLittleEndian(in.Sequence,4)...)
	
	return result
}

// here we need to use the GetHeader command
func (in TxInput) Value() uint{
	chain := BlockChain{}
	tx,err := chain.FindTransaction(in.PrevTxID)
	Handle(err)
	return tx.Outputs[in.Out].Amount
}
func (in TxInput) ScriptpubKey(testnet bool) []byte{
	chain := BlockChain{}
	tx,err := chain.FindTransaction(in.PrevTxID)
	Handle(err)
	return tx.Outputs[in.Out].ScriptPubKey
}

func DeserializeInput(data []byte) (TxInput,int){
	var txin TxInput
	var lensc uint
	txin.PrevTxID = utils.ToLittleEndian(data[:33],32)
	txin.Out = uint(binary.BigEndian.Uint64(utils.ToLittleEndian(data[33:37],4)))
	txin.ScriptSig = data[37:]
	txin.Sequence = utils.ToLittleEndian(data[lensc+33 : lensc+37],4)
	return txin,len(data)
}

type TxOutput struct{
	Amount uint
	ScriptPubKey []byte
}
func (out TxOutput) Serialize()[]byte{
	amount := out.Amount
	result := utils.ToLittleEndian([]byte{byte(amount)},8)
	result = append(result, byte(len(out.ScriptPubKey)))
	result = append(result, out.ScriptPubKey...)

	return result
}

func DeserializeOutput(data []byte) (TxOutput,int){
	var txout TxOutput
	
	txout.Amount = uint(binary.BigEndian.Uint64(data[:8]))
	
	txout.ScriptPubKey = data[8:]
	return txout, len(data)
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool{
	lockingHash := wallet.PublicKeyHash(in.ScriptSig)
	return bytes.Compare(lockingHash,pubKeyHash) == 0
}

func (out *TxOutput) IsLockedWithKey(scriptPubKey []byte) bool{
	return bytes.Equal(out.ScriptPubKey,scriptPubKey)
}



