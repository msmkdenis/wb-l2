package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var MONTHS = map[string]int{
	"jan":  1,
	"feb":  2,
	"mar":  3,
	"apr":  4,
	"may":  5,
	"june": 6,
	"july": 7,
	"aug":  8,
	"sep":  9,
	"oct":  10,
	"nov":  11,
	"dec":  12,
}

type config struct {
	column      int
	number      bool
	reverse     bool
	unique      bool
	month       bool
	checkSorted bool
	file        string
}

func (c *config) parseFlags() {
	var column int
	flag.IntVar(&column, "k", 1, "column to sort")

	var number bool
	flag.BoolVar(&number, "n", false, "number sort")

	var reverse bool
	flag.BoolVar(&reverse, "r", false, "bool reverse sort (default false)")

	var unique bool
	flag.BoolVar(&unique, "u", false, "bool show only unique (default false)")

	var month bool
	flag.BoolVar(&month, "M", false, "bool sort by month (default false)")

	var checkSorted bool
	flag.BoolVar(&checkSorted, "c", false, "bool check if text is sorted (default false)")

	var file string
	flag.StringVar(&file, "f", "", "file")

	flag.Parse()

	// проверим взаимоисключающие флаги
	if number && month {
		fmt.Println("Error: -n and -M flags cannot be used together")
		os.Exit(1)
	}

	if column <= 0 {
		fmt.Println("Error: invalid flag -k must be > 0")
		os.Exit(1)
	}

	c.column = column - 1
	c.number = number
	c.reverse = reverse
	c.unique = unique
	c.month = month
	c.checkSorted = checkSorted
	c.file = file
}

// go run task.go -k=6 -u=true -M=true -r=true -f=text_1.txt

/*
		  Сортировка поддерживает флаг -f с именем файла.

		  Флаги -n и -M не могут быть вызваны одновременно.
		  Флаг -k должен быть > 0.

		  Сортировка не учитывает пустые "ячейки" в строках - всегда располагает их в конце.

		  Флаги сортировки работают следующим образом:
		    -n     сортирует по числовому значению:
		           сверху - вниз:
				    сначала цифры по увеличению
				    слова в лексикографическом порядке
			        пустые ячейки в конце

		    -n -r  сортирует по числовому значению реверсивно:
		           сверху - вниз:
				    сначала цифры по уменьшению
			        пустые ячейки в конце

		    -M     сортирует по месяцам:
		           сверху - вниз:
				    сначала месяцы по увеличению
				    цифры по увеличению
				    слова в лексикографическом порядке
			        пустые ячейки в конце

		    -M -r  сортирует по месяцам:
		           сверху - вниз:
				     сначала месяцы по уменьшению
				     цифры по уменьшению
				     слова в обратном лексикографическом порядке
			         пустые ячейки в конце

		    -u     сортирует по заданной схеме
	               выводит только уникальные строки

			-c     проверяет, отсортирован ли текст
				   выводит только текст: sorted / not sorted
*/
func main() {
	cfg := config{}
	cfg.parseFlags()
	input := getInputData(cfg.file)
	result := sortInput(input, cfg)
	fmt.Println(result)
}

func sortInput(input []string, cfg config) string {
	in := &strings.Builder{}
	data := make([][]string, 0)
	var maxLen int
	for i, v := range input {
		line := strings.Fields(v)
		if len(line) > maxLen {
			maxLen = len(line)
		}
		if i == len(input)-1 {
			in.WriteString(strings.Join(line, " "))
		} else {
			in.WriteString(strings.Join(line, " ") + "\n")
		}
		data = append(data, line)
	}

	res := sortData(data, cfg)
	out := &strings.Builder{}
	if cfg.checkSorted {
		resCompare := make([]string, 0)
		for _, v := range res {
			resCompare = append(resCompare, v...)
		}
		if strings.Join(input, " ") == strings.Join(resCompare, " ") {
			out.WriteString("sorted")
		} else {
			out.WriteString("not sorted")
		}
	} else {
		for i, v := range res {
			if i == len(res)-1 {
				out.WriteString(strings.Join(v, " "))
			} else {
				out.WriteString(strings.Join(v, " ") + "\n")
			}
		}
	}
	return out.String()
}

func sortByMonth(i, j, k int, result [][]string) bool {
	checkLengthA := k >= len(result[i])
	checkLengthB := k >= len(result[j])
	if checkLengthA && checkLengthB {
		return false
	} else {
		if checkLengthA {
			return false
		} else if checkLengthB {
			return true
		}
	}

	a := strings.ToLower(result[i][k])
	b := strings.ToLower(result[j][k])
	aMonthIndex, okA := MONTHS[a]
	bMonthIndex, okB := MONTHS[b]
	if okA && okB {
		return aMonthIndex < bMonthIndex
	} else {
		if okA {
			return true
		} else if okB {
			return false
		} else {
			return result[i][k] < result[j][k]
		}
	}
}

func sortByMonthReversed(i, j, k int, result [][]string) bool {
	checkLengthA := k >= len(result[i])
	checkLengthB := k >= len(result[j])
	if checkLengthA && checkLengthB {
		return false
	} else {
		if checkLengthA {
			return false
		} else if checkLengthB {
			return true
		}
	}

	a := strings.ToLower(result[i][k])
	b := strings.ToLower(result[j][k])
	aMonthIndex, okA := MONTHS[a]
	bMonthIndex, okB := MONTHS[b]
	if okA && okB {
		return aMonthIndex > bMonthIndex
	} else {
		if okA {
			return true
		} else if okB {
			return false
		} else {
			return result[i][k] > result[j][k]
		}
	}
}

func sortByNumberReversed(i, j, k int, result [][]string) bool {
	checkLengthA := k >= len(result[i])
	checkLengthB := k >= len(result[j])
	if checkLengthA && checkLengthB {
		return false
	} else {
		if checkLengthA {
			return false
		} else if checkLengthB {
			return true
		}
	}
	a, errA := strconv.ParseFloat(result[i][k], 64)
	b, errB := strconv.ParseFloat(result[j][k], 64)
	if errA == nil && errB == nil {
		return a > b
	} else {
		if errA != nil {
			return false
		} else if errB != nil {
			return true
		}
		return result[i][k] > result[j][k]

	}
}

func sortByNumber(i, j, k int, result [][]string) bool {
	checkLengthA := k >= len(result[i])
	checkLengthB := k >= len(result[j])
	if checkLengthA && checkLengthB {
		return false
	} else {
		if checkLengthA {
			return false
		} else if checkLengthB {
			return true
		}
	}
	a, errA := strconv.ParseFloat(result[i][k], 64)
	b, errB := strconv.ParseFloat(result[j][k], 64)
	if errA == nil && errB == nil {
		return a < b
	} else {
		return result[i][k] < result[j][k]
	}
}

func defaultSortReversed(i, j, k int, result [][]string) bool {
	checkLengthA := k >= len(result[i])
	checkLengthB := k >= len(result[j])
	if checkLengthA && checkLengthB {
		return false
	} else {
		if checkLengthA {
			return false
		} else if checkLengthB {
			return true
		}
	}
	return result[i][k] > result[j][k]
}

func defaultSort(i, j, k int, result [][]string) bool {
	checkLengthA := k >= len(result[i])
	checkLengthB := k >= len(result[j])
	if checkLengthA && checkLengthB {
		return false
	} else {
		if checkLengthA {
			return false
		} else if checkLengthB {
			return true
		}
	}
	return result[i][k] < result[j][k]
}

func sortData(data [][]string, cfg config) [][]string {
	result := slices.Clone(data)

	if !cfg.reverse {
		switch {
		case cfg.month:
			sort.Slice(result, func(i, j int) bool {
				return sortByMonth(i, j, cfg.column, result)
			})
		case cfg.number:
			sort.Slice(result, func(i, j int) bool {
				return sortByNumber(i, j, cfg.column, result)
			})
		default:
			sort.Slice(result, func(i, j int) bool {
				return defaultSort(i, j, cfg.column, result)
			})
		}
	} else {
		switch {
		case cfg.month:
			sort.Slice(result, func(i, j int) bool {
				return sortByMonthReversed(i, j, cfg.column, result)
			})
		case cfg.number:
			sort.Slice(result, func(i, j int) bool {
				return sortByNumberReversed(i, j, cfg.column, result)
			})
		default:
			sort.Slice(result, func(i, j int) bool {
				return defaultSortReversed(i, j, cfg.column, result)
			})
		}
	}

	if cfg.unique {
		index := len(result)
		unique := make(map[string]struct{})
		for i := 0; i < index; i++ {
			s := strings.Join(result[i], " ")
			_, ok := unique[s]
			if !ok {
				result = append(result, result[i])
			}
			unique[s] = struct{}{}
		}
		result = slices.Delete(result, 0, index)
	}
	return result
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
