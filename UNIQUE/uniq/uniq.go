package uniq

import (
	"errors"
	"fmt"
	"strings"
)

type flaging int

const (
	flagNo flaging = iota
	flagC
	flagD
	flagU
)

type Options struct {
	Count      bool
	Duplicates bool
	Unique     bool
	IgnoreCase bool
	Fields     int
	Chars      int
}

func appendGroup(result []string, count int, line string, fl flaging) []string {
	switch fl {
	case flagNo:
		result = append(result, line)
	case flagC:
		result = append(result, fmt.Sprintf("%d %s", count, line))
	case flagD:
		if count > 1 {
			result = append(result, line)
		}
	case flagU:
		if count == 1 {
			result = append(result, line)
		}
	}
	return result
}

func Uniq(lines []string, op Options) ([]string, error) {
	flagCount := 0

	if op.Count {
		flagCount++
	}
	if op.Duplicates {
		flagCount++
	}
	if op.Unique {
		flagCount++
	}

	if flagCount > 1 {
		return nil, errors.New("Usage: uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]")
	}

	if op.Fields < 0 || op.Chars < 0 {
		return nil, errors.New("ошибка: -f и -s не могут быть отрицательными")
	}

	fl := flagNo

	switch {
	case op.Count:
		fl = flagC
	case op.Duplicates:
		fl = flagD
	case op.Unique:
		fl = flagU
	}

	result := make([]string, 0, len(lines))
	var currentLine, currentCutLine string
	count := 0
	hasCurrent := false

	for _, line := range lines {
		cutLine := skip(line, op)

		if hasCurrent && cutLine == currentCutLine {
			count++
			continue
		}

		if hasCurrent {
			result = appendGroup(result, count, currentLine, fl)
		}

		currentLine = line
		currentCutLine = cutLine
		count = 1
		hasCurrent = true
	}

	if hasCurrent {
		result = appendGroup(result, count, currentLine, fl)
	}

	return result, nil
}

func skip(line string, op Options) string {
	if op.Fields > 0 {
		fields := strings.Fields(line)
		if op.Fields >= len(fields) {
			line = ""
		} else {
			line = strings.Join(fields[op.Fields:], " ")
		}
	}

	if op.Chars > 0 {
		runes := []rune(line)
		if op.Chars >= len(runes) {
			line = ""
		} else {
			line = string(runes[op.Chars:])
		}
	}

	if op.IgnoreCase {
		line = strings.ToLower(line)
	}

	return line
}
