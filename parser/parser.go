package parser

import (
	"mssfmt/ast"
	"mssfmt/scanner"
	"mssfmt/token"
)

type Words []ast.Word

// Function ScanWords scans till EOF and accumulates all SQL tokens and literals
// as Words.
func ScanWords(s scanner.Scanner) Words {
	words := make(Words, 0, 1000)
	for {
		tok, litt := s.Scan()
		if tok == token.EOF {
			return words
		}
		words = append(words, ast.Word{tok, litt})
	}
}

// TODO
type Parser struct {
	fileName string
	source   Words
	word     ast.Word
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
		p.word = ast.Word{token.EOF, ""}
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

// Method peek returns next Word in the script but don't move forward.
// Peek peeks next non-comment token.
func (p *Parser) peek() ast.Word {
	i := 1

	for {
		if p.offset+i == len(p.source) {
			return ast.Word{token.EOF, ""}
		}
		if p.source[p.offset+i].Token == token.COMMENT {
			i++
			continue
		}
		return p.source[p.offset+i]
	}
	return ast.Word{token.EOF, ""}
}
