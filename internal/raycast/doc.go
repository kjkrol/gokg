// Package raycast is the state and the algorithm behind the public View of package aabbworld:
// what one observer sees through a cone, worked out over any index that can be queried by
// rectangle.
//
// # QueryableSpace
//
// [QueryableSpace] is the slice of a spatial index a scan needs: a range query, the box of one
// entity by id, and the world's size and which axes wrap. The public Space satisfies it through an
// adapter; the tests use a fake.
//
// # Cone
//
// A [Cone] is a direction, a half-angle either side of it below π, a radius, and how see-through
// each entity is. The scan queries the index for the rectangles covering the cone, split at each
// wrapping seam, and keeps only the candidates the cone's two edges do not dismiss.
//
// # View
//
// [View.Scan] sweeps the candidates by angle: every angle where the reach of sight can change shape
// gets a sample, the nearest blocking candidate along it or where the budget runs out. What is
// behind a nearer box is in its shadow and never sampled. The samples are then read three ways:
// [View.Entities] lists what was hit, nearest first; [View.Depths] resamples the reach at evenly
// spaced angles; [View.Outline] traces the lit region as a fan of points from the observer. A View
// keeps its buffers between scans and serves one goroutine.
//
// # Budget
//
// A ray starts with the radius as its budget and pays one per unit of empty way. A see-through
// candidate (transparency τ in (0, 1]) is not a hit: the way through it costs 1/τ per unit, so the
// ray gets less far behind it. The cast along one angle gathers the see-through boxes it crosses,
// sorted by entry, walks them charging each stretch once, and stops at the first blocking box if the
// budget lasts that far. The reach behind a see-through box bends with the way through it, so its
// span is sampled every two degrees on top of its edges; Depths, cast per angle, is exact. A
// see-through box around the eye dims the whole cone from where the eye stands; a blocking one
// around the eye is left as it was.
//
// # Files
//
// sweep.go is the angular sweep and the budgeted cast along one angle; shadow.go the arcs a box
// subtends and the wedge test that dismisses a box outside the cone; slab.go the ray-box entry and
// exit distances; space.go the search rectangles, the nearest wrapped image of a box and the
// distances on a torus.
package raycast
