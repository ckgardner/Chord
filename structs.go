package main

import(
	"sync"
)

// Node has stuff
type Node struct{
	MyAddress	string
	Port		string
	Successor	string
	Bucket		map[string]string
	Ring 		bool
	kill 		chan struct{}
	Lock		sync.Mutex
}

// Nothing will do nothing
type Nothing struct{}

// Pair has a key-value pair
type Pair struct{
	Key 	string
	Value	string
}