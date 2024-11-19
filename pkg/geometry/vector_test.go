package geometry

import (
	"testing"
)

func TestVector2dReduce(t *testing.T) {
	for _, test := range []struct {
		arg1     Vec[int]
		arg2     Vec[int]
		expected Vec[int]
	}{
		{Vec[int]{2, 3}, Vec[int]{1, 2}, Vec[int]{1, 1}},
		{Vec[int]{1, 2}, Vec[int]{1, 2}, Vec[int]{0, 0}},
		{Vec[int]{0, 0}, Vec[int]{4, 4}, Vec[int]{-4, -4}},
		{Vec[int]{4, 0}, Vec[int]{1, 0}, Vec[int]{3, 0}},
		{Vec[int]{1, 0}, Vec[int]{4, 0}, Vec[int]{-3, 0}},
	} {
		if output := test.arg1.Sub(test.arg2); !output.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", test.arg1, test.expected)
		}
	}
}
