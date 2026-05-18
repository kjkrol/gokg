package spatial

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/stretchr/testify/assert"
)

func TestBucketGrid_Buckets_Len(t *testing.T) {

	idx, _ := NewBucketGrid(
		Size128x128,
		Size32x32,
		WithBucketCapacityFactor(1),
	)
	bucketGrid := idx.(*bucketGrid)
	bucketsLen := len(bucketGrid.buckets)

	if bucketsLen != 16 {
		t.Errorf("Buckets len %v is not equal to expected 16", bucketsLen)
	}

}

func TestBucketGrid_Buckets_QueryRange(t *testing.T) {

	// given
	bucketGrid, _ := NewBucketGrid(
		Size128x128,
		Size16x16,
		WithBucketCapacityFactor(1),
	)

	entries := []Entry{
		{AABB: NewAABBAt(NewVec(15, 15), 3, 3), Id: NewEntryID(1, 0)},
		{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: NewEntryID(2, 0)},
		{AABB: NewAABBAt(NewVec(30, 30), 3, 3), Id: NewEntryID(3, 0)},
		{AABB: NewAABBAt(NewVec(38, 38), 3, 3), Id: NewEntryID(4, 0)},
		{AABB: NewAABBAt(NewVec(39, 39), 3, 3), Id: NewEntryID(5, 0)},
	}
	bucketGrid.BulkInsert(entries)

	out := make([]uint64, 0, len(entries))
	aabb := geom.NewAABBAt(Vec{X: 28, Y: 28}, 10, 10)
	expected := []uint64{2, 3, 4}

	// when
	bucketGrid.QueryRange(aabb, func(u uint64, frag plane.FragPosition) { out = append(out, u) })

	// then
	assert.ElementsMatch(t, expected, out, "Should have exactly same elements")
}

func TestBucketGrid_Buckets_BulkMove(t *testing.T) {

	// given
	bucketGrid, _ := NewBucketGrid(
		Size128x128,
		Size16x16,
		WithBucketCapacityFactor(1),
	)

	entries := []Entry{
		{AABB: NewAABBAt(NewVec(15, 15), 3, 3), Id: NewEntryID(1, 0)},
		{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: NewEntryID(2, 0)},
		{AABB: NewAABBAt(NewVec(30, 30), 3, 3), Id: NewEntryID(3, 0)},
		{AABB: NewAABBAt(NewVec(38, 38), 3, 3), Id: NewEntryID(4, 0)},
		{AABB: NewAABBAt(NewVec(39, 39), 3, 3), Id: NewEntryID(5, 0)},
	}

	bucketGrid.BulkInsert(entries)

	out := make([]uint64, 0, len(entries))
	aabb := geom.NewAABBAt(Vec{X: 28, Y: 28}, 10, 10)
	expected := []uint64{1, 3}

	entriesMove := EntriesMove{
		Old: []Entry{
			{AABB: NewAABBAt(NewVec(15, 15), 3, 3), Id: NewEntryID(1, 0)},
			{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: NewEntryID(2, 0)},
			{AABB: NewAABBAt(NewVec(38, 38), 3, 3), Id: NewEntryID(4, 0)},
		},
		New: []Entry{
			{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: NewEntryID(1, 0)},
			{AABB: NewAABBAt(NewVec(60, 25), 3, 3), Id: NewEntryID(2, 0)},
			{AABB: NewAABBAt(NewVec(70, 38), 3, 3), Id: NewEntryID(4, 0)},
		},
	}
	bucketGrid.BulkMove(entriesMove)

	// when
	bucketGrid.QueryRange(aabb, func(u uint64, frag plane.FragPosition) { out = append(out, u) })

	// then
	assert.ElementsMatch(t, expected, out, "Should have exactly same elements")
}
