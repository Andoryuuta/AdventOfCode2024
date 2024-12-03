package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type MulExpression [2]int

func extractMulExpressions(data []byte) ([]MulExpression, error) {
	// TODO: don't recompile RE
	re := regexp.MustCompile(`mul\((\d\d?\d?),(\d\d?\d?)\)`)
	matches := re.FindAllSubmatch(data, -1)
	if matches == nil {
		return nil, fmt.Errorf("failed to match any mul expressions in input")
	}

	var expressions []MulExpression
	for _, matchGroups := range matches {
		lefthand, err := strconv.Atoi(string(matchGroups[1]))
		if err != nil {
			return nil, err
		}

		righthand, err := strconv.Atoi(string(matchGroups[2]))
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, MulExpression{lefthand, righthand})
	}

	return expressions, nil
}

func calcMulExpressionSum(data []byte) (int, error) {
	expressions, err := extractMulExpressions(data)
	if err != nil {
		return 0, err
	}

	sum := 0
	for _, expression := range expressions {
		sum += (expression[0] * expression[1])
	}

	return sum, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("must provide input filename as an argument")
		return
	}
	filename := os.Args[1]

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("cannot open input file: %v\n", err)
	}

	expressionSum, err := calcMulExpressionSum(data)
	if err != nil {
		log.Fatalf("error calcuating mul expression sum: %v\n", err)
	}

	fmt.Printf("Sum of all mul express products: %d\n", expressionSum)

}
