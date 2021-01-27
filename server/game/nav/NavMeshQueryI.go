//
// Copyright (c) 2009-2010 Mikko Mononen memon@inside.org
//
// This software is provided 'as-is', without any express or implied
// warranty.  In no event will the authors be held liable for any damages
// arising from the use of this software.
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would be
//    appreciated but is not required.
// 2. Altered source versions must be plainly marked as such, and must not be
//    misrepresented as being the original software.
// 3. This notice may not be removed or altered from any source distribution.
//

package detour

import (
	"math"
	"unsafe"
)

/// @class dtQueryFilter
///
/// <b>The Default Implementation</b>
///
/// At construction: All area costs default to 1.0.  All flags are included
/// and none are excluded.
///
/// If a polygon has both an include and an exclude flag, it will be excluded.
///
/// The way filtering works, a navigation mesh polygon must have at least one flag
/// set to ever be considered by a query. So a polygon with no flags will never
/// be considered.
///
/// Setting the include flags to 0 will result in all polygons being excluded.
///
/// <b>Custom Implementations</b>
///
/// DT_VIRTUAL_QUERYFILTER must be defined in order to extend this class.
///
/// Implement a custom query filter by overriding the virtual passFilter()
/// and getCost() functions. If this is done, both functions should be as
/// fast as possible. Use cached local copies of data rather than accessing
/// your own objects where possible.
///
/// Custom implementations do not need to adhere to the flags or cost logic
/// used by the default implementation.
///
/// In order for A* searches to work properly, the cost should be proportional to
/// the travel distance. Implementing a cost modifier less than 1.0 is likely
/// to lead to problems during pathfinding.
///
/// @see dtNavMeshQuery

func (this *DtQueryFilter) constructor() {
	this.m_includeFlags = 0xffff
	this.m_excludeFlags = 0
	for i := 0; i < DT_MAX_AREAS; i++ {
		this.m_areaCost[i] = 1.0
	}
}

func (this *DtQueryFilter) destructor() {
}

func (this *DtQueryFilter) PassFilter(_ DtPolyRef, _ *DtMeshTile, poly *DtPoly) bool {
	return (poly.Flags&this.m_includeFlags) != 0 && (poly.Flags&this.m_excludeFlags) == 0
}

func (this *DtQueryFilter) GetCost(pa, pb []float32,
	_ DtPolyRef, _ *DtMeshTile, _ *DtPoly,
	_ DtPolyRef, _ *DtMeshTile, curPoly *DtPoly,
	_ DtPolyRef, _ *DtMeshTile, _ *DtPoly) float32 {
	return DtVdist(pa, pb) * this.m_areaCost[curPoly.GetArea()]
}

const H_SCALE float32 = 0.999 // Search heuristic scale.

//////////////////////////////////////////////////////////////////////////////////////////

/// @class dtNavMeshQuery
///
/// For methods that support undersized buffers, if the buffer is too small
/// to hold the entire result set the return status of the method will include
/// the #DT_BUFFER_TOO_SMALL flag.
///
/// Constant member functions can be used by multiple clients without side
/// effects. (E.g. No change to the closed list. No impact on an in-progress
/// sliced path query. Etc.)
///
/// Walls and portals: A @e wall is a polygon segment that is
/// considered impassable. A @e portal is a passable segment between polygons.
/// A portal may be treated as a wall based on the dtQueryFilter used for a query.
///
/// @see dtNavMesh, dtQueryFilter, #dtAllocNavMeshQuery(), #dtAllocNavMeshQuery()

/// @name Getters and setters for the default implementation data.
///@{

/// Returns the traversal cost of the area.
///  @param[in]		i		The id of the area.
/// @returns The traversal cost of the area.
func (this *DtQueryFilter) GetAreaCost(i int) float32 {
	return this.m_areaCost[i]
}

/// Sets the traversal cost of the area.
///  @param[in]		i		The id of the area.
///  @param[in]		cost	The new cost of traversing the area.
func (this *DtQueryFilter) SetAreaCost(i int, cost float32) {
	this.m_areaCost[i] = cost
}

/// Returns the include flags for the filter.
/// Any polygons that include one or more of these flags will be
/// included in the operation.
func (this *DtQueryFilter) GetIncludeFlags() uint16 {
	return this.m_includeFlags
}

/// Sets the include flags for the filter.
/// @param[in]		flags	The new flags.
func (this *DtQueryFilter) SetIncludeFlags(flags uint16) {
	this.m_includeFlags = flags
}

/// Returns the exclude flags for the filter.
/// Any polygons that include one ore more of these flags will be
/// excluded from the operation.
func (this *DtQueryFilter) GetExcludeFlags() uint16 {
	return this.m_excludeFlags
}

/// Sets the exclude flags for the filter.
/// @param[in]		flags		The new flags.
func (this *DtQueryFilter) SetExcludeFlags(flags uint16) {
	this.m_excludeFlags = flags
}

///@}
/// Gets the node pool.
/// @returns The node pool.
func (this *DtNavMeshQuery) GetNodePool() *DtNodePool {
	return this.m_nodePool
}

/// Gets the navigation mesh the query object is using.
/// @return The navigation mesh the query object is using.
func (this *DtNavMeshQuery) GetAttachedNavMesh() *DtNavMesh {
	return this.m_nav
}

func (this *DtNavMeshQuery) constructor() {

}

func (this *DtNavMeshQuery) destructor() {
	if this.m_tinyNodePool != nil {
		DtFreeNodePool(this.m_tinyNodePool)
		this.m_tinyNodePool = nil
	}
	if this.m_nodePool != nil {
		DtFreeNodePool(this.m_nodePool)
		this.m_nodePool = nil
	}
	if this.m_openList != nil {
		DtFreeNodeQueue(this.m_openList)
		this.m_openList = nil
	}
}

/// Initializes the query object.
///  @param[in]		nav			Pointer to the dtNavMesh object to use for all queries.
///  @param[in]		maxNodes	Maximum number of search nodes. [Limits: 0 < value <= 65535]
/// @returns The status flags for the query.
/// @par
///
/// Must be the first function called after construction, before other
/// functions are used.
///
/// This function can be used multiple times.
func (this *DtNavMeshQuery) Init(nav *DtNavMesh, maxNodes int) DtStatus {
	if maxNodes > int(DT_NULL_IDX) || maxNodes > int((1<<DT_NODE_PARENT_BITS)-1) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	this.m_nav = nav

	if this.m_nodePool == nil || this.m_nodePool.GetMaxNodes() < uint32(maxNodes) {
		if this.m_nodePool != nil {
			DtFreeNodePool(this.m_nodePool)
			this.m_nodePool = nil
		}
		this.m_nodePool = DtAllocNodePool(uint32(maxNodes), DtNextPow2(uint32(maxNodes/4)))
		if this.m_nodePool == nil {
			return DT_FAILURE | DT_OUT_OF_MEMORY
		}
	} else {
		this.m_nodePool.Clear()
	}

	if this.m_tinyNodePool == nil {
		this.m_tinyNodePool = DtAllocNodePool(64, 32)
		if this.m_tinyNodePool == nil {
			return DT_FAILURE | DT_OUT_OF_MEMORY
		}
	} else {
		this.m_tinyNodePool.Clear()
	}

	if this.m_openList == nil || this.m_openList.GetCapacity() < maxNodes {
		if this.m_openList != nil {
			DtFreeNodeQueue(this.m_openList)
			this.m_openList = nil
		}
		this.m_openList = DtAllocNodeQueue(maxNodes)
		if this.m_openList == nil {
			return DT_FAILURE | DT_OUT_OF_MEMORY
		}
	} else {
		this.m_openList.Clear()
	}

	return DT_SUCCESS
}

/// Returns random location on navmesh.
/// Polygons are chosen weighted by area. The search runs in linear related to number of polygon.
///  @param[in]		filter			The polygon filter to apply to the query.
///  @param[in]		frand			Function returning a random number [0..1).
///  @param[out]	randomRef		The reference id of the random location.
///  @param[out]	randomPt		The random location.
/// @returns The status flags for the query.
func (this *DtNavMeshQuery) FindRandomPoint(filter *DtQueryFilter, frand func() float32,
	randomRef *DtPolyRef, randomPt []float32) DtStatus {
	DtAssert(this.m_nav != nil)

	// Randomly pick one tile. Assume that all tiles cover roughly the same area.
	var tile *DtMeshTile
	var tsum float32
	for i := 0; i < int(this.m_nav.GetMaxTiles()); i++ {
		t := this.m_nav.GetTile(i)
		if t == nil || t.Header == nil {
			continue
		}

		// Choose random tile using reservoi sampling.
		const area float32 = 1.0 // Could be tile area too.
		tsum += area
		u := frand()
		if u*tsum <= area {
			tile = t
		}
	}
	if tile == nil {
		return DT_FAILURE
	}
	// Randomly pick one polygon weighted by polygon area.
	var poly *DtPoly
	var polyRef DtPolyRef
	base := this.m_nav.GetPolyRefBase(tile)

	var areaSum float32
	for i := 0; i < int(tile.Header.PolyCount); i++ {
		p := &tile.Polys[i]
		// Do not return off-mesh connection polygons.
		if p.GetType() != DT_POLYTYPE_GROUND {
			continue
		}
		// Must pass filter
		ref := base | (DtPolyRef)(i)
		if !filter.PassFilter(ref, tile, p) {
			continue
		}
		// Calc area of the polygon.
		var polyArea float32
		for j := 2; j < int(p.VertCount); j++ {
			va := tile.Verts[p.Verts[0]*3:]
			vb := tile.Verts[p.Verts[j-1]*3:]
			vc := tile.Verts[p.Verts[j]*3:]
			polyArea += DtTriArea2D(va, vb, vc)
		}

		// Choose random polygon weighted by area, using reservoi sampling.
		areaSum += polyArea
		u := frand()
		if u*areaSum <= polyArea {
			poly = p
			polyRef = ref
		}
	}

	if poly == nil {
		return DT_FAILURE
	}
	// Randomly pick point on polygon.
	v := tile.Verts[poly.Verts[0]*3:]
	var verts [3 * DT_VERTS_PER_POLYGON]float32
	var areas [DT_VERTS_PER_POLYGON]float32
	DtVcopy(verts[0*3:], v)
	for j := 1; j < int(poly.VertCount); j++ {
		v = tile.Verts[poly.Verts[j]*3:]
		DtVcopy(verts[j*3:], v)
	}

	s := frand()
	t := frand()

	var pt [3]float32
	DtRandomPointInConvexPoly(verts[:], int(poly.VertCount), areas[:], s, t, pt[:])

	var h float32
	status := this.GetPolyHeight(polyRef, pt[:], &h)
	if DtStatusFailed(status) {
		return status
	}
	pt[1] = h

	DtVcopy(randomPt, pt[:])
	*randomRef = polyRef

	return DT_SUCCESS
}

/// Returns random location on navmesh within the reach of specified location.
/// Polygons are chosen weighted by area. The search runs in linear related to number of polygon.
/// The location is not exactly constrained by the circle, but it limits the visited polygons.
///  @param[in]		startRef		The reference id of the polygon where the search starts.
///  @param[in]		centerPos		The center of the search circle. [(x, y, z)]
///  @param[in]		filter			The polygon filter to apply to the query.
///  @param[in]		frand			Function returning a random number [0..1).
///  @param[out]	randomRef		The reference id of the random location.
///  @param[out]	randomPt		The random location. [(x, y, z)]
/// @returns The status flags for the query.
func (this *DtNavMeshQuery) FindRandomPointAroundCircle(startRef DtPolyRef, centerPos []float32, maxRadius float32,
	filter *DtQueryFilter, frand func() float32,
	randomRef *DtPolyRef, randomPt []float32) DtStatus {
	DtAssert(this.m_nav != nil)
	DtAssert(this.m_nodePool != nil)
	DtAssert(this.m_openList != nil)

	// Validate input
	if startRef == 0 || !this.m_nav.IsValidPolyRef(startRef) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	var startTile *DtMeshTile
	var startPoly *DtPoly
	this.m_nav.GetTileAndPolyByRefUnsafe(startRef, &startTile, &startPoly)
	if !filter.PassFilter(startRef, startTile, startPoly) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	this.m_nodePool.Clear()
	this.m_openList.Clear()

	startNode := this.m_nodePool.GetNode(startRef, 0)
	DtVcopy(startNode.Pos[:], centerPos)
	startNode.Pidx = 0
	startNode.Cost = 0
	startNode.Total = 0
	startNode.Id = startRef
	startNode.Flags = DT_NODE_OPEN
	this.m_openList.Push(startNode)

	status := DT_SUCCESS

	radiusSqr := DtSqrFloat32(maxRadius)
	var areaSum float32

	var randomTile *DtMeshTile
	var randomPoly *DtPoly
	var randomPolyRef DtPolyRef

	for !this.m_openList.Empty() {
		bestNode := this.m_openList.Pop()
		bestNode.Flags &= ^DT_NODE_OPEN
		bestNode.Flags |= DT_NODE_CLOSED

		// Get poly and tile.
		// The API input has been cheked already, skip checking internal data.
		bestRef := bestNode.Id
		var bestTile *DtMeshTile
		var bestPoly *DtPoly
		this.m_nav.GetTileAndPolyByRefUnsafe(bestRef, &bestTile, &bestPoly)

		// Place random locations on on ground.
		if bestPoly.GetType() == DT_POLYTYPE_GROUND {
			// Calc area of the polygon.
			var polyArea float32
			for j := 2; j < int(bestPoly.VertCount); j++ {
				va := bestTile.Verts[bestPoly.Verts[0]*3:]
				vb := bestTile.Verts[bestPoly.Verts[j-1]*3:]
				vc := bestTile.Verts[bestPoly.Verts[j]*3:]
				polyArea += DtTriArea2D(va, vb, vc)
			}
			// Choose random polygon weighted by area, using reservoi sampling.
			areaSum += polyArea
			u := frand()
			if u*areaSum <= polyArea {
				randomTile = bestTile
				randomPoly = bestPoly
				randomPolyRef = bestRef
			}
		}

		// Get parent poly and tile.
		var parentRef DtPolyRef
		var parentTile *DtMeshTile
		var parentPoly *DtPoly
		if bestNode.Pidx != 0 {
			parentRef = this.m_nodePool.GetNodeAtIdx(bestNode.Pidx).Id
		}
		if parentRef != 0 {
			this.m_nav.GetTileAndPolyByRefUnsafe(parentRef, &parentTile, &parentPoly)
		}
		for i := bestPoly.FirstLink; i != DT_NULL_LINK; i = bestTile.Links[i].Next {
			link := &bestTile.Links[i]
			neighbourRef := link.Ref
			// Skip invalid neighbours and do not follow back to parent.
			if neighbourRef == 0 || neighbourRef == parentRef {
				continue
			}
			// Expand to neighbour
			var neighbourTile *DtMeshTile
			var neighbourPoly *DtPoly
			this.m_nav.GetTileAndPolyByRefUnsafe(neighbourRef, &neighbourTile, &neighbourPoly)

			// Do not advance if the polygon is excluded by the filter.
			if !filter.PassFilter(neighbourRef, neighbourTile, neighbourPoly) {
				continue
			}
			// Find edge and calc distance to the edge.
			var va, vb [3]float32
			if stat := this.getPortalPoints2(bestRef, bestPoly, bestTile, neighbourRef, neighbourPoly, neighbourTile, va[:], vb[:]); DtStatusFailed(stat) {
				continue
			}
			// If the circle is not touching the next polygon, skip it.
			var tseg float32
			distSqr := DtDistancePtSegSqr2D(centerPos, va[:], vb[:], &tseg)
			if distSqr > radiusSqr {
				continue
			}
			neighbourNode := this.m_nodePool.GetNode(neighbourRef, 0)
			if neighbourNode == nil {
				status |= DT_OUT_OF_NODES
				continue
			}

			if (neighbourNode.Flags & DT_NODE_CLOSED) != 0 {
				continue
			}
			// Cost
			if neighbourNode.Flags == 0 {
				DtVlerp(neighbourNode.Pos[:], va[:], vb[:], 0.5)
			}
			total := bestNode.Total + DtVdist(bestNode.Pos[:], neighbourNode.Pos[:])

			// The node is already in open list and the new result is worse, skip.
			if (neighbourNode.Flags&DT_NODE_OPEN) != 0 && total >= neighbourNode.Total {
				continue
			}
			neighbourNode.Id = neighbourRef
			neighbourNode.Flags = neighbourNode.Flags & ^DT_NODE_CLOSED
			neighbourNode.Pidx = this.m_nodePool.GetNodeIdx(bestNode)
			neighbourNode.Total = total

			if (neighbourNode.Flags & DT_NODE_OPEN) != 0 {
				this.m_openList.Modify(neighbourNode)
			} else {
				neighbourNode.Flags = DT_NODE_OPEN
				this.m_openList.Push(neighbourNode)
			}
		}
	}

	if randomPoly == nil {
		return DT_FAILURE
	}
	// Randomly pick point on polygon.
	v := randomTile.Verts[randomPoly.Verts[0]*3:]
	var verts [3 * DT_VERTS_PER_POLYGON]float32
	var areas [DT_VERTS_PER_POLYGON]float32
	DtVcopy(verts[0*3:], v[:])
	for j := 1; j < int(randomPoly.VertCount); j++ {
		v = randomTile.Verts[randomPoly.Verts[j]*3:]
		DtVcopy(verts[j*3:], v[:])
	}

	s := frand()
	t := frand()

	var pt [3]float32
	DtRandomPointInConvexPoly(verts[:], int(randomPoly.VertCount), areas[:], s, t, pt[:])

	var h float32
	stat := this.GetPolyHeight(randomPolyRef, pt[:], &h)
	if DtStatusFailed(status) {
		return stat
	}
	pt[1] = h

	DtVcopy(randomPt, pt[:])
	*randomRef = randomPolyRef

	return DT_SUCCESS
}

//////////////////////////////////////////////////////////////////////////////////////////

/// Finds the closest point on the specified polygon.
///  @param[in]		ref			The reference id of the polygon.
///  @param[in]		pos			The position to check. [(x, y, z)]
///  @param[out]	closest		The closest point on the polygon. [(x, y, z)]
///  @param[out]	posOverPoly	True of the position is over the polygon.
/// @returns The status flags for the query.
/// @par
///
/// Uses the detail polygons to find the surface height. (Most accurate.)
///
/// @p pos does not have to be within the bounds of the polygon or navigation mesh.
///
/// See closestPointOnPolyBoundary() for a limited but faster option.
///
func (this *DtNavMeshQuery) ClosestPointOnPoly(ref DtPolyRef, pos, closest []float32, posOverPoly *bool) DtStatus {
	DtAssert(this.m_nav != nil)
	var tile *DtMeshTile
	var poly *DtPoly
	if DtStatusFailed(this.m_nav.GetTileAndPolyByRef(ref, &tile, &poly)) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if tile == nil {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	// Off-mesh connections don't have detail polygons.
	if poly.GetType() == DT_POLYTYPE_OFFMESH_CONNECTION {
		v0 := tile.Verts[poly.Verts[0]*3:]
		v1 := tile.Verts[poly.Verts[1]*3:]
		d0 := DtVdist(pos, v0)
		d1 := DtVdist(pos, v1)
		u := d0 / (d0 + d1)
		DtVlerp(closest, v0, v1, u)
		if posOverPoly != nil {
			*posOverPoly = false
		}
		return DT_SUCCESS
	}

	polyBase := uintptr(unsafe.Pointer(&(tile.Polys[0])))
	current := uintptr(unsafe.Pointer(poly))
	ip := (uint32)(current-polyBase) / sizeofPoly
	pd := &tile.DetailMeshes[ip]

	// Clamp point to be inside the polygon.
	var verts [DT_VERTS_PER_POLYGON * 3]float32
	var edged [DT_VERTS_PER_POLYGON]float32
	var edget [DT_VERTS_PER_POLYGON]float32
	nv := int(poly.VertCount)
	for i := 0; i < nv; i++ {
		DtVcopy(verts[i*3:], tile.Verts[poly.Verts[i]*3:])
	}
	DtVcopy(closest, pos)
	if !DtDistancePtPolyEdgesSqr(pos, verts[:], nv, edged[:], edget[:]) {
		// Point is outside the polygon, dtClamp to nearest edge.
		dmin := edged[0]
		imin := 0
		for i := 1; i < nv; i++ {
			if edged[i] < dmin {
				dmin = edged[i]
				imin = i
			}
		}
		va := verts[imin*3:]
		vb := verts[((imin+1)%nv)*3:]
		DtVlerp(closest, va, vb, edget[imin])

		if posOverPoly != nil {
			*posOverPoly = false
		}
	} else {
		if posOverPoly != nil {
			*posOverPoly = true
		}
	}

	// Find height at the location.
	for j := 0; j < int(pd.TriCount); j++ {
		t := tile.DetailTris[(int(pd.TriBase)+j)*4:]
		var v [3][]float32
		for k := 0; k < 3; k++ {
			if t[k] < poly.VertCount {
				v[k] = tile.Verts[poly.Verts[t[k]]*3:]
			} else {
				v[k] = tile.DetailVerts[(pd.VertBase+uint32(t[k]-poly.VertCount))*3:]
			}
		}
		var h float32
		if DtClosestHeightPointTriangle(closest, v[0], v[1], v[2], &h) {
			closest[1] = h
			break
		}
	}

	return DT_SUCCESS
}

/// Returns a point on the boundary closest to the source point if the source point is outside the
/// polygon's xz-bounds.
///  @param[in]		ref			The reference id to the polygon.
///  @param[in]		pos			The position to check. [(x, y, z)]
///  @param[out]	closest		The closest point. [(x, y, z)]
/// @returns The status flags for the query.
/// @par
///
/// Much faster than closestPointOnPoly().
///
/// If the provided position lies within the polygon's xz-bounds (above or below),
/// then @p pos and @p closest will be equal.
///
/// The height of @p closest will be the polygon boundary.  The height detail is not used.
///
/// @p pos does not have to be within the bounds of the polybon or the navigation mesh.
///
func (this *DtNavMeshQuery) ClosestPointOnPolyBoundary(ref DtPolyRef, pos, closest []float32) DtStatus {
	DtAssert(this.m_nav != nil)
	var tile *DtMeshTile
	var poly *DtPoly
	if DtStatusFailed(this.m_nav.GetTileAndPolyByRef(ref, &tile, &poly)) {
		return DT_FAILURE | DT_INVALID_PARAM
	}

	// Collect vertices.
	var verts [DT_VERTS_PER_POLYGON * 3]float32
	var edged [DT_VERTS_PER_POLYGON]float32
	var edget [DT_VERTS_PER_POLYGON]float32
	nv := 0
	for i := 0; i < int(poly.VertCount); i++ {
		DtVcopy(verts[nv*3:], tile.Verts[poly.Verts[i]*3:])
		nv++
	}

	inside := DtDistancePtPolyEdgesSqr(pos, verts[:], nv, edged[:], edget[:])
	if inside {
		// Point is inside the polygon, return the point.
		DtVcopy(closest, pos)
	} else {
		// Point is outside the polygon, dtClamp to nearest edge.
		dmin := edged[0]
		imin := 0
		for i := 1; i < nv; i++ {
			if edged[i] < dmin {
				dmin = edged[i]
				imin = i
			}
		}
		va := verts[imin*3:]
		vb := verts[((imin+1)%nv)*3:]
		DtVlerp(closest, va, vb, edget[imin])
	}

	return DT_SUCCESS
}

/// Gets the height of the polygon at the provided position using the height detail. (Most accurate.)
///  @param[in]		ref			The reference id of the polygon.
///  @param[in]		pos			A position within the xz-bounds of the polygon. [(x, y, z)]
///  @param[out]	height		The height at the surface of the polygon.
/// @returns The status flags for the query.
/// @par
///
/// Will return #DT_FAILURE if the provided position is outside the xz-bounds
/// of the polygon.
///
func (this *DtNavMeshQuery) GetPolyHeight(ref DtPolyRef, pos []float32, height *float32) DtStatus {
	DtAssert(this.m_nav != nil)

	var tile *DtMeshTile
	var poly *DtPoly
	if DtStatusFailed(this.m_nav.GetTileAndPolyByRef(ref, &tile, &poly)) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if poly.GetType() == DT_POLYTYPE_OFFMESH_CONNECTION {
		v0 := tile.Verts[poly.Verts[0]*3:]
		v1 := tile.Verts[poly.Verts[1]*3:]
		d0 := DtVdist2D(pos, v0)
		d1 := DtVdist2D(pos, v1)
		u := d0 / (d0 + d1)
		if height != nil {
			*height = v0[1] + (v1[1]-v0[1])*u
		}
		return DT_SUCCESS
	} else {
		polyBase := uintptr(unsafe.Pointer(&(tile.Polys[0])))
		current := uintptr(unsafe.Pointer(poly))
		ip := (uint32)(current-polyBase) / sizeofPoly
		pd := &tile.DetailMeshes[ip]
		for j := 0; j < int(pd.TriCount); j++ {
			t := tile.DetailTris[(int(pd.TriBase)+j)*4:]
			var v [3][]float32
			for k := 0; k < 3; k++ {
				if t[k] < poly.VertCount {
					v[k] = tile.Verts[poly.Verts[t[k]]*3:]
				} else {
					v[k] = tile.DetailVerts[(pd.VertBase+uint32(t[k]-poly.VertCount))*3:]
				}
			}
			var h float32
			if DtClosestHeightPointTriangle(pos, v[0], v[1], v[2], &h) {
				if height != nil {
					*height = h
				}
				return DT_SUCCESS
			}
		}
	}

	return DT_FAILURE | DT_INVALID_PARAM
}

type dtFindNearestPolyQuery struct {
	m_query              *DtNavMeshQuery
	m_center             []float32
	m_nearestDistanceSqr float32
	m_nearestRef         DtPolyRef
	m_nearestPoint       [3]float32
}

func (this *dtFindNearestPolyQuery) constructor(query *DtNavMeshQuery, center []float32) {
	this.m_query = query
	this.m_center = center
	this.m_nearestDistanceSqr = float32(math.MaxFloat32)
}

func (this *dtFindNearestPolyQuery) nearestRef() DtPolyRef   { return this.m_nearestRef }
func (this *dtFindNearestPolyQuery) nearestPoint() []float32 { return this.m_nearestPoint[:] }

func (this *dtFindNearestPolyQuery) Process(tile *DtMeshTile, polys []*DtPoly, refs []DtPolyRef, count int) {
	//DtIgnoreUnused(polys);
	for i := 0; i < count; i++ {
		ref := refs[i]
		var closestPtPoly [3]float32
		var diff [3]float32
		posOverPoly := false
		var d float32
		this.m_query.ClosestPointOnPoly(ref, this.m_center, closestPtPoly[:], &posOverPoly)

		// If a point is directly over a polygon and closer than
		// climb height, favor that instead of straight line nearest point.
		DtVsub(diff[:], this.m_center, closestPtPoly[:])
		if posOverPoly {
			d = DtAbsFloat32(diff[1]) - tile.Header.WalkableClimb
			if d > 0 {
				d = d * d
			} else {
				d = 0
			}
		} else {
			d = DtVlenSqr(diff[:])
		}

		if d < this.m_nearestDistanceSqr {
			DtVcopy(this.m_nearestPoint[:], closestPtPoly[:])

			this.m_nearestDistanceSqr = d
			this.m_nearestRef = ref
		}
	}
}

/// Finds the polygon nearest to the specified center point.
///  @param[in]		center		The center of the search box. [(x, y, z)]
///  @param[in]		halfExtents		The search distance along each axis. [(x, y, z)]
///  @param[in]		filter		The polygon filter to apply to the query.
///  @param[out]	nearestRef	The reference id of the nearest polygon.
///  @param[out]	nearestPt	The nearest point on the polygon. [opt] [(x, y, z)]
/// @returns The status flags for the query.
/// @par
///
/// @note If the search box does not intersect any polygons the search will
/// return #DT_SUCCESS, but @p nearestRef will be zero. So if in doubt, check
/// @p nearestRef before using @p nearestPt.
///
func (this *DtNavMeshQuery) FindNearestPoly(center, halfExtents []float32,
	filter *DtQueryFilter,
	nearestRef *DtPolyRef, nearestPt []float32) DtStatus {
	DtAssert(this.m_nav != nil)

	if nearestRef == nil {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	query := dtFindNearestPolyQuery{}
	query.constructor(this, center)

	status := this.QueryPolygons2(center, halfExtents, filter, &query)
	if DtStatusFailed(status) {
		return status
	}
	*nearestRef = query.nearestRef()
	// Only override nearestPt if we actually found a poly so the nearest point
	// is valid.
	if nearestPt != nil && (*nearestRef) != 0 {
		DtVcopy(nearestPt, query.nearestPoint())
	}
	return DT_SUCCESS
}

/// Queries polygons within a tile.
func (this *DtNavMeshQuery) queryPolygonsInTile(tile *DtMeshTile, qmin, qmax []float32,
	filter *DtQueryFilter, query DtPolyQuery) {
	DtAssert(this.m_nav != nil)
	const batchSize int = 32
	var polyRefs [batchSize]DtPolyRef
	var polys [batchSize]*DtPoly
	n := 0

	if tile.BvTree != nil {
		nodeIndex := 0
		endIndex := int(tile.Header.BvNodeCount)
		tbmin := tile.Header.Bmin[:]
		tbmax := tile.Header.Bmax[:]
		qfac := tile.Header.BvQuantFactor

		// Calculate quantized box
		var bmin, bmax [3]uint16
		// dtClamp query box to world box.
		minx := DtClampFloat32(qmin[0], tbmin[0], tbmax[0]) - tbmin[0]
		miny := DtClampFloat32(qmin[1], tbmin[1], tbmax[1]) - tbmin[1]
		minz := DtClampFloat32(qmin[2], tbmin[2], tbmax[2]) - tbmin[2]
		maxx := DtClampFloat32(qmax[0], tbmin[0], tbmax[0]) - tbmin[0]
		maxy := DtClampFloat32(qmax[1], tbmin[1], tbmax[1]) - tbmin[1]
		maxz := DtClampFloat32(qmax[2], tbmin[2], tbmax[2]) - tbmin[2]
		// Quantize
		bmin[0] = (uint16)(qfac*minx) & 0xfffe
		bmin[1] = (uint16)(qfac*miny) & 0xfffe
		bmin[2] = (uint16)(qfac*minz) & 0xfffe
		bmax[0] = (uint16)(qfac*maxx+1) | 1
		bmax[1] = (uint16)(qfac*maxy+1) | 1
		bmax[2] = (uint16)(qfac*maxz+1) | 1

		// Traverse tree
		base := this.m_nav.GetPolyRefBase(tile)
		for nodeIndex < endIndex {
			node := &tile.BvTree[nodeIndex]
			overlap := DtOverlapQuantBounds(bmin[:], bmax[:], node.Bmin[:], node.Bmax[:])
			isLeafNode := (node.I >= 0)

			if isLeafNode && overlap {
				ref := base | (DtPolyRef)(node.I)
				if filter.PassFilter(ref, tile, &tile.Polys[node.I]) {
					polyRefs[n] = ref
					polys[n] = &tile.Polys[node.I]

					if n == batchSize-1 {
						query.Process(tile, polys[:], polyRefs[:], batchSize)
						n = 0
					} else {
						n++
					}
				}
			}

			if overlap || isLeafNode {
				nodeIndex++
			} else {
				escapeIndex := int(-node.I)
				nodeIndex += escapeIndex
			}
		}
	} else {
		var bmin, bmax [3]float32
		base := this.m_nav.GetPolyRefBase(tile)
		for i := 0; i < int(tile.Header.PolyCount); i++ {
			p := &tile.Polys[i]
			// Do not return off-mesh connection polygons.
			if p.GetType() == DT_POLYTYPE_OFFMESH_CONNECTION {
				continue
			}
			// Must pass filter
			ref := base | (DtPolyRef)(i)
			if !filter.PassFilter(ref, tile, p) {
				continue
			}
			// Calc polygon bounds.
			v := tile.Verts[p.Verts[0]*3:]
			DtVcopy(bmin[:], v)
			DtVcopy(bmax[:], v)
			for j := 1; j < int(p.VertCount); j++ {
				v = tile.Verts[p.Verts[j]*3:]
				DtVmin(bmin[:], v)
				DtVmax(bmax[:], v)
			}
			if DtOverlapBounds(qmin, qmax, bmin[:], bmax[:]) {
				polyRefs[n] = ref
				polys[n] = p

				if n == batchSize-1 {
					query.Process(tile, polys[:], polyRefs[:], batchSize)
					n = 0
				} else {
					n++
				}
			}
		}
	}

	// Process the last polygons that didn't make a full batch.
	if n > 0 {
		query.Process(tile, polys[:], polyRefs[:], n)
	}
}

type dtCollectPolysQuery struct {
	m_polys        []DtPolyRef
	m_maxPolys     int
	m_numCollected int
	m_overflow     bool
}

func (this *dtCollectPolysQuery) constructor(polys []DtPolyRef, maxPolys int) {
	this.m_polys = polys
	this.m_maxPolys = maxPolys
}

func (this *dtCollectPolysQuery) numCollected() int { return this.m_numCollected }
func (this *dtCollectPolysQuery) overflowed() bool  { return this.m_overflow }

func (this *dtCollectPolysQuery) Process(tile *DtMeshTile, polys []*DtPoly, refs []DtPolyRef, count int) {
	//dtIgnoreUnused(tile);
	//dtIgnoreUnused(polys);
	numLeft := this.m_maxPolys - this.m_numCollected
	toCopy := count
	if toCopy > numLeft {
		this.m_overflow = true
		toCopy = numLeft
	}
	copy(this.m_polys[this.m_numCollected:], refs[0:toCopy])
	this.m_numCollected += toCopy
}

/// Finds polygons that overlap the search box.
///  @param[in]		center		The center of the search box. [(x, y, z)]
///  @param[in]		halfExtents		The search distance along each axis. [(x, y, z)]
///  @param[in]		filter		The polygon filter to apply to the query.
///  @param[out]	polys		The reference ids of the polygons that overlap the query box.
///  @param[out]	polyCount	The number of polygons in the search result.
///  @param[in]		maxPolys	The maximum number of polygons the search result can hold.
/// @returns The status flags for the query.
/// @par
///
/// If no polygons are found, the function will return #DT_SUCCESS with a
/// @p polyCount of zero.
///
/// If @p polys is too small to hold the entire result set, then the array will
/// be filled to capacity. The method of choosing which polygons from the
/// full set are included in the partial result set is undefined.
///
func (this *DtNavMeshQuery) QueryPolygons(center, halfExtents []float32,
	filter *DtQueryFilter,
	polys []DtPolyRef, polyCount *int, maxPolys int) DtStatus {
	if polys == nil || polyCount == nil || maxPolys < 0 {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	collector := dtCollectPolysQuery{}
	collector.constructor(polys, maxPolys)

	status := this.QueryPolygons2(center, halfExtents, filter, &collector)
	if DtStatusFailed(status) {
		return status
	}
	*polyCount = collector.numCollected()
	if collector.overflowed() {
		return DT_SUCCESS | DT_BUFFER_TOO_SMALL
	} else {
		return DT_SUCCESS
	}
}

/// Finds polygons that overlap the search box.
///  @param[in]		center		The center of the search box. [(x, y, z)]
///  @param[in]		halfExtents		The search distance along each axis. [(x, y, z)]
///  @param[in]		filter		The polygon filter to apply to the query.
///  @param[in]		query		The query. Polygons found will be batched together and passed to this query.
/// @par
///
/// The query will be invoked with batches of polygons. Polygons passed
/// to the query have bounding boxes that overlap with the center and halfExtents
/// passed to this function. The dtPolyQuery::process function is invoked multiple
/// times until all overlapping polygons have been processed.
///
func (this *DtNavMeshQuery) QueryPolygons2(center, halfExtents []float32,
	filter *DtQueryFilter, query DtPolyQuery) DtStatus {
	DtAssert(this.m_nav != nil)

	if center == nil || halfExtents == nil || filter == nil || query == nil {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	var bmin, bmax [3]float32
	DtVsub(bmin[:], center, halfExtents)
	DtVadd(bmax[:], center, halfExtents)

	// Find tiles the query touches.
	var minx, miny, maxx, maxy int32
	this.m_nav.CalcTileLoc(bmin[:], &minx, &miny)
	this.m_nav.CalcTileLoc(bmax[:], &maxx, &maxy)

	const MAX_NEIS int = 32
	var neis [MAX_NEIS]*DtMeshTile

	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			nneis := this.m_nav.GetTilesAt(x, y, neis[:], MAX_NEIS)
			for j := 0; j < nneis; j++ {
				this.queryPolygonsInTile(neis[j], bmin[:], bmax[:], filter, query)
			}
		}
	}

	return DT_SUCCESS
}

/// Finds a path from the start polygon to the end polygon.
///  @param[in]		startRef	The refrence id of the start polygon.
///  @param[in]		endRef		The reference id of the end polygon.
///  @param[in]		startPos	A position within the start polygon. [(x, y, z)]
///  @param[in]		endPos		A position within the end polygon. [(x, y, z)]
///  @param[in]		filter		The polygon filter to apply to the query.
///  @param[out]	path		An ordered list of polygon references representing the path. (Start to end.)
///  							[(polyRef) * @p pathCount]
///  @param[out]	pathCount	The number of polygons returned in the @p path array.
///  @param[in]		maxPath		The maximum number of polygons the @p path array can hold. [Limit: >= 1]
/// @par
///
/// If the end polygon cannot be reached through the navigation graph,
/// the last polygon in the path will be the nearest the end polygon.
///
/// If the path array is to small to hold the full result, it will be filled as
/// far as possible from the start polygon toward the end polygon.
///
/// The start and end positions are used to calculate traversal costs.
/// (The y-values impact the result.)
///
func (this *DtNavMeshQuery) FindPath(startRef, endRef DtPolyRef,
	startPos, endPos []float32,
	filter *DtQueryFilter,
	path []DtPolyRef, pathCount *int, maxPath int) DtStatus {
	DtAssert(this.m_nav != nil)
	DtAssert(this.m_nodePool != nil)
	DtAssert(this.m_openList != nil)

	if pathCount != nil {
		*pathCount = 0
	}
	// Validate input
	if !this.m_nav.IsValidPolyRef(startRef) || !this.m_nav.IsValidPolyRef(endRef) ||
		startPos == nil || endPos == nil || filter == nil || maxPath <= 0 || path == nil || pathCount == nil {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if startRef == endRef {
		path[0] = startRef
		*pathCount = 1
		return DT_SUCCESS
	}

	this.m_nodePool.Clear()
	this.m_openList.Clear()

	startNode := this.m_nodePool.GetNode(startRef, 0)
	DtVcopy(startNode.Pos[:], startPos)
	startNode.Pidx = 0
	startNode.Cost = 0
	startNode.Total = DtVdist(startPos, endPos) * H_SCALE
	startNode.Id = startRef
	startNode.Flags = DT_NODE_OPEN
	this.m_openList.Push(startNode)

	lastBestNode := startNode
	lastBestNodeCost := startNode.Total

	outOfNodes := false

	for !this.m_openList.Empty() {
		// Remove node from open list and put it in closed list.
		bestNode := this.m_openList.Pop()
		bestNode.Flags &= ^DT_NODE_OPEN
		bestNode.Flags |= DT_NODE_CLOSED

		// Reached the goal, stop searching.
		if bestNode.Id == endRef {
			lastBestNode = bestNode
			break
		}

		// Get current poly and tile.
		// The API input has been cheked already, skip checking internal data.
		bestRef := bestNode.Id
		var bestTile *DtMeshTile
		var bestPoly *DtPoly
		this.m_nav.GetTileAndPolyByRefUnsafe(bestRef, &bestTile, &bestPoly)

		// Get parent poly and tile.
		var parentRef DtPolyRef
		var parentTile *DtMeshTile
		var parentPoly *DtPoly
		if bestNode.Pidx != 0 {
			parentRef = this.m_nodePool.GetNodeAtIdx(bestNode.Pidx).Id
		}
		if parentRef != 0 {
			this.m_nav.GetTileAndPolyByRefUnsafe(parentRef, &parentTile, &parentPoly)
		}

		for i := bestPoly.FirstLink; i != DT_NULL_LINK; i = bestTile.Links[i].Next {
			neighbourRef := bestTile.Links[i].Ref

			// Skip invalid ids and do not expand back to where we came from.
			if neighbourRef == 0 || neighbourRef == parentRef {
				continue
			}
			// Get neighbour poly and tile.
			// The API input has been cheked already, skip checking internal data.
			var neighbourTile *DtMeshTile
			var neighbourPoly *DtPoly
			this.m_nav.GetTileAndPolyByRefUnsafe(neighbourRef, &neighbourTile, &neighbourPoly)

			if !filter.PassFilter(neighbourRef, neighbourTile, neighbourPoly) {
				continue
			}
			// deal explicitly with crossing tile boundaries
			var crossSide uint8
			if bestTile.Links[i].Side != 0xff {
				crossSide = (bestTile.Links[i].Side >> 1)
			}
			// get the node
			neighbourNode := this.m_nodePool.GetNode(neighbourRef, crossSide)
			if neighbourNode == nil {
				outOfNodes = true
				continue
			}

			// If the node is visited the first time, calculate node position.
			if neighbourNode.Flags == 0 {
				this.getEdgeMidPoint2(bestRef, bestPoly, bestTile,
					neighbourRef, neighbourPoly, neighbourTile,
					neighbourNode.Pos[:])
			}

			// Calculate cost and heuristic.
			var cost float32
			var heuristic float32

			// Special case for last node.
			if neighbourRef == endRef {
				// Cost
				curCost := filter.GetCost(bestNode.Pos[:], neighbourNode.Pos[:],
					parentRef, parentTile, parentPoly,
					bestRef, bestTile, bestPoly,
					neighbourRef, neighbourTile, neighbourPoly)
				endCost := filter.GetCost(neighbourNode.Pos[:], endPos,
					bestRef, bestTile, bestPoly,
					neighbourRef, neighbourTile, neighbourPoly,
					0, nil, nil)

				cost = bestNode.Cost + curCost + endCost
				heuristic = 0
			} else {
				// Cost
				curCost := filter.GetCost(bestNode.Pos[:], neighbourNode.Pos[:],
					parentRef, parentTile, parentPoly,
					bestRef, bestTile, bestPoly,
					neighbourRef, neighbourTile, neighbourPoly)
				cost = bestNode.Cost + curCost
				heuristic = DtVdist(neighbourNode.Pos[:], endPos) * H_SCALE
			}

			total := cost + heuristic

			// The node is already in open list and the new result is worse, skip.
			if (neighbourNode.Flags&DT_NODE_OPEN) != 0 && total >= neighbourNode.Total {
				continue
			}
			// The node is already visited and process, and the new result is worse, skip.
			if (neighbourNode.Flags&DT_NODE_CLOSED) != 0 && total >= neighbourNode.Total {
				continue
			}
			// Add or update the node.
			neighbourNode.Pidx = this.m_nodePool.GetNodeIdx(bestNode)
			neighbourNode.Id = neighbourRef
			neighbourNode.Flags = (neighbourNode.Flags & ^DT_NODE_CLOSED)
			neighbourNode.Cost = cost
			neighbourNode.Total = total

			if (neighbourNode.Flags & DT_NODE_OPEN) != 0 {
				// Already in open, update node location.
				this.m_openList.Modify(neighbourNode)
			} else {
				// Put the node in open list.
				neighbourNode.Flags |= DT_NODE_OPEN
				this.m_openList.Push(neighbourNode)
			}

			// Update nearest node to target so far.
			if heuristic < lastBestNodeCost {
				lastBestNodeCost = heuristic
				lastBestNode = neighbourNode
			}
		}
	}

	status := this.getPathToNode(lastBestNode, path, pathCount, maxPath)

	if lastBestNode.Id != endRef {
		status |= DT_PARTIAL_RESULT
	}
	if outOfNodes {
		status |= DT_OUT_OF_NODES
	}
	return status
}

// Gets the path leading to the specified end node.
func (this *DtNavMeshQuery) getPathToNode(endNode *DtNode, path []DtPolyRef, pathCount *int, maxPath int) DtStatus {
	// Find the length of the entire path.
	curNode := endNode
	length := 0
	for curNode != nil {
		length++
		curNode = this.m_nodePool.GetNodeAtIdx(curNode.Pidx)
	}

	// If the path cannot be fully stored then advance to the last node we will be able to store.
	curNode = endNode
	var writeCount int
	for writeCount = length; writeCount > maxPath; writeCount-- {
		DtAssert(curNode != nil)
		curNode = this.m_nodePool.GetNodeAtIdx(curNode.Pidx)
	}

	// Write path
	for i := writeCount - 1; i >= 0; i-- {
		DtAssert(curNode != nil)
		path[i] = curNode.Id
		curNode = this.m_nodePool.GetNodeAtIdx(curNode.Pidx)
	}

	DtAssert(curNode == nil)

	*pathCount = int(DtMinInt32(int32(length), int32(maxPath)))

	if length > maxPath {
		return DT_SUCCESS | DT_BUFFER_TOO_SMALL
	}
	return DT_SUCCESS
}

/// Intializes a sliced path query.
///  @param[in]		startRef	The refrence id of the start polygon.
///  @param[in]		endRef		The reference id of the end polygon.
///  @param[in]		startPos	A position within the start polygon. [(x, y, z)]
///  @param[in]		endPos		A position within the end polygon. [(x, y, z)]
///  @param[in]		filter		The polygon filter to apply to the query.
///  @param[in]		options		query options (see: #dtFindPathOptions)
/// @returns The status flags for the query.
/// @par
///
/// @warning Calling any non-slice methods before calling finalizeSlicedFindPath()
/// or finalizeSlicedFindPathPartial() may result in corrupted data!
///
/// The @p filter pointer is stored and used for the duration of the sliced
/// path query.
///
func (this *DtNavMeshQuery) InitSlicedFindPath(startRef, endRef DtPolyRef,
	startPos, endPos []float32,
	filter *DtQueryFilter, options DtFindPathOptions) DtStatus {
	DtAssert(this.m_nav != nil)
	DtAssert(this.m_nodePool != nil)
	DtAssert(this.m_openList != nil)

	// Init path state.
	this.m_query = dtQueryData{}
	this.m_query.status = DT_FAILURE
	this.m_query.startRef = startRef
	this.m_query.endRef = endRef
	DtVcopy(this.m_query.startPos[:], startPos)
	DtVcopy(this.m_query.endPos[:], endPos)
	this.m_query.filter = filter
	this.m_query.options = options
	this.m_query.raycastLimitSqr = float32(math.MaxFloat32)

	if startRef == 0 || endRef == 0 {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	// Validate input
	if !this.m_nav.IsValidPolyRef(startRef) || !this.m_nav.IsValidPolyRef(endRef) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	// trade quality with performance?
	if (options & DT_FINDPATH_ANY_ANGLE) != 0 {
		// limiting to several times the character radius yields nice results. It is not sensitive
		// so it is enough to compute it from the first tile.
		tile := this.m_nav.GetTileByRef(DtTileRef(startRef))
		agentRadius := tile.Header.WalkableRadius
		this.m_query.raycastLimitSqr = DtSqrFloat32(agentRadius * DT_RAY_CAST_LIMIT_PROPORTIONS)
	}

	if startRef == endRef {
		this.m_query.status = DT_SUCCESS
		return DT_SUCCESS
	}

	this.m_nodePool.Clear()
	this.m_openList.Clear()

	startNode := this.m_nodePool.GetNode(startRef, 0)
	DtVcopy(startNode.Pos[:], startPos)
	startNode.Pidx = 0
	startNode.Cost = 0
	startNode.Total = DtVdist(startPos, endPos) * H_SCALE
	startNode.Id = startRef
	startNode.Flags = DT_NODE_OPEN
	this.m_openList.Push(startNode)

	this.m_query.status = DT_IN_PROGRESS
	this.m_query.lastBestNode = startNode
	this.m_query.lastBestNodeCost = startNode.Total

	return this.m_query.status
}

/// Updates an in-progress sliced path query.
///  @param[in]		maxIter		The maximum number of iterations to perform.
///  @param[out]	doneIters	The actual number of iterations completed. [opt]
/// @returns The status flags for the query.
func (this *DtNavMeshQuery) UpdateSlicedFindPath(maxIter int, doneIters *int) DtStatus {
	if !DtStatusInProgress(this.m_query.status) {
		return this.m_query.status
	}
	// Make sure the request is still valid.
	if !this.m_nav.IsValidPolyRef(this.m_query.startRef) || !this.m_nav.IsValidPolyRef(this.m_query.endRef) {
		this.m_query.status = DT_FAILURE
		return DT_FAILURE
	}

	rayHit := DtRaycastHit{}
	rayHit.MaxPath = 0

	iter := 0
	for iter < maxIter && !this.m_openList.Empty() {
		iter++

		// Remove node from open list and put it in closed list.
		bestNode := this.m_openList.Pop()
		bestNode.Flags &= ^DT_NODE_OPEN
		bestNode.Flags |= DT_NODE_CLOSED

		// Reached the goal, stop searching.
		if bestNode.Id == this.m_query.endRef {
			this.m_query.lastBestNode = bestNode
			details := this.m_query.status & DT_STATUS_DETAIL_MASK
			this.m_query.status = DT_SUCCESS | details
			if doneIters != nil {
				*doneIters = iter
			}
			return this.m_query.status
		}

		// Get current poly and tile.
		// The API input has been cheked already, skip checking internal data.
		bestRef := bestNode.Id
		var bestTile *DtMeshTile
		var bestPoly *DtPoly
		if DtStatusFailed(this.m_nav.GetTileAndPolyByRef(bestRef, &bestTile, &bestPoly)) {
			// The polygon has disappeared during the sliced query, fail.
			this.m_query.status = DT_FAILURE
			if doneIters != nil {
				*doneIters = iter
			}
			return this.m_query.status
		}

		// Get parent and grand parent poly and tile.
		var parentRef, grandpaRef DtPolyRef
		var parentTile *DtMeshTile
		var parentPoly *DtPoly
		var parentNode *DtNode
		if bestNode.Pidx != 0 {
			parentNode = this.m_nodePool.GetNodeAtIdx(bestNode.Pidx)
			parentRef = parentNode.Id
			if parentNode.Pidx != 0 {
				grandpaRef = this.m_nodePool.GetNodeAtIdx(parentNode.Pidx).Id
			}
		}
		if parentRef != 0 {
			invalidParent := DtStatusFailed(this.m_nav.GetTileAndPolyByRef(parentRef, &parentTile, &parentPoly))
			if invalidParent || (grandpaRef != 0 && !this.m_nav.IsValidPolyRef(grandpaRef)) {
				// The polygon has disappeared during the sliced query, fail.
				this.m_query.status = DT_FAILURE
				if doneIters != nil {
					*doneIters = iter
				}
				return this.m_query.status
			}
		}

		// decide whether to test raycast to previous nodes
		tryLOS := false
		if (this.m_query.options & DT_FINDPATH_ANY_ANGLE) != 0 {
			if (parentRef != 0) && (DtVdistSqr(parentNode.Pos[:], bestNode.Pos[:]) < this.m_query.raycastLimitSqr) {
				tryLOS = true
			}
		}

		for i := bestPoly.FirstLink; i != DT_NULL_LINK; i = bestTile.Links[i].Next {
			neighbourRef := bestTile.Links[i].Ref

			// Skip invalid ids and do not expand back to where we came from.
			if neighbourRef == 0 || neighbourRef == parentRef {
				continue
			}
			// Get neighbour poly and tile.
			// The API input has been cheked already, skip checking internal data.
			var neighbourTile *DtMeshTile
			var neighbourPoly *DtPoly
			this.m_nav.GetTileAndPolyByRefUnsafe(neighbourRef, &neighbourTile, &neighbourPoly)

			if !this.m_query.filter.PassFilter(neighbourRef, neighbourTile, neighbourPoly) {
				continue
			}
			// get the neighbor node
			neighbourNode := this.m_nodePool.GetNode(neighbourRef, 0)
			if neighbourNode == nil {
				this.m_query.status |= DT_OUT_OF_NODES
				continue
			}

			// do not expand to nodes that were already visited from the same parent
			if neighbourNode.Pidx != 0 && neighbourNode.Pidx == bestNode.Pidx {
				continue
			}
			// If the node is visited the first time, calculate node position.
			if neighbourNode.Flags == 0 {
				this.getEdgeMidPoint2(bestRef, bestPoly, bestTile,
					neighbourRef, neighbourPoly, neighbourTile,
					neighbourNode.Pos[:])
			}

			// Calculate cost and heuristic.
			var cost float32
			var heuristic float32

			// raycast parent
			foundShortCut := false
			rayHit.PathCost = 0
			rayHit.T = 0
			if tryLOS {
				this.Raycast2(parentRef, parentNode.Pos[:], neighbourNode.Pos[:], this.m_query.filter, DT_RAYCAST_USE_COSTS, &rayHit, grandpaRef)
				foundShortCut = (rayHit.T >= 1.0)
			}

			// update move cost
			if foundShortCut {
				// shortcut found using raycast. Using shorter cost instead
				cost = parentNode.Cost + rayHit.PathCost
			} else {
				// No shortcut found.
				curCost := this.m_query.filter.GetCost(bestNode.Pos[:], neighbourNode.Pos[:],
					parentRef, parentTile, parentPoly,
					bestRef, bestTile, bestPoly,
					neighbourRef, neighbourTile, neighbourPoly)
				cost = bestNode.Cost + curCost
			}

			// Special case for last node.
			if neighbourRef == this.m_query.endRef {
				endCost := this.m_query.filter.GetCost(neighbourNode.Pos[:], this.m_query.endPos[:],
					bestRef, bestTile, bestPoly,
					neighbourRef, neighbourTile, neighbourPoly,
					0, nil, nil)

				cost = cost + endCost
				heuristic = 0
			} else {
				heuristic = DtVdist(neighbourNode.Pos[:], this.m_query.endPos[:]) * H_SCALE
			}

			total := cost + heuristic

			// The node is already in open list and the new result is worse, skip.
			if (neighbourNode.Flags&DT_NODE_OPEN) != 0 && total >= neighbourNode.Total {
				continue
			}
			// The node is already visited and process, and the new result is worse, skip.
			if (neighbourNode.Flags&DT_NODE_CLOSED) != 0 && total >= neighbourNode.Total {
				continue
			}
			// Add or update the node.
			if foundShortCut {
				neighbourNode.Pidx = bestNode.Pidx
			} else {
				neighbourNode.Pidx = this.m_nodePool.GetNodeIdx(bestNode)
			}
			neighbourNode.Id = neighbourRef
			neighbourNode.Flags = (neighbourNode.Flags & ^(DT_NODE_CLOSED | DT_NODE_PARENT_DETACHED))
			neighbourNode.Cost = cost
			neighbourNode.Total = total
			if foundShortCut {
				neighbourNode.Flags = (neighbourNode.Flags | DT_NODE_PARENT_DETACHED)
			}
			if (neighbourNode.Flags & DT_NODE_OPEN) != 0 {
				// Already in open, update node location.
				this.m_openList.Modify(neighbourNode)
			} else {
				// Put the node in open list.
				neighbourNode.Flags |= DT_NODE_OPEN
				this.m_openList.Push(neighbourNode)
			}

			// Update nearest node to target so far.
			if heuristic < this.m_query.lastBestNodeCost {
				this.m_query.lastBestNodeCost = heuristic
				this.m_query.lastBestNode = neighbourNode
			}
		}
	}

	// Exhausted all nodes, but could not find path.
	if this.m_openList.Empty() {
		details := this.m_query.status & DT_STATUS_DETAIL_MASK
		this.m_query.status = DT_SUCCESS | details
	}

	if doneIters != nil {
		*doneIters = iter
	}

	return this.m_query.status
}

/// Finalizes and returns the results of a sliced path query.
///  @param[out]	path		An ordered list of polygon references representing the path. (Start to end.)
///  							[(polyRef) * @p pathCount]
///  @param[out]	pathCount	The number of polygons returned in the @p path array.
///  @param[in]		maxPath		The max number of polygons the path array can hold. [Limit: >= 1]
/// @returns The status flags for the query.
func (this *DtNavMeshQuery) FinalizeSlicedFindPath(path []DtPolyRef, pathCount *int, maxPath int) DtStatus {
	*pathCount = 0

	if DtStatusFailed(this.m_query.status) {
		// Reset query.
		this.m_query = dtQueryData{}
		return DT_FAILURE
	}

	n := 0

	if this.m_query.startRef == this.m_query.endRef {
		// Special case: the search starts and ends at same poly.
		path[n] = this.m_query.startRef
		n++
	} else {
		// Reverse the path.
		DtAssert(this.m_query.lastBestNode != nil)

		if this.m_query.lastBestNode.Id != this.m_query.endRef {
			this.m_query.status |= DT_PARTIAL_RESULT
		}

		var prev *DtNode
		node := this.m_query.lastBestNode
		var prevRay DtNodeFlags
		for node != nil {
			next := this.m_nodePool.GetNodeAtIdx(node.Pidx)
			node.Pidx = this.m_nodePool.GetNodeIdx(prev)
			prev = node
			nextRay := node.Flags & DT_NODE_PARENT_DETACHED                // keep track of whether parent is not adjacent (i.e. due to raycast shortcut)
			node.Flags = (node.Flags & ^DT_NODE_PARENT_DETACHED) | prevRay // and store it in the reversed path's node
			prevRay = nextRay
			node = next
		}

		// Store path
		node = prev
		for node != nil {
			next := this.m_nodePool.GetNodeAtIdx(node.Pidx)
			var status DtStatus
			if (node.Flags & DT_NODE_PARENT_DETACHED) != 0 {
				var t float32
				var normal [3]float32
				var m int
				status = this.Raycast(node.Id, node.Pos[:], next.Pos[:], this.m_query.filter, &t, normal[:], path[n:], &m, maxPath-n)
				n += m
				// raycast ends on poly boundary and the path might include the next poly boundary.
				if path[n-1] == next.Id {
					n-- // remove to avoid duplicates
				}
			} else {
				path[n] = node.Id
				n++
				if n >= maxPath {
					status = DT_BUFFER_TOO_SMALL
				}
			}

			if (status & DT_STATUS_DETAIL_MASK) != 0 {
				this.m_query.status |= status & DT_STATUS_DETAIL_MASK
				break
			}
			node = next
		}
	}

	details := this.m_query.status & DT_STATUS_DETAIL_MASK

	// Reset query.
	this.m_query = dtQueryData{}
	*pathCount = n

	return DT_SUCCESS | details
}

/// Finalizes and returns the results of an incomplete sliced path query, returning the path to the furthest
/// polygon on the existing path that was visited during the search.
///  @param[in]		existing		An array of polygon references for the existing path.
///  @param[in]		existingSize	The number of polygon in the @p existing array.
///  @param[out]	path			An ordered list of polygon references representing the path. (Start to end.)
///  								[(polyRef) * @p pathCount]
///  @param[out]	pathCount		The number of polygons returned in the @p path array.
///  @param[in]		maxPath			The max number of polygons the @p path array can hold. [Limit: >= 1]
/// @returns The status flags for the query.
func (this *DtNavMeshQuery) FinalizeSlicedFindPathPartial(existing []DtPolyRef, existingSize int,
	path []DtPolyRef, pathCount *int, maxPath int) DtStatus {
	*pathCount = 0

	if existingSize == 0 {
		return DT_FAILURE
	}

	if DtStatusFailed(this.m_query.status) {
		// Reset query.
		this.m_query = dtQueryData{}
		return DT_FAILURE
	}

	n := 0

	if this.m_query.startRef == this.m_query.endRef {
		// Special case: the search starts and ends at same poly.
		path[n] = this.m_query.startRef
		n++
	} else {
		// Find furthest existing node that was visited.
		var prev *DtNode
		var node *DtNode
		for i := existingSize - 1; i >= 0; i-- {
			var tempNode [1]*DtNode
			bFind := this.m_nodePool.FindNodes(existing[i], tempNode[:], 1)
			if bFind > 0 {
				node = tempNode[0]
				break
			}
		}

		if node == nil {
			this.m_query.status |= DT_PARTIAL_RESULT
			DtAssert(this.m_query.lastBestNode != nil)
			node = this.m_query.lastBestNode
		}

		// Reverse the path.
		var prevRay DtNodeFlags
		for node != nil {
			next := this.m_nodePool.GetNodeAtIdx(node.Pidx)
			node.Pidx = this.m_nodePool.GetNodeIdx(prev)
			prev = node
			nextRay := node.Flags & DT_NODE_PARENT_DETACHED                // keep track of whether parent is not adjacent (i.e. due to raycast shortcut)
			node.Flags = (node.Flags & ^DT_NODE_PARENT_DETACHED) | prevRay // and store it in the reversed path's node
			prevRay = nextRay
			node = next
		}

		// Store path
		node = prev
		for node != nil {
			next := this.m_nodePool.GetNodeAtIdx(node.Pidx)
			var status DtStatus
			if (node.Flags & DT_NODE_PARENT_DETACHED) != 0 {
				var t float32
				var normal [3]float32
				var m int
				status = this.Raycast(node.Id, node.Pos[:], next.Pos[:], this.m_query.filter, &t, normal[:], path[n:], &m, maxPath-n)
				n += m
				// raycast ends on poly boundary and the path might include the next poly boundary.
				if path[n-1] == next.Id {
					n-- // remove to avoid duplicates
				}
			} else {
				path[n] = node.Id
				n++
				if n >= maxPath {
					status = DT_BUFFER_TOO_SMALL
				}
			}

			if (status & DT_STATUS_DETAIL_MASK) != 0 {
				this.m_query.status |= status & DT_STATUS_DETAIL_MASK
				break
			}
			node = next
		}
	}

	details := this.m_query.status & DT_STATUS_DETAIL_MASK

	// Reset query.
	this.m_query = dtQueryData{}
	*pathCount = n

	return DT_SUCCESS | details
}

// Appends vertex to a straight path
func (this *DtNavMeshQuery) appendVertex(pos []float32, flags DtStraightPathFlags, ref DtPolyRef,
	straightPath []float32, straightPathFlags []DtStraightPathFlags, straightPathRefs []DtPolyRef,
	straightPathCount *int, maxStraightPath int) DtStatus {
	if (*straightPathCount) > 0 && DtVequal(straightPath[((*straightPathCount)-1)*3:], pos) {
		// The vertices are equal, update flags and poly.
		if straightPathFlags != nil {
			straightPathFlags[(*straightPathCount)-1] = flags
		}
		if straightPathRefs != nil {
			straightPathRefs[(*straightPathCount)-1] = ref
		}
	} else {
		// Append new vertex.
		DtVcopy(straightPath[(*straightPathCount)*3:], pos)
		if straightPathFlags != nil {
			straightPathFlags[(*straightPathCount)] = flags
		}
		if straightPathRefs != nil {
			straightPathRefs[(*straightPathCount)] = ref
		}
		(*straightPathCount)++

		// If there is no space to append more vertices, return.
		if (*straightPathCount) >= maxStraightPath {
			return DT_SUCCESS | DT_BUFFER_TOO_SMALL
		}

		// If reached end of path, return.
		if flags == DT_STRAIGHTPATH_END {
			return DT_SUCCESS
		}
	}
	return DT_IN_PROGRESS
}

// Appends intermediate portal points to a straight path.
func (this *DtNavMeshQuery) appendPortals(startIdx, endIdx int, endPos []float32, path []DtPolyRef,
	straightPath []float32, straightPathFlags []DtStraightPathFlags, straightPathRefs []DtPolyRef,
	straightPathCount *int, maxStraightPath int, options DtStraightPathOptions) DtStatus {
	startPos := straightPath[(*straightPathCount-1)*3:]
	// Append or update last vertex
	var stat DtStatus
	for i := startIdx; i < endIdx; i++ {
		// Calculate portal
		from := path[i]
		var fromTile *DtMeshTile
		var fromPoly *DtPoly
		if DtStatusFailed(this.m_nav.GetTileAndPolyByRef(from, &fromTile, &fromPoly)) {
			return DT_FAILURE | DT_INVALID_PARAM
		}
		to := path[i+1]
		var toTile *DtMeshTile
		var toPoly *DtPoly
		if DtStatusFailed(this.m_nav.GetTileAndPolyByRef(to, &toTile, &toPoly)) {
			return DT_FAILURE | DT_INVALID_PARAM
		}
		var left, right [3]float32
		if DtStatusFailed(this.getPortalPoints2(from, fromPoly, fromTile, to, toPoly, toTile, left[:], right[:])) {
			break
		}
		if (options & DT_STRAIGHTPATH_AREA_CROSSINGS) != 0 {
			// Skip intersection if only area crossings are requested.
			if fromPoly.GetArea() == toPoly.GetArea() {
				continue
			}
		}

		// Append intersection
		var s, t float32
		if DtIntersectSegSeg2D(startPos, endPos, left[:], right[:], &s, &t) {
			var pt [3]float32
			DtVlerp(pt[:], left[:], right[:], t)

			stat = this.appendVertex(pt[:], 0, path[i+1],
				straightPath, straightPathFlags, straightPathRefs,
				straightPathCount, maxStraightPath)
			if stat != DT_IN_PROGRESS {
				return stat
			}
		}
	}
	return DT_IN_PROGRESS
}

/// Finds the straight path from the start to the end position within the polygon corridor.
///  @param[in]		startPos			Path start position. [(x, y, z)]
///  @param[in]		endPos				Path end position. [(x, y, z)]
///  @param[in]		path				An array of polygon references that represent the path corridor.
///  @param[in]		pathSize			The number of polygons in the @p path array.
///  @param[out]	straightPath		Points describing the straight path. [(x, y, z) * @p straightPathCount].
///  @param[out]	straightPathFlags	Flags describing each point. (See: #dtStraightPathFlags) [opt]
///  @param[out]	straightPathRefs	The reference id of the polygon that is being entered at each point. [opt]
///  @param[out]	straightPathCount	The number of points in the straight path.
///  @param[in]		maxStraightPath		The maximum number of points the straight path arrays can hold.  [Limit: > 0]
///  @param[in]		options				Query options. (see: #dtStraightPathOptions)
/// @returns The status flags for the query.
/// @par
///
/// This method peforms what is often called 'string pulling'.
///
/// The start position is clamped to the first polygon in the path, and the
/// end position is clamped to the last. So the start and end positions should
/// normally be within or very near the first and last polygons respectively.
///
/// The returned polygon references represent the reference id of the polygon
/// that is entered at the associated path position. The reference id associated
/// with the end point will always be zero.  This allows, for example, matching
/// off-mesh link points to their representative polygons.
///
/// If the provided result buffers are too small for the entire result set,
/// they will be filled as far as possible from the start toward the end
/// position.
///
func (this *DtNavMeshQuery) FindStraightPath(startPos, endPos []float32,
	path []DtPolyRef, pathSize int,
	straightPath []float32, straightPathFlags []DtStraightPathFlags, straightPathRefs []DtPolyRef,
	straightPathCount *int, maxStraightPath int, options DtStraightPathOptions) DtStatus {
	DtAssert(this.m_nav != nil)

	*straightPathCount = 0

	if maxStraightPath == 0 {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if path[0] == 0 {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	var stat DtStatus

	// TODO: Should this be callers responsibility?
	var closestStartPos [3]float32
	if DtStatusFailed(this.ClosestPointOnPolyBoundary(path[0], startPos, closestStartPos[:])) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	var closestEndPos [3]float32
	if DtStatusFailed(this.ClosestPointOnPolyBoundary(path[pathSize-1], endPos, closestEndPos[:])) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	// Add start point.
	stat = this.appendVertex(closestStartPos[:], DT_STRAIGHTPATH_START, path[0],
		straightPath, straightPathFlags, straightPathRefs,
		straightPathCount, maxStraightPath)
	if stat != DT_IN_PROGRESS {
		return stat
	}
	if pathSize > 1 {
		var portalApex, portalLeft, portalRight [3]float32
		DtVcopy(portalApex[:], closestStartPos[:])
		DtVcopy(portalLeft[:], portalApex[:])
		DtVcopy(portalRight[:], portalApex[:])
		var apexIndex int
		var leftIndex int
		var rightIndex int

		var leftPolyType DtPolyTypes
		var rightPolyType DtPolyTypes

		leftPolyRef := path[0]
		rightPolyRef := path[0]

		for i := 0; i < pathSize; i++ {
			var left, right [3]float32
			var toType DtPolyTypes

			if i+1 < pathSize {
				var fromType DtPolyTypes // fromType is ignored.

				// Next portal.
				if DtStatusFailed(this.getPortalPoints(path[i], path[i+1], left[:], right[:], &fromType, &toType)) {
					// Failed to get portal points, in practice this means that path[i+1] is invalid polygon.
					// Clamp the end point to path[i], and return the path so far.

					if DtStatusFailed(this.ClosestPointOnPolyBoundary(path[i], endPos, closestEndPos[:])) {
						// This should only happen when the first polygon is invalid.
						return DT_FAILURE | DT_INVALID_PARAM
					}

					// Apeend portals along the current straight path segment.
					if (options & (DT_STRAIGHTPATH_AREA_CROSSINGS | DT_STRAIGHTPATH_ALL_CROSSINGS)) != 0 {
						// Ignore status return value as we're just about to return anyway.
						this.appendPortals(apexIndex, i, closestEndPos[:], path,
							straightPath, straightPathFlags, straightPathRefs,
							straightPathCount, maxStraightPath, options)
					}

					// Ignore status return value as we're just about to return anyway.
					this.appendVertex(closestEndPos[:], 0, path[i],
						straightPath, straightPathFlags, straightPathRefs,
						straightPathCount, maxStraightPath)

					if *straightPathCount >= maxStraightPath {
						return DT_SUCCESS | DT_PARTIAL_RESULT | DT_BUFFER_TOO_SMALL
					} else {
						return DT_SUCCESS | DT_PARTIAL_RESULT
					}
				}

				// If starting really close the portal, advance.
				if i == 0 {
					var t float32
					if DtDistancePtSegSqr2D(portalApex[:], left[:], right[:], &t) < DtSqrFloat32(0.001) {
						continue
					}
				}
			} else {
				// End of the path.
				DtVcopy(left[:], closestEndPos[:])
				DtVcopy(right[:], closestEndPos[:])

				toType = DT_POLYTYPE_GROUND
			}

			// Right vertex.
			if DtTriArea2D(portalApex[:], portalRight[:], right[:]) <= 0.0 {
				if DtVequal(portalApex[:], portalRight[:]) || DtTriArea2D(portalApex[:], portalLeft[:], right[:]) > 0.0 {
					DtVcopy(portalRight[:], right[:])
					if i+1 < pathSize {
						rightPolyRef = path[i+1]
					} else {
						rightPolyRef = 0
					}
					rightPolyType = toType
					rightIndex = i
				} else {
					// Append portals along the current straight path segment.
					if (options & (DT_STRAIGHTPATH_AREA_CROSSINGS | DT_STRAIGHTPATH_ALL_CROSSINGS)) != 0 {
						stat = this.appendPortals(apexIndex, leftIndex, portalLeft[:], path,
							straightPath, straightPathFlags, straightPathRefs,
							straightPathCount, maxStraightPath, options)
						if stat != DT_IN_PROGRESS {
							return stat
						}
					}

					DtVcopy(portalApex[:], portalLeft[:])
					apexIndex = leftIndex

					var flags DtStraightPathFlags
					if leftPolyRef == 0 {
						flags = DT_STRAIGHTPATH_END
					} else if leftPolyType == DT_POLYTYPE_OFFMESH_CONNECTION {
						flags = DT_STRAIGHTPATH_OFFMESH_CONNECTION
					}
					ref := leftPolyRef

					// Append or update vertex
					stat = this.appendVertex(portalApex[:], flags, ref,
						straightPath, straightPathFlags, straightPathRefs,
						straightPathCount, maxStraightPath)
					if stat != DT_IN_PROGRESS {
						return stat
					}
					DtVcopy(portalLeft[:], portalApex[:])
					DtVcopy(portalRight[:], portalApex[:])
					leftIndex = apexIndex
					rightIndex = apexIndex

					// Restart
					i = apexIndex

					continue
				}
			}

			// Left vertex.
			if DtTriArea2D(portalApex[:], portalLeft[:], left[:]) >= 0.0 {
				if DtVequal(portalApex[:], portalLeft[:]) || DtTriArea2D(portalApex[:], portalRight[:], left[:]) < 0.0 {
					DtVcopy(portalLeft[:], left[:])
					if i+1 < pathSize {
						leftPolyRef = path[i+1]
					} else {
						leftPolyRef = 0
					}
					leftPolyType = toType
					leftIndex = i
				} else {
					// Append portals along the current straight path segment.
					if (options & (DT_STRAIGHTPATH_AREA_CROSSINGS | DT_STRAIGHTPATH_ALL_CROSSINGS)) != 0 {
						stat = this.appendPortals(apexIndex, rightIndex, portalRight[:], path,
							straightPath, straightPathFlags, straightPathRefs,
							straightPathCount, maxStraightPath, options)
						if stat != DT_IN_PROGRESS {
							return stat
						}
					}

					DtVcopy(portalApex[:], portalRight[:])
					apexIndex = rightIndex

					var flags DtStraightPathFlags
					if rightPolyRef == 0 {
						flags = DT_STRAIGHTPATH_END
					} else if rightPolyType == DT_POLYTYPE_OFFMESH_CONNECTION {
						flags = DT_STRAIGHTPATH_OFFMESH_CONNECTION
					}
					ref := rightPolyRef

					// Append or update vertex
					stat = this.appendVertex(portalApex[:], flags, ref,
						straightPath, straightPathFlags, straightPathRefs,
						straightPathCount, maxStraightPath)
					if stat != DT_IN_PROGRESS {
						return stat
					}
					DtVcopy(portalLeft[:], portalApex[:])
					DtVcopy(portalRight[:], portalApex[:])
					leftIndex = apexIndex
					rightIndex = apexIndex

					// Restart
					i = apexIndex

					continue
				}
			}
		}

		// Append portals along the current straight path segment.
		if (options & (DT_STRAIGHTPATH_AREA_CROSSINGS | DT_STRAIGHTPATH_ALL_CROSSINGS)) != 0 {
			stat = this.appendPortals(apexIndex, pathSize-1, closestEndPos[:], path,
				straightPath, straightPathFlags, straightPathRefs,
				straightPathCount, maxStraightPath, options)
			if stat != DT_IN_PROGRESS {
				return stat
			}
		}
	}

	// Ignore status return value as we're just about to return anyway.
	this.appendVertex(closestEndPos[:], DT_STRAIGHTPATH_END, 0,
		straightPath, straightPathFlags, straightPathRefs,
		straightPathCount, maxStraightPath)

	if *straightPathCount >= maxStraightPath {
		return DT_SUCCESS | DT_BUFFER_TOO_SMALL
	} else {
		return DT_SUCCESS
	}
}

/// Moves from the start to the end position constrained to the navigation mesh.
///  @param[in]		startRef		The reference id of the start polygon.
///  @param[in]		startPos		A position of the mover within the start polygon. [(x, y, x)]
///  @param[in]		endPos			The desired end position of the mover. [(x, y, z)]
///  @param[in]		filter			The polygon filter to apply to the query.
///  @param[out]	resultPos		The result position of the mover. [(x, y, z)]
///  @param[out]	visited			The reference ids of the polygons visited during the move.
///  @param[out]	visitedCount	The number of polygons visited during the move.
///  @param[in]		maxVisitedSize	The maximum number of polygons the @p visited array can hold.
/// @returns The status flags for the query.
/// @par
///
/// This method is optimized for small delta movement and a small number of
/// polygons. If used for too great a distance, the result set will form an
/// incomplete path.
///
/// @p resultPos will equal the @p endPos if the end is reached.
/// Otherwise the closest reachable position will be returned.
///
/// @p resultPos is not projected onto the surface of the navigation
/// mesh. Use #getPolyHeight if this is needed.
///
/// This method treats the end position in the same manner as
/// the #raycast method. (As a 2D point.) See that method's documentation
/// for details.
///
/// If the @p visited array is too small to hold the entire result set, it will
/// be filled as far as possible from the start position toward the end
/// position.
///
func (this *DtNavMeshQuery) MoveAlongSurface(startRef DtPolyRef, startPos, endPos []float32,
	filter *DtQueryFilter,
	resultPos []float32, visited []DtPolyRef, visitedCount *int, maxVisitedSize int,
	bHit *bool) DtStatus {
	DtAssert(this.m_nav != nil)
	DtAssert(this.m_tinyNodePool != nil)

	*visitedCount = 0

	// Validate input
	if startRef == 0 {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if !this.m_nav.IsValidPolyRef(startRef) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	status := DT_SUCCESS

	const MAX_STACK int = 48
	var stack [MAX_STACK]*DtNode
	var nstack int

	this.m_tinyNodePool.Clear()

	startNode := this.m_tinyNodePool.GetNode(startRef, 0)
	startNode.Pidx = 0
	startNode.Cost = 0
	startNode.Total = 0
	startNode.Id = startRef
	startNode.Flags = DT_NODE_CLOSED
	stack[nstack] = startNode
	nstack++

	var bestPos [3]float32
	bestDist := float32(math.MaxFloat32)
	var bestNode *DtNode
	DtVcopy(bestPos[:], startPos)

	// Search constraints
	var searchPos [3]float32
	var searchRadSqr float32
	DtVlerp(searchPos[:], startPos, endPos, 0.5)
	searchRadSqr = DtSqrFloat32(DtVdist(startPos, endPos)/2.0 + 0.001)

	var verts [DT_VERTS_PER_POLYGON * 3]float32

	var wallNode *DtNode
	for nstack != 0 {
		// Pop front.
		curNode := stack[0]
		for i := 0; i < nstack-1; i++ {
			stack[i] = stack[i+1]
		}
		nstack--

		// Get poly and tile.
		// The API input has been cheked already, skip checking internal data.
		curRef := curNode.Id
		var curTile *DtMeshTile
		var curPoly *DtPoly
		this.m_nav.GetTileAndPolyByRefUnsafe(curRef, &curTile, &curPoly)

		// Collect vertices.
		nverts := int(curPoly.VertCount)
		for i := 0; i < nverts; i++ {
			DtVcopy(verts[i*3:], curTile.Verts[curPoly.Verts[i]*3:])
		}
		// If target is inside the poly, stop search.
		if DtPointInPolygon(endPos, verts[:], nverts) {
			bestNode = curNode
			DtVcopy(bestPos[:], endPos)
			break
		}

		// Find wall edges and find nearest point inside the walls.
		for i, j := 0, (int)(curPoly.VertCount-1); i < (int)(curPoly.VertCount); j, i = i, i+1 {
			// Find links to neighbours.
			const MAX_NEIS int = 8
			nneis := 0
			var neis [MAX_NEIS]DtPolyRef

			if (curPoly.Neis[j] & DT_EXT_LINK) != 0 {
				// Tile border.
				for k := curPoly.FirstLink; k != DT_NULL_LINK; k = curTile.Links[k].Next {
					link := &curTile.Links[k]
					if link.Edge == uint8(j) {
						if link.Ref != 0 {
							var neiTile *DtMeshTile
							var neiPoly *DtPoly
							this.m_nav.GetTileAndPolyByRefUnsafe(link.Ref, &neiTile, &neiPoly)
							if filter.PassFilter(link.Ref, neiTile, neiPoly) {
								if nneis < MAX_NEIS {
									neis[nneis] = link.Ref
									nneis++
								}
							}
						}
					}
				}
			} else if curPoly.Neis[j] != 0 {
				idx := (uint32)(curPoly.Neis[j] - 1)
				ref := this.m_nav.GetPolyRefBase(curTile) | DtPolyRef(idx)
				if filter.PassFilter(ref, curTile, &curTile.Polys[idx]) {
					// Internal edge, encode id.
					neis[nneis] = ref
					nneis++
				}
			}

			if nneis == 0 {
				// Wall edge, calc distance.
				vj := verts[j*3:]
				vi := verts[i*3:]
				var tseg float32
				distSqr := DtDistancePtSegSqr2D(endPos, vj, vi, &tseg)
				if distSqr < bestDist {
					// Update nearest distance.
					DtVlerp(bestPos[:], vj, vi, tseg)
					bestDist = distSqr
					bestNode = curNode
					wallNode = curNode
				}
			} else {
				for k := 0; k < nneis; k++ {
					// Skip if no node can be allocated.
					neighbourNode := this.m_tinyNodePool.GetNode(neis[k], 0)
					if neighbourNode == nil {
						continue
					}
					// Skip if already visited.
					if (neighbourNode.Flags & DT_NODE_CLOSED) != 0 {
						continue
					}
					// Skip the link if it is too far from search constraint.
					// TODO: Maybe should use getPortalPoints(), but this one is way faster.
					vj := verts[j*3:]
					vi := verts[i*3:]
					var tseg float32
					distSqr := DtDistancePtSegSqr2D(searchPos[:], vj, vi, &tseg)
					if distSqr > searchRadSqr {
						continue
					}
					// Mark as the node as visited and push to queue.
					if nstack < MAX_STACK {
						neighbourNode.Pidx = this.m_tinyNodePool.GetNodeIdx(curNode)
						neighbourNode.Flags |= DT_NODE_CLOSED
						stack[nstack] = neighbourNode
						nstack++
					}
				}
			}
		}
	}

	var n int
	if bestNode != nil {
		// Reverse the path.
		var prev *DtNode
		node := bestNode
		for node != nil {
			next := this.m_tinyNodePool.GetNodeAtIdx(node.Pidx)
			node.Pidx = this.m_tinyNodePool.GetNodeIdx(prev)
			prev = node
			node = next
		}

		// Store result
		node = prev
		for node != nil {
			visited[n] = node.Id
			n++
			if n >= maxVisitedSize {
				status |= DT_BUFFER_TOO_SMALL
				break
			}
			node = this.m_tinyNodePool.GetNodeAtIdx(node.Pidx)
		}
	}

	*bHit = (wallNode != nil && wallNode == bestNode)

	DtVcopy(resultPos, bestPos[:])

	*visitedCount = n

	return status
}

func (this *DtNavMeshQuery) getPortalPoints(from, to DtPolyRef, left, right []float32,
	fromType, toType *DtPolyTypes) DtStatus {
	DtAssert(this.m_nav != nil)

	var fromTile *DtMeshTile
	var fromPoly *DtPoly
	if DtStatusFailed(this.m_nav.GetTileAndPolyByRef(from, &fromTile, &fromPoly)) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	*fromType = fromPoly.GetType()

	var toTile *DtMeshTile
	var toPoly *DtPoly
	if DtStatusFailed(this.m_nav.GetTileAndPolyByRef(to, &toTile, &toPoly)) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	*toType = toPoly.GetType()

	return this.getPortalPoints2(from, fromPoly, fromTile, to, toPoly, toTile, left, right)
}

// Returns portal points between two polygons.
func (this *DtNavMeshQuery) getPortalPoints2(from DtPolyRef, fromPoly *DtPoly, fromTile *DtMeshTile,
	to DtPolyRef, toPoly *DtPoly, toTile *DtMeshTile,
	left, right []float32) DtStatus {
	// Find the link that points to the 'to' polygon.
	var link *DtLink
	for i := fromPoly.FirstLink; i != DT_NULL_LINK; i = fromTile.Links[i].Next {
		if fromTile.Links[i].Ref == to {
			link = &fromTile.Links[i]
			break
		}
	}
	if link == nil {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	// Handle off-mesh connections.
	if fromPoly.GetType() == DT_POLYTYPE_OFFMESH_CONNECTION {
		// Find link that points to first vertex.
		for i := fromPoly.FirstLink; i != DT_NULL_LINK; i = fromTile.Links[i].Next {
			if fromTile.Links[i].Ref == to {
				v := fromTile.Links[i].Edge
				DtVcopy(left, fromTile.Verts[fromPoly.Verts[v]*3:])
				DtVcopy(right, fromTile.Verts[fromPoly.Verts[v]*3:])
				return DT_SUCCESS
			}
		}
		return DT_FAILURE | DT_INVALID_PARAM
	}

	if toPoly.GetType() == DT_POLYTYPE_OFFMESH_CONNECTION {
		for i := toPoly.FirstLink; i != DT_NULL_LINK; i = toTile.Links[i].Next {
			if toTile.Links[i].Ref == from {
				v := toTile.Links[i].Edge
				DtVcopy(left, toTile.Verts[toPoly.Verts[v]*3:])
				DtVcopy(right, toTile.Verts[toPoly.Verts[v]*3:])
				return DT_SUCCESS
			}
		}
		return DT_FAILURE | DT_INVALID_PARAM
	}

	// Find portal vertices.
	v0 := fromPoly.Verts[link.Edge]
	v1 := fromPoly.Verts[int(link.Edge+1)%(int)(fromPoly.VertCount)]
	DtVcopy(left, fromTile.Verts[v0*3:])
	DtVcopy(right, fromTile.Verts[v1*3:])

	// If the link is at tile boundary, dtClamp the vertices to
	// the link width.
	if link.Side != 0xff {
		// Unpack portal limits.
		if link.Bmin != 0 || link.Bmax != 255 {
			s := float32(1.0 / 255.0)
			tmin := float32(link.Bmin) * s
			tmax := float32(link.Bmax) * s
			DtVlerp(left, fromTile.Verts[v0*3:], fromTile.Verts[v1*3:], tmin)
			DtVlerp(right, fromTile.Verts[v0*3:], fromTile.Verts[v1*3:], tmax)
		}
	}

	return DT_SUCCESS
}

// Returns edge mid point between two polygons.
func (this *DtNavMeshQuery) getEdgeMidPoint(from, to DtPolyRef, mid []float32) DtStatus {
	var left, right [3]float32
	var fromType, toType DtPolyTypes
	if DtStatusFailed(this.getPortalPoints(from, to, left[:], right[:], &fromType, &toType)) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	mid[0] = (left[0] + right[0]) * 0.5
	mid[1] = (left[1] + right[1]) * 0.5
	mid[2] = (left[2] + right[2]) * 0.5
	return DT_SUCCESS
}

func (this *DtNavMeshQuery) getEdgeMidPoint2(from DtPolyRef, fromPoly *DtPoly, fromTile *DtMeshTile,
	to DtPolyRef, toPoly *DtPoly, toTile *DtMeshTile,
	mid []float32) DtStatus {
	var left, right [3]float32
	if DtStatusFailed(this.getPortalPoints2(from, fromPoly, fromTile, to, toPoly, toTile, left[:], right[:])) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	mid[0] = (left[0] + right[0]) * 0.5
	mid[1] = (left[1] + right[1]) * 0.5
	mid[2] = (left[2] + right[2]) * 0.5
	return DT_SUCCESS
}

/// Casts a 'walkability' ray along the surface of the navigation mesh from
/// the start position toward the end position.
/// @note A wrapper around raycast(..., RaycastHit*). Retained for backward compatibility.
///  @param[in]		startRef	The reference id of the start polygon.
///  @param[in]		startPos	A position within the start polygon representing
///  							the start of the ray. [(x, y, z)]
///  @param[in]		endPos		The position to cast the ray toward. [(x, y, z)]
///  @param[out]	t			The hit parameter. (FLT_MAX if no wall hit.)
///  @param[out]	hitNormal	The normal of the nearest wall hit. [(x, y, z)]
///  @param[in]		filter		The polygon filter to apply to the query.
///  @param[out]	path		The reference ids of the visited polygons. [opt]
///  @param[out]	pathCount	The number of visited polygons. [opt]
///  @param[in]		maxPath		The maximum number of polygons the @p path array can hold.
/// @returns The status flags for the query.
/// @par
///
/// This method is meant to be used for quick, short distance checks.
///
/// If the path array is too small to hold the result, it will be filled as
/// far as possible from the start postion toward the end position.
///
/// <b>Using the Hit Parameter (t)</b>
///
/// If the hit parameter is a very high value (FLT_MAX), then the ray has hit
/// the end position. In this case the path represents a valid corridor to the
/// end position and the value of @p hitNormal is undefined.
///
/// If the hit parameter is zero, then the start position is on the wall that
/// was hit and the value of @p hitNormal is undefined.
///
/// If 0 < t < 1.0 then the following applies:
///
/// @code
/// distanceToHitBorder = distanceToEndPosition * t
/// hitPoint = startPos + (endPos - startPos) * t
/// @endcode
///
/// <b>Use Case Restriction</b>
///
/// The raycast ignores the y-value of the end position. (2D check.) This
/// places significant limits on how it can be used. For example:
///
/// Consider a scene where there is a main floor with a second floor balcony
/// that hangs over the main floor. So the first floor mesh extends below the
/// balcony mesh. The start position is somewhere on the first floor. The end
/// position is on the balcony.
///
/// The raycast will search toward the end position along the first floor mesh.
/// If it reaches the end position's xz-coordinates it will indicate FLT_MAX
/// (no wall hit), meaning it reached the end position. This is one example of why
/// this method is meant for short distance checks.
///
func (this *DtNavMeshQuery) Raycast(startRef DtPolyRef, startPos, endPos []float32,
	filter *DtQueryFilter,
	t *float32, hitNormal []float32, path []DtPolyRef, pathCount *int, maxPath int) DtStatus {
	var hit DtRaycastHit
	hit.Path = path
	hit.MaxPath = int32(maxPath)

	status := this.Raycast2(startRef, startPos, endPos, filter, 0, &hit, 0)

	*t = hit.T
	if hitNormal != nil {
		DtVcopy(hitNormal, hit.HitNormal[:])
	}
	if pathCount != nil {
		*pathCount = int(hit.PathCount)
	}
	return status
}

/// Casts a 'walkability' ray along the surface of the navigation mesh from
/// the start position toward the end position.
///  @param[in]		startRef	The reference id of the start polygon.
///  @param[in]		startPos	A position within the start polygon representing
///  							the start of the ray. [(x, y, z)]
///  @param[in]		endPos		The position to cast the ray toward. [(x, y, z)]
///  @param[in]		filter		The polygon filter to apply to the query.
///  @param[in]		flags		govern how the raycast behaves. See dtRaycastOptions
///  @param[out]	hit			Pointer to a raycast hit structure which will be filled by the results.
///  @param[in]		prevRef		parent of start ref. Used during for cost calculation [opt]
/// @returns The status flags for the query.
/// @par
///
/// This method is meant to be used for quick, short distance checks.
///
/// If the path array is too small to hold the result, it will be filled as
/// far as possible from the start postion toward the end position.
///
/// <b>Using the Hit Parameter t of RaycastHit</b>
///
/// If the hit parameter is a very high value (FLT_MAX), then the ray has hit
/// the end position. In this case the path represents a valid corridor to the
/// end position and the value of @p hitNormal is undefined.
///
/// If the hit parameter is zero, then the start position is on the wall that
/// was hit and the value of @p hitNormal is undefined.
///
/// If 0 < t < 1.0 then the following applies:
///
/// @code
/// distanceToHitBorder = distanceToEndPosition * t
/// hitPoint = startPos + (endPos - startPos) * t
/// @endcode
///
/// <b>Use Case Restriction</b>
///
/// The raycast ignores the y-value of the end position. (2D check.) This
/// places significant limits on how it can be used. For example:
///
/// Consider a scene where there is a main floor with a second floor balcony
/// that hangs over the main floor. So the first floor mesh extends below the
/// balcony mesh. The start position is somewhere on the first floor. The end
/// position is on the balcony.
///
/// The raycast will search toward the end position along the first floor mesh.
/// If it reaches the end position's xz-coordinates it will indicate FLT_MAX
/// (no wall hit), meaning it reached the end position. This is one example of why
/// this method is meant for short distance checks.
///
func (this *DtNavMeshQuery) Raycast2(startRef DtPolyRef, startPos, endPos []float32,
	filter *DtQueryFilter, options DtRaycastOptions,
	hit *DtRaycastHit, prevRef DtPolyRef) DtStatus {
	DtAssert(this.m_nav != nil)

	hit.T = 0
	hit.PathCount = 0
	hit.PathCost = 0

	// Validate input
	if startRef == 0 || !this.m_nav.IsValidPolyRef(startRef) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if prevRef != 0 && !this.m_nav.IsValidPolyRef(prevRef) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	var dir, curPos, lastPos [3]float32
	var verts [DT_VERTS_PER_POLYGON*3 + 3]float32
	n := 0

	DtVcopy(curPos[:], startPos)
	DtVsub(dir[:], endPos, startPos)
	DtVset(hit.HitNormal[:], 0, 0, 0)

	status := DT_SUCCESS

	var prevTile, tile, nextTile *DtMeshTile
	var prevPoly, poly, nextPoly *DtPoly
	var curRef DtPolyRef

	// The API input has been checked already, skip checking internal data.
	curRef = startRef
	tile = nil
	poly = nil
	this.m_nav.GetTileAndPolyByRefUnsafe(curRef, &tile, &poly)
	nextTile = tile
	prevTile = tile
	nextPoly = poly
	prevPoly = poly
	if prevRef != 0 {
		this.m_nav.GetTileAndPolyByRefUnsafe(prevRef, &prevTile, &prevPoly)
	}
	for curRef != 0 {
		// Cast ray against current polygon.

		// Collect vertices.
		nv := 0
		for i := 0; i < (int)(poly.VertCount); i++ {
			DtVcopy(verts[nv*3:], tile.Verts[poly.Verts[i]*3:])
			nv++
		}

		var tmin, tmax float32
		var segMin, segMax int
		if !DtIntersectSegmentPoly2D(startPos, endPos, verts[:], nv, &tmin, &tmax, &segMin, &segMax) {
			// Could not hit the polygon, keep the old t and report hit.
			hit.PathCount = int32(n)
			return status
		}

		hit.HitEdgeIndex = int32(segMax)

		// Keep track of furthest t so far.
		if tmax > hit.T {
			hit.T = tmax
		}
		// Store visited polygons.
		if n < int(hit.MaxPath) {
			hit.Path[n] = curRef
			n++
		} else {
			status |= DT_BUFFER_TOO_SMALL
		}
		// Ray end is completely inside the polygon.
		if segMax == -1 {
			hit.T = float32(math.MaxFloat32)
			hit.PathCount = int32(n)

			// add the cost
			if (options & DT_RAYCAST_USE_COSTS) != 0 {
				hit.PathCost += filter.GetCost(curPos[:], endPos, prevRef, prevTile, prevPoly, curRef, tile, poly, curRef, tile, poly)
			}
			return status
		}

		// Follow neighbours.
		var nextRef DtPolyRef

		for i := poly.FirstLink; i != DT_NULL_LINK; i = tile.Links[i].Next {
			link := &tile.Links[i]

			// Find link which contains this edge.
			if (int)(link.Edge) != segMax {
				continue
			}
			// Get pointer to the next polygon.
			nextTile = nil
			nextPoly = nil
			this.m_nav.GetTileAndPolyByRefUnsafe(link.Ref, &nextTile, &nextPoly)

			// Skip off-mesh connections.
			if nextPoly.GetType() == DT_POLYTYPE_OFFMESH_CONNECTION {
				continue
			}
			// Skip links based on filter.
			if !filter.PassFilter(link.Ref, nextTile, nextPoly) {
				continue
			}
			// If the link is internal, just return the ref.
			if link.Side == 0xff {
				nextRef = link.Ref
				break
			}

			// If the link is at tile boundary,

			// Check if the link spans the whole edge, and accept.
			if link.Bmin == 0 && link.Bmax == 255 {
				nextRef = link.Ref
				break
			}

			// Check for partial edge links.
			v0 := poly.Verts[link.Edge]
			v1 := poly.Verts[(link.Edge+1)%poly.VertCount]
			left := tile.Verts[v0*3:]
			right := tile.Verts[v1*3:]

			// Check that the intersection lies inside the link portal.
			if link.Side == 0 || link.Side == 4 {
				// Calculate link size.
				s := float32(1.0 / 255.0)
				lmin := left[2] + (right[2]-left[2])*(float32(link.Bmin)*s)
				lmax := left[2] + (right[2]-left[2])*(float32(link.Bmax)*s)
				if lmin > lmax {
					DtSwapFloat32(&lmin, &lmax)
				}

				// Find Z intersection.
				z := startPos[2] + (endPos[2]-startPos[2])*tmax
				if z >= lmin && z <= lmax {
					nextRef = link.Ref
					break
				}
			} else if link.Side == 2 || link.Side == 6 {
				// Calculate link size.
				s := float32(1.0 / 255.0)
				lmin := left[0] + (right[0]-left[0])*(float32(link.Bmin)*s)
				lmax := left[0] + (right[0]-left[0])*(float32(link.Bmax)*s)
				if lmin > lmax {
					DtSwapFloat32(&lmin, &lmax)
				}

				// Find X intersection.
				x := startPos[0] + (endPos[0]-startPos[0])*tmax
				if x >= lmin && x <= lmax {
					nextRef = link.Ref
					break
				}
			}
		}

		// add the cost
		if (options & DT_RAYCAST_USE_COSTS) != 0 {
			// compute the intersection point at the furthest end of the polygon
			// and correct the height (since the raycast moves in 2d)
			DtVcopy(lastPos[:], curPos[:])
			DtVmad(curPos[:], startPos, dir[:], hit.T)
			e1 := verts[segMax*3:]
			e2 := verts[((segMax+1)%nv)*3:]
			var eDir, diff [3]float32
			DtVsub(eDir[:], e2, e1)
			DtVsub(diff[:], curPos[:], e1)
			var s float32
			if DtSqrFloat32(eDir[0]) > DtSqrFloat32(eDir[2]) {
				s = diff[0] / eDir[0]
			} else {
				s = diff[2] / eDir[2]
			}
			curPos[1] = e1[1] + eDir[1]*s

			hit.PathCost += filter.GetCost(lastPos[:], curPos[:], prevRef, prevTile, prevPoly, curRef, tile, poly, nextRef, nextTile, nextPoly)
		}

		if nextRef == 0 {
			// No neighbour, we hit a wall.

			// Calculate hit normal.
			a := segMax
			var b int
			if segMax+1 < nv {
				b = segMax + 1
			} else {
				b = 0
			}
			va := verts[a*3:]
			vb := verts[b*3:]
			dx := vb[0] - va[0]
			dz := vb[2] - va[2]
			hit.HitNormal[0] = dz
			hit.HitNormal[1] = 0
			hit.HitNormal[2] = -dx
			DtVnormalize(hit.HitNormal[:])

			hit.PathCount = int32(n)
			return status
		}

		// No hit, advance to neighbour polygon.
		prevRef = curRef
		curRef = nextRef
		prevTile = tile
		tile = nextTile
		prevPoly = poly
		poly = nextPoly
	}

	hit.PathCount = int32(n)

	return status
}

/// Finds the polygons along the navigation graph that touch the specified circle.
///  @param[in]		startRef		The reference id of the polygon where the search starts.
///  @param[in]		centerPos		The center of the search circle. [(x, y, z)]
///  @param[in]		radius			The radius of the search circle.
///  @param[in]		filter			The polygon filter to apply to the query.
///  @param[out]	resultRef		The reference ids of the polygons touched by the circle. [opt]
///  @param[out]	resultParent	The reference ids of the parent polygons for each result.
///  								Zero if a result polygon has no parent. [opt]
///  @param[out]	resultCost		The search cost from @p centerPos to the polygon. [opt]
///  @param[out]	resultCount		The number of polygons found. [opt]
///  @param[in]		maxResult		The maximum number of polygons the result arrays can hold.
/// @returns The status flags for the query.
/// @par
///
/// At least one result array must be provided.
///
/// The order of the result set is from least to highest cost to reach the polygon.
///
/// A common use case for this method is to perform Dijkstra searches.
/// Candidate polygons are found by searching the graph beginning at the start polygon.
///
/// If a polygon is not found via the graph search, even if it intersects the
/// search circle, it will not be included in the result set. For example:
///
/// polyA is the start polygon.
/// polyB shares an edge with polyA. (Is adjacent.)
/// polyC shares an edge with polyB, but not with polyA
/// Even if the search circle overlaps polyC, it will not be included in the
/// result set unless polyB is also in the set.
///
/// The value of the center point is used as the start position for cost
/// calculations. It is not projected onto the surface of the mesh, so its
/// y-value will effect the costs.
///
/// Intersection tests occur in 2D. All polygons and the search circle are
/// projected onto the xz-plane. So the y-value of the center point does not
/// effect intersection tests.
///
/// If the result arrays are to small to hold the entire result set, they will be
/// filled to capacity.
///
func (this *DtNavMeshQuery) FindPolysAroundCircle(startRef DtPolyRef, centerPos []float32, radius float32,
	filter *DtQueryFilter,
	resultRef, resultParent []DtPolyRef, resultCost []float32,
	resultCount *int, maxResult int) DtStatus {
	DtAssert(this.m_nav != nil)
	DtAssert(this.m_nodePool != nil)
	DtAssert(this.m_openList != nil)

	*resultCount = 0

	// Validate input
	if startRef == 0 || !this.m_nav.IsValidPolyRef(startRef) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	this.m_nodePool.Clear()
	this.m_openList.Clear()

	startNode := this.m_nodePool.GetNode(startRef, 0)
	DtVcopy(startNode.Pos[:], centerPos)
	startNode.Pidx = 0
	startNode.Cost = 0
	startNode.Total = 0
	startNode.Id = startRef
	startNode.Flags = DT_NODE_OPEN
	this.m_openList.Push(startNode)

	status := DT_SUCCESS

	n := 0

	radiusSqr := DtSqrFloat32(radius)

	for !this.m_openList.Empty() {
		bestNode := this.m_openList.Pop()
		bestNode.Flags &= ^DT_NODE_OPEN
		bestNode.Flags |= DT_NODE_CLOSED

		// Get poly and tile.
		// The API input has been cheked already, skip checking internal data.
		bestRef := bestNode.Id
		var bestTile *DtMeshTile
		var bestPoly *DtPoly
		this.m_nav.GetTileAndPolyByRefUnsafe(bestRef, &bestTile, &bestPoly)

		// Get parent poly and tile.
		var parentRef DtPolyRef
		var parentTile *DtMeshTile
		var parentPoly *DtPoly
		if bestNode.Pidx != 0 {
			parentRef = this.m_nodePool.GetNodeAtIdx(bestNode.Pidx).Id
		}
		if parentRef != 0 {
			this.m_nav.GetTileAndPolyByRefUnsafe(parentRef, &parentTile, &parentPoly)
		}
		if n < maxResult {
			if resultRef != nil {
				resultRef[n] = bestRef
			}
			if resultParent != nil {
				resultParent[n] = parentRef
			}
			if resultCost != nil {
				resultCost[n] = bestNode.Total
			}
			n++
		} else {
			status |= DT_BUFFER_TOO_SMALL
		}

		for i := bestPoly.FirstLink; i != DT_NULL_LINK; i = bestTile.Links[i].Next {
			link := &bestTile.Links[i]
			neighbourRef := link.Ref
			// Skip invalid neighbours and do not follow back to parent.
			if neighbourRef == 0 || neighbourRef == parentRef {
				continue
			}
			// Expand to neighbour
			var neighbourTile *DtMeshTile
			var neighbourPoly *DtPoly
			this.m_nav.GetTileAndPolyByRefUnsafe(neighbourRef, &neighbourTile, &neighbourPoly)

			// Do not advance if the polygon is excluded by the filter.
			if !filter.PassFilter(neighbourRef, neighbourTile, neighbourPoly) {
				continue
			}
			// Find edge and calc distance to the edge.
			var va, vb [3]float32
			if stat := this.getPortalPoints2(bestRef, bestPoly, bestTile, neighbourRef, neighbourPoly, neighbourTile, va[:], vb[:]); DtStatusFailed(stat) {
				continue
			}
			// If the circle is not touching the next polygon, skip it.
			var tseg float32
			distSqr := DtDistancePtSegSqr2D(centerPos, va[:], vb[:], &tseg)
			if distSqr > radiusSqr {
				continue
			}
			neighbourNode := this.m_nodePool.GetNode(neighbourRef, 0)
			if neighbourNode == nil {
				status |= DT_OUT_OF_NODES
				continue
			}

			if (neighbourNode.Flags & DT_NODE_CLOSED) != 0 {
				continue
			}
			// Cost
			if neighbourNode.Flags == 0 {
				DtVlerp(neighbourNode.Pos[:], va[:], vb[:], 0.5)
			}
			cost := filter.GetCost(
				bestNode.Pos[:], neighbourNode.Pos[:],
				parentRef, parentTile, parentPoly,
				bestRef, bestTile, bestPoly,
				neighbourRef, neighbourTile, neighbourPoly)

			total := bestNode.Total + cost

			// The node is already in open list and the new result is worse, skip.
			if (neighbourNode.Flags&DT_NODE_OPEN) != 0 && total >= neighbourNode.Total {
				continue
			}
			neighbourNode.Id = neighbourRef
			neighbourNode.Pidx = this.m_nodePool.GetNodeIdx(bestNode)
			neighbourNode.Total = total

			if (neighbourNode.Flags & DT_NODE_OPEN) != 0 {
				this.m_openList.Modify(neighbourNode)
			} else {
				neighbourNode.Flags = DT_NODE_OPEN
				this.m_openList.Push(neighbourNode)
			}
		}
	}

	*resultCount = n

	return status
}

/// Finds the polygons along the naviation graph that touch the specified convex polygon.
///  @param[in]		startRef		The reference id of the polygon where the search starts.
///  @param[in]		verts			The vertices describing the convex polygon. (CCW)
///  								[(x, y, z) * @p nverts]
///  @param[in]		nverts			The number of vertices in the polygon.
///  @param[in]		filter			The polygon filter to apply to the query.
///  @param[out]	resultRef		The reference ids of the polygons touched by the search polygon. [opt]
///  @param[out]	resultParent	The reference ids of the parent polygons for each result. Zero if a
///  								result polygon has no parent. [opt]
///  @param[out]	resultCost		The search cost from the centroid point to the polygon. [opt]
///  @param[out]	resultCount		The number of polygons found.
///  @param[in]		maxResult		The maximum number of polygons the result arrays can hold.
/// @returns The status flags for the query.
/// @par
///
/// The order of the result set is from least to highest cost.
///
/// At least one result array must be provided.
///
/// A common use case for this method is to perform Dijkstra searches.
/// Candidate polygons are found by searching the graph beginning at the start
/// polygon.
///
/// The same intersection test restrictions that apply to findPolysAroundCircle()
/// method apply to this method.
///
/// The 3D centroid of the search polygon is used as the start position for cost
/// calculations.
///
/// Intersection tests occur in 2D. All polygons are projected onto the
/// xz-plane. So the y-values of the vertices do not effect intersection tests.
///
/// If the result arrays are is too small to hold the entire result set, they will
/// be filled to capacity.
///
func (this *DtNavMeshQuery) FindPolysAroundShape(startRef DtPolyRef, verts []float32, nverts int,
	filter *DtQueryFilter,
	resultRef, resultParent []DtPolyRef, resultCost []float32,
	resultCount *int, maxResult int) DtStatus {
	DtAssert(this.m_nav != nil)
	DtAssert(this.m_nodePool != nil)
	DtAssert(this.m_openList != nil)

	*resultCount = 0

	// Validate input
	if startRef == 0 || !this.m_nav.IsValidPolyRef(startRef) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	this.m_nodePool.Clear()
	this.m_openList.Clear()

	var centerPos [3]float32
	for i := 0; i < nverts; i++ {
		DtVadd(centerPos[:], centerPos[:], verts[i*3:])
	}
	DtVscale(centerPos[:], centerPos[:], 1.0/float32(nverts))

	startNode := this.m_nodePool.GetNode(startRef, 0)
	DtVcopy(startNode.Pos[:], centerPos[:])
	startNode.Pidx = 0
	startNode.Cost = 0
	startNode.Total = 0
	startNode.Id = startRef
	startNode.Flags = DT_NODE_OPEN
	this.m_openList.Push(startNode)

	status := DT_SUCCESS

	n := 0

	for !this.m_openList.Empty() {
		bestNode := this.m_openList.Pop()
		bestNode.Flags &= ^DT_NODE_OPEN
		bestNode.Flags |= DT_NODE_CLOSED

		// Get poly and tile.
		// The API input has been cheked already, skip checking internal data.
		bestRef := bestNode.Id
		var bestTile *DtMeshTile
		var bestPoly *DtPoly
		this.m_nav.GetTileAndPolyByRefUnsafe(bestRef, &bestTile, &bestPoly)

		// Get parent poly and tile.
		var parentRef DtPolyRef
		var parentTile *DtMeshTile
		var parentPoly *DtPoly
		if bestNode.Pidx != 0 {
			parentRef = this.m_nodePool.GetNodeAtIdx(bestNode.Pidx).Id
		}
		if parentRef != 0 {
			this.m_nav.GetTileAndPolyByRefUnsafe(parentRef, &parentTile, &parentPoly)
		}
		if n < maxResult {
			if resultRef != nil {
				resultRef[n] = bestRef
			}
			if resultParent != nil {
				resultParent[n] = parentRef
			}
			if resultCost != nil {
				resultCost[n] = bestNode.Total
			}

			n++
		} else {
			status |= DT_BUFFER_TOO_SMALL
		}

		for i := bestPoly.FirstLink; i != DT_NULL_LINK; i = bestTile.Links[i].Next {
			link := &bestTile.Links[i]
			neighbourRef := link.Ref
			// Skip invalid neighbours and do not follow back to parent.
			if neighbourRef == 0 || neighbourRef == parentRef {
				continue
			}
			// Expand to neighbour
			var neighbourTile *DtMeshTile
			var neighbourPoly *DtPoly
			this.m_nav.GetTileAndPolyByRefUnsafe(neighbourRef, &neighbourTile, &neighbourPoly)

			// Do not advance if the polygon is excluded by the filter.
			if !filter.PassFilter(neighbourRef, neighbourTile, neighbourPoly) {
				continue
			}
			// Find edge and calc distance to the edge.
			var va, vb [3]float32
			if stat := this.getPortalPoints2(bestRef, bestPoly, bestTile, neighbourRef, neighbourPoly, neighbourTile, va[:], vb[:]); DtStatusFailed(stat) {
				continue
			}
			// If the poly is not touching the edge to the next polygon, skip the connection it.
			var tmin, tmax float32
			var segMin, segMax int
			if !DtIntersectSegmentPoly2D(va[:], vb[:], verts, nverts, &tmin, &tmax, &segMin, &segMax) {
				continue
			}
			if tmin > 1.0 || tmax < 0.0 {
				continue
			}
			neighbourNode := this.m_nodePool.GetNode(neighbourRef, 0)
			if neighbourNode == nil {
				status |= DT_OUT_OF_NODES
				continue
			}

			if (neighbourNode.Flags & DT_NODE_CLOSED) != 0 {
				continue
			}
			// Cost
			if neighbourNode.Flags == 0 {
				DtVlerp(neighbourNode.Pos[:], va[:], vb[:], 0.5)
			}
			cost := filter.GetCost(
				bestNode.Pos[:], neighbourNode.Pos[:],
				parentRef, parentTile, parentPoly,
				bestRef, bestTile, bestPoly,
				neighbourRef, neighbourTile, neighbourPoly)

			total := bestNode.Total + cost

			// The node is already in open list and the new result is worse, skip.
			if (neighbourNode.Flags&DT_NODE_OPEN) != 0 && total >= neighbourNode.Total {
				continue
			}
			neighbourNode.Id = neighbourRef
			neighbourNode.Pidx = this.m_nodePool.GetNodeIdx(bestNode)
			neighbourNode.Total = total

			if (neighbourNode.Flags & DT_NODE_OPEN) != 0 {
				this.m_openList.Modify(neighbourNode)
			} else {
				neighbourNode.Flags = DT_NODE_OPEN
				this.m_openList.Push(neighbourNode)
			}
		}
	}

	*resultCount = n

	return status
}

/// Gets a path from the explored nodes in the previous search.
///  @param[in]		endRef		The reference id of the end polygon.
///  @param[out]	path		An ordered list of polygon references representing the path. (Start to end.)
///  							[(polyRef) * @p pathCount]
///  @param[out]	pathCount	The number of polygons returned in the @p path array.
///  @param[in]		maxPath		The maximum number of polygons the @p path array can hold. [Limit: >= 0]
///  @returns		The status flags. Returns DT_FAILURE | DT_INVALID_PARAM if any parameter is wrong, or if
///  				@p endRef was not explored in the previous search. Returns DT_SUCCESS | DT_BUFFER_TOO_SMALL
///  				if @p path cannot contain the entire path. In this case it is filled to capacity with a partial path.
///  				Otherwise returns DT_SUCCESS.
///  @remarks		The result of this function depends on the state of the query object. For that reason it should only
///  				be used immediately after one of the two Dijkstra searches, findPolysAroundCircle or findPolysAroundShape.
func (this *DtNavMeshQuery) GetPathFromDijkstraSearch(endRef DtPolyRef, path []DtPolyRef, pathCount *int, maxPath int) DtStatus {
	if !this.m_nav.IsValidPolyRef(endRef) || path == nil || pathCount == nil || maxPath < 0 {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	*pathCount = 0

	var endNode [1]*DtNode
	if this.m_nodePool.FindNodes(endRef, endNode[:], 1) != 1 ||
		(endNode[0].Flags&DT_NODE_CLOSED) == 0 {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	return this.getPathToNode(endNode[0], path, pathCount, maxPath)
}

/// Finds the non-overlapping navigation polygons in the local neighbourhood around the center position.
///  @param[in]		startRef		The reference id of the polygon where the search starts.
///  @param[in]		centerPos		The center of the query circle. [(x, y, z)]
///  @param[in]		radius			The radius of the query circle.
///  @param[in]		filter			The polygon filter to apply to the query.
///  @param[out]	resultRef		The reference ids of the polygons touched by the circle.
///  @param[out]	resultParent	The reference ids of the parent polygons for each result.
///  								Zero if a result polygon has no parent. [opt]
///  @param[out]	resultCount		The number of polygons found.
///  @param[in]		maxResult		The maximum number of polygons the result arrays can hold.
/// @returns The status flags for the query.
/// @par
///
/// This method is optimized for a small search radius and small number of result
/// polygons.
///
/// Candidate polygons are found by searching the navigation graph beginning at
/// the start polygon.
///
/// The same intersection test restrictions that apply to the findPolysAroundCircle
/// mehtod applies to this method.
///
/// The value of the center point is used as the start point for cost calculations.
/// It is not projected onto the surface of the mesh, so its y-value will effect
/// the costs.
///
/// Intersection tests occur in 2D. All polygons and the search circle are
/// projected onto the xz-plane. So the y-value of the center point does not
/// effect intersection tests.
///
/// If the result arrays are is too small to hold the entire result set, they will
/// be filled to capacity.
///
func (this *DtNavMeshQuery) FindLocalNeighbourhood(startRef DtPolyRef, centerPos []float32, radius float32,
	filter *DtQueryFilter,
	resultRef, resultParent []DtPolyRef,
	resultCount *int, maxResult int) DtStatus {
	DtAssert(this.m_nav != nil)
	DtAssert(this.m_tinyNodePool != nil)

	*resultCount = 0

	// Validate input
	if startRef == 0 || !this.m_nav.IsValidPolyRef(startRef) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	const MAX_STACK int = 48
	var stack [MAX_STACK]*DtNode
	nstack := 0

	this.m_tinyNodePool.Clear()

	startNode := this.m_tinyNodePool.GetNode(startRef, 0)
	startNode.Pidx = 0
	startNode.Id = startRef
	startNode.Flags = DT_NODE_CLOSED
	stack[nstack] = startNode
	nstack++

	radiusSqr := DtSqrFloat32(radius)

	var pa [DT_VERTS_PER_POLYGON * 3]float32
	var pb [DT_VERTS_PER_POLYGON * 3]float32

	status := DT_SUCCESS

	n := 0
	if n < maxResult {
		resultRef[n] = startNode.Id
		if resultParent != nil {
			resultParent[n] = 0
		}
		n++
	} else {
		status |= DT_BUFFER_TOO_SMALL
	}

	for nstack != 0 {
		// Pop front.
		curNode := stack[0]
		for i := 0; i < nstack-1; i++ {
			stack[i] = stack[i+1]
		}
		nstack--

		// Get poly and tile.
		// The API input has been cheked already, skip checking internal data.
		curRef := curNode.Id
		var curTile *DtMeshTile
		var curPoly *DtPoly
		this.m_nav.GetTileAndPolyByRefUnsafe(curRef, &curTile, &curPoly)

		for i := curPoly.FirstLink; i != DT_NULL_LINK; i = curTile.Links[i].Next {
			link := &curTile.Links[i]
			neighbourRef := link.Ref
			// Skip invalid neighbours.
			if neighbourRef == 0 {
				continue
			}
			// Skip if cannot alloca more nodes.
			neighbourNode := this.m_tinyNodePool.GetNode(neighbourRef, 0)
			if neighbourNode == nil {
				continue
			}
			// Skip visited.
			if (neighbourNode.Flags & DT_NODE_CLOSED) != 0 {
				continue
			}
			// Expand to neighbour
			var neighbourTile *DtMeshTile
			var neighbourPoly *DtPoly
			this.m_nav.GetTileAndPolyByRefUnsafe(neighbourRef, &neighbourTile, &neighbourPoly)

			// Skip off-mesh connections.
			if neighbourPoly.GetType() == DT_POLYTYPE_OFFMESH_CONNECTION {
				continue
			}
			// Do not advance if the polygon is excluded by the filter.
			if !filter.PassFilter(neighbourRef, neighbourTile, neighbourPoly) {
				continue
			}
			// Find edge and calc distance to the edge.
			var va, vb [3]float32
			if stat := this.getPortalPoints2(curRef, curPoly, curTile, neighbourRef, neighbourPoly, neighbourTile, va[:], vb[:]); DtStatusFailed(stat) {
				continue
			}
			// If the circle is not touching the next polygon, skip it.
			var tseg float32
			distSqr := DtDistancePtSegSqr2D(centerPos, va[:], vb[:], &tseg)
			if distSqr > radiusSqr {
				continue
			}
			// Mark node visited, this is done before the overlap test so that
			// we will not visit the poly again if the test fails.
			neighbourNode.Flags |= DT_NODE_CLOSED
			neighbourNode.Pidx = this.m_tinyNodePool.GetNodeIdx(curNode)

			// Check that the polygon does not collide with existing polygons.

			// Collect vertices of the neighbour poly.
			npa := int(neighbourPoly.VertCount)
			for k := 0; k < npa; k++ {
				DtVcopy(pa[k*3:], neighbourTile.Verts[neighbourPoly.Verts[k]*3:])
			}
			overlap := false
			for j := 0; j < n; j++ {
				pastRef := resultRef[j]

				// Connected polys do not overlap.
				connected := false
				for k := curPoly.FirstLink; k != DT_NULL_LINK; k = curTile.Links[k].Next {
					if curTile.Links[k].Ref == pastRef {
						connected = true
						break
					}
				}
				if connected {
					continue
				}
				// Potentially overlapping.
				var pastTile *DtMeshTile
				var pastPoly *DtPoly
				this.m_nav.GetTileAndPolyByRefUnsafe(pastRef, &pastTile, &pastPoly)

				// Get vertices and test overlap
				npb := int(pastPoly.VertCount)
				for k := 0; k < npb; k++ {
					DtVcopy(pb[k*3:], pastTile.Verts[pastPoly.Verts[k]*3:])
				}
				if DtOverlapPolyPoly2D(pa[:], npa, pb[:], npb) {
					overlap = true
					break
				}
			}
			if overlap {
				continue
			}
			// This poly is fine, store and advance to the poly.
			if n < maxResult {
				resultRef[n] = neighbourRef
				if resultParent != nil {
					resultParent[n] = curRef
				}
				n++
			} else {
				status |= DT_BUFFER_TOO_SMALL
			}

			if nstack < MAX_STACK {
				stack[nstack] = neighbourNode
				nstack++
			}
		}
	}

	*resultCount = n

	return status
}

type dtSegInterval struct {
	ref        DtPolyRef
	tmin, tmax int16
}

func insertInterval(ints []dtSegInterval, nints *int, maxInts int,
	tmin, tmax int16, ref DtPolyRef) {
	if *nints+1 > maxInts {
		return
	}
	// Find insertion point.
	idx := 0
	for idx < *nints {
		if tmax <= ints[idx].tmin {
			break
		}
		idx++
	}
	// Move current results.
	if (*nints - idx) != 0 {
		for i := *nints; i >= idx+1; i-- {
			ints[i] = ints[i-1]
		}
	}
	// Store
	ints[idx].ref = ref
	ints[idx].tmin = tmin
	ints[idx].tmax = tmax
	(*nints)++
}

/// Returns the segments for the specified polygon, optionally including portals.
///  @param[in]		ref				The reference id of the polygon.
///  @param[in]		filter			The polygon filter to apply to the query.
///  @param[out]	segmentVerts	The segments. [(ax, ay, az, bx, by, bz) * segmentCount]
///  @param[out]	segmentRefs		The reference ids of each segment's neighbor polygon.
///  								Or zero if the segment is a wall. [opt] [(parentRef) * @p segmentCount]
///  @param[out]	segmentCount	The number of segments returned.
///  @param[in]		maxSegments		The maximum number of segments the result arrays can hold.
/// @returns The status flags for the query.
/// @par
///
/// If the @p segmentRefs parameter is provided, then all polygon segments will be returned.
/// Otherwise only the wall segments are returned.
///
/// A segment that is normally a portal will be included in the result set as a
/// wall if the @p filter results in the neighbor polygon becoomming impassable.
///
/// The @p segmentVerts and @p segmentRefs buffers should normally be sized for the
/// maximum segments per polygon of the source navigation mesh.
///
func (this *DtNavMeshQuery) GetPolyWallSegments(ref DtPolyRef, filter *DtQueryFilter,
	segmentVerts []float32, segmentRefs []DtPolyRef, segmentCount *int,
	maxSegments int) DtStatus {
	DtAssert(this.m_nav != nil)

	*segmentCount = 0

	var tile *DtMeshTile
	var poly *DtPoly
	if DtStatusFailed(this.m_nav.GetTileAndPolyByRef(ref, &tile, &poly)) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	n := 0
	const MAX_INTERVAL int = 16
	var ints [MAX_INTERVAL]dtSegInterval
	var nints int

	storePortals := (segmentRefs != nil)

	status := DT_SUCCESS

	for i, j := 0, (int)(poly.VertCount-1); i < (int)(poly.VertCount); j, i = i, i+1 {
		// Skip non-solid edges.
		nints = 0
		if (poly.Neis[j] & DT_EXT_LINK) != 0 {
			// Tile border.
			for k := poly.FirstLink; k != DT_NULL_LINK; k = tile.Links[k].Next {
				link := &tile.Links[k]
				if link.Edge == uint8(j) {
					if link.Ref != 0 {
						var neiTile *DtMeshTile
						var neiPoly *DtPoly
						this.m_nav.GetTileAndPolyByRefUnsafe(link.Ref, &neiTile, &neiPoly)
						if filter.PassFilter(link.Ref, neiTile, neiPoly) {
							insertInterval(ints[:], &nints, MAX_INTERVAL, int16(link.Bmin), int16(link.Bmax), link.Ref)
						}
					}
				}
			}
		} else {
			// Internal edge
			var neiRef DtPolyRef
			if poly.Neis[j] != 0 {
				idx := (uint32)(poly.Neis[j] - 1)
				neiRef = this.m_nav.GetPolyRefBase(tile) | DtPolyRef(idx)
				if !filter.PassFilter(neiRef, tile, &tile.Polys[idx]) {
					neiRef = 0
				}
			}

			// If the edge leads to another polygon and portals are not stored, skip.
			if neiRef != 0 && !storePortals {
				continue
			}
			if n < maxSegments {
				vj := tile.Verts[poly.Verts[j]*3:]
				vi := tile.Verts[poly.Verts[i]*3:]
				seg := segmentVerts[n*6:]
				DtVcopy(seg[0:], vj)
				DtVcopy(seg[3:], vi)
				if segmentRefs != nil {
					segmentRefs[n] = neiRef
				}
				n++
			} else {
				status |= DT_BUFFER_TOO_SMALL
			}

			continue
		}

		// Add sentinels
		insertInterval(ints[:], &nints, MAX_INTERVAL, -1, 0, 0)
		insertInterval(ints[:], &nints, MAX_INTERVAL, 255, 256, 0)

		// Store segments.
		vj := tile.Verts[poly.Verts[j]*3:]
		vi := tile.Verts[poly.Verts[i]*3:]
		for k := 1; k < nints; k++ {
			// Portal segment.
			if storePortals && ints[k].ref != 0 {
				tmin := ints[k].tmin / 255.0
				tmax := ints[k].tmax / 255.0
				if n < maxSegments {
					seg := segmentVerts[n*6:]
					DtVlerp(seg[0:], vj, vi, float32(tmin))
					DtVlerp(seg[3:], vj, vi, float32(tmax))
					if segmentRefs != nil {
						segmentRefs[n] = ints[k].ref
					}
					n++
				} else {
					status |= DT_BUFFER_TOO_SMALL
				}
			}

			// Wall segment.
			imin := ints[k-1].tmax
			imax := ints[k].tmin
			if imin != imax {
				tmin := imin / 255.0
				tmax := imax / 255.0
				if n < maxSegments {
					seg := segmentVerts[n*6:]
					DtVlerp(seg[0:], vj, vi, float32(tmin))
					DtVlerp(seg[3:], vj, vi, float32(tmax))
					if segmentRefs != nil {
						segmentRefs[n] = 0
					}
					n++
				} else {
					status |= DT_BUFFER_TOO_SMALL
				}
			}
		}
	}

	*segmentCount = n

	return status
}

/// Finds the distance from the specified position to the nearest polygon wall.
///  @param[in]		startRef		The reference id of the polygon containing @p centerPos.
///  @param[in]		centerPos		The center of the search circle. [(x, y, z)]
///  @param[in]		maxRadius		The radius of the search circle.
///  @param[in]		filter			The polygon filter to apply to the query.
///  @param[out]	hitDist			The distance to the nearest wall from @p centerPos.
///  @param[out]	hitPos			The nearest position on the wall that was hit. [(x, y, z)]
///  @param[out]	hitNormal		The normalized ray formed from the wall point to the
///  								source point. [(x, y, z)]
/// @returns The status flags for the query.
/// @par
///
/// @p hitPos is not adjusted using the height detail data.
///
/// @p hitDist will equal the search radius if there is no wall within the
/// radius. In this case the values of @p hitPos and @p hitNormal are
/// undefined.
///
/// The normal will become unpredicable if @p hitDist is a very small number.
///
func (this *DtNavMeshQuery) FindDistanceToWall(startRef DtPolyRef, centerPos []float32, maxRadius float32,
	filter *DtQueryFilter,
	hitDist *float32, hitPos []float32, hitNormal []float32) DtStatus {
	DtAssert(this.m_nav != nil)
	DtAssert(this.m_nodePool != nil)
	DtAssert(this.m_openList != nil)

	// Validate input
	if startRef == 0 || !this.m_nav.IsValidPolyRef(startRef) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	this.m_nodePool.Clear()
	this.m_openList.Clear()

	startNode := this.m_nodePool.GetNode(startRef, 0)
	DtVcopy(startNode.Pos[:], centerPos)
	startNode.Pidx = 0
	startNode.Cost = 0
	startNode.Total = 0
	startNode.Id = startRef
	startNode.Flags = DT_NODE_OPEN
	this.m_openList.Push(startNode)

	radiusSqr := DtSqrFloat32(maxRadius)

	status := DT_SUCCESS

	for !this.m_openList.Empty() {
		bestNode := this.m_openList.Pop()
		bestNode.Flags &= ^DT_NODE_OPEN
		bestNode.Flags |= DT_NODE_CLOSED

		// Get poly and tile.
		// The API input has been cheked already, skip checking internal data.
		bestRef := bestNode.Id
		var bestTile *DtMeshTile
		var bestPoly *DtPoly
		this.m_nav.GetTileAndPolyByRefUnsafe(bestRef, &bestTile, &bestPoly)

		// Get parent poly and tile.
		var parentRef DtPolyRef
		var parentTile *DtMeshTile
		var parentPoly *DtPoly
		if bestNode.Pidx != 0 {
			parentRef = this.m_nodePool.GetNodeAtIdx(bestNode.Pidx).Id
		}
		if parentRef != 0 {
			this.m_nav.GetTileAndPolyByRefUnsafe(parentRef, &parentTile, &parentPoly)
		}
		// Hit test walls.
		for i, j := 0, (int)(bestPoly.VertCount-1); i < (int)(bestPoly.VertCount); j, i = i, i+1 {
			// Skip non-solid edges.
			if (bestPoly.Neis[j] & DT_EXT_LINK) != 0 {
				// Tile border.
				solid := true
				for k := bestPoly.FirstLink; k != DT_NULL_LINK; k = bestTile.Links[k].Next {
					link := &bestTile.Links[k]
					if link.Edge == uint8(j) {
						if link.Ref != 0 {
							var neiTile *DtMeshTile
							var neiPoly *DtPoly
							this.m_nav.GetTileAndPolyByRefUnsafe(link.Ref, &neiTile, &neiPoly)
							if filter.PassFilter(link.Ref, neiTile, neiPoly) {
								solid = false
							}
						}
						break
					}
				}
				if !solid {
					continue
				}
			} else if bestPoly.Neis[j] != 0 {
				// Internal edge
				idx := (uint32)(bestPoly.Neis[j] - 1)
				ref := this.m_nav.GetPolyRefBase(bestTile) | DtPolyRef(idx)
				if filter.PassFilter(ref, bestTile, &bestTile.Polys[idx]) {
					continue
				}
			}

			// Calc distance to the edge.
			vj := bestTile.Verts[bestPoly.Verts[j]*3:]
			vi := bestTile.Verts[bestPoly.Verts[i]*3:]
			var tseg float32
			distSqr := DtDistancePtSegSqr2D(centerPos, vj, vi, &tseg)

			// Edge is too far, skip.
			if distSqr > radiusSqr {
				continue
			}
			// Hit wall, update radius.
			radiusSqr = distSqr
			// Calculate hit pos.
			hitPos[0] = vj[0] + (vi[0]-vj[0])*tseg
			hitPos[1] = vj[1] + (vi[1]-vj[1])*tseg
			hitPos[2] = vj[2] + (vi[2]-vj[2])*tseg
		}

		for i := bestPoly.FirstLink; i != DT_NULL_LINK; i = bestTile.Links[i].Next {
			link := &bestTile.Links[i]
			neighbourRef := link.Ref
			// Skip invalid neighbours and do not follow back to parent.
			if neighbourRef == 0 || neighbourRef == parentRef {
				continue
			}
			// Expand to neighbour.
			var neighbourTile *DtMeshTile
			var neighbourPoly *DtPoly
			this.m_nav.GetTileAndPolyByRefUnsafe(neighbourRef, &neighbourTile, &neighbourPoly)

			// Skip off-mesh connections.
			if neighbourPoly.GetType() == DT_POLYTYPE_OFFMESH_CONNECTION {
				continue
			}
			// Calc distance to the edge.
			va := bestTile.Verts[bestPoly.Verts[link.Edge]*3:]
			vb := bestTile.Verts[bestPoly.Verts[(link.Edge+1)%bestPoly.VertCount]*3:]
			var tseg float32
			distSqr := DtDistancePtSegSqr2D(centerPos, va, vb, &tseg)

			// If the circle is not touching the next polygon, skip it.
			if distSqr > radiusSqr {
				continue
			}
			if !filter.PassFilter(neighbourRef, neighbourTile, neighbourPoly) {
				continue
			}
			neighbourNode := this.m_nodePool.GetNode(neighbourRef, 0)
			if neighbourNode == nil {
				status |= DT_OUT_OF_NODES
				continue
			}

			if (neighbourNode.Flags & DT_NODE_CLOSED) != 0 {
				continue
			}

			// Cost
			if neighbourNode.Flags == 0 {
				this.getEdgeMidPoint2(bestRef, bestPoly, bestTile,
					neighbourRef, neighbourPoly, neighbourTile, neighbourNode.Pos[:])
			}

			total := bestNode.Total + DtVdist(bestNode.Pos[:], neighbourNode.Pos[:])

			// The node is already in open list and the new result is worse, skip.
			if (neighbourNode.Flags&DT_NODE_OPEN) != 0 && total >= neighbourNode.Total {
				continue
			}
			neighbourNode.Id = neighbourRef
			neighbourNode.Flags = (neighbourNode.Flags & ^DT_NODE_CLOSED)
			neighbourNode.Pidx = this.m_nodePool.GetNodeIdx(bestNode)
			neighbourNode.Total = total

			if (neighbourNode.Flags & DT_NODE_OPEN) != 0 {
				this.m_openList.Modify(neighbourNode)
			} else {
				neighbourNode.Flags |= DT_NODE_OPEN
				this.m_openList.Push(neighbourNode)
			}
		}
	}

	// Calc hit normal.
	DtVsub(hitNormal, centerPos, hitPos)
	DtVnormalize(hitNormal)

	*hitDist = DtMathSqrtf(radiusSqr)

	return status
}

/// Returns true if the polygon reference is valid and passes the filter restrictions.
///  @param[in]		ref			The polygon reference to check.
///  @param[in]		filter		The filter to apply.
func (this *DtNavMeshQuery) IsValidPolyRef(ref DtPolyRef, filter *DtQueryFilter) bool {
	var tile *DtMeshTile
	var poly *DtPoly
	status := this.m_nav.GetTileAndPolyByRef(ref, &tile, &poly)
	// If cannot get polygon, assume it does not exists and boundary is invalid.
	if DtStatusFailed(status) {
		return false
	}
	// If cannot pass filter, assume flags has changed and boundary is invalid.
	if !filter.PassFilter(ref, tile, poly) {
		return false
	}
	return true
}

/// Returns true if the polygon reference is in the closed list.
///  @param[in]		ref		The reference id of the polygon to check.
/// @returns True if the polygon is in closed list.
/// @par
///
/// The closed list is the list of polygons that were fully evaluated during
/// the last navigation graph search. (A* or Dijkstra)
///
func (this *DtNavMeshQuery) IsInClosedList(ref DtPolyRef) bool {
	if this.m_nodePool == nil {
		return false
	}

	var nodes [DT_MAX_STATES_PER_NODE]*DtNode
	n := this.m_nodePool.FindNodes(ref, nodes[:], uint32(DT_MAX_STATES_PER_NODE))

	for i := 0; i < int(n); i++ {
		if (nodes[i].Flags & DT_NODE_CLOSED) != 0 {
			return true
		}
	}

	return false
}