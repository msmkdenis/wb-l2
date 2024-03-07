package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type config struct {
	fields    string
	delimiter string
	separated bool
}

func (c *config) parseFlags() {
	var fields string
	flag.StringVar(&fields, "f", "", "choose field")

	var delimiter string
	flag.StringVar(&delimiter, "d", "\t", "choose delimiter")

	var separated bool
	flag.BoolVar(&separated, "s", false, "only with delimiter")

	flag.Parse()

	c.fields = fields
	c.delimiter = delimiter
	c.separated = separated
}

func main() {
	cfg := config{}
	cfg.parseFlags()
	indexes := parseFields(cfg.fields)
	input := readText(cfg)

	if cfg.fields != "" {
		for _, line := range input {
			if len(indexes) == 1 {
				var s []string
				for _, v := range line {
					s = append(s, v)
				}
				fmt.Println(strings.Join(s, " "))
			} else {
				if indexes[0] == math.MinInt {
					var s []string
					for k, v := range line {
						if k <= indexes[1] {
							s = append(s, v)
						}
					}
					fmt.Println(strings.Join(s, " "))
				} else if indexes[1] == math.MaxInt {
					var s []string
					for k, v := range line {
						if k >= indexes[0] {
							s = append(s, v)
						}
					}
					fmt.Println(strings.Join(s, " "))
				} else {
					var s []string
					for _, v := range indexes {
						if v < len(line) {
							s = append(s, line[v])
						}
					}
					fmt.Println(strings.Join(s, " "))
				}
			}
		}
	}

}

func parseFields(fields string) []int {
	var indexes []int
	if strings.Contains(fields, "-") {
		s := strings.Split(fields, "-")
		if len(s) == 2 {
			if s[0] == "" {
				fmt.Println(s[1])
				a, err := strconv.Atoi(s[1])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				indexes = append(indexes, math.MinInt)
				indexes = append(indexes, a)
			} else if s[1] == "" {
				fmt.Println(s[0])
				a, err := strconv.Atoi(s[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				indexes = append(indexes, a)
				indexes = append(indexes, math.MaxInt)
			} else {
				fmt.Println(s[0], s[1])
				a, err := strconv.Atoi(s[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				b, err := strconv.Atoi(s[1])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if a > b {
					fmt.Println("invalid fields input")
					os.Exit(1)
				}
				for i := a; i < b; i++ {
					indexes = append(indexes, i)
				}
			}
		} else {
			fmt.Println("invalid fields input")
			os.Exit(1)
		}
	} else if strings.Contains(fields, ",") {
		s := strings.Split(fields, ",")
		for _, v := range s {
			a, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			indexes = append(indexes, a)
		}
	} else {
		a, err := strconv.Atoi(fields)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		indexes = append(indexes, a)
	}
	return indexes
}

func readText(cfg config) [][]string {
	scanner := bufio.NewScanner(os.Stdin)
	text := make([][]string, 0)
	fmt.Println("Enter text, end with empty line")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		if cfg.separated {
			if strings.Contains(line, cfg.delimiter) {
				text = append(text, strings.Split(line, cfg.delimiter))
			}
		} else {
			text = append(text, strings.Split(line, cfg.delimiter))
		}
	}
	return text
}
