package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
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
	numFields := op.Fields
	numChars := op.Chars
	var result string
	nextIsSpace := false
	for _, c := range line {
		if c != ' ' {
			nextIsSpace = true
		}
		if numFields > 0 && c == ' ' && nextIsSpace {
			numFields--
			nextIsSpace = false
			continue
		}
		if numChars > 0 && numFields == 0 {
			numChars--
			continue
		}
		if numFields == 0 && numChars == 0 {
			result += string(c)
		}
	}

	if op.IgnoreCase {
		result = strings.ToLower(result)
	}

	return result
}

func main() {
	c := flag.Bool("c", false, "to count")
	d := flag.Bool("d", false, "only dublicate")
	u := flag.Bool("u", false, "only unique")
	f := flag.Int("f", 0, "skip first N fields")
	s := flag.Int("s", 0, "skip first N chars")
	i := flag.Bool("i", false, "without case")
	flag.Parse()

	op := Options{
		Count:      *c,
		Duplicates: *d,
		Unique:     *u,
		IgnoreCase: *i,
		Fields:     *f,
		Chars:      *s,
	}

	var in io.Reader = os.Stdin
	var out io.Writer = os.Stdout

	args := flag.Args()

	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "ошибка открытия файла: ", err)
			os.Exit(1)
		}
		defer file.Close()
		in = file
		if len(args) > 1 {
			file, err := os.Create(args[1])
			if err != nil {
				fmt.Fprintln(os.Stderr, "ошибка создания файла: ", err)
				os.Exit(1)
			}
			defer file.Close()
			out = file
		}
	}

	var lines []string
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, "ошибка чтения файла: ", scanner.Err())
		os.Exit(1)
	}

	result, err := Uniq(lines, op)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, line := range result {
		fmt.Fprintln(out, line)
	}
}
