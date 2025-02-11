package gostr_test

import (
	"reflect"
	"testing"

	"github.com/mailund/gostr/gostr"
	"github.com/mailund/gostr/testutils"
)

func Test_BorderarrayBasics(t *testing.T) {
	tests := []struct {
		name string
		x    string
		want []int
	}{
		{"(empty string)", "", []int{}},
		{"a", "a", []int{0}},
		{"aaa", "aaa", []int{0, 1, 2}},
		{"aaaba", "aaaba", []int{0, 1, 2, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gostr.Borderarray(tt.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Borderarray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_StrictBorderarrayBasics(t *testing.T) {
	tests := []struct {
		name string
		x    string
		want []int
	}{
		{"(empty string)", "", []int{}},
		{"a", "a", []int{0}},
		{"aaa", "aaa", []int{0, 0, 2}},
		{"aaaba", "aaaba", []int{0, 0, 2, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gostr.StrictBorderarray(tt.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrictBorderarray() = %v, want %v", got, tt.want)
			}
		})
	}
}

// FIXME: check if the border is the *longest* as well
func checkBorders(t *testing.T, x string, ba []int) {
	t.Helper()

	for i, b := range ba {
		if b != 0 && x[:b] != x[i-b+1:i+1] {
			t.Errorf(`x[:%d] == %q is not a border of %q`, b, x[:b], x[:i+1])
			t.Fatalf(`x = %q, ba = %v`, x, ba)
		}
	}
}

func Test_Borderarray(t *testing.T) {
	rng := testutils.NewRandomSeed(t)
	testutils.GenerateTestStrings(10, 20, rng,
		func(x string) {
			checkBorders(t, x, gostr.Borderarray(x))
		})
}

func checkStrict(t *testing.T, x string, ba []int) bool {
	t.Helper()

	for i, b := range ba[:len(ba)-1] {
		if b > 0 && x[b] == x[i+1] {
			t.Errorf(`x[:%d] == %q[%q] is not a strict border of %q[%q]`, b, x[:b], x[b], x[:i+1], x[i+1])
			t.Errorf(`x[%d] == %q == x[%d+1] (should be different)`, b, x[b], i)
			t.Fatalf(`x = %q, ba = %v`, x, ba)

			return false
		}
	}

	return true
}

func Test_StrictBorderarray(t *testing.T) {
	rng := testutils.NewRandomSeed(t)
	testutils.GenerateTestStrings(10, 20, rng,
		func(x string) {
			ba := gostr.StrictBorderarray(x)
			checkBorders(t, x, ba)
			checkStrict(t, x, ba)
		})
}
