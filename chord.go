package main

import (
	"fmt"
	"os"
)

const (
	defaultHost = "localHost"
	defaultPort = "3410"
)

func main() {
	myNode := new(Node)
	myNode.kill = make(chan struct{})
	myNode.Bucket = make(map[string]string)
	myNode.Ip = getLocalAddress()
	myNode.Port = ":" + defaultPort
	myNode.MyAddress = myNode.Ip + myNode.Port
	fmt.Printf("myNode Address: %v", myNode.MyAddress)

	go func(){
		<-myNode.kill
		os.Exit(0)
	}()

	mainCommands(myNode)
}

// Done by Cade Gardner & Andrew Nelson