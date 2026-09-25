package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccess(t *testing.T) {
	music := []string{
		"I love music.",
		"I love music.",
		"I love music.",
		"",
		"I love music of Kartik.",
		"I love music of Kartik.",
		"Thanks.",
		"I love music of Kartik.",
		"I love music of Kartik.",
	}

	tests := []struct {
		name  string
		input []string
		op    Options
		want  []string
	}{
		{
			name:  "без флагов",
			input: music,
			op:    Options{},
			want: []string{
				"I love music.",
				"",
				"I love music of Kartik.",
				"Thanks.",
				"I love music of Kartik.",
			},
		},
		{
			name:  "-c",
			input: music,
			op:    Options{Count: true},
			want: []string{
				"3 I love music.",
				"1 ",
				"2 I love music of Kartik.",
				"1 Thanks.",
				"2 I love music of Kartik.",
			},
		},
		{
			name:  "-d",
			input: music,
			op:    Options{Duplicates: true},
			want: []string{
				"I love music.",
				"I love music of Kartik.",
				"I love music of Kartik.",
			},
		},
		{
			name:  "-u",
			input: music,
			op:    Options{Unique: true},
			want: []string{
				"",
				"Thanks.",
			},
		},
		{
			name: "-i",
			input: []string{
				"I LOVE MUSIC.",
				"I love music.",
				"I LoVe MuSiC.",
				"",
				"I love MuSIC of Kartik.",
				"I love music of kartik.",
				"Thanks.",
				"I love music of kartik.",
				"I love MuSIC of Kartik.",
			},
			op: Options{IgnoreCase: true},
			want: []string{
				"I LOVE MUSIC.",
				"",
				"I love MuSIC of Kartik.",
				"Thanks.",
				"I love music of kartik.",
			},
		},
		{
			name: "-f 1",
			input: []string{
				"We love music.",
				"I love music.",
				"They love music.",
				"",
				"I love music of Kartik.",
				"We love music of Kartik.",
				"Thanks.",
			},
			op: Options{Fields: 1},
			want: []string{
				"We love music.",
				"",
				"I love music of Kartik.",
				"Thanks.",
			},
		},
		{
			name:  "пустой вход",
			input: nil,
			op:    Options{},
			want:  []string{},
		},
		{
			name:  "одна строка с -c",
			input: []string{"a"},
			op:    Options{Count: true},
			want:  []string{"1 a"},
		},
		{
			name:  "пустая первая строка не теряется",
			input: []string{"", "a", "a"},
			op:    Options{},
			want:  []string{"", "a"},
		},
		{
			name:  "несмежные повторы не схлопываются",
			input: []string{"a", "b", "a"},
			op:    Options{Count: true},
			want:  []string{"1 a", "1 b", "1 a"},
		},
		{
			name:  "-f с несколькими пробелами между полями",
			input: []string{"aa   bb ccdef", "xx yy ccdef"},
			op:    Options{Fields: 2},
			want:  []string{"aa   bb ccdef"},
		},
		{
			name:  "-f больше числа полей",
			input: []string{"a b", "c d"},
			op:    Options{Fields: 5},
			want:  []string{"a b"},
		},
		{
			name:  "-s",
			input: []string{"xxabc", "yyabc", "zzabd"},
			op:    Options{Chars: 2},
			want:  []string{"xxabc", "zzabd"},
		},
		{
			name:  "-s больше длины строки",
			input: []string{"abc", "de"},
			op:    Options{Chars: 10},
			want:  []string{"abc"},
		},
		{
			name:  "-f и -s вместе",
			input: []string{"aa bb ccdef", "xx yy zzdef", "xx yy zzdeg"},
			op:    Options{Fields: 2, Chars: 2},
			want:  []string{"aa bb ccdef", "xx yy zzdeg"},
		},
		{
			name:  "-c, -i и -f вместе",
			input: []string{"1 Foo", "2 foo", "3 bar"},
			op:    Options{Count: true, IgnoreCase: true, Fields: 1},
			want:  []string{"2 1 Foo", "1 3 bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Uniq(tt.input, tt.op)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name string
		op   Options
	}{
		{name: "-c и -d", op: Options{Count: true, Duplicates: true}},
		{name: "-c и -u", op: Options{Count: true, Unique: true}},
		{name: "-d и -u", op: Options{Duplicates: true, Unique: true}},
		{name: "-c, -d и -u", op: Options{Count: true, Duplicates: true, Unique: true}},
		{name: "отрицательный -f", op: Options{Fields: -1}},
		{name: "отрицательный -s", op: Options{Chars: -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Uniq([]string{"a", "a"}, tt.op)
			require.Error(t, err)
			require.Nil(t, got)
		})
	}
}
