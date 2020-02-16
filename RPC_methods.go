package main

func (node *Node) Ping(nothing Nothing, response *Nothing) error{
	return nil
}

func (node *Node) Put(pair Pair, reply *Nothing) error{
	node.Lock.Lock()
	node.Bucket[pair.Key] = pair.Value
	node.Lock.Unlock()
	return nil
}

func (node *Node) Get(Key string, res *string) error{
	node.Lock.Lock()
	for i := range node.Bucket{
		if i == Key{
			*res = node.Bucket[i]
		}
	}
	node.Lock.Unlock()
	return nil
}

func (node *Node) Delete(Key string, res *Nothing)error{
	node.Lock.Lock()
	delete(node.Bucket, Key)
	node.Lock.Unlock()
	return nil
}