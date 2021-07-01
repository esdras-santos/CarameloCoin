package main

import (
	//"os"
	//"gochain/cli"

	//"fmt"

	"fmt"
	"strconv"

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
	add := "127.0.0.1"
	addb := AddressToBytes(add)
	for i := 0;i<len(addb);i++{
		fmt.Printf("%x ",addb[i])
	}

	fmt.Println(AddressToString(addb))
}	

func AddressToString(addr []byte) string{//IPv4
	ip := fmt.Sprintf("%d.%d.%d.%d",addr[0],addr[1],addr[2],addr[3])
	return ip
}

func AddressToBytes(address string)[]byte{
	var number []rune
	var bytesIp []byte
	address = fmt.Sprintf("%s%s",address,".")
	for _,c := range address{
		if c == '.'{	
			i, err := strconv.Atoi(string(number))
			if err != nil{
				log.Panic(err)
			}
			
			bytesIp = append(bytesIp, byte(i))
			number = nil
			
		}else{
			number = append(number, c)
		}
	}
	return bytesIp
}