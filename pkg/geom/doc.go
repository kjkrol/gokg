// Package geom provides generic 2D geometry primitives shared across the GOK
// modules. It defines Vec[T] with numeric constraints plus vector operations
// (add, subtract, clamp, wrap) that are specialised per numeric kind via
// VectorMath. AABB supplies axis-aligned bounding boxes with containment,
// intersection, and quad-splitting helpers that higher-level packages wrap in
// plane-aware types.
package geom
