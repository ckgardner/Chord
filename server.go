package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func server(myNode *Node) {
	location := myNode.Port
	rpc.Register(myNode)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", location)
	if err != nil {
		log.Fatal("Error thrown while listening: ", err)
	}
	fmt.Printf("Listening %v\n", location)

	go func() {
		if err := http.Serve(listener, nil); err != nil {
			log.Fatalf("Serving: %v", err)
		}
	}()
	fmt.Println("Server is on")
}

func join(myNode *Node, location string) error {
	var nothing Nothing
	if err := call(location, "Node.Ping", nothing, &nothing); err != nil {
		log.Printf("Connection not working: %v\n", err)
		return err
	}
	myNode.Successors[0] = location
	return nil
}

func call(address string, method string, request interface{}, reply interface{}) error {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Printf("rpc.DialHTTP: %v", err)
		return err
	}

	defer client.Close()

	return client.Call(method, request, reply)
}
