package expand

import (
	"reflect"
	"strings"
	"testing"
)

var testExpandCases = []struct {
	args, expected string
}{
	{"1 3 7", "$e1\t$e3\t$e7"},
	{"1-3 6", "$e1\t$e2\t$e3\t$e6"},
	{"seven 2-5 1", "seven\t$e2\t$e3\t$e4\t$e5\t$e1"},
}

func TestExpand(t *testing.T) {
	for _, tc := range testExpandCases {
		// split here to emulate what Cobra will pass us but still write test with
		// normal strings
		args := strings.Split(tc.args, " ")
		actual := expand(args)
		if actual != tc.expected {
			t.Fatalf("ExpandArgs(%v): expected %v, actual %v", tc.args, tc.expected, actual)
		}
	}
}

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
