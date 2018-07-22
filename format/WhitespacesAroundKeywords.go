package format

import (
	"regexp"

	"github.com/DSkrzypiec/mssfmt/read"
	"github.com/DSkrzypiec/mssfmt/stringF"
)

type WsAroundKeywords struct {
	InNewLine    bool
	NewLineAfter bool
	IndentAfter  bool
}

// TODO: doesnt work yet. Make done algorithm for whitespaces before and after
// keywords.
func (w WsAroundKeywords) Format(script *read.Script) {
	//	words := make([]string, 0, 1000)
	//	prevWord := ""
	formatWSInMultiwordKeywords(script)
	//
	//	for wId, word := range script.Words {
	//		if !(*script.Flags)[wId].IsMainKeyword || (*script.Flags)[wId].IsComment {
	//			words = append(words, word)
	//			prevWord = word
	//			continue
	//		}
	//		for _, newWord := range w.formatSingleKeywordWS(word, prevWord, (*script.Flags)[wId]) {
	//			words = append(words, newWord)
	//		}
	//	}
	//
	//	script.Words = words
	//	script.RawContent = strings.Join(script.Words, "")
}

// This function handle multi words keyword like "group by" or "union all".
// Function parse any whitespace between to components of the keyword into
// single space.
// Example: keyword "group     \t  \n by" is converted into "group by".
func formatWSInMultiwordKeywords(script *read.Script) {
	multiKeyRegExp := read.KeywordsRegexpsForWSReplace()

	for regEx, replace := range multiKeyRegExp {
		re := regexp.MustCompile(regEx)
		script.RawContent = re.ReplaceAllString(script.RawContent, replace)
	}

	script.Words = stringF.SplitWithSep(script.RawContent)
}

func (w WsAroundKeywords) formatSingleKeywordWS(word, prevWord string,
	flags read.WordFlag) []string {

	keyWS := make([]string, 0, 5)

	if prevWord != "\n" { // TODO: doesn't work correctly for case \n [space] KEYWORD
		keyWS = append(keyWS, "\n")
		addIndents(&keyWS, flags.LineIndentLvl)
	}
	keyWS = append(keyWS, word)

	if w.NewLineAfter {
		keyWS = append(keyWS, "\n")
	}
	if w.IndentAfter {
		addIndents(&keyWS, flags.LineIndentLvl+1)
	}
	return keyWS
}

func addIndents(words *[]string, indentLvl int) {
	for i := 0; i < indentLvl; i++ {
		*words = append(*words, "\t")
	}
}
