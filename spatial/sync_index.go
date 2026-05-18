package spatial

import (
	"sync"

	"github.com/kjkrol/gokg/plane"
)

type syncIndex struct {
	Index
	mu sync.RWMutex
}

var _ Index = (*syncIndex)(nil)

func SyncIndex(index Index) Index {
	return &syncIndex{
		Index: index,
	}
}

func (s *syncIndex) BulkInsert(entries []Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Index.BulkInsert(entries)
}
func (s *syncIndex) BulkRemove(entries []Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Index.BulkRemove(entries)
}
func (s *syncIndex) BulkMove(moves EntriesMove) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Index.BulkMove(moves)
}

// Collector cannot modify Index.
func (s *syncIndex) QueryRange(aabb AABB, collector func(uint64, plane.FragPosition)) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Index.QueryRange(aabb, collector)
}
func (s *syncIndex) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Index.Count()
}
func (s *syncIndex) Optimize() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Index.Optimize()
}
func (s *syncIndex) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Index.Clear()
}

func (s *syncIndex) CalculateGridIndex(vec Vec) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if indexer, ok := s.Index.(GridIndexer); ok {
		return indexer.CalculateGridIndex(vec)
	}
	return -1
}
