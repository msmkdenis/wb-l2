package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortText(t *testing.T) {
	testCases := []struct {
		name   string
		cfg    config
		input  []string
		result string
	}{
		{name: "default sort by first column", cfg: config{
			column:      0,
			number:      false,
			reverse:     false,
			unique:      false,
			month:       false,
			checkSorted: false,
		}, input: []string{"b g", "a b", "k f"}, result: "a b\nb g\nk f"},

		{name: "default sort by first column reversed", cfg: config{
			column:      0,
			number:      false,
			reverse:     true,
			unique:      false,
			month:       false,
			checkSorted: false,
		}, input: []string{"b g", "a b", "k f"}, result: "k f\nb g\na b"},

		{name: "default sort by first column reversed", cfg: config{
			column:      0,
			number:      false,
			reverse:     true,
			unique:      false,
			month:       false,
			checkSorted: false,
		}, input: []string{"b g", "a b", "k f"}, result: "k f\nb g\na b"},

		{name: "month sort by first column", cfg: config{
			column:      0,
			number:      false,
			reverse:     false,
			unique:      false,
			month:       true,
			checkSorted: false,
		}, input: []string{"aug b g", "JULY a b", " jan k f"}, result: "jan k f\nJULY a b\naug b g"},

		{name: "month sort by first column reversed", cfg: config{
			column:      0,
			number:      false,
			reverse:     true,
			unique:      false,
			month:       true,
			checkSorted: false,
		}, input: []string{"aug b g", "JULY a b", " jan k f"}, result: "aug b g\nJULY a b\njan k f"},

		{name: "number sort by first column", cfg: config{
			column:      0,
			number:      true,
			reverse:     false,
			unique:      false,
			month:       false,
			checkSorted: false,
		}, input: []string{"55.57 b g", "22.45 a b", "33 k f", "44.44 c d", "22 e f"}, result: "22 e f\n22.45 a b\n33 k f\n44.44 c d\n55.57 b g"},

		{name: "number sort by first column reversed", cfg: config{
			column:      0,
			number:      true,
			reverse:     true,
			unique:      false,
			month:       false,
			checkSorted: false,
		}, input: []string{"55.57 b g", "22.45 a b", "33 k f", "44.44 c d", "22 e f"}, result: "55.57 b g\n44.44 c d\n33 k f\n22.45 a b\n22 e f"},

		{name: "check of sorted with given config", cfg: config{
			column:      0,
			number:      true,
			reverse:     true,
			unique:      false,
			month:       false,
			checkSorted: true,
		}, input: []string{"55.57 b g", "44.44 c d", "33 k f", "22.45 a b", "22 e f"}, result: "sorted"},

		{name: "check of sorted with given config (not sorted)", cfg: config{
			column:      0,
			number:      true,
			reverse:     true,
			unique:      false,
			month:       false,
			checkSorted: true,
		}, input: []string{"55.57 b g", "22.45 a b", "33 k f", "44.44 c d", "22 e f"}, result: "not sorted"},

		{name: "number sort by first column reversed only unique", cfg: config{
			column:      0,
			number:      true,
			reverse:     true,
			unique:      true,
			month:       false,
			checkSorted: false,
		}, input: []string{"55.57 b g", "22.45 a b", "33 k f", "44.44 c d", "22 e f", "44.44 c d", "44.44 c d"}, result: "55.57 b g\n44.44 c d\n33 k f\n22.45 a b\n22 e f"},

		{name: "number sort by third column with white empty cells and chars", cfg: config{
			column:      2,
			number:      true,
			reverse:     false,
			unique:      false,
			month:       false,
			checkSorted: false,
		}, input: []string{"b g 55.57", "a b 22.45", "k f 33", "c d 44.44", "e f 22", "c d a", "c d c", "c d", "k m"}, result: "e f 22\na b 22.45\nk f 33\nc d 44.44\nb g 55.57\nc d a\nc d c\nc d\nk m"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			res := sortInput(test.input, test.cfg)
			assert.Equal(t, test.result, res)
		})
	}
}
