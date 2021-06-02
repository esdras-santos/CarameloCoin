package script

import (
	"bytes"
	"crypto/sha256"
	"gochain/utils"
	"gochain/wallet"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

// receives an byte and return a function
var OP_CODE_FUNCTIONS = map[byte]interface{}{
	0x00: OP_0,
	// 0x51:"op_1",
	// 0x60:"op_16",
	0x76: OP_DUP,
	// 0x93:"op_add",
	0xa9: OP_HASH160,
	0xaa: OP_HASH256,
	0xac: OP_CHECKSIG,
}

type Script struct {
	Stack *Stack
	Cmd [][]byte
}

/*before you pass the argument "transaction" you have to convert the transaction 
to string like that "dataToVerify := fmt.Sprintf("%x\n", transaction)" and then cast 
to array of bytes and pass as argument like that "script.Script.Evaluate([]byte(dataToVerify))"
*/
func (s *Script) Evaluate(transaction []byte) bool{
	cmds  := s.Cmd[:]
	stack := Stack{}
	altstack := Stack{}
	for len(cmds) > 0{
		cmd := append(cmds[:0],cmds[1:]...)
		if _,ok := OP_CODE_FUNCTIONS[cmd[0][0]]; ok {
			operation := OP_CODE_FUNCTIONS[cmd[0][0]]
			if 99 <= cmd[0][0] && cmd[0][0] <= 100{
				if !operation.(func(Stack,[][]byte)bool)(stack,cmds){
					log.Printf("bad op: %x",cmd[0])
					return false
				}
			}else if 107 <= cmd[0][0] && cmd[0][0] <= 108 {
				if !operation.(func(Stack,Stack)bool)(stack,altstack){
					log.Printf("bad op: %x",cmd[0])
					return false
				}
			}else if 172 <= cmd[0][0] && cmd[0][0] <= 175 {
				if !operation.(func(Stack,[]byte)bool)(stack,transaction){
					log.Printf("bad op: %x",cmd[0])
					return false
				}
			}else{
				if !operation.(func(Stack)bool)(stack){
					log.Printf("bad op: %x",cmd[0])
					return false
				}
			}
		}else{
			stack.Push(cmd[0])
		}
	}
	if stack.Empty(){
		return false
	}
	if c,_ := stack.Front();bytes.Equal(c,[]byte{}){
		return false
	}
	return true
}

//this is called before the Evaluate() function
func (scr Script) Add(other Script) Script {
	var s Script
	s.Cmd = append(scr.Cmd, other.Cmd...) 
	return s
}

//if is between 0x01 and 0x4b this is an element not an opcode
func (scr *Script) ScriptParser(s []byte) ([][]byte, uint) {
	var length uint
	utils.ReadVarint(s,&length)
	Cmd := [][]byte{}
	count := 0
	start := 0
	if length >= 0xfd{
		start = 2
	}else if length >= 0xfe{
		start = 3
	}else if length >= 0xff{
		start = 4
	}else{
		start = 1
	}
	for count < int(length){
		current := s[start]
		count++
		if current >= 1 && current <= 75{
			n := current
			Cmd = append(Cmd, s[start:n+1])
			count += int(n)
			start += count
		}else if current == 76{
			dataLen := current
			Cmd = append(Cmd, s[start:dataLen+1])
			count += int(dataLen)+1
			start += count
		}else if current == 77{
			dataLen := s[start+1]
			Cmd = append(Cmd, s[start:dataLen+1])
			count += int(dataLen)+2
			start += count
		}else{
			opcode := current
			Cmd = append(Cmd, []byte{opcode})
			start++
		}

	}
	return Cmd, length
}
func (src *Script) Serialize() []byte{
	var result []byte
	for _,Cmd := range src.Cmd{
		if _,ok := OP_CODE_FUNCTIONS[Cmd[0]]; ok {
			result = append(result, Cmd[0])
		}else{
			length := len(Cmd)
			if length < 75{
				result = append([]byte{byte(length)},result...)
			}else if length > 75 && length < 0x100{
				result = append(result, 76)
				result = append(utils.ToLittleEndian([]byte{byte(length)},1),result...)
			}else if length >= 0x100 && length <= 520{
				result = append(result, 77)
				result = append(utils.ToLittleEndian([]byte{byte(length)},2),result...)
			}else{
				log.Panic("too long an Cmd")
			}
			result = append(result, Cmd...)
		}
	}
	total := big.NewInt(int64(len(result)))
	buf := []byte{}
	utils.EncodeVarint(*total,&buf)
	result = append(buf,result...)
	return result
}

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
		return false
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