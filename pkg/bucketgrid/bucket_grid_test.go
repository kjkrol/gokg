package bucketgrid

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
	"github.com/kjkrol/gokg/pkg/spatial"

	"github.com/stretchr/testify/assert"
)

var (
	NewAABBAt = geom.NewAABBAt[uint32]
	NewVec    = geom.NewVec[uint32]
)

func TestBucketGrid_Buckets_Len(t *testing.T) {

	overalRes := spatial.Size128x128
	bucketRes := spatial.Size32x32
	idx := NewBucketGrid(overalRes, bucketRes, 1)
	bucketGrid := idx.(*BucketGrid)
	bucketsLen := len(bucketGrid.buckets)

	if bucketsLen != 16 {
		t.Errorf("Buckets len %v is not equal to expected 16", bucketsLen)
	}

}

func TestBucketGrid_Buckets_QueryRange(t *testing.T) {

	// given
	overalRes := spatial.Size128x128
	bucketRes := spatial.Size16x16
	bucketGrid := NewBucketGrid(overalRes, bucketRes, 1)

	entries := []Entry{
		{AABB: NewAABBAt(NewVec(15, 15), 3, 3), Id: 1},
		{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: 2},
		{AABB: NewAABBAt(NewVec(30, 30), 3, 3), Id: 3},
		{AABB: NewAABBAt(NewVec(38, 38), 3, 3), Id: 4},
		{AABB: NewAABBAt(NewVec(39, 39), 3, 3), Id: 5},
	}
	bucketGrid.BulkInsert(entries)

	out := make([]uint64, 0, len(entries))
	seen := make(map[uint64]bool, len(entries))
	aabb := geom.NewAABBAt(Vec{X: 28, Y: 28}, 10, 10)
	expected := []uint64{2, 3, 4}

	// when
	bucketGrid.QueryRange(aabb, func(u uint64) {
		if !seen[u] {
			out = append(out, u)
			seen[u] = true
		}
	})

	// then
	assert.ElementsMatch(t, expected, out, "Should have exactly same elements")
}

func TestBucketGrid_Buckets_BulkMove(t *testing.T) {

	// given
	overalRes := spatial.Size128x128
	bucketRes := spatial.Size16x16
	bucketGrid := NewBucketGrid(overalRes, bucketRes, 1)

	entries := []Entry{
		{AABB: NewAABBAt(NewVec(15, 15), 3, 3), Id: 1},
		{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: 2},
		{AABB: NewAABBAt(NewVec(30, 30), 3, 3), Id: 3},
		{AABB: NewAABBAt(NewVec(38, 38), 3, 3), Id: 4},
		{AABB: NewAABBAt(NewVec(39, 39), 3, 3), Id: 5},
	}

	bucketGrid.BulkInsert(entries)

	out := make([]uint64, 0, len(entries))
	seen := make(map[uint64]bool, len(entries))
	aabb := geom.NewAABBAt(Vec{X: 28, Y: 28}, 10, 10)
	expected := []uint64{1, 3}

	entriesMove := EntriesMove{
		Old: []Entry{
			{AABB: NewAABBAt(NewVec(15, 15), 3, 3), Id: 1},
			{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: 2},
			{AABB: NewAABBAt(NewVec(38, 38), 3, 3), Id: 4},
		},
		New: []Entry{
			{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: 1},
			{AABB: NewAABBAt(NewVec(60, 25), 3, 3), Id: 2},
			{AABB: NewAABBAt(NewVec(70, 38), 3, 3), Id: 4},
		},
	}
	bucketGrid.BulkMove(entriesMove)

	// when
	bucketGrid.QueryRange(aabb, func(u uint64) {
		if !seen[u] {
			out = append(out, u)
			seen[u] = true
		}
	})

	// then
	assert.ElementsMatch(t, expected, out, "Should have exactly same elements")
}
