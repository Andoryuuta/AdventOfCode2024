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

func getMinimalOrderingRules(orderingRules map[PageID][]PageID, pages PageList) map[PageID][]PageID {
	// Build a map lookup rather than array scanning
	contains := map[PageID]bool{}
	for _, pageID := range pages {
		contains[pageID] = true
	}

	// Filter down to a minimal set of ordering rules that are relevant.
	minimalOrderingRules := map[PageID][]PageID{}
	for _, pageID := range pages {
		for _, depPageID := range orderingRules[pageID] {
			if contains[depPageID] {
				minimalOrderingRules[pageID] = append(minimalOrderingRules[pageID], depPageID)

			}
		}
	}

	return minimalOrderingRules
}

func generateMermaidGraph(orderingRules map[PageID][]PageID, pages PageList) string {
	minimalOrderingRules := getMinimalOrderingRules(orderingRules, pages)
	output := "graph LR\n"
	joined := ""
	for idx, pageID := range pages {
		if idx == len(pages)-1 {
			joined += fmt.Sprintf("%d", pageID)
		} else {
			joined += fmt.Sprintf("%d, ", pageID)
		}
	}
	output += fmt.Sprintf("subgraph \"Pages: %v\"\n", joined)

	// Emit all page IDs (required to have a node for any non-connected node)
	for _, pageID := range pages {
		output += fmt.Sprintf("%d\n", pageID)
	}

	// Emit all page dependencies
	for pageID, dependencies := range minimalOrderingRules {
		for _, depPageID := range dependencies {
			output += fmt.Sprintf("%d --> %d\n", pageID, depPageID)
		}
	}

	output += "end\n"

	return output
}

// pseudoTopSortPageList attempts a pseudo-topological sort to find a solution.
// This is "pseudo" because I don't know graph theory.
// TODO(Andoryuuta): learn graph theory?
func pseudoTopSortPageList(orderingRules map[PageID][]PageID, invalidPageList PageList) (PageList, error) {
	// 0. Create list of each node and it's dependencies
	// 1. Emit whichever node doesn't have any dependencies
	// 2. Remove the node emitted in step 1 with from the dependency list of any other node
	// 3. Repeat 1&2 until list is empty.
	// 4. Hope and pray to whatever god you believe in that there aren't any circular dependencies

	// Filter down to a minimal set of ordering rules that are relevant.
	// The keys of this map will
	minimalOrderingRules := getMinimalOrderingRules(orderingRules, invalidPageList)

	// Build a map of pageID -> map dep pageID
	graph := map[PageID]map[PageID]bool{}
	for _, pageID := range invalidPageList {
		graph[pageID] = map[PageID]bool{}
	}
	for pageID, deps := range minimalOrderingRules {
		for _, depPageID := range deps {
			graph[pageID][depPageID] = true
		}
	}

	var solution []PageID
	graphLength := len(graph)
	for i := 0; i < graphLength; i++ {
		// fmt.Printf("Iter %d, graph: %+v\n", i, graph)
		var nodeToEmit *PageID

		for pageID, deps := range graph {
			if len(deps) == 0 {
				nodeToEmit = &pageID
				break
			}
		}

		if nodeToEmit == nil {
			break
		}

		// fmt.Printf("Removing node %v\n", *nodeToEmit)

		// Remove it from list
		delete(graph, *nodeToEmit)

		// Remove it deps of other nodes
		for _, deps := range graph {
			_, ok := deps[*nodeToEmit]
			if ok {
				delete(deps, *nodeToEmit)
			}
		}

		solution = append(solution, *nodeToEmit)
	}

	if len(graph) != 0 {
		return nil, fmt.Errorf("failed to solution satisfying ordering rules")
	}

	return solution, nil
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

// calculatePartOneSolution calculates the part 1 solution by filtering the
// input by valid "update" page lists, then summing the middle page number of each.
func calculatePartTwoSolution(orderingRules map[PageID][]PageID, updates []PageList) (int, error) {
	var invalidPageLists []PageList
	for _, pageList := range updates {
		if !isPageListCompliant(orderingRules, pageList) {
			invalidPageLists = append(invalidPageLists, pageList)
		}
	}

	sum := 0
	for _, pageList := range invalidPageLists {
		correctedSolution, err := pseudoTopSortPageList(orderingRules, pageList)
		if err != nil {
			return 0, err
		}
		sum += int(correctedSolution[len(correctedSolution)/2])
	}
	return sum, nil
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

	partTwoSolution, err := calculatePartTwoSolution(updateSummary.orderingRules, updateSummary.updates)
	if err != nil {
		log.Fatalf("error calculating part 2 solution: %v\n", err)
	}
	fmt.Printf("Part 2 solution: %v\n", partTwoSolution)

	// var mermaidGraph string
	// mermaidGraph = generateMermaidGraph(updateSummary.orderingRules, []PageID{75, 97, 47, 61, 53})
	// fmt.Printf("Mermaid Graph:\n%v\n", mermaidGraph)
	// mermaidGraph = generateMermaidGraph(updateSummary.orderingRules, []PageID{61, 13, 29})
	// fmt.Printf("Mermaid Graph:\n%v\n", mermaidGraph)
	// mermaidGraph = generateMermaidGraph(updateSummary.orderingRules, []PageID{97, 13, 75, 29, 47})
	// fmt.Printf("Mermaid Graph:\n%v\n", mermaidGraph)
}
