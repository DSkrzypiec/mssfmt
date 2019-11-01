package parser

import (
	"mssfmt/scanner"
	"mssfmt/token"
	"testing"
)

func TestParserNext(t *testing.T) {
	src := []byte("select distinct X from Y")
	var s scanner.Scanner
	var p Parser
	s.Init("s", src)
	words := ScanWords(s)
	p.Init("p", words)

	expToken := []token.Token{token.SELECT, token.DISTINCT, token.IDENT,
		token.FROM, token.IDENT}
	expLit := []string{"select", "distinct", "X", "from", "Y"}
	i := 0

	for {
		if p.word.Token == token.EOF {
			break
		}

		if p.word.Token != expToken[i] {
			t.Errorf("Expected token [%s], got: [%s]", p.word.Token, expToken[i])
		}
		if p.word.Literal != expLit[i] {
			t.Errorf("Expected literal [%s], got: [%s]", p.word.Literal, expLit[i])
		}
		i++
		p.next()
	}
}
