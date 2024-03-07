package main

import (
	"fmt"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	Separator = `\`
)

func main() {
	fmt.Println(UnpackEscape(`a4bc2d5e`), UnpackEscape(`a4bc2d5e`) == `aaaabccddddde`)
	fmt.Println(UnpackEscape(`abcd`), UnpackEscape(`abcd`) == `abcd`)
	fmt.Println(UnpackEscape(`45`), UnpackEscape(`45`) == ``)
	fmt.Println(UnpackEscape(``), UnpackEscape(``) == ``)
	fmt.Println(UnpackEscape(`qwe\4\5`), UnpackEscape(`qwe\4\5`) == `qwe45`)
	fmt.Println(UnpackEscape(`qwe\45`), UnpackEscape(`qwe\45`) == `qwe44444`)
	fmt.Println(UnpackEscape(`qwe\\5`), UnpackEscape(`qwe\\5`) == `qwe\\\\\`)
	fmt.Println(UnpackEscape(`qwe\\5`), UnpackEscape(`qwe\\5`) == `qwe\\\\\`)
}

func UnpackEscape(str string) string {
	if str == "" {
		return ""
	}

	// Переводим в руны
	s := []rune(str)
	out := &strings.Builder{}

	// Строка не может начинаться с цифры
	if unicode.IsDigit(s[0]) {
		return ""
	}

	// Флаг экранирования
	var escape bool
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			if escape == false {
				escape = true
				continue
			} else {
				out.WriteRune(s[i])
				escape = false
				continue
			}
		}

		if s[i] != '\\' && !unicode.IsDigit(s[i]) {
			out.WriteRune(s[i])
		}

		if unicode.IsDigit(s[i]) {
			if escape == false {
				for k := 0; k < int(s[i]-'1'); k++ {
					out.WriteRune(s[i-1])
				}
			}
			if escape == true {
				out.WriteRune(s[i])
				escape = false
			}
		}
	}

	return out.String()
}
