package main

import(
	"log"
	"math/big"
	"net/rpc"
	"fmt"
)
// Stabilize method exported
func (node *Node) stabilize() error{
	var predecessor string
	var successors []string

	if err := call(node.Successors[0], "Node.GetSuccessors", struct{}{}, &successors); err != nil{
		log.Printf("could not get successors %v", err)
	}else{
		node.Successors[1] = successors[0]
		node.Successors[2] = successors[1]
	}

	if node.Successors[0] == ""{
		fmt.Println("Ending node is now", node.MyAddress)
		node.Successors[0] = node.MyAddress
	}else{
		node.Successors[0] = node.Successors[1]
		node.Successors[1] = node.Successors[2]
		node.Successors[2] = ""
	}

	pred := ""
	call(node.Successors[0], "Node.GetPredecessor", struct{}{}, &pred)

	if between(hashString(node.MyAddress),
	hashString(pred),
	hashString(node.Successors[0]),
	false) && pred != ""{
		node.Successors[0] = pred
	}
	

	if predecessor != node.MyAddress {
		//You Only Notify Successor If Current Successor Predecessor Should be You Instead
		if predecessor < node.MyAddress {
			if err := call(node.Successors[0], "Node.Notify", node.MyAddress, &struct{}{}); err != nil {
				log.Printf("notifing the successor failed: %v", err)
			}
		}else{
			node.Successors[0] = predecessor
			fmt.Println("the nodes successor is predecessor")
		}
	}

	return nil
}

func (node *Node) check_predecessor() error{
	if node.Predecessor != "" {
		client, err := rpc.DialHTTP("tcp", node.Predecessor)
		if err != nil{
			fmt.Println("Predecessor has failed: ", node.Predecessor)
			node.Predecessor = ""
		}else{
			client.Close()
		}
	}
	return nil
}

func (node *Node) GetPredecessor(empty *struct{}, predecessor *string) error{
	*predecessor = node.Predecessor
	return nil
}

func (node *Node) GetSuccessors(empty *struct{}, successors *[]string)error{
	*successors = node.Successors[:]
	return nil
}

// ClosestPrecedingNode is exported
func (node *Node) closestPrecedingNode(id *big.Int) string{
	for i := len(node.Finger) - 1; i > 0; i--{
		if between(hashString(node.MyAddress), hashString(node.Finger[i]), id, false){
			return node.Finger[i]
		}
	}
	
	return node.Successors[0]
}

func (node *Node) find(key string) string{
	foundNode := FoundNode{
		Found: false,
		Node: "",
	}
	max_steps := 30
	foundNode.Node = node.Successors[0]
	for !foundNode.Found{
		if max_steps > 0 {
			err := call(foundNode.Node, "Node.FindSuccessor", hashString(key), &foundNode)
			if err == nil{
				max_steps--
			}else{
				max_steps = 0
			}
		}else{
			return ""
		}
	}
	return foundNode.Node
}

func (node *Node) fix_fingers() error {
	node.Next++
	if node.Next > len(node.Finger)-1{
		node.Next = 0
	}
	bigInt := jump(node.MyAddress, node.Next)
	bigString := bigInt.String()
	address := node.find(bigString) //
	if node.Finger[node.Next] != address && address != ""{
		node.Finger[node.Next] = address
	}
	for{
		node.Next++
		if node.Next > len(node.Finger)-1{
			node.Next = 0
			return nil
		}
		if between(hashString(node.MyAddress), jump(node.MyAddress, node.Next), hashString(address), false) && address != ""{
			node.Finger[node.Next] = address
		}else{
			node.Next--
			return nil
		}
	}
}
