package token

import "strconv"

// TODO...
type Position struct {
	FileName string
	Line     int // starting at 1
	Column   int // starting at 1
	Offset   int // wordId, starting at 0
}

// TODO: think about it
type Pos int

// IsValid method checks if Position is legal.
func (pos Position) IsValid() bool {
	return pos.Line > 0 && pos.Column > 0
}

// String representation of token Position.
func (pos Position) String() string {
	if pos.IsValid() && pos.FileName == "" {
		return strconv.Itoa(pos.Line) + ":" + strconv.Itoa(pos.Column)
	}
	if pos.IsValid() {
		return pos.FileName + ":" + strconv.Itoa(pos.Line) + ":" +
			strconv.Itoa(pos.Column)
	}
	return ""
}
