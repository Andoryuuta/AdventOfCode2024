# AdventOfCode2024
Advent of Code 2024 


# Running specific day/solution
```bash
$ go run ./cmd/day1/ ./challenge_data/day1/input_example
Total distance: 11
Similarity score: 31
```

# Testing
```bash
$ go test ./cmd/... -v
=== RUN   TestParseLocationListEmpty
=== RUN   TestParseLocationListEmpty/Empty_input
=== RUN   TestParseLocationListEmpty/Single_valid_pair
=== RUN   TestParseLocationListEmpty/Multi_valid_pair
=== RUN   TestParseLocationListEmpty/Require_specific_separator_-_not_single_space_character
=== RUN   TestParseLocationListEmpty/Require_specific_separator_-_not_tab
=== RUN   TestParseLocationListEmpty/Disallow_empty_lines
=== RUN   TestParseLocationListEmpty/Disallow_signed_numbers
=== RUN   TestParseLocationListEmpty/Disallow_hexideciaml
--- PASS: TestParseLocationListEmpty (0.00s)
    --- PASS: TestParseLocationListEmpty/Empty_input (0.00s)
    --- PASS: TestParseLocationListEmpty/Single_valid_pair (0.00s)
    --- PASS: TestParseLocationListEmpty/Multi_valid_pair (0.00s)
    --- PASS: TestParseLocationListEmpty/Require_specific_separator_-_not_single_space_character (0.00s)
    --- PASS: TestParseLocationListEmpty/Require_specific_separator_-_not_tab (0.00s)
    --- PASS: TestParseLocationListEmpty/Disallow_empty_lines (0.00s)
    --- PASS: TestParseLocationListEmpty/Disallow_signed_numbers (0.00s)
    --- PASS: TestParseLocationListEmpty/Disallow_hexideciaml (0.00s)
=== RUN   TestCalcSimilarityScore
=== RUN   TestCalcSimilarityScore/test_case_0
=== RUN   TestCalcSimilarityScore/test_case_1
--- PASS: TestCalcSimilarityScore (0.00s)
    --- PASS: TestCalcSimilarityScore/test_case_0 (0.00s)
    --- PASS: TestCalcSimilarityScore/test_case_1 (0.00s)
=== RUN   TestCalcListDistance
=== RUN   TestCalcListDistance/test_case_0
=== RUN   TestCalcListDistance/test_case_1
--- PASS: TestCalcListDistance (0.00s)
    --- PASS: TestCalcListDistance/test_case_0 (0.00s)
    --- PASS: TestCalcListDistance/test_case_1 (0.00s)
PASS
ok      github.com/Andoryuuta/AdventOfCode2024/cmd/day1 0.002s
```