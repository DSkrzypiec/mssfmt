package format

import (
	"fmt"
	"mssfmt/script"
	"strings"
)

type MainKeywords struct {
	InNewLine    bool
	NewLineAfter bool
	Uppercase    bool
}

// Format method formats SQL main keywords and whitespaces around them.
func (mk MainKeywords) Format(script *script.SQL) {
	newWords := make([]string, 0, len(script.Words))
	flags := *script.Flags
	skipTill := -1

	for wId, word := range script.Words {
		if wId <= skipTill {
			continue
		}

		flagKey := flags[wId].IsMainKeyword
		if flagKey.Is && !flags[wId].IsComment {
			keywordFmt := mk.KeywordFormat(flagKey.Keyword)
			newWords = append(newWords, keywordFmt)
			skipTill = flagKey.WordIdEnd
			continue
		}
		newWords = append(newWords, word)
	}

	script.RawContent = strings.Join(newWords, "")
	script.Words = newWords
	script.InitFlags()
	script.MarkMainKeywords()
	script.MarkLineNumbers()
	script.MarkLineIndentLvl()
	script.MarkComments()
}

// Method KeywordFormat formats given keyword according to MainKeywords
// configuration.
func (mk MainKeywords) KeywordFormat(keyword string) string {
	var fKeyword string = keyword

	if len(keyword) <= 2 {
		return keyword // TODO: change it. Its workaround for case -- comment select \n\n\t\n
	}
	if mk.Uppercase {
		fKeyword = strings.ToUpper(fKeyword)
	}
	if mk.InNewLine {
		fKeyword = fmt.Sprintf("\n%s", fKeyword)
	}
	if mk.NewLineAfter {
		fKeyword = fmt.Sprintf("%s\n", fKeyword)
	}
	return fKeyword
}
