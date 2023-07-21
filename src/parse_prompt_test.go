package src

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeToken(t *testing.T) {
	token := "\\abc"
	assert.Equal(t, "abc", escapeToken(token))

	token = "\\\\abc"
	assert.Equal(t, "\\abc", escapeToken(token))
}

func TestParseTagPrompt(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result Prompt
	}{
		{"abc", "abc", Prompt{kind: "tag", name: "abc", tokens: []string{"abc"}}},
		{"abc xyz", "abc xyz", Prompt{kind: "tag", name: "abc xyz", tokens: []string{"abc", "xyz"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := NewTokenReader(TokenizePrompt(test.input))
			result, err := parseTagPrompt(reader)
			assert.Equal(t, nil, err)
			assert.True(t, reflect.DeepEqual(result, test.result))
		})
	}

	// Test error
	input := "(abc)"
	reader := NewTokenReader(TokenizePrompt(input))
	result, err := parseTagPrompt(reader)
	assert.EqualError(t, err, "tag expected")
	assert.True(t, reflect.DeepEqual(Prompt{}, result))
}

func TestParsePositivePrompt(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result Prompt
	}{
		{"(abc)", "(abc)", Prompt{kind: "pw", contents: []Prompt{{kind: "tag", name: "abc", tokens: []string{"abc"}}}}},
		{"((abc))", "((abc))", Prompt{kind: "pw", contents: []Prompt{{kind: "pw", contents: []Prompt{{kind: "tag", name: "abc", tokens: []string{"abc"}}}}}}},
		{"(abc, xyz)", "(abc, xyz)", Prompt{kind: "pw", contents: []Prompt{{kind: "tag", name: "abc", tokens: []string{"abc"}}, {kind: "tag", name: "xyz", tokens: []string{"xyz"}}}}},
		{"(abc:xyz)", "(abc:xyz)", Prompt{kind: "pw", contents: []Prompt{{kind: "tag", name: "abc", tokens: []string{"abc"}}, {kind: "tag", name: "xyz", tokens: []string{"xyz"}}}}},
		{"(abc:1.5)", "(abc:1.5)", Prompt{kind: "ew", weight: 1.5, weightText: "1.5", contents: []Prompt{{kind: "tag", name: "abc", tokens: []string{"abc"}}}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := NewTokenReader(TokenizePrompt(test.input))
			result, err := parsePositivePrompt(reader)
			assert.Equal(t, nil, err)
			assert.True(t, reflect.DeepEqual(result, test.result))
		})
	}
}

func TestParseNegativePrompt(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result Prompt
	}{
		{"[abc]", "[abc]", Prompt{kind: "nw", contents: []Prompt{{kind: "tag", name: "abc", tokens: []string{"abc"}}}}},
		{"[[abc]]", "[[abc]]", Prompt{kind: "nw", contents: []Prompt{{kind: "nw", contents: []Prompt{{kind: "tag", name: "abc", tokens: []string{"abc"}}}}}}},
		{"[abc, xyz]", "[abc, xyz]", Prompt{kind: "nw", contents: []Prompt{{kind: "tag", name: "abc", tokens: []string{"abc"}}, {kind: "tag", name: "xyz", tokens: []string{"xyz"}}}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := NewTokenReader(TokenizePrompt(test.input))
			result, err := parseNegativePrompt(reader)
			assert.Equal(t, nil, err)
			assert.True(t, reflect.DeepEqual(result, test.result))
		})
	}
}

func TestParseContentToken(t *testing.T) {
	input := "abc"
	reader := NewTokenReader(TokenizePrompt(input))
	result, err := parseContentToken(reader, "Name")
	assert.Equal(t, nil, err)
	assert.Equal(t, "abc", result)

	input = "(abc)"
	reader = NewTokenReader(TokenizePrompt(input))
	result, err = parseContentToken(reader, "Name")
	assert.EqualError(t, err, "Name expected")
	assert.Equal(t, "", result)
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		result     float64
		resultText string
	}{
		{"1", "1", 1, "1"},
		{"1.5", "1.5", 1.5, "1.5"},
		{"0,5", "0,5", 0.5, "0.5"},
		{"1. 5", "1. 5", 1.5, "1.5"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := NewTokenReader(TokenizePrompt(test.input))
			result, resultText, err := parseNumber(reader, "Name")
			assert.Equal(t, nil, err)
			assert.Equal(t, result, test.result)
			assert.Equal(t, resultText, test.resultText)
		})
	}
}

func TestParseAnglePrompt(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		inputName string
		result    Prompt
	}{
		{"<lora:file>", "<lora:file>", "lora", Prompt{kind: "lora", filename: "file"}},
		{"<hypernet:file:0.5>", "<hypernet:file:0.5>", "hypernet", Prompt{kind: "hypernet", filename: "file", multiplier: 0.5, multiplierText: "0.5"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := NewTokenReader(TokenizePrompt(test.input))
			result, err := parseAnglePrompt(reader, test.inputName)
			assert.Equal(t, nil, err)
			assert.True(t, reflect.DeepEqual(result, test.result))
		})
	}
}

func TestParsePromptContent(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		inputTopLevel bool
		result        Prompt
	}{
		{"<lora:file>", "<lora:file>", true, Prompt{kind: "lora", filename: "file"}},
		{"<abc:xyz>", "<abc:xyz>", true, Prompt{}},
		{"<abc:1.5>", "<abc:1.5>", true, Prompt{}},
		{"abc:1.5", "abc:1.5", true, Prompt{
			kind:       "ew",
			weight:     1.5,
			weightText: "1.5",
			contents: []Prompt{
				{
					kind:   "tag",
					name:   "abc",
					tokens: []string{"abc"},
				},
			},
		},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := NewTokenReader(TokenizePrompt(test.input))
			result, err := parsePromptContent(reader, test.inputTopLevel)
			assert.Equal(t, nil, err)
			assert.True(t, reflect.DeepEqual(result, test.result))
		})
	}
}
