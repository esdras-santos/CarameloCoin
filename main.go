package main

import (
	"fmt"
	"time"
)

// 	"gochain/cli"

// 	"os"



func main(){
	// defer os.Exit(0) 
	// cmd := cli.CommandLine{}
	// cmd.Run()
	r := make(chan string,1)
	
	go func(){
		fmt.Println(<-r)
		fmt.Println("teste")
	}()
    fmt.Println("inicio")
	time.Sleep(time.Second * 5)
	r <- "meu pau"
	time.Sleep(time.Second * 1)
}	

