package script

import (
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)


type Script struct {
	stack *Stack
}

var OP_CODE_FUNCTIONS = map[byte]string{
	0x00:"op_0",
	0x51:"op_1",
	0x60:"op_16",
	0x76:"op_dup",
	0x93:"op_add",
	0xa9:"op_hash160",
	0xaa:"op_hash256",
	0xac:"op_checksig",
}

//if is between 0x01 and 0x4b this is an element not an opcode
func (s *Script) ScriptParser(s []byte) ([]byte, int) {
	return nil, 1
}

func opDup(stack *Stack) bool {
	if stack.Size() < 1 {
		return false
	}
	opcode, err := stack.Front()
	handle(err)
	stack.Push(opcode)
	return true
}

func opHash256(stack *Stack) bool{
	if stack.Size() < 1{
		return false
	}
	data,err := stack.Front()
	handle(err)
	stack.Pop()
	hash := sha256.Sum256(data) 
	stack.Push(hash[:])
	return false
}

func opHash160(stack *Stack) bool{
	if stack.Size() < 1{
		return false
	}
	data,err := stack.Front()
	handle(err)
	hasher := ripemd160.New()
	_,err = hasher.Write(data[:])
	handle(err)
	hash := hasher.Sum(nil)
	stack.Push(hash[:])
	return true
}

func handle(err error){
	if err != nil {
		log.Panic(err)
	}
}