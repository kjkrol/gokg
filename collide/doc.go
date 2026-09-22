// Package collide is the contract of a Space's collision engine: what it is told ([Config]), whom it
// reports to ([Handler]), and what a caller drives ([Engine]). Space.CollideEngine of package
// aabbworld builds one over a Space's CanCollide boxes.
//
// # Handler
//
// One [Engine.Tick] talks to the Handler in a fixed order. Touch is asked once for every pair of
// boxes that overlap, with the penetration between them; the handler may confirm the pair as it is,
// refine the penetration (a round shape inside a square box touches later than the box does) or
// veto the pair for this tick. Contact is told each confirmed pair once, with the penetration it is
// being separated by. Moved is told every box the engine pushed, once, where it came to rest — the
// caller writes that box back to whatever owns it. A box pushed out through an open edge is not
// Moved; it is listed by [Engine.Left] until the next Tick.
//
// # Config
//
// Reach is how far a box may travel in one tick, as a fraction of its shorter side; boxes whose
// reaches do not touch are never paired, so it bounds the work of a tick. Iterations caps the passes
// one Tick spends on chained overlaps: separating one pair may push a box into a third, which the
// next pass resolves.
//
// # Engine
//
// An Engine serves one goroutine. Tick reads the boxes of its Space as they are now — the Space
// need only have been Rebuilt since the items changed — and leaves them pushed apart; the Space
// answers its next Query with the pushed boxes on its own.
package collide
