package main

import(
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"net"
	"net/rpc"
	"net/http"
)

func mainCommands(myNode *Node){
	log.Printf("Chord is running")
	log.Printf("Command options: help, port, create, join, quit")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		words := strings.SplitN(line, " ", 3)
		if len(words) == 0{
			continue
		}
		var nothing Nothing
		switch words[0]{
		
		case "help":
			fmt.Println("Command options: help, port, create, join, quit")

		case "port":
			if len(words) == 1{
				log.Printf("Current port: %v", myNode.Port)
			} else if len(words) == 2{
				if myNode.Ring{
					myNode.Port =  ":" + words[1]
					myNode.MyAddress = defaultHost + ":" + words[1]
					fmt.Printf("Port is now %v", myNode.MyAddress)
				} else{
					fmt.Printf("Error, this node is already part of the ring")
				}
			} else{
				fmt.Println("Using", myNode.MyAddress)
			}

		case "create":
			var Error error
			if myNode.Ring{
				log.Printf("This node exists already!")
			} else{
				server(myNode)
				if Error = join(myNode, myNode.MyAddress); Error == nil{
					myNode.Ring = true
					log.Printf("\nAddress: %v", myNode.MyAddress)
				}
			}

		case "join":

		case "get":
			if len(words) != 2{
				log.Printf("key is <Key>")
				continue
			}
			var line string
			first_line := line
			if err := call(myNode.MyAddress, "myNode.Get", words[1], &line); err != nil{
				log.Printf("calling myNode.Get %v", err)
			}
			if line == first_line{
				fmt.Println("Not found")
			} else{
				fmt.Println("Value is: ", line)
			}

		case "delete":
			if len(words) != 2{
				log.Printf("delete <Key>")
				continue
			}
			var line string
			if err := call(myNode.MyAddress, "myNode.Delete", words[1], &line); err != nil{
				log.Printf("myNode.Delete %v", err)
			} else{
				log.Printf("Key deleted")
			}

		case "put":
			if len(words) != 3{
				log.Printf("put <Key> <Value>")
				continue
			}
			pair := Pair{words[1], words[2]}

			if err := call(myNode.MyAddress, "myNode.Put", pair, &nothing); err != nil{
				log.Printf("myNode.Post: %v", err)
			}else{
				log.Printf("This was inserted: %v", pair.Key, pair.Value)
			}

		case "quit":
			if myNode.Successor == myNode.MyAddress{
				log.Printf("This ring is now shutting down: %v", myNode.MyAddress)
				myNode.kill <- nothing
			} else{
				log.Printf("Killing node: %V", myNode.MyAddress)
				myNode.kill <- nothing
			}
			}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("in main command loop: %v")
	}
}
	
func server(myNode *Node){
	location := myNode.Port
	rpc.Register(myNode)
	rpc.HandleHTTP()
	listener, err := net.Listen("ip", location)
	if err != nil{
		log.Fatal("Error thrown while listening: ", err)
	}
	fmt.Printf("Listening %v", location)

	go func(){
		if err := http.Serve(listener, nil); err != nil {
			log.Fatalf("Serving: %V", err)
		}
	}()
	fmt.Println("Server is on")
}

func join(myNode *Node, location string) error {
	var nothing Nothing
	if err := call(location, "ping", nothing, &nothing); err != nil{
		log.Printf("Connection not working: %V", err)
		return err
	}
	myNode.Successor = location
	if err := call(myNode.MyAddress, "notify", myNode.MyAddress, &nothing); err != nil{
		log.Printf("Could not join", err)
	}
	return nil
}

func call(address string, method string, request interface{}, reply interface{}) error{
	user, err := rpc.DialHTTP("ip", address)
	if err != nil{
		log.Printf("rpc.DialHTTP:", err)
		return err
	}

	defer user.Close()

	if err = user.Call(method, request, reply); err != nil{
		log.Printf("user.Call %s: %v", method, err)
		return err
	}
	return nil
}