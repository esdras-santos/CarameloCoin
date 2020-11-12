package main

import (
	"os"
	"gochain/cli"
)

func main(){
	defer os.Exit(0) 
	cmd := cli.CommandLine{}
	cmd.Run()

}