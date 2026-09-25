package collide

import (
	pcollide "github.com/kjkrol/aabbworld/collide"
	"github.com/kjkrol/aabbworld/geom"
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/aabbworld/internal/spatial"
)

// Field is how a solve meets the solid ground: the boxes round an item, and whom to tell of them.
type Field struct {
	Solid   pcollide.SolidField
	Touch   func(item int32, cell uint64, pen geom.Vec) (geom.Vec, bool)
	Contact func(item int32, cell uint64, pen geom.Vec)
}

// fieldKey is one item against one cell of the ground.
type fieldKey struct {
	item int32
	cell uint64
}

// What a tick has settled of an item against a cell.
const (
	fieldReported uint8 = 1 + iota
	fieldVetoed
)

// fieldPass is the state of pushing one item out of the ground.
type fieldPass struct {
	field   *Field
	surface *iplane.Surface
	item    int32
	it      *spatial.Item
	sensor  bool
	pushed  bool
	met     map[fieldKey]uint8
	visit   func(pcollide.FieldBox) bool
}

// meet pushes the item out of one solid box it overlaps, through the shallowest open side.
func (p *fieldPass) meet(fb pcollide.FieldBox) bool {
	box := &p.it.Box
	a := geom.AABB{TopLeft: box.TopLeft, BottomRight: box.TopLeft.Add(box.Size)}
	b := fb.Box
	if a.BottomRight.X <= b.TopLeft.X || a.TopLeft.X >= b.BottomRight.X || a.BottomRight.Y <= b.TopLeft.Y || a.TopLeft.Y >= b.BottomRight.Y {
		return true
	}
	push := pushOut(a, b, fb.Open)
	key := fieldKey{p.item, fb.Cell}
	switch p.met[key] {
	case fieldVetoed:
		return true
	case 0:
		if p.field.Touch != nil {
			pen, ok := p.field.Touch(p.item, fb.Cell, push)
			if !ok {
				p.met[key] = fieldVetoed
				return true
			}
			push = pen
		}
		p.met[key] = fieldReported
		if p.field.Contact != nil {
			p.field.Contact(p.item, fb.Cell, push)
		}
	}
	if p.sensor || push == (geom.Vec{}) {
		return true
	}
	p.surface.Translate(box, push)
	p.pushed = true
	return true
}

// pushOut is the shortest way out of b for a through one of b's open sides; with none open, the
// shortest way out at all.
func pushOut(a, b geom.AABB, open pcollide.Sides) geom.Vec {
	ways := [4]struct {
		side pcollide.Sides
		v    geom.Vec
		d    float64
	}{
		{pcollide.Left, geom.NewVec(b.TopLeft.X-a.BottomRight.X, 0), a.BottomRight.X - b.TopLeft.X},
		{pcollide.Right, geom.NewVec(b.BottomRight.X-a.TopLeft.X, 0), b.BottomRight.X - a.TopLeft.X},
		{pcollide.Top, geom.NewVec(0, b.TopLeft.Y-a.BottomRight.Y), a.BottomRight.Y - b.TopLeft.Y},
		{pcollide.Bottom, geom.NewVec(0, b.BottomRight.Y-a.TopLeft.Y), b.BottomRight.Y - a.TopLeft.Y},
	}
	best := -1
	for i, w := range ways {
		if open&w.side != 0 && (best < 0 || w.d < ways[best].d) {
			best = i
		}
	}
	if best < 0 {
		for i, w := range ways {
			if best < 0 || w.d < ways[best].d {
				best = i
			}
		}
	}
	return ways[best].v
}

// movable reports whether an item takes part in collisions and something may shift it.
func movable(it *spatial.Item) bool {
	return it.Caps&spatial.CanCollide != 0 && it.Caps&spatial.Static == 0
}
