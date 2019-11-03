package parser

import (
	"mssfmt/scanner"
	"mssfmt/token"
)

// Word represents single "word" in SQL script. It's a pair of token.Token and
// corresponding literal.
type Word struct {
	Token   token.Token
	Literal string
}
type Words []Word

// Function ScanWords scans till EOF and accumulates all SQL tokens and literals
// as Words.
func ScanWords(s scanner.Scanner) Words {
	words := make(Words, 0, 1000)
	for {
		tok, litt := s.Scan()
		if tok == token.EOF {
			return words
		}
		words = append(words, Word{tok, litt})
	}
}

// TODO
type Parser struct {
	fileName string
	source   Words
	word     Word
	offset   int
}

func (p *Parser) Init(name string, src Words) {
	p.fileName = name
	p.source = src
	p.word = src[0] // TODO
	p.offset = 0
}

// Method next jumps to next Word in the SQL script.
func (p *Parser) next() {
	if p.offset+1 == len(p.source) {
		p.offset = len(p.source)
		p.word = Word{token.EOF, ""}
		return
	}

	p.offset++
	p.word = p.source[p.offset]

	// During parsing Parser ommit comments. Comment will be added
	// during printing the tree.
	if p.word.Token == token.COMMENT {
		p.next()
	}
}
