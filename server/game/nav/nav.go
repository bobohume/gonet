package detour

import (
	"gonet/base"
	"gonet/base/vector"
	"gonet/server/cm/lmath"
	"io/ioutil"
	"math"
	"strings"
	"sync"
	"unsafe"
)

// 坐标系
// Z				Y(0)
// *			  *
// *        *
// *     *
// *  *
// *  *  *  *  *  X(Pi/2)
type (
	Detour struct {
		mMaxNode int
		mMesh    *DtNavMesh
		mQuery   *DtNavMeshQuery
	}

	IDetour interface {
		Load(path string) int                                                       //加载mesh
		FindPath(start, end lmath.Point3F, path *vector.Vector[lmath.Point3F]) bool //巡径
		RayCast(start, end lmath.Point3F, path *vector.Vector[lmath.Point3F]) bool  //射线
		RandomPosition() (bool, lmath.Point3F)                                      //随机点
		GetPoly(pos lmath.Point3F) bool                                             //获取路径点
		GetPolyPos(pos lmath.Point3F) (bool, lmath.Point3F)                         //获取路径带高度
		TryMove(start, end lmath.Point3F) (bool, lmath.Point3F)

		LineTestCloseToEnd(start, end lmath.Point3F, pos *lmath.Point3F) bool //直线距离
		CanReach(pos lmath.Point3F) bool

		GetAreaWidth() float32                   //获取单个区域大小
		GetAreaNumX() int                        //x轴区域总数
		GetAreaNumY() int                        //y轴区域总数
		GetAreaNum() int                         //获取区域总数
		GetAreaPos(pos lmath.Point3F) (int, int) //坐标转到区域坐标
	}
)

var (
	mStaticMesh      map[string]*DtNavMesh = make(map[string]*DtNavMesh)
	mStaticMeshMutex sync.Mutex
)

const (
	FILE_SUFFIX_0 string = ".tile.bin"
	FILE_SUFFIX_1 string = ".nav"
	MAX_POLYS     int    = 256
)

// local x,y,z  unity is x,z,y
func SetF(pos *lmath.Point3F, f []float32) {
	pos.X, pos.Y, pos.Z = f[0], f[2], f[1]
}

func NewDetour(maxNode uint16) *Detour {
	DtAssert(maxNode != 0)
	this := &Detour{
		mMaxNode: int(maxNode),
	}
	return this
}

func (this *Detour) Load(path string) int {
	DtAssert(strings.HasSuffix(path, FILE_SUFFIX_0) || strings.HasSuffix(path, FILE_SUFFIX_1))
	errCode := 0
	this.mMesh = this.createStaticMesh(path, &errCode)

	if errCode != 0 {
		return errCode
	}

	this.mQuery = DtAllocNavMeshQuery()
	if this.mQuery == nil {
		return 3
	}

	status := this.mQuery.Init(this.mMesh, this.mMaxNode)
	if !DtStatusSucceed(status) {
		return 4
	}
	return 0
}

func (this *Detour) FindPath(start, end lmath.Point3F, path *vector.Vector[lmath.Point3F]) bool {
	filter := DtAllocDtQueryFilter()
	extents := []float32{1, 100, 1}

	straightPath := [MAX_POLYS * 3]float32{}
	straightPathFlags := [MAX_POLYS]DtStraightPathFlags{}
	straightPathPolys := [MAX_POLYS]DtPolyRef{}
	nstraightPath := 0

	var orgRef, dstRef DtPolyRef
	var org, dst [3]float32
	if !DtStatusSucceed(this.mQuery.FindNearestPoly(start.ToF(), extents, filter, &orgRef, org[:])) {
		return false
	}

	if !DtStatusSucceed(this.mQuery.FindNearestPoly(end.ToF(), extents, filter, &dstRef, dst[:])) {
		return false
	}

	path1 := make([]DtPolyRef, MAX_POLYS)
	pathCount := 0
	if !DtStatusSucceed(this.mQuery.FindPath(orgRef, dstRef, org[:], dst[:], filter, path1, &pathCount, MAX_POLYS)) {
		return false
	}

	if pathCount != 0 {
		epos1 := [3]float32{}
		//SetF(end, dst[:])

		if path1[pathCount-1] != dstRef {
			this.mQuery.ClosestPointOnPoly(path1[pathCount-1], dst[:], epos1[:], nil)
		}

		this.mQuery.FindStraightPath(org[:], dst[:], path1, pathCount, straightPath[:], straightPathFlags[:], straightPathPolys[:], &nstraightPath, MAX_POLYS, 0)
		for i := 0; i < nstraightPath*3; {
			pos := lmath.Point3F{}
			pos.X = straightPath[i]
			i++
			pos.Z = straightPath[i]
			i++
			pos.Y = straightPath[i]
			i++
			//fmt.Println(pos)
			path.PushBack(pos)
		}
	}

	return true
}

func (this *Detour) RayCast(start, end lmath.Point3F, path *vector.Vector[lmath.Point3F]) bool {
	if this.mQuery == nil {
		return false
	}

	filter := DtAllocDtQueryFilter()
	extents := []float32{1, 100, 1}
	endPos := end.ToF()

	var t float32 = 0
	var hitNormal [3]float32
	var polys [MAX_POLYS]DtPolyRef
	var startPolyRef DtPolyRef
	var startPos [3]float32
	if !DtStatusSucceed(this.mQuery.FindNearestPoly(start.ToF(), extents, filter, &startPolyRef, startPos[:])) {
		return false
	}

	npolys := 0
	status := this.mQuery.Raycast(startPolyRef, startPos[:], endPos, filter,
		&t, hitNormal[:], polys[:], &npolys, MAX_POLYS)
	if !DtStatusSucceed(status) {
		return false
	}
	bHit := (t <= 1)
	hitPos := [3]float32{}
	if bHit {
		DtVlerp(hitPos[:], startPos[:], endPos, t)
		if npolys > 0 {
			var h float32 = 0
			this.mQuery.GetPolyHeight(polys[npolys-1], hitPos[:], &h)
			hitPos[1] = h
		}
	}
	pos := lmath.Point3F{}
	SetF(&pos, hitPos[:])
	path.PushBack(pos)
	return true
}

func (this *Detour) RandomPosition() (bool, lmath.Point3F) {
	if this.mQuery == nil {
		return false, lmath.Point3F{}
	}

	randomPt := [3]float32{}
	filter := DtAllocDtQueryFilter()
	var ref DtPolyRef
	status := this.mQuery.FindRandomPoint(filter, func() float32 {
		return base.RandF[float32](0, 1)
	}, &ref, randomPt[:])
	if !DtStatusSucceed(status) {
		return false, lmath.Point3F{}
	}
	pos := lmath.Point3F{}
	SetF(&pos, randomPt[:])
	return true, pos
}

func (this *Detour) GetPoly(pos lmath.Point3F) bool {
	if this.mQuery == nil {
		return false
	}

	filter := DtAllocDtQueryFilter()
	extents := []float32{1, 100, 1}
	nearestPt := []float32{0, 0, 0}
	var nRef DtPolyRef
	status := this.mQuery.FindNearestPoly(pos.ToF(), extents, filter, &nRef, nearestPt)
	if !DtStatusSucceed(status) {
		return false
	}

	return nRef != 0
}

func (this *Detour) GetPolyPos(pos lmath.Point3F) (bool, lmath.Point3F) {
	if this.mQuery == nil {
		return false, lmath.Point3F{}
	}

	filter := DtAllocDtQueryFilter()
	extents := []float32{1, 100, 1}
	nearestPt := []float32{0, 0, 0}
	var nRef DtPolyRef
	status := this.mQuery.FindNearestPoly(pos.ToF(), extents, filter, &nRef, nearestPt)
	if !DtStatusSucceed(status) {
		return false, lmath.Point3F{}
	} else if nRef == 0 {
		return false, lmath.Point3F{}
	}

	nearestPt[0] = pos.X
	nearestPt[2] = pos.Y
	pos1 := lmath.Point3F{}
	SetF(&pos1, nearestPt)
	return true, pos1
}

func (this *Detour) TryMove(start, end lmath.Point3F) (bool, lmath.Point3F) {
	bHit := false
	filter := DtAllocDtQueryFilter()
	extents := []float32{1, 1, 1}
	if this.mQuery == nil {
		return false, lmath.Point3F{}
	}

	var startPolyRef, realEndPolyRef DtPolyRef
	startPos := []float32{0, 0, 0}
	realEndPos := []float32{0, 0, 0}
	if !DtStatusSucceed(this.mQuery.FindNearestPoly(start.ToF(), extents, filter, &startPolyRef, startPos)) {
		return false, lmath.Point3F{}
	}

	const visitedNodeCount = 16
	var visited [visitedNodeCount]DtPolyRef
	nvisited := 0
	status := this.mQuery.MoveAlongSurface(
		startPolyRef,
		startPos,
		end.ToF(),
		filter,
		realEndPos,
		visited[:],
		&nvisited,
		visitedNodeCount,
		&bHit)

	if DtStatusDetail(status, DT_INVALID_PARAM) {
		var tempRef DtPolyRef
		var tempPos [3]float32
		this.mQuery.FindNearestPoly(startPos, extents, filter, &tempRef, tempPos[:])
		startPolyRef = tempRef
		DtVcopy(startPos, tempPos[:])

		status = this.mQuery.MoveAlongSurface(
			startPolyRef,
			startPos,
			end.ToF(),
			filter,
			realEndPos,
			visited[:],
			&nvisited,
			visitedNodeCount,
			&bHit)
	}

	if !DtStatusSucceed(status) {
		return false, lmath.Point3F{}
	}

	realEndPolyRef = startPolyRef
	if nvisited > 0 {
		realEndPolyRef = visited[nvisited-1]
	}

	var h float32 = 0
	this.mQuery.GetPolyHeight(realEndPolyRef, realEndPos, &h)
	realEndPos[1] = h
	pos := lmath.Point3F{}
	SetF(&pos, realEndPos)
	return true, pos
}

func (this *Detour) LineTestCloseToEnd(start, end lmath.Point3F, pos *lmath.Point3F) bool {
	path := vector.New[lmath.Point3F]()
	if this.RayCast(start, end, path) != true {
		return false
	}

	*pos = path.Back()
	return true
}

func (this *Detour) CanReach(pos lmath.Point3F) bool {
	return this.GetPoly(pos)
}

func (this *Detour) createStaticMesh(path string, errCode *int) *DtNavMesh {
	mStaticMeshMutex.Lock()
	defer mStaticMeshMutex.Unlock()
	if m, ok := mStaticMesh[path]; ok {
		return m
	} else {
		mesh := this.loadStaticMesh(path, errCode)
		if *errCode == 0 {
			mStaticMesh[path] = mesh
		}
		return mesh
	}
}

func (this *Detour) GetAreaWidth() float32 {
	return float32(math.Max(float64(this.mMesh.mTileWidth), float64(1)))
}

func (this *Detour) GetAreaNumX() int {
	return int(this.mMesh.mBounds.Len_x()/this.mMesh.mTileWidth) + 1
}

func (this *Detour) GetAreaNumY() int {
	return int(this.mMesh.mBounds.Len_y()/this.mMesh.mTileWidth) + 1
}

func (this *Detour) GetAreaNum() int {
	return this.GetAreaNumX() * this.GetAreaNumY()
}

func (this *Detour) GetAreaPos(pos lmath.Point3F) (int, int) {
	x := int(math.Floor(float64(pos.X-this.mMesh.mOrig.X) / float64(this.mMesh.mTileWidth)))
	y := int(math.Floor(float64(pos.Y-this.mMesh.mOrig.Y) / float64(this.mMesh.mTileWidth)))
	return x, y
}

type NavMeshSetHeader struct {
	//magic    int32
	version  int32
	numTiles int32
	params   DtNavMeshParams
}
type NavMeshSetHeaderExt struct {
	boundsMinX float32
	boundsMinY float32
	boundsMinZ float32
	boundsMaxX float32
	boundsMaxY float32
	boundsMaxZ float32
}

type NavMeshTileHeader struct {
	tileRef  DtTileRef
	dataSize int32
}

const NAVMESHSET_MAGIC_RAW int32 = int32('M')<<24 | int32('S')<<16 | int32('E')<<8 | int32('T')
const NAVMESHSET_MAGIC_EXT int32 = int32('M')<<24 | int32('S')<<16 | int32('A')<<8 | int32('T')
const NAVMESHSET_VERSION int32 = 1
const TILECACHESET_MAGIC_RAW int32 = int32('T')<<24 | int32('S')<<16 | int32('E')<<8 | int32('T')
const TILECACHESET_MAGIC_EXT int32 = int32('T')<<24 | int32('S')<<16 | int32('A')<<8 | int32('T')
const TILECACHESET_VERSION int32 = 1

func (this *Detour) loadStaticMesh(path string, errCode *int) *DtNavMesh {
	*errCode = 0
	meshData, err := ioutil.ReadFile(path)
	if err != nil {
		*errCode = 101
		return nil
	}

	// Read header.
	header := (*NavMeshSetHeader)(unsafe.Pointer(&(meshData[0])))
	/*if header.magic != NAVMESHSET_MAGIC_RAW && header.magic != NAVMESHSET_MAGIC_EXT {
		*errCode = 103
		return nil
	}*/
	if header.version != NAVMESHSET_VERSION {
		*errCode = 104
		return nil
	}

	d := int32(unsafe.Sizeof(*header))
	/*if header.magic == NAVMESHSET_MAGIC_EXT {
		headerExt := (*NavMeshSetHeaderExt)(unsafe.Pointer(&(meshData[d])))
		d += int32(unsafe.Sizeof(*headerExt))
		this.mBounds.Min.X = headerExt.boundsMinX
		this.mBounds.Min.Z = headerExt.boundsMinY
		this.mBounds.Min.Y = headerExt.boundsMinZ
		this.mBounds.Max.X = headerExt.boundsMaxX
		this.mBounds.Max.Z = headerExt.boundsMaxY
		this.mBounds.Max.Y = headerExt.boundsMaxZ
	}*/

	mesh := DtAllocNavMesh()
	if mesh == nil {
		*errCode = 105
		return nil
	}
	state := mesh.Init(&header.params)
	if DtStatusFailed(state) {
		*errCode = 106
		return nil
	}

	// Read tiles.
	for i := 0; i < int(header.numTiles); i++ {
		tileHeader := (*NavMeshTileHeader)(unsafe.Pointer(&(meshData[d])))
		if tileHeader.tileRef == 0 || tileHeader.dataSize == 0 {
			break
		}
		d += int32(unsafe.Sizeof(*tileHeader))
		data := meshData[d : d+tileHeader.dataSize]
		state = mesh.AddTile(data, int(tileHeader.dataSize), DT_TILE_FREE_DATA, tileHeader.tileRef, nil)
		if DtStatusFailed(state) {
			*errCode = 108
			return nil
		}
		d += tileHeader.dataSize
	}

	mesh.mTileWidth = float32(math.Max(float64(header.params.TileWidth), float64(header.params.TileHeight)))
	mesh.mBounds.Min.SetMax(lmath.Point3F{0xFFFFFFF, 0xFFFFFFF, 0xFFFFFFF})
	mesh.mOrig.SetF(header.params.Orig[:])
	// 获取地图大小
	for _, v := range mesh.m_tiles {
		if v.Header == nil {
			continue
		}
		if mesh.mBounds.Min.X > v.Header.Bmin[0] {
			mesh.mBounds.Min.X = v.Header.Bmin[0]
		}
		if mesh.mBounds.Min.Z > v.Header.Bmin[1] {
			mesh.mBounds.Min.Z = v.Header.Bmin[1]
		}
		if mesh.mBounds.Min.Y > v.Header.Bmin[2] {
			mesh.mBounds.Min.Y = v.Header.Bmin[2]
		}
		if mesh.mBounds.Max.X < v.Header.Bmax[0] {
			mesh.mBounds.Max.X = v.Header.Bmax[0]
		}
		if mesh.mBounds.Max.Z < v.Header.Bmax[1] {
			mesh.mBounds.Max.Z = v.Header.Bmax[1]
		}
		if mesh.mBounds.Max.Y < v.Header.Bmax[2] {
			mesh.mBounds.Max.Y = v.Header.Bmax[2]
		}
	}
	mesh.mOrig = mesh.mBounds.Min

	return mesh
}
