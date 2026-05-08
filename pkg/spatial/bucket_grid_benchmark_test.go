package spatial

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func generateEntries(count int) []Entry {
	entries := make([]Entry, count)
	for i := range count {
		pos := NewVec(uint32(i%128), uint32((i/128)%128))
		entries[i] = Entry{
			Id:   EntryId(i),
			AABB: NewAABBAt(pos, 2, 2),
		}
	}
	return entries
}

func BenchmarkBucketGrid_BulkInsert_1000(b *testing.B) {
	entries := generateEntries(1000)

	bucketGrid, _ := NewBucketGrid(
		Size128x128,
		Size16x16,
		WithBucketCapacityFactor(1),
	)

	for b.Loop() {
		bucketGrid.Clear()
		bucketGrid.BulkInsert(entries)
	}
}

func BenchmarkBucketGrid_QueryRange(b *testing.B) {
	bucketGrid, _ := NewBucketGrid(
		Size128x128,
		Size16x16,
		WithBucketCapacityFactor(1),
	)

	entries := generateEntries(2000)
	bucketGrid.BulkInsert(entries)

	aabb := geom.NewAABBAt(NewVec(64, 64), 20, 20)

	b.Run("WithCollection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out := make([]uint64, 0, 100)
			bucketGrid.QueryRange(aabb, func(u uint64) { out = append(out, u) })
		}
	})

	b.Run("NoCollection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			count := 0
			bucketGrid.QueryRange(aabb, func(u uint64) { count++ })
		}
	})
}

func BenchmarkBucketGrid_BulkMove_100(b *testing.B) {
	bucketGrid, _ := NewBucketGrid(
		Size128x128,
		Size16x16,
		WithBucketCapacityFactor(1),
	)

	initialEntries := generateEntries(500)
	bucketGrid.BulkInsert(initialEntries)
	moveCount := 100
	oldEntries := initialEntries[:moveCount]
	newEntries := make([]Entry, moveCount)
	for i := range moveCount {
		newEntries[i] = oldEntries[i]
		newEntries[i].AABB = NewAABBAt(NewVec(uint32(i%100)+5, 5), 2, 2)
	}

	moveData := EntriesMove{
		Old: oldEntries,
		New: newEntries,
	}

	for b.Loop() {
		bucketGrid.BulkMove(moveData)
	}
}
