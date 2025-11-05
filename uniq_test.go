package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUniq_IgnoreCase(t *testing.T) {
	input := []string{"Apple", "apple", "BANANA", "banana", "Cherry"}
	opts := Options{IgnoreCase: true}

	expected := []string{"Apple", "BANANA", "Cherry"}

	result, err := Uniq(input, opts)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

//pltcmcmcmmc
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
            // После пропуска 1 поля: "apple red", "apple green", "banana yellow", "apple red"
            // Группы: "apple red" (2 раза), "apple green" (1 раз), "banana yellow" (1 раз)
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
            // После пропуска 2 полей: "apple", "apple", "banana"  
            // Группы: "apple" (2 раза), "banana" (1 раз)
            expected: []string{"1 2 apple", "5 6 banana"},
        },
        {
            name: "skip more fields than available",
            input: []string{
                "a b",
                "c d",
            },
            opts:     Options{NumFields: 5},
            // После пропуска 5 полей: обе строки становятся пустыми
            // Все строки в одной группе, выводим только первую
            expected: []string{"a b"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Uniq(tt.input, tt.opts)
            require.NoError(t, err)
            assert.Equal(t, tt.expected, result)
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
            // После пропуска 5 символов: "123", "456", "789"
            // Все строки уникальны
            expected: []string{"apple123", "apple456", "banana789"},
        },
        {
            name: "skip more chars than available", 
            input: []string{
                "abc",
                "def",
            },
            opts:     Options{NumChars: 10},
            // После пропуска 10 символов: обе строки становятся пустыми
            // Все строки в одной группе, выводим только первую
            expected: []string{"abc"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Uniq(tt.input, tt.opts)
            require.NoError(t, err)
            assert.Equal(t, tt.expected, result)
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
            // После пропуска 2 полей: "abc123", "def456", "abc123"
            // После пропуска 3 символов: "123", "456", "123"  
            // Группы: "123" (2 раза), "456" (1 раз)
            expected: []string{"field1 field2 abc123", "field1 field2 def456"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Uniq(tt.input, tt.opts)
            require.NoError(t, err)
            assert.Equal(t, tt.expected, result)
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
			assert.ErrorIs(t, err, tt.want)
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
			assert.Equal(t, tt.expected, result)
		})
	}
}