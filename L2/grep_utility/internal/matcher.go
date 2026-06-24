package internal

import (
	"fmt"
	"regexp"
	"strings"
)

type Matcher func(line string) bool

func BuildMatcher(pattern string, opts Options) (Matcher, error) {
	if opts.ExactMatch {
		if opts.IgnoreCase {
			pattern = strings.ToLower(pattern)

			return func(line string) bool {
				return strings.Contains(strings.ToLower(line), pattern)
			}, nil
		}

		return func(line string) bool {
			return strings.Contains(line, pattern)
		}, nil
	}

	if opts.IgnoreCase {
		pattern = "(?i)" + pattern
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		// In case of regex compilation error, return a matcher that matches nothing
		return nil, err
	}

	return func(line string) bool {
		return regex.MatchString(line)
	}, nil
}

func GrepLines(lines []string, pattern string, opts Options) ([]*Line, error) {
	matcher, err := BuildMatcher(pattern, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to build matcher: %w", err)
	}

	matchedLines := make([]*Line, 0, len(lines))

	for i, text := range lines {
		isMatch := matcher(text)
		if opts.InvertMatch {
			isMatch = !isMatch
		}

		line := NewLine(i+1, text, isMatch)
		matchedLines = append(matchedLines, line)
	}

	applyContext(matchedLines, opts)
	answer := make([]*Line, 0, len(matchedLines))

	for _, line := range matchedLines {
		if line.Printed {
			answer = append(answer, line)
		}
	}

	return answer, nil
}

func applyContext(lines []*Line, opts Options) {
	totalLines := len(lines)

	for i, line := range lines {
		if line.Match {
			// Apply before context
			for j := 1; j <= opts.LinesBeforeFound; j++ {
				if i-j >= 0 && !lines[i-j].Printed {
					lines[i-j].Printed = true
				}
			}
			// Apply after context
			for j := 1; j <= opts.LinesAfterFound; j++ {
				if i+j < totalLines && !lines[i+j].Printed {
					lines[i+j].Printed = true
				}
			}

			line.Printed = true
		}
	}
}
