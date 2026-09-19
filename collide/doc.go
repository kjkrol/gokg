// Package collide separates overlapping boxes. It is given pairs that may be
// in contact and pushes them apart until they are not, under the boundary
// rules of the space they live in — so a pair meeting across a toroidal seam
// separates the short way, not across the whole world.
//
// The solver knows geometry and nothing else: no identifiers, no spatial
// index, no notion of mass or velocity. Pairs arrive as pointers to boxes the
// caller owns and the solver writes through them, so whatever the caller keeps
// those boxes in — an ECS chunk, a slice, a struct field — stays the caller's
// business. What a contact *means* is the caller's business too: the solver
// reports each one once and moves on.
//
// Finding which pairs are worth handing over is the spatial index's job, not
// this package's — see Space.Neighbours.
package collide
