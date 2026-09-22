// Package plane is a box as a Space holds it: a [geom.AABB] that knows its own size and how far it
// runs past the far edges of a wrapping world, so that its wrapped pieces follow from it. It is the
// public face of the geometry a Space normalises; the normalising itself is done by the Space.
//
// # AABB
//
// An [AABB] embeds a [geom.AABB] and adds a Size, which survives clamping at an edge, and an
// Overhang: how far the box reaches past the right and bottom edge of the world once its top-left
// corner has been folded inside. A box with no Overhang lies whole in the world. [NewAABB] builds
// one not yet folded by any Space; Space.WrapAABB, Space.Place and Space.Move of package aabbworld
// fold it.
//
// # Fragments
//
// A box straddling a wrapping seam is its main piece plus up to three wrapped pieces, named by
// [FragPosition]: [FRAG_RIGHT] for what came back through the left edge, [FRAG_BOTTOM] for what came
// back through the top, [FRAG_BOTTOM_RIGHT] for the corner that did both. [AABB.VisitFragments]
// hands each wrapped piece to a [FragVisitor] until it says stop; the main piece is the embedded
// geom.AABB itself.
//
// # Overlap
//
// [AABB.DeepestOverlapWith] compares every piece of one box with every piece of another and returns
// the deepest [Overlap]: the two pieces that meet and the penetration between them.
package plane
