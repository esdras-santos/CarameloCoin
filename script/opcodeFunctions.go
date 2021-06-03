package script

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"gochain/wallet"
	"log"

	"golang.org/x/crypto/ripemd160"
)

func OP_0(stack *Stack)bool{
	stack.Push([]byte{0})
	return true
}

func OP_CHECKSIG(stack *Stack, transaction []byte) bool{
	if stack.Size() < 2{
		return false
	}
	pubkey,err := stack.Front()
	handle(err)
	err = stack.Pop()
	handle(err)
	signature,err := stack.Front()
	handle(err)
	leng := len(signature)
	signature = append(signature[:leng-1],signature[leng:]...)
	err = stack.Pop()
	handle(err)
	if wallet.VerifySignature(transaction,pubkey,signature){
		stack.Push([]byte{1})
	}else{
		stack.Push([]byte{0})
	}
	return true
}

func OP_DUP(stack *Stack) bool {
	if stack.Empty() {
		return false
	}
	opcode, err := stack.Front()
	handle(err)
	stack.Push(opcode)
	return true
}

func OP_HASH256(stack *Stack) bool{
	if stack.Size() < 1{
		return false
	}
	data,err := stack.Front()
	handle(err)
	stack.Pop()
	hash := sha256.Sum256(data) 
	stack.Push(hash[:])
	return true
}

func OP_HASH160(stack *Stack) bool{
	if stack.Size() < 1{
		return false
	}
	data,err := stack.Front()
	handle(err)
	pubHash := sha256.Sum256(data[:])
	hasher := ripemd160.New()
	_,err = hasher.Write(pubHash[:])
	handle(err)
	hash := hasher.Sum(nil)
	stack.Push(hash[:])
	return true
}

//0x88
func OP_EQUALVERIFY(stack *Stack) bool {
	hash1,err := stack.Front()
	handle(err)
	err = stack.Pop()
	handle(err)
	hash2,err := stack.Front()
	handle(err)
	err = stack.Pop()
	handle(err)

	if bytes.Equal(hash1,hash2){
		return true
	}else{
		log.Panic("diferent hashs")
		return false
	}
}

func OP_ADD(stack *Stack) bool{
	data1,err := stack.Front()
	number1 := binary.BigEndian.Uint64(data1)
	handle(err)
	err = stack.Pop()
	handle(err)
	data2,err := stack.Front()
	number2 := binary.BigEndian.Uint64(data2)
	handle(err)
	err = stack.Pop()
	handle(err)
	sum :=  number1 + number2
	stack.Push([]byte{byte(sum)})
	return true
}

func OP_MUL(stack *Stack) bool{
	data1,err := stack.Front()
	number1 := binary.BigEndian.Uint64(data1)
	handle(err)
	err = stack.Pop()
	handle(err)
	data2,err := stack.Front()
	number2 := binary.BigEndian.Uint64(data2)
	handle(err)
	err = stack.Pop()
	handle(err)
	sum :=  number1 * number2
	stack.Push([]byte{byte(sum)})
	return true
}

func OP_EQUAL(stack *Stack) bool {
	data1,err := stack.Front()
	handle(err)
	err = stack.Pop()
	handle(err)
	data2,err := stack.Front()
	handle(err)
	err = stack.Pop()
	handle(err)

	if (bytes.Equal(data1,data2)){
		stack.Push([]byte{1})
	}else{
		stack.Push([]byte{0})
	}
	return true
}
