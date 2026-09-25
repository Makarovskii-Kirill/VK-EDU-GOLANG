package calc

import (
	"fmt"
	"strconv"
)

const leftBracket = "("
const rightBracket = ")"
const plus = "+"
const minus = "-"
const multi = "*"
const div = "/"
const space = " "

func StringToSlice(expression string) ([]string, error) {
	var result []string
	var number string
	bracketsCounter := 0
	for _, c := range expression {
		if (c >= '0' && c <= '9') || c == '.' {
			number += string(c)
			continue
		}
		if number != "" {
			result = append(result, number)
			number = ""
		}

		switch string(c) {
		case leftBracket:
			result = append(result, string(c))
			bracketsCounter++
		case rightBracket:
			result = append(result, string(c))
			bracketsCounter--
			if bracketsCounter < 0 {
				return nil, fmt.Errorf("закрывающая скобка раньше открывающей")
			}
		case plus:
			result = append(result, string(c))
		case minus:
			result = append(result, string(c))
		case multi:
			result = append(result, string(c))
		case div:
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

func multiplier(chars []string, pos int) (float64, int, error) {
	if pos >= len(chars) {
		return 0, pos, fmt.Errorf("за пределами выражения")
	}

	if chars[pos] == minus {
		result, newPos, err := multiplier(chars, pos+1)
		return (-1 * result), newPos, err
	}
	if chars[pos] == leftBracket {
		result, newPos, err := expression(chars, pos+1)
		if err != nil {
			return 0, newPos, err
		}
		if newPos >= len(chars) || chars[newPos] != rightBracket {
			return 0, newPos, fmt.Errorf("где закрывающая скобка")
		}
		return result, newPos + 1, nil
	}
	result, err := strconv.ParseFloat(chars[pos], 64)
	if err != nil {
		return 0, pos, fmt.Errorf("где число")
	}
	return result, pos + 1, nil

}

func term(chars []string, pos int) (float64, int, error) {
	result, pos, err := multiplier(chars, pos)
	if err != nil {
		return 0, pos, err
	}
	for pos < len(chars) && (chars[pos] == "*" || chars[pos] == "/") {
		right, newPos, err := multiplier(chars, pos+1)
		if err != nil {
			return 0, pos, err
		}
		if chars[pos] == "*" {
			result *= right
		}
		if chars[pos] == "/" {
			if right == 0 {
				return 0, pos, fmt.Errorf("деление на ноль")
			}
			result /= right
		}
		pos = newPos
	}
	return result, pos, nil

}

func expression(chars []string, pos int) (float64, int, error) {
	result, pos, err := term(chars, pos)
	if err != nil {
		return 0, pos, err
	}
	for pos < len(chars) && (chars[pos] == "+" || chars[pos] == "-") {
		right, newPos, err := term(chars, pos+1)
		if err != nil {
			return 0, pos, err
		}
		if chars[pos] == "+" {
			result += right
		}
		if chars[pos] == "-" {
			result -= right
		}
		pos = newPos
	}
	return result, pos, nil
}

func Calc(expr string) (float64, error) {
	chars, err := StringToSlice(expr)
	if err != nil {
		return 0, err
	}

	result, pos, err := expression(chars, 0)
	if err != nil {
		return 0, err
	}

	if pos != len(chars) {
		return 0, fmt.Errorf("лишние символы")
	}

	return result, nil
}
