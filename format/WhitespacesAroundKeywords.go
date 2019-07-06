package format

import (
	"mssfmt/read"
)

type WsAroundKeywords struct {
	InNewLine    bool
	NewLineAfter bool
	IndentAfter  bool
}

func (w WsAroundKeywords) Format(script *read.Script) {
}
