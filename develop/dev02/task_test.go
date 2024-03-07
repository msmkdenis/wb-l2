package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackEscape(t *testing.T) {
	testCases := []struct {
		name   string
		line   string
		result string
	}{
		{name: `a4bc2d5e`, line: `a4bc2d5e`, result: `aaaabccddddde`},
		{name: `abcd`, line: `abcd`, result: `abcd`},
		{name: `45`, line: `45`, result: ``},
		{name: ``, line: ``, result: ``},
		{name: `qwe\4\5`, line: `qwe\4\5`, result: `qwe45`},
		{name: `qwe\45`, line: `qwe\45`, result: `qwe44444`},
		{name: `qwe\\5`, line: `qwe\\5`, result: `qwe\\\\\`},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			res := UnpackEscape(test.line)
			assert.Equal(t, res, test.result)
		})
	}
}
