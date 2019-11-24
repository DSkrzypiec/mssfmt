package scanner

import (
	"mssfmt/token"
	"strings"
	"unicode"
	"unicode/utf8"
)

const bom = 0xFEFF     // byte order mark, only permitted as very first character
const singleQuote = 39 // value for single quote character
const doubleQuote = 34 // value for double quote character
const hashSign = 35    // value for '#' sign
const atSign = 64      // value for '@' sign

// Scanner represents current state of scanning .sql file char by char. In
// source filed SQL script content is stored as slice of bytes. Field char
// contains current character, offset is number of character in the file and
// rdOffset is position after current offset.
type Scanner struct {
	fileName string
	source   []byte // source content that is being scanned
	char     rune   // current character
	offset   int    // character offset
	rdOffset int    // reading offset - position after current char
}

// Scan method scans T-SQL script and returns T-SQL tokens defined in token
// package. Single call of Scan scans single token. Scan ends scanning after
// token.EOF.
func (s *Scanner) Scan() (token.Token, string) {
	var tok token.Token
	var literal string
	s.skipWhitespace()

	switch ch := s.char; {
	case isLetter(ch) || ch == hashSign || ch == atSign:
		literal = s.scanIdentifier()
		ucLit := strings.ToUpper(literal)
		if len(literal) > 1 {
			if howMany, isMulti := token.MultiwordKeywords[ucLit]; isMulti {
				return s.handleMultiwordKeyword(literal, howMany)
			}

			tok = token.KeywordLookup(strings.ToUpper(literal))
			return tok, literal
		} else {
			return token.IDENT, literal
		}

	case isDigit(ch) || (ch == '.' && isDigit(rune(s.peek()))) ||
		((ch == '-' || ch == '+') && isDigit(rune(s.peek()))):
		return s.scanNumber()
	case ch == '-' && s.peek() == '-':
		return token.COMMENT, s.scanLineComment()
	case ch == '/' && s.peek() == '*':
		return token.COMMENT, s.scanBlockComment()
	default:
		s.next()
		switch ch {
		case -1:
			return token.EOF, ""
		case singleQuote:
			return token.STRING, s.scanSQLString()
		case '+':
			return token.ADD, "+"
		case '-':
			return token.SUB, "-"
		case '*':
			return token.MUL, "*"
		case '/':
			return token.DIV, "/"
		case '%':
			return token.MOD, "%"
		case '=':
			return token.ASSIGN, "="
		case '.':
			return token.PERIOD, "."
		case ',':
			return token.COMMA, ","
		case ';':
			return token.SEMICOLON, ";"
		case '(':
			return token.LPAREN, "("
		case ')':
			return token.RPAREN, ")"
		default:
			return token.ILLEGAL, ""
		}
	}
	return token.ILLEGAL, ""
}

// Method handleMultiwordKeyword scans the rest (after first word) part of
// multi-word keyword.
func (s *Scanner) handleMultiwordKeyword(firstWord string, nWords int) (token.Token, string) {
	words := make([]string, nWords+1)
	words[0] = firstWord

	for i := 0; i < nWords; i++ {
		s.skipWhitespace()
		lit := s.scanIdentifier()
		words[i+1] = lit
	}

	keyword := strings.Join(words, " ")
	tok := token.KeywordLookup(strings.ToUpper(keyword))
	return tok, keyword
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

	if isLetter(s.char) || s.char == hashSign || s.char == atSign {
		for isLetter(s.char) || isDigit(s.char) || isSpecialInsideIden(s.char) {
			s.next()
		}
		return string(s.source[startOffset:s.offset])
	}
	return ""
}

// Method scanSQLString scans T-SQL string literal. Result also includes opening
// and closing single quote - '. It also includes single quote escapement which
// in T-SQL occurs as double single quote - ''. This metod assumes that the
// opening single quote was already scanned.
func (s *Scanner) scanSQLString() string {
	startOffset := s.offset

	for {
		s.next()
		if s.isEscapedSQ() {
			s.next()
			s.next()
		}
		if s.char == singleQuote && !s.isEscapedSQ() {
			s.next()
			break
		}
	}
	return string(s.source[startOffset:s.offset])
}

// Method scanNumber scans number literals. It includes integers, floats and
// decimals, scientific notation and hexidecimal format (TODO). This method
// assumes that s.char is a digit.
func (s *Scanner) scanNumber() (token.Token, string) {
	startOffset := s.offset
	var tok token.Token = token.INT

	for {
		s.next()
		if !isDigit(s.char) && s.char != '.' && s.char != 'e' &&
			s.char != 'E' && s.char != '+' && s.char != '-' {
			break
		}

		if s.char == '.' {
			tok = token.FLOAT
		}
	}
	return tok, string(s.source[startOffset:s.offset])
}

// Method scanLineComment scans line comment in T-SQL which starts from "--" and
// ends at line break. This method assumes that s.char == '-' and s.peek() ==
// '-', so it's a line comment start.
func (s *Scanner) scanLineComment() string {
	startOffset := s.offset

	for s.char != '\n' && s.char != '\r' {
		s.next()
	}
	return string(s.source[startOffset:s.offset])
}

// Method scanBlockComment scans block comment in T-SQL which starts from "/*"
// and ends at "*/". Block comments in T-SQL supports nested block comments.
// This method assumes that it's on start of block comment - s.char == '/' &&
// s.peek() == '*'.
func (s *Scanner) scanBlockComment() string {
	startOffset := s.offset
	nestingLvl := 1

	for !(s.char == '*' && s.peek() == '/' && nestingLvl == 0) {
		s.next()
		if s.char == '/' && s.peek() == '*' {
			nestingLvl++
		}
		if s.char == '*' && s.peek() == '/' && nestingLvl > 0 {
			nestingLvl--
		}
	}
	s.next() // to accumulate closing "*/"
	s.next() // to accumulate closing "*/"
	return string(s.source[startOffset:s.offset])
}

// Method isEscapedSQ verifies if current character (s.char) is escaped single
// quote inside T-SQL string.
func (s *Scanner) isEscapedSQ() bool {
	return s.char == singleQuote && s.peek() == singleQuote
}

func isSpecialInsideIden(char rune) bool {
	return char == '_' || char == atSign || char == hashSign
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
