package world

import (
	"fmt"
	"math"
)

type(
	//十字链路
	GameNode struct {
		xPrev *GameNode
		xNext *GameNode
		yPrev *GameNode
		yNext *GameNode

		X int
		Y int
	}
)

var(
	m_XNode *GameNode
	m_YNode *GameNode
)

//添加到十字链路
func AddNode(node* GameNode){
	//x handle
	var tail *GameNode
	bFind := false
	if m_XNode == nil || m_YNode == nil{
		m_XNode, m_YNode = node, node
		return
	}
	for curNode := m_XNode; curNode != nil; curNode = curNode.xNext{
		// insert data
		if curNode.X > node.X{
			node.xNext = curNode
			if 	curNode.xPrev != nil{
				node.xPrev = curNode.xPrev
				curNode.xPrev.xNext = node
			}else{
				m_XNode = node
			}
			curNode.xPrev = node
			bFind = true
			break
		}
		tail = curNode
	}

	if tail != nil && !bFind{
		tail.xNext = node
		node.xPrev = tail
	}

	tail = nil
	bFind = false
	//y handle
	for curNode := m_YNode; curNode != nil; curNode = curNode.yNext{
		// insert data
		if curNode.Y > node.Y{
			node.yNext = curNode
			if 	curNode.yPrev != nil {
				node.yPrev = curNode.yPrev
				curNode.yPrev.yNext = node
			} else{
				m_YNode = node
			}
			curNode.yPrev = node
			bFind = true
			break
		}
		tail = curNode
	}

	if tail != nil && !bFind{
		tail.yNext = node
		node.yPrev = tail
	}
}

//删除到十字链路
func  LeaveNode(node* GameNode){
	if node == m_XNode{
		if node.xNext != nil{
			m_XNode = node.xNext
			if m_XNode.xPrev != nil{
				m_XNode.xPrev = nil
			}
		}else{
			m_XNode = nil
		}
	}else if node.xPrev != nil && node.xNext != nil{
		node.xPrev.xNext = node.xNext
		node.xNext.xPrev = node.xPrev
	}else if node.xPrev != nil{
		node.xPrev.xNext = nil
	}

	if node == m_YNode{
		if node.yNext != nil{
			m_YNode = node.yNext
			if m_YNode.yPrev != nil{
				m_YNode.yPrev = nil
			}
		}else{
			m_YNode = nil
		}
	}else if node.yPrev != nil && node.yNext != nil{
		node.yPrev.yNext = node.yNext
		node.yNext.yPrev = node.yPrev
	}else if node.yPrev != nil{
		node.yPrev.yNext = nil
	}

	node.xPrev = nil
	node.xNext = nil
	node.yPrev = nil
	node.yNext = nil
}

//AOI
func AOI(node* GameNode, xLen, yLen float32) {
	// 往后找
	for cur := node.xNext; cur != nil; cur = cur.xNext{
		if float64(cur.X - node.X) > float64(xLen){
			break
		}else{
			if math.Abs(float64(cur.Y - node.Y)) <= float64(yLen) {
				fmt.Println(cur.X, cur.Y)
			}
		}
	}

	// 往前找
	for cur := node.xPrev; cur != nil; cur = cur.xPrev{
		if float64(node.X - cur.X) > float64(xLen){
			break
		}else{
			if math.Abs(float64(cur.Y - node.Y)) <= float64(yLen) {
				fmt.Println(cur.X, cur.Y)
			}
		}
	}
}

//move
//先LeaveNode
//再设置位置
//最后addNode

func testList(){
	node := [6]*GameNode{}
	node[0] = &GameNode{X:1, Y:5}
	node[1] = &GameNode{X:6, Y:6}
	node[2] = &GameNode{X:3, Y:1}
	node[3] = &GameNode{X:2, Y:2}
	node[4] = &GameNode{X:5, Y:3}
	node[5] = &GameNode{X:3, Y:3}
	for i := 0; i <6; i++{
		AddNode(node[i])
	}
	AOI(node[5], 2, 2)
	LeaveNode(node[4])
	LeaveNode(node[3])
	LeaveNode(node[2])
	LeaveNode(node[1])
	LeaveNode(node[0])
	LeaveNode(node[0])
}

