package main

import (
	//"os"
	//"gochain/cli"
	"encoding/binary"
	"fmt"
	
)

func main(){
	//defer os.Exit(0) 
	//cmd := cli.CommandLine{}
	//cmd.Run()
	var t []byte
	t = nil
	fmt.Println(t)
	s := []byte{0xfd,0xff,0x00}
	var b uint
	ReadVarint(s,&b)
	fmt.Println(b)
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