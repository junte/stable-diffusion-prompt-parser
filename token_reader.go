package main

import "errors"

type TokenReader struct {
	index  int
	tokens []Token
	length int
}

func NewTokenReader(tokens []Token) *TokenReader {
	return &TokenReader{
		index:  0,
		tokens: tokens,
		length: len(tokens),
	}
}

func (reader *TokenReader) getToken() string {
	if reader.index < reader.length {
		return reader.tokens[reader.index].value
	}
	return ""
}

func (reader *TokenReader) getMultipleTokens(count int) ([]string, error) {
	if reader.index+count <= reader.length {
		values := make([]string, count)

		for i := 0; i < count; i++ {
			values[i] = reader.tokens[reader.index+i].value
		}
		return values, nil
	}
	return nil, errors.New("count out of range")
}

func (reader *TokenReader) nextToken() {
	if reader.index < reader.length {
		reader.index++
	}
}
