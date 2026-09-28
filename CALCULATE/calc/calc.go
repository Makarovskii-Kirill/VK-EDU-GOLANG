package calc

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	leftBracket  = '('
	rightBracket = ')'
	plus         = '+'
	minus        = '-'
	multi        = '*'
	div          = '/'
	space        = ' '
)

func tokenize(expr string) ([]string, error) {
	var result []string
	var number string
	bracketsCounter := 0
	for _, c := range expr {
		if (c >= '0' && c <= '9') || c == '.' {
			number += string(c)
			continue
		}
		if number != "" {
			result = append(result, number)
			number = ""
		}

		switch c {
		case leftBracket:
			result = append(result, string(c))
			bracketsCounter++
		case rightBracket:
			result = append(result, string(c))
			bracketsCounter--
			if bracketsCounter < 0 {
				return nil, fmt.Errorf("закрывающая скобка раньше открывающей")
			}
		case plus, minus, multi, div:
			result = append(result, string(c))
		case space:
			continue
		default:
			return nil, fmt.Errorf("неизвестный символ")
		}
	}
	if number != "" {
		result = append(result, number)
	}
	if bracketsCounter != 0 {
		return nil, fmt.Errorf("неверное число скобок")
	}
	return result, nil
}

func parseFactor(tokens []string, pos int) (float64, int, error) {
	if pos >= len(tokens) {
		return 0, pos, fmt.Errorf("за пределами выражения")
	}

	if tokens[pos] == string(minus) {
		result, newPos, err := parseFactor(tokens, pos+1)
		return (-1 * result), newPos, err
	}
	if tokens[pos] == string(leftBracket) {
		result, newPos, err := parseExpression(tokens, pos+1)
		if err != nil {
			return 0, newPos, err
		}
		if newPos >= len(tokens) || tokens[newPos] != string(rightBracket) {
			return 0, newPos, fmt.Errorf("где закрывающая скобка")
		}
		return result, newPos + 1, nil
	}
	result, err := strconv.ParseFloat(tokens[pos], 64)
	if err != nil {
		return 0, pos, fmt.Errorf("где число")
	}
	return result, pos + 1, nil
}

func parseTerm(tokens []string, pos int) (float64, int, error) {
	result, pos, err := parseFactor(tokens, pos)
	if err != nil {
		return 0, pos, err
	}
	for pos < len(tokens) && (tokens[pos] == string(multi) || tokens[pos] == string(div)) {
		right, newPos, err := parseFactor(tokens, pos+1)
		if err != nil {
			return 0, pos, err
		}
		if tokens[pos] == string(multi) {
			result *= right
		}
		if tokens[pos] == string(div) {
			if right == 0 {
				return 0, pos, fmt.Errorf("деление на ноль")
			}
			result /= right
		}
		pos = newPos
	}
	return result, pos, nil
}

func parseExpression(tokens []string, pos int) (float64, int, error) {
	result, pos, err := parseTerm(tokens, pos)
	if err != nil {
		return 0, pos, err
	}
	for pos < len(tokens) && (tokens[pos] == string(plus) || tokens[pos] == string(minus)) {
		right, newPos, err := parseTerm(tokens, pos+1)
		if err != nil {
			return 0, pos, err
		}
		if tokens[pos] == string(plus) {
			result += right
		}
		if tokens[pos] == string(minus) {
			result -= right
		}
		pos = newPos
	}
	return result, pos, nil
}

func Calc(expr string) (float64, error) {
	if strings.TrimSpace(expr) == "" {
		return 0, nil
	}

	tokens, err := tokenize(expr)
	if err != nil {
		return 0, err
	}

	result, pos, err := parseExpression(tokens, 0)
	if err != nil {
		return 0, err
	}

	if pos != len(tokens) {
		return 0, fmt.Errorf("лишние символы")
	}

	return result, nil
}
