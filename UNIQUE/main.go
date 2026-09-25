package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-uniq/uniq"
	"io"
	"os"
)

func main() {
	c := flag.Bool("c", false, "to count")
	d := flag.Bool("d", false, "only dublicate")
	u := flag.Bool("u", false, "only unique")
	f := flag.Int("f", 0, "skip first N fields")
	s := flag.Int("s", 0, "skip first N chars")
	i := flag.Bool("i", false, "without case")
	flag.Parse()

	op := uniq.Options{
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

	result, err := uniq.Uniq(lines, op)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, line := range result {
		fmt.Fprintln(out, line)
	}
}
