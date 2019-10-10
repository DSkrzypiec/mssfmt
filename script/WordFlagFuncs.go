package script

import (
	"log"
	"regexp"
	"strings"
)

// Method InitFlags initialize ScriptFlags.
// This function is used during Script object creation.
func (s *SQL) InitFlags() {
	sFlags := make(ScriptFlags, len(s.Words))
	*s.Flags = sFlags
}

// Method MarkLineNumbers assigns for each word it's line number.
func (s *SQL) MarkLineNumbers() {
	lineNumber := 1

	for wordId, word := range s.Words {
		(*s.Flags)[wordId].LineNumber = lineNumber
		if word == "\n" {
			lineNumber++
		}
	}
}

// MarkCharIds for each word in the script marks it's start and end character
// indexes.
//
// Example: For script content "select top 10 * from table"
// Words[0] = "select" {CharIdStart = 0, CharIdEnd = 5}
// Words[1] = " " {CharIdStart = 6, CharIdEnd = 6}
func (s *SQL) MarkCharIds() {
	charId := -1
	for wId, word := range s.Words {
		(*s.Flags)[wId].CharIdStart = charId + 1
		(*s.Flags)[wId].CharIdEnd = charId + len(word)
		charId += len(word)
	}
}

// Method MarkMainKeywords marks keywords in the script. Furthermore marks
// KeywordEnd flag.
func (s *SQL) MarkMainKeywords() {
	s.MarkCharIds()
	mainKeywords := keywordsRegexpsForWSFormat()

	for keywordRegexp, keyword := range mainKeywords {
		reg := regexp.MustCompile(keywordRegexp)
		idx := reg.FindAllStringIndex(s.RawContent, -1)
		s.markIsMainKeyword(idx, keyword)
	}
}

// Method markIsMainKeyword marks IsMainKeyword flag for given result of
// FindAllStringIndex applied in script for main keyword regex.
func (s *SQL) markIsMainKeyword(indexes [][]int, keyword string) {
	for _, idx := range indexes {

		for wId, flag := range *s.Flags {
			if idx[0] == idx[1] && flag.CharIdStart >= idx[0] &&
				flag.CharIdStart <= idx[1] {

				wordIdStart := s.wordId(idx[0])
				wordIdEnd := s.wordId(idx[1])
				(*s.Flags)[wId].IsMainKeyword = MainKeyword{true, keyword,
					wordIdStart, wordIdEnd}
			}
			if idx[0] < idx[1] && flag.CharIdStart >= idx[0] &&
				flag.CharIdStart < idx[1] {

				wordIdStart := s.wordId(idx[0])
				wordIdEnd := s.wordId(idx[1] - 1)
				(*s.Flags)[wId].IsMainKeyword = MainKeyword{true, keyword,
					wordIdStart, wordIdEnd}
			}
			if flag.CharIdStart > idx[1] {
				break
			}
		}
	}
}

// CharId to WordId translator. If given charId couldn't be found that -1 is
// returned.
func (s *SQL) wordId(charId int) int {
	for wId, f := range *s.Flags {
		if f.CharIdStart == f.CharIdEnd && charId == f.CharIdStart {
			return wId
		}
		if f.CharIdStart < f.CharIdEnd && charId >= f.CharIdStart &&
			charId < f.CharIdEnd {
			return wId
		}
	}
	return -1
}

// Function MarkComments determines whenever current word in the script is
// comment or not. It updates scripts Flags. This function is used during
// creating Script from RawScript.
func (s *SQL) MarkComments() {
	isGlobalComment := false
	isLineComment := false

	for wordId, word := range s.Words {
		inlineComm, _ := regexp.MatchString("^--", word)
		globalStart, _ := regexp.MatchString("^/\\*", word)
		globalEnd, _ := regexp.MatchString("\\*/$", word)

		if !isLineComment && inlineComm {
			isLineComment = true
		}
		if word == "\n" || word == "\r" {
			isLineComment = false
		}
		if !isLineComment && globalStart {
			isGlobalComment = true
		}
		if !isLineComment && isGlobalComment && word == "*/" {
			isGlobalComment = false
		}

		(*s.Flags)[wordId].IsComment = isLineComment || isGlobalComment

		// 'word*/' is still a comment, therefore isGlobalComment cannot be set
		// to false before IsComment flag assignment, becuase then it won't be a
		// comment
		if !isLineComment && isGlobalComment && globalEnd {
			isGlobalComment = false
		}
	}
}

// Method MarkLineIndentLvl marks, for each word, flag "LineIndentLvl". That
// means it assign for each word in the line it's level of indentation <=>
// number of TABs at the front of line.
func (s *SQL) MarkLineIndentLvl() {
	linesIndentLvls := s.LineIndentDepth()

	for wId, wFlag := range *s.Flags {
		ll, ok := linesIndentLvls[wFlag.LineNumber]
		if !ok {
			// TODO!
			log.Printf("Couldn't find LineIndevntLevel for line = %d \n.",
				wFlag.LineNumber)
		}
		(*s.Flags)[wId].LineIndentLvl = ll
	}
}

// Method LineIndentDepth returns a slice of indentation depth for each line in
// given T-SQL script.
//
// Example: For the following script
//
// [1.] select
// [2.]		*
// [3.] from
// [4.]		tableName t
// [5.] where
// [6.]			t.A = 1
//
// Method returns {1:0, 2:1, 3:0, 4:1, 5:0, 6:2}
func (s *SQL) LineIndentDepth() map[int]int {
	lineDepth := make(map[int]int, 0)

	for id, e := range strings.Split(s.RawContent, "\n") {
		lineDepth[id+1] = calcLineDepth(e)
	}
	return lineDepth
}

// Function calcLineDepth returns level of indentation for given script "line".
func calcLineDepth(line string) int {
	depth := 0

	for _, char := range line {
		if char != '\t' {
			return depth
		}
		depth += 1
	}
	return depth
}

// MarkSelectList method marks the SELECT column list.
func (s *SQL) MarkSelectList() {
	starts := s.findSelectListStarts()
	ends := s.findSelectListEnds()
	startEndMap := matchStartAndEnds(starts, ends)

	for startId, endId := range startEndMap {
		for i := startId; i <= endId; i++ {
			(*s.Flags)[i].SelectList = SelectColList{true, startId, endId}
		}
	}
}

// Method findSelectListStarts finds and returns ids of words which starts
// select column list in SQL script.
func (s *SQL) findSelectListStarts() []int {
	startIds := make([]int, 0)

	for wId, wFlag := range *s.Flags {
		if wId == 0 {
			continue
		}

		isPrevSelect := (*s.Flags)[wId-1].IsMainKeyword.Keyword == "select"
		if isPrevSelect && (!wFlag.IsMainKeyword.Is ||
			wFlag.IsMainKeyword.Keyword != "select") {
			startIds = append(startIds, wId)
		}
	}

	return startIds
}

// Method fundSelectListEnds finds and returns ids of words which ends select
// column list in SQL script.
func (s *SQL) findSelectListEnds() []int {
	endIds := make([]int, 0)

	for wId, wFlag := range *s.Flags {
		if wId == 0 {
			continue
		}

		isIntoKeyword := strings.ToLower(s.Words[wId]) == "into"
		isFromKeyword := wFlag.IsMainKeyword.Keyword == "from"

		if !wFlag.IsComment && (isIntoKeyword || isFromKeyword) {
			endIds = append(endIds, wId-1)
		}
	}

	return endIds
}

// This function matches start and ends and returns it in form of map of startId
// onto endIds.
func matchStartAndEnds(starts, ends []int) map[int]int {
	matches := make(map[int]int)

	for _, startId := range starts {
		for _, endId := range ends {
			if endId >= startId {
				// takes first endId after startId
				matches[startId] = endId
				break
			}
		}
	}
	return matches
}
