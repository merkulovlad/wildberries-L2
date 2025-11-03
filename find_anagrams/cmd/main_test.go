package main

import (
	"testing"
	"slices"
	"reflect"
)
type Test struct {
	Name string
	in []string
	out map[string][]string
}
func TestFindAnagrams(t *testing.T) {
	tests := []Test{
		{
			Name: "basic wb sample",
			in:   []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"},
			out: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			Name: "two groups",
			in:   []string{"кот", "ток", "окт", "нос", "сон", "сно"},
			out: map[string][]string{
				"кот": {"кот", "окт", "ток"},
				"нос": {"нос", "сно", "сон"},
			},
		},
		{
			Name: "no anagrams",
			in:   []string{"дом", "лес", "река"},
			out:  map[string][]string{},
		},
	}

	for _, tc := range	 tests {
		t.Run(tc.Name, func(t *testing.T) {
			got := findAnagrams(tc.in)
			if len(got) != len(tc.out) {
				t.Fatalf("expected %d groups, got %d", len(tc.out), len(got))
			}
			for k, wantSlice := range tc.out {
				gotSlice, ok := got[k]
				if !ok {
					t.Fatalf("expected key %q not found", k)
				}
				slices.Sort(gotSlice)
				slices.Sort(wantSlice)
				if !reflect.DeepEqual(gotSlice, wantSlice) {
					t.Fatalf("group %q mismatch\n got: %#v\nwant: %#v", k, gotSlice, wantSlice)
				}
			}
		})
	}
}