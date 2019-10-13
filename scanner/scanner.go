package scanner

import "mssfmt/token"

type Scanner struct {
	// SourceFile
	source   []byte // source content that is being scanned
	char     rune   // current character
	offset   int    // character offset
	rdOffset int    // reading offset - position after current char
	line     int    // current line number
}

// TODO...
func (s *Scanner) Scan() (token.Pos, token.Token, string) {
	return 0, token.AND, ""
}
