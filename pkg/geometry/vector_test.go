package geometry_test

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geometry"
)

func TestVecSub(t *testing.T) {
	for _, test := range []struct {
		arg1     geometry.Vec[int]
		arg2     geometry.Vec[int]
		expected geometry.Vec[int]
	}{
		{geometry.Vec[int]{2, 3}, geometry.Vec[int]{1, 2}, geometry.Vec[int]{1, 1}},
		{geometry.Vec[int]{1, 2}, geometry.Vec[int]{1, 2}, geometry.Vec[int]{0, 0}},
		{geometry.Vec[int]{0, 0}, geometry.Vec[int]{4, 4}, geometry.Vec[int]{-4, -4}},
		{geometry.Vec[int]{4, 0}, geometry.Vec[int]{1, 0}, geometry.Vec[int]{3, 0}},
		{geometry.Vec[int]{1, 0}, geometry.Vec[int]{4, 0}, geometry.Vec[int]{-3, 0}},
	} {
		if output := test.arg1.Sub(test.arg2); !output.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", test.arg1, test.expected)
		}
	}
}
