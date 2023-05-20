package board_validation

import (
	"reflect"
	"testing"
)

func TestTwoPassConnectedComponentLabeling(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name   string
		img    [][]int
		labels [][]int
	}{
		{
			name: "Simple Case",
			img: [][]int{
				{0, 1, 0},
				{1, 1, 1},
				{0, 1, 0},
			},
			labels: [][]int{
				{0, 1, 0},
				{1, 1, 1},
				{0, 1, 0},
			},
		},
		{
			name: "Two Components",
			img: [][]int{
				{0, 1, 0, 1, 0},
				{0, 1, 1, 1, 0},
				{0, 0, 0, 1, 0},
				{0, 1, 0, 0, 0},
				{0, 1, 0, 1, 1},
			},
			labels: [][]int{
				{0, 1, 0, 1, 0},
				{0, 1, 1, 1, 0},
				{0, 0, 0, 1, 0},
				{0, 2, 0, 0, 0},
				{0, 2, 0, 3, 3},
			},
		},
		{
			name: "Two Components sharing corner",
			img: [][]int{
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 1, 0, 0},
				{0, 0, 1, 0, 0},
				{0, 0, 0, 1, 1},
			},
			labels: [][]int{
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 1, 0, 0},
				{0, 0, 1, 0, 0},
				{0, 0, 0, 2, 2},
			},
		},
		{
			name: "Five components",
			img: [][]int{
				{1, 0, 1},
				{0, 1, 0},
				{1, 0, 1},
			},
			labels: [][]int{
				{1, 0, 2},
				{0, 3, 0},
				{4, 0, 5},
			},
		},
	}

	// Iterate over the test cases and run the tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ConnectedComponentLabeling(tc.img)

			// Check that the output labels match the expected labels
			if !reflect.DeepEqual(got, tc.labels) {
				t.Errorf("Expected labels:\n%v\nGot labels:\n%v", tc.labels, got)
			}
		})
	}
}
