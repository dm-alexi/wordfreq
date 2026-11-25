package wordfreq

import (
	"errors"
	"io"
	"maps"
	"strings"
	"testing"
	"testing/iotest"
)

func TestCount_General(t *testing.T) {
	tests := []struct {
		name string
		text string
		want map[string]int
	}{
		{
			name: "Words lowercased and counted, whitespaces and puntuation stripped",
			text: `	How much wood would a woodchuck chuck if a WoodChuck could
		 CHUCK WooD? `,
			want: map[string]int{
				"how":       1,
				"much":      1,
				"wood":      2,
				"would":     1,
				"could":     1,
				"a":         2,
				"woodchuck": 2,
				"chuck":     2,
				"if":        1,
			},
		},
		{
			name: "Non-letter character break words",
			text: "Can't can’T+ fifty-fifty=fifty-fifty covid-19!!! take'em rock'n'roll well--known. B2C --Hello-",
			want: map[string]int{
				"can":   2,
				"t":     2,
				"fifty": 4,
				"take":  1,
				"em":    1,
				"covid": 1,
				"rock":  1,
				"n":     1,
				"roll":  1,
				"well":  1,
				"b":     1,
				"c":     1,
				"known": 1,
				"hello": 1,
			},
		},
		{
			name: "empty input",
			text: "",
			want: map[string]int{},
		},
		{
			name: "Multilanguage",
			text: "ПРИВЕТ! 東京",
			want: map[string]int{
				"привет": 1,
				"東京":     1,
			},
		},
		{
			name: "Buffer border word",
			text: strings.Repeat(" ", 4095) + "привет привет привет",
			want: map[string]int{
				"привет": 3,
			},
		},
		{
			name: "Mark character",
			text: "cafe\u0301\u0301",
			want: map[string]int{
				"café́": 1,
			},
		},
		{
			name: "Invalid encoding",
			text: "x\xffword word",
			want: map[string]int{
				"x":    1,
				"word": 2,
			},
		},
		{
			name: "Invalid encoding, trailing",
			text: "	abc\xe2\x80",
			want: map[string]int{"abc": 1},
		},
		{
			name: "Leading combining mark",
			text: " ́de",
			want: map[string]int{"de": 1},
		},
		{
			name: "Long line",
			text: strings.Repeat("hi ", 350000),
			want: map[string]int{"hi": 350000},
		},
		{
			name: "Long word",
			text: strings.Repeat("h", MaxWordLength),
			want: map[string]int{strings.Repeat("h", MaxWordLength): 1},
		},
	}

	readers := []struct {
		name   string
		reader func(string) io.Reader
	}{
		{
			name:   "Default reader",
			reader: func(s string) io.Reader { return strings.NewReader(s) },
		},
		{
			name:   "DataErr reader",
			reader: func(s string) io.Reader { return iotest.DataErrReader(strings.NewReader(s)) },
		},
		{
			name:   "One byte reader",
			reader: func(s string) io.Reader { return iotest.OneByteReader(strings.NewReader(s)) },
		},
	}

	for _, rr := range readers {
		for _, tt := range tests {
			t.Run(rr.name+" "+tt.name, func(t *testing.T) {
				got, err := Count(rr.reader(tt.text))
				assertNoError(t, tt.text, err)
				assertMapsEqual(t, tt.text, got, tt.want)
			})
		}
	}
}

func TestCount_Error(t *testing.T) {
	tests := []struct {
		name   string
		reader io.Reader
		want   error
	}{
		{
			name:   "Token too long",
			reader: strings.NewReader(strings.Repeat("x", MaxWordLength+1)),
			want:   ErrTooLong,
		},
		{
			name:   "Reader fault",
			reader: iotest.TimeoutReader(strings.NewReader("test test")),
			want:   iotest.ErrTimeout,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Count(tt.reader)
			if got != nil {
				t.Errorf("Count got non-nil value, want nil")
			}
			if !errors.Is(err, tt.want) {
				t.Errorf("Count error = %v, want %v", err, tt.want)
			}
		})
	}
}

func assertNoError(t *testing.T, input string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Count(%q), no error expected, got %v", truncate(input), err)
	}
}

const maxMapDiff = 10

func assertMapsEqual(t *testing.T, input string, got, want map[string]int) {
	t.Helper()
	counter := 0
	if !maps.Equal(got, want) {
		t.Errorf("Count(%q) mismatch:", truncate(input))
		for k, v := range want {
			if g, ok := got[k]; !ok {
				t.Errorf("\tkey: %q, missing, want %v", truncate(k), v)
				if counter++; counter == maxMapDiff {
					return
				}
			} else if g != v {
				t.Errorf("\tkey: %q, got %v, want %v", truncate(k), got[k], v)
				if counter++; counter == maxMapDiff {
					return
				}
			}
		}
		for k, v := range got {
			if _, ok := want[k]; !ok {
				t.Errorf("\tkey: %q, got %v, unexpected", truncate(k), v)
				if counter++; counter == maxMapDiff {
					return
				}
			}
		}
	}
}

func truncate(s string) string {
	if len([]rune(s)) > 32 {
		s = string([]rune(s)[:29]) + "..."
	}
	return s
}
