package tile

import (
	"gonet/base"
	"container/heap"
	"math"
	"gonet/server/world/game/lmath"
)
//坐标系
//Z				Y
//*			  *
//*        *
//*     *
//*  *
//*  *  *  *  *  X
const(
	NaviGridSize = 1.0//网格大小
	AREA_SIZE = 10.0//区域大小
	// 单张地图的数据
	SingleTableWidth	     = (250 / NaviGridSize)
	SingleTableSize		 = (SingleTableWidth * SingleTableWidth)	// SingleTableWidth*SingleTableWidth
	AreaTableWidth	     = (SingleTableWidth / AREA_SIZE + 1)
	ms_maxSearchNode  	 = 100//寻路节点的限制
	NAV_LAYER_COUNT 		 = 8
)

type(
	Grid struct {
		flag int
	}

	//网格管理类
	NavigationMesh struct {
		m_Tile []*Grid
		m_OpenList *OpenHeap
	}

	INavigationMesh interface {
		Open(fileName string) bool//读取网格信息
		FindPath(start, end lmath.Point3F, path *base.Vector) bool
		GetGridFlag(row, col int) int
		GetTile(tile *Tile) base.Vector//a星获取周边网格
		CanReach(lmath.Point3F) bool//能够移动到网格
		LineTestCloseToEnd(start, end lmath.Point3F, pos *lmath.Point3F) bool
	}
)

//位置转网格
func PosToGrid(pos lmath.Point3F)(int, int){
	fx, fy := pos.X, pos.Y
	return  int(fx / NaviGridSize), int(fy / NaviGridSize)
}

func GeneratePosition(x, y int) lmath.Point3F{
	var pos lmath.Point3F
	pos.X = float32(x) * NaviGridSize
	pos.Y = float32(y) * NaviGridSize
	pos.Z = 0
	return pos
}

/*func GetGridId(x, y int) int{
	return ((x & 0x0000ffff) << 16) + y & 0x0000ffff
}*/

func GetGridId(x, y int) int{
	return x * int(SingleTableWidth) + y
}
//--------------NavigationMesh------------------//
func (this *NavigationMesh) Init(){
	this.m_Tile = make([]*Grid, SingleTableSize)
	for i := 0; i < SingleTableSize; i++{
		this.m_Tile[i] = &Grid{1}
	}
	this.m_OpenList = &OpenHeap{}
}

func (this *NavigationMesh) GetAreaWidth() float32{
	return AREA_SIZE
}

func (this *NavigationMesh) GetAreaNum() int{
	return AreaTableWidth
}

func (this *NavigationMesh) GetAreaPos(pos lmath.Point3F)(int, int){
	x := int(math.Floor(float64(pos.X-0) / float64(AREA_SIZE)))
	y := int(math.Floor(float64(pos.Y-0) / float64(AREA_SIZE)))
	return x, y
}

func (this *NavigationMesh) Open(fileName string) bool{
	return false
}

func (this *NavigationMesh) CanReach(pos lmath.Point3F) bool{
	row, col := PosToGrid(pos)

	if row < 0 || col < 0 || row >= SingleTableWidth || col >= SingleTableWidth {
		return false
	}

	if this.m_Tile == nil{
		return false
	}
	return this.GetGridFlag(row, col) != 0
}

func (this *NavigationMesh) GetGridFlag(row, col int) int{
	if row < 0 || col < 0 || row >= SingleTableWidth || col >= SingleTableWidth {
		return 0
	}
	if this.m_Tile != nil{
		return this.m_Tile[row * int(SingleTableWidth) + col].flag
	}
	return 0
}

func (this *NavigationMesh) GetTile(tile *Tile) (roundVec base.Vector) {
	xmin, xmax := lmath.ClampI(tile.x - 1, 0, SingleTableWidth - 1), lmath.ClampI(tile.x + 1, 0, SingleTableWidth - 1)
	ymin, ymax := lmath.ClampI(tile.y - 1, 0, SingleTableWidth - 1), lmath.ClampI(tile.y + 1, 0, SingleTableWidth - 1)
	for x := xmin; x <= xmax; x++{
		for y := ymin; y <= ymax; y++{
			if x != tile.x || y != tile.y{
				roundVec.Push_back(&Tile{x, y})
			}
		}
	}
	return roundVec
	/*for x := lmath.ClampI(tile.x - 1, 0, SingleTableWidth); x >= 0 && x < SingleTableWidth && x <= tile.x+1; x++{
		for y := lmath.ClampI(tile.y - 1, 0, SingleTableWidth); y >= 0 && y < SingleTableWidth && y <= tile.y+1; y++{
			if x != tile.x || y != tile.y{
				roundVec.Push_back(&Tile{x, y})
			}
		}
	}

	return roundVec*/
}

func (this *NavigationMesh) FindPath(start, end lmath.Point3F, path *base.Vector) bool {
	//openList := &OpenHeap{}
	openList := this.m_OpenList
	openList.m_Nodes.Clear()
	//closeList  := make(map[int] *ATile)
	closeList := [int(SingleTableWidth*SingleTableWidth)]bool{}
	//openSet := make(map[int] *ATile)
	heap.Init(openList)
	sx, sy := PosToGrid(start)
	ex, ey := PosToGrid(end)
	sx = lmath.ClampI(sx, 0, SingleTableWidth - 1)
	sy = lmath.ClampI(sy, 0, SingleTableWidth - 1)
	ex = lmath.ClampI(ex, 0, SingleTableWidth - 1)
	ey = lmath.ClampI(ey, 0, SingleTableWidth - 1)
	searchNum := 0
	endTile := NewATile(Tile{ex, ey}, nil, nil)
	heap.Push(openList, NewATile(Tile{sx, sy}, nil, nil))// 首先把起点加入开放列表
	for openList.Len() > 0{
		// 将节点从开放列表移到关闭列表当中。
		v := heap.Pop(openList)
		curPoint := v.(*ATile)
		id := GetGridId(curPoint.x, curPoint.y)
		//closeList[id] = curPoint
		closeList[id] = true
		roundVec := this.GetTile(&curPoint.Tile)
		searchNum++
		//超出最大寻径
		if searchNum > ms_maxSearchNode{
			for curPoint.father != nil{
				path.Push_front(GeneratePosition(curPoint.x, curPoint.y))
				curPoint = curPoint.father
			}
			return true
		}

		for _, t := range roundVec.Array(){
			tile := *t.(*Tile)
			curTile := NewATile(tile, curPoint, endTile)
			id := GetGridId(tile.x, tile.y)
			if curTile.IsEqual(endTile){
				// 找出路径了, 标记路径
				for curTile.father != nil{
					path.Push_front(GeneratePosition(curTile.x, curTile.y))
					curTile = curTile.father
				}
				return true
			}

			if this.GetGridFlag(curTile.x, curTile.y) == 0{
				//closeList[id] = curTile
				closeList[id] = true
				continue
			}

			//_, ok := closeList[id]
			ok := closeList[id]
			if ok{
				continue
			}else{
				heap.Push(openList, curTile)
			}

			//closeList[id] = curTile
			closeList[id] = true
			/*existPoint, ok := openSet[GetGridId(tile.x, tile.y)]
			if !ok {
				heap.Push(openList, curTile)
				openSet[GetGridId(tile.x, tile.y)] = curTile
			} else {
				oldGVal, oldFather := existPoint.gVal, existPoint.father
				existPoint.father = curTile
				existPoint.calcGVal()
				// 如果新的节点的G值还不如老的节点就恢复老的节点
				if existPoint.gVal > oldGVal {
					// restore father
					existPoint.father = oldFather
					existPoint.gVal = oldGVal
				}
			}*/
		}
	}
	return false
}

func (this *NavigationMesh) LineTestCloseToEnd(start, end lmath.Point3F, pos *lmath.Point3F) bool{
	sx, sy := PosToGrid(start)
	ex, ey := PosToGrid(end)
	sx = lmath.ClampI(sx, 0, SingleTableWidth - 1)
	sy = lmath.ClampI(sy, 0, SingleTableWidth - 1)
	ex = lmath.ClampI(ex, 0, SingleTableWidth - 1)
	ey = lmath.ClampI(ey, 0, SingleTableWidth - 1)
	xLast, yLast, xOffset, yOffset := -1, -1, 0, 0
	if sx > ex{
		xOffset = -1
	} else if sx < ex{
		xOffset = 1
	}

	if sy > ey{
		yOffset = -1
	} else if sy < ey{
		yOffset = 1
	}

	*pos = end
	// 两点下x坐标相近，这里可能有细微的误差
	if sx == ex || math.Abs(float64(start.X - end.X)) < 0.00001{
		if yOffset == 0{
			return true
		}

		if yOffset > 0 {
			for j :=sy;  j<= ey; j += yOffset {
				if this.GetGridFlag(j,sx) == 0{
					if yLast != -1 {
						*pos = GeneratePosition(sx, yLast)
						pos.X += 0.5 * NaviGridSize
						pos.Y += 0.5 * NaviGridSize
					} else{
						*pos = start
					}
					return false
				} else{
					yLast = j
				}
			}
		}else{
			for j :=sy;  j>= ey; j += yOffset {
				if this.GetGridFlag(j,sx) == 0{
					if yLast != -1 {
						*pos = GeneratePosition(sx, yLast)
						pos.X += 0.5 * NaviGridSize
						pos.Y += 0.5 * NaviGridSize
					} else{
						*pos = start
					}
					return false
				} else{
					yLast = j
				}
			}
		}
		return true
	}
	// 以下xOffset不可能为零
	if xOffset == 0{
		return false
	}

	yStart, yEnd := sy, sy
	k := (start.Y - end.Y)/(start.X - end.X)	// 斜率
	constant := (start.Y) - k * (start.X )		// 常数

	if xOffset > 0 {
		for i :=sx; i<=ex; i+=xOffset{
			if i >= SingleTableWidth || i < 0{
				continue
			}

			if xOffset > 0{
				yEnd = int((k * float32(i +  1) * NaviGridSize + constant) / NaviGridSize)
			}else{
				yEnd = int((k * float32(i +  0) * NaviGridSize + constant) / NaviGridSize)
			}
			// 最有一个点可能会超出ey的范围
			if yOffset>0{
				yEnd = lmath.ClampI(yEnd, sy, ey)
			} else{
				yEnd = lmath.ClampI(yEnd, ey, sy)
			}
			if yOffset > 0 {
				for j := yStart; j <= yEnd; j += yOffset {
					if j >= SingleTableWidth || j < 0{
						if 0 == yOffset{// yOffset为零直接跳出
							break
						}
						continue
					}
					if this.GetGridFlag(j, i) == 0 {
						if xLast != -1 && yLast != -1 && (sx != xLast || sy != yLast){
							*pos = GeneratePosition(xLast, yLast)
							pos.X += 0.5 * NaviGridSize
							pos.Y += 0.5 * NaviGridSize
						} else {
							*pos = start
						}
						return false
					} else {
						xLast = i
						yLast = j
					}
					if 0 == yOffset{// yOffset为零直接跳出
						break
					}
				}
			}else{
				for j := yStart; j >= yEnd; j += yOffset {
					if j >= SingleTableWidth || j < 0{
						if 0 == yOffset{// yOffset为零直接跳出
							break
						}
						continue
					}
					if this.GetGridFlag(j, i) == 0 {
						if xLast != -1 && yLast != -1 && (sx != xLast || sy != yLast){
							*pos = GeneratePosition(xLast, yLast)
							pos.X += 0.5 * NaviGridSize
							pos.Y += 0.5 * NaviGridSize
						} else {
							*pos = start
						}
						return false
					} else {
						xLast = i
						yLast = j
					}
					if 0 == yOffset{// yOffset为零直接跳出
						break
					}
				}
			}
			yStart = yEnd
		}
	}else{
		for i :=sx; i>=ex; i+=xOffset{
			if i >= SingleTableWidth || i < 0{
				continue
			}

			if xOffset > 0{
				yEnd = int((k * float32(i +  1) * NaviGridSize + constant) / NaviGridSize)
			}else{
				yEnd = int((k * float32(i +  0) * NaviGridSize + constant) / NaviGridSize)
			}
			// 最有一个点可能会超出ey的范围
			if yOffset>0{
				yEnd = lmath.ClampI(yEnd, sy, ey)
			} else{
				yEnd = lmath.ClampI(yEnd, ey, sy)
			}
			if yOffset > 0 {
				for j := yStart; j <= yEnd; j += yOffset {
					if j >= SingleTableWidth || j < 0{
						if 0 == yOffset{// yOffset为零直接跳出
							break
						}
						continue
					}
					if this.GetGridFlag(j, i) == 0 {
						if xLast != -1 && yLast != -1 && (sx != xLast || sy != yLast){
							*pos = GeneratePosition(xLast, yLast)
							pos.X += 0.5 * NaviGridSize
							pos.Y += 0.5 * NaviGridSize
						} else {
							*pos = start
						}
						return false
					} else {
						xLast = i
						yLast = j
					}
					if 0 == yOffset{// yOffset为零直接跳出
						break
					}
				}
			}else{
				for j := yStart; j >= yEnd; j += yOffset {
					if j >= SingleTableWidth || j < 0{
						if 0 == yOffset{// yOffset为零直接跳出
							break
						}
						continue
					}
					if this.GetGridFlag(j, i) == 0 {
						if xLast != -1 && yLast != -1 && (sx != xLast || sy != yLast){
							*pos = GeneratePosition(xLast, yLast)
							pos.X += 0.5 * NaviGridSize
							pos.Y += 0.5 * NaviGridSize
						} else {
							*pos = start
						}
						return false
					} else {
						xLast = i
						yLast = j
					}
					if 0 == yOffset{// yOffset为零直接跳出
						break
					}
				}
			}
			yStart = yEnd
		}
	}

	return true
}