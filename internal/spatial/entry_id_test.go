package spatial

import (
	"testing"

	"github.com/kjkrol/uid"
)

func TestWithFrag_RoundTrip(t *testing.T) {
	var pool uid.UID64Pool
	pool.Init(8, 4)
	id := pool.Next()

	for frag := uint8(0); frag < 4; frag++ {
		tagged := withFrag(id, frag)

		if got := fragOf(tagged); got != frag {
			t.Errorf("fragOf(withFrag(id, %d)) = %d, want %d", frag, got, frag)
		}
		if got := withoutFrag(tagged); got != id {
			t.Errorf("withoutFrag(withFrag(id, %d)) = %v, want %v", frag, got, id)
		}
	}
}

func TestWithFrag_DistinctPerFragment(t *testing.T) {
	var pool uid.UID64Pool
	pool.Init(8, 4)
	id := pool.Next()

	seen := make(map[uid.UID64]bool)
	for frag := uint8(0); frag < 4; frag++ {
		tagged := withFrag(id, frag)
		if seen[tagged] {
			t.Fatalf("withFrag(id, %d) collided with a previous fragment's tagged value", frag)
		}
		seen[tagged] = true
	}
}

func TestWithFrag_DifferentGenerationsDiffer(t *testing.T) {
	var pool uid.UID64Pool
	pool.Init(8, 4)
	id := pool.Next()
	pool.Release(id)
	reissued := pool.Next()

	if withFrag(id, 0) == withFrag(reissued, 0) {
		t.Error("expected different generations of the same index to produce different tagged values")
	}
}
