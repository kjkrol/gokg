package spatial

import (
	"fmt"
	"slices"
	"sync"

	"github.com/kjkrol/gokg/plane"
)

type (
	bucketGrid struct {
		resolution        Resolution
		bucketsResolution Resolution
		bucketCapacity    int
		gridResolution    Resolution
		gridCellCodec     CellCodec
		aabbById          map[EntryId]AABB
		bounds            AABB
		buckets           []bucket
		optimizer         *memoryOptimizer
	}

	bucket struct {
		ids []EntryId
	}

	memoryOptimizer struct {
		dirtyQueue []int
		isDirty    []bool
	}

	Option func(*bucketGrid) error
)

var (
	_            Index = (*bucketGrid)(nil)
	queryMapPool       = sync.Pool{
		New: func() any {
			return make(map[EntryId]struct{}, 1024)
		},
	}
)

func NewBucketGrid(
	overallResolution Resolution,
	bucketsResolution Resolution,
	opts ...Option,
) (Index, error) {
	gridResolution := NewResolution(uint8(overallResolution - bucketsResolution))
	gridCellCodec := NewLinearCodec(gridResolution)
	side := overallResolution.Side()

	bg := &bucketGrid{
		resolution:        overallResolution,
		bucketsResolution: bucketsResolution,
		gridResolution:    gridResolution,
		gridCellCodec:     gridCellCodec,
		bounds:            NewAABBAt(NewVec(0, 0), side, side),
	}
	for _, opt := range opts {
		err := opt(bg)
		if err != nil {
			return nil, err
		}
	}

	if bg.aabbById == nil {
		return nil, fmt.Errorf("Initialize bucket capacity first")
	}

	return bg, nil
}

func WithBucketCapacityFactor(capacityFactor float64) Option {
	return func(bg *bucketGrid) error {
		cellsCount := float64(bg.bucketsResolution.Cells())
		bucketCapacity := max(int(cellsCount*capacityFactor), 2)
		return WithBucketCapacity(bucketCapacity)(bg)
	}
}

func WithBucketCapacity(bucketCapacity int) Option {
	return func(bg *bucketGrid) error {
		bucketsNumber := bg.gridResolution.Cells()
		overallCapacity := bucketCapacity * bucketsNumber
		bg.bucketCapacity = bucketCapacity

		bg.buckets = make([]bucket, bucketsNumber)
		bg.optimizer = newMemoryOptimizer(int(bucketsNumber))
		bg.aabbById = make(map[EntryId]AABB, overallCapacity)
		return nil
	}
}

func (bg *bucketGrid) BulkInsert(entries []Entry) {
	for _, entry := range entries {
		tlIdx := bg.CalculateGridIndex(entry.AABB.TopLeft)
		brIdx := bg.CalculateGridIndex(entry.AABB.BottomRight)
		if tlIdx < 0 || brIdx < 0 || tlIdx >= len(bg.buckets) || brIdx >= len(bg.buckets) {
			continue
		}

		if tlIdx == brIdx {
			bg.buckets[tlIdx].Add(entry.Id, bg.bucketCapacity)
		} else {
			x1, y1 := bg.gridCellCodec.Decode(tlIdx)
			x2, y2 := bg.gridCellCodec.Decode(brIdx)
			for y := y1; y <= y2; y++ {
				for x := x1; x <= x2; x++ {
					idx, err := bg.gridCellCodec.Encode(x, y)
					if err != nil || idx < 0 || idx >= len(bg.buckets) {
						continue
					}
					bg.buckets[idx].Add(entry.Id, bg.bucketCapacity)
				}
			}
		}
		bg.aabbById[entry.Id] = entry.AABB
	}
}

// BulkRemove – remove whatever is stored at the given positions.
func (bg *bucketGrid) BulkRemove(entries []Entry) {
	for _, entry := range entries {
		tlIdx := bg.CalculateGridIndex(entry.AABB.TopLeft)
		brIdx := bg.CalculateGridIndex(entry.AABB.BottomRight)
		if tlIdx < 0 || brIdx < 0 || tlIdx >= len(bg.buckets) || brIdx >= len(bg.buckets) {
			continue
		}

		if tlIdx == brIdx {
			if bg.buckets[tlIdx].Remove(entry.Id) {
				bg.optimizer.mark(tlIdx, len(bg.buckets[tlIdx].ids) == 0)
			}
		} else {
			x1, y1 := bg.gridCellCodec.Decode(tlIdx)
			x2, y2 := bg.gridCellCodec.Decode(brIdx)
			for y := y1; y <= y2; y++ {
				for x := x1; x <= x2; x++ {
					idx, err := bg.gridCellCodec.Encode(x, y)
					if err != nil || idx < 0 || idx >= len(bg.buckets) {
						continue
					}
					if bg.buckets[idx].Remove(entry.Id) {
						bg.optimizer.mark(idx, len(bg.buckets[idx].ids) == 0)
					}
				}
			}
		}
		delete(bg.aabbById, entry.Id)
	}
}

// BulkMove – update objects (typically same Value, different XY).
func (bg *bucketGrid) BulkMove(moves EntriesMove) {
	for i := range moves.New {
		oldEntry := moves.Old[i]
		newEntry := moves.New[i]
		if oldEntry == newEntry {
			continue
		}

		oldTl := bg.CalculateGridIndex(oldEntry.AABB.TopLeft)
		oldBr := bg.CalculateGridIndex(oldEntry.AABB.BottomRight)
		newTl := bg.CalculateGridIndex(newEntry.AABB.TopLeft)
		newBr := bg.CalculateGridIndex(newEntry.AABB.BottomRight)
		if oldTl < 0 || oldBr < 0 || newTl < 0 || newBr < 0 {
			continue
		}
		if oldTl >= len(bg.buckets) || oldBr >= len(bg.buckets) || newTl >= len(bg.buckets) || newBr >= len(bg.buckets) {
			continue
		}

		if oldTl != newTl || oldBr != newBr {
			if oldTl == oldBr {
				bg.buckets[oldTl].Remove(oldEntry.Id)
			} else {
				x1, y1 := bg.gridCellCodec.Decode(oldTl)
				x2, y2 := bg.gridCellCodec.Decode(oldBr)
				for y := y1; y <= y2; y++ {
					for x := x1; x <= x2; x++ {
						idx, err := bg.gridCellCodec.Encode(x, y)
						if err != nil || idx < 0 || idx >= len(bg.buckets) {
							continue
						}
						bg.buckets[idx].Remove(oldEntry.Id)
					}
				}
			}

			if newTl == newBr {
				bg.buckets[newTl].Add(newEntry.Id, bg.bucketCapacity)
			} else {
				x1, y1 := bg.gridCellCodec.Decode(newTl)
				x2, y2 := bg.gridCellCodec.Decode(newBr)
				for y := y1; y <= y2; y++ {
					for x := x1; x <= x2; x++ {
						idx, err := bg.gridCellCodec.Encode(x, y)
						if err != nil || idx < 0 || idx >= len(bg.buckets) {
							continue
						}
						bg.buckets[idx].Add(newEntry.Id, bg.bucketCapacity)
					}
				}
			}
		}
		bg.aabbById[newEntry.Id] = newEntry.AABB
	}
}

// QueryRange – all objects within the AABB.
func (bg *bucketGrid) QueryRange(aabb AABB, collector func(uint64, plane.FragPosition)) int {

	tlIdx := bg.CalculateGridIndex(aabb.TopLeft)
	brIdx := bg.CalculateGridIndex(aabb.BottomRight)
	if tlIdx < 0 || brIdx < 0 || tlIdx >= len(bg.buckets) || brIdx >= len(bg.buckets) {
		return 0
	}
	counter := 0

	// Optimal path
	if tlIdx == brIdx {
		bucket := bg.buckets[tlIdx]
		for _, id := range bucket.ids {
			itemAABB, ok := bg.aabbById[id]
			if ok && aabb.Intersects(itemAABB) {
				collector(id.OriginalID(), plane.FragPosition(id.ExtractFrag()))
				counter++
			}
		}
		return counter
	}

	// Full path
	seen := queryMapPool.Get().(map[EntryId]struct{})
	defer queryMapPool.Put(seen)
	for k := range seen {
		delete(seen, k)
	}

	x1, y1 := bg.gridCellCodec.Decode(tlIdx)
	x2, y2 := bg.gridCellCodec.Decode(brIdx)

	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			idx, err := bg.gridCellCodec.Encode(x, y)
			if err != nil || idx < 0 || idx >= len(bg.buckets) {
				continue
			}
			bucket := bg.buckets[idx]

			for _, id := range bucket.ids {
				if _, ok := seen[id]; ok {
					continue
				}

				itemAABB, ok := bg.aabbById[id]
				if ok && aabb.Intersects(itemAABB) {
					seen[id] = struct{}{}
					collector(id.OriginalID(), plane.FragPosition(id.ExtractFrag()))
					counter++
				}
			}
		}
	}

	return counter
}

// Count – number of objects in the structure.
func (bg *bucketGrid) Count() int { return len(bg.aabbById) }

// Bounds – global bounds of the handled space.
func (bg *bucketGrid) Bounds() AABB { return bg.bounds }

func (bg *bucketGrid) Clear() {
	bg.aabbById = make(map[EntryId]AABB, len(bg.aabbById))
	for i := range bg.buckets {
		bg.buckets[i].ids = nil
	}
	bg.optimizer.clear()
}

func (bg *bucketGrid) Optimize() {
	bg.optimizer.collect(bg.buckets)
}

func (bg *bucketGrid) CalculateGridIndex(vec Vec) int {
	xHead := vec.X >> bg.bucketsResolution
	yHead := vec.Y >> bg.bucketsResolution
	idx, err := bg.gridCellCodec.Encode(xHead, yHead)
	if err != nil {
		return -1
	}
	return idx
}

func (bg *bucketGrid) forEachBucketIndex(aabb AABB, fn func(uint32)) {
	if bg == nil {
		return
	}
	tlIdx := bg.CalculateGridIndex(aabb.TopLeft)
	brIdx := bg.CalculateGridIndex(aabb.BottomRight)
	if tlIdx < 0 || brIdx < 0 || tlIdx >= len(bg.buckets) || brIdx >= len(bg.buckets) {
		return
	}

	if tlIdx == brIdx {
		fn(uint32(tlIdx))
		return
	}

	x1, y1 := bg.gridCellCodec.Decode(tlIdx)
	x2, y2 := bg.gridCellCodec.Decode(brIdx)
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			idx, err := bg.gridCellCodec.Encode(x, y)
			if err != nil || idx < 0 || idx >= len(bg.buckets) {
				continue
			}
			fn(uint32(idx))
		}
	}
}

// -----------------------------------------------------------
// bucket

func (b *bucket) Add(id EntryId, initialCap int) bool {
	if b.ids == nil {
		b.ids = make([]EntryId, 0, initialCap)
	} else if slices.Contains(b.ids, id) {
		return false
	}
	b.ids = append(b.ids, id)
	return true
}

func (b *bucket) Remove(id EntryId) bool {
	if b.ids == nil {
		return false
	}
	for i, existingId := range b.ids {
		if existingId == id {
			// Swap & Pop
			lastIdx := len(b.ids) - 1
			b.ids[i] = b.ids[lastIdx]
			b.ids = b.ids[:lastIdx]
			return true
		}
	}
	return false
}

// -----------------------------------------------------------
// memoryOptimizer

func newMemoryOptimizer(cells int) *memoryOptimizer {
	return &memoryOptimizer{
		dirtyQueue: make([]int, 0, 128),
		isDirty:    make([]bool, cells),
	}
}

func (mo *memoryOptimizer) mark(idx int, isEmpty bool) {
	if isEmpty && !mo.isDirty[idx] {
		mo.isDirty[idx] = true
		mo.dirtyQueue = append(mo.dirtyQueue, idx)
	}
}

func (mo *memoryOptimizer) collect(buckets []bucket) {
	for _, idx := range mo.dirtyQueue {
		mo.isDirty[idx] = false
		if len(buckets[idx].ids) == 0 {
			buckets[idx].ids = nil
		}
	}
	mo.dirtyQueue = mo.dirtyQueue[:0]
}

func (mo *memoryOptimizer) clear() {
	for i := range mo.isDirty {
		mo.isDirty[i] = false
	}
	mo.dirtyQueue = mo.dirtyQueue[:0]
}
