package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculatePacks(t *testing.T) {
	testCases := []struct {
		input  int
		output map[int]int
	}{
		{
			input:  1,
			output: map[int]int{250: 1},
		},
		{
			input:  250,
			output: map[int]int{250: 1},
		},
		{
			input:  251,
			output: map[int]int{500: 1},
		},
		{
			input:  501,
			output: map[int]int{500: 1, 250: 1},
		},
		{
			input:  750,
			output: map[int]int{500: 1, 250: 1},
		},
		{
			input:  3000,
			output: map[int]int{2000: 1, 1000: 1},
		},
		{
			input:  3001,
			output: map[int]int{2000: 1, 1000: 1, 250: 1},
		},
		{
			input:  12001,
			output: map[int]int{250: 1, 2000: 1, 5000: 2},
		},
		{
			input:  20000,
			output: map[int]int{5000: 4},
		},
	}

	for _, tc := range testCases {
		t.Logf("Testing %d", tc.input)

		t.Run("calculatePacks", func(t *testing.T) {
			packs := calculatePacks(tc.input)
			require.Equal(t, tc.output, packs, "outputs should be equal")
		})
	}
}
