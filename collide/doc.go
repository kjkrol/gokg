// Package collide finds and separates overlapping boxes of a Space: BroadPhase names the pairs
// that may touch, a NarrowPhase gathers them as Pairs, and Separate tests each exactly and
// pushes the overlapping apart.
package collide
