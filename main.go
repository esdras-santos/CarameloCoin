package main

import (
	//"os"
	//"gochain/cli"

	"bytes"
	"encoding/binary"
	"fmt"

	"log"
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
	add := "127.127.0.1"
	addb := []byte(add)
	fmt.Printf("%x",len(addb))
}

func ToHex(num int64) []byte{
	buff := new(bytes.Buffer)
	err := binary.Write(buff,binary.BigEndian,num)
	if err != nil{
		log.Panic(err)
	}	

	return buff.Bytes()
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


func toLittleEndian(bytes []byte) []byte{
	var le []byte
	for i := len(bytes)-1;i >= 0;i--{
		le = append(le, bytes[i]) 
	}
	return le
}