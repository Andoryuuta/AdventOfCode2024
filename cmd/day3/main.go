package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Opcode uint

const (
	OP_DO Opcode = iota
	OP_DONT
	OP_MUL
)

type Instruction struct {
	op   Opcode
	args []uint
}

func extractInstructions(data []byte) ([]Instruction, error) {
	re := regexp.MustCompile(`(do\(\)|don't\(\)|mul\((\d\d?\d?),(\d\d?\d?)\))`)
	matches := re.FindAllSubmatch(data, -1)
	if matches == nil {
		return nil, fmt.Errorf("failed to match any expressions in input")
	}

	var instructions []Instruction
	for _, matchGroups := range matches {
		var instruction Instruction

		fullMatchString := string(matchGroups[0])
		if fullMatchString == "do()" {
			instruction.op = OP_DO
		} else if fullMatchString == "don't()" {
			instruction.op = OP_DONT
		} else if strings.HasPrefix(fullMatchString, "mul(") {
			instruction.op = OP_MUL

			lefthand, err := strconv.Atoi(string(matchGroups[2]))
			if err != nil {
				return nil, err
			}
			instruction.args = append(instruction.args, uint(lefthand))

			righthand, err := strconv.Atoi(string(matchGroups[3]))
			if err != nil {
				return nil, err
			}
			instruction.args = append(instruction.args, uint(righthand))
		}
		instructions = append(instructions, instruction)
	}

	return instructions, nil
}

func evaluateProgram(instructions []Instruction, evalToggleInstructions bool) uint {
	mulEnabled := true
	accumulator := uint(0)
	for _, instruction := range instructions {
		switch instruction.op {
		case OP_DO:
			if evalToggleInstructions {
				mulEnabled = true
			}
		case OP_DONT:
			if evalToggleInstructions {
				mulEnabled = false
			}
		case OP_MUL:
			if mulEnabled {
				product := instruction.args[0] * instruction.args[1]
				accumulator += product
			}
		}
	}

	return accumulator
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

	instructions, err := extractInstructions(data)
	if err != nil {
		log.Fatalf("error extracting instructions from input: %v\n", err)
	}

	expressionSumNoToggles := evaluateProgram(instructions, false)
	fmt.Printf("Sum of all 'mul' instruction products (toggle instruction not evaluated - part 1): %d\n", expressionSumNoToggles)

	expressionSumWithToggles := evaluateProgram(instructions, true)
	fmt.Printf("Sum of all 'mul' instruction products (toggle instruction evaluated - part 2): %d\n", expressionSumWithToggles)

}
