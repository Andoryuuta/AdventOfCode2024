package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseReportList(t *testing.T) {
	var tests = []struct {
		name            string
		input           string
		expectedReports *[]Report
		expectedError   bool
	}{
		{
			"single valid report, level pair",
			"1 2",
			&[]Report{
				{1, 2},
			},
			false,
		},
		{
			"Single valid report, multiple levels",
			"1 2 3 4 5 6 7 8 9 10",
			&[]Report{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			false,
		},
		{
			"multiple valid reports, level pairs",
			"1 2\n3 4",
			&[]Report{
				{1, 2},
				{3, 4},
			},
			false,
		},
		{
			"multiple valid reports, multiple levels",
			"1 2 3 2 1\n3 4 5 6 7",
			&[]Report{
				{1, 2, 3, 2, 1},
				{3, 4, 5, 6, 7},
			},
			false,
		},
		{
			"allow trailing newline",
			"1 2\n",
			&[]Report{
				{1, 2},
			},
			false,
		},
		{
			"no empty input",
			"",
			nil,
			true,
		},
		{
			"requires two levels at minimum",
			"1",
			nil,
			true,
		},
		{
			"no hexidecimal",
			"0x1 2 3 2 1",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			gotReports, gotErr := parseReportList(reader)

			if tt.expectedError && gotErr == nil {
				t.Errorf("got %v, expected nil", gotErr)
				return
			} else if !tt.expectedError && gotErr != nil {
				t.Errorf("got %v, expected !nil", gotErr)
				return
			}

			if tt.expectedReports != nil {
				if len(gotReports) != len(*tt.expectedReports) {
					t.Errorf("unexpected report count, got %v, expected %v", len(gotReports), len(*tt.expectedReports))
					return
				}

				for idx, expectedReport := range *tt.expectedReports {
					if !reflect.DeepEqual(gotReports[idx], expectedReport) {
						t.Errorf("got %v, expected %v", gotReports[idx], expectedReport)
						return
					}
				}
			}
		})
	}
}

func TestIsReportSafe(t *testing.T) {
	var tests = []struct {
		report                 Report
		problemDampenerEnabled bool
		expected               bool
	}{
		// Part 1 examples
		{
			Report{7, 6, 4, 2, 1},
			false,
			true,
		},
		{
			Report{1, 2, 7, 8, 9},
			false,
			false,
		},
		{
			Report{9, 7, 6, 2, 1},
			false,
			false,
		},
		{
			Report{1, 3, 2, 4, 5},
			false,
			false,
		},
		{
			Report{8, 6, 4, 4, 1},
			false,
			false,
		},
		{
			Report{1, 3, 6, 7, 9},
			false,
			true,
		},

		// Part 2 examples
		{
			Report{7, 6, 4, 2, 1},
			true,
			true,
		},
		{
			Report{1, 2, 7, 8, 9},
			true,
			false,
		},
		{
			Report{9, 7, 6, 2, 1},
			true,
			false,
		},
		{
			Report{1, 3, 2, 4, 5},
			true,
			true,
		},
		{
			Report{8, 6, 4, 4, 1},
			true,
			true,
		},
		{
			Report{1, 3, 6, 7, 9},
			true,
			true,
		},
	}

	for idx, tt := range tests {
		testname := fmt.Sprintf("test_case_%v", idx)
		t.Run(testname, func(t *testing.T) {
			result := isReportSafe(tt.report, tt.problemDampenerEnabled)
			if result != tt.expected {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCalcSafeReports(t *testing.T) {
	var tests = []struct {
		reports                []Report
		problemDampenerEnabled bool
		expected               uint64
	}{
		{
			[]Report{
				{7, 6, 4, 2, 1},
				{1, 2, 7, 8, 9},
				{9, 7, 6, 2, 1},
				{1, 3, 2, 4, 5},
				{8, 6, 4, 4, 1},
				{1, 3, 6, 7, 9},
			},
			false,
			2,
		},
		{
			[]Report{
				{7, 6, 4, 2, 1},
				{1, 2, 7, 8, 9},
				{9, 7, 6, 2, 1},
				{1, 3, 2, 4, 5},
				{8, 6, 4, 4, 1},
				{1, 3, 6, 7, 9},
			},
			true,
			4,
		},
	}

	for idx, tt := range tests {
		testname := fmt.Sprintf("test_case_%v", idx)
		t.Run(testname, func(t *testing.T) {
			result := calcSafeReports(tt.reports, tt.problemDampenerEnabled)
			if result != tt.expected {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}
