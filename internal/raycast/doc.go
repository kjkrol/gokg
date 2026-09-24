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
// A [Cone] is a direction, a half-angle either side of it below π, a radius, how see-through each
// entity is and, when sight has heights, the eye's height, each entity's bottom and top and the
// ground's height at a point. The scan queries the index for the rectangles covering the cone,
// split at each wrapping seam, and keeps only the candidates the cone's two edges do not dismiss.
//
// # View
//
// [View.Scan] sweeps the candidates by angle: every angle where the reach of sight can change shape
// gets a sample, the nearest blocking candidate along it or where the budget runs out. What is
// behind a nearer box is in its shadow and never sampled. Each cast notes on the candidates it
// reaches how near it saw them. The samples are then read three ways: [View.Entities] lists every
// candidate some cast reached, nearest first; [View.Depths] resamples the reach at evenly spaced
// angles without touching what was seen; [View.Outline] traces the lit region as a fan of points
// from the observer. A View keeps its buffers between scans and serves one goroutine.
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
// around the eye is left as it was. A see-through box the ray enters within its budget is seen,
// like anything else the ray reaches.
//
// # Elevation
//
// With Cone.Elevation or Cone.Ground set, sight has heights. The eye is at Cone.Eye; an entity
// spans its bottom to its top (nil Elevation: every height); the ground lies at Cone.Ground of a
// point (nil: flat at 0) and is sampled every Cone.GroundStep along a ray, plus once at the radius.
// A target — the top of an entity's near face, or a ground sample — is in sight when the straight
// line from the eye to it stays above every nearer ground sample, passes through no nearer blocking
// box's band over the stretch it crosses, and the budget lasts: a see-through box charges 1/τ per
// unit only for the stretch the line spends inside its band, so a hawk looking down over a forest
// pays nothing for it. Entities are seen whether or not they block: a walker under a hawk is seen
// beside it. The reach along an angle is the farthest lit ground: the last ground sample in sight,
// or the foot of a box standing on the ground — one whose bottom is at the ground under its centre
// — when that foot is in sight; ground hidden behind a crest or in a wall's shadow is passed over
// and lower ground farther on may be lit again. A wall taller than the eye thus cuts the reach at
// its foot as on a plane; a wall lower than the eye is looked over and shades the ground only for
// the stretch its top hides; a hill hides the plain behind it from the lowland and not from a hawk.
// Over uneven ground the whole cone is sampled every two degrees; with Ground nil only each box's
// span is.
//
// # Files
//
// sweep.go is the angular sweep and the budgeted cast along one angle on a plane; elevation.go the
// cast with heights; shadow.go the arcs a box subtends and the wedge test that dismisses a box
// outside the cone; slab.go the ray-box entry and exit distances; space.go the search rectangles,
// the nearest wrapped image of a box and the distances on a torus.
package raycast
