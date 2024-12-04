package main

import (
	"slices"
	"strings"
	"testing"
)

func TestParseWordsearch(t *testing.T) {
	var tests = []struct {
		name           string
		input          string
		expectedOutput [][]rune
		expectedError  bool
	}{
		{
			"valid square",
			"AB\nCD\n",
			[][]rune{
				{'A', 'B'},
				{'C', 'D'},
			},
			false,
		},
		{
			"invalid input (rows different length)",
			"AB\nCDE\n",
			[][]rune{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			gotOutput, gotErr := parseWordSearch(reader)

			if tt.expectedError && gotErr == nil {
				t.Errorf("got error %v, expected nil", gotErr)
				return
			}

			if len(gotOutput) != len(tt.expectedOutput) {
				t.Errorf("got output length %v, expected %v", len(gotOutput), len(tt.expectedOutput))
				return
			}

			for i := 0; i < len(gotOutput); i++ {
				got := gotOutput[i]
				expected := tt.expectedOutput[i]
				if !slices.Equal(got, expected) {
					t.Errorf("got args %+v, expected %+v", got, expected)
					return
				}
			}
		})
	}
}

func TestSearchShape(t *testing.T) {
	var tests = []struct {
		name            string
		inputWordSearch [][]rune
		inputShape      [][]rune
		inputShapeMask  [][]bool
		expectedOutput  []Point
	}{
		{
			"Valid row matching",
			[][]rune{
				{'A', 'B', 'C', 'D'},
				{'X', 'X', 'X', 'X'},
				{'A', 'B', 'C', 'D'},
				{'X', 'X', 'X', 'X'},
			},
			[][]rune{
				{'A', 'B', 'C', 'D'},
			},
			[][]bool{
				{true, true, true, true},
			},
			[]Point{
				{0, 0},
				{2, 0},
			},
		},
		{
			"Valid inner square matching",
			[][]rune{
				{'X', 'X', 'X', 'X'},
				{'X', 'A', 'A', 'X'},
				{'X', 'A', 'A', 'X'},
				{'X', 'X', 'X', 'X'},
			},
			[][]rune{
				{'A', 'A'},
				{'A', 'A'},
			},
			[][]bool{
				{true, true},
				{true, true},
			},
			[]Point{
				{1, 1},
			},
		},
		{
			"Valid inner masked shape matching",
			[][]rune{
				{'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', '@', 'X', 'X', 'X', 'X'},
				{'X', '@', 'X', 'X', 'X', 'X'},
				{'X', '@', '@', '@', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', '@', 'X', 'X', 'X'},
				{'X', 'X', '@', 'X', 'X', 'X'},
				{'X', 'X', '@', '@', '@', 'X'},
			},
			[][]rune{
				{'@', ' ', ' '},
				{'@', ' ', ' '},
				{'@', '@', '@'},
			},
			[][]bool{
				{true, false, false},
				{true, false, false},
				{true, true, true},
			},
			[]Point{
				{2, 1},
				{9, 2},
			},
		},
		{
			"Valid inner masked shape matching with overlap",
			[][]rune{
				{'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', '@', 'X', 'X', 'X', 'X'},
				{'X', '@', 'X', 'X', 'X', 'X'},
				{'X', '@', '@', '@', 'X', 'X'},
				{'X', 'X', '@', 'X', 'X', 'X'},
				{'X', 'X', '@', '@', '@', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X'},
			},
			[][]rune{
				{'@', ' ', ' '},
				{'@', ' ', ' '},
				{'@', '@', '@'},
			},
			[][]bool{
				{true, false, false},
				{true, false, false},
				{true, true, true},
			},
			[]Point{
				{2, 1},
				{4, 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput := searchShape(tt.inputWordSearch, tt.inputShape, tt.inputShapeMask)

			if len(gotOutput) != len(tt.expectedOutput) {
				t.Errorf("got output length %v, expected %v", len(gotOutput), len(tt.expectedOutput))
				return
			}

			got := gotOutput
			expected := tt.expectedOutput
			if !slices.Equal(got, expected) {
				t.Errorf("got args %+v, expected %+v", got, expected)
				return
			}
		})
	}
}
