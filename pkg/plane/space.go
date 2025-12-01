package plane

import "github.com/kjkrol/gokg/pkg/geom"

const (
	modeEuclidean2D = "Euclidean2D"
	modeToroidal2D  = "Toroidal2D"
)

type (
	Space2D[T geom.Numeric] interface {
		WrapAABB(aabb geom.AABB[T]) AABB[T]
		WrapVec(vec geom.Vec[T]) AABB[T]
		Expand(aabb *AABB[T], margin T)
		Translate(aabb *AABB[T], delta geom.Vec[T])
		AABBDistance() AABBDistance[T]
		Name() string
		Viewport() geom.AABB[T]
	}

	Metric[T geom.Numeric] func(vec1, vec2 geom.Vec[T]) T
)

type space2d[T geom.Numeric] struct {
	size       geom.Vec[T]
	vectorMath geom.VectorMath[T]
	viewport   geom.AABB[T]
}
