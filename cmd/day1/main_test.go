package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseLocationListEmpty(t *testing.T) {
	var tests = []struct {
		name                         string
		input                        string
		expectedList1, expectedList2 *[]uint64
		expectedError                bool
	}{
		{
			"Empty input",
			"",
			nil,
			nil,
			false,
		},
		{
			"Single valid pair",
			"1   2",
			&[]uint64{1},
			&[]uint64{2},
			false,
		},
		{
			"Multi valid pair",
			"1   2\n3   4\n5   6",
			&[]uint64{1, 3, 5},
			&[]uint64{2, 4, 6},
			false,
		},
		{
			"Require specific separator - not single space character",
			"1 6",
			nil,
			nil,
			true,
		},
		{
			"Require specific separator - not tab",
			"1\t6",
			nil,
			nil,
			true,
		},
		{
			"Disallow empty lines",
			"1   2\n\n5   6",
			nil,
			nil,
			true,
		},
		{
			"Disallow signed numbers",
			"-1   2",
			nil,
			nil,
			true,
		},
		{
			"Disallow hexideciaml",
			"0x1   2",
			nil,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			gotList1, gotList2, gotErr := parseLocationList(reader)

			if tt.expectedList1 != nil && !reflect.DeepEqual(gotList1, *tt.expectedList1) {
				t.Errorf("got %v, expected %v", gotList1, *tt.expectedList1)
			}

			if tt.expectedList2 != nil && !reflect.DeepEqual(gotList2, *tt.expectedList2) {
				t.Errorf("got %v, expected %v", gotList2, *tt.expectedList2)
			}

			if tt.expectedError && gotErr == nil {
				t.Errorf("got %v, expected %v", gotErr, nil)
			}
		})
	}

}

func TestCalcSimilarityScore(t *testing.T) {
	var tests = []struct {
		list1, list2 []uint64
		expected     uint64
	}{
		{
			[]uint64{},
			[]uint64{},
			0,
		},
		{
			[]uint64{3, 4, 2, 1, 3, 3},
			[]uint64{4, 3, 5, 3, 9, 3},
			31,
		},
	}

	for idx, tt := range tests {
		testname := fmt.Sprintf("test_case_%v", idx)
		t.Run(testname, func(t *testing.T) {
			result := calcSimilarityScore(tt.list1, tt.list2)
			if result != tt.expected {
				t.Errorf("got %d, expected %d", result, tt.expected)
			}
		})
	}
}

func TestCalcListDistance(t *testing.T) {
	var tests = []struct {
		list1, list2 []uint64
		expected     uint64
	}{
		{
			[]uint64{},
			[]uint64{},
			0,
		},
		{
			[]uint64{3, 4, 2, 1, 3, 3},
			[]uint64{4, 3, 5, 3, 9, 3},
			11,
		},
	}

	for idx, tt := range tests {
		testname := fmt.Sprintf("test_case_%v", idx)
		t.Run(testname, func(t *testing.T) {
			result := calcListDistance(tt.list1, tt.list2)
			if result != tt.expected {
				t.Errorf("got %d, expected %d", result, tt.expected)
			}
		})
	}
}
