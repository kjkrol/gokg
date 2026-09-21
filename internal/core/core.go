// Package core is the inside of a Space, shared by the packages of this module that work on one.
package core

import (
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"

	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/aabbworld/internal/spatial"
)

// Space is the surface boxes move on and the index kept over them.
type Space struct {
	Index   *spatial.GridIndexManager
	Surface *iplane.Surface
}

// Follow queues the index to take box as where id is now; false once it has left by an open edge.
func (s *Space) Follow(id uid.UID64, box *plane.AABB) bool {
	if s.Surface.Left(box) {
		s.Index.QueueRemove(id)
		return false
	}
	s.Index.QueueUpdate(id, *box, true)
	return true
}

// Of returns the inside of an *aabbworld.Space; package aabbworld sets it.
var Of func(space any) *Space
