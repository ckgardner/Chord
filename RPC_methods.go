package main

import (
	"fmt"
	"time"
)

// Ping methods will all be exported
func (node *Node) Ping(nothing Nothing, response *Nothing) error {
	return nil
}

// Set method exported
func (node *Node) Set(pair Pair, reply *Nothing) error {
	node.Lock.Lock()
	node.Bucket[pair.Key] = pair.Value
	node.Lock.Unlock()
	return nil
}

// Get method exported
func (node *Node) Get(Key string, res *string) error {
	node.Lock.Lock()
	for i := range node.Bucket {
		if i == Key {
			*res = node.Bucket[i]
		}
	}
	node.Lock.Unlock()
	return nil
}

// Delete method exported
func (node *Node) Delete(Key string, res *Nothing) error {
	node.Lock.Lock()
	delete(node.Bucket, Key)
	node.Lock.Unlock()
	return nil
}

func (node *Node) Notify(predecessor string, response *Nothing) {

	node.Lock.Lock()
	for true {

		time.Sleep(time.millisecond * 1333)
	}

	if node.Predecessor == "" {

		node.Predecessor = predecessor
	}
	node.Lock.Lock()
	return nil
}

func (node *Node) Stabilize(nothing Nothing, response *Nothing) {

	var predecessor string

	//calling predecessor
	if error := call(node.Successor, "Node.GetPredecessor", nothing, &predecessor); error != nil {
		log.printf("failed to get a connection: %v", error)
		return error
	}

	if predeccessor != "" {

		node.Successor = predecessor
	}

	//calling notify
	if error := call(node.Successor, "Node.Notify", node.address, &nothing); error != nil {
		log.printf("notifing the successor failed: %v", error)
		return error
	}

	log.printf("Successor Notified: %v", node.Successor)
	return nil
}

func (node *Node) GetPredecessor(nothing Nothing, response *string) {

	node.Lock.Lock()
	*response = node.Predecessor
	node.Lock.Lock()
	return nil
}
