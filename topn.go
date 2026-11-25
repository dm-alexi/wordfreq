package wordfreq

import (
	"cmp"
	"slices"
	"strings"
)

// A WordCount represents a word and its count
type WordCount struct {
	Word  string
	Count int
}

// Top returns a sorted list of n most common words, ties are broken alphabetically.
// Full list returned if n is greater than the size of the words collection;
// n <= 0 results in empty list.
func Top(counts map[string]int, n int) []WordCount {
	if n <= 0 {
		return []WordCount{}
	}
	wordList := make([]WordCount, 0, len(counts))
	for k, v := range counts {
		wordList = append(wordList, WordCount{Word: k, Count: v})
	}
	slices.SortFunc(wordList, func(a, b WordCount) int {
		return cmp.Or(cmp.Compare(b.Count, a.Count), strings.Compare(a.Word, b.Word))
	})
	if n < len(wordList) {
		wordList = wordList[:n]
	}
	return wordList
}
