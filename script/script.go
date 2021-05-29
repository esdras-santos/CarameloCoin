package script

import (
	"crypto/sha256"
	"gochain/utils"
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
func (scr *Script) ScriptParser(s []byte) ([][]byte, uint) {
	var length uint
	utils.ReadVarint(s,&length)
	cmd := [][]byte{}
	count := 0
	start := 0
	if length >= 0xfd{
		start = 3
	}else if length >= 0xfe{
		start = 4
	}else if length >= 0xff{
		start = 5
	}else{
		start = 1
	}
	for count < int(length){
		current := s[start]
		count++
		if current >= 1 && current <= 75{
			n := current
			cmd = append(cmd, s[start:n+1])
			count += int(n)
			start += count
		}else if current == 76{
			dataLen := current
			cmd = append(cmd, s[start:dataLen+1])
			count += int(dataLen)+1
			start += count
		}else if current == 77{
			dataLen := s[start+1]
			cmd = append(cmd, s[start:dataLen+1])
			count += int(dataLen)+2
			start += count
		}else{
			opcode := current
			cmd = append(cmd, []byte{opcode})
			start++
		}

	}
	return cmd, length
}
func (src *Script) Serialize() []byte{
	var result []byte
	for _,cmd := range src.cmd{
		if _,ok := OP_CODE_FUNCTIONS[cmd[0]]; ok {
			result = append(result, cmd[0])
		}else{
			length := len(cmd)
			if length < 75{
				result = append([]byte{byte(length)},result...)
			}else if length > 75 && length < 0x100{
				result = append(result, 76)
				result = append(utils.ToLittleEndian([]byte{byte(length)},1),result...)
			}else if length >= 0x100 && length <= 520{
				result = append(result, 77)
				result = append(utils.ToLittleEndian([]byte{byte(length)},2),result...)
			}else{
				log.Panic("too long an cmd")
			}
			result = append(result, cmd...)
		}
	}
	total := big.NewInt(int64(len(result)))
	buf := []byte{}
	utils.EncodeVarint(*total,&buf)
	result = append(buf,result...)
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