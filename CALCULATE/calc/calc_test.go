package calc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalcSuccess(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want float64
	}{
		{name: "пример", expr: "(1+2)-3", want: 0},
		{name: "пример", expr: "(1+2)*3", want: 9},
		{name: "пример", expr: "7/2", want: 3.5},
		{name: "пример", expr: "1+2*3", want: 7},
		{name: "пример", expr: "8-3-2", want: 3},
		{name: "пример", expr: "-(1+2)", want: -3},
		{name: "пример", expr: "2*-1.5", want: -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calc(tt.expr)
			require.NoError(t, err)
			require.InDelta(t, tt.want, got, 1e-9)
		})
	}
}

func TestCalcError(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{name: "деление на ноль", expr: "(1/0)"},
		{name: "неизвестный символ", expr: "2x"},
		{name: "незакрытая скобка", expr: "(1+2"},
		{name: "закрывающая скобка раньше открывающей", expr: ")1("},
		{name: "оператор в конце", expr: "1+2*"},
		{name: "два числа без оператора", expr: "1 2"},
		{name: "два числа внутри скобок", expr: "(1 2)"},
		{name: "неправильное число", expr: "1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Calc(tt.expr)
			require.Error(t, err)
		})
	}
}
