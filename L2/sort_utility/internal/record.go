package internal

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Record struct {
	OriginalLine      string
	TrimmedKey        string
	NumericValue      float64
	NumValid          bool
	MonthValue        int
	MonthValid        bool
	HumanNumericValue float64
	HumanNumericValid bool
}

var monthMap = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4,
	"may": 5, "jun": 6, "jul": 7, "aug": 8,
	"sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

func NewRecord(line string, opts Options) *Record {
	key := line
	fields := strings.Split(line, "\t")

	if opts.Column != 0 && opts.Column <= len(fields) {
		key = fields[opts.Column-1]
	}

	if opts.IgnoreTrailing {
		key = strings.TrimRightFunc(key, unicode.IsSpace)
	}

	numValid := true
	numValue := 0.0

	if opts.Numeric {
		num, err := strconv.ParseFloat(key, 64)
		if err != nil || math.IsNaN(num) || math.IsInf(num, 0) {
			numValid = false
		} else {
			numValue = num
		}
	}

	monthValue := 0
	monthValid := true

	if opts.Month {
		if len(key) < 3 {
			monthValid = false
		} else {
			lowerKey := strings.ToLower(key[0:3])
			if val, ok := monthMap[lowerKey]; ok {
				monthValue = val
			} else {
				monthValid = false
			}
		}
	}

	humanNumValue := 0.0
	humanNumValid := true

	if opts.HumanNumeric {
		humanNum, err := parseHumanNumeric(key)
		if err != nil {
			humanNumValid = false
		} else {
			humanNumValue = humanNum
		}
	}

	return &Record{
		OriginalLine:      line,
		TrimmedKey:        key,
		NumericValue:      numValue,
		NumValid:          numValid,
		MonthValue:        monthValue,
		MonthValid:        monthValid,
		HumanNumericValue: humanNumValue,
		HumanNumericValid: humanNumValid,
	}
}

func parseHumanNumeric(s string) (float64, error) {
	if len(s) == 0 {
		return 0, strconv.ErrSyntax
	}

	multiplier := 1.0
	lastChar := unicode.ToLower(rune(s[len(s)-1]))

	const (
		kib = 1024.0
		mib = kib * 1024
		gib = mib * 1024
		tib = gib * 1024
		pib = tib * 1024
		eib = pib * 1024
		zib = eib * 1024
		yib = zib * 1024
	)

	switch lastChar {
	case 'k':
		multiplier = kib
		s = s[:len(s)-1]
	case 'm':
		multiplier = mib
		s = s[:len(s)-1]
	case 'g':
		multiplier = gib
		s = s[:len(s)-1]
	case 't':
		multiplier = tib
		s = s[:len(s)-1]
	case 'p':
		multiplier = pib
		s = s[:len(s)-1]
	case 'e':
		multiplier = eib
		s = s[:len(s)-1]
	case 'z':
		multiplier = zib
		s = s[:len(s)-1]
	case 'y':
		multiplier = yib
		s = s[:len(s)-1]
	}

	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return value * multiplier, nil
}
