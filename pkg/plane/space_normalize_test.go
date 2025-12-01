package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestCartesian_normalizeBox(t *testing.T) {
	runCartesianNormalizeBoxTest[int](t, "int")
	runCartesianNormalizeBoxTest[uint32](t, "uint32")
	runCartesianNormalizeBoxTest[float64](t, "float64")
}

func runCartesianNormalizeBoxTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		cartesian := NewCartesian(T(10), T(10))

		for _, tc := range cartesianNormalizeScenarios[T]() {
			t.Run(tc.name, func(t *testing.T) {
				aabb := newAABB(vec[T](tc.topLeft.X, tc.topLeft.Y), T(tc.width), T(tc.height))
				cartesian.(space2d[T]).normalizeAABB(&aabb)
				expectAABBState(t, aabb,
					vec[T](tc.expectedTopLeft.X, tc.expectedTopLeft.Y),
					vec[T](tc.expectedBottomRight.X, tc.expectedBottomRight.Y),
					map[FragPosition][2]geom.Vec[T]{},
				)
			})
		}
	})
}

type normalizeScenario struct {
	name                string
	topLeft             geom.Vec[int]
	width, height       int
	expectedTopLeft     geom.Vec[int]
	expectedBottomRight geom.Vec[int]
}

func cartesianNormalizeScenarios[T geom.Numeric]() []normalizeScenario {
	return []normalizeScenario{
		{
			name:                "keeps_box_inside_viewport",
			topLeft:             geom.NewVec(1, 1),
			width:               2,
			height:              2,
			expectedTopLeft:     geom.NewVec(1, 1),
			expectedBottomRight: geom.NewVec(3, 3),
		},
		{
			name:                "clamps_bottom_right",
			topLeft:             geom.NewVec(9, 9),
			width:               2,
			height:              2,
			expectedTopLeft:     geom.NewVec(9, 9),
			expectedBottomRight: geom.NewVec(10, 10),
		},
		{
			name:                "clamps_box_outside_viewport",
			topLeft:             geom.NewVec(12, 12),
			width:               2,
			height:              2,
			expectedTopLeft:     geom.NewVec(10, 10),
			expectedBottomRight: geom.NewVec(10, 10),
		},
		{
			name:                "truncates_height_at_boundary",
			topLeft:             geom.NewVec(3, 9),
			width:               2,
			height:              2,
			expectedTopLeft:     geom.NewVec(3, 9),
			expectedBottomRight: geom.NewVec(5, 10),
		},
		{
			name:                "clamps_large_box",
			topLeft:             geom.NewVec(8, 8),
			width:               5,
			height:              5,
			expectedTopLeft:     geom.NewVec(8, 8),
			expectedBottomRight: geom.NewVec(10, 10),
		},
		{
			name:                "clamps_top_left_outside_viewport",
			topLeft:             geom.NewVec(-2, -2),
			width:               2,
			height:              2,
			expectedTopLeft:     geom.NewVec(0, 0),
			expectedBottomRight: geom.NewVec(0, 0),
		},
		{
			name:                "moves_negative_box_inside_viewport",
			topLeft:             geom.NewVec(-1, -1),
			width:               2,
			height:              2,
			expectedTopLeft:     geom.NewVec(0, 0),
			expectedBottomRight: geom.NewVec(1, 1),
		},
		{
			name:                "ignores_far_negative_box",
			topLeft:             geom.NewVec(-11, -11),
			width:               2,
			height:              2,
			expectedTopLeft:     geom.NewVec(0, 0),
			expectedBottomRight: geom.NewVec(0, 0),
		},
		{
			name:                "clamps_far_outside_viewport",
			topLeft:             geom.NewVec(19, 19),
			width:               2,
			height:              2,
			expectedTopLeft:     geom.NewVec(10, 10),
			expectedBottomRight: geom.NewVec(10, 10),
		}}
}

func TestTorus_normalizeBox(t *testing.T) {
	runTorusNormalizeBoxTest[int](t, "int")
	runTorusNormalizeBoxTest[uint32](t, "uint32")
	runTorusNormalizeBoxTest[float64](t, "float64")
}

func runTorusNormalizeBoxTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		torus := NewTorus(T(10), T(10))

		for _, tc := range torusNormalizeScenarios[T]() {
			t.Run(tc.name, func(t *testing.T) {
				aabb := newAABB(vec[T](tc.topLeft.X, tc.topLeft.Y), T(tc.width), T(tc.height))
				torus.(space2d[T]).normalizeAABB(&aabb)
				expectAABBState(t, aabb,
					vec[T](tc.expectedTopLeft.X, tc.expectedTopLeft.Y),
					vec[T](tc.expectedBottomRight.X, tc.expectedBottomRight.Y),
					convertFragments[T](tc.expectedFragments),
				)
			})
		}
	})
}

type torusNormalizeScenario struct {
	normalizeScenario
	expectedFragments map[FragPosition][2]geom.Vec[int]
}

func torusNormalizeScenarios[T geom.Numeric]() []torusNormalizeScenario {
	return []torusNormalizeScenario{
		{
			normalizeScenario: normalizeScenario{
				name:                "wraps_box_into_view",
				topLeft:             geom.NewVec(12, 12),
				width:               2,
				height:              2,
				expectedTopLeft:     geom.NewVec(2, 2),
				expectedBottomRight: geom.NewVec(4, 4),
			},
			expectedFragments: map[FragPosition][2]geom.Vec[int]{},
		},
		{
			normalizeScenario: normalizeScenario{
				name:                "fragments_on_edges",
				topLeft:             geom.NewVec(9, 9),
				width:               2,
				height:              2,
				expectedTopLeft:     geom.NewVec(9, 9),
				expectedBottomRight: geom.NewVec(10, 10),
			},
			expectedFragments: map[FragPosition][2]geom.Vec[int]{
				FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
				FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
				FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
			},
		},
		{
			normalizeScenario: normalizeScenario{
				name:                "wraps_far_outside",
				topLeft:             geom.NewVec(19, 19),
				width:               2,
				height:              2,
				expectedTopLeft:     geom.NewVec(9, 9),
				expectedBottomRight: geom.NewVec(10, 10),
			},
			expectedFragments: map[FragPosition][2]geom.Vec[int]{
				FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
				FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
				FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
			},
		},
		{
			normalizeScenario: normalizeScenario{
				name:                "fragments_on_bottom_edge",
				topLeft:             geom.NewVec(0, 9),
				width:               2,
				height:              2,
				expectedTopLeft:     geom.NewVec(0, 9),
				expectedBottomRight: geom.NewVec(2, 10),
			},
			expectedFragments: map[FragPosition][2]geom.Vec[int]{
				FRAG_BOTTOM: {geom.NewVec(0, 0), geom.NewVec(2, 1)},
			},
		},
		{
			normalizeScenario: normalizeScenario{
				name:                "wraps_negative_box_into_view",
				topLeft:             geom.NewVec(-2, -2),
				width:               2,
				height:              2,
				expectedTopLeft:     geom.NewVec(8, 8),
				expectedBottomRight: geom.NewVec(10, 10),
			},
			expectedFragments: map[FragPosition][2]geom.Vec[int]{},
		},
		{
			normalizeScenario: normalizeScenario{
				name:                "wraps_from_negative_corner_with_fragments",
				topLeft:             geom.NewVec(-1, -1),
				width:               2,
				height:              2,
				expectedTopLeft:     geom.NewVec(9, 9),
				expectedBottomRight: geom.NewVec(10, 10),
			},
			expectedFragments: map[FragPosition][2]geom.Vec[int]{
				FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
				FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
				FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
			},
		},
		{
			normalizeScenario: normalizeScenario{
				name:                "wraps_far_negative_box",
				topLeft:             geom.NewVec(-11, -11),
				width:               2,
				height:              2,
				expectedTopLeft:     geom.NewVec(9, 9),
				expectedBottomRight: geom.NewVec(10, 10),
			},
			expectedFragments: map[FragPosition][2]geom.Vec[int]{
				FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
				FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
				FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
			},
		}}
}

func convertFragments[T geom.Numeric](frags map[FragPosition][2]geom.Vec[int]) map[FragPosition][2]geom.Vec[T] {
	if frags == nil {
		return nil
	}

	converted := make(map[FragPosition][2]geom.Vec[T], len(frags))
	for pos, vecs := range frags {
		converted[pos] = [2]geom.Vec[T]{
			geom.NewVec(T(vecs[0].X), T(vecs[0].Y)),
			geom.NewVec(T(vecs[1].X), T(vecs[1].Y)),
		}
	}

	return converted
}
