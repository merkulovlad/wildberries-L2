package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalid = errors.New("invalid value")

func unpackString(s string) (string, error) {

	var lastRune rune
	var havePrev bool
	var output strings.Builder
	var alrDigit bool
	var esc bool

	for _, r := range s {
		if !havePrev && unicode.IsDigit(r) {
			return "", ErrInvalid
		} else if havePrev && unicode.IsDigit(r) && !esc {
			if alrDigit {
				return "", ErrInvalid
			}
			digit := int(r - '0')
			if digit < 0 || digit > 9 {
				return "", ErrInvalid
			}

			for i := 0; i < digit-1; i++ {
				output.WriteRune(lastRune)
			}
			alrDigit = true
		} else if r == '\\' && !esc {
			esc = true
			havePrev = true
			alrDigit = false
		} else if esc {
			output.WriteRune(r)
			havePrev = true
			lastRune = r
			esc = false
			alrDigit = false
		} else {
			output.WriteRune(r)
			havePrev = true
			lastRune = r
			alrDigit = false
		}
	}
	if esc {
		return "", ErrInvalid
	}
	return output.String(), nil
}

func main() {
	first, err := unpackString("qwe\\4\\5")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(first)
	}
}
