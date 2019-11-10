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

// Test for parsing "WITH TIES" keyword. It's a bit special because "WITH" is
// also a single word keyword.
func TestParseWithTies(t *testing.T) {
	src1 := []byte("With ties")
	src2 := []byte("with /* comment */ TIES")
	src3 := []byte("with x as ()")
	var s1, s2, s3 scanner.Scanner
	var p1, p2, p3 Parser
	s1.Init("s1", src1)
	s2.Init("s2", src2)
	s3.Init("s3", src3)
	words1 := ScanWords(s1)
	words2 := ScanWords(s2)
	words3 := ScanWords(s3)
	p1.Init("p1", words1)
	p2.Init("p2", words2)
	p3.Init("p3", words3)

	if p1.word.Token != token.WITH || p1.peek().Token != token.TIES {
		t.Errorf("Expected tokens [WITH TIES], got [%s %s]",
			p1.word.Token, p1.peek().Token)
	}
	if p2.word.Token != token.WITH || p2.peek().Token != token.TIES {
		t.Errorf("Expected tokens [WITH TIES], got [%s %s]",
			p2.word.Token, p2.peek().Token)
	}
	if p3.word.Token != token.WITH || p3.peek().Token == token.TIES {
		t.Errorf("Expected token WITH withou TIES, got [%s %s]",
			p3.word.Token, p3.peek().Token)
	}
}

// Test for peek method. In particular skipping comments.
func TestParsePeek(t *testing.T) {
	src := []byte("select top /* comment */ /* another*/ 10 * from x")
	var s scanner.Scanner
	var p Parser
	s.Init("s", src)
	words := ScanWords(s)
	p.Init("p", words)
	p.next()

	if p.word.Token != token.TOP {
		t.Errorf("Expected token TOP, got: [%s]", p.word.Token)
	}
	peek := p.peek()
	if peek.Token != token.INT {
		t.Errorf("Expected token INT, got: [%s]", peek.Token)
	}
	if peek.Literal != "10" {
		t.Errorf("Expected token TOP, got: [%s]", peek.Literal)
	}
}
