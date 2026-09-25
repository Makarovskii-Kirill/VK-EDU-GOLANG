package main

import (
	"fmt"
	"go-calc/calc"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	var expr string

	if len(os.Args) > 1 {
		expr = strings.Join(os.Args[1:], " ")
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ошибка чтения:", err)
			os.Exit(1)
		}
		expr = strings.TrimSpace(string(data))
	}

	result, err := calc.Calc(expr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}

	fmt.Println(strconv.FormatFloat(result, 'f', -1, 64))
}
