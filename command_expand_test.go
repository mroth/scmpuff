package main

import (
	"strings"
	"testing"
)

var testExpandArgsCases = []struct {
	args, expected string
}{
	{"1 3 7", "$e1 $e3 $e7"},
	{"1-3 6", "$e1 $e2 $e3 $e6"},
	{"seven 2-5 1", "seven $e2 $e3 $e4 $e5 $e1"},
}

func TestExpandArgs(t *testing.T) {
	for _, tc := range testExpandArgsCases {
		// split here to emulate what Cobra will pass us but still write test with
		// normal strings
		args := strings.Split(tc.args, " ")
		actual := ExpandArgs(args)
		if actual != tc.expected {
			t.Fatalf("ExpandArgs(%v): expected %v, actual %v", tc.args, tc.expected, actual)
		}
	}
}

var testExpandArgCases = []struct {
	arg, expected string
}{
	{"1", "$e1"},           // single digit
	{"1-3", "$e1 $e2 $e3"}, // range
	{"seven", "seven"},     // no moleste
}

func TestExpandArg(t *testing.T) {
	for _, tc := range testExpandArgCases {
		actual := ExpandArg(tc.arg)
		if actual != tc.expected {
			t.Fatalf("ExpandArg(%v): expected %v, actual %v", tc.arg, tc.expected, actual)
		}
	}
}
