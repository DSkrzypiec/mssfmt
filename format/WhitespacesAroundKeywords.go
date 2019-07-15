package format

import (
	"mssfmt/read"
    "strings"
)

type WsAroundKeywords struct {
	InNewLine    bool
	NewLineAfter bool
	IndentAfter  bool
}

// Format method formats SQL keywords and whitespaces around them.
func (w WsAroundKeywords) Format(script *read.Script) {
	//newWords := make([]string, 0, len(script.Words))
    //var prevWordIsKeyword bool = false


	//for wId, word := range script.Words {
	//	f := (*script.Flags)[wId]
	//	if f.IsMainKeyword && !f.IsComment {
	//	}
	//	newWords = append(newWords, word)
	//}

	//script.Words = newWords
    //script.RawContent = strings.Join(newWords, "")
}

func indent(level int) string {
    if level <= 0 {
        return ""
    }

    s := make([]string, level * 4)
    for i := 0; i < level * 4; i++ {
        s[i] = " "
    }
    return strings.Join(s, "")
}
