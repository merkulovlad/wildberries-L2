package main

import (
	"errors"
	"testing"
)

func TestUnpacking(t *testing.T) {
	type tc struct {
		name string
		in   string
		out  string
		err  error
	}
	tests := []tc{
		{name: "empty", in: "", out: "", err: nil},
		{name: "simple", in: "abcd", out: "abcd", err: nil},
		{name: "only number", in: "45", out: "", err: ErrInvalid},
		{name: "good example", in: "a4bc2d5e", out: "aaaabccddddde", err: nil},
		{name: "handle \\", in: "qwe\\45", out: "qwe44444", err: nil},
		{name: "another handle \\", in: "qwe\\4\\5", out: "qwe45", err: nil},
		{name: "leading digit", in: "3abc", out: "", err: ErrInvalid},
		{name: "trailing backslash", in: `abc\`, out: "", err: ErrInvalid},
		{name: "escaped backslash then count", in: `\\3`, out: `\\\`, err: nil},
		{name: "escaped digit then count", in: `x\93`, out: "x999", err: nil},
		{name: "escaped letter then count", in: `a\b3`, out: "abbb", err: nil},
		{name: "emoji repeat", in: "😺3", out: "😺😺😺", err: nil},
		{name: "zero repeat", in: "a0", out: "a", err: nil},
		{name: "two digits in a row (policy)", in: "a23", out: "", err: ErrInvalid},
		{name: "non-ascii digit", in: "a\uFF15", out: "", err: ErrInvalid},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := unpackString(test.in)
			if test.err != nil {
				if err == nil {
					t.Errorf("unpackString(%q) expected error, got nothing", test.in)
				} else {
					if !errors.Is(err, test.err) {
						t.Errorf("unpackString(%q) expected error %q, got %q", test.in, err, test.err)
					}
				}
				return
			}
			if err != nil {
				t.Errorf("unpackString(%q) expected no error, got %q", test.in, err)
			}
			if got != test.out {
				t.Errorf("unpackString(%q) expected %q, got %q", test.in, test.out, got)
			}
		})
	}

}
