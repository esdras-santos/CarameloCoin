package script_test

import (
	"bytes"
	"crypto/sha256"
	"gochain/script"
	"log"
	"testing"

	"golang.org/x/crypto/ripemd160"
)

var stack *script.Stack = new(script.Stack)
// var cmd = [][]byte{{0x01},{0x02}}
// var script = script.Script{&stack,cmd}

func TestSerialize(t *testing.T){

}

func TestScriptParser(t *testing.T){

}

func TestOpDup(t *testing.T){
	stack.Push([]byte{0x01})
	if stack.Size() != 1{
		t.Error("wrong size")
	}
	if  front, _ := stack.Front(); !bytes.Equal(front,[]byte{0x01}){
		t.Error("fail")
	}
	ok := script.OpDup(stack)
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
	hash := sha256.Sum256([]byte{0x01}) 
	ok := script.OpHash256(stack)
	if !ok{
		t.Error("fail")
	}
	if stack.Size() != 2{
		t.Error("fail")
	}
	if front, _ := stack.Front(); !bytes.Equal(front,hash[:]){
		t.Error("fail")
	}
}

func TestOpHash160(t *testing.T){
	hasher := ripemd160.New()
	_,err := hasher.Write([]byte{0x01})
	if err != nil{
		log.Panic(err)
	}
	hash := hasher.Sum(nil)
	stack.Pop()
	if stack.Size() != 1{
		t.Error("fail")
	}
	ok := script.OpHash160(stack)
	if !ok{
		t.Error("fail")
	}
	if front, _ := stack.Front(); !bytes.Equal(front,hash[:]){
		t.Error("fail")
	}
}