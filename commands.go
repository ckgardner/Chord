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
		var nothing Nothing
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
					myNode.MyAddress = defaultHost + ":" + parts[1]
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
				if Error = join(myNode, myNode.MyAddress); Error == nil {
					myNode.Ring = true
					log.Printf("\nAddress: %v", myNode.MyAddress)
				}
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
			//	if myNode.Ring == true {
			//		if len(parts) == 2 {
			//			fmt.Println("This appears before the get call")
			//			call(myNode.MyAddress, "myNode.Get", parts[1], myNode.Bucket[parts[1]])
			//		} else {
			//			fmt.Println("Get did not work")
			//		}
			//	}

			if len(parts) != 2 {
				log.Printf("Use get <key>")
				continue
			}
			var line string
			firstLine := line
			if err := call(myNode.MyAddress, "Node.Get", parts[1], &line); err != nil {
				log.Printf("calling myNode.Get %v", err)
			}
			if line == firstLine {
				fmt.Println("Not found")
			} else {
				fmt.Println("Value is: " + line)
			}

		case "delete":
			if len(parts) != 2 {
				log.Printf("delete <key>")
				continue
			}
			var line string
			if err := call(myNode.MyAddress, "Node.Delete", parts[1], &line); err != nil {
				log.Printf("myNode.Delete %v", err)
			} else {
				log.Printf("Key deleted")
			}

		case "put":
			//	if myNode.Ring == true {
			//		if len(parts) == 3 {
			//			pairval := Pair{parts[1], parts[2]}
			//			call(myNode.MyAddress, "myNode.Set", pairval, "put successful")
			//			fmt.Println("you added something to the ring")
			//		} else {
			//			fmt.Printf("put did not work")
			//		}
			//	}
			if len(parts) != 3 {
				log.Printf("put <key> <value>")
				continue
			}
			pair := Pair{parts[1], parts[2]}

			if err := call(myNode.MyAddress, "Node.Set", pair, &nothing); err != nil {
				log.Printf("myNode.Set: %v", err)
			} else {
				log.Printf("This was inserted to the Node: {%v:%v}", pair.Key, pair.Value)
			}

		case "dump":
			fmt.Printf("\nNode info\nLocal Node: %v\nSuccessor: %v\nBucket: \n", myNode.MyAddress, myNode.Successor)
			for i := range myNode.Bucket {
				fmt.Printf("\n{%v : %v} \n", i, myNode.Bucket[i])
			}

		case "quit":
			if myNode.Successor == myNode.MyAddress {
				log.Printf("This ring is now shutting down: %v", myNode.MyAddress)
				myNode.kill <- nothing
			} else {
				log.Printf("Killing node: %v", myNode.MyAddress)
				myNode.kill <- nothing
			}

		case "dumpkey":
			if len(parts) != 2 {
				log.Printf("Specify a port number")
				continue
			} else {
				for index := range myNode.Bucket {
					fmt.Printf("\t%v", index)

				}
			}

		default:
			log.Printf("I don't recognize this command")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error in main command loop: %v", err)
	}
}
