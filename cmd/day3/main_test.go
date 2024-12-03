package main

import (
	"slices"
	"testing"
)

func TestExtractInstructions(t *testing.T) {
	var tests = []struct {
		name                 string
		input                string
		expectedInstructions *[]Instruction
		expectedError        bool
	}{
		{
			"valid instructions extracted (part1 example)",
			`xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))`,
			&[]Instruction{
				{op: OP_MUL, args: []uint{2, 4}},
				{op: OP_MUL, args: []uint{5, 5}},
				{op: OP_MUL, args: []uint{11, 8}},
				{op: OP_MUL, args: []uint{8, 5}},
			},
			false,
		},
		{
			"valid instructions extracted (part2 example)",
			`xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))`,
			&[]Instruction{
				{op: OP_MUL, args: []uint{2, 4}},
				{op: OP_DONT, args: []uint{}},
				{op: OP_MUL, args: []uint{5, 5}},
				{op: OP_MUL, args: []uint{11, 8}},
				{op: OP_DO, args: []uint{}},
				{op: OP_MUL, args: []uint{8, 5}},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInstructions, gotErr := extractInstructions([]byte(tt.input))

			if tt.expectedError && gotErr == nil {
				t.Errorf("got %v, expected nil", gotErr)
				return
			}

			if tt.expectedInstructions != nil {
				if len(gotInstructions) != len(*tt.expectedInstructions) {
					t.Errorf("unexpected instruction count, got %v, expected %v", len(gotInstructions), len(*tt.expectedInstructions))
					return
				}

				for i := 0; i < len(*tt.expectedInstructions); i++ {
					got := gotInstructions[i]
					expected := (*tt.expectedInstructions)[i]

					if got.op != expected.op {
						t.Errorf("got op %v, expected %v", got.op, expected.op)
					}

					if !slices.Equal(got.args, expected.args) {
						t.Errorf("got args %+v, expected %+v", got.args, expected.args)
						return
					}
				}
			}
		})
	}
}

func TestEvaluateProgram(t *testing.T) {
	var tests = []struct {
		name                       string
		program                    []Instruction
		evaluateToggleInstructions bool
		expectedProductSum         uint
	}{
		{
			"valid example sum (without evaluating toggle instructions)",
			[]Instruction{
				{op: OP_MUL, args: []uint{2, 4}},
				{op: OP_DONT, args: []uint{}},
				{op: OP_MUL, args: []uint{5, 5}},
				{op: OP_MUL, args: []uint{11, 8}},
				{op: OP_DO, args: []uint{}},
				{op: OP_MUL, args: []uint{8, 5}},
			},
			false,
			161,
		},
		{
			"valid example sum (with evaluating toggle instructions)",
			[]Instruction{
				{op: OP_MUL, args: []uint{2, 4}},
				{op: OP_DONT, args: []uint{}},
				{op: OP_MUL, args: []uint{5, 5}},
				{op: OP_MUL, args: []uint{11, 8}},
				{op: OP_DO, args: []uint{}},
				{op: OP_MUL, args: []uint{8, 5}},
			},
			true,
			48,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSum := evaluateProgram(tt.program, tt.evaluateToggleInstructions)

			if gotSum != tt.expectedProductSum {
				t.Errorf("got product sum %v, expected %v", gotSum, tt.expectedProductSum)
			}
		})
	}
}
