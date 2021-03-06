package pretty

import (
	"testing"
)

type difftest struct {
	a   interface{}
	b   interface{}
	exp []string
}

type S struct {
	A int
	S *S
	I interface{}
	C []int
	u int
}

var diffs = []difftest{
	{a: nil, b: nil},
	{a: S{A: 1}, b: S{A: 1}},
	{a: S{u: 1}, b: S{u: 1}},

	{0, "", []string{`int != string`}},
	{0, 1, []string{`0 != 1`}},
	{S{}, new(S), []string{`pretty.S != *pretty.S`}},
	{"a", "b", []string{`"a" != "b"`}},
	{S{}, S{A: 1}, []string{`A: 0 != 1`}},
	{new(S), &S{A: 1}, []string{`A: 0 != 1`}},
	{S{u: 1}, S{u: 2}, []string{`unexported fields don't match`}},
	{S{S: new(S)}, S{S: &S{A: 1}}, []string{`S.A: 0 != 1`}},
	{S{}, S{I: 0}, []string{`I: nil != 0`}},
	{S{I: 1}, S{I: "x"}, []string{`I: int != string`}},
	{S{}, S{C: []int{1}}, []string{`C: []int[0] != []int[1]`}},
	{S{C: []int{}}, S{C: []int{1}}, []string{`C: []int[0] != []int[1]`}},
	{S{C: []int{1, 2, 3}}, S{C: []int{1, 2, 4}}, []string{`C[2]: 3 != 4`}},
	{S{}, S{A: 1, S: new(S)}, []string{`A: 0 != 1`, `S: nil != &{0 <nil> <nil> [] 0}`}},
}

func TestDiff(t *testing.T) {
	for _, tt := range diffs {
		got := Diff(tt.a, tt.b)
		eq := len(got) == len(tt.exp)
		if eq {
			for i := range got {
				eq = eq && got[i] == tt.exp[i]
			}
		}
		if !eq {
			t.Errorf("diffing % #v", tt.a)
			t.Errorf("with    % #v", tt.b)
			diffdiff(t, got, tt.exp)
			continue
		}
	}
}

func diffdiff(t *testing.T, got, exp []string) {
	minus(t, "unexpected:", got, exp)
	minus(t, "missing:", exp, got)
}

func minus(t *testing.T, s string, a, b []string) {
	var i, j int
	for i = 0; i < len(a); i++ {
		for j = 0; j < len(b); j++ {
			if a[i] == b[j] {
				break
			}
		}
		if j == len(b) {
			t.Error(s, a[i])
		}
	}
}
