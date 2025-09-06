package anagram

import (
	"sort"
	"strings"
)

func Search(words []string) map[string][]string {
	anagramSets := make(map[string][]string)
	firstOccurrence := make(map[string]string)

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		sorted := sortRunes(lowerWord)

		if first, exists := firstOccurrence[sorted]; exists {
			anagramSets[first] = append(anagramSets[first], lowerWord)
		} else {
			firstOccurrence[sorted] = lowerWord
			anagramSets[lowerWord] = []string{lowerWord}
		}
	}

	result := make(map[string][]string)
	for key, words := range anagramSets {
		if len(words) > 1 {
			sort.Strings(words)
			result[key] = words
		}
	}

	return result
}

func sortRunes(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}
