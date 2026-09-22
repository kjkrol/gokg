# Changelog

## v1.5.0 — 2026-09-22

The public API settled since v1.4.0 and is documented end to end.

**Collisions**
- One collision engine: `Space.CollideEngine(handler, collide.Config)` returns a `collide.Engine`
  whose `Tick` pairs every two `CanCollide` boxes within `Reach`, asks the `Handler` to confirm each
  overlap (`Touch`), separates confirmed pairs over `Iterations` passes (`Contact`), reports each
  pushed box once (`Moved`) and lists who left through an open edge (`Left`).
- Simpler solver: a pair is measured again only when one of its boxes has moved.

**Spatial index**
- `Rebuild` is lazy: it records the items and the id map; cells are built at the first `Query`.
  Several rebuilds in a tick cost one indexing.
- The engine invalidates the cells after pushing boxes, so `Query` and `Scan` see the pushed boxes
  without another `Rebuild`.
- Pairs are found by position in the item slice; the engine no longer looks pairs up by id
  (collision tick 2–18% faster on the broad-phase scenes).
- Merging the query cells and the pair sweep into one layout was measured and rejected; the
  reasoning is in [BENCHMARKS.md](BENCHMARKS.md#key-takeaways).

**Project**
- Every package has a `doc.go`; the root one carries the concepts, the tick lifecycle and the
  package graph.
- All benchmarks moved to `bench/` (one `bench_test` package); `Makefile` with `test`, `bench`,
  `bench-save`; results and method in `BENCHMARKS.md`.
- README rewritten around what the library is; the boundary-handling section moved under *Edges*.

## v1.4.0

New public API: `Space`, `Item`, `Capability`, `Edges`, `Config`, `View`, `Cone`; the earlier
`Space2D` implementations became the internal surface behind a `Space`.
