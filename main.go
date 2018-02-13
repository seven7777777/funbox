package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Key int

func (k Key) Less(than Key) bool {
	return k < than
}

type KeyValue struct {
	key   Key
	value Key
	term  Node
}

type Node interface {
	Parent() Node
	Child() Node
	K() Key
	Insert(Key) *KeyValue
	Split(*KeyValue, Node, Node)
}

type TreeNode struct {
	keyValue *KeyValue
	kvs      []*KeyValue
	parent   Node
	bound    []Node
	degree   int
	child    Node
	tail     Node
}

func (this *TreeNode) Parent() Node {
	return this.parent
}
func (this *TreeNode) Child() Node {
	return this.child
}
func (this *TreeNode) K() Key {
	return this.keyValue.key
}

//back to parent after split
func (this *TreeNode) Split(kv *KeyValue, left Node, right Node) {
	if kv == nil {
		return
	}
	if len(this.bound) < this.degree {
		if kv == nil {
			return
		}
	}
	if this.bound[0].Child() == nil {
		upNode := this.bound[this.degree/2]
		kv := upNode.(*TreeNode).keyValue
		tmpBound := this.bound
		left = this
		left.(*TreeNode).bound = tmpBound[:this.degree/2]
		right = &TreeNode{}
		right.(*TreeNode).degree = this.degree
		right.(*TreeNode).bound = tmpBound[this.degree/2:]
		fmt.Println("after split----------")
		for _, node := range left.(*TreeNode).bound {
			fmt.Print(node.(*TreeNode).K(), " ")
		}
		fmt.Print("--")
		for _, node := range right.(*TreeNode).bound {
			fmt.Print(node.(*TreeNode).K(), " ")
		}
		fmt.Print("------------------------")
		left.(*TreeNode).tail = right
		if _, ok := this.parent.(*TreeNode); ok {
			fmt.Printf("leef parents bound len:%v %v", len(this.parent.(*TreeNode).bound), this.parent)
		}
		fmt.Println()
		this.parent.Split(kv, left, right)
	} else {
		var i int = 0
		var node Node
		for _, node = range this.bound {
			if node.Child() == left {
				break
			}
			i = i + 1
		}
		if i == 0 {
			newKV := &KeyValue{key: kv.key}
			newNode := &TreeNode{keyValue: newKV}
			newNode.child = left
			left.(*TreeNode).parent = this
			this.bound[0].(*TreeNode).child = right
			right.(*TreeNode).parent = this
			this.bound = append([]Node{newNode}, this.bound...)
		} else if i == len(this.bound) {
			this.tail = right
			right.(*TreeNode).parent = this
			newKV := &KeyValue{key: kv.key}
			newNode := &TreeNode{keyValue: newKV}
			newNode.child = left
			left.(*TreeNode).parent = this
			this.bound = append(this.bound, newNode)
		} else {
			node.(*TreeNode).child = right
			right.(*TreeNode).parent = this
			newKV := &KeyValue{key: kv.key}
			newNode := &TreeNode{keyValue: newKV}
			newNode.child = left
			left.(*TreeNode).parent = newNode

			preBound := append([]Node{}, this.bound[:i]...)
			postBound := append([]Node{}, this.bound[i:]...)
			this.bound = append(preBound, newNode)
			this.bound = append(this.bound, postBound...)
		}
		if len(this.bound) >= this.degree {
			upNode := this.bound[this.degree/2]
			kv := upNode.(*TreeNode).keyValue
			tmpBound := this.bound
			right = &TreeNode{}
			this.tail.(*TreeNode).parent = right
			right.(*TreeNode).tail = this.tail
			left = this
			left.(*TreeNode).tail = tmpBound[this.degree/2].(*TreeNode).child
			this.bound[this.degree/2].(*TreeNode).child = nil
			left.(*TreeNode).bound = tmpBound[:this.degree/2]
			right.(*TreeNode).degree = this.degree
			right.(*TreeNode).bound = tmpBound[this.degree/2+1:]
			if _, ok := this.parent.(*TreeNode); ok {
				fmt.Println("parents bound len", len(this.parent.(*TreeNode).bound))
			}
			this.parent.Split(kv, left, right)
		}
	}
	return
}

func (this *TreeNode) Insert(key Key) *KeyValue {
	kv := &KeyValue{key: key, value: key}
	newNode := &TreeNode{}
	newNode.keyValue = kv
	if this.bound[0].Child() == nil {

		if len(this.bound) >= this.degree {
			return nil
		}
		i := 0
		for _, node := range this.bound {
			if key.Less(node.(*TreeNode).K()) {
				break
			}
			i = i + 1
		}
		preBound := append([]Node{}, this.bound[:i]...)
		postBound := append([]Node{}, this.bound[i:]...)
		this.bound = append(preBound, newNode)
		this.bound = append(this.bound, postBound...)
		if len(this.bound) >= this.degree {
			skv := &KeyValue{}
			if _, ok := this.parent.(*TreeNode); ok {
				fmt.Println("parents bound len", len(this.parent.(*TreeNode).bound))
			}
			this.Split(skv, nil, nil)
		}
		return kv
	} else {
		for _, node := range this.bound {
			if key.Less(node.(*TreeNode).K()) {
				return node.Child().Insert(key)
			}
		}
		return this.tail.Insert(key)
	}
	return nil
}

type TreeRoot struct {
	count  int64
	degree int
	height int
	tree   Node
}

func (this *TreeRoot) Parent() Node {
	return nil
}
func (this *TreeRoot) Child() Node {
	return this.tree
}
func (this *TreeRoot) K() Key {
	return Key(0)
}
func (this *TreeRoot) Split(kv *KeyValue, left Node, right Node) {
	nbkv := &KeyValue{key: kv.key}
	newBoundNode := &TreeNode{}
	newBoundNode.keyValue = nbkv
	newBoundNode.child = left
	newBoundNode.degree = this.degree

	nkv := &KeyValue{key: kv.key}
	newNode := &TreeNode{}
	newNode.degree = this.degree
	newNode.keyValue = nkv
	left.(*TreeNode).parent = newNode
	newNode.tail = right
	newNode.bound = append(newNode.bound, newBoundNode)
	right.(*TreeNode).parent = newNode
	newNode.parent = this
	this.tree = newNode
	return
}
func (this *TreeRoot) Insert(key Key) (kv *KeyValue) {
	if this.tree == nil {
		kv := &KeyValue{key: key, value: key}
		treeNode := &TreeNode{}
		treeNode.keyValue = kv
		kv.term = treeNode
		boundNode := &TreeNode{}
		boundNode.degree = this.degree
		boundNode.keyValue = kv
		treeNode.kvs = append(treeNode.kvs, kv)
		treeNode.degree = this.degree
		treeNode.bound = append(treeNode.bound, boundNode)
		treeNode.parent = this
		this.tree = treeNode
		return nil
	} else {
		return this.tree.Insert(key)
	}
	return nil
}

func NewTree(degree int) *TreeRoot {
	tree := &TreeRoot{degree: degree}
	return tree
}
func (this *TreeRoot) Print() {
	if this.tree == nil {
		return
	}

	curNodes := []Node{this.tree}
	for {
		fmt.Println("-------layer-----------")
		var tmpCurNodes []Node
		if len(curNodes) == 0 {
			return
		}
		for _, node := range curNodes {
			isLeef := false
			for _, bnode := range node.(*TreeNode).bound {
				fmt.Print(bnode.K(), " ")
				if bnode.Child() != nil {
					tmpCurNodes = append(tmpCurNodes, bnode.Child())
				} else {
					isLeef = true
				}
			}
			fmt.Print("*")
			if node.(*TreeNode).tail != nil && !isLeef {
				tmpCurNodes = append(tmpCurNodes, node.(*TreeNode).tail)
			}
		}
		fmt.Println()
		curNodes = tmpCurNodes

	}
}
func main() {
	fmt.Println("hello, tree")
	tree := NewTree(3)
	tree.Insert(87)
	tree.Insert(61)
	tree.Insert(96)
	tree.Insert(46)
	tree.Insert(58)
	tree.Insert(70)
	tree.Insert(24)
	tree.Insert(99)
	tree.Insert(69)
	//	tree.Insert(30)
	//	tree.Insert(42)
	tree.Print()
	return
	rand.Seed(time.Now().Unix())
	for i := 1; i < 10; i++ {
		key := rand.Intn(100)
		fmt.Println("insert ", key)
		tree.Insert(Key(key))
	}
	tree.Print()
}
