package main

import (
	//"os"
	//"gochain/cli"

	"encoding/binary"
	"fmt"
	"math/big"
)

var FUNCTIONS = map[int]interface{}{
	// 0x00:"op_0",
	// 0x51:"op_1",
	// 0x60:"op_16",
	1: Add,
	// 0x93:"op_add",
	// 0xa9:"op_hash160",
	// 0xaa:"op_hash256",
	// 0xac:"op_checksig",
}

func Add(a int,b int) int{
	return a+b
}

func main(){
	//defer os.Exit(0) 
	//cmd := cli.CommandLine{}
	//cmd.Run()
	a := big.NewInt(510)
	ar := a.Bytes()
	for _,e := range ar{
		fmt.Printf("%x",e)
	}
	

}

func mapkey(m map[byte]string, value string) (key byte, ok bool) {
	for k, v := range m {
	  if v == value { 
		key = k
		ok = true
		return
	  }
	}
	return
}

func ReadVarint(s []byte, buf *uint){
	i := s[0]
	if i == 0xfd{
		a := binary.LittleEndian.Uint16(s[1:3])
		*buf = uint(a)
	}else if i == 0xfe{
		a := binary.LittleEndian.Uint32(s[1:5])
		*buf = uint(a)
	}else if i == 0xff{
		a := binary.LittleEndian.Uint64(s[1:9])
		*buf = uint(a)
	}else{
		*buf = uint(i)
	}
}

func toLittleEndian(bytes []byte) []byte{
	var le []byte
	for i := len(bytes)-1;i >= 0;i--{
		le = append(le, bytes[i]) 
	}
	return le
}