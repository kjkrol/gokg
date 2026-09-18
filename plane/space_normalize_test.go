package plane

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
)

func TestEuclidean2D_normalizeBox(t *testing.T) {
	euclidean := NewEuclidean2D(10, 10)

	for _, tc := range euclideanNormalizeScenarios() {
		t.Run(tc.name, func(t *testing.T) {
			aabb := NewAABB(vec(tc.topLeft.X, tc.topLeft.Y), tc.width, tc.height)
			euclidean.(*euclidean2d).normalizeAABB(&aabb)
			expectAABBState(t, aabb,
				vec(tc.expectedTopLeft.X, tc.expectedTopLeft.Y),
				vec(tc.expectedBottomRight.X, tc.expectedBottomRight.Y),
				map[FragPosition][2]geom.Vec{},
			)
		})
	}
}

type normalizeScenario struct {
	name                string
	topLeft             geom.Vec
	width, height       float64
	expectedTopLeft     geom.Vec
	expectedBottomRight geom.Vec
}

func euclideanNormalizeScenarios() []normalizeScenario {
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

func TestToroidal2D_normalizeBox(t *testing.T) {
	toroidal := NewToroidal2D(10, 10)

	for _, tc := range toroidalNormalizeScenarios() {
		t.Run(tc.name, func(t *testing.T) {
			aabb := NewAABB(vec(tc.topLeft.X, tc.topLeft.Y), tc.width, tc.height)
			toroidal.(*toroidal2d).normalizeAABB(&aabb)
			expectAABBState(t, aabb,
				vec(tc.expectedTopLeft.X, tc.expectedTopLeft.Y),
				vec(tc.expectedBottomRight.X, tc.expectedBottomRight.Y),
				convertFragments(tc.expectedFragments),
			)
		})
	}
}

type toroidalNormalizeScenario struct {
	normalizeScenario
	expectedFragments map[FragPosition][2]geom.Vec
}

func toroidalNormalizeScenarios() []toroidalNormalizeScenario {
	return []toroidalNormalizeScenario{
		{
			normalizeScenario: normalizeScenario{
				name:                "wraps_box_into_view",
				topLeft:             geom.NewVec(12, 12),
				width:               2,
				height:              2,
				expectedTopLeft:     geom.NewVec(2, 2),
				expectedBottomRight: geom.NewVec(4, 4),
			},
			expectedFragments: map[FragPosition][2]geom.Vec{},
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
			expectedFragments: map[FragPosition][2]geom.Vec{
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
			expectedFragments: map[FragPosition][2]geom.Vec{
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
			expectedFragments: map[FragPosition][2]geom.Vec{
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
			expectedFragments: map[FragPosition][2]geom.Vec{},
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
			expectedFragments: map[FragPosition][2]geom.Vec{
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
			expectedFragments: map[FragPosition][2]geom.Vec{
				FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
				FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
				FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
			},
		}}
}

func convertFragments(frags map[FragPosition][2]geom.Vec) map[FragPosition][2]geom.Vec {
	if frags == nil {
		return nil
	}

	converted := make(map[FragPosition][2]geom.Vec, len(frags))
	for pos, vecs := range frags {
		converted[pos] = [2]geom.Vec{
			geom.NewVec(float64(vecs[0].X), float64(vecs[0].Y)),
			geom.NewVec(float64(vecs[1].X), float64(vecs[1].Y)),
		}
	}

	return converted
}
