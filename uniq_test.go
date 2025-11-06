package main

import (
	"reflect"
	"testing"
)

func TestUniq_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		opts     Options
		expected []string
	}{
		{
			name:     "basic unique",
			input:    []string{"a", "b", "a", "c", "b"},
			opts:     Options{},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with count",
			input:    []string{"a", "b", "a", "c", "b"},
			opts:     Options{Count: true},
			expected: []string{"2 a", "2 b", "1 c"},
		},
		{
			name:     "duplicate only",
			input:    []string{"a", "b", "a", "c", "b"},
			opts:     Options{Duplicate: true},
			expected: []string{"a", "b"},
		},
		{
			name:     "unique only",
			input:    []string{"a", "b", "a", "c", "b"},
			opts:     Options{Unique: true},
			expected: []string{"c"},
		},
		{
			name:     "empty input",
			input:    []string{},
			opts:     Options{},
			expected: []string{},
		},
		{
			name:     "single line",
			input:    []string{"hello"},
			opts:     Options{},
			expected: []string{"hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Uniq(tt.input, tt.opts)
			if err != nil {
				t.Errorf("Uniq() error = %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Uniq() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniq_IgnoreCase(t *testing.T) {
	input := []string{"Apple", "apple", "BANANA", "banana", "Cherry"}
	opts := Options{IgnoreCase: true}
	expected := []string{"Apple", "BANANA", "Cherry"}

	result, err := Uniq(input, opts)
	if err != nil {
		t.Errorf("Uniq() error = %v", err)
		return
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Uniq() = %v, want %v", result, expected)
	}
}

func TestUniq_NumFields(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		opts     Options
		expected []string
	}{
		{
			name: "skip 1 field",
			input: []string{
				"1 apple red",
				"2 apple green",
				"3 banana yellow",
				"4 apple red",
			},
			opts:     Options{NumFields: 1},
			expected: []string{"1 apple red", "2 apple green", "3 banana yellow"},
		},
		{
			name: "skip 2 fields",
			input: []string{
				"1 2 apple",
				"3 4 apple",
				"5 6 banana",
			},
			opts:     Options{NumFields: 2},
			expected: []string{"1 2 apple", "5 6 banana"},
		},
		{
			name: "skip more fields than available",
			input: []string{
				"a b",
				"c d",
			},
			opts:     Options{NumFields: 5},
			expected: []string{"a b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Uniq(tt.input, tt.opts)
			if err != nil {
				t.Errorf("Uniq() error = %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Uniq() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniq_NumChars(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		opts     Options
		expected []string
	}{
		{
			name: "skip 5 chars",
			input: []string{
				"apple123",
				"apple456",
				"banana789",
			},
			opts:     Options{NumChars: 5},
			expected: []string{"apple123", "apple456", "banana789"},
		},
		{
			name: "skip more chars than available",
			input: []string{
				"abc",
				"def",
			},
			opts:     Options{NumChars: 10},
			expected: []string{"abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Uniq(tt.input, tt.opts)
			if err != nil {
				t.Errorf("Uniq() error = %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Uniq() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniq_CombinedOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		opts     Options
		expected []string
	}{
		{
			name: "ignore case with count",
			input: []string{"Apple", "apple", "Banana", "banana"},
			opts:  Options{Count: true, IgnoreCase: true},
			expected: []string{
				"2 Apple",
				"2 Banana",
			},
		},
		{
			name: "fields and chars",
			input: []string{
				"field1 field2 abc123",
				"field1 field2 def456",
				"field1 field3 abc123",
			},
			opts:     Options{NumFields: 2, NumChars: 3},
			expected: []string{"field1 field2 abc123", "field1 field2 def456"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Uniq(tt.input, tt.opts)
			if err != nil {
				t.Errorf("Uniq() error = %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Uniq() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniq_Validation(t *testing.T) {
	tests := []struct {
		name string
		opts Options
		want error
	}{
		{
			name: "mutually exclusive c and d",
			opts: Options{Count: true, Duplicate: true},
			want: ErrMutuallyExclusive,
		},
		{
			name: "mutually exclusive c and u",
			opts: Options{Count: true, Unique: true},
			want: ErrMutuallyExclusive,
		},
		{
			name: "mutually exclusive d and u",
			opts: Options{Duplicate: true, Unique: true},
			want: ErrMutuallyExclusive,
		},
		{
			name: "negative num fields",
			opts: Options{NumFields: -1},
			want: ErrNegativeNumber,
		},
		{
			name: "negative num chars",
			opts: Options{NumChars: -5},
			want: ErrNegativeNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Uniq([]string{"test"}, tt.opts)
			if err != tt.want {
				t.Errorf("Uniq() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestSkipFields(t *testing.T) {
	tests := []struct {
		input    string
		n        int
		expected string
	}{
		{"a b c d", 2, "c d"},
		{"a b c d", 0, "a b c d"},
		{"a b", 3, ""},
		{"   a   b   c  ", 1, "b c"},
		{"", 1, ""},
		{"single", 1, ""},
		{"one two three", 1, "two three"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := skipFields(tt.input, tt.n)
			if result != tt.expected {
				t.Errorf("skipFields(%q, %d) = %q, want %q", tt.input, tt.n, result, tt.expected)
			}
		})
	}
}