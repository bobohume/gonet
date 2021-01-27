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
	"reflect"
	"unsafe"
)

var(
	sizeofMeshTile uint32 = (uint32)(unsafe.Sizeof(DtMeshTile{}))
	sizeofPoly uint32 = (uint32)(unsafe.Sizeof(DtPoly{}))
)

///--------------------------------------------------------------  ///
///	DtPoly
///---------------------------------------------------------------  ///
/// Sets the user defined area id. [Limit: < #DT_MAX_AREAS]
func (this *DtPoly) SetArea(a uint8) {
	this.AreaAndtype = (this.AreaAndtype & 0xc0) | (a & 0x3f)
}

/// Sets the polygon type. (See: #dtPolyTypes.)
func (this *DtPoly) SetType(t DtPolyTypes) {
	this.AreaAndtype = (this.AreaAndtype & 0x3f) | (uint8(t) << 6)
}

/// Gets the user defined area id.
func (this *DtPoly) GetArea() uint8 {
	return this.AreaAndtype & 0x3f
}

/// Gets the polygon type. (See: #dtPolyTypes)
func (this *DtPoly) GetType() DtPolyTypes {
	return DtPolyTypes(this.AreaAndtype >> 6)
}

///--------------------------------------------------------------  ///
///	DtNavMesh
///---------------------------------------------------------------  ///
/// @{
/// @name Encoding and Decoding
/// These functions are generally meant for internal use only.

/// Derives a standard polygon reference.
///  @note This function is generally meant for internal use only.
///  @param[in]	salt	The tile's salt value.
///  @param[in]	it		The index of the tile.
///  @param[in]	ip		The index of the polygon within the tile.
func (this *DtNavMesh) EncodePolyId(salt, it, ip uint32) DtPolyRef {
	return DtPolyRef((salt << (this.m_polyBits + this.m_tileBits)) | (it << this.m_polyBits) | ip)
}

/// Decodes a standard polygon reference.
///  @note This function is generally meant for internal use only.
///  @param[in]	ref   The polygon reference to decode.
///  @param[out]	salt	The tile's salt value.
///  @param[out]	it		The index of the tile.
///  @param[out]	ip		The index of the polygon within the tile.
///  @see #encodePolyId
func (this *DtNavMesh) DecodePolyId(ref DtPolyRef, salt, it, ip *uint32) {
	saltMask := (uint32(1) << this.m_saltBits) - 1
	tileMask := (uint32(1) << this.m_tileBits) - 1
	polyMask := (uint32(1) << this.m_polyBits) - 1
	*salt = ((uint32(ref) >> (this.m_polyBits + this.m_tileBits)) & saltMask)
	*it = ((uint32(ref) >> this.m_polyBits) & tileMask)
	*ip = (uint32(ref) & polyMask)
}

/// Extracts a tile's salt value from the specified polygon reference.
///  @note This function is generally meant for internal use only.
///  @param[in]	ref		The polygon reference.
///  @see #encodePolyId
func (this *DtNavMesh) DecodePolyIdSalt(ref DtPolyRef) uint32 {
	saltMask := (uint32(1) << this.m_saltBits) - 1
	return ((uint32(ref) >> (this.m_polyBits + this.m_tileBits)) & saltMask)
}

/// Extracts the tile's index from the specified polygon reference.
///  @note This function is generally meant for internal use only.
///  @param[in]	ref		The polygon reference.
///  @see #encodePolyId
func (this *DtNavMesh) DecodePolyIdTile(ref DtPolyRef) uint32 {
	tileMask := (uint32(1) << this.m_tileBits) - 1
	return ((uint32(ref) >> this.m_polyBits) & tileMask)
}

/// Extracts the polygon's index (within its tile) from the specified polygon reference.
///  @note This function is generally meant for internal use only.
///  @param[in]	ref		The polygon reference.
///  @see #encodePolyId
func (this *DtNavMesh) DecodePolyIdPoly(ref DtPolyRef) uint32 {
	polyMask := (uint32(1) << this.m_polyBits) - 1
	return (uint32(ref) & polyMask)
}

/**
@class dtNavMesh

The navigation mesh consists of one or more tiles defining three primary types of structural data:

A polygon mesh which defines most of the navigation graph. (See rcPolyMesh for its structure.)
A detail mesh used for determining surface height on the polygon mesh. (See rcPolyMeshDetail for its structure.)
Off-mesh connections, which define custom point-to-point edges within the navigation graph.

The general build process is as follows:

-# Create rcPolyMesh and rcPolyMeshDetail data using the Recast build pipeline.
-# Optionally, create off-mesh connection data.
-# Combine the source data into a dtNavMeshCreateParams structure.
-# Create a tile data array using dtCreateNavMeshData().
-# Allocate at dtNavMesh object and initialize it. (For single tile navigation meshes,
   the tile data is loaded during this step.)
-# For multi-tile navigation meshes, load the tile data using dtNavMesh::addTile().

Notes:

- This class is usually used in conjunction with the dtNavMeshQuery class for pathfinding.
- Technically, all navigation meshes are tiled. A 'solo' mesh is simply a navigation mesh initialized
  to have only a single tile.
- This class does not implement any asynchronous methods. So the ::dtStatus result of all methods will
  always contain either a success or failure flag.

@see dtNavMeshQuery, dtCreateNavMeshData, dtNavMeshCreateParams, #dtAllocNavMesh, #dtFreeNavMesh
*/
func (this *DtNavMesh) constructor() {

}

func (this *DtNavMesh) destructor() {
	for i := 0; i < int(this.m_maxTiles); i++ {
		if (this.m_tiles[i].Flags & DT_TILE_FREE_DATA) != 0 {
			this.m_tiles[i].Data = nil
			this.m_tiles[i].DataSize = 0
		}
	}
	this.m_posLookup = nil
	this.m_tiles = nil
}

/// @{
/// @name Initialization and Tile Management

/// Initializes the navigation mesh for tiled use.
///  @param[in]	params		Initialization parameters.
/// @return The status flags for the operation.
func (this *DtNavMesh) Init(params *DtNavMeshParams) DtStatus {
	this.m_params = *params
	DtVcopy(this.m_orig[:], params.Orig[:])
	this.m_tileWidth = params.TileWidth
	this.m_tileHeight = params.TileHeight

	// Init tiles
	this.m_maxTiles = int32(params.MaxTiles)
	this.m_tileLutSize = int32(DtNextPow2(params.MaxTiles / 4))
	if this.m_tileLutSize == 0 {
		this.m_tileLutSize = 1
	}
	this.m_tileLutMask = this.m_tileLutSize - 1

	this.m_tiles = make([]DtMeshTile, this.m_maxTiles)
	if this.m_tiles == nil {
		return DT_FAILURE | DT_OUT_OF_MEMORY
	}
	this.m_posLookup = make([]*DtMeshTile, this.m_tileLutSize)
	if this.m_posLookup == nil {
		return DT_FAILURE | DT_OUT_OF_MEMORY
	}

	this.m_nextFree = nil
	for i := int(this.m_maxTiles - 1); i >= 0; i-- {
		this.m_tiles[i].Salt = 1
		this.m_tiles[i].Next = this.m_nextFree
		this.m_nextFree = &this.m_tiles[i]
	}

	// Init ID generator values.
	this.m_tileBits = DtIlog2(DtNextPow2(params.MaxTiles))
	this.m_polyBits = DtIlog2(DtNextPow2(params.MaxPolys))
	// Only allow 31 salt bits, since the salt mask is calculated using 32bit uint and it will overflow.
	this.m_saltBits = DtMinUInt32(31, 32-this.m_tileBits-this.m_polyBits)

	if this.m_saltBits < 10 {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	return DT_SUCCESS
}

/// Initializes the navigation mesh for single tile use.
///  @param[in]	data		Data of the new tile. (See: #dtCreateNavMeshData)
///  @param[in]	dataSize	The data size of the new tile.
///  @param[in]	flags		The tile flags. (See: #dtTileFlags)
/// @return The status flags for the operation.
///  @see dtCreateNavMeshData
func (this *DtNavMesh) Init2(data []byte, dataSize int, flags DtTileFlags) DtStatus {
	// Make sure the data is in right format.
	if dataSize < DtAlign4(int(unsafe.Sizeof(DtMeshHeader{}))) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	header := (*DtMeshHeader)(unsafe.Pointer(&(data[0])))
	if header.Magic != DT_NAVMESH_MAGIC {
		return DT_FAILURE | DT_WRONG_MAGIC
	}
	if header.Version != DT_NAVMESH_VERSION {
		return DT_FAILURE | DT_WRONG_VERSION
	}
	var params DtNavMeshParams
	DtVcopy(params.Orig[:], header.Bmin[:])
	params.TileWidth = header.Bmax[0] - header.Bmin[0]
	params.TileHeight = header.Bmax[2] - header.Bmin[2]
	params.MaxTiles = 1
	params.MaxPolys = uint32(header.PolyCount)

	status := this.Init(&params)
	if DtStatusFailed(status) {
		return status
	}
	return this.AddTile(data, dataSize, flags, 0, nil)
}

/// The navigation mesh initialization params.
/// @par
///
/// @note The parameters are created automatically when the single tile
/// initialization is performed.
func (this *DtNavMesh) GetParams() *DtNavMeshParams {
	return &this.m_params
}

//////////////////////////////////////////////////////////////////////////////////////////

/// Returns all polygons in neighbour tile based on portal defined by the segment.
func (this *DtNavMesh) findConnectingPolys(va, vb []float32, tile *DtMeshTile, side int, con []DtPolyRef, conarea []float32, maxcon int) int {
	if tile == nil {
		return 0
	}

	var amin, amax [2]float32
	calcSlabEndPoints(va, vb, amin[:], amax[:], side)
	apos := getSlabCoord(va, side)

	// Remove links pointing to 'side' and compact the links array.
	var bmin, bmax [2]float32
	m := DT_EXT_LINK | (uint16)(side)
	n := 0

	base := this.GetPolyRefBase(tile)

	for i := 0; i < int(tile.Header.PolyCount); i++ {
		poly := &tile.Polys[i]
		nv := int(poly.VertCount)
		for j := 0; j < nv; j++ {
			// Skip edges which do not point to the right side.
			if poly.Neis[j] != m {
				continue
			}

			vc := tile.Verts[poly.Verts[j]*3:]
			vd := tile.Verts[poly.Verts[(j+1)%nv]*3:]
			bpos := getSlabCoord(vc, side)

			// Segments are not close enough.
			if DtAbsFloat32(apos-bpos) > 0.01 {
				continue
			}
			// Check if the segments touch.
			calcSlabEndPoints(vc, vd, bmin[:], bmax[:], side)

			if !overlapSlabs(amin[:], amax[:], bmin[:], bmax[:], 0.01, tile.Header.WalkableClimb) {
				continue
			}

			// Add return value.
			if n < maxcon {
				conarea[n*2+0] = DtMaxFloat32(amin[0], bmin[0])
				conarea[n*2+1] = DtMinFloat32(amax[0], bmax[0])
				con[n] = base | (DtPolyRef)(i)
				n++
			}
			break
		}
	}
	return n
}

/// Removes external links at specified side.
func (this *DtNavMesh) unconnectLinks(tile, target *DtMeshTile) {
	if tile == nil || target == nil {
		return
	}

	targetNum := this.DecodePolyIdTile(DtPolyRef(this.GetTileRef(target)))

	for i := 0; i < int(tile.Header.PolyCount); i++ {
		poly := &tile.Polys[i]
		j := poly.FirstLink
		pj := DT_NULL_LINK
		for j != DT_NULL_LINK {
			if this.DecodePolyIdTile(tile.Links[j].Ref) == targetNum {
				// Remove link.
				nj := tile.Links[j].Next
				if pj == DT_NULL_LINK {
					poly.FirstLink = nj
				} else {
					tile.Links[pj].Next = nj
				}
				freeLink(tile, j)
				j = nj
			} else {
				// Advance
				pj = j
				j = tile.Links[j].Next
			}
		}
	}
}

/// Builds external polygon links for a tile.
func (this *DtNavMesh) connectExtLinks(tile, target *DtMeshTile, side int) {
	if tile == nil {
		return
	}

	// Connect border links.
	for i := 0; i < int(tile.Header.PolyCount); i++ {
		poly := &tile.Polys[i]

		// Create new links.
		//		unsigned short m = DT_EXT_LINK | (unsigned short)side;

		nv := int(poly.VertCount)
		for j := 0; j < nv; j++ {
			// Skip non-portal edges.
			if (poly.Neis[j] & DT_EXT_LINK) == 0 {
				continue
			}
			dir := (int)(poly.Neis[j] & 0xff)
			if side != -1 && dir != side {
				continue
			}
			// Create new links
			va := tile.Verts[poly.Verts[j]*3:]
			vb := tile.Verts[poly.Verts[(j+1)%nv]*3:]
			var nei [4]DtPolyRef
			var neia [4 * 2]float32
			nnei := this.findConnectingPolys(va, vb, target, DtOppositeTile(dir), nei[:], neia[:], 4)
			for k := 0; k < nnei; k++ {
				idx := allocLink(tile)
				if idx != DT_NULL_LINK {
					link := &tile.Links[idx]
					link.Ref = nei[k]
					link.Edge = (uint8)(j)
					link.Side = (uint8)(dir)

					link.Next = poly.FirstLink
					poly.FirstLink = idx

					// Compress portal limits to a byte value.
					if dir == 0 || dir == 4 {
						tmin := (neia[k*2+0] - va[2]) / (vb[2] - va[2])
						tmax := (neia[k*2+1] - va[2]) / (vb[2] - va[2])
						if tmin > tmax {
							DtSwapFloat32(&tmin, &tmax)
						}
						link.Bmin = (uint8)(DtClampFloat32(tmin, 0.0, 1.0) * 255.0)
						link.Bmax = (uint8)(DtClampFloat32(tmax, 0.0, 1.0) * 255.0)
					} else if dir == 2 || dir == 6 {
						tmin := (neia[k*2+0] - va[0]) / (vb[0] - va[0])
						tmax := (neia[k*2+1] - va[0]) / (vb[0] - va[0])
						if tmin > tmax {
							DtSwapFloat32(&tmin, &tmax)
						}
						link.Bmin = (uint8)(DtClampFloat32(tmin, 0.0, 1.0) * 255.0)
						link.Bmax = (uint8)(DtClampFloat32(tmax, 0.0, 1.0) * 255.0)
					}
				}
			}
		}
	}
}

/// Builds external polygon links for a tile.
func (this *DtNavMesh) connectExtOffMeshLinks(tile, target *DtMeshTile, side int) {
	if tile == nil {
		return
	}

	// Connect off-mesh links.
	// We are interested on links which land from target tile to this tile.
	var oppositeSide uint16
	if side == -1 {
		oppositeSide = 0xff
	} else {
		oppositeSide = uint16(DtOppositeTile(side))
	}

	for i := 0; i < int(target.Header.OffMeshConCount); i++ {
		targetCon := &target.OffMeshCons[i]
		if uint16(targetCon.Side) != oppositeSide {
			continue
		}
		targetPoly := &target.Polys[targetCon.Poly]
		// Skip off-mesh connections which start location could not be connected at all.
		if targetPoly.FirstLink == DT_NULL_LINK {
			continue
		}
		halfExtents := [3]float32{targetCon.Rad, target.Header.WalkableClimb, targetCon.Rad}

		// Find polygon to connect to.
		p := targetCon.Pos[3:]
		var nearestPt [3]float32
		ref := this.findNearestPolyInTile(tile, p, halfExtents[:], nearestPt[:])
		if ref == 0 {
			continue
		}
		// findNearestPoly may return too optimistic results, further check to make sure.
		if DtSqrFloat32(nearestPt[0]-p[0])+DtSqrFloat32(nearestPt[2]-p[2]) > DtSqrFloat32(targetCon.Rad) {
			continue
		}
		// Make sure the location is on current mesh.
		v := target.Verts[targetPoly.Verts[1]*3:]
		DtVcopy(v, nearestPt[:])

		// Link off-mesh connection to target poly.
		idx := allocLink(target)
		if idx != DT_NULL_LINK {
			link := &target.Links[idx]
			link.Ref = ref
			link.Edge = 1
			link.Side = uint8(oppositeSide)
			link.Bmin = 0
			link.Bmax = 0
			// Add to linked list.
			link.Next = targetPoly.FirstLink
			targetPoly.FirstLink = idx
		}

		// Link target poly to off-mesh connection.
		if (targetCon.Flags & DT_OFFMESH_CON_BIDIR) != 0 {
			tidx := allocLink(tile)
			if tidx != DT_NULL_LINK {
				landPolyIdx := (uint16)(this.DecodePolyIdPoly(ref))
				landPoly := &tile.Polys[landPolyIdx]
				link := &tile.Links[tidx]
				link.Ref = this.GetPolyRefBase(target) | (DtPolyRef)(targetCon.Poly)
				link.Edge = 0xff
				if side == -1 {
					link.Side = 0xff
				} else {
					link.Side = uint8(side)
				}
				link.Bmin = 0
				link.Bmax = 0
				// Add to linked list.
				link.Next = landPoly.FirstLink
				landPoly.FirstLink = tidx
			}
		}
	}
}

/// Builds internal polygons links for a tile.
func (this *DtNavMesh) connectIntLinks(tile *DtMeshTile) {
	if tile == nil {
		return
	}

	base := this.GetPolyRefBase(tile)

	for i := 0; i < int(tile.Header.PolyCount); i++ {
		poly := &tile.Polys[i]
		poly.FirstLink = DT_NULL_LINK

		if poly.GetType() == DT_POLYTYPE_OFFMESH_CONNECTION {
			continue
		}

		// Build edge links backwards so that the links will be
		// in the linked list from lowest index to highest.
		for j := int(poly.VertCount - 1); j >= 0; j-- {
			// Skip hard and non-internal edges.
			if poly.Neis[j] == 0 || (poly.Neis[j]&DT_EXT_LINK) != 0 {
				continue
			}

			idx := allocLink(tile)
			if idx != DT_NULL_LINK {
				link := &tile.Links[idx]
				link.Ref = base | (DtPolyRef)(poly.Neis[j]-1)
				link.Edge = (uint8)(j)
				link.Side = 0xff
				link.Bmin = 0
				link.Bmax = 0
				// Add to linked list.
				link.Next = poly.FirstLink
				poly.FirstLink = idx
			}
		}
	}
}

/// Builds internal polygons links for a tile.
func (this *DtNavMesh) baseOffMeshLinks(tile *DtMeshTile) {
	if tile == nil {
		return
	}

	base := this.GetPolyRefBase(tile)

	// Base off-mesh connection start points.
	for i := 0; i < int(tile.Header.OffMeshConCount); i++ {
		con := &tile.OffMeshCons[i]
		poly := &tile.Polys[con.Poly]

		halfExtents := [3]float32{con.Rad, tile.Header.WalkableClimb, con.Rad}

		// Find polygon to connect to.
		p := con.Pos[:] // First vertex
		var nearestPt [3]float32
		ref := this.findNearestPolyInTile(tile, p, halfExtents[:], nearestPt[:])
		if ref == 0 {
			continue
		}
		// findNearestPoly may return too optimistic results, further check to make sure.
		if DtSqrFloat32(nearestPt[0]-p[0])+DtSqrFloat32(nearestPt[2]-p[2]) > DtSqrFloat32(con.Rad) {
			continue
		}
		// Make sure the location is on current mesh.
		v := tile.Verts[poly.Verts[0]*3:]
		DtVcopy(v, nearestPt[:])

		// Link off-mesh connection to target poly.
		idx := allocLink(tile)
		if idx != DT_NULL_LINK {
			link := &tile.Links[idx]
			link.Ref = ref
			link.Edge = 0
			link.Side = 0xff
			link.Bmin = 0
			link.Bmax = 0
			// Add to linked list.
			link.Next = poly.FirstLink
			poly.FirstLink = idx
		}

		// Start end-point is always connect back to off-mesh connection.
		tidx := allocLink(tile)
		if tidx != DT_NULL_LINK {
			landPolyIdx := (uint16)(this.DecodePolyIdPoly(ref))
			landPoly := &tile.Polys[landPolyIdx]
			link := &tile.Links[tidx]
			link.Ref = base | (DtPolyRef)(con.Poly)
			link.Edge = 0xff
			link.Side = 0xff
			link.Bmin = 0
			link.Bmax = 0
			// Add to linked list.
			link.Next = landPoly.FirstLink
			landPoly.FirstLink = tidx
		}
	}
}

/// Returns closest point on polygon.
func (this *DtNavMesh) closestPointOnPoly(ref DtPolyRef, pos, closest []float32, posOverPoly *bool) {
	var tile *DtMeshTile = nil
	var poly *DtPoly = nil
	this.GetTileAndPolyByRefUnsafe(ref, &tile, &poly)

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
		return
	}

	polysBase := uintptr(unsafe.Pointer(&(tile.Polys[0])))
	current := uintptr(unsafe.Pointer(poly))
	ip := uint32(current-polysBase) / sizeofPoly
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
}

/// Find nearest polygon within a tile.
func (this *DtNavMesh) findNearestPolyInTile(tile *DtMeshTile, center, halfExtents, nearestPt []float32) DtPolyRef {
	var bmin, bmax [3]float32
	DtVsub(bmin[:], center, halfExtents)
	DtVadd(bmax[:], center, halfExtents)

	// Get nearby polygons from proximity grid.
	var polys [128]DtPolyRef
	polyCount := this.queryPolygonsInTile(tile, bmin[:], bmax[:], polys[:], 128)

	// Find nearest polygon amongst the nearby polygons.
	var nearest DtPolyRef = 0
	var nearestDistanceSqr float32 = math.MaxFloat32
	for i := 0; i < polyCount; i++ {
		ref := polys[i]
		var closestPtPoly [3]float32
		var diff [3]float32
		posOverPoly := false
		var d float32
		this.closestPointOnPoly(ref, center, closestPtPoly[:], &posOverPoly)

		// If a point is directly over a polygon and closer than
		// climb height, favor that instead of straight line nearest point.
		DtVsub(diff[:], center, closestPtPoly[:])
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

		if d < nearestDistanceSqr {
			DtVcopy(nearestPt, closestPtPoly[:])
			nearestDistanceSqr = d
			nearest = ref
		}
	}

	return nearest
}

/// Queries polygons within a tile.
func (this *DtNavMesh) queryPolygonsInTile(tile *DtMeshTile, qmin, qmax []float32, polys []DtPolyRef, maxPolys int) int {
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
		bmin[0] = uint16(qfac*minx) & 0xfffe
		bmin[1] = uint16(qfac*miny) & 0xfffe
		bmin[2] = uint16(qfac*minz) & 0xfffe
		bmax[0] = uint16(qfac*maxx+1) | 1
		bmax[1] = uint16(qfac*maxy+1) | 1
		bmax[2] = uint16(qfac*maxz+1) | 1

		// Traverse tree
		base := this.GetPolyRefBase(tile)
		n := 0
		for nodeIndex < endIndex {
			node := &tile.BvTree[nodeIndex]
			overlap := DtOverlapQuantBounds(bmin[:], bmax[:], node.Bmin[:], node.Bmax[:])
			isLeafNode := (node.I >= 0)

			if isLeafNode && overlap {
				if n < maxPolys {
					polys[n] = base | (DtPolyRef)(node.I)
					n = n + 1
				}
			}

			if overlap || isLeafNode {
				nodeIndex++
			} else {
				escapeIndex := int(-node.I)
				nodeIndex += escapeIndex
			}
		}

		return n
	} else {
		var bmin, bmax [3]float32
		n := 0
		base := this.GetPolyRefBase(tile)
		for i := 0; i < int(tile.Header.PolyCount); i++ {
			p := &tile.Polys[i]
			// Do not return off-mesh connection polygons.
			if p.GetType() == DT_POLYTYPE_OFFMESH_CONNECTION {
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
				if n < maxPolys {
					polys[n] = base | (DtPolyRef)(i)
					n = n + 1
				}
			}
		}
		return n
	}
}

/// Adds a tile to the navigation mesh.
///  @param[in]		data		Data for the new tile mesh. (See: #dtCreateNavMeshData)
///  @param[in]		dataSize	Data size of the new tile mesh.
///  @param[in]		flags		Tile flags. (See: #dtTileFlags)
///  @param[in]		lastRef		The desired reference for the tile. (When reloading a tile.) [opt] [Default: 0]
///  @param[out]	result		The tile reference. (If the tile was succesfully added.) [opt]
/// @return The status flags for the operation.
/// @par
///
/// The add operation will fail if the data is in the wrong format, the allocated tile
/// space is full, or there is a tile already at the specified reference.
///
/// The lastRef parameter is used to restore a tile with the same tile
/// reference it had previously used.  In this case the #dtPolyRef's for the
/// tile will be restored to the same values they were before the tile was
/// removed.
///
/// The nav mesh assumes exclusive access to the data passed and will make
/// changes to the dynamic portion of the data. For that reason the data
/// should not be reused in other nav meshes until the tile has been successfully
/// removed from this nav mesh.
///
/// @see dtCreateNavMeshData, #removeTile
func (this *DtNavMesh) AddTile(data []byte, dataSize int, flags DtTileFlags,
	lastRef DtTileRef, result *DtTileRef) DtStatus {

	// Make sure the data is in right format.
	if dataSize < DtAlign4(int(unsafe.Sizeof(DtMeshHeader{}))) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	header := (*DtMeshHeader)(unsafe.Pointer(&(data[0])))
	if header.Magic != DT_NAVMESH_MAGIC {
		return DT_FAILURE | DT_WRONG_MAGIC
	}
	if header.Version != DT_NAVMESH_VERSION {
		return DT_FAILURE | DT_WRONG_VERSION
	}

	// Make sure the location is free.
	if this.GetTileAt(header.X, header.Y, header.Layer) != nil {
		return DT_FAILURE | DT_ALREADY_OCCUPIED
	}

	// Allocate a tile.
	var tile *DtMeshTile = nil
	if lastRef == 0 {
		if this.m_nextFree != nil {
			tile = this.m_nextFree
			this.m_nextFree = tile.Next
			tile.Next = nil
		}
	} else {
		// Try to relocate the tile to specific index with same salt.
		tileIndex := this.DecodePolyIdTile(DtPolyRef(lastRef))
		if tileIndex >= uint32(this.m_maxTiles) {
			return DT_FAILURE | DT_OUT_OF_MEMORY
		}
		// Try to find the specific tile id from the free list.
		target := &this.m_tiles[tileIndex]
		var prev *DtMeshTile = nil
		tile = this.m_nextFree
		for tile != nil && tile != target {
			prev = tile
			tile = tile.Next
		}
		// Could not find the correct location.
		if tile != target {
			return DT_FAILURE | DT_OUT_OF_MEMORY
		}
		// Remove from freelist
		if prev == nil {
			this.m_nextFree = tile.Next
		} else {
			prev.Next = tile.Next
		}
		// Restore salt.
		tile.Salt = this.DecodePolyIdSalt(DtPolyRef(lastRef))
	}

	// Make sure we could allocate a tile.
	if tile == nil {
		return DT_FAILURE | DT_OUT_OF_MEMORY
	}
	// Insert tile into the position lut.
	h := computeTileHash(header.X, header.Y, int32(this.m_tileLutMask))
	tile.Next = this.m_posLookup[h]
	this.m_posLookup[h] = tile

	// Patch header pointers.
	headerSize := DtAlign4(int(unsafe.Sizeof(DtMeshHeader{})))
	vertsSize := DtAlign4(int(unsafe.Sizeof(float32(1.0))) * 3 * int(header.VertCount))
	polysSize := DtAlign4(int(unsafe.Sizeof(DtPoly{})) * int(header.PolyCount))
	linksSize := DtAlign4(int(unsafe.Sizeof(DtLink{})) * int(header.MaxLinkCount))
	detailMeshesSize := DtAlign4(int(unsafe.Sizeof(DtPolyDetail{})) * int(header.DetailMeshCount))
	detailVertsSize := DtAlign4(int(unsafe.Sizeof(float32(1.0))) * 3 * int(header.DetailVertCount))
	detailTrisSize := DtAlign4(int(unsafe.Sizeof(uint8(1))) * 4 * int(header.DetailTriCount))
	bvtreeSize := DtAlign4(int(unsafe.Sizeof(DtBVNode{})) * int(header.BvNodeCount))
	offMeshLinksSize := DtAlign4(int(unsafe.Sizeof(DtOffMeshConnection{})) * int(header.OffMeshConCount))

	d := 0 + headerSize

	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&(tile.Verts))))
	sliceHeader.Cap = 3 * int(header.VertCount)
	sliceHeader.Len = 3 * int(header.VertCount)
	sliceHeader.Data = uintptr(unsafe.Pointer(&(data[d])))
	d += vertsSize

	sliceHeader = (*reflect.SliceHeader)((unsafe.Pointer(&(tile.Polys))))
	sliceHeader.Cap = int(header.PolyCount)
	sliceHeader.Len = int(header.PolyCount)
	sliceHeader.Data = uintptr(unsafe.Pointer(&(data[d])))
	d += polysSize

	sliceHeader = (*reflect.SliceHeader)((unsafe.Pointer(&(tile.Links))))
	sliceHeader.Cap = int(header.MaxLinkCount)
	sliceHeader.Len = int(header.MaxLinkCount)
	sliceHeader.Data = uintptr(unsafe.Pointer(&(data[d])))
	d += linksSize

	if header.DetailMeshCount != 0 {
		sliceHeader = (*reflect.SliceHeader)((unsafe.Pointer(&(tile.DetailMeshes))))
		sliceHeader.Cap = int(header.DetailMeshCount)
		sliceHeader.Len = int(header.DetailMeshCount)
		sliceHeader.Data = uintptr(unsafe.Pointer(&(data[d])))
		d += detailMeshesSize
	}

	if header.DetailVertCount != 0 {
		sliceHeader = (*reflect.SliceHeader)((unsafe.Pointer(&(tile.DetailVerts))))
		sliceHeader.Cap = 3 * int(header.DetailVertCount)
		sliceHeader.Len = 3 * int(header.DetailVertCount)
		sliceHeader.Data = uintptr(unsafe.Pointer(&(data[d])))
		d += detailVertsSize
	}

	if header.DetailTriCount != 0 {
		sliceHeader = (*reflect.SliceHeader)((unsafe.Pointer(&(tile.DetailTris))))
		sliceHeader.Cap = 4 * int(header.DetailTriCount)
		sliceHeader.Len = 4 * int(header.DetailTriCount)
		sliceHeader.Data = uintptr(unsafe.Pointer(&(data[d])))
		d += detailTrisSize
	}

	if header.BvNodeCount != 0 {
		sliceHeader = (*reflect.SliceHeader)((unsafe.Pointer(&(tile.BvTree))))
		sliceHeader.Cap = int(header.BvNodeCount)
		sliceHeader.Len = int(header.BvNodeCount)
		sliceHeader.Data = uintptr(unsafe.Pointer(&(data[d])))
		d += bvtreeSize
	}

	if header.OffMeshConCount != 0 {
		sliceHeader = (*reflect.SliceHeader)((unsafe.Pointer(&(tile.OffMeshCons))))
		sliceHeader.Cap = int(header.OffMeshConCount)
		sliceHeader.Len = int(header.OffMeshConCount)
		sliceHeader.Data = uintptr(unsafe.Pointer(&(data[d])))
		d += offMeshLinksSize
	}

	// If there are no items in the bvtree, reset the tree pointer.
	if header.BvNodeCount == 0 {
		tile.BvTree = nil
	}

	// Build links freelist
	tile.LinksFreeList = 0
	tile.Links[header.MaxLinkCount-1].Next = DT_NULL_LINK
	for i := 0; i < int(header.MaxLinkCount-1); i++ {
		tile.Links[i].Next = uint32(i + 1)
	}

	// Init tile.
	tile.Header = header
	tile.Data = data
	tile.DataSize = int32(dataSize)
	tile.Flags = flags

	this.connectIntLinks(tile)

	// Base off-mesh connections to their starting polygons and connect connections inside the tile.
	this.baseOffMeshLinks(tile)
	this.connectExtOffMeshLinks(tile, tile, -1)

	// Create connections with neighbour tiles.
	const MAX_NEIS int = 32
	var neis [MAX_NEIS]*DtMeshTile
	var nneis int

	// Connect with layers in current tile.
	nneis = this.GetTilesAt(header.X, header.Y, neis[:], MAX_NEIS)
	for j := 0; j < nneis; j++ {
		if neis[j] == tile {
			continue
		}
		this.connectExtLinks(tile, neis[j], -1)
		this.connectExtLinks(neis[j], tile, -1)
		this.connectExtOffMeshLinks(tile, neis[j], -1)
		this.connectExtOffMeshLinks(neis[j], tile, -1)
	}

	// Connect with neighbour tiles.
	for i := 0; i < 8; i++ {
		nneis = this.GetNeighbourTilesAt(header.X, header.Y, i, neis[:], MAX_NEIS)
		for j := 0; j < nneis; j++ {
			this.connectExtLinks(tile, neis[j], i)
			this.connectExtLinks(neis[j], tile, DtOppositeTile(i))
			this.connectExtOffMeshLinks(tile, neis[j], i)
			this.connectExtOffMeshLinks(neis[j], tile, DtOppositeTile(i))
		}
	}

	if result != nil {
		*result = this.GetTileRef(tile)
	}
	return DT_SUCCESS
}

/// Gets the tile at the specified grid location.
///  @param[in]	x		The tile's x-location. (x, y, layer)
///  @param[in]	y		The tile's y-location. (x, y, layer)
///  @param[in]	layer	The tile's layer. (x, y, layer)
/// @return The tile, or null if the tile does not exist.
func (this *DtNavMesh) GetTileAt(x, y, layer int32) *DtMeshTile {
	// Find tile based on hash.
	h := computeTileHash(x, y, int32(this.m_tileLutMask))
	tile := this.m_posLookup[h]
	for tile != nil {
		if tile.Header != nil &&
			tile.Header.X == int32(x) &&
			tile.Header.Y == int32(y) &&
			tile.Header.Layer == int32(layer) {
			return tile
		}
		tile = tile.Next
	}
	return nil
}

/// Returns neighbour tile based on side.
func (this *DtNavMesh) GetNeighbourTilesAt(x, y int32, side int, tiles []*DtMeshTile, maxTiles int) int {
	nx := x
	ny := y
	switch side {
	case 0:
		nx++
	case 1:
		nx++
		ny++
	case 2:
		ny++
	case 3:
		nx--
		ny++
	case 4:
		nx--
	case 5:
		nx--
		ny--
	case 6:
		ny--
	case 7:
		nx++
		ny--
	}
	return this.GetTilesAt(nx, ny, tiles, maxTiles)
}

/// Gets all tiles at the specified grid location. (All layers.)
///  @param[in]		x			The tile's x-location. (x, y)
///  @param[in]		y			The tile's y-location. (x, y)
///  @param[out]	tiles		A pointer to an array of tiles that will hold the result.
///  @param[in]		maxTiles	The maximum tiles the tiles parameter can hold.
/// @return The number of tiles returned in the tiles array.
/// @par
///
/// This function will not fail if the tiles array is too small to hold the
/// entire result set.  It will simply fill the array to capacity.
func (this *DtNavMesh) GetTilesAt(x, y int32, tiles []*DtMeshTile, maxTiles int) int {
	n := 0

	// Find tile based on hash.
	h := computeTileHash(x, y, this.m_tileLutMask)
	tile := this.m_posLookup[h]
	for tile != nil {
		if tile.Header != nil &&
			tile.Header.X == x &&
			tile.Header.Y == y {
			if n < maxTiles {
				tiles[n] = tile
				n++
			}
		}
		tile = tile.Next
	}

	return n
}

/// Gets the tile reference for the tile at specified grid location.
///  @param[in]	x		The tile's x-location. (x, y, layer)
///  @param[in]	y		The tile's y-location. (x, y, layer)
///  @param[in]	layer	The tile's layer. (x, y, layer)
/// @return The tile reference of the tile, or 0 if there is none.
func (this *DtNavMesh) GetTileRefAt(x, y, layer int32) DtTileRef {
	// Find tile based on hash.
	h := computeTileHash(x, y, this.m_tileLutMask)
	tile := this.m_posLookup[h]
	for tile != nil {
		if tile.Header != nil &&
			tile.Header.X == x &&
			tile.Header.Y == y &&
			tile.Header.Layer == layer {
			return this.GetTileRef(tile)
		}
		tile = tile.Next
	}
	return 0
}

/// Gets the tile for the specified tile reference.
///  @param[in]	ref		The tile reference of the tile to retrieve.
/// @return The tile for the specified reference, or null if the
///		reference is invalid.
func (this *DtNavMesh) GetTileByRef(ref DtTileRef) *DtMeshTile {
	if ref == 0 {
		return nil
	}
	tileIndex := this.DecodePolyIdTile((DtPolyRef)(ref))
	tileSalt := this.DecodePolyIdSalt((DtPolyRef)(ref))
	if (int32)(tileIndex) >= this.m_maxTiles {
		return nil
	}
	tile := &this.m_tiles[tileIndex]
	if tile.Salt != tileSalt {
		return nil
	}
	return tile
}

/// The maximum number of tiles supported by the navigation mesh.
/// @return The maximum number of tiles supported by the navigation mesh.
func (this *DtNavMesh) GetMaxTiles() int32 {
	return this.m_maxTiles
}

/// Returns pointer to tile in the tile array.
func (this *DtNavMesh) GetTile(i int) *DtMeshTile {
	return &this.m_tiles[i]
}

/// Calculates the tile grid location for the specified world position.
///  @param[in]	pos  The world position for the query. [(x, y, z)]
///  @param[out]	tx		The tile's x-location. (x, y)
///  @param[out]	ty		The tile's y-location. (x, y)
func (this *DtNavMesh) CalcTileLoc(pos []float32, tx, ty *int32) {
	*tx = (int32)(math.Floor(float64(pos[0]-this.m_orig[0]) / float64(this.m_tileWidth)))
	*ty = (int32)(math.Floor(float64(pos[2]-this.m_orig[2]) / float64(this.m_tileHeight)))
}

/// Gets the tile and polygon for the specified polygon reference.
///  @param[in]		ref		The reference for the a polygon.
///  @param[out]	tile	The tile containing the polygon.
///  @param[out]	poly	The polygon.
/// @return The status flags for the operation.
func (this *DtNavMesh) GetTileAndPolyByRef(ref DtPolyRef, tile **DtMeshTile, poly **DtPoly) DtStatus {
	if ref == 0 {
		return DT_FAILURE
	}
	var salt, it, ip uint32
	this.DecodePolyId(ref, &salt, &it, &ip)
	if it >= (uint32)(this.m_maxTiles) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if this.m_tiles[it].Salt != salt || this.m_tiles[it].Header == nil {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if ip >= (uint32)(this.m_tiles[it].Header.PolyCount) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	*tile = &(this.m_tiles[it])
	*poly = &(this.m_tiles[it].Polys[ip])
	return DT_SUCCESS
}

/// Returns the tile and polygon for the specified polygon reference.
///  @param[in]		ref		A known valid reference for a polygon.
///  @param[out]	tile	The tile containing the polygon.
///  @param[out]	poly	The polygon.
/// @par
///
/// @warning Only use this function if it is known that the provided polygon
/// reference is valid. This function is faster than #getTileAndPolyByRef, but
/// it does not validate the reference.
func (this *DtNavMesh) GetTileAndPolyByRefUnsafe(ref DtPolyRef, tile **DtMeshTile, poly **DtPoly) {
	var salt, it, ip uint32
	this.DecodePolyId(ref, &salt, &it, &ip)
	*tile = &(this.m_tiles[it])
	*poly = &(this.m_tiles[it].Polys[ip])
}

/// Checks the validity of a polygon reference.
///  @param[in]	ref		The polygon reference to check.
/// @return True if polygon reference is valid for the navigation mesh.
func (this *DtNavMesh) IsValidPolyRef(ref DtPolyRef) bool {
	if ref == 0 {
		return false
	}
	var salt, it, ip uint32
	this.DecodePolyId(ref, &salt, &it, &ip)
	if it >= (uint32)(this.m_maxTiles) {
		return false
	}
	if this.m_tiles[it].Salt != salt || this.m_tiles[it].Header == nil {
		return false
	}
	if ip >= (uint32)(this.m_tiles[it].Header.PolyCount) {
		return false
	}
	return true
}

/// Removes the specified tile from the navigation mesh.
///  @param[in]		ref			The reference of the tile to remove.
///  @param[out]	data		Data associated with deleted tile.
///  @param[out]	dataSize	Size of the data associated with deleted tile.
/// @return The status flags for the operation.
/// @par
///
/// This function returns the data for the tile so that, if desired,
/// it can be added back to the navigation mesh at a later point.
func (this *DtNavMesh) RemoveTile(ref DtTileRef, data *[]byte, dataSize *int) DtStatus {
	if ref == 0 {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	tileIndex := this.DecodePolyIdTile((DtPolyRef)(ref))
	tileSalt := this.DecodePolyIdSalt((DtPolyRef)(ref))
	if (int32)(tileIndex) >= this.m_maxTiles {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	tile := &this.m_tiles[tileIndex]
	if tile.Salt != tileSalt {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	// Remove tile from hash lookup.
	h := computeTileHash(tile.Header.X, tile.Header.Y, this.m_tileLutMask)
	var prev *DtMeshTile = nil
	cur := this.m_posLookup[h]
	for cur != nil {
		if cur == tile {
			if prev != nil {
				prev.Next = cur.Next
			} else {
				this.m_posLookup[h] = cur.Next
			}
			break
		}
		prev = cur
		cur = cur.Next
	}

	// Remove connections to neighbour tiles.
	const MAX_NEIS int = 32
	var neis [MAX_NEIS]*DtMeshTile
	var nneis int

	// Disconnect from other layers in current tile.
	nneis = this.GetTilesAt(tile.Header.X, tile.Header.Y, neis[:], MAX_NEIS)
	for j := 0; j < nneis; j++ {
		if neis[j] == tile {
			continue
		}
		this.unconnectLinks(neis[j], tile)
	}

	// Disconnect from neighbour tiles.
	for i := 0; i < 8; i++ {
		nneis = this.GetNeighbourTilesAt(tile.Header.X, tile.Header.Y, i, neis[:], MAX_NEIS)
		for j := 0; j < nneis; j++ {
			this.unconnectLinks(neis[j], tile)
		}
	}

	// Reset tile.
	if (tile.Flags & DT_TILE_FREE_DATA) != 0 {
		// Owns data
		tile.Data = nil
		tile.DataSize = 0
		if data != nil {
			*data = nil
		}
		if dataSize != nil {
			*dataSize = 0
		}
	} else {
		if data != nil {
			*data = tile.Data
		}
		if dataSize != nil {
			*dataSize = int(tile.DataSize)
		}
	}

	tile.Header = nil
	tile.Flags = 0
	tile.LinksFreeList = 0
	tile.Polys = nil
	tile.Verts = nil
	tile.Links = nil
	tile.DetailMeshes = nil
	tile.DetailVerts = nil
	tile.DetailTris = nil
	tile.BvTree = nil
	tile.OffMeshCons = nil

	// Update salt, salt should never be zero.
	tile.Salt = (tile.Salt + 1) & ((1 << this.m_saltBits) - 1)
	if tile.Salt == 0 {
		tile.Salt++
	}
	// Add to free list.
	tile.Next = this.m_nextFree
	this.m_nextFree = tile

	return DT_SUCCESS
}

/// Gets the tile reference for the specified tile.
///  @param[in]	tile	The tile.
/// @return The tile reference of the tile.
func (this *DtNavMesh) GetTileRef(tile *DtMeshTile) DtTileRef {
	if tile == nil {
		return 0
	}
	tileBase := uintptr(unsafe.Pointer(&(this.m_tiles[0])))
	current := uintptr(unsafe.Pointer(tile))
	it := (uint32)(current-tileBase) / sizeofMeshTile
	return (DtTileRef)(this.EncodePolyId(tile.Salt, it, 0))
}

/// Gets the polygon reference for the tile's base polygon.
///  @param[in]	tile		The tile.
/// @return The polygon reference for the base polygon in the specified tile.
/// @par
///
/// Example use case:
/// @code
///
/// const dtPolyRef base = navmesh->getPolyRefBase(tile);
/// for (int i = 0; i < tile->header->polyCount; ++i)
/// {
///     const dtPoly* p = &tile->polys[i];
///     const dtPolyRef ref = base | (dtPolyRef)i;
///
///     // Use the reference to access the polygon data.
/// }
/// @endcode
func (this *DtNavMesh) GetPolyRefBase(tile *DtMeshTile) DtPolyRef {
	if tile == nil {
		return 0
	}
	tileBase := uintptr(unsafe.Pointer(&(this.m_tiles[0])))
	current := uintptr(unsafe.Pointer(tile))
	it := (uint32)(current-tileBase) / sizeofMeshTile
	return this.EncodePolyId(tile.Salt, it, 0)
}

type dtTileState struct {
	magic   int32     // Magic number, used to identify the data.
	version int32     // Data version number.
	ref     DtTileRef // Tile ref at the time of storing the data.
}

type dtPolyState struct {
	flags uint16 // Flags (see dtPolyFlags).
	area  uint8  // Area ID of the polygon.
}

/// Gets the size of the buffer required by #storeTileState to store the specified tile's state.
///  @param[in]	tile	The tile.
/// @return The size of the buffer required to store the state.
func (this *DtNavMesh) GetTileStateSize(tile *DtMeshTile) int {
	if tile == nil {
		return 0
	}
	headerSize := DtAlign4(int(unsafe.Sizeof(dtTileState{})))
	polyStateSize := DtAlign4(int(unsafe.Sizeof(dtPolyState{})) * int(tile.Header.PolyCount))
	return headerSize + polyStateSize
}

/// Stores the non-structural state of the tile in the specified buffer. (Flags, area ids, etc.)
///  @param[in]		tile			The tile.
///  @param[out]	data			The buffer to store the tile's state in.
///  @param[in]		maxDataSize		The size of the data buffer. [Limit: >= #getTileStateSize]
/// @return The status flags for the operation.
/// @par
///
/// Tile state includes non-structural data such as polygon flags, area ids, etc.
/// @note The state data is only valid until the tile reference changes.
/// @see #getTileStateSize, #restoreTileState
func (this *DtNavMesh) StoreTileState(tile *DtMeshTile, data []byte, maxDataSize int) DtStatus {
	// Make sure there is enough space to store the state.
	sizeReq := this.GetTileStateSize(tile)
	if maxDataSize < sizeReq {
		return DT_FAILURE | DT_BUFFER_TOO_SMALL
	}

	tileState := (*dtTileState)(unsafe.Pointer(&(data[0])))
	var polyStates []dtPolyState
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&polyStates)))
	sliceHeader.Cap = int(tile.Header.PolyCount)
	sliceHeader.Len = int(tile.Header.PolyCount)
	sliceHeader.Data = uintptr(unsafe.Pointer(&(data[DtAlign4(int(unsafe.Sizeof(dtTileState{})))])))

	// Store tile state.
	tileState.magic = DT_NAVMESH_STATE_MAGIC
	tileState.version = DT_NAVMESH_STATE_VERSION
	tileState.ref = this.GetTileRef(tile)

	// Store per poly state.
	for i := 0; i < int(tile.Header.PolyCount); i++ {
		p := &tile.Polys[i]
		s := &polyStates[i]
		s.flags = p.Flags
		s.area = p.GetArea()
	}

	return DT_SUCCESS
}

/// Restores the state of the tile.
///  @param[in]	tile			The tile.
///  @param[in]	data			The new state. (Obtained from #storeTileState.)
///  @param[in]	maxDataSize		The size of the state within the data buffer.
/// @return The status flags for the operation.
/// @par
///
/// Tile state includes non-structural data such as polygon flags, area ids, etc.
/// @note This function does not impact the tile's #dtTileRef and #dtPolyRef's.
/// @see #storeTileState
func (this *DtNavMesh) RestoreTileState(tile *DtMeshTile, data []byte, maxDataSize int) DtStatus {
	// Make sure there is enough space to store the state.
	sizeReq := this.GetTileStateSize(tile)
	if maxDataSize < sizeReq {
		return DT_FAILURE | DT_BUFFER_TOO_SMALL
	}

	tileState := (*dtTileState)(unsafe.Pointer(&(data[0])))
	var polyStates []dtPolyState
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&polyStates)))
	sliceHeader.Cap = int(tile.Header.PolyCount)
	sliceHeader.Len = int(tile.Header.PolyCount)
	sliceHeader.Data = uintptr(unsafe.Pointer(&(data[DtAlign4(int(unsafe.Sizeof(dtTileState{})))])))

	// Check that the restore is possible.
	if tileState.magic != DT_NAVMESH_STATE_MAGIC {
		return DT_FAILURE | DT_WRONG_MAGIC
	}
	if tileState.version != DT_NAVMESH_STATE_VERSION {
		return DT_FAILURE | DT_WRONG_VERSION
	}
	if tileState.ref != this.GetTileRef(tile) {
		return DT_FAILURE | DT_INVALID_PARAM
	}

	// Restore per poly state.
	for i := 0; i < int(tile.Header.PolyCount); i++ {
		p := &tile.Polys[i]
		s := &polyStates[i]
		p.Flags = s.flags
		p.SetArea(s.area)
	}

	return DT_SUCCESS
}

/// Gets the endpoints for an off-mesh connection, ordered by "direction of travel".
///  @param[in]		prevRef		The reference of the polygon before the connection.
///  @param[in]		polyRef		The reference of the off-mesh connection polygon.
///  @param[out]	startPos	The start position of the off-mesh connection. [(x, y, z)]
///  @param[out]	endPos		The end position of the off-mesh connection. [(x, y, z)]
/// @return The status flags for the operation.
/// @par
///
/// Off-mesh connections are stored in the navigation mesh as special 2-vertex
/// polygons with a single edge. At least one of the vertices is expected to be
/// inside a normal polygon. So an off-mesh connection is "entered" from a
/// normal polygon at one of its endpoints. This is the polygon identified by
/// the prevRef parameter.
func (this *DtNavMesh) GetOffMeshConnectionPolyEndPoints(prevRef, polyRef DtPolyRef, startPos, endPos []float32) DtStatus {
	var salt, it, ip uint32

	if polyRef == 0 {
		return DT_FAILURE
	}
	// Get current polygon
	this.DecodePolyId(polyRef, &salt, &it, &ip)
	if it >= (uint32)(this.m_maxTiles) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if this.m_tiles[it].Salt != salt || this.m_tiles[it].Header == nil {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	tile := &this.m_tiles[it]
	if ip >= (uint32)(tile.Header.PolyCount) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	poly := &tile.Polys[ip]

	// Make sure that the current poly is indeed off-mesh link.
	if poly.GetType() != DT_POLYTYPE_OFFMESH_CONNECTION {
		return DT_FAILURE
	}
	// Figure out which way to hand out the vertices.
	idx0 := 0
	idx1 := 1

	// Find link that points to first vertex.
	for i := poly.FirstLink; i != DT_NULL_LINK; i = tile.Links[i].Next {
		if tile.Links[i].Edge == 0 {
			if tile.Links[i].Ref != prevRef {
				idx0 = 1
				idx1 = 0
			}
			break
		}
	}

	DtVcopy(startPos, tile.Verts[poly.Verts[idx0]*3:])
	DtVcopy(endPos, tile.Verts[poly.Verts[idx1]*3:])

	return DT_SUCCESS
}

/// Gets the specified off-mesh connection.
///  @param[in]	ref		The polygon reference of the off-mesh connection.
/// @return The specified off-mesh connection, or null if the polygon reference is not valid.
func (this *DtNavMesh) GetOffMeshConnectionByRef(ref DtPolyRef) *DtOffMeshConnection {
	var salt, it, ip uint32

	if ref == 0 {
		return nil
	}

	// Get current polygon
	this.DecodePolyId(ref, &salt, &it, &ip)
	if it >= (uint32)(this.m_maxTiles) {
		return nil
	}
	if this.m_tiles[it].Salt != salt || this.m_tiles[it].Header == nil {
		return nil
	}
	tile := &this.m_tiles[it]
	if ip >= (uint32)(tile.Header.PolyCount) {
		return nil
	}
	poly := &tile.Polys[ip]

	// Make sure that the current poly is indeed off-mesh link.
	if poly.GetType() != DT_POLYTYPE_OFFMESH_CONNECTION {
		return nil
	}

	idx := ip - uint32(tile.Header.OffMeshBase)
	DtAssert(idx < uint32(tile.Header.OffMeshConCount))
	return &tile.OffMeshCons[idx]
}

/// Sets the user defined flags for the specified polygon.
///  @param[in]	ref		The polygon reference.
///  @param[in]	flags	The new flags for the polygon.
/// @return The status flags for the operation.
func (this *DtNavMesh) SetPolyFlags(ref DtPolyRef, flags uint16) DtStatus {
	if ref == 0 {
		return DT_FAILURE
	}
	var salt, it, ip uint32
	this.DecodePolyId(ref, &salt, &it, &ip)
	if it >= (uint32)(this.m_maxTiles) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if this.m_tiles[it].Salt != salt || this.m_tiles[it].Header == nil {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	tile := &this.m_tiles[it]
	if ip >= (uint32)(tile.Header.PolyCount) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	poly := &tile.Polys[ip]

	// Change flags.
	poly.Flags = flags

	return DT_SUCCESS
}

/// Gets the user defined flags for the specified polygon.
///  @param[in]		ref				The polygon reference.
///  @param[out]	resultFlags		The polygon flags.
/// @return The status flags for the operation.
func (this *DtNavMesh) GetPolyFlags(ref DtPolyRef, resultFlags *uint16) DtStatus {
	if ref == 0 {
		return DT_FAILURE
	}
	var salt, it, ip uint32
	this.DecodePolyId(ref, &salt, &it, &ip)
	if it >= (uint32)(this.m_maxTiles) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if this.m_tiles[it].Salt != salt || this.m_tiles[it].Header == nil {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	tile := &this.m_tiles[it]
	if ip >= (uint32)(tile.Header.PolyCount) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	poly := &tile.Polys[ip]

	*resultFlags = poly.Flags

	return DT_SUCCESS
}

/// Sets the user defined area for the specified polygon.
///  @param[in]	ref		The polygon reference.
///  @param[in]	area	The new area id for the polygon. [Limit: < #DT_MAX_AREAS]
/// @return The status flags for the operation.
func (this *DtNavMesh) SetPolyArea(ref DtPolyRef, area uint8) DtStatus {
	if ref == 0 {
		return DT_FAILURE
	}
	var salt, it, ip uint32
	this.DecodePolyId(ref, &salt, &it, &ip)
	if it >= (uint32)(this.m_maxTiles) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if this.m_tiles[it].Salt != salt || this.m_tiles[it].Header == nil {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	tile := &this.m_tiles[it]
	if ip >= (uint32)(tile.Header.PolyCount) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	poly := &tile.Polys[ip]

	poly.SetArea(area)

	return DT_SUCCESS
}

/// Gets the user defined area for the specified polygon.
///  @param[in]		ref			The polygon reference.
///  @param[out]	resultArea	The area id for the polygon.
/// @return The status flags for the operation.
func (this *DtNavMesh) GetPolyArea(ref DtPolyRef, resultArea *uint8) DtStatus {
	if ref == 0 {
		return DT_FAILURE
	}
	var salt, it, ip uint32
	this.DecodePolyId(ref, &salt, &it, &ip)
	if it >= (uint32)(this.m_maxTiles) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	if this.m_tiles[it].Salt != salt || this.m_tiles[it].Header == nil {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	tile := &this.m_tiles[it]
	if ip >= (uint32)(tile.Header.PolyCount) {
		return DT_FAILURE | DT_INVALID_PARAM
	}
	poly := &tile.Polys[ip]

	*resultArea = poly.GetArea()

	return DT_SUCCESS
}


func overlapSlabs(amin, amax, bmin, bmax []float32, px, py float32) bool {
	// Check for horizontal overlap.
	// The segment is shrunken a little so that slabs which touch
	// at end points are not connected.
	minx := DtMaxFloat32(amin[0]+px, bmin[0]+px)
	maxx := DtMinFloat32(amax[0]-px, bmax[0]-px)
	if minx > maxx {
		return false
	}
	// Check vertical overlap.
	ad := (amax[1] - amin[1]) / (amax[0] - amin[0])
	ak := amin[1] - ad*amin[0]
	bd := (bmax[1] - bmin[1]) / (bmax[0] - bmin[0])
	bk := bmin[1] - bd*bmin[0]
	aminy := ad*minx + ak
	amaxy := ad*maxx + ak
	bminy := bd*minx + bk
	bmaxy := bd*maxx + bk
	dmin := bminy - aminy
	dmax := bmaxy - amaxy

	// Crossing segments always overlap.
	if dmin*dmax < 0 {
		return true
	}
	// Check for overlap at endpoints.
	thr := DtSqrFloat32(py * 2)
	if dmin*dmin <= thr || dmax*dmax <= thr {
		return true
	}
	return false
}

func getSlabCoord(va []float32, side int) float32 {
	if side == 0 || side == 4 {
		return va[0]
	} else if side == 2 || side == 6 {
		return va[2]
	}
	return 0
}

func calcSlabEndPoints(va, vb, bmin, bmax []float32, side int) {
	if side == 0 || side == 4 {
		if va[2] < vb[2] {
			bmin[0] = va[2]
			bmin[1] = va[1]
			bmax[0] = vb[2]
			bmax[1] = vb[1]
		} else {
			bmin[0] = vb[2]
			bmin[1] = vb[1]
			bmax[0] = va[2]
			bmax[1] = va[1]
		}
	} else if side == 2 || side == 6 {
		if va[0] < vb[0] {
			bmin[0] = va[0]
			bmin[1] = va[1]
			bmax[0] = vb[0]
			bmax[1] = vb[1]
		} else {
			bmin[0] = vb[0]
			bmin[1] = vb[1]
			bmax[0] = va[0]
			bmax[1] = va[1]
		}
	}
}

func computeTileHash(x, y, mask int32) int32 {
	h1 := uint32(0x8da6b343) // Large multiplicative constants;
	h2 := uint32(0xd8163841) // here arbitrarily chosen primes
	n := h1*uint32(x) + h2*uint32(y)
	return int32(n & uint32(mask))
}

func allocLink(tile *DtMeshTile) uint32 {
	if tile.LinksFreeList == DT_NULL_LINK {
		return DT_NULL_LINK
	}
	link := tile.LinksFreeList
	tile.LinksFreeList = tile.Links[link].Next
	return link
}

func freeLink(tile *DtMeshTile, link uint32) {
	tile.Links[link].Next = tile.LinksFreeList
	tile.LinksFreeList = link
}