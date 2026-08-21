package spatial

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/uid"
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

func TestBucketGrid_DefaultCellCodec(t *testing.T) {
	idx, _ := NewBucketGrid(Size128x128, Size32x32, WithBucketCapacity(1))
	bucketGrid := idx.(*bucketGrid)
	if _, ok := bucketGrid.gridCellCodec.(LinearCodeCodec); !ok {
		t.Errorf("default codec = %T, want LinearCodeCodec", bucketGrid.gridCellCodec)
	}
}

func TestBucketGrid_WithMortonCodec(t *testing.T) {
	idx, _ := NewBucketGrid(Size128x128, Size32x32, WithBucketCapacity(1), WithMortonCodec())
	bucketGrid := idx.(*bucketGrid)
	if _, ok := bucketGrid.gridCellCodec.(MortonCodeCodec); !ok {
		t.Errorf("codec = %T, want MortonCodeCodec", bucketGrid.gridCellCodec)
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
		{AABB: NewAABBAt(NewVec(15, 15), 3, 3), Id: withFrag(uid.UID64(1), 0)},
		{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: withFrag(uid.UID64(2), 0)},
		{AABB: NewAABBAt(NewVec(30, 30), 3, 3), Id: withFrag(uid.UID64(3), 0)},
		{AABB: NewAABBAt(NewVec(38, 38), 3, 3), Id: withFrag(uid.UID64(4), 0)},
		{AABB: NewAABBAt(NewVec(39, 39), 3, 3), Id: withFrag(uid.UID64(5), 0)},
	}
	bucketGrid.BulkInsert(entries)

	out := make([]uid.UID64, 0, len(entries))
	aabb := geom.NewAABBAt(Vec{X: 28, Y: 28}, 10, 10)
	expected := []uid.UID64{2, 3, 4}

	// when
	bucketGrid.QueryRange(aabb, func(u uid.UID64, frag plane.FragPosition) { out = append(out, u) })

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
		{AABB: NewAABBAt(NewVec(15, 15), 3, 3), Id: withFrag(uid.UID64(1), 0)},
		{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: withFrag(uid.UID64(2), 0)},
		{AABB: NewAABBAt(NewVec(30, 30), 3, 3), Id: withFrag(uid.UID64(3), 0)},
		{AABB: NewAABBAt(NewVec(38, 38), 3, 3), Id: withFrag(uid.UID64(4), 0)},
		{AABB: NewAABBAt(NewVec(39, 39), 3, 3), Id: withFrag(uid.UID64(5), 0)},
	}

	bucketGrid.BulkInsert(entries)

	out := make([]uid.UID64, 0, len(entries))
	aabb := geom.NewAABBAt(Vec{X: 28, Y: 28}, 10, 10)
	expected := []uid.UID64{1, 3}

	entriesMove := EntriesMove{
		Old: []Entry{
			{AABB: NewAABBAt(NewVec(15, 15), 3, 3), Id: withFrag(uid.UID64(1), 0)},
			{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: withFrag(uid.UID64(2), 0)},
			{AABB: NewAABBAt(NewVec(38, 38), 3, 3), Id: withFrag(uid.UID64(4), 0)},
		},
		New: []Entry{
			{AABB: NewAABBAt(NewVec(25, 25), 3, 3), Id: withFrag(uid.UID64(1), 0)},
			{AABB: NewAABBAt(NewVec(60, 25), 3, 3), Id: withFrag(uid.UID64(2), 0)},
			{AABB: NewAABBAt(NewVec(70, 38), 3, 3), Id: withFrag(uid.UID64(4), 0)},
		},
	}
	bucketGrid.BulkMove(entriesMove)

	// when
	bucketGrid.QueryRange(aabb, func(u uid.UID64, frag plane.FragPosition) { out = append(out, u) })

	// then
	assert.ElementsMatch(t, expected, out, "Should have exactly same elements")
}
