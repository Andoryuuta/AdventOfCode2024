package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type PageID int
type PageList []PageID

type UpdateSummary struct {
	// mapping of page number -> any dependency page numbers.
	orderingRules map[PageID][]PageID

	updates []PageList
}

func parseUpdateSummary(reader io.Reader) (*UpdateSummary, error) {
	summary := &UpdateSummary{
		orderingRules: map[PageID][]PageID{},
		updates:       []PageList{},
	}

	parsingOrderingRules := true
	scanner := bufio.NewScanner(reader)
	lineNumber := 1
	for scanner.Scan() {
		lineNumber += 1
		line := scanner.Text()

		// First empty line is the separator between the ordering rules
		// and the update page list.
		if parsingOrderingRules && line == "" {
			parsingOrderingRules = false
			continue
		}

		if parsingOrderingRules {
			split := strings.Split(line, "|")
			if len(split) != 2 {
				return nil, fmt.Errorf("expected page ordering rule on line %d, got:%v", lineNumber, line)
			}

			depPageID, err := strconv.Atoi(split[0])
			if err != nil {
				return nil, fmt.Errorf("error on line %d: %v", err)
			}

			pageID, err := strconv.Atoi(split[1])
			if err != nil {
				return nil, fmt.Errorf("error on line %d: %v", err)
			}

			typedPageID := PageID(pageID)
			summary.orderingRules[typedPageID] = append(summary.orderingRules[typedPageID], PageID(depPageID))
		} else {
			var update PageList
			split := strings.Split(line, ",")
			for _, pageID := range split {
				parsedPageID, err := strconv.Atoi(pageID)
				if err != nil {
					return nil, fmt.Errorf("error on line %d: %v", err)
				}
				update = append(update, PageID(parsedPageID))
			}
			summary.updates = append(summary.updates, update)
		}
	}
	err := scanner.Err()
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func isPageListCompliant(orderingRules map[PageID][]PageID, pages PageList) bool {
	// Build a map lookup rather than array scanning
	contains := map[PageID]bool{}
	for _, pageID := range pages {
		contains[pageID] = true
	}

	// Verify all applicable rules.
	seen := map[PageID]bool{}
	for _, pageID := range pages {
		for _, depPageID := range orderingRules[pageID] {

			// Only apply the rule if both pages (the pageID and depdency pageID) are in the list.
			// pageID is (obviously) in the list, so we just need to check the deps.
			if !contains[depPageID] {
				continue
			}

			if !seen[depPageID] {
				// fmt.Printf("PageID %d requires %d, which hasn't been seen (list: %v)\n", pageID, depPageID, pages)
				return false
			}
		}
		seen[pageID] = true
	}

	return true
}

// calculatePartOneSolution calculates the part 1 solution by filtering the
// input by valid "update" page lists, then summing the middle page number of each.
func calculatePartOneSolution(orderingRules map[PageID][]PageID, updates []PageList) int {
	var validPageLists []PageList
	for _, pageList := range updates {
		if isPageListCompliant(orderingRules, pageList) {
			validPageLists = append(validPageLists, pageList)
		}
	}

	sum := 0
	for _, pageList := range validPageLists {
		sum += int(pageList[len(pageList)/2])
	}
	return sum
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

	updateSummary, err := parseUpdateSummary(file)
	if err != nil {
		log.Fatalf("error parsing update summary input: %v\n", err)
	}
	// fmt.Printf("Summary: %+v\n", updateSummary)

	partOneSolution := calculatePartOneSolution(updateSummary.orderingRules, updateSummary.updates)
	fmt.Printf("Part 1 solution: %v\n", partOneSolution)

}
