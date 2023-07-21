package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizePrompt(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result []Token
	}{
		{"abc", "abc", []Token{{"abc", 0}}},
		{"abc xyz", "abc xyz", []Token{{"abc", 0}, {"xyz", 4}}},
		{"abc, xyz", "abc, xyz", []Token{{"abc", 0}, {",", 3}, {"xyz", 5}}},
		{"(abc)", "(abc)", []Token{{"(", 0}, {"abc", 1}, {")", 4}}},
		{"[abc:0.5]", "[abc:0.5]", []Token{{"[", 0}, {"abc", 1}, {":", 4}, {"0.5", 5}, {"]", 8}}},
		{"(abc|xyz)", "(abc|xyz)", []Token{{"(", 0}, {"abc", 1}, {"|", 4}, {"xyz", 5}, {")", 8}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TokenizePrompt(test.input)
			assert.True(t, reflect.DeepEqual(result, test.result))
		})
	}
}
