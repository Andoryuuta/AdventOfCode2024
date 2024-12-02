package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// parseLocationList parses a location list and returns two lists (one for each column).
// The two returned lists are guaranteed to have the same number of elements.
//
// This function assumes the input conforms to the following grammar:
//
//	LocationList ::= (LocationListPair ('\n')? )*
//	LocationListPair ::= (Digits) ('   ') (Digits)
//	Digits ::= #'[0-9]+'
func parseLocationList(reader io.Reader) (locationList1 []uint64, locationList2 []uint64, err error) {
	lineNumber := 1
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lineNumber += 1

		splitStrings := strings.Split(scanner.Text(), "   ")
		if len(splitStrings) != 2 {
			return nil, nil, fmt.Errorf("unexcepted format on input line %d", lineNumber)
		}

		first, err := strconv.ParseUint(splitStrings[0], 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("unexcepted format of first decimal on input line %d", lineNumber)
		}
		locationList1 = append(locationList1, first)

		second, err := strconv.ParseUint(splitStrings[1], 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("unexcepted format of second decimal on input line %d", lineNumber)
		}
		locationList2 = append(locationList2, second)
	}

	return
}

func distance(a uint64, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}

// calcListDistance calculates the total distance of the two provided lists.
// (This is for part 1 of the AOC challenge)
func calcListDistance(left []uint64, right []uint64) uint64 {
	// Sort the lists in order to comply with the matching requirement:
	//
	// "... pair up the numbers and measure how far apart they are.
	// Pair up the smallest number in the left list with the smallest number in the right list,
	// then the second-smallest left number with the second-smallest right number, and so on."
	sort.Slice(left, func(i, j int) bool { return left[i] < left[j] })
	sort.Slice(right, func(i, j int) bool { return right[i] < right[j] })

	// Calculate the total distance:
	//
	// "... Within each pair, figure out how far apart the two numbers are;
	// you'll need to add up all of those distances."
	var totalDistance uint64 = 0
	for i := 0; i < len(left); i++ {
		totalDistance += distance(left[i], right[i])
	}

	return totalDistance
}

// calcSimilarityScore calculates the similarity score of the two provided lists.
// (This is for part 2 of the AOC challenge)
func calcSimilarityScore(left []uint64, right []uint64) uint64 {
	rightListOccurances := make(map[uint64]uint64)
	for _, rv := range right {
		rightListOccurances[rv] += 1
	}

	var similarityScore uint64
	for _, lv := range left {
		// If the left value (lv) has not occured in the right list,
		// this map lookup will return 0, which makes the part 2 example(s).
		similarityScore += lv * rightListOccurances[lv]
	}
	return similarityScore
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
	defer file.Close()

	locationList1, locationList2, err := parseLocationList(file)
	if err != nil {
		log.Fatalf("error parsing location list: %v\n", err)
		return
	}

	totalDistance := calcListDistance(locationList1, locationList2)
	fmt.Printf("Total distance: %d\n", totalDistance)

	similarityScore := calcSimilarityScore(locationList1, locationList2)
	fmt.Printf("Similarity score: %d\n", similarityScore)
}
