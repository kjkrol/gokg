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
// A [Cone] is a direction, a half-angle either side of it below π, and a radius. The scan queries
// the index for the rectangles covering the cone, split at each wrapping seam, and keeps only the
// candidates the cone's two edges do not dismiss.
//
// # View
//
// [View.Scan] sweeps the candidates by angle: every angle where the reach of sight can change shape
// gets a sample, the nearest candidate along it or the radius. What is behind a nearer box is in
// its shadow and never sampled. The samples are then read three ways: [View.Entities] lists what was
// hit, nearest first; [View.Depths] resamples the reach at evenly spaced angles; [View.Outline]
// traces the lit region as a fan of points from the observer. A View keeps its buffers between
// scans and serves one goroutine.
//
// # Files
//
// sweep.go is the angular sweep and the cast along one angle; shadow.go the arcs a box subtends and
// the wedge test that dismisses a box outside the cone; slab.go the ray-box hit distance; space.go
// the search rectangles, the nearest wrapped image of a box and the distances on a torus.
package raycast
