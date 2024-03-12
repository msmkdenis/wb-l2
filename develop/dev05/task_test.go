package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomGrep(t *testing.T) {
	testCases := []struct {
		name    string
		cfg     config
		pattern string
		input   []string
		result  string
	}{
		{name: "grep lineNum",
			cfg: config{
				after:      0,
				before:     0,
				context:    0,
				count:      false,
				ignoreCase: false,
				invert:     false,
				fixed:      false,
				lineNum:    true,
			},
			pattern: "KAPPA",
			input:   []string{"alpha betta gamma", "sun ray god", "KAPPA lambda", "gnome kDe mate", "kappa final end"},
			result:  "3 KAPPA lambda\n"},
		{name: "grep lineNum ignoreCase",
			cfg: config{
				after:      0,
				before:     0,
				context:    0,
				count:      false,
				ignoreCase: true,
				invert:     false,
				fixed:      false,
				lineNum:    true,
			},
			pattern: "KAPPA",
			input:   []string{"alpha betta gamma", "sun ray god", "KAPPA lambda", "gnome kDe mate", "kappa final end"},
			result:  "3 KAPPA lambda\n5 kappa final end"},
		{name: "grep lineNum ignoreCase",
			cfg: config{
				after:      1,
				before:     1,
				context:    0,
				count:      false,
				ignoreCase: true,
				invert:     false,
				fixed:      false,
				lineNum:    true,
			},
			pattern: "KAPPA",
			input:   []string{"alpha betta gamma", "sun ray god", "KAPPA lambda", "gnome kDe mate", "kappa final end"},
			result:  "2 sun ray god\n3 KAPPA lambda\n4 gnome kDe mate\n5 kappa final end"},
		{name: "grep lineNum ignoreCase",
			cfg: config{
				after:      1,
				before:     1,
				context:    1,
				count:      false,
				ignoreCase: true,
				invert:     false,
				fixed:      false,
				lineNum:    true,
			},
			pattern: "KAPPA",
			input:   []string{"alpha betta gamma", "sun ray god", "KAPPA lambda", "gnome kDe mate", "kappa final end"},
			result:  "2 sun ray god\n3 KAPPA lambda\n4 gnome kDe mate\n5 kappa final end"},
		{name: "grep lineNum ignoreCase",
			cfg: config{
				after:      0,
				before:     0,
				context:    0,
				count:      false,
				ignoreCase: true,
				invert:     true,
				fixed:      false,
				lineNum:    true,
			},
			pattern: "KAPPA",
			input:   []string{"alpha betta gamma", "sun ray god", "KAPPA lambda", "gnome kDe mate", "kappa final end"},
			result:  "1 alpha betta gamma\n2 sun ray god\n4 gnome kDe mate\n"},
		{name: "grep lineNum ignoreCase",
			cfg: config{
				after:      0,
				before:     0,
				context:    0,
				count:      false,
				ignoreCase: true,
				invert:     false,
				fixed:      true,
				lineNum:    true,
			},
			pattern: "sun ray god",
			input:   []string{"alpha betta gamma", "sun ray god", "KAPPA lambda", "gnome kDe mate", "kappa final end"},
			result:  "2 sun ray god\n"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			res := customGrep(test.cfg, test.pattern, test.input)
			assert.Equal(t, test.result, res)
			clear(dict)
		})
	}

}
