package parser

import (
	"reflect"
	"testing"

	"github.com/junte/stable-diffusion-prompt-parser/src/reader"
	"github.com/stretchr/testify/assert"
)

func TestEscapeToken(t *testing.T) {
	parser := NewPromptParser()

	token := "\\abc"
	assert.Equal(t, "abc", parser.escapeToken(token))

	token = "\\\\abc"
	assert.Equal(t, "\\abc", parser.escapeToken(token))
}

func TestParseTagPrompt(t *testing.T) {
	tests := []struct {
		input  string
		result prompt
	}{
		{
			"abc",
			prompt{
				kind:   "tag",
				name:   "abc",
				tokens: []string{"abc"},
			},
		},
		{
			"abc xyz",
			prompt{
				kind:   "tag",
				name:   "abc xyz",
				tokens: []string{"abc", "xyz"},
			},
		},
	}

	parser := NewPromptParser()

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := reader.NewTokenReader(test.input)
			result, err := parser.parseTagPrompt(reader)
			assert.Equal(t, nil, err)
			assert.True(t, reflect.DeepEqual(*result, test.result))
		})
	}

	input := "(abc)"
	t.Run(input, func(t *testing.T) {
		reader := reader.NewTokenReader(input)
		result, err := parser.parseTagPrompt(reader)
		assert.EqualError(t, err, "tag expected")
		assert.True(t, reflect.DeepEqual(*result, prompt{}))
	})
}

func TestParsePositivePrompt(t *testing.T) {
	tests := []struct {
		input  string
		result prompt
	}{
		{
			"(abc)",
			prompt{
				kind: "pw",
				contents: []*prompt{
					{
						kind:   "tag",
						name:   "abc",
						tokens: []string{"abc"},
					},
				},
			},
		},
		{
			"((abc))",
			prompt{
				kind: "pw",
				contents: []*prompt{
					{
						kind: "pw",
						contents: []*prompt{
							{
								kind:   "tag",
								name:   "abc",
								tokens: []string{"abc"},
							},
						},
					},
				},
			},
		},
		{
			"(abc, xyz)",
			prompt{
				kind: "pw",
				contents: []*prompt{
					{
						kind:   "tag",
						name:   "abc",
						tokens: []string{"abc"},
					},
					{
						kind:   "tag",
						name:   "xyz",
						tokens: []string{"xyz"},
					},
				},
			},
		},
		{
			"(abc:xyz)",
			prompt{
				kind: "pw",
				contents: []*prompt{
					{
						kind:   "tag",
						name:   "abc",
						tokens: []string{"abc"},
					}, {
						kind:   "tag",
						name:   "xyz",
						tokens: []string{"xyz"},
					},
				},
			},
		},
		{
			"(abc:1.5)",
			prompt{
				kind:   "cw",
				weight: 1.5,
				contents: []*prompt{
					{
						kind:   "tag",
						name:   "abc",
						tokens: []string{"abc"},
					},
				},
			},
		},
		{
			"(abc:.1.5)",
			prompt{
				kind:   "cw",
				weight: 1.5,
				contents: []*prompt{
					{
						kind:   "tag",
						name:   "abc",
						tokens: []string{"abc"},
					},
				},
			},
		},
		{
			"(abc:1.5.1)",
			prompt{
				kind:   "cw",
				weight: 1.5,
				contents: []*prompt{
					{
						kind:   "tag",
						name:   "abc",
						tokens: []string{"abc"},
					},
				},
			},
		},
		{
			"(abc:)",
			prompt{
				kind: "pw",
				contents: []*prompt{
					{
						kind:   "tag",
						name:   "abc",
						tokens: []string{"abc"},
					},
				},
			},
		},
		{
			"(abc:1.5, xyz)",
			prompt{
				kind:   "cw",
				weight: 1.5,
				contents: []*prompt{
					{
						kind:   "tag",
						name:   "abc",
						tokens: []string{"abc"},
					},
				},
			},
		},
	}

	parser := NewPromptParser()

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := reader.NewTokenReader(test.input)
			result, err := parser.parsePositivePrompt(reader)
			assert.Equal(t, nil, err)
			assert.True(t, reflect.DeepEqual(*result, test.result))
		})
	}
}

func TestParseNegativePrompt(t *testing.T) {
	tests := []struct {
		input  string
		result prompt
	}{
		{
			"[abc]",
			prompt{
				kind: "nw",
				contents: []*prompt{
					{
						kind:   "tag",
						name:   "abc",
						tokens: []string{"abc"},
					},
				},
			},
		},
		{
			"[[abc]]",
			prompt{
				kind: "nw",
				contents: []*prompt{
					{
						kind: "nw",
						contents: []*prompt{
							{
								kind:   "tag",
								name:   "abc",
								tokens: []string{"abc"},
							},
						},
					},
				},
			},
		},
		{
			"[abc, xyz]",
			prompt{
				kind: "nw",
				contents: []*prompt{
					{
						kind:   "tag",
						name:   "abc",
						tokens: []string{"abc"},
					},
					{
						kind:   "tag",
						name:   "xyz",
						tokens: []string{"xyz"},
					},
				},
			},
		},
	}

	parser := NewPromptParser()

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := reader.NewTokenReader(test.input)
			result, err := parser.parseNegativePrompt(reader)
			assert.Equal(t, nil, err)
			assert.True(t, reflect.DeepEqual(*result, test.result))
		})
	}
}

func TestParseContentToken(t *testing.T) {
	parser := NewPromptParser()

	input := "abc"
	tokenReader := reader.NewTokenReader(input)
	result, err := parser.parseContentToken(tokenReader, "Name")
	assert.Equal(t, nil, err)
	assert.Equal(t, "abc", result)

	input = "(abc)"
	tokenReader = reader.NewTokenReader(input)
	result, err = parser.parseContentToken(tokenReader, "Name")
	assert.EqualError(t, err, "Name expected")
	assert.Equal(t, "", result)
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input  string
		result float64
	}{
		{"1", 1},
		{"1.5", 1.5},
		{".5", 0.5},
		{"1. 5", 1.5},
		{"0,5", 0.5},
		{"0, 5", 0.5},
		{"1..5", 1.5},
		{".1.5", 1.5},
		{"1.5.1", 1.5},
	}

	parser := NewPromptParser()

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := reader.NewTokenReader(test.input)
			result, err := parser.parseNumber(reader, "Name")
			assert.Equal(t, nil, err)
			assert.Equal(t, result, test.result)
		})
	}
}

func TestParseAnglePrompt(t *testing.T) {
	tests := []struct {
		input     string
		inputName string
		result    prompt
	}{
		{
			"<lora:file>",
			"lora",
			prompt{
				kind:       "lora",
				filename:   "file",
				multiplier: 0.5,
			},
		},
		{
			"<hypernet:file:0.5>",
			"hypernet",
			prompt{
				kind:       "hypernet",
				filename:   "file",
				multiplier: 0.5,
			},
		},
		{
			"<lora:file name:.5>",
			"lora",
			prompt{
				kind:       "lora",
				filename:   "file name",
				multiplier: 0.5,
			},
		},
		{
			"<lora:file.name-v.1.5:.5>",
			"lora",
			prompt{
				kind:       "lora",
				filename:   "file.name-v.1.5",
				multiplier: 0.5,
			},
		},
		{
			"<lora:file.name-v.1.5:1..5>",
			"lora",
			prompt{
				kind:       "lora",
				filename:   "file.name-v.1.5",
				multiplier: 1.5,
			},
		},
		{
			"<lora:file.name-v.1.5:.1.5>",
			"lora",
			prompt{
				kind:       "lora",
				filename:   "file.name-v.1.5",
				multiplier: 1.5,
			},
		},
		{
			"<lora:file.name-v.1.5:1.5.1>",
			"lora",
			prompt{
				kind:       "lora",
				filename:   "file.name-v.1.5",
				multiplier: 1.5,
			},
		},
		{
			"<lora:file.name-v.1.5:.>",
			"lora",
			prompt{
				kind:       "lora",
				filename:   "file.name-v.1.5",
				multiplier: 0.5,
			},
		},
		{
			"<lora:file.name-v.1.5:>",
			"lora",
			prompt{
				kind:       "lora",
				filename:   "file.name-v.1.5",
				multiplier: 0.5,
			},
		},
		{
			"< lora : ?file,name[v.1]  : .5 >",
			"lora",
			prompt{
				kind:       "lora",
				filename:   "?file,name[v.1]",
				multiplier: 0.5,
			},
		},
	}

	parser := NewPromptParser()

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := reader.NewTokenReader(test.input)
			result, err := parser.parseAnglePrompt(reader, test.inputName)
			assert.Equal(t, nil, err)
			assert.True(t, reflect.DeepEqual(*result, test.result))
		})
	}
}

func TestParsePromptContent(t *testing.T) {
	tests := []struct {
		input         string
		inputTopLevel bool
		result        prompt
	}{
		{
			"<lora:file>",
			true,
			prompt{
				kind:       "lora",
				filename:   "file",
				multiplier: 0.5,
			},
		},
		{
			"<",
			true,
			prompt{},
		},
		{
			"<abc:xyz>",
			true,
			prompt{},
		},
		{
			"<abc:1.5>",
			true,
			prompt{},
		},
		{
			"abc:1.5",
			true,
			prompt{
				kind:   "cw",
				weight: 1.5,
				contents: []*prompt{
					{
						kind:   "tag",
						name:   "abc",
						tokens: []string{"abc"},
					},
				},
			},
		},
	}

	parser := NewPromptParser()

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := reader.NewTokenReader(test.input)
			result, err := parser.parsePromptContent(reader, test.inputTopLevel)
			assert.Equal(t, nil, err)
			assert.True(t, reflect.DeepEqual(*result, test.result))
		})
	}
}

func TestParsePrompt(t *testing.T) {
	tests := []struct {
		input  string
		result ParsedPrompt
	}{
		{
			"abc",
			ParsedPrompt{
				Tags: []*PromptTag{{Tag: "abc", Weight: 1}},
			},
		},
		{
			"(abc)",
			ParsedPrompt{
				Tags: []*PromptTag{{Tag: "abc", Weight: 1.1}},
			},
		},
		{
			"((abc))",
			ParsedPrompt{
				Tags: []*PromptTag{{Tag: "abc", Weight: 1.2100000000000002}},
			},
		},
		{
			"(abc:1.5)",
			ParsedPrompt{
				Tags: []*PromptTag{{Tag: "abc", Weight: 1.5}},
			},
		},
		{
			"(abc:1.5.2)",
			ParsedPrompt{
				Tags: []*PromptTag{{Tag: "abc", Weight: 1.5}},
			},
		},
		{
			"[abc]",
			ParsedPrompt{
				Tags: []*PromptTag{{Tag: "abc", Weight: 0.9090909090909091}},
			},
		},
		{
			"[[abc]]",
			ParsedPrompt{
				Tags: []*PromptTag{{Tag: "abc", Weight: 0.8264462809917354}},
			},
		},
		{
			"abc xyz",
			ParsedPrompt{
				Tags: []*PromptTag{{Tag: "abc xyz", Weight: 1}},
			},
		},
		{
			"abc, xyz",
			ParsedPrompt{
				Tags: []*PromptTag{
					{
						Tag:    "abc",
						Weight: 1,
					},
					{
						Tag:    "xyz",
						Weight: 1,
					},
				},
			},
		},
		{
			"(abc, xyz)",
			ParsedPrompt{
				Tags: []*PromptTag{
					{
						Tag:    "abc",
						Weight: 1.1,
					},
					{
						Tag:    "xyz",
						Weight: 1.1,
					},
				},
			},
		},
		{
			"(abc: xyz)",
			ParsedPrompt{
				Tags: []*PromptTag{
					{
						Tag:    "abc",
						Weight: 1.1,
					},
					{
						Tag:    "xyz",
						Weight: 1.1,
					},
				},
			},
		},
		{
			"<lora:file>",
			ParsedPrompt{
				Loras: []*PromptModel{{Filename: "file", Multiplier: 0.5}},
			},
		},
		{
			"<lora:file:>",
			ParsedPrompt{
				Loras: []*PromptModel{{Filename: "file", Multiplier: 0.5}},
			},
		},
		{
			"<lora:file:.>",
			ParsedPrompt{
				Loras: []*PromptModel{{Filename: "file", Multiplier: 0.5}},
			},
		},
		{
			"<lora:file:0.5>",
			ParsedPrompt{
				Loras: []*PromptModel{{Filename: "file", Multiplier: 0.5}},
			},
		},
		{
			"<lora:file:.1.5>",
			ParsedPrompt{
				Loras: []*PromptModel{{Filename: "file", Multiplier: 1.5}},
			},
		},
		{
			"<lora:file:1.5.2>",
			ParsedPrompt{
				Loras: []*PromptModel{{Filename: "file", Multiplier: 1.5}},
			},
		},
		{
			"<hypernet:file>",
			ParsedPrompt{
				Hypernets: []*PromptModel{{Filename: "file", Multiplier: 0.5}},
			},
		},
		{
			"<hypernet:file:1.5>",
			ParsedPrompt{
				Hypernets: []*PromptModel{{Filename: "file", Multiplier: 1.5}},
			},
		},
		{
			"abc, [[mno]], (xyz), <hypernet:file>, <lora:file:1.5>",
			ParsedPrompt{
				Tags: []*PromptTag{
					{
						Tag:    "abc",
						Weight: 1,
					},
					{
						Tag:    "mno",
						Weight: 0.8264462809917354,
					},
					{
						Tag:    "xyz",
						Weight: 1.1,
					},
				},
				Loras: []*PromptModel{
					{
						Filename:   "file",
						Multiplier: 1.5,
					},
				},
				Hypernets: []*PromptModel{
					{
						Filename:   "file",
						Multiplier: 0.5,
					},
				},
			},
		},
		{
			"abc,,,xyz",
			ParsedPrompt{
				Tags: []*PromptTag{
					{
						Tag:    "abc",
						Weight: 1,
					},
					{
						Tag:    "xyz",
						Weight: 1,
					},
				},
			},
		},
	}

	parser := NewPromptParser()

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := parser.ParsePrompt(test.input)
			assert.Equal(t, nil, err)
			assert.True(t, reflect.DeepEqual(*result, test.result))
		})
	}
}

func TestPromptToString(t *testing.T) {
	tests := []struct {
		input  string
		result string
	}{
		{"abc,,,,,xyz", "abc, xyz"},
		{"( (abc ) )", "((abc))"},
		{"[ [abc ] ]", "[[abc]]"},
		{"(abc:)", "(abc)"},
		{"(abc:.1.5)", "(abc:1.5)"},
		{"(abc:1.5.1)", "(abc:1.5)"},
		{"( abc : 1, 5 )", "(abc:1.5)"},
		{"(abc) [xyz]", "(abc), [xyz]"},
		{"(abc)mno[xyz]", "(abc), mno, [xyz]"},
		{"< lora : file name : .5 >", "<lora:file name:.5>"},
		{"< hypernet : file name>", "<hypernet:file name:.5>"},
		{"abc <lora:filename> mno xyz", "abc, <lora:filename:.5>, mno xyz"},
		{"(abc,, <lora:filename>,, <lora:filename>)", "(abc, <lora:filename:.5>, <lora:filename:.5>)"},
		{"<lora:file.name:.1.5>", "<lora:file.name:1.5>"},
		{"<lora:file.name:1.5.1>", "<lora:file.name:1.5>"},
	}

	parser := NewPromptParser()

	for _, test := range tests {
		t.Run(test.result, func(t *testing.T) {
			result, err := parser.BeautifyPrompt(test.input)
			assert.Equal(t, nil, err)
			assert.Equal(t, result, test.result)
		})
	}
}
