package format

import (
	"mssfmt/read"
)

type WsAroundKeywords struct {
	InNewLine    bool
	NewLineAfter bool
	IndentAfter  bool
}

// Format method formats SQL keywords and whitespaces around them.
func (w WsAroundKeywords) Format(script *read.Script) {
	newWords := make([]string, 0, len(script.Words))

	for wId, w := range script.Words {
		f := (*script.Flags)[wId]
		if f.IsMainKeyword && !f.IsComment {
			newWord, ok := w.formatKeyword(w, f)
			if ok {
				newWords = append(newWords, newWord)
			}
			continue
		}
		newWords = append(newWords, w)
	}

	// Script update
}

func (w WsAroundKeywords) formatKeyword(w string, flags read.ScriptFlags) (string, bool) {
	// TODO
}
