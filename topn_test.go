package wordfreq

import (
	"slices"
	"testing"
)

func TestTop(t *testing.T) {
	tests := []struct {
		name string
		dict map[string]int
		n    int
		want []WordCount
	}{
		{
			name: "List sorted and trimmed",
			dict: map[string]int{
				"how":       1,
				"much":      1,
				"wood":      2,
				"would":     1,
				"could":     1,
				"a":         2,
				"woodchuck": 2,
				"chuck":     2,
				"if":        1,
				"first":     14,
			},
			n:    7,
			want: []WordCount{{"first", 14}, {"a", 2}, {"chuck", 2}, {"wood", 2}, {"woodchuck", 2}, {"could", 1}, {"how", 1}},
		},
		{
			name: "N exceeds dict size",
			dict: map[string]int{
				"how":  1,
				"much": 1,
			},
			n:    7,
			want: []WordCount{{"how", 1}, {"much", 1}},
		},
		{
			name: "Negarive N",
			dict: map[string]int{
				"how":  1,
				"much": 1,
			},
			n:    -2,
			want: []WordCount{},
		},
		{
			name: "Zero N",
			dict: map[string]int{
				"how":  1,
				"much": 1,
			},
			n:    0,
			want: []WordCount{},
		},
		{
			name: "Empty dict",
			dict: map[string]int{},
			n:    7,
			want: []WordCount{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Top(tt.dict, tt.n)
			if !slices.Equal(got, tt.want) {
				t.Errorf("Top(%v, %v) = %v, want %v", tt.dict, tt.n, got, tt.want)
			}
		})
	}
}
