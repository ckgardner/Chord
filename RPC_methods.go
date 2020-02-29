package main

import (
	"log"
	"math/big"
)

// Ping methods will all be exported
func (node *Node) Ping(nothing Nothing, response *Nothing) error {
	return nil
}

// Set method exported
func (node *Node) Set(pair Pair, reply *Nothing) error {
	node.Bucket[pair.Key] = pair.Value
	return nil
}

// Get method exported
func (node *Node) Get(Key string, res *string) error {
	for i := range node.Bucket {
		if i == Key {
			*res = node.Bucket[i]
		}
	}
	return nil
}

// Delete method exported
func (node *Node) Delete(Key string, res *Nothing) error {
	delete(node.Bucket, Key)
	return nil
}

// Notify method exported
func (node *Node) Notify(predecessor string, response *Nothing) error {
	if node.Predecessor == "" || between(hashString(node.Predecessor), hashString(predecessor), hashString(node.MyAddress), false){
		node.Predecessor = predecessor
	}
	return nil
}

func (node *Node) FindSuccessor(id *big.Int, res *FoundNode) error{
	var value string
	last := value
	if err := call(node.Successors[0], "Node.Get", id, &value); err != nil{
		log.Printf("find successor node.get: %v", err)
	}else{
		if value != last{
			res.Found = true
			res.Node = node.Successors[0]
			log.Printf("Successor found: %v", node.Successors[0])
		}else{
			node.ClosestPrecedingNode(id, &res.Node)
			log.Printf("else: the next node is: %v", res.Node)
		}
	}
	return nil
}