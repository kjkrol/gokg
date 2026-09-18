package plane

import "github.com/kjkrol/gokg/geom"

const (
	modeEuclidean2D = "Euclidean2D"
	modeToroidal2D  = "Toroidal2D"
)

type (
	Space2D interface {
		Normalize(aabb geom.AABB) geom.AABB
		WrapAABB(aabb geom.AABB) AABB
		WrapVec(vec geom.Vec) AABB
		Expand(aabb *AABB, margin float64)
		Translate(aabb *AABB, delta geom.Vec)
		Reposition(aabb *AABB, delta geom.Vec)
		AABBDistance() AABBDistance
		Name() string
		Viewport() geom.AABB
	}

	Metric func(vec1, vec2 geom.Vec) float64
)

type space2d struct {
	size     geom.Vec
	viewport geom.AABB
}
