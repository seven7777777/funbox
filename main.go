package main

import(
	"fmt"
	"strconv"
)


func main(){
	fmt.Println("Hello")
	tree:=NewTree(4)
	smp:=[]int{3,8,20,30,40,50,60,70,1,45,2,25,22,100,42,43,80,90,91,92,93,94,95}
	for _,v:=range smp{
		value:=[]byte("hello, tree at key "+ strconv.Itoa(v))
		tree.Insert(Key(v),value)
	}
	k,v:=tree.Insert(95,[]byte("new 95"))
	fmt.Println(k,string(v))
	k,v=tree.Insert(95,nil)
	fmt.Println(k,string(v))

	//tree.checkRight()
	//tree.Insert(30,nil)
	//tree.Insert(40,nil)
	//tree.Insert(50,nil)
	//tree.Insert(60,nil)
	tree.PrintTreeFinal()
	
//	x:=tree.tree.tail
//	fmt.Println("&&&&","len bound",len(x.(*TreeNode).bound), len(tree.tree.bound),len(tree.tree.bound[0].child.(*TreeNode).bound))
//	for _,bn:=range x.(*TreeNode).bound{
//		fmt.Println("&&&&&",bn.key)
//	}
	
//	x=x.(*TreeNode).tail
//		for _,bn:=range x.(*TreeNode).bound{
//		fmt.Println("*****",bn.key)
//	}
//	fmt.Println(x.(*TreeNode).parent,tree.tree.tail)
//	for _,bn:=range x.(*TreeNode).parent.(*TreeNode).bound{
//		fmt.Println("&&&&&",bn.key)
//	}
}