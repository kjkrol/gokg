package spatial

import (
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"testing"
)

func TestNewGridIndexManager_CellCodec(t *testing.T) {
	space := iplane.NewEuclidean2D(128, 128)

	m, err := NewGridIndexManager(space, GridIndexConfig{
		Resolution:       Size128x128,
		BucketResolution: Size32x32,
		BucketCapacity:   1,
	})
	if err != nil {
		t.Fatalf("NewGridIndexManager (default) error: %v", err)
	}
	if _, ok := m.bucketGrid.gridCellCodec.(LinearCodeCodec); !ok {
		t.Errorf("default CellCodec = %T, want LinearCodeCodec", m.bucketGrid.gridCellCodec)
	}

	m, err = NewGridIndexManager(space, GridIndexConfig{
		Resolution:       Size128x128,
		BucketResolution: Size32x32,
		BucketCapacity:   1,
		CellCodec:        MortonCellCodec,
	})
	if err != nil {
		t.Fatalf("NewGridIndexManager (morton) error: %v", err)
	}
	if _, ok := m.bucketGrid.gridCellCodec.(MortonCodeCodec); !ok {
		t.Errorf("CellCodec: MortonCellCodec = %T, want MortonCodeCodec", m.bucketGrid.gridCellCodec)
	}
}
