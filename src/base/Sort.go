package base

import (
	"fmt"
	"math"
)
//-----------quick sort----------//
func QuickSort(arr []int, left int, right int){
	i,j := left, right
	key := arr[i]
	if i >= j{
		return
	}

	for i < j{
		for i < j && arr[j] > key{
			j--
		}
		arr[i] = arr[j]
		for i < j && arr[i] < key{
			i++
		}
		arr[j] = arr[i]
	}

	arr[i] = key
	QuickSort(arr, left, i-1)
	QuickSort(arr, i+1, right)
}

//-----------insert sort----------//
func InsertSort(arr []int){
	for i:=1; i < len(arr); i++{
		j:=i-1
		key := arr[i]
		for ; j>=0 && arr[j] < key;j--{
			arr[j+1] = arr[j]
		}
		arr[j+1] = key
	}
}

//------------select sort--------//
func SelectSort(arr []int){
	for i:=0; i < len(arr); i++ {
		key := i
		for j := i+1; j < len(arr); j++{
			if arr[j] > arr[key]{
				key = j
			}
		}
		swap(&arr[i], &arr[key])
	}
}

//-----------heap sort-----------//
func swap(a *int, b *int){
	c := *a
	*a = *b
	*b = c
}

func MaxHeap(arr []int){
	nLen := len(arr)
	if nLen <= 1{
		return
	}
	for i := nLen/2-1; i >= 0; i--{
		if 2*i+1 < nLen && arr[2*i+1] > arr[i]{
			swap(&arr[2*i+1], &arr[i])
		}
		if 2*i+2 < nLen && arr[2*i+2] > arr[i]{
			swap(&arr[2*i+2], &arr[i])
		}
	}
}

func MinHeap(arr []int){
	nLen := len(arr)
	if nLen <= 1{
		return
	}
	for i := nLen/2-1; i >= 0; i--{
		if 2*i+1 < nLen && arr[2*i+1] < arr[i]{
			swap(&arr[2*i+1], &arr[i])
		}
		if 2*i+2 < nLen && arr[2*i+2] < arr[i]{
			swap(&arr[2*i+2], &arr[i])
		}
	}
}

func PopHeap(arr []int, bMax bool) []int{
	nLen := len(arr)
	if nLen <= 1{
		return arr
	}
	swap(&arr[0], &arr[nLen-1])
	arr = arr[0:nLen-1]
	if bMax{
		MaxHeap(arr)
	}else{
		MinHeap(arr)
	}
	return arr
}

//-----------bittreee-------------------//
type(
	BitTree struct {
		data int
		lchild *BitTree
		rchild *BitTree
	}

	IBitTree interface {
		Less(int) bool
		Equal(int) bool
	}
)

func (this *BitTree) Less(data int) bool{
	return data < this.data
}

func (this *BitTree) Equal(data int) bool{
	return data == this.data
}

func InsertBitTree(tree *BitTree, data int)  *BitTree{
	if tree == nil{
		tree = &BitTree{}
		tree.lchild, tree.rchild = nil, nil
		tree.data = data
		return tree
	}

	if tree.Less(data){
		tree.lchild = InsertBitTree(tree.lchild, data)
	}else{
		tree.rchild = InsertBitTree(tree.rchild, data)
	}

	return tree
}

func DeleteTree(tree *BitTree, data int) *BitTree{
	return deletetree(tree, tree, &tree, data)
}

func MiddlePrint(tree *BitTree){
	if tree != nil{
		MiddlePrint(tree.lchild)
		fmt.Println(tree.data)
		MiddlePrint(tree.rchild)
	}
}

func deletetree(tree *BitTree, parent *BitTree, root **BitTree, data int) *BitTree{
	if tree != nil{
		if tree.Equal(data){
			if *root == tree {
				if tree.lchild == nil {
					*root = tree.rchild
					parent = *root
				} else if tree.rchild == nil {
					*root = tree.lchild
					parent = *root
				} else {
					if tree.lchild.rchild == nil {
						tree.lchild.rchild = tree.rchild
						*root = tree.lchild
						parent = *root
					} else {
						tree1, tree2 := tree.lchild, tree.lchild
						for ; tree1.rchild != nil; {
							tree2 = tree1
							tree1 = tree1.rchild
						}

						if tree1 != tree2 {
							if tree1.lchild != nil {
								tree2.rchild = tree1.lchild
							}

							tree2.rchild = nil
							tree1.lchild = tree.lchild
							tree1.rchild = tree.rchild
							*root = tree1
							parent = *root
						} else {
							tree1.rchild = tree.rchild
							*root = tree1
							parent = *root
						}
					}
				}
				tree = *root
			}
			if tree.lchild == nil{
				if parent.Less(data){
					parent.lchild  = tree.rchild
				}else{
					parent.rchild = tree.rchild
				}
			}else if tree.rchild == nil{
				if parent.Less(data){
					parent.lchild  = tree.rchild
				}else{
					parent.rchild = tree.rchild
				}
			}else{
				parent.lchild = tree.lchild
				tree.lchild.rchild = tree.rchild
			}
		} else if tree.Less(data){
			deletetree(tree.lchild, tree, root, data)
		}else{
			deletetree(tree.rchild, tree, root, data)
		}
	}
	return *root
}

//--------------avltree---------//
type(
	AvlBitTree struct {
		data int
		lchild *AvlBitTree
		rchild *AvlBitTree
		bt int
	}

	IAvlBitTree interface {
		Less(int) bool
		Equal(int) bool
	}
)

func (this *AvlBitTree) Less(data int) bool{
	return data < this.data
}

func (this *AvlBitTree) Equal(data int) bool{
	return data == this.data
}

func GetAvlHeight(tree *AvlBitTree) int{
	if tree == nil{
		return 0
	}

	return int(math.Max( float64(GetAvlHeight(tree.lchild)), float64(GetAvlHeight(tree.rchild)))) +1
}

func DeleteAvlTree(tree *AvlBitTree, data int) *AvlBitTree{
	tree1 := &tree
	delteavl(tree, tree, tree1, data)
	return *tree1
}

func InsertAvlBitTree(tree *AvlBitTree, data int)  *AvlBitTree {
	tree1 := &tree
	insertavl(tree, &tree, data)
	return *tree1
}

func R_Rorate(tree *AvlBitTree, root **AvlBitTree)*AvlBitTree{
	p := tree.lchild
	tree.lchild = p.rchild
	p.rchild = tree
	p.bt = GetAvlHeight(p.lchild) - GetAvlHeight(p.rchild)
	tree.bt = GetAvlHeight(tree.lchild) - GetAvlHeight(tree.rchild)
	if tree == *root{
		*root = p
	}
	tree = p
	return tree
}

func L_Rorate(tree *AvlBitTree, root **AvlBitTree)*AvlBitTree{
	p := tree.rchild
	tree.rchild = p.lchild
	p.lchild = tree
	p.bt = GetAvlHeight(p.lchild) - GetAvlHeight(p.rchild)
	tree.bt = GetAvlHeight(tree.lchild) - GetAvlHeight(tree.rchild)
	if tree == *root{
		*root = p
	}
	tree = p
	return tree
}

func LR_Rorate(tree *AvlBitTree, root **AvlBitTree)*AvlBitTree{
	tree.lchild = L_Rorate(tree.lchild, root)
	return R_Rorate(tree, root)
}

func RL_Rorate(tree *AvlBitTree, root **AvlBitTree)*AvlBitTree{
	tree.rchild = R_Rorate(tree.rchild, root)
	return L_Rorate(tree, root)
}

func MiddleAvlPrint(tree *AvlBitTree){
	if tree != nil{
		MiddleAvlPrint(tree.lchild)
		fmt.Println(tree.data, tree)
		MiddleAvlPrint(tree.rchild)
	}
}

func blanceavl(tree* AvlBitTree,root **AvlBitTree, data int)*AvlBitTree{
	tree.bt = GetAvlHeight(tree.lchild) - GetAvlHeight(tree.rchild)
	if tree.bt >= 2{
		if data < tree.lchild.data{
			tree = R_Rorate(tree, root)
		}else{
			tree = LR_Rorate(tree, root)
		}
	}

	if tree.bt <= -2{
		if data >= tree.rchild.data{
			tree = L_Rorate(tree, root)
		}else{
			tree = RL_Rorate(tree, root)
		}
	}
	return tree
}

func delteavl(tree *AvlBitTree, parent *AvlBitTree, root **AvlBitTree, data int) *AvlBitTree{
	if tree != nil{
		if tree.Equal(data){
			if *root == tree{
				if tree.lchild == nil{
					*root = tree.rchild
					parent = *root
				}else if tree.rchild == nil{
					*root = tree.lchild
					parent = *root
				}else{
					if tree.lchild.rchild == nil{
						tree.lchild.rchild = tree.rchild
						*root = tree.lchild
						parent = *root
					}else{
						tree1, tree2 := tree.lchild, tree.lchild
						for ;tree1.rchild != nil; {
							tree2 = tree1
							tree1 = tree1.rchild
						}

						if tree1 != tree2{
							if tree1.lchild != nil {
								tree2.rchild = tree1.lchild
							}

							tree2.rchild = nil
							tree1.lchild = tree.lchild
							tree1.rchild = tree.rchild
							*root = tree1
							parent = *root
						}else{
							tree1.rchild = tree.rchild
							*root = tree1
							parent = *root
						}
					}
				}
				tree = *root
			} else if tree.lchild == nil{
				if parent.Less(data){
					parent.lchild  = tree.rchild
				}else{
					parent.rchild = tree.rchild
				}
			}else if tree.rchild == nil {
				if parent.Less(data) {
					parent.lchild = tree.rchild
				} else {
					parent.rchild = tree.rchild
				}
			}else{
				parent.lchild = tree.lchild
				tree.lchild.rchild = tree.rchild
			}
		} else if tree.Less(data){
			delteavl(tree.lchild, tree, root, data)
		}else{
			delteavl(tree.rchild, tree, root, data)
		}
		blanceavl(tree, root, data)
	}
	return tree
}

func insertavl(tree *AvlBitTree, root **AvlBitTree, data int)  *AvlBitTree{
	if tree == nil{
		tree = &AvlBitTree{}
		tree.lchild, tree.rchild = nil, nil
		tree.data = data
		tree.bt = 0
		if *root == nil{
			*root = tree
		}
		return tree
	}

	if tree.Less(data){
		tree.lchild = InsertAvlBitTree(tree.lchild, data)
	}else{
		tree.rchild = InsertAvlBitTree(tree.rchild, data)
	}

	return blanceavl(tree, root, data)
}

