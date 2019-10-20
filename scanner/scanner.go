package scanner

import (
	"mssfmt/token"
	"unicode"
	"unicode/utf8"
)

const bom = 0xFEFF     // byte order mark, only permitted as very first character
const singleQuote = 39 // value for single quote character
const doubleQuote = 34 // value for double quote character

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

// Method peek returns the byte following the most recently read character without
// advancing the scanner. If the scanner is at EOF, peek returns 0.
// This method is copied from Go scanner package from standard library -
// scanner/scanner.go - #line 88-95.
func (s *Scanner) peek() byte {
	if s.rdOffset < len(s.source) {
		return s.source[s.rdOffset]
	}
	return 0
}

// Method scanIdentifier scans T-SQL identifiers. Including regular one and
// delimited identifiers. Keywords and function names are just special case of
// identifiers.
func (s *Scanner) scanIdentifier() string {
	startOffset := s.offset
	if s.char == '[' {
		for s.char != ']' {
			s.next()
		}
		s.next()
		return string(s.source[startOffset:s.offset])
	}

	if s.char == doubleQuote {
		s.next()
		for s.char != doubleQuote {
			s.next()
		}
		s.next()
		return string(s.source[startOffset:s.offset])
	}

	if isLetter(s.char) {
		for isLetter(s.char) || isDigit(s.char) || isSpecialInsideIden(s.char) {
			s.next()
		}
		return string(s.source[startOffset:s.offset])
	}
	return ""
}

// Method scanSQLString scans T-SQL string literal. Result also includes opening
// and closing single quote - '. It also includes single quote escapement which
// in T-SQL occurs as double single quote - ''.
func (s *Scanner) scanSQLString() string {
	return "TODO"
}

func isSpecialInsideIden(char rune) bool {
	return char == '_' || char == '@' || char == '#'
}

func isLetter(char rune) bool {
	return 'a' <= lower(char) && lower(char) <= 'z' || char == '_' ||
		char >= utf8.RuneSelf && unicode.IsLetter(char)
}

func isDigit(char rune) bool {
	return isDecimal(char) || char >= utf8.RuneSelf && unicode.IsDigit(char)
}

func isDecimal(char rune) bool {
	return '0' <= char && char <= '9'
}

func isWhitespace(char rune) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

func isHex(char rune) bool {
	return '0' <= char && char <= '9' || 'a' <= lower(char) && lower(char) <= 'f'
}

// returns lower-case char iff ch is ASCII letter
func lower(char rune) rune {
	return ('a' - 'A') | char
}
