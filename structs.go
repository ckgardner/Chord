package main

import (
	//"sync"
)

// Node has stuff
type Node struct {
	MyAddress 	string
	Port      	string
	Finger 		[]string
	Successors	[3]string
	Predecessor string
	Bucket    	map[string]string
	Ring      	bool
	kill      	chan struct{}
	//Lock      	sync.Mutex
	Ip        	string
	Next		int
}

// Nothing will do nothing
type Nothing struct{}

// Pair has a key-value pair
type Pair struct {
	Key   string
	Value string
}

// Commnad is a struct that has a verb and a function
type command struct {
}

// Finger is an address of a node
type Finger struct {
}

// FoundNode is a boolean
type FoundNode struct{
	Found bool
	Node string
}

// KeyFound = struct
type KeyFound struct{
	Found bool
	Address string 
}