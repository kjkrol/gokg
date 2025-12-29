package bucketgrid

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
	"github.com/kjkrol/gokg/pkg/spatial"
)

// Helper do generowania dużej ilości danych testowych
func generateEntries(count int) []Entry {
	entries := make([]Entry, count)
	for i := range count {
		// Rozrzucamy elementy losowo/równomiernie w obszarze 128x128
		pos := NewVec(uint32(i%128), uint32((i/128)%128))
		entries[i] = Entry{
			Id:   uint64(i),
			AABB: NewAABBAt(pos, 2, 2),
		}
	}
	return entries
}

// Benchmark wstawiania masowego (BulkInsert)
func BenchmarkBucketGrid_BulkInsert_1000(b *testing.B) {
	overalRes := spatial.Size128x128
	bucketRes := spatial.Size16x16
	entries := generateEntries(1000)

	// Resetujemy licznik, by nie liczyć czasu generowania danych
	for b.Loop() {
		grid := NewBucketGrid(overalRes, bucketRes, 1)
		grid.BulkInsert(entries)
	}
}

// Benchmark zapytania QueryRange (najważniejszy pod kątem wydajności kolektora)
func BenchmarkBucketGrid_QueryRange(b *testing.B) {
	overalRes := spatial.Size128x128
	bucketRes := spatial.Size16x16
	grid := NewBucketGrid(overalRes, bucketRes, 1)

	entries := generateEntries(2000)
	grid.BulkInsert(entries)

	aabb := geom.NewAABBAt(NewVec(64, 64), 20, 20)

	b.Run("WithMapDeduplication", func(b *testing.B) {
		// To testuje scenariusz, o który pytałeś wcześniej (mapa w środku)
		for i := 0; i < b.N; i++ {
			out := make([]uint64, 0, 100)
			seen := make(map[uint64]bool)
			grid.QueryRange(aabb, func(u uint64) {
				if !seen[u] {
					out = append(out, u)
					seen[u] = true
				}
			})
		}
	})

	b.Run("NoDeduplication", func(b *testing.B) {
		// Czysty czas trawersowania bez logiki usuwania duplikatów
		for i := 0; i < b.N; i++ {
			count := 0
			grid.QueryRange(aabb, func(u uint64) {
				count++
			})
		}
	})
}

// Benchmark aktualizacji pozycji (BulkMove)
func BenchmarkBucketGrid_BulkMove_100(b *testing.B) {
	overalRes := spatial.Size128x128
	bucketRes := spatial.Size16x16
	grid := NewBucketGrid(overalRes, bucketRes, 1)

	initialEntries := generateEntries(500)
	grid.BulkInsert(initialEntries)

	// Przygotowujemy dane do ruchu
	moveCount := 100
	oldEntries := initialEntries[:moveCount]
	newEntries := make([]Entry, moveCount)
	for i := range moveCount {
		newEntries[i] = oldEntries[i]
		// Przesuwamy o wektor (5,5)
		newEntries[i].AABB = NewAABBAt(NewVec(uint32(i%100)+5, 5), 2, 2)
	}

	moveData := EntriesMove{
		Old: oldEntries,
		New: newEntries,
	}

	for b.Loop() {
		grid.BulkMove(moveData)
	}
}
