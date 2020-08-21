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

type(

	/// Provides information about raycast hit
	/// filled by dtNavMeshQuery::raycast
	/// @ingroup detour
	DtRaycastHit struct {
		/// The hit parameter. (FLT_MAX if no wall hit.)
		T float32

		/// hitNormal	The normal of the nearest wall hit. [(x, y, z)]
		HitNormal [3]float32

		/// The index of the edge on the final polygon where the wall was hit.
		HitEdgeIndex int32

		/// Pointer to an array of reference ids of the visited polygons. [opt]
		Path []DtPolyRef

		/// The number of visited polygons. [opt]
		PathCount int32

		/// The maximum number of polygons the @p path array can hold.
		MaxPath int32

		///  The cost of the path until hit.
		PathCost float32
	}

	/// Provides custom polygon query behavior.
	/// Used by dtNavMeshQuery::queryPolygons.
	/// @ingroup detour
	DtPolyQuery interface {
		/// Called for each batch of unique polygons touched by the search area in dtNavMeshQuery::queryPolygons.
		/// This can be called multiple times for a single query.
		Process(tile *DtMeshTile, polys []*DtPoly, refs []DtPolyRef, count int)
	}

	dtQueryData struct {
		status           DtStatus
		lastBestNode     *DtNode
		lastBestNodeCost float32
		startRef         DtPolyRef
		endRef           DtPolyRef
		startPos         [3]float32
		endPos           [3]float32
		filter           *DtQueryFilter
		options          DtFindPathOptions
		raycastLimitSqr  float32
	}

	DtNavMeshQuery struct {
		m_nav          *DtNavMesh   ///< Pointer to navmesh data.
		m_query        dtQueryData  ///< Sliced query state.
		m_tinyNodePool *DtNodePool  ///< Pointer to small node pool.
		m_nodePool     *DtNodePool  ///< Pointer to node pool.
		m_openList     *DtNodeQueue ///< Pointer to open list queue.
	}

	/// Defines polygon filtering and traversal costs for navigation mesh query operations.
	/// @ingroup detour
	DtQueryFilter struct {
		m_areaCost     [DT_MAX_AREAS]float32 ///< Cost per area type. (Used by default implementation.)
		m_includeFlags uint16                ///< Flags for polygons that can be visited. (Used by default implementation.)
		m_excludeFlags uint16                ///< Flags for polygons that should not be visted. (Used by default implementation.)
	}
)


/// Allocates a query object using the Detour allocator.
/// @return An allocated query object, or null on failure.
/// @ingroup detour
func DtAllocNavMeshQuery() *DtNavMeshQuery {
	query := &DtNavMeshQuery{}
	query.constructor()
	return query
}

/// Frees the specified query object using the Detour allocator.
///  @param[in]		query		A query object allocated using #dtAllocNavMeshQuery
/// @ingroup detour
func DtFreeNavMeshQuery(query *DtNavMeshQuery) {
	if query == nil {
		return
	}
	query.destructor()
}

func DtAllocDtQueryFilter() *DtQueryFilter {
	filter := &DtQueryFilter{}
	filter.constructor()
	return filter
}

func DtFreeDtQueryFilter(filter *DtQueryFilter) {
	if filter == nil {
		return
	}
	filter.destructor()
}