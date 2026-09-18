package plane

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
)

func TestEuclidean2D_Reposition(t *testing.T) {
	euclidean := NewEuclidean2D(10, 10)

	t.Run("box filling the whole world does not shrink or move", func(t *testing.T) {
		box := NewAABB(vec(0, 0), 10, 10)
		euclidean.Reposition(&box, vec(1, 1))
		expectAABBState(t, box, vec(0, 0), vec(10, 10), map[FragPosition][2]geom.Vec{})
	})

	t.Run("box with room clamps position, keeps size", func(t *testing.T) {
		box := NewAABB(vec(7, 7), 2, 2)
		euclidean.Reposition(&box, vec(5, 5))
		expectAABBState(t, box, vec(8, 8), vec(10, 10), map[FragPosition][2]geom.Vec{})
	})

	t.Run("large negative delta clamps to origin, keeps size", func(t *testing.T) {
		box := NewAABB(vec(7, 7), 2, 2)
		euclidean.Reposition(&box, vec(-20, -20))
		expectAABBState(t, box, vec(0, 0), vec(2, 2), map[FragPosition][2]geom.Vec{})
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
