package bucketgrid

import (
	"slices"

	"github.com/kjkrol/gokg/pkg/geom"
	"github.com/kjkrol/gokg/pkg/spatial"
)

type (
	Entry       = spatial.Entry
	EntriesMove = spatial.EntriesMove
	Resolution  = spatial.Resolution
	CellCodec   = spatial.CellCodec
	AABB        = geom.AABB[uint32]
	Vec         = geom.Vec[uint32]

	Bucket struct {
		AABB
		entriesIds []uint64
	}

	BucketGrid struct {
		gridResolution    Resolution
		bucketsResolution Resolution
		gridCellCodec     CellCodec
		buckets           []Bucket
		entriesById       map[uint64]Entry
		bounds            AABB
	}
)

func NewBucket(capacity uint32, aabb AABB) Bucket {
	return Bucket{
		AABB:       aabb,
		entriesIds: make([]uint64, 0, capacity),
	}
}

func (b *Bucket) Add(id uint64) bool {
	if slices.Contains(b.entriesIds, id) {
		return false
	}
	b.entriesIds = append(b.entriesIds, id)
	return true
}

func (b *Bucket) Remove(id uint64) bool {
	for i, existingId := range b.entriesIds {
		if existingId == id {
			// Swap & Pop
			lastIdx := len(b.entriesIds) - 1
			b.entriesIds[i] = b.entriesIds[lastIdx]
			b.entriesIds = b.entriesIds[:lastIdx]
			return true
		}
	}
	return false
}

var _ spatial.Index = (*BucketGrid)(nil)

func NewBucketGrid(
	overallResolution Resolution,
	bucketsResolution Resolution,
	capacityFactor float64,
) spatial.Index {
	gridResolution := spatial.NewResolution(uint8(overallResolution - bucketsResolution))
	gridCellCodec := spatial.NewLinearCodec(gridResolution)
	buckets := make([]Bucket, gridResolution.Cells())

	cellsCount := float64(bucketsResolution.Cells())
	bucketCapacity := uint32(cellsCount * capacityFactor)

	initBuckets(buckets, bucketCapacity, bucketsResolution, gridResolution, gridCellCodec)
	entriesById := make(map[uint64]Entry, bucketCapacity*uint32(bucketsResolution.Cells()))
	bounds := geom.NewAABBAt(
		geom.NewVec[uint32](0, 0),
		overallResolution.Side(),
		overallResolution.Side(),
	)
	return &BucketGrid{
		gridResolution:    gridResolution,
		bucketsResolution: bucketsResolution,
		gridCellCodec:     gridCellCodec,
		buckets:           buckets,
		entriesById:       entriesById,
		bounds:            bounds,
	}
}

func initBuckets(
	buckets []Bucket,
	bucketCapacity uint32,
	bucketsResolution Resolution,
	gridResolution Resolution,
	gridCellCodec CellCodec,
) {
	for y := range gridResolution.Side() {
		for x := range gridResolution.Side() {
			index := gridCellCodec.Encode(x, y)
			topLeft := geom.NewVec(x*bucketsResolution.Side(), y*bucketsResolution.Side())
			bottomRight := topLeft.Add(geom.NewVec(bucketsResolution.Side(), bucketsResolution.Side()))
			aabb := AABB(geom.NewAABB(topLeft, bottomRight))
			buckets[index] = NewBucket(bucketCapacity, aabb)
		}
	}
}

func (bg *BucketGrid) BulkInsert(entries []Entry) {
	for _, entry := range entries {
		bg.insert(entry)
		bg.entriesById[entry.Id] = entry
	}
}

func (bg *BucketGrid) insert(entry Entry) {
	tlGridIndex := bg.calculateGridIndex(entry.AABB.TopLeft)
	brGridIndex := bg.calculateGridIndex(entry.AABB.BottomRight)
	if brGridIndex != tlGridIndex {
		bg.traverseRange(tlGridIndex, brGridIndex, func(i uint64) { bg.buckets[i].Add(entry.Id) })
	} else {
		bg.buckets[tlGridIndex].Add(entry.Id)
	}
}

// BulkRemove – remove whatever is stored at the given positions.
func (bg *BucketGrid) BulkRemove(entries []Entry) {
	for _, entry := range entries {
		bg.remove(entry)
		delete(bg.entriesById, entry.Id)
	}
}

func (bg *BucketGrid) remove(entry Entry) {
	tlGridIndex := bg.calculateGridIndex(entry.AABB.TopLeft)
	brGridIndex := bg.calculateGridIndex(entry.AABB.BottomRight)
	if brGridIndex != tlGridIndex {
		bg.traverseRange(tlGridIndex, brGridIndex, func(i uint64) { bg.buckets[i].Remove(entry.Id) })
	} else {
		bg.buckets[tlGridIndex].Remove(entry.Id)
	}
}

// BulkMove – update objects (typically same Value, different XY).
func (bg *BucketGrid) BulkMove(moves EntriesMove) {
	for i := range moves.New {
		oldEntry := moves.Old[i]
		newEntry := moves.New[i]
		if oldEntry == newEntry {
			continue
		}
		oldTlGridIndex := bg.calculateGridIndex(oldEntry.AABB.TopLeft)
		oldBrGridIndex := bg.calculateGridIndex(oldEntry.AABB.BottomRight)
		newTlGridIndex := bg.calculateGridIndex(newEntry.AABB.TopLeft)
		newBrGridIndex := bg.calculateGridIndex(newEntry.AABB.BottomRight)
		if oldTlGridIndex != newTlGridIndex || oldBrGridIndex != newBrGridIndex {
			bg.remove(oldEntry)
			bg.insert(newEntry)
		}
		bg.entriesById[newEntry.Id] = newEntry
	}
}

// QueryRange – all objects within the AABB.
func (bg *BucketGrid) QueryRange(aabb AABB, collector func(uint64)) int {
	tlGridIndex := bg.calculateGridIndex(aabb.TopLeft)
	brGridIndex := bg.calculateGridIndex(aabb.BottomRight)
	counter := 0
	if brGridIndex != tlGridIndex {
		bg.traverseRange(tlGridIndex, brGridIndex, func(i uint64) { counter += bg.query(i, aabb, collector) })
	} else {
		return bg.query(tlGridIndex, aabb, collector)
	}
	return counter
}

func (bg *BucketGrid) query(index uint64, aabb AABB, collector func(uint64)) int {
	counter := 0
	bucket := bg.buckets[index]
	for _, entryId := range bucket.entriesIds {
		entry, ok := bg.entriesById[entryId]
		if !ok {
			continue
		}
		if aabb.Intersects(entry.AABB) {
			collector(entryId)
			counter++
		}
	}
	return counter
}

// Count – number of objects in the structure.
func (bg *BucketGrid) Count() int { return len(bg.entriesById) }

// Bounds – global bounds of the handled space.
func (bg *BucketGrid) Bounds() AABB { return bg.bounds }

func (bg *BucketGrid) calculateGridIndex(vec Vec) uint64 {
	xHead := vec.X >> bg.bucketsResolution
	yHead := vec.Y >> bg.bucketsResolution
	return bg.gridCellCodec.Encode(xHead, yHead)
}

func (bg *BucketGrid) traverseRange(from, to uint64, fn func(uint64)) {
	x1, y1 := bg.gridCellCodec.Decode(from)
	x2, y2 := bg.gridCellCodec.Decode(to)
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			i := bg.gridCellCodec.Encode(x, y)
			fn(i)
		}
	}
}
