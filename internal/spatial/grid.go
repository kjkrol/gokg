package spatial

import (
	"github.com/kjkrol/aabbworld/geom"
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

// Item is one box a Grid is told of at Rebuild.
type Item struct {
	ID   uid.UID64
	Box  plane.AABB
	Caps Capability
}

// Grid is a cell grid over the items of its last Rebuild. Pair sweeps read the items as they are;
// the cells are built from them at the first Query after a Rebuild or an Invalidate.
type Grid struct {
	surface  *iplane.Surface
	cellRes  Resolution
	side     uint32
	maxCoord uint32

	items    []Item
	rebuilds uint32
	slotOf   []int32
	stamps   []uint32
	stale    bool

	pieces []gridPiece
	starts []uint32
	cursor []uint32
	cells  []gridCell
	sweep  pairSweep

	listing int32
	onImage plane.FragVisitor
}

// gridPiece is one image of an item as the grid holds it, and the cells it covers.
type gridPiece struct {
	gridCell
	x1, y1, x2, y2 uint32
}

type gridCell struct {
	box  geom.AABB
	at   int32
	caps Capability
}

// NewGrid builds an empty grid over surface, worldRes on a side, cellRes per cell.
func NewGrid(surface *iplane.Surface, worldRes, cellRes Resolution) *Grid {
	g := &Grid{
		surface:  surface,
		cellRes:  cellRes,
		side:     NewResolution(uint8(worldRes - cellRes)).Side(),
		maxCoord: worldRes.MaxCoord(),
	}
	g.onImage = g.addImage
	g.Rebuild(nil)
	return g
}

// Rebuild replaces what the grid holds with items, which must stay put until the next Rebuild.
// Their boxes may change in the meantime; Invalidate tells the grid so. The cells are built at
// the first Query that needs them, so several Rebuilds in a row cost one indexing.
func (g *Grid) Rebuild(items []Item) {
	g.items = items
	g.rebuilds++
	for i := range items {
		g.remember(items[i].ID, int32(i))
	}
	g.stale = true
}

// Invalidate marks the cells as behind the boxes of the items, which have moved since they were built.
func (g *Grid) Invalidate() { g.stale = true }

// index builds the cells from the items as they are now.
func (g *Grid) index() {
	g.stale = false
	items := g.items
	g.pieces = g.pieces[:0]
	cells := int(g.side) * int(g.side)
	if cap(g.starts) < cells+1 {
		g.starts = make([]uint32, cells+1)
		g.cursor = make([]uint32, cells)
	}
	g.starts = g.starts[:cells+1]
	g.cursor = g.cursor[:cells]
	clear(g.starts)

	for i := range items {
		e := &items[i]
		if main, ok := g.clamp(e.Box.AABB); ok {
			g.add(main, int32(i), e.Caps)
		}
		if e.Box.Overhang != (geom.Vec{}) {
			g.listing = int32(i)
			e.Box.VisitFragments(g.onImage)
		}
	}

	for i := 1; i < len(g.starts); i++ {
		g.starts[i] += g.starts[i-1]
	}
	copy(g.cursor, g.starts)
	total := int(g.starts[cells])
	if cap(g.cells) < total {
		g.cells = make([]gridCell, total)
	}
	g.cells = g.cells[:total]
	for i := range g.pieces {
		p := &g.pieces[i]
		for y := p.y1; y <= p.y2; y++ {
			for x := p.x1; x <= p.x2; x++ {
				cell := y*g.side + x
				g.cells[g.cursor[cell]] = p.gridCell
				g.cursor[cell]++
			}
		}
	}
}

func (g *Grid) addImage(_ plane.FragPosition, image geom.AABB) bool {
	if frag, ok := g.clamp(image); ok {
		g.add(frag, g.listing, g.items[g.listing].Caps)
	}
	return true
}

func (g *Grid) add(box geom.AABB, at int32, caps Capability) {
	p := gridPiece{
		gridCell: gridCell{box: box, at: at, caps: caps},
		x1:       g.cellOf(box.TopLeft.X), y1: g.cellOf(box.TopLeft.Y),
		x2: g.cellOf(box.BottomRight.X), y2: g.cellOf(box.BottomRight.Y),
	}
	for y := p.y1; y <= p.y2; y++ {
		for x := p.x1; x <= p.x2; x++ {
			g.starts[y*g.side+x+1]++
		}
	}
	g.pieces = append(g.pieces, p)
}

// remember maps the entity id names to its item for this Rebuild.
func (g *Grid) remember(id uid.UID64, at int32) {
	i := int(id.Index())
	for i >= len(g.slotOf) {
		g.slotOf = append(g.slotOf, 0)
		g.stamps = append(g.stamps, 0)
	}
	g.slotOf[i], g.stamps[i] = at, g.rebuilds
}

// Item is the item the entity id names, if the last Rebuild was told of it.
func (g *Grid) Item(id uid.UID64) (*Item, bool) {
	i := int(id.Index())
	if i >= len(g.slotOf) || g.stamps[i] != g.rebuilds {
		return nil, false
	}
	e := &g.items[g.slotOf[i]]
	if e.ID != id {
		return nil, false
	}
	return e, true
}

// EntryAABB is the main box of the entity id names as the grid holds it.
func (g *Grid) EntryAABB(id uid.UID64) (geom.AABB, bool) {
	e, ok := g.Item(id)
	if !ok {
		return geom.AABB{}, false
	}
	return g.clamp(e.Box.AABB)
}

// Items is what the last Rebuild was told of, in order.
func (g *Grid) Items() []Item { return g.items }

// Query calls fn for every piece sharing want that intersects box, each once, and returns how many.
// It reads the boxes as they are now, building the cells first if they are behind.
func (g *Grid) Query(box geom.AABB, want Capability, fn func(uid.UID64)) int {
	area, ok := g.clamp(box)
	if !ok {
		return 0
	}
	if g.stale {
		g.index()
	}
	x1, y1 := g.cellOf(area.TopLeft.X), g.cellOf(area.TopLeft.Y)
	x2, y2 := g.cellOf(area.BottomRight.X), g.cellOf(area.BottomRight.Y)
	found := 0
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			cell := y*g.side + x
			for _, c := range g.cells[g.starts[cell]:g.starts[cell+1]] {
				if !c.caps.matches(want) || !area.Intersects(c.box) {
					continue
				}
				if max(g.cellOf(c.box.TopLeft.X), x1) != x || max(g.cellOf(c.box.TopLeft.Y), y1) != y {
					continue
				}
				fn(g.items[c.at].ID)
				found++
			}
		}
	}
	return found
}

// Pairs calls fn once, lower position first, for every two items sharing want whose boxes grown
// by reach of their shorter side touch. It reads the boxes as they are now.
func (g *Grid) Pairs(reach float64, want Capability, fn func(a, b int32)) {
	s := &g.sweep
	s.shift, s.side = g.cellRes, g.side
	s.listItems(g.surface, g.items, reach, want)
	s.sort()
	s.visit(fn)
	s.flushSeam(fn)
}

// cellOf is the column or row holding v, held inside the grid.
func (g *Grid) cellOf(v float64) uint32 {
	if !(v >= 0) {
		return 0
	}
	return min(uint32(v)>>g.cellRes, g.side-1)
}

// clamp holds box inside the grid; false when nothing of it is.
func (g *Grid) clamp(box geom.AABB) (geom.AABB, bool) {
	minX := clampToGrid(box.TopLeft.X, g.maxCoord)
	minY := clampToGrid(box.TopLeft.Y, g.maxCoord)
	maxX := clampToGrid(box.BottomRight.X, g.maxCoord)
	maxY := clampToGrid(box.BottomRight.Y, g.maxCoord)
	if maxX < minX || maxY < minY {
		return geom.AABB{}, false
	}
	return geom.NewAABB(geom.NewVec(minX, minY), geom.NewVec(maxX, maxY)), true
}

// clampToGrid holds a coordinate inside the indexable grid, NaN included.
func clampToGrid(val float64, max uint32) float64 {
	if !(val >= 0) {
		return 0
	}
	if val > float64(max) {
		return float64(max)
	}
	return val
}
