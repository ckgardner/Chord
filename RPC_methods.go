package main

import (
	"math/big"
	"fmt"
	"strconv"
)

// Ping methods will all be exported
func (node *Node) Ping(nothing Nothing, response *Nothing) error {
	return nil
}

// Set method exported
func (node *Node) Set(pair *Pair, reply *struct{}) error {
	node.Bucket[pair.Key] = pair.Value
	fmt.Println(*pair, "added to node")
	return nil
}

// Get method exported
func (node *Node) Get(Key string, res *string) error {
	if val, ok := node.Bucket[Key]; ok{
		*res = val
		return nil
	}
	return fmt.Errorf("pair does not exist for:", Key)
}

// Delete method exported
func (node *Node) Delete(pair *Pair, res *Nothing) error {
	if val, ok := node.Bucket[pair.Key]; ok{
		delete(node.Bucket, pair.Key)
		fmt.Println("Removed Pair: ", pair.Key, val)
		return nil
	}
	return fmt.Errorf("pair does not exist for pair:", pair)
}

// Notify method exported
func (node *Node) Notify(predecessor string, response *Nothing) error {
	if node.Predecessor == "" || between(hashString(node.Predecessor), hashString(predecessor), hashString(node.MyAddress), false){
		node.Predecessor = predecessor
	}
	return nil
}

func (node *Node) Dump(emp *struct{}, dump *Node)error{
	dump.MyAddress = node.MyAddress
	dump.Predecessor = node.Predecessor
	dump.Successors = node.Successors
	dump.Bucket = node.Bucket
	var old string
	for i := 0; i < len(node.Finger); i++{
		if old != node.Finger[i]{
			dump.Finger = append(dump.Finger, strconv.Itoa(i)+":\t", node.Finger[i], "\n\t\t\t")
			old = node.Finger[i]
		}
	}
	return nil
}
// FindSuccessor RPC
func (node *Node) FindSuccessor(id *big.Int, res *FoundNode) error{
	if between(hashString(node.MyAddress), id, hashString(node.Successors[0]), true){
		res.Node = node.Successors[0]
		res.Found = true
		return nil
	}
	res.Node = node.closestPrecedingNode(id)
	return nil
}
// PutAll RPC
func (node *Node) PutAll(bucket map[string]string, empty *struct{}) error{
	for k, v := range bucket{
		node.Bucket[k] = v
	}
	return nil
}
// GetAll RPC
func (node *Node) GetAll(address string, empty *struct{}) error{
	bucket := make(map[string]string)
	for k, v := range node.Bucket{
		if between(hashString(node.Predecessor), hashString(string(k)), hashString(address), false){
			bucket[k] = v
			delete(node.Bucket, k)
		}
	}
	call(address, "Node.PutAll", bucket, &struct{}{})
	return nil
}