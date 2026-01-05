package spatial

import "github.com/kjkrol/gokg/pkg/plane"

// toroidGridHelper

func (h *toroidGridHelper) VisitWrappedAABB(aabb AABB, visit func(AABB)) {
	if h.m == nil || h.m.space == nil {
		return
	}
	wrapped := h.m.space.WrapAABB(aabb)
	if idxAABB, ok := h.m.indexAABB(wrapped.AABB); ok {
		visit(idxAABB)
	}
	wrapped.VisitFragments(func(_ plane.FragPosition, frag AABB) bool {
		if idxAABB, ok := h.m.indexAABB(frag); ok {
			visit(idxAABB)
		}
		return true
	})
}

func (h *toroidGridHelper) BuildFragments(shape AABB) (uint8, [4]AABB) {
	var frags [4]AABB
	mask := uint8(0)
	if h.m == nil || h.m.space == nil {
		return mask, frags
	}
	wrapped := h.m.space.WrapAABB(shape)
	if base, ok := h.m.indexAABB(wrapped.AABB); ok {
		frags[0] = base
		mask |= 1
	}
	wrapped.VisitFragments(func(pos plane.FragPosition, aabb AABB) bool {
		idx := uint8(pos) + 1
		if frag, ok := h.m.indexAABB(aabb); ok {
			frags[idx] = frag
			mask |= 1 << idx
		}
		return true
	})
	return mask, frags
}

// standardGridHelper

func (h *standardGridHelper) VisitWrappedAABB(aabb AABB, visit func(AABB)) {
	if h.m == nil {
		return
	}
	if idxAABB, ok := h.m.indexAABB(aabb); ok {
		visit(idxAABB)
	}
}

func (h *standardGridHelper) BuildFragments(shape AABB) (uint8, [4]AABB) {
	var frags [4]AABB
	if h.m == nil {
		return 0, frags
	}
	if base, ok := h.m.indexAABB(shape); ok {
		frags[0] = base
		return 1, frags
	}
	return 0, frags
}
