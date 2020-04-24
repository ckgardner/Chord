package main

import (
	"fmt"
	"os"
	"time"
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
	myNode.Successors = [3]string{getLocalAddress() + myNode.Port}
	myNode.Finger = make([]string, 161)
	myNode.Next = 0
	fmt.Printf("myNode Address: %v\n", myNode.MyAddress)

	go func () {
        for {
            if myNode.Ring{
                time.Sleep(time.Millisecond * 1333)
                myNode.stabilize()
                time.Sleep(time.Millisecond * 1333)
                myNode.check_predecessor()
                time.Sleep(time.Millisecond * 1333)
                myNode.fix_fingers()
            }
        }
    }()

	go func(){
		<-myNode.kill
		os.Exit(0)
	}()
	

	mainCommands(myNode)
}

// Done by Cade Gardner & Andrew Nelson