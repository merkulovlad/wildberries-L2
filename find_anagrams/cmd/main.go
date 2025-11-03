package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

func findAnagrams(words []string) map[string][]string {
	firstSeen := make(map[string]string)
	anagrams := make(map[string][]string)
	for _, line := range words {
		runes := []rune(line)
		slices.Sort(runes)
		sortedKey := string(runes)
		if v, ok := firstSeen[sortedKey]; !ok {
			firstSeen[sortedKey] = line
		} else {
			if anagrams[v] == nil {
				anagrams[v] = make([]string, 1)
				anagrams[v][0] = v
			}
			anagrams[v] = append(anagrams[v], line)
		}
	}
	for _, v := range anagrams {
		if len(v) > 1 {
			slices.Sort(v)
		}
	}
	return anagrams
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := make([]string, 0, 20)
	anagrams := make(map[string][]string)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	anagrams = findAnagrams(lines)
	keys := make([]string, 0, len(anagrams))
	for k := range anagrams {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Printf("%q: %v\n", k, anagrams[k])
	}
}
