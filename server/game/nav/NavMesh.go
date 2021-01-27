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

import "gonet/server/game/lmath"

/// A handle to a polygon within a navigation mesh tile.
/// @ingroup detour
type DtPolyRef uint32

/// A handle to a tile within a navigation mesh.
/// @ingroup detour
type DtTileRef uint32

/// The maximum number of vertices per navigation polygon.
/// @ingroup detour
const DT_VERTS_PER_POLYGON int32 = 6

/// @{
/// @name Tile Serialization Constants
/// These constants are used to detect whether a navigation tile's data
/// and state format is compatible with the current build.
///

/// A magic number used to detect compatibility of navigation tile data.
const DT_NAVMESH_MAGIC int32 = 'D'<<24 | 'N'<<16 | 'A'<<8 | 'V'

/// A version number used to detect compatibility of navigation tile data.
const DT_NAVMESH_VERSION int32 = 7

/// A magic number used to detect the compatibility of navigation tile states.
const DT_NAVMESH_STATE_MAGIC int32 = 'D'<<24 | 'N'<<16 | 'M'<<8 | 'S'

/// A version number used to detect compatibility of navigation tile states.
const DT_NAVMESH_STATE_VERSION int32 = 1

/// @}

/// A flag that indicates that an entity links to an external entity.
/// (E.g. A polygon edge is a portal that links to another polygon.)
const DT_EXT_LINK uint16 = 0x8000

/// A value that indicates the entity does not link to anything.
const DT_NULL_LINK uint32 = 0xffffffff

/// A flag that indicates that an off-mesh connection can be traversed in both directions. (Is bidirectional.)
const DT_OFFMESH_CON_BIDIR uint8 = 1

/// The maximum number of user defined area ids.
/// @ingroup detour
const DT_MAX_AREAS int = 64

/// Tile flags used for various functions and fields.
/// For an example, see dtNavMesh::addTile().
type DtTileFlags int

const (
	/// The navigation mesh owns the tile memory and is responsible for freeing it.
	DT_TILE_FREE_DATA DtTileFlags = 0x01
)

/// Vertex flags returned by dtNavMeshQuery::findStraightPath.
type DtStraightPathFlags uint8

const (
	DT_STRAIGHTPATH_START              DtStraightPathFlags = 0x01 ///< The vertex is the start position in the path.
	DT_STRAIGHTPATH_END                DtStraightPathFlags = 0x02 ///< The vertex is the end position in the path.
	DT_STRAIGHTPATH_OFFMESH_CONNECTION DtStraightPathFlags = 0x04 ///< The vertex is the start of an off-mesh connection.
)

/// Options for dtNavMeshQuery::findStraightPath.
type DtStraightPathOptions int

const (
	DT_STRAIGHTPATH_AREA_CROSSINGS DtStraightPathOptions = 0x01 ///< Add a vertex at every polygon edge crossing where area changes.
	DT_STRAIGHTPATH_ALL_CROSSINGS  DtStraightPathOptions = 0x02 ///< Add a vertex at every polygon edge crossing.
)

/// Options for dtNavMeshQuery::initSlicedFindPath and updateSlicedFindPath
type DtFindPathOptions int

const (
	DT_FINDPATH_ANY_ANGLE DtFindPathOptions = 0x02 ///< use raycasts during pathfind to "shortcut" (raycast still consider costs)
)

/// Options for dtNavMeshQuery::raycast
type DtRaycastOptions int

const (
	DT_RAYCAST_USE_COSTS DtRaycastOptions = 0x01 ///< Raycast should calculate movement cost along the ray and fill RaycastHit::cost
)

/// Limit raycasting during any angle pahfinding
/// The limit is given as a multiple of the character radius
const DT_RAY_CAST_LIMIT_PROPORTIONS float32 = 50.0

/// Flags representing the type of a navigation mesh polygon.
type DtPolyTypes uint8

const (
	/// The polygon is a standard convex polygon that is part of the surface of the mesh.
	DT_POLYTYPE_GROUND DtPolyTypes = 0
	/// The polygon is an off-mesh connection consisting of two vertices.
	DT_POLYTYPE_OFFMESH_CONNECTION DtPolyTypes = 1
)

type(
	/// Defines a polygon within a dtMeshTile object.
	/// @ingroup detour
	DtPoly struct {
		/// Index to first link in linked list. (Or #DT_NULL_LINK if there is no link.)
		FirstLink uint32

		/// The indices of the polygon's vertices.
		/// The actual vertices are located in dtMeshTile::verts.
		Verts [DT_VERTS_PER_POLYGON]uint16

		/// Packed data representing neighbor polygons references and flags for each edge.
		Neis [DT_VERTS_PER_POLYGON]uint16

		/// The user defined polygon flags.
		Flags uint16

		/// The number of vertices in the polygon.
		VertCount uint8

		/// The bit packed area id and polygon type.
		/// @note Use the structure's set and get methods to acess this value.
		AreaAndtype uint8
	}

	/// Defines the location of detail sub-mesh data within a dtMeshTile.
	DtPolyDetail struct {
		VertBase  uint32 ///< The offset of the vertices in the dtMeshTile::detailVerts array.
		TriBase   uint32 ///< The offset of the triangles in the dtMeshTile::detailTris array.
		VertCount uint8  ///< The number of vertices in the sub-mesh.
		TriCount  uint8  ///< The number of triangles in the sub-mesh.
	}

	/// Defines a link between polygons.
	/// @note This structure is rarely if ever used by the end user.
	/// @see dtMeshTile
	DtLink struct {
		Ref  DtPolyRef ///< Neighbour reference. (The neighbor that is linked to.)
		Next uint32    ///< Index of the next link.
		Edge uint8     ///< Index of the polygon edge that owns this link.
		Side uint8     ///< If a boundary link, defines on which side the link is.
		Bmin uint8     ///< If a boundary link, defines the minimum sub-edge area.
		Bmax uint8     ///< If a boundary link, defines the maximum sub-edge area.
	}

	/// Bounding volume node.
	/// @note This structure is rarely if ever used by the end user.
	/// @see dtMeshTile
	DtBVNode struct {
		Bmin [3]uint16 ///< Minimum bounds of the node's AABB. [(x, y, z)]
		Bmax [3]uint16 ///< Maximum bounds of the node's AABB. [(x, y, z)]
		I    int32     ///< The node's index. (Negative for escape sequence.)
	}

	/// Defines an navigation mesh off-mesh connection within a dtMeshTile object.
	/// An off-mesh connection is a user defined traversable connection made up to two vertices.
	DtOffMeshConnection struct {
		/// The endpoints of the connection. [(ax, ay, az, bx, by, bz)]
		Pos [6]float32

		/// The radius of the endpoints. [Limit: >= 0]
		Rad float32

		/// The polygon reference of the connection within the tile.
		Poly uint16

		/// Link flags.
		/// @note These are not the connection's user defined flags. Those are assigned via the
		/// connection's dtPoly definition. These are link flags used for internal purposes.
		Flags uint8

		/// End point side.
		Side uint8

		/// The id of the offmesh connection. (User assigned when the navigation mesh is built.)
		UserId uint32
	}

	/// Provides high level information related to a dtMeshTile object.
	/// @ingroup detour
	DtMeshHeader struct {
		Magic           int32  ///< Tile magic number. (Used to identify the data format.)
		Version         int32  ///< Tile data format version number.
		X               int32  ///< The x-position of the tile within the dtNavMesh tile grid. (x, y, layer)
		Y               int32  ///< The y-position of the tile within the dtNavMesh tile grid. (x, y, layer)
		Layer           int32  ///< The layer of the tile within the dtNavMesh tile grid. (x, y, layer)
		UserId          uint32 ///< The user defined id of the tile.
		PolyCount       int32  ///< The number of polygons in the tile.
		VertCount       int32  ///< The number of vertices in the tile.
		MaxLinkCount    int32  ///< The number of allocated links.
		DetailMeshCount int32  ///< The number of sub-meshes in the detail mesh.

		/// The number of unique vertices in the detail mesh. (In addition to the polygon vertices.)
		DetailVertCount int32

		DetailTriCount  int32      ///< The number of triangles in the detail mesh.
		BvNodeCount     int32      ///< The number of bounding volume nodes. (Zero if bounding volumes are disabled.)
		OffMeshConCount int32      ///< The number of off-mesh connections.
		OffMeshBase     int32      ///< The index of the first polygon which is an off-mesh connection.
		WalkableHeight  float32    ///< The height of the agents using the tile.
		WalkableRadius  float32    ///< The radius of the agents using the tile.
		WalkableClimb   float32    ///< The maximum climb height of the agents using the tile.
		Bmin            [3]float32 ///< The minimum bounds of the tile's AABB. [(x, y, z)]
		Bmax            [3]float32 ///< The maximum bounds of the tile's AABB. [(x, y, z)]

		/// The bounding volume quantization factor.
		BvQuantFactor float32
	}

	/// Defines a navigation mesh tile.
	/// @ingroup detour
	DtMeshTile struct {
		Salt uint32 ///< Counter describing modifications to the tile.

		LinksFreeList uint32         ///< Index to the next free link.
		Header        *DtMeshHeader  ///< The tile header.
		Polys         []DtPoly       ///< The tile polygons. [Size: dtMeshHeader::polyCount]
		Verts         []float32      ///< The tile vertices. [Size: dtMeshHeader::vertCount]
		Links         []DtLink       ///< The tile links. [Size: dtMeshHeader::maxLinkCount]
		DetailMeshes  []DtPolyDetail ///< The tile's detail sub-meshes. [Size: dtMeshHeader::detailMeshCount]

		/// The detail mesh's unique vertices. [(x, y, z) * dtMeshHeader::detailVertCount]
		DetailVerts []float32

		/// The detail mesh's triangles. [(vertA, vertB, vertC) * dtMeshHeader::detailTriCount]
		DetailTris []uint8

		/// The tile bounding volume nodes. [Size: dtMeshHeader::bvNodeCount]
		/// (Will be null if bounding volumes are disabled.)
		BvTree []DtBVNode

		OffMeshCons []DtOffMeshConnection ///< The tile off-mesh connections. [Size: dtMeshHeader::offMeshConCount]

		Data     []byte      ///< The tile data. (Not directly accessed under normal situations.)
		DataSize int32       ///< Size of the tile data.
		Flags    DtTileFlags ///< Tile flags. (See: #dtTileFlags)
		Next     *DtMeshTile ///< The next free tile, or the next tile in the spatial grid.
	}

	/// Configuration parameters used to define multi-tile navigation meshes.
	/// The values are used to allocate space during the initialization of a navigation mesh.
	/// @see dtNavMesh::init()
	/// @ingroup detour
	DtNavMeshParams struct {
		Orig       [3]float32 ///< The world space origin of the navigation mesh's tile space. [(x, y, z)]
		TileWidth  float32    ///< The width of each tile. (Along the x-axis.)
		TileHeight float32    ///< The height of each tile. (Along the z-axis.)
		MaxTiles   uint32     ///< The maximum number of tiles the navigation mesh can contain.
		MaxPolys   uint32     ///< The maximum number of polygons each tile can contain.
	}

	/// A navigation mesh based on tiles of convex polygons.
	/// @ingroup detour
	DtNavMesh struct {
		m_params                  DtNavMeshParams ///< Current initialization params. TODO: do not store this info twice.
		m_orig                    [3]float32     ///< Origin of the tile (0,0)
		m_tileWidth, m_tileHeight float32         ///< Dimensions of each tile.
		m_maxTiles                int32           ///< Max number of tiles.
		m_tileLutSize             int32           ///< Tile hash lookup size (must be pot).
		m_tileLutMask             int32           ///< Tile hash lookup mask.

		m_posLookup []*DtMeshTile ///< Tile hash lookup.
		m_nextFree  *DtMeshTile   ///< Freelist of tiles.
		m_tiles     []DtMeshTile  ///< List of tiles.

		m_saltBits uint32 ///< Number of salt bits in the tile ID.
		m_tileBits uint32 ///< Number of tile bits in the tile ID.
		m_polyBits uint32 ///< Number of poly bits in the tile ID.

		mBounds 	  	  lmath.Box3F //地图大小
		mTileWidth  	  float32    ///< The width of each tile. (the max x-axis or y-axis )
		mOrig       	  lmath.Point3F ///< The world space origin of the navigation mesh's tile space. [(x, y, z)]
	}
)
/// @}

/// Allocates a navigation mesh object using the Detour allocator.
/// @return A navigation mesh that is ready for initialization, or null on failure.
///  @ingroup detour
func DtAllocNavMesh() *DtNavMesh {
	navmesh := &DtNavMesh{}
	navmesh.constructor()
	return navmesh
}

/// Frees the specified navigation mesh object using the Detour allocator.
///  @param[in]	navmesh		A navigation mesh allocated using #dtAllocNavMesh
///  @ingroup detour
func DtFreeNavMesh(navmesh *DtNavMesh) {
	if navmesh == nil {
		return
	}
	navmesh.destructor()
}
///////////////////////////////////////////////////////////////////////////

// This section contains detailed documentation for members that don't have
// a source file. It reduces clutter in the main section of the header.

/**

@typedef dtPolyRef
@par

Polygon references are subject to the same invalidate/preserve/restore
rules that apply to #dtTileRef's.  If the #dtTileRef for the polygon's
tile changes, the polygon reference becomes invalid.

Changing a polygon's flags, area id, etc. does not impact its polygon
reference.

@typedef dtTileRef
@par

The following changes will invalidate a tile reference:

- The referenced tile has been removed from the navigation mesh.
- The navigation mesh has been initialized using a different set
  of #dtNavMeshParams.

A tile reference is preserved/restored if the tile is added to a navigation
mesh initialized with the original #dtNavMeshParams and is added at the
original reference location. (E.g. The lastRef parameter is used with
dtNavMesh::addTile.)

Basically, if the storage structure of a tile changes, its associated
tile reference changes.


@var unsigned short dtPoly::neis[DT_VERTS_PER_POLYGON]
@par

Each entry represents data for the edge starting at the vertex of the same index.
E.g. The entry at index n represents the edge data for vertex[n] to vertex[n+1].

A value of zero indicates the edge has no polygon connection. (It makes up the
border of the navigation mesh.)

The information can be extracted as follows:
@code
neighborRef = neis[n] & 0xff; // Get the neighbor polygon reference.

if (neis[n] & #DT_EX_LINK)
{
    // The edge is an external (portal) edge.
}
@endcode

@var float dtMeshHeader::bvQuantFactor
@par

This value is used for converting between world and bounding volume coordinates.
For example:
@code
const float cs = 1.0f / tile->header->bvQuantFactor;
const dtBVNode* n = &tile->bvTree[i];
if (n->i >= 0)
{
    // This is a leaf node.
    float worldMinX = tile->header->bmin[0] + n->bmin[0]*cs;
    float worldMinY = tile->header->bmin[0] + n->bmin[1]*cs;
    // Etc...
}
@endcode

@struct dtMeshTile
@par

Tiles generally only exist within the context of a dtNavMesh object.

Some tile content is optional.  For example, a tile may not contain any
off-mesh connections.  In this case the associated pointer will be null.

If a detail mesh exists it will share vertices with the base polygon mesh.
Only the vertices unique to the detail mesh will be stored in #detailVerts.

@warning Tiles returned by a dtNavMesh object are not guarenteed to be populated.
For example: The tile at a location might not have been loaded yet, or may have been removed.
In this case, pointers will be null.  So if in doubt, check the polygon count in the
tile's header to determine if a tile has polygons defined.

@var float dtOffMeshConnection::pos[6]
@par

For a properly built navigation mesh, vertex A will always be within the bounds of the mesh.
Vertex B is not required to be within the bounds of the mesh.

*/
