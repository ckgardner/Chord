package main

import (
	"sync"
)

// Node has stuff
type Node struct {
	MyAddress string
	Port      string
	Successor string
	Bucket    map[string]string
	Ring      bool
	kill      chan struct{}
	Lock      sync.Mutex
	Ip        string
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
type finger struct {
}
