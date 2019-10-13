package token

type Position struct {
	FileName string
	Line     int // starting at 1
	Column   int // starting at 1
}

// IsValid method checks if Position is legal.
func (pos Position) IsValid() bool {
	return pos.Line > 0 && pos.Column > 0
}

// String representation of token Position.
func (pos Position) String() string {
	if pos.IsValid() && pos.FileName == "" {
		return pos.Line + ":" + pos.Column
	}
	if pos.IsValid() {
		return pos.FileName + ":" + pos.Line + ":" + pos.Column
	}
	return ""
}
