package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type config struct {
	after      int
	before     int
	context    int
	count      int
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
}

func (c *config) parseFlags() {
	var after int
	flag.IntVar(&after, "A", 0, "печатать +N строк после совпадения")

	var before int
	flag.IntVar(&before, "B", 0, "печатать +N строк до совпадения")

	var context int
	flag.IntVar(&context, "C", 0, "(A+B) печатать ±N строк вокруг совпадения")

	var count int
	flag.IntVar(&count, "c", 0, "количество строк")

	var ignoreCase bool
	flag.BoolVar(&ignoreCase, "i", false, "игнорировать регистр")

	var invert bool
	flag.BoolVar(&invert, "v", false, "вместо совпадения, исключать")

	var fixed bool
	flag.BoolVar(&fixed, "F", false, "точное совпадение со строкой, не паттерн")

	var lineNum bool
	flag.BoolVar(&lineNum, "n", false, "печатать номер строки")

	flag.Parse()

	c.after = after
	c.before = before
	c.context = c.after + c.before
	c.count = count
	c.ignoreCase = ignoreCase
	c.invert = invert
	c.fixed = fixed
	c.lineNum = lineNum
}

func main() {
	cfg := config{}
	cfg.parseFlags()

	str := flag.Args()
	if len(str) < 2 {
		slog.Error("not enough arguments")
		os.Exit(1)
	}

	input := getInputData(str[1])

	customGrep(cfg, str[0], input)
}

func customGrep(cfg config, str string, input []string) {

	var pattern *regexp.Regexp
	if cfg.ignoreCase {
		pattern = regexp.MustCompile(strings.ToLower(str))
	} else {
		pattern = regexp.MustCompile(str)
	}

	for i, v := range input {
		if pattern.Match([]byte(v)) {
			fmt.Println(i+1, v)
		}
	}
}

func getInputData(file string) []string {
	input, err := os.Open(file)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	var s []string
	for scanner.Scan() {
		bufStr := scanner.Text()
		s = append(s, bufStr)
	}

	return s
}
