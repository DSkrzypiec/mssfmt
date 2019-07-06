package read

import (
	"log"
	"regexp"
	"strings"
)

// Method initFlags initialize ScriptFlags.
// This function is used during Script object creation.
func (s *Script) initFlags() {
	sFlags := make(ScriptFlags, len(s.Words))
	*s.Flags = sFlags
}

// Method markLineNumbers assigns for each word it's line number.
func (s *Script) markLineNumbers() {
	lineNumber := 1

	for wordId, word := range s.Words {
		(*s.Flags)[wordId].LineNumber = lineNumber
		if word == "\n" || word == "\r" {
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
func (s *Script) markCharIds() {
	charId := -1
	for wId, word := range s.Words {
		(*s.Flags)[wId].CharIdStart = charId + 1
		(*s.Flags)[wId].CharIdEnd = charId + len(word)
		charId += len(word)
	}
}

// Method markAndFormatKeywords marks keywords in the script. Furthermore marks
// KeywordEnd flag.
func (s *Script) markMainKeywords() {
	s.markCharIds()
	mainKeywords := keywordsRegexpsForWSFormat()

	for _, keywordRegexp := range mainKeywords {
		reg := regexp.MustCompile(keywordRegexp)
		idx := reg.FindAllStringIndex(s.RawContent, -1)
		s.markIsMainKeyword(idx)
	}
}

// Method markIsMainKeyword marks IsMainKeyword flag for given result of
// FindAllStringIndex applied in script for main keyword regex.
func (s *Script) markIsMainKeyword(indexes [][]int) {
	for _, idx := range indexes {
		for wId, flag := range *s.Flags {
			if idx[0] == idx[1] && flag.CharIdStart >= idx[0] &&
				flag.CharIdStart <= idx[1] {
				(*s.Flags)[wId].IsMainKeyword = true
			}
			if idx[0] < idx[1] && flag.CharIdStart >= idx[0] &&
				flag.CharIdStart < idx[1] {
				(*s.Flags)[wId].IsMainKeyword = true
			}
			if flag.CharIdStart > idx[1] {
				break
			}
		}
	}
}

// Function markComments determines whenever current word in the script is
// comment or not. It updates scripts Flags. This function is used during
// creating Script from RawScript.
func (s *Script) markComments() {
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

// Method markLineIndentLvl marks, for each word, flag "LineIndentLvl". That
// means it assign for each word in the line it's level of indentation <=>
// number of TABs at the front of line.
func (s *Script) markLineIndentLvl() {
	linesIndentLvls := s.lineIndentDepth()

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
func (s *Script) lineIndentDepth() map[int]int {
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
