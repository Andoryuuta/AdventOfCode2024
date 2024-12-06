package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Direction int

const (
	DIRECTION_UP Direction = iota
	DIRECTION_LEFT
	DIRECTION_DOWN
	DIRECTION_RIGHT
)

type Point struct {
	Row int
	Col int
}

type PuzzleMap struct {
	MapData             [][]rune
	GuardStartPosition  Point
	GuardStartDirection Direction
}

// Parses the puzzle map input into a 2D slice of runes.
func parseMap(reader io.Reader) (*PuzzleMap, error) {
	puzzleMap := &PuzzleMap{}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		rowText := scanner.Text()
		row := []rune(rowText)

		// If this row has a guard, save it's position and direction, then remove.
		directionIndex := strings.IndexAny(rowText, "^<v>")
		if directionIndex != -1 {
			puzzleMap.GuardStartPosition.Row = len(puzzleMap.MapData)
			puzzleMap.GuardStartPosition.Col = directionIndex

			switch row[directionIndex] {
			case '^':
				puzzleMap.GuardStartDirection = DIRECTION_UP
			case '<':
				puzzleMap.GuardStartDirection = DIRECTION_LEFT
			case 'v':
				puzzleMap.GuardStartDirection = DIRECTION_DOWN
			case '>':
				puzzleMap.GuardStartDirection = DIRECTION_RIGHT
			}

			row[directionIndex] = '.'
		}

		// Verify all rows are the same length
		if len(puzzleMap.MapData) > 0 && len(puzzleMap.MapData[0]) != len(row) {
			return nil, fmt.Errorf("expected all rows in the input to be the same length")
		}

		puzzleMap.MapData = append(puzzleMap.MapData, row)
	}
	err := scanner.Err()
	if err != nil {
		return nil, err
	}
	return puzzleMap, nil
}

// simulateGuardPatrol simluates the guard patrol of the provided puzzle map.
// Returns the distinct points walked by the guard, and whether the
// guard entered an infinte loop.
func simulateGuardPatrol(puzzleMap *PuzzleMap) (map[Point]bool, bool) {
	trailLogKey := func(pos Point, dir Direction) string {
		return fmt.Sprintf("%d-%d-%d", pos.Row, pos.Col, dir)
	}

	mapRuneAt := func(row int, col int) *rune {
		if row >= 0 && row < len(puzzleMap.MapData) && col >= 0 && col < len(puzzleMap.MapData[row]) {
			return &puzzleMap.MapData[row][col]
		}
		return nil
	}

	forwardPosition := func(pos Point, dir Direction) *Point {
		switch dir {
		case DIRECTION_UP:
			return &Point{pos.Row - 1, pos.Col}
		case DIRECTION_LEFT:
			return &Point{pos.Row, pos.Col - 1}
		case DIRECTION_DOWN:
			return &Point{pos.Row + 1, pos.Col}
		case DIRECTION_RIGHT:
			return &Point{pos.Row, pos.Col + 1}
		default:
			panic(fmt.Sprintf("invalid direction %v", dir))
		}
	}

	curPosition := puzzleMap.GuardStartPosition
	curDir := puzzleMap.GuardStartDirection
	trailLog := make(map[string]bool)
	seenPoints := make(map[Point]bool)
	infiniteLoop := false
	for {
		if trailLog[trailLogKey(curPosition, curDir)] {
			infiniteLoop = true
			break
		}
		trailLog[trailLogKey(curPosition, curDir)] = true
		seenPoints[curPosition] = true

		nextPos := forwardPosition(curPosition, curDir)
		nextPosRune := mapRuneAt(nextPos.Row, nextPos.Col)
		if nextPosRune == nil {
			// Guard left map
			break
		} else if *nextPosRune == '#' {
			// Guard hit obstruction, turn right 90 deg
			switch curDir {
			case DIRECTION_UP:
				curDir = DIRECTION_RIGHT
			case DIRECTION_RIGHT:
				curDir = DIRECTION_DOWN
			case DIRECTION_DOWN:
				curDir = DIRECTION_LEFT
			case DIRECTION_LEFT:
				curDir = DIRECTION_UP
			}
		} else {
			// Guard is free to move foward
			curPosition = *nextPos
		}
	}

	return seenPoints, infiniteLoop
}

// Literally brute-force currently.
// TODO(Andoryuuta): Make a real solution before the day ends!
func findAllLoopingOptions(puzzleMap *PuzzleMap) []Point {
	// Simulate it once to get the list of points walked by the guard.
	normallySeenPoints, _ := simulateGuardPatrol(puzzleMap)

	// Create a map of possible obstruction points.
	// We just add one for each cardinal direction from the point here,
	// could by optimized by keeping track of seen directions for each point.
	possibleObstructionPoints := make(map[Point]bool, len(normallySeenPoints)*5)
	for point, _ := range normallySeenPoints {
		possibleObstructionPoints[Point{point.Row, point.Col}] = true

		// Up
		if point.Row-1 >= 0 {
			possibleObstructionPoints[Point{point.Row - 1, point.Col}] = true
		}
		// Down
		if point.Row+1 < len(puzzleMap.MapData) {
			possibleObstructionPoints[Point{point.Row + 1, point.Col}] = true
		}
		// Left
		if point.Col-1 >= 0 {
			possibleObstructionPoints[Point{point.Row, point.Col - 1}] = true
		}
		// Right
		if len(puzzleMap.MapData) > 0 && point.Col+1 < len(puzzleMap.MapData[0]) {
			possibleObstructionPoints[Point{point.Row, point.Col + 1}] = true
		}
	}
	fmt.Printf("Testing %d possible obstruction points\n", len(possibleObstructionPoints))

	var loopCausingObstructions []Point
	for point, _ := range possibleObstructionPoints {
		// Only try to add obstructions where there aren't any existing, and not in the starting position.
		isStartingPos := point.Row == puzzleMap.GuardStartPosition.Row && point.Col == puzzleMap.GuardStartPosition.Col
		if puzzleMap.MapData[point.Row][point.Col] != '#' && !isStartingPos {
			newPuzzleMap := &PuzzleMap{
				MapData:             puzzleMap.MapData,
				GuardStartPosition:  puzzleMap.GuardStartPosition,
				GuardStartDirection: puzzleMap.GuardStartDirection,
			}

			// We swap the rune in the map data, simulate, then put it back.
			// This keeps us from having to deep copy the large map data for
			// each possible solution.
			originalRune := newPuzzleMap.MapData[point.Row][point.Col]
			newPuzzleMap.MapData[point.Row][point.Col] = '#'
			_, infiniteLoop := simulateGuardPatrol(newPuzzleMap)
			newPuzzleMap.MapData[point.Row][point.Col] = originalRune

			// fmt.Printf("point.Row:%d,point.Col:%d,distinctPositions:%d,loop:%v\n", point.Row, point.Col, len(distinctPositions), infiniteLoop)
			if infiniteLoop {
				loopCausingObstructions = append(loopCausingObstructions, Point{point.Row, point.Col})
			}
		}
	}
	return loopCausingObstructions
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

	puzzleMap, err := parseMap(file)
	if err != nil {
		log.Fatalf("error parsing map input: %v\n", err)
	}
	// puzzleMapJson, _ := json.MarshalIndent(puzzleMap, "", "\t")
	// fmt.Printf("Puzzle Input: %v\n", string(puzzleMapJson))

	distinctPositions, infiniteLoop := simulateGuardPatrol(puzzleMap)
	fmt.Printf("Distinct positions in simulated guard path: %d, infinite loop: %v\n", len(distinctPositions), infiniteLoop)

	loopCausingObstructions := findAllLoopingOptions(puzzleMap)
	fmt.Printf("Loop-causing obstruction options: %v\n", len(loopCausingObstructions))
}
