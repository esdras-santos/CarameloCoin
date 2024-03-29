package blockchain

// import (
// 	"bytes"
// 	"encoding/binary"
	
	
// 	"gochain/utils"
// )



// type TxInput struct {
// 	PrevTxID []byte // 32 bytes little-endian
// 	Out uint // 4 bytes little-endian
// 	ScriptSig []byte
// 	Sequence []byte // 4 bytes little-endian
// }
// func (in *TxInput) NewInput(prevTx,prevIndex,scriptSig,sequence []byte) {
// 	in.PrevTxID = prevTx
// 	in.Out = uint(binary.BigEndian.Uint64(prevIndex)) 
// 	in.ScriptSig = scriptSig
// 	in.Sequence = sequence
// }
// func (in TxInput) Serialize() []byte{
// 	result := utils.ToLittleEndian(in.PrevTxID)
// 	result = append(result, utils.ToLittleEndian([]byte{byte(in.Out)})...)
// 	result = append(result, in.ScriptSig...)
// 	result = append(result, utils.ToLittleEndian(in.Sequence)...)
	
// 	return result
// }

// // here we need to use the GetHeader command
// func (in TxInput) Value() uint{
// 	chain := BlockChain{}
// 	tx,err := chain.FindTransaction(in.PrevTxID)
// 	Handle(err)
// 	return tx.Outputs[in.Out].Amount
// }


// func (in TxInput) ScriptpubKey() []byte{
// 	chain := BlockChain{}
// 	tx,err := chain.FindTransaction(in.PrevTxID)
// 	Handle(err)
// 	return tx.Outputs[in.Out].ScriptPubKey
// }

// func DeserializeInput(data []byte) (TxInput,int){
// 	var txin TxInput
// 	var lensc uint
// 	txin.PrevTxID = utils.ToLittleEndian(data[:33])
// 	txin.Out = uint(binary.BigEndian.Uint64(utils.ToLittleEndian(data[33:37])))
// 	txin.ScriptSig = data[37:]
// 	txin.Sequence = utils.ToLittleEndian(data[lensc+33 : lensc+37])
// 	return txin,len(data)
// }

// type TxOutput struct{
// 	Amount uint
// 	ScriptPubKey []byte
// }
// func (out TxOutput) Serialize()[]byte{
// 	amount := out.Amount
// 	result := utils.ToLittleEndian([]byte{byte(amount)})
// 	result = append(result, byte(len(out.ScriptPubKey)))
// 	result = append(result, out.ScriptPubKey...)

// 	return result
// }

// func SerializeOutputs(outs []TxOutput) []byte{
// 	var result []byte
// 	for _,i := range outs{
// 		utils.EncodeVarint(int64(len(i.Serialize())),&result)
// 		result = append(result, i.Serialize()...)
// 	}
// 	result = append(utils.ToHex(int64(len(outs))),result...)
// 	return result
// }

// func (out *TxOutput) Parse(data []byte) (TxOutput,int){
// 	var txout TxOutput
	
// 	txout.Amount = uint(binary.BigEndian.Uint64(data[:8]))
	
// 	txout.ScriptPubKey = data[8:]
// 	return txout, len(data)
// }

// func ParseOutputs(data []byte) ([]TxOutput){
// 	var start int
// 	var lenOut int
// 	var out TxOutput
// 	var outs []TxOutput
// 	utils.ReadVarint(data[1:],&lenOut)
// 	if lenOut <= 253{
// 		start = 2
// 	}else if lenOut <= 254{
// 		start = 3
// 	}else if lenOut <= 255{
// 		start = 4
// 	}

// 	end := start + lenOut
// 	lenOuts := binary.BigEndian.Uint64([]byte{data[0]})

// 	for i := 0;i<int(lenOuts);i++{
// 		o,_ := out.Parse(data[start:end])
// 		outs = append(outs, o)
// 		utils.ReadVarint(data[end+1:],&lenOut)
// 		if lenOut <= 253{
// 			start = end+2
// 		}else if lenOut <= 254{
// 			start = end+3
// 		}else if lenOut <= 255{
// 			start = end+4
// 		}
// 		end = start + lenOut
// 	}
// 	return outs
// }

// // func (in *TxInput) UsesKey(pubKeyHash []byte) bool{
// // 	lockingHash := wallet.PublicKeyHash(in.ScriptSig)
// // 	return bytes.Compare(lockingHash,pubKeyHash) == 0
// // }

// func (out *TxOutput) IsLockedWithKey(scriptPubKey []byte) bool{
// 	return bytes.Equal(out.ScriptPubKey,scriptPubKey)
// }
