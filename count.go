// Package wordfreq implements word counter and topN words operation on a io.Reader.
package wordfreq

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

// MaxWordLength is a maximum length of a word in bytes
const MaxWordLength = 1600

// ErrTooLong is returned when a word exceeds maximum word length
var ErrTooLong = errors.New("wordfreq: word too long")

// Count counts words from reader. Maximum word length (in bytes) is MaxWordLength.
// Words are counted in their lower case form, full case folding is not performed.
// On error returns nil map, partial counts are discarded.
func Count(r io.Reader) (map[string]int, error) {
	m := make(map[string]int)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(nil, MaxWordLength+1)
	scanner.Split(scanWords)
	for scanner.Scan() {
		m[strings.ToLower(scanner.Text())]++
	}
	if err := scanner.Err(); err != nil {
		if errors.Is(err, bufio.ErrTooLong) {
			return nil, ErrTooLong
		}
		return nil, fmt.Errorf("wordfreq: %w", err)
	}
	return m, nil
}

// scanWords is a bufio.Split function used to split input by words.
// A word is made of Unicode letters. Intraword punctuation (hyphens, apostrophes), numbers, invalid characters break words.
// Unicode marks considered part of the word. Words longer than MaxWordLength are not supported.
func scanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	for width := 0; start < len(data); start += width {
		if !utf8.FullRune(data[start:]) {
			return start, nil, nil
		}
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if unicode.IsLetter(r) {
			break
		}
	}
	for width, i := 0, start; i < len(data); i += width {
		if !utf8.FullRune(data[i:]) {
			if atEOF {
				return len(data), data[start:i], nil
			}
			return start, nil, nil
		}
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if !(unicode.IsLetter(r) || unicode.IsMark(r)) {
			return i + width, data[start:i], nil
		}
	}
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	return start, nil, nil
}
