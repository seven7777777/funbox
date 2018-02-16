package main

import(
	"fmt"
)
type Key int
func (k Key)Less(than Key)bool{
	return k<than
}

type Node interface{
	//Child() Node
	Split(*TreeNode)
	Insert(Key, []byte)(Key,[]byte)
	SetNext(Node)
}

type Root struct{
	degree int
	final Node
	tree *TreeNode
}
type BoundNode struct{
	child Node
	final *FinalNode
	key Key
}

type TreeNode struct{
	degree int
	bound []*BoundNode
	parent Node //父亲节点
	tail Node //中间节点
	next Node //叶子节点的后驱指针
	pre Node //叶子节点的前驱指针
}

type KeyValue struct{
	key Key
	value []byte
}

type FinalNode struct{
	kv *KeyValue
}


func (tree *Root) PrintTreeFinal(){
	finalNode:=tree.final.(*TreeNode)
	for finalNode!=nil{
		fmt.Print("[")
		for _,kv:=range finalNode.bound{
			fmt.Print(" ",kv.key)
		}
		fmt.Print("]")
		//fmt.Println("node.pre == nil",finalNode.pre==nil)
		if finalNode.next !=nil{
			finalNode = finalNode.next.(*TreeNode)
		}else{
			finalNode = nil
		}
	}
}

func NewTree(degree int)*Root{
	if degree<3||degree>9999{
		return nil
	}
	root:=&Root{}
	tn:=&TreeNode{}
	tn.degree = degree
	tn.parent=root
	root.tree=tn
	root.final = tn
	root.degree=degree
	tn.pre=root
	return root
}

func (this *Root) Insert(key Key, value[]byte)(Key,[]byte){
	return this.tree.Insert(key,value)
}
func (this *Root) Split(node *TreeNode){
	//node = this.tree
	isTerminal := node.bound[0].child==nil

	left:=&TreeNode{degree:this.degree}
	right:=&TreeNode{degree:this.degree}
	newTree:=&TreeNode{degree:this.degree}
	if isTerminal{
		left.bound=append(left.bound,node.bound[:len(node.bound)/2]...)
		right.bound=append(right.bound,node.bound[len(node.bound)/2:]...)
	}else{
		left.bound=append(left.bound,node.bound[:len(node.bound)/2]...)
		right.bound=append(right.bound,node.bound[len(node.bound)/2+1:]...)
	}
	right.parent=newTree
	left.parent=newTree
	left.tail=node.bound[len(node.bound)/2].child
	right.tail=node.tail
	if node.tail!=nil{
		node.tail.(*TreeNode).parent=right
	}
	if isTerminal{
		this.SetNext(left)
		left.pre=this
		left.SetNext(right)
		right.pre=left
		right.SetNext(node.next)
		if right.next!=nil{
			right.next.(*TreeNode).pre=right
		}
	}else{
		for _,bn:=range left.bound{
			bn.child.(*TreeNode).parent=left
		}
		left.tail.(*TreeNode).parent=left
		for _,bn:=range right.bound{
			bn.child.(*TreeNode).parent=right
		}
		right.tail.(*TreeNode).parent=right
	}

	nbn:=NewBoundNode(false,node.bound[len(node.bound)/2].key,nil)
	nbn.child=left
	newTree.tail=right
	newTree.bound=append(newTree.bound,nbn)
	newTree.parent=this
	this.tree=newTree
	return
}
func (this *Root) SetNext(node Node){
	this.final=node
}
func NewBoundNode(isTerm bool, key Key, value []byte) *BoundNode{
	bn:=&BoundNode{}
	bn.key = key
	if isTerm{
		fn:=&FinalNode{}
		fn.kv=&KeyValue{key:key, value:value,}
		bn.final = fn
		bn.child = nil

	}
	return bn
}

func (this *TreeNode) Insert(key Key, value []byte)(Key, []byte){
	var i int=0
	var bn *BoundNode = nil
	var ins2node Node
	var flag bool
	for i,bn=range this.bound{
		if key.Less(bn.key){
			flag=true
			break
		}
		bn=nil
	}
	
	isTerm:=(len(this.bound)==0)|| this.bound[0].child==nil
	
	if bn!=nil{
		ins2node=bn.child
	}else{
		ins2node = this.tail
	}
	

	
	if isTerm{
		if !flag{
			i=len(this.bound)
		}
		if i-1>=0 && this.bound[i-1].key==key{
			oldValue:=this.bound[i-1].final.kv.value
			this.bound[i-1].final.kv.value=value
			return key,oldValue
		}
		nbn:=NewBoundNode(true,key,value)
		tmp:=append([]*BoundNode{},this.bound[:i]...)
		tmp=append(tmp,nbn)
		tmp=append(tmp,this.bound[i:]...)
		this.bound = tmp
	}else{
		return ins2node.Insert(key,value)
	}
	
	if len(this.bound)>=this.degree{
	this.parent.Split(this)
	}
	return key,value
}
func (this *TreeNode) SetNext(node Node){
	if len(this.bound)>0&&this.bound[0].final!=nil{
		this.next=node
	}
}

func (this *TreeNode) Split(node *TreeNode){
	var bn *BoundNode = nil
	var pos int=0
	var flag bool=false
	for pos,bn=range this.bound{
		if bn.child==node{
			flag=true
			break
		}
	}
	if flag ==false{
		pos=len(this.bound)
	}

	isTerminal := node.bound[0].child==nil

	left:=&TreeNode{degree:this.degree}
	right:=&TreeNode{degree:this.degree}

	if isTerminal{
		left.bound=append(left.bound,node.bound[:len(node.bound)/2]...)
		right.bound=append(right.bound,node.bound[len(node.bound)/2:]...)
	}else{
		left.bound=append(left.bound,node.bound[:len(node.bound)/2]...)
		right.bound=append(right.bound,node.bound[len(node.bound)/2+1:]...)
	}

	right.parent=this
	left.parent=this
	nbn:=NewBoundNode(false,node.bound[len(node.bound)/2].key,nil)
	if !isTerminal{
		left.tail=node.bound[len(node.bound)/2].child
		left.tail.(*TreeNode).parent=left
		for _,bn:=range left.bound{
			bn.child.(*TreeNode).parent=left
		}
		for _,bn:=range right.bound{
			bn.child.(*TreeNode).parent=right
		}
		right.tail=node.tail
		right.tail.(*TreeNode).parent=right
	}else{
		node.pre.SetNext(left)
		left.pre=node.pre
		left.SetNext(right)
		right.pre=left
		right.SetNext(node.next)
		if right.next!=nil{
			right.next.(*TreeNode).pre=right
		}
	}

	//insert new bound node
	if pos+1==len(this.bound) && false == flag{
		pos =len(this.bound)
	}
	tmpBound:=append([]*BoundNode{},this.bound[:pos]...)
	tmpBound=append(tmpBound,nbn)
	tmpBound=append(tmpBound,this.bound[pos:]...)
	nbn.child=left
	if pos==len(this.bound){
		this.tail=right
	}else{
		this.bound[pos].child=right
	}
	this.bound=tmpBound

	if len(this.bound)>=this.degree{
		this.parent.Split(this)
	}

}

func (this *Root) checkRight()bool{
	a:=this.tree
	i:=0
	for a!=nil{
		fmt.Println("layer:", i, "*****bad", len(a.bound)==0)
		for _,bn:=range a.bound{
			fmt.Print(bn.key, "  ")
		}
		fmt.Println()
		i++
		if a.tail!=nil{
		a=a.tail.(*TreeNode)
		}else{
			a=nil
		}
	}
	fmt.Println("\n------------")
	return true

}