package scanner

import (
	"testing"
)

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
