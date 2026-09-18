package spatial

import (
	"math"
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/uid"
)

func quantGrid(t *testing.T) *GridIndexManager {
	t.Helper()
	m, err := NewGridIndexManager(plane.NewEuclidean2D(256, 256), GridIndexConfig{
		Resolution:       Size256x256,
		BucketResolution: Size32x32,
		BucketCapacity:   8,
	})
	if err != nil {
		t.Fatalf("NewGridIndexManager: %v", err)
	}
	return m
}

// The world used to be unsigned, so the type ruled out coordinates that have no
// cell. Now the code has to, and getting this wrong would not be a wrong answer
// — converting a negative or NaN float to uint32 is implementation-defined, so
// it would be a silently arbitrary one.
func TestCalculateGridIndex_RefusesCoordinatesWithNoCell(t *testing.T) {
	bg := quantGrid(t).bucketGrid

	cases := map[string]geom.Vec{
		"negative on X":     geom.NewVec(-1, 10),
		"negative on Y":     geom.NewVec(10, -1),
		"both negative":     geom.NewVec(-5, -5),
		"NaN":               geom.NewVec(math.NaN(), 10),
		"NaN on both":       geom.NewVec(math.NaN(), math.NaN()),
		"beyond uint32":     geom.NewVec(math.MaxUint32*2, 10),
		"positive infinity": geom.NewVec(math.Inf(1), 10),
	}
	for name, v := range cases {
		t.Run(name, func(t *testing.T) {
			if got := bg.CalculateGridIndex(v); got != -1 {
				t.Errorf("CalculateGridIndex(%v) = %d, want -1", v, got)
			}
		})
	}
}

// A coordinate is in the cell it falls inside, and the fraction does not move
// it: the boundary belongs to the cell that starts there.
func TestCalculateGridIndex_PlacesFractionsInTheCellTheyFallIn(t *testing.T) {
	bg := quantGrid(t).bucketGrid

	inFirst := bg.CalculateGridIndex(geom.NewVec(0, 0))
	stillFirst := bg.CalculateGridIndex(geom.NewVec(31.999, 0))
	nextOne := bg.CalculateGridIndex(geom.NewVec(32, 0))

	if inFirst != stillFirst {
		t.Errorf("0 and 31.999 landed in cells %d and %d, want the same one", inFirst, stillFirst)
	}
	if nextOne == inFirst {
		t.Errorf("32 landed in cell %d, the same as 0 — the boundary should start a new cell", nextOne)
	}
	if just := bg.CalculateGridIndex(geom.NewVec(32.5, 0)); just != nextOne {
		t.Errorf("32.5 landed in cell %d, want %d — the fraction must not move it", just, nextOne)
	}
}

// An entity placed at a fractional position is still findable, which is the
// whole point of letting the world be continuous.
func TestIndex_HoldsEntitiesAtFractionalPositions(t *testing.T) {
	m := quantGrid(t)
	id := uid.UID64(1)
	box := plane.NewAABB(geom.NewVec(40.5, 40.25), 8, 8)

	m.QueueInsert(id, box)
	m.Flush(nil)

	found := false
	m.QueryRange(geom.NewAABB(geom.NewVec(40, 40), geom.NewVec(50, 50)), func(got uid.UID64, _ plane.FragPosition) {
		found = found || got == id
	})
	if !found {
		t.Error("an entity at a fractional position is not findable where it stands")
	}
}
