package format

import (
	"mssfmt/script"
	"strings"
)

type WsAroundKeywords struct {
	InNewLine    bool
	NewLineAfter bool
	IndentAfter  bool
}

// Format method formats SQL keywords and whitespaces around them.
func (w WsAroundKeywords) Format(script *script.SQL) {
}

func indent(level int) string {
	if level <= 0 {
		return ""
	}

	s := make([]string, level*4)
	for i := 0; i < level*4; i++ {
		s[i] = " "
	}
	return strings.Join(s, "")
}
