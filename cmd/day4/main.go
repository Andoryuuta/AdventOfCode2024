package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// "parses" the word-search input into a 2D array of runes.
// This function requires that all input rows are the same length.
func parseWordSearch(reader io.Reader) ([][]rune, error) {
	// Convert data into 2D rune arrays
	var data2d [][]rune

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		row := []rune(scanner.Text())

		// Verify all rows are the same length
		if len(data2d) > 0 && len(data2d[0]) != len(row) {
			return nil, fmt.Errorf("expected all rows in the input to be the same length")
		}

		data2d = append(data2d, row)
	}
	err := scanner.Err()
	if err != nil {
		return nil, err
	}
	return data2d, nil
}

// Point represents a single point in a 2D array
// This struct uses row/col (rather than x/y) for clarity.
type Point struct {
	row uint
	col uint
}

func shapeMatch(data2d [][]rune, shape [][]rune, mask [][]bool, matchRow uint, matchCol uint) bool {
	runeAt := func(row uint, col uint) *rune {
		if row < uint(len(data2d)) && col < uint(len(data2d[row])) {
			return &data2d[row][col]
		}
		return nil
	}

	for srow := uint(0); srow < uint(len(shape)); srow++ {
		for scol := uint(0); scol < uint(len(shape[srow])); scol++ {
			// This part of the shape is masked off, don't check it.
			if !mask[srow][scol] {
				continue
			}

			dRune := runeAt(matchRow+srow, matchCol+uint(scol))
			sRune := shape[srow][scol]
			if dRune == nil || *dRune != sRune {
				return false
			}
		}
	}

	return true
}

func searchShape(data2d [][]rune, shape [][]rune, mask [][]bool) []Point {
	var matches []Point
	for row := uint(0); row < uint(len(data2d)); row++ {
		for col := uint(0); col < uint(len(data2d[row])); col++ {
			matchResult := shapeMatch(data2d, shape, mask, row, col)
			if matchResult {
				matches = append(matches, Point{row, col})
			}
		}
	}

	return matches
}

// countXmasShapePart1 counts instances of "XMAS" in the input
// (horizontal/vertical/diagonal, allowing reverse spelling)
func countXmasShapePart1(data2d [][]rune) uint {
	shapes := []struct {
		shape [][]rune
		mask  [][]bool
	}{
		// Horizontal
		{
			[][]rune{
				{'X', 'M', 'A', 'S'},
			},
			[][]bool{
				{true, true, true, true},
			},
		},
		{
			[][]rune{
				{'S', 'A', 'M', 'X'},
			},
			[][]bool{
				{true, true, true, true},
			},
		},

		// Vertical
		{
			[][]rune{
				{'X'},
				{'M'},
				{'A'},
				{'S'},
			},
			[][]bool{
				{true},
				{true},
				{true},
				{true},
			},
		},
		{
			[][]rune{
				{'S'},
				{'A'},
				{'M'},
				{'X'},
			},
			[][]bool{
				{true},
				{true},
				{true},
				{true},
			},
		},

		// Diag Right
		{
			[][]rune{
				{'X', ' ', ' ', ' '},
				{' ', 'M', ' ', ' '},
				{' ', ' ', 'A', ' '},
				{' ', ' ', ' ', 'S'},
			},
			[][]bool{
				{true, false, false, false},
				{false, true, false, false},
				{false, false, true, false},
				{false, false, false, true},
			},
		},
		{
			[][]rune{
				{'S', ' ', ' ', ' '},
				{' ', 'A', ' ', ' '},
				{' ', ' ', 'M', ' '},
				{' ', ' ', ' ', 'X'},
			},
			[][]bool{
				{true, false, false, false},
				{false, true, false, false},
				{false, false, true, false},
				{false, false, false, true},
			},
		},

		// Diag Left
		{
			[][]rune{
				{' ', ' ', ' ', 'X'},
				{' ', ' ', 'M', ' '},
				{' ', 'A', ' ', ' '},
				{'S', ' ', ' ', ' '},
			},
			[][]bool{
				{false, false, false, true},
				{false, false, true, false},
				{false, true, false, false},
				{true, false, false, false},
			},
		},
		{
			[][]rune{
				{' ', ' ', ' ', 'S'},
				{' ', ' ', 'A', ' '},
				{' ', 'M', ' ', ' '},
				{'X', ' ', ' ', ' '},
			},
			[][]bool{
				{false, false, false, true},
				{false, false, true, false},
				{false, true, false, false},
				{true, false, false, false},
			},
		},
	}

	count := 0
	for _, shapeGroup := range shapes {
		points := searchShape(
			data2d,
			shapeGroup.shape,
			shapeGroup.mask,
		)
		// fmt.Printf("Shape group %d matching points: %+v\n", i, points)
		count += len(points)
	}

	return uint(count)
}

// countXmasShapePart1 counts instance of cross "MAS" shapes in the input
func countXmasShapePart2(data2d [][]rune) uint {
	crossShapeMask := [][]bool{
		{true, false, true},
		{false, true, false},
		{true, false, true},
	}

	shapes := []struct {
		shape [][]rune
		mask  [][]bool
	}{
		{
			[][]rune{
				{'M', ' ', 'M'},
				{' ', 'A', ' '},
				{'S', ' ', 'S'},
			},
			crossShapeMask,
		},
		{
			[][]rune{
				{'S', ' ', 'M'},
				{' ', 'A', ' '},
				{'S', ' ', 'M'},
			},
			crossShapeMask,
		},
		{
			[][]rune{
				{'M', ' ', 'S'},
				{' ', 'A', ' '},
				{'M', ' ', 'S'},
			},
			crossShapeMask,
		},
		{
			[][]rune{
				{'S', ' ', 'S'},
				{' ', 'A', ' '},
				{'M', ' ', 'M'},
			},
			crossShapeMask,
		},
	}

	count := 0
	for _, shapeGroup := range shapes {
		points := searchShape(
			data2d,
			shapeGroup.shape,
			shapeGroup.mask,
		)
		count += len(points)
	}

	return uint(count)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("must provide input filename as an argument")
		return
	}
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open input file: %v\n", err)
	}

	wordSearch, err := parseWordSearch(file)
	if err != nil {
		log.Fatalf("error parsing word-search input: %v\n", err)
	}

	countPart1 := countXmasShapePart1(wordSearch)
	fmt.Printf("Horizontal/Vertical/Diagonal count (part 1): %d\n", countPart1)

	countPart2 := countXmasShapePart2(wordSearch)
	fmt.Printf("Cross-'MAS' shape count (part 2): %d\n", countPart2)
}
