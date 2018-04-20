package main

import (
	"fmt"
	"time"
	"base"
)

type(
	stBitTree struct{
		data int
		lchild *stBitTree
		rchild *stBitTree
	}
)

func inserTree(tree *stBitTree, key int) *stBitTree{
	if tree == nil{
		tree = &stBitTree{}
		tree.data = key
		return tree
	}

	if key < tree.data{
		tree.lchild = inserTree(tree.lchild, key)
	}else{
		tree.rchild = inserTree(tree.rchild, key)
	}
	return  tree
}

func midleOrder(tree *stBitTree){
	if tree != nil{
		midleOrder(tree.lchild)
		fmt.Println(tree.data)
		midleOrder(tree.rchild)
	}
}

func delteTree(tree *stBitTree, root *stBitTree, key int){
	if tree != nil{
		if tree.data == key {
			bLeft := root.lchild == tree
			if tree.rchild == nil{
				if bLeft{
					root.lchild = tree.lchild
				}else{
					root.rchild = tree.lchild
				}
			}else if(tree.lchild == nil){
				if bLeft{
					root.lchild = tree.rchild
				}else{
					root.rchild = tree.rchild
				}
			}else{
				if bLeft{
					root.lchild = tree.lchild
					tree.rchild = tree.rchild
				}else{
					root.rchild = tree.lchild
					tree.rchild = tree.rchild
				}
			}
		}

		delteTree(tree.lchild, tree, key)
		delteTree(tree.rchild, tree, key)
	}
}

func swap(a *int, b *int){
	temp := *a
	*a = *b
	*b = temp
}

func maxHeap(arr []int, len int){
	for i := len/2 -1; i >= 0; i-- {
		if 2 * i + 1 < len {
			if arr[2 * i + 1] >= arr[i]{
				swap(&arr[i], &arr[2 * i + 1])
			}
		}

		if  2 * i + 2 < len{
			if arr[2 * i + 2] >= arr[i]{
				swap(&arr[i], &arr[2 * i + 2])
			}
		}
	}
}

func minHeap(arr []int, len int){
	for i := len/2 - 1; i >= 0; i-- {
		if 2*i+1 < len {
			if arr[2*i+1] <= arr[i] {
				swap(&arr[i], &arr[2 * i + 1])
			}
		}

		if 2*i+2 < len {
			if arr[2*i+2] <= arr[i] {
				swap(&arr[i], &arr[2 * i + 2])
			}
		}
	}
}


func main() {
	n := []int{8, 2, 3, 5, 7, 1}
	base.SelectSort(n)
	fmt.Println(n)
	//base.QuickSort(n, 0, len(n)-1)
	fmt.Println(n)
	var tree *base.AvlBitTree
	tree = base.InsertAvlBitTree(tree, 3)
	tree = base.InsertAvlBitTree(tree, 2)
	tree = base.InsertAvlBitTree(tree, 1)
	tree = base.InsertAvlBitTree(tree, 0)
	tree = base.InsertAvlBitTree(tree, -1)
	fmt.Println("---------------")
	base.MiddleAvlPrint(tree)
	tree = base.DeleteAvlTree(tree, 2)
	//tree = base.DeleteAvlTree(tree, 3)
	fmt.Println("---------------")
	base.MiddleAvlPrint(tree)
	//base.InsertAvlBitTree(tree, 4)
	//base.InsertAvlBitTree(tree, 5)
	//base.InsertAvlBitTree(tree, 6)
	//base.MiddleAvlPrint(tree)
	for {
		time.Sleep(1000)
	}
}
