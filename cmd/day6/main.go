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
	deepCopyMapData := func(puzzleMap *PuzzleMap) [][]rune {
		rows := [][]rune{}
		for row := 0; row < len(puzzleMap.MapData); row++ {
			v := make([]rune, len(puzzleMap.MapData[row]))
			copy(v, puzzleMap.MapData[row])
			rows = append(rows, v)
		}
		return rows
	}

	var loopCausingObstructions []Point
	for row := 0; row < len(puzzleMap.MapData); row++ {
		for col := 0; col < len(puzzleMap.MapData[row]); col++ {
			// Only try to add obstructions where there aren't any existing, and not in the starting position.
			isStartingPos := row == puzzleMap.GuardStartPosition.Row && col == puzzleMap.GuardStartPosition.Col
			if puzzleMap.MapData[row][col] != '#' && !isStartingPos {
				newPuzzleMap := &PuzzleMap{
					MapData:             deepCopyMapData(puzzleMap),
					GuardStartPosition:  puzzleMap.GuardStartPosition,
					GuardStartDirection: puzzleMap.GuardStartDirection,
				}

				newPuzzleMap.MapData[row][col] = '#'

				_, infiniteLoop := simulateGuardPatrol(newPuzzleMap)
				// fmt.Printf("row:%d,col:%d,distinctPositions:%d,loop:%v\n", row, col, len(distinctPositions), infiniteLoop)
				if infiniteLoop {
					loopCausingObstructions = append(loopCausingObstructions, Point{row, col})
				}
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
