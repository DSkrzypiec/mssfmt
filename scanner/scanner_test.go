package scanner

import (
	"testing"
)

func TestScanSQLString(t *testing.T) {
	src := []byte("'stringValue'  'with '' escape'   'Cox'''")
	var s Scanner
	s.Init("s", src)

	string1 := s.scanSQLString()
	s.skipWhitespace()
	string2 := s.scanSQLString()
	s.skipWhitespace()
	string3 := s.scanSQLString()

	if string1 != "'stringValue'" {
		t.Errorf("Expected 'stringValue', got: %s", string1)
	}
	if string2 != "'with '' escape'" {
		t.Errorf("Expected 'with '' escape', got: %s", string2)
	}
	if string3 != "'Cox'''" {
		t.Errorf("Expected 'Cos''', got: %s", string3)
	}
}

func TestIdentifierMany(t *testing.T) {
	src1 := []byte("       GOSIA    \t DamiansTable\n")
	var s Scanner
	s.Init("s", src1)

	s.skipWhitespace()
	firstId := s.scanIdentifier()
	s.skipWhitespace()
	secondId := s.scanIdentifier()

	if firstId != "GOSIA" {
		t.Errorf("Expected <GOSIA>, got: <%s>", firstId)
	}

	if secondId != "DamiansTable" {
		t.Errorf("Expected <DamiansTable>, got: <%s>", secondId)
	}
}

// Tests for scanning delimited identifiers - [SELECT AS ColName] or
// "Illegal Var % Name!".
func TestIdentifierDelimited(t *testing.T) {
	src1 := []byte("[ illegal var __name]")
	src2 := []byte(`"Another illegal<!>"`)
	var s1, s2 Scanner
	s1.Init("s1", src1)
	s2.Init("s2", src2)

	id1 := s1.scanIdentifier()
	id2 := s2.scanIdentifier()
	const eid1 = "[ illegal var __name]"
	const eid2 = `"Another illegal<!>"`

	if id1 != eid1 {
		t.Errorf("Expected identifier <%s>, got <%s>", eid1, id1)
	}
	if id2 != eid2 {
		t.Errorf("Expected identifier <%s>, got <%s>", eid2, id2)
	}
}

// Tests for scanning regular identifiers in T-SQL.
func TestIdentifierRegular(t *testing.T) {
	src1 := []byte("VarNa#m3")
	src2 := []byte("Ciąg ")
	var s1, s2 Scanner
	s1.Init("s1", src1)
	s2.Init("s2", src2)

	id1 := s1.scanIdentifier()
	id2 := s2.scanIdentifier()
	const eid1 = "VarNa#m3"
	const eid2 = "Ciąg"

	if id1 != eid1 {
		t.Errorf("Expected identifier [%s], got [%s]", eid1, id1)
	}
	if id2 != eid2 {
		t.Errorf("Expected identifier [%s], got [%s]", eid2, id2)
	}

}

func TestSkipWhitespace(t *testing.T) {
	src1 := []byte(" X")
	src2 := []byte(" \t  \nX")
	src3 := []byte("XY")

	ss := make([]Scanner, 3)

	ss[0].Init("s1", src1)
	ss[1].Init("s2", src2)
	ss[2].Init("s3", src3)

	expOffset := []int{1, 5, 0}
	expChar := []byte("XXX")

	for i, s := range ss {
		s.skipWhitespace()

		if s.offset != expOffset[i] {
			t.Errorf("Expected offset: %d, got: %d \n", expOffset[i], s.offset)
		}
		if byte(s.char) != expChar[i] {
			t.Errorf("Expected char: %q, got: %q \n", s.char, expChar[i])
		}
	}

}

func TestPeek(t *testing.T) {
	src := []byte("SELECT")
	var s Scanner
	s.Init("s1", src)

	if string(s.peek()) != "E" {
		t.Errorf("Expected 'E', got: '%s'", string(s.peek()))
	}

	s.next()
	if string(s.peek()) != "L" {
		t.Errorf("Expected 'L', got: '%s'", string(s.peek()))
	}
}
