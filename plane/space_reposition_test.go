package plane

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
)

func TestEuclidean2D_Reposition(t *testing.T) {
	runEuclidean2DRepositionTest[int](t, "int")
	runEuclidean2DRepositionTest[uint32](t, "uint32")
	runEuclidean2DRepositionTest[float64](t, "float64")
}

func runEuclidean2DRepositionTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		euclidean := NewEuclidean2D(T(10), T(10))

		t.Run("box filling the whole world does not shrink or move", func(t *testing.T) {
			box := NewAABB(vec[T](0, 0), T(10), T(10))
			euclidean.Reposition(&box, vec[T](1, 1))
			expectAABBState(t, box, vec[T](0, 0), vec[T](10, 10), map[FragPosition][2]geom.Vec[T]{})
		})

		t.Run("box with room clamps position, keeps size", func(t *testing.T) {
			box := NewAABB(vec[T](7, 7), T(2), T(2))
			euclidean.Reposition(&box, vec[T](5, 5))
			expectAABBState(t, box, vec[T](8, 8), vec[T](10, 10), map[FragPosition][2]geom.Vec[T]{})
		})

		t.Run("large negative delta clamps to origin, keeps size", func(t *testing.T) {
			box := NewAABB(vec[T](7, 7), T(2), T(2))
			euclidean.Reposition(&box, vec[T](-20, -20))
			expectAABBState(t, box, vec[T](0, 0), vec[T](2, 2), map[FragPosition][2]geom.Vec[T]{})
		})
	})
}

func TestToroidal2D_Reposition_MatchesTranslate(t *testing.T) {
	toroidal := NewToroidal2D(10, 10)

	viaReposition := NewAABB(geom.NewVec(9, 9), 2, 2)
	toroidal.Reposition(&viaReposition, geom.NewVec(3, 3))

	viaTranslate := NewAABB(geom.NewVec(9, 9), 2, 2)
	toroidal.Translate(&viaTranslate, geom.NewVec(3, 3))

	if viaReposition.TopLeft != viaTranslate.TopLeft || viaReposition.BottomRight != viaTranslate.BottomRight {
		t.Errorf("Reposition = (%+v,%+v), want same as Translate = (%+v,%+v)",
			viaReposition.TopLeft, viaReposition.BottomRight, viaTranslate.TopLeft, viaTranslate.BottomRight)
	}
}
