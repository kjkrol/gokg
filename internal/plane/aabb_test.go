package plane

import (
	"github.com/kjkrol/aabbworld/plane"
	"testing"
)

func TestAABB_NewAABB(t *testing.T) {
	aabb := plane.NewAABB(vec(0, 0), 10, 10)
	expected := vec(10, 10)
	if aabb.BottomRight != expected {
		t.Errorf("center %v not equal to expected %v", aabb.BottomRight, expected)
	}
}
