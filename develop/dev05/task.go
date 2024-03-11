package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
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

	fmt.Println(cfg)
	customGrep(cfg, str[0], input)
}

func customGrep(cfg config, str string, input []string) {

	var pattern *regexp.Regexp
	if cfg.ignoreCase {
		pattern = regexp.MustCompile(strings.ToLower(str))
	} else {
		pattern = regexp.MustCompile(str)
	}

	answer := &strings.Builder{}
	for i, v := range input {
		if cfg.ignoreCase {
			if cfg.invert {
				if !cfg.fixed && !pattern.Match([]byte(strings.ToLower(v))) {
					writeAnswer(answer, cfg, i, v, input)
				}
				if cfg.fixed && strings.ToLower(str) != strings.ToLower(v) {
					writeAnswer(answer, cfg, i, v, input)
				}
			} else {
				if !cfg.fixed && pattern.Match([]byte(strings.ToLower(v))) {
					writeAnswer(answer, cfg, i, v, input)
				}
				if cfg.fixed && strings.ToLower(str) == strings.ToLower(v) {
					writeAnswer(answer, cfg, i, v, input)
				}
			}
		} else {
			if cfg.invert {
				if !cfg.fixed && !pattern.Match([]byte(v)) {
					writeAnswer(answer, cfg, i, v, input)
				}
				if cfg.fixed && str != v {
					writeAnswer(answer, cfg, i, v, input)
				}
			} else {
				if !cfg.fixed && pattern.Match([]byte(v)) {
					writeAnswer(answer, cfg, i, v, input)
				}
				if cfg.fixed && str != v {
					writeAnswer(answer, cfg, i, v, input)
				}
			}
		}
	}
	fmt.Println(answer.String())
}

var dict = make(map[string]struct{})

func writeAnswer(answer *strings.Builder, cfg config, i int, v string, input []string) {

	if cfg.before >= i+1 {
		for j := 0; j < i; j++ {
			s := fmt.Sprintf(getFormat(j, len(input)), getLineNum(j, cfg.lineNum), input[j])
			if _, ok := dict[s]; !ok {
				dict[s] = struct{}{}
				answer.WriteString(s)
			}
		}
	} else {
		for j := i - cfg.before; j < i; j++ {
			s := fmt.Sprintf(getFormat(j, len(input)), getLineNum(j, cfg.lineNum), input[j])
			if _, ok := dict[s]; !ok {
				dict[s] = struct{}{}
				answer.WriteString(s)
			}
		}
	}

	s := fmt.Sprintf(getFormat(i, len(input)), getLineNum(i, cfg.lineNum), v)
	if _, ok := dict[s]; !ok {
		dict[s] = struct{}{}
		answer.WriteString(s)
	}

	//answer.WriteString(fmt.Sprintf(getFormat(i, len(input)), getLineNum(i, cfg.lineNum), v))

	if cfg.after >= len(input)-i-1 {
		for j := i + 1; j < len(input); j++ {
			s := fmt.Sprintf(getFormat(j, len(input)), getLineNum(j, cfg.lineNum), input[j])
			if _, ok := dict[s]; !ok {
				dict[s] = struct{}{}
				answer.WriteString(s)
			}
		}
	} else {
		for j := i + 1; j < i+1+cfg.after; j++ {
			s := fmt.Sprintf(getFormat(j, len(input)), getLineNum(j, cfg.lineNum), input[j])
			if _, ok := dict[s]; !ok {
				dict[s] = struct{}{}
				answer.WriteString(s)
			}
		}
	}
}

func getLineNum(i int, isLineNum bool) string {
	if isLineNum {
		return strconv.Itoa(i+1) + " "
	}
	return ""
}

func getFormat(i, len int) string {
	if i == len-1 {
		return "%s%s"
	}
	return "%s%s\n"
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
