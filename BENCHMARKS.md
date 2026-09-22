# ⏱️ aabbworld Benchmarks

[← Back to README](./README.md)

> What each benchmark measures, the numbers, and what they taught us about the design.

## Environment

- **CPU:** Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz (4 cores, 8 threads)
- **Go version:** 1.27.1
- **OS:** Linux

Every number below is the **median of 5 runs** (`make bench-save`, `-count=5`) at commit `d6d935d`,
2026-09-22. This machine drifts by up to ~10% between runs, so a change is judged by running the
old and the new tree **alternately**, never one run each; see [How to benchmark](#how-to-benchmark).

All benchmarks live in [`bench/`](bench/) and drive the library through its importable API only.
"0 allocs/op" throughout means a steady tick allocates nothing; the non-zero B/op on the crowd
benchmarks is the amortised growth of buffers the grid, the solver and the view keep between calls.

## Broad phase — `Benchmark_BroadPhase`

Two ways of finding who is near whom in a crowd of `CanCollide` boxes on a 2048×2048 torus:

- **probe-per-entity**: one `Space.Query` per entity over its box grown by the reach; the
  straightforward way to do a broad phase from outside.
- **tick**: one `Engine.Tick` whose `Touch` refuses every pair, so it measures pairing plus one
  overlap test per pair and nothing moves.

The boxes are scattered at random, so some overlap. "Covering 20%" means the boxes' total area is
20% of the world's: 8,388 boxes of 10×10 cover 838,800 of the torus's 4,194,304 square units. The
scenes are sized to hold that coverage constant across box sizes, so the same crowding is measured
with different boxes. Bucket size is what a caller sizing buckets from its largest entity would pick.

| Scene | bucket | probe-per-entity | pairs found | tick | overlaps found |
|:---|---:|---:|---:|---:|---:|
| 8,388 boxes, 10×10 each, covering 20% of the world | 32 | 2.27 ms | 13,471 | **1.78 ms** | 3,342 |
| 16,777 boxes, 10×10 each, covering 40% | 32 | 7.25 ms | 53,723 | **6.20 ms** | 13,398 |
| 33,554 boxes, 5×5 each, covering 20% | 16 | 11.46 ms | 53,866 | **8.61 ms** | 13,372 |
| 12,000 boxes 8×8 plus 20 boxes 100×100 | 256 | 14.69 ms | 18,571 | **7.42 ms** | 5,005 |

The tick wins everywhere, and by most in the mixed-size crowd: its sweep places each box by its own
grown size, while a per-entity probe over 256-unit buckets scans every neighbour of every 100-box.

## Spatial grid — `Benchmark_Grid_RebuildPairsQuery`

One tick of the index alone, on a 1024×1024 torus with 16-unit cells: every box moves a little,
`Rebuild`, `Pairs` at reach 0.5 over the two thirds of items that collide, and eight 100×100
`Query` probes.

| Scene | per tick |
|:---|---:|
| 2,000 boxes, 5×5 each | 242 µs |
| 8,388 boxes, 5×5 each | 1.23 ms |
| 8,388 boxes of mixed sizes: 4×4, 5×5, 16×16 and 60×60 | 8.48 ms |

Profile of the 8,388 boxes 5×5 case: building the query cells ≈31%, building the pair sweep ≈24%,
testing pairs ≈32%, the eight queries ≈2%.

## Collision solver — `Benchmark_Solver_*`

The solver alone, fed pre-computed pairs.

| Benchmark | Scene | per solve |
|:---|:---|---:|
| `Solver_Field` | a 100×50 lattice of 10×10 boxes, 9,850 neighbour pairs, 1 box in 25 nudged into its neighbour | 221 µs |
| `Solver_Dense` | 8,388 boxes, 5×5 each, crowded on a 1024×1024 torus: 13,547 candidate pairs, 3,763 real contacts | 4.25 ms |

## Sight — `Benchmark_View_*`

An observer on a 4000×4000 plane with 256-unit buckets scans a 90° cone of radius 800.

| Read | 100 entities | 1,000 entities |
|:---|---:|---:|
| `Scan` + `Entities` (nearest first) | 2.88 µs | 24.8 µs |
| `Scan` + `Outline` (lit region as a fan) | 4.09 µs | 26.4 µs |
| `Scan` + both | 4.17 µs | 28.0 µs |
| `Scan` + `Depths` (63 angles) | 4.78 µs | 26.8 µs |

The scan dominates; each read of the samples adds a microsecond or two.

## Primitives — `Benchmark_AABB_*`, `Benchmark_Surface_*`, `Benchmark_Vec_*`

| Operation | euclidean | toroidal |
|:---|---:|---:|
| `Surface.WrapAABB` (box across the far corner) | 25.7 ns | 30.6 ns |
| `Surface.Expand` (grow by 5 near the far corner) | 11.2 ns | 16.6 ns |
| `Surface.Translate` (across the far edge and back) | 13.5 ns | 17.7 ns |

| Operation | ns |
|:---|---:|
| `AABB.Contains`, `AABB.Intersects` | 1.43 |
| `Clamp` a point | 1.50 |
| `Wrap` a point: already inside / one step over / just below zero | 5.5 / 6.6 / 7.1 |
| `Wrap` a point from far out (99,999) | 73 |

## Key takeaways

* **`Rebuild` is lazy, and that is free.** It records the slice and the id map and leaves cell
  building to the first `Query`. The check is one predictable branch per query; the engine
  invalidates the cells after pushing boxes, so queries see the pushes without another rebuild.
* **The grid keeps two layouts on purpose.** Query cells hold every item at its actual size; the
  pair sweep holds only the items the engine wants, grown by their reach. Merging them into one
  reach-grown layout was implemented and measured against a baseline copy, alternately:
  +10…+30% on the grid tick, +4…+13% on the collision tick, +12…+39% on per-entity probes.
  The shared part (a counting sort) is cheap; what costs is per-piece geometry, which a merge does
  not remove, while both readers then scan a superset (pair entries 45k → 90k in the mixed-size
  scene). Two specialised layouts beat one general one.
* **Pairs are reported by position, not id.** The engine used to look each pair's two ids up in the
  grid; reporting positions in the item slice removed that and shrank a sweep entry from 48 to 40
  bytes, for a 2–18% faster collision tick across the broad-phase scenes.
* **Zero allocations on the tick.** Rebuild, Pairs, Solve, Query and Scan all report 0 allocs/op once
  their buffers have grown to the scene.

## How to benchmark

```bash
make bench        # the whole suite once, with allocations
make bench-save   # 5 repeats, raw output under bench_results/ (ignored by git)
```

To judge a change on this or any drifting machine, keep a copy of the baseline tree and run the two
alternately, then compare medians:

```bash
cp -r . /tmp/before          # before the change
for i in 1 2; do
  (cd /tmp/before && go test -run xxx -bench . -count 4 ./bench/... | sed 's/^/BEFORE /')
  go test -run xxx -bench . -count 4 ./bench/... | sed 's/^/AFTER  /'
done
```
