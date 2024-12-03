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

type Report []uint64

// parseLocationList parses a list of reports.
// This function assumes the input conforms to the following grammar:
//
//	ReportList ::= (Report ('\n')?)+
//	Report ::= (Level (' ')?)+
//	Level ::= (Digits)
//	Digits ::= #'[0-9]+'
//
// Each report is gauranteed to have at least two level values.
func parseReportList(reader io.Reader) (reports []Report, err error) {
	lineNumber := 1
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lineNumber += 1
		var report Report

		splitStrings := strings.Split(scanner.Text(), " ")
		if len(splitStrings) < 2 {
			return nil, fmt.Errorf("unexcepted format of report, expected at least two level values (line %d)", lineNumber)
		}

		for idx, s := range splitStrings {
			parsedLevel, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("unexcepted format of report (line %d, level idx: %d): %v", lineNumber, idx, err)
			}
			report = append(report, parsedLevel)
		}

		reports = append(reports, report)
	}

	if scanner.Err() != nil {
		return nil, err
	}

	if len(reports) == 0 {
		return nil, fmt.Errorf("expected at least one report")
	}

	return
}

func absoluteDifference(a uint64, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}

func isReportSafeRaw(report Report) bool {
	if len(report) < 2 {
		panic("expected each report to have at least 2 levels")
	}

	increasing := report[0] < report[1]
	checkLevelChangeSafe := func(x, y uint64) bool {
		// "The levels are either all increasing or all decreasing."
		directionChanged := (x < y) != increasing

		// "Any two adjacent levels differ by at least one and at most three."
		diff := absoluteDifference(x, y)
		toleranceExceeded := diff < 1 || diff > 3

		return !(directionChanged || toleranceExceeded)
	}

	for i := 0; i < len(report)-1; i++ {
		safe := checkLevelChangeSafe(report[i], report[i+1])

		if !safe {
			return false
		}
	}

	return true
}

func isReportSafe(report Report, problemDampenerEnabled bool) bool {
	if problemDampenerEnabled {
		// NOTE(Andoryuuta): This permutation logic is gross, but I struggled to find a way
		// to incorporate the "Problem Dampener" (skipping a single bad level) into the normal
		// `isReportSafeRaw` function.
		//
		// Definitely worth looking at other people's solutions for this later to see how
		// this can be improved / done correctly.
		for i := 0; i < len(report); i++ {
			permutation := make(Report, len(report)-1)
			copy(permutation[:i], report[:i])
			copy(permutation[i:], report[i+1:])
			if isReportSafeRaw(permutation) {
				return true
			}
		}
	}

	return isReportSafeRaw(report)
}

func calcSafeReports(reports []Report, problemDampenerEnabled bool) uint64 {
	var safeCount uint64
	for _, report := range reports {
		if isReportSafe(report, problemDampenerEnabled) {
			safeCount += 1
		}
	}
	return uint64(safeCount)
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

	reports, err := parseReportList(file)
	if err != nil {
		log.Fatalf("error parsing location list: %v\n", err)
		return
	}

	safeReportsCountNoDampener := calcSafeReports(reports, false)
	fmt.Printf("Safe reports - No Problem Dampener (Part 1): %d\n", safeReportsCountNoDampener)

	safeReportsCountWithDampener := calcSafeReports(reports, true)
	fmt.Printf("Safe reports - With Problem Dampener (Part 2): %d\n", safeReportsCountWithDampener)
}
