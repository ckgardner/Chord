package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func mainCommands(myNode *Node) {
	log.Printf("Chord is running")
	log.Printf("Command options: help, ping, port, create, join, dump, quit")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		parts := strings.SplitN(line, " ", 3)
		if len(parts) == 0 {
			continue
		}
		var nothing *Nothing
		switch parts[0] {

		case "help":
			fmt.Println("Commnad Options:")
			fmt.Println("ping, port, create, join, dump, quit")

		case "ping":
			if len(parts) == 1 {
				if err := call(myNode.MyAddress, "Node.Ping", myNode, &nothing); err != nil {
					log.Printf("server not reachable: %v", err)
				} else {
					log.Println("server responded!")
				}
			} else if len(parts) == 2 {
				if err := call(myNode.Ip+":"+parts[1], "Node.Ping", myNode, &nothing); err != nil {
					log.Printf("server not reachable: %v", err)
				} else {
					log.Printf("server responded!")
				}
			}
		case "port":
			if len(parts) == 1 {
				log.Printf("Current port: %v", myNode.Port)
			} else if len(parts) == 2 {
				if !myNode.Ring {
					myNode.Port = ":" + parts[1]
					myNode.MyAddress = myNode.Ip + myNode.Port
					log.Printf("Port is now %v", myNode.MyAddress)
				} else {
					log.Printf("Error, this node is already part of the ring")
				}
			} else {
				log.Println("Usage: port <number>")
			}

		case "create":
			var Error error
			if myNode.Ring == false {
				server(myNode)
				myNode.Ring = true
			} else {
				log.Printf("Ring already exists: %v\n", Error)
			}

		case "join":
			if len(parts) < 2 || len(parts) > 2 {
				log.Printf("join <address>")
				continue
			}
			if myNode.Ring == false {
				server(myNode)
				if err := join(myNode, parts[1]); err == nil {
					myNode.Ring = true
				} else {
					log.Printf("Invalid address: %v\n", err)
				}
			} else {
				log.Println("This ring already exists.")
			}

		case "get":
			if len(parts) == 2 {
				pair := Pair{parts[1], ""}
				call(myNode.find(parts[1]), "Node.Get", pair.Key, &pair.Value)
				println("The value for " + pair.Key + " is " + pair.Value)
			} else {
				fmt.Println("Get did not work")
			}

		case "delete":
			if len(parts) == 2 {
				pair := Pair{parts[1], ""}
				call(myNode.find(string(parts[1])), "Node.Delete", pair, &pair.Value)
				fmt.Println("Successfully removed:", pair.Key, pair.Value)
			}
		case "put":
			if len(parts) == 3 {
				pair := Pair{parts[1], parts[2]}
				call(myNode.find(parts[1]), "Node.Set", pair, &struct{}{})
				fmt.Println("You put", pair.Key, " & ", pair.Value, "on the ring")
			} else {
				fmt.Printf("put is not working")
			}

		case "dump":
			if len(parts) == 1{
				var dumpNode Node
				err := call(myNode.MyAddress, "Node.Dump", &struct{}{}, &dumpNode)
				if err == nil{
					fmt.Println("Address:	", dumpNode.MyAddress)
					fmt.Println("Predecessor:	", dumpNode.Predecessor)
					fmt.Println("Successors:	", dumpNode.Successors)
					fmt.Println("Bucket:		", dumpNode.Bucket)
					fmt.Println("Fingertable:	", dumpNode.Finger)
				}
			}

		case "quit":
			call(myNode.Successors[0], "Node.PutAll", myNode.Bucket, &struct{}{})
			os.Exit(3)

		default:
			log.Printf("I don't recognize this command")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error in main command loop: %v", err)
	}
}
