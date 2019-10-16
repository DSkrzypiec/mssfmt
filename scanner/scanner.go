package scanner

import (
	"mssfmt/token"
	"unicode/utf8"
)

const bom = 0xFEFF // byte order mark, only permitted as very first character

type Scanner struct {
	fileName string
	source   []byte // source content that is being scanned
	char     rune   // current character
	offset   int    // character offset
	rdOffset int    // reading offset - position after current char
}

// TODO...
func (s *Scanner) Scan() (token.Pos, token.Token, string) {
	// MAIN METHOD
	return 0, token.AND, ""
}

// Init method prepares Scanner for start of the source file for scanning its
// content.
func (s *Scanner) Init(fName string, src []byte) {
	s.fileName = fName
	s.source = src
	s.char = ' '
	s.offset = 0
	s.rdOffset = 0

	s.next()
	if s.char == bom {
		s.next() // ignore BOM at file beginning
	}
}

// Method next reads next Unicode character into s.char. Case when s.char < 0
// means EOF.
func (s *Scanner) next() {
	if s.rdOffset < len(s.source) {
		s.offset = s.rdOffset
		r := rune(s.source[s.rdOffset])
		w := 1

		switch {
		case r == 0:
			// s.error(s.offset, "illegal character NUL")
		case r >= utf8.RuneSelf:
			// not ASCII
			r, w = utf8.DecodeRune(s.source[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				// TODO: how to handle errors? s.error(s.offset, "illegal UTF-8 encoding")
			} else if r == bom && s.offset > 0 {
				// TODO: how to habdle errors? s.error(s.offset, "illegal byte order mark")
			}
		}
		s.rdOffset += w
		s.char = r
		return
	}

	s.offset = len(s.source)
	s.char = -1 // eof
}

// Method skipWhitespace skips all whitespace until first non whitespace
// character.
func (s *Scanner) skipWhitespace() {
	for s.char == ' ' || s.char == '\t' || s.char == '\n' || s.char == '\r' {
		s.next()
	}
}
