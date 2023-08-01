package reader

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizePrompt(t *testing.T) {
	tests := []struct {
		input  string
		result []string
	}{
		{
			"abc",
			[]string{"abc"},
		},
		{
			"abc xyz",
			[]string{"abc", "xyz"},
		},
		{
			"abc, xyz",
			[]string{"abc", ",", "xyz"},
		},
		{
			"(abc)",
			[]string{"(", "abc", ")"},
		},
		{
			"((abc))",
			[]string{"(", "(", "abc", ")", ")"},
		},
		{
			"(abc:0.5)",
			[]string{"(", "abc", ":", "0.5", ")"},
		},
		{
			"[[abc,xyz]]",
			[]string{"[", "[", "abc", ",", "xyz", "]", "]"},
		},
		{
			"<lora:file name:1.5>",
			[]string{"<", "lora", ":", "file name", ":", "1.5", ">"},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := tokenizeInput(test.input)
			assert.True(t, reflect.DeepEqual(result, test.result))
		})
	}
}
