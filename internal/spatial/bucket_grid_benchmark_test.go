package spatial

import (
	"testing"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/uid"
)

func generateEntries(count int) []Entry {
	entries := make([]Entry, count)
	for i := range count {
		pos := NewVec(float64(i%128), float64((i/128)%128))
		entries[i] = Entry{
			Id:   uid.UID64(i),
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
			out := make([]uid.UID64, 0, 100)
			bucketGrid.QueryRange(aabb, func(u uid.UID64) { out = append(out, u) })
		}
	})

	b.Run("NoCollection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			count := 0
			bucketGrid.QueryRange(aabb, func(u uid.UID64) { count++ })
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
		newEntries[i].AABB = NewAABBAt(NewVec(float64(i%100)+5, 5), 2, 2)
	}

	moveData := EntriesMove{
		Old: oldEntries,
		New: newEntries,
	}

	for b.Loop() {
		bucketGrid.BulkMove(moveData)
	}
}

func BenchmarkBucketGrid_BulkMove(b *testing.B) {
	index, _ := NewBucketGrid(Size128x128, Size16x16, WithBucketCapacityFactor(1))
	bucketGrid := index.(*bucketGrid)

	entries := generateEntries(2000)
	bucketGrid.BulkInsert(entries)

	moves := EntriesMove{Old: make([]Entry, len(entries)), New: make([]Entry, len(entries))}
	for i, e := range entries {
		moves.Old[i] = e
		moves.New[i] = Entry{Id: e.Id, AABB: geom.NewAABB(
			e.AABB.TopLeft.Add(NewVec(0.5, 0)), e.AABB.BottomRight.Add(NewVec(0.5, 0)))}
	}

	for b.Loop() {
		bucketGrid.BulkMove(moves)
		moves.Old, moves.New = moves.New, moves.Old
	}
}
