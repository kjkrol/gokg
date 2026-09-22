package collide_test

import (
	"github.com/kjkrol/aabbworld"
	"github.com/kjkrol/aabbworld/collide"
	"math/rand/v2"
	"testing"

	"github.com/kjkrol/uid"
)

// broadPhaseScene is a population and the bucket size a caller sizing its
// buckets from its largest entity would have picked for it.
type broadPhaseScene struct {
	name   string
	bucket uint32
	margin float64 // what a probe per entity reaches: the smallest side
	groups []sceneGroup
}

type sceneGroup struct {
	count int
	side  float64
}

var broadPhaseScenes = []broadPhaseScene{
	{name: "side=10,fill=20%", bucket: 32, margin: 10, groups: []sceneGroup{{8388, 10}}},
	{name: "side=10,fill=40%", bucket: 32, margin: 10, groups: []sceneGroup{{16777, 10}}},
	{name: "side=5,fill=20%", bucket: 16, margin: 5, groups: []sceneGroup{{33554, 5}}},
	{name: "crowd=8,few=100", bucket: 256, margin: 8, groups: []sceneGroup{{12000, 8}, {20, 100}}},
}

func (sc broadPhaseScene) build(b *testing.B) *crowd {
	b.Helper()
	space, err := aabbworld.NewSpace(aabbworld.Config{
		Width: 2048, Height: 2048, Edges: aabbworld.Torus,
		BucketSize: sc.bucket, BucketCapacity: 8,
	})
	if err != nil {
		b.Fatalf("NewSpace: %v", err)
	}
	c := &crowd{space: space}
	rng := rand.New(rand.NewPCG(41, 43))
	for _, g := range sc.groups {
		for range g.count {
			c.add(rng.Float64()*2048, rng.Float64()*2048, g.side, g.side, aabbworld.CanCollide)
		}
	}
	space.Flush(nil)
	return c
}

func BenchmarkBroadPhase(b *testing.B) {
	for _, sc := range broadPhaseScenes {
		b.Run(sc.name+"/probe-per-entity", func(b *testing.B) {
			c := sc.build(b)
			var self uid.UID64
			found := 0
			onFound := func(other uid.UID64) {
				if other.Index() > self.Index() {
					found++
				}
			}
			for b.Loop() {
				found = 0
				for i, id := range c.ids {
					self = id
					c.space.Query(grown(c.boxes[i], sc.margin), aabbworld.CanCollide, onFound)
				}
			}
			b.ReportMetric(float64(found), "pairs")
		})
		b.Run(sc.name+"/pairs", func(b *testing.B) {
			c := sc.build(b)
			found := 0
			var e collide.Engine
			onPair := func(_, _ uid.UID64) (collide.Body, collide.Body, bool) {
				found++
				return collide.Body{}, collide.Body{}, false
			}
			for b.Loop() {
				found = 0
				e.Tick(c.space, pairsReach, aabbworld.CanCollide, 0, onPair, nil, nil)
			}
			b.ReportMetric(float64(found), "pairs")
		})
	}
}
