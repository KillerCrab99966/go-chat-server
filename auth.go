package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type codeGenerator struct {
	words   []string
	codeLen int
}

func newCodeGenerator(words []string, codeLen int) *codeGenerator {
	return &codeGenerator{
		words:   words,
		codeLen: codeLen,
	}
}

// generateCode generates a cryptographically secure code consisting of codeLen-many words hyphenated.
func (g *codeGenerator) generateCode() (string, error) {
	codeWords := make([]string, 0, g.codeLen)
	for range g.codeLen {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(g.words))))
		if err != nil {
			return "", fmt.Errorf("generating random digit: %w", err)
		}

		codeWords = append(codeWords, g.words[num.Int64()])
	}

	return strings.Join(codeWords, "-"), nil
}

type token struct {
	value   string
	room    *room
	created time.Time
}

func newToken(rm *room, tknLen int) (token, error) {
	var builder strings.Builder

	for range tknLen {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(10)))
		if err != nil {
			return token{}, fmt.Errorf("generating random digit: %w", err)
		}

		builder.WriteString(strconv.Itoa(int(num.Int64())))
	}

	return token{
		value:   builder.String(),
		room:    rm,
		created: time.Now(),
	}, nil
}
