# Changelog

## v1.7.0 — 2026-09-24

Sight with heights: an eye, entities with a bottom and a top, ground that rises and falls.

**Sight**
- `Cone.Eye`, `Cone.Elevation func(id) (bottom, top float64)`, `Cone.Ground func(p geom.Vec)
  float64` and `Cone.GroundStep` give sight heights. An entity is seen when the straight line from
  the eye to its top stays above every nearer ground sample, passes through no nearer blocking
  entity's band, and the budget lasts — a see-through entity charges `1/τ` per unit only for the
  stretch the line spends inside its band. Entities are seen whether or not they block: a walker
  under a hawk is seen beside it. `Elevation` nil spans every entity over all heights, `Ground` nil
  is flat ground at 0, `GroundStep` 0 is a sixteenth of the radius; with both funcs nil the scan is
  the one from v1.6.0.
- The reach of an angle (`View.Depths`, `View.Outline`) is the farthest lit ground: the last ground
  sample in sight or the foot of a box standing on the ground. A wall taller than the eye cuts the
  reach at its foot as before; a wall lower than the eye is looked over; a hill hides the plain
  behind it from the lowland and not from a hawk. Over uneven ground the whole cone is sampled
  every two degrees.
- **Changed**: a see-through entity the ray enters within its budget is now listed by
  `View.Entities` — a forest looked into is seen. Casts note what they see on the candidates, so
  `Depths` no longer changes what `Entities` reports.
- `Scan` stays at 0 allocs/op; the scan without heights measures as in v1.6.0 (0–6% over six
  alternating runs, one row at p < 0.05). Against the same scan without heights, bands alone cost
  2.5–2.8×, a ground raster sampled every 32 of the radius 800 3.2–5.2×, at the default step
  2.3–4.1× — the price is the angles and ground points cast, see
  [BENCHMARKS.md](BENCHMARKS.md#sight-with-heights--benchmark_view_elevated).

## v1.6.0 — 2026-09-24

Sight through see-through entities.

**Sight**
- `Cone.Transparency func(id uid.UID64) float64` gives each entity a transparency τ. A ray has a
  budget of the cone's radius: an empty stretch costs its length, a stretch through an entity of
  τ in (0, 1] costs its length divided by τ, an entity at τ ≤ 0 blocks where it is met. A forest at
  τ = 0.5 shortens sight by its depth; a wall cuts it as before. Nil (the default) keeps every
  entity blocking, so existing callers see no change.
- See-through entities shorten the reach but are not listed by `View.Entities`; `View.Depths` and
  `View.Outline` show the shortened reach. The outline behind a see-through entity is sampled every
  two degrees; `Depths` is exact per angle.
- A see-through entity around the observer dims the whole cone from where the observer stands.
- `Scan` stays at 0 allocs/op; the opaque-only scan measures the same as in v1.5.0 (no significant
  difference over six alternating runs), a scene with three entities in ten see-through costs about
  8% more — see [BENCHMARKS.md](BENCHMARKS.md#sight--benchmark_view_).

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
