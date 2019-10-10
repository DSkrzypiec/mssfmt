package format

import (
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
	// TODO fix it after Replace method.
	//skipTill := -1
	//for wId, flag := range *script.Flags {
	//	if wId <= skipTill {
	//		continue
	//	}

	//	mKey := flag.IsMainKeyword
	//	if mKey.Is && !flag.IsComment {
	//		newWords, newFlags := mk.KeywordFormat(mKey.Keyword, wId)
	//		l := mKey.WordIdEnd - mKey.WordIdStart
	//		script.Replace(wId, wId+l, newWords, newFlags)
	//		skipTill = wId + len(newWords) - 1
	//	}
	//}

	//script.MarkCharIds()
	//script.MarkLineNumbers()
	//script.MarkLineIndentLvl()
}

// Method KeywordFormat formats given keyword according to MainKeywords
// configuration.
func (mk MainKeywords) KeywordFormat(keyword string,
	wIdStart int) ([]string, script.ScriptFlags) {

	fKeyword := make([]string, 0, 3)
	flags := make(script.ScriptFlags, 0, 3)
	maink := script.MainKeyword{true, strings.ToLower(keyword), wIdStart,
		wIdStart + 2}
	keyFlag := script.WordFlag{false, maink, 0, 0, 0, 0,
		script.SelectColList{}}

	if len(keyword) <= 2 {
		return fKeyword, flags // TODO: change it. Its workaround for case -- comment select \n\n\t\n
	}

	if mk.InNewLine {
		fKeyword = append(fKeyword, "\n")
		flags = append(flags, keyFlag)
	} else {
		fKeyword = append(fKeyword, " ")
		flags = append(flags, keyFlag)
	}

	if mk.Uppercase {
		fKeyword = append(fKeyword, strings.ToUpper(keyword))
		flags = append(flags, keyFlag)
	} else {
		fKeyword = append(fKeyword, strings.ToLower(keyword))
		flags = append(flags, keyFlag)
	}

	if mk.NewLineAfter {
		fKeyword = append(fKeyword, "\n")
		flags = append(flags, keyFlag)
	} else {
		fKeyword = append(fKeyword, " ")
		flags = append(flags, keyFlag)
	}

	return fKeyword, flags
}
