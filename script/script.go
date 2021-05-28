package script

import (
	"crypto/sha256"
	"gochain/transaction"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)


type Script struct {
	stack *Stack
	cmd [][]byte
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
func (scr *Script) ScriptParser(s []byte) ([]byte, int) {
	length := len(s) 
	cmd := []byte{}
	count := 0
	
	for count < length{
		current := s[0]
		count++
		if current >= 1 && current <= 75{
			n := current
			cmd = append(cmd, n)
			count += int(n)
		}else if current == 76{
			
		}
	}
	return nil, length
}
func (src *Script) Serialize() []byte{
	var result []byte
	for _,cmd := range src.cmd{
		if _,ok := OP_CODE_FUNCTIONS[cmd[0]]; ok {
			result = append(result, cmd[0])
		}else{
			length := len(cmd)
			if length < 75{
				result = append(result, byte(length))
			}else if length > 75 && length < 0x100{
				result = append(result, 76)
				result = append(result, transaction.ToLittleEndian(length,1))
			}else if length >= 0x100 && length <= 520{
				result = append(result, 77)
				result = append(result, transaction.ToLittleEndian(length,2))
			}else{
				log.Panic("too long an cmd")
			}
			result = append(result, cmd...)
		}
	}
	total := len(result)
	result = append([]byte{byte(total)},result...)
	return result
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

func getkey(m map[byte]string, value string) (key byte, ok bool) {
	for k, v := range m {
	  if v == value { 
		key = k
		ok = true
		return
	  }
	}
	return
}

func handle(err error){
	if err != nil {
		log.Panic(err)
	}
}