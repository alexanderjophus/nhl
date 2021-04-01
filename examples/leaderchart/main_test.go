package main

import (
	"testing"

	"github.com/wcharczuk/go-chart"
)

func testEq(a, b []float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
func TestGetLines(t *testing.T) {
	var scenarios = []struct {
		min, max, count float64
		expected        []float64
	}{
		{
			min:      6,
			max:      10,
			count:    3,
			expected: []float64{0, 4, 8, 12},
		},
		{
			min:      1,
			max:      1,
			count:    1,
			expected: []float64{0, 1},
		},
		{
			min:      -4,
			max:      7,
			count:    4,
			expected: []float64{-4, -1, 2, 5, 8},
		},
	}

	for _, scenario := range scenarios {
		results := getLines(scenario.min, scenario.max, scenario.count)
		if !testEq(scenario.expected, results) {
			t.Errorf("expected %v, found %v", scenario.expected, results)
		}
	}
}

func TestGetFileExtension(t *testing.T) {
	var testCases = []struct {
		input    string
		expected outputFormat
	}{
		{
			input:    "svg",
			expected: outputFormat{name: "svg", renderer: chart.SVG},
		},
		{
			input:    "SVG",
			expected: outputFormat{name: "svg", renderer: chart.SVG},
		},
		{
			input:    "png",
			expected: outputFormat{name: "png", renderer: chart.PNG},
		},
		{
			input:    "unrecognised file extension",
			expected: outputFormat{name: "svg", renderer: chart.SVG},
		},
	}
	for _, testCase := range testCases {
		result := getFileExtension(testCase.input)
		if testCase.expected.name != result.name { //would be good to check renderer too
			t.Errorf("expected %v, got %v", testCase.expected, result)
		}
	}
}
