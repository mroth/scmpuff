package exec

import (
	"reflect"
	"strings"
	"testing"
)

// Expansion of multiple args at the same time
var testExpandCases = []struct {
	args, expected string
}{
	{"1 3 7", "$e1 $e3 $e7"},
	{"1-3 6", "$e1 $e2 $e3 $e6"},
	{"seven 2-5 1", "seven $e2 $e3 $e4 $e5 $e1"},
}

func TestExpand(t *testing.T) {
	for _, tc := range testExpandCases {
		// split here to emulate what Cobra will pass us but still write tests with
		// normal looking strings
		args := strings.Split(tc.args, " ")
		expected := strings.Split(tc.expected, " ")
		actual := expand(args)
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("ExpandArgs(%v): expected %v, actual %v", tc.args, expected, actual)
		}
	}
}

// Expansion of a single arg, which might still be a range
var testExpandArgCases = []struct {
	arg      string
	expected []string
}{
	{"1", []string{"$e1"}},                 // single digit
	{"1-3", []string{"$e1", "$e2", "$e3"}}, // range
	{"seven", []string{"seven"}},           // no moleste
}

func TestExpandArg(t *testing.T) {
	for _, tc := range testExpandArgCases {
		actual := expandArg(tc.arg)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Fatalf("ExpandArg(%v): expected %v, actual %v", tc.arg, tc.expected, actual)
		}
	}
}
