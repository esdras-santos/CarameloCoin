package script_test

import (
	"bytes"
	"crypto/sha256"

	"gochain/script"
	"log"
	"testing"

	"golang.org/x/crypto/ripemd160"
)

// var cmd = [][]byte{{0x01},{0x02}}
// var script = script.Script{&stack,cmd}

func TestSerialize(t *testing.T){

}

func TestScriptParser(t *testing.T){

}

func TestOpDup(t *testing.T){
	var stack script.Stack
	stack.Push([]byte{0x01})
	if stack.Size() != 1{
		t.Error("wrong size")
	}
	if  front, _ := stack.Front(); !bytes.Equal(front,[]byte{0x01}){
		t.Error("fail")
	}
	ok := script.OP_DUP(&stack)
	if !ok{
		t.Error("fail")
	}
	if stack.Size() != 2{
		t.Error("fail")
	}
	if front, _ := stack.Front(); !bytes.Equal(front,[]byte{0x01}){
		t.Error("fail")
	}
}

func TestOpHash256(t *testing.T){
	var stack script.Stack
	stack.Push([]byte{0x01})
	hash := sha256.Sum256([]byte{0x01}) 
	ok := script.OP_HASH256(&stack)
	if !ok{
		t.Error("fail")
	}
	if stack.Size() != 1{
		t.Error("fail")
	}
	if front, _ := stack.Front(); !bytes.Equal(front,hash[:]){
		t.Error("fail")
	}
}

func TestOpHash160(t *testing.T){
	var stack script.Stack
	data := sha256.Sum256([]byte{0x01})
	hasher := ripemd160.New()
	_,err := hasher.Write(data[:])
	if err != nil{
		log.Panic(err)
	}
	hash := hasher.Sum(nil)

	stack.Push([]byte{0x01})

	if stack.Size() != 1{
		t.Error("fail")
	}
	ok := script.OP_HASH160(&stack)
	if !ok{
		t.Error("fail")
	}
	front, _ := stack.Front()
	
	if !bytes.Equal(front,hash[:]){
		t.Error("fail")
	}
}

func TestOpEqualVerify(t *testing.T){
	var stack script.Stack
	stack.Push([]byte{0x01})
	stack.Push([]byte{0x01})

	if stack.Size() != 2{
		t.Error("fail")
	}
	if !script.OP_EQUALVERIFY(&stack){
		t.Error("fail")
	}
	if !stack.Empty(){
		t.Error("fail")
	} 
}

func TestOpAdd(t *testing.T){
	
}