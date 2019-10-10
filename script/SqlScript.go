package script

import (
	"fmt"
	"mssfmt/read"
	"mssfmt/stringF"
	"strings"
)

// Type (script).SQL is the main object representing T-SQL script.
// Field RawContent is a content of the script as a single string.
// TODO:Extend description
type SQL struct {
	Name       string
	FullPath   string
	RawContent string
	Words      []string
	Flags      *ScriptFlags
}

// ScriptFlags corresponds to Script.Words - it has the same length
// and i-th element of ScriptFlags correspond to i-th word in Script.Words.
type ScriptFlags []WordFlag

// WordFlag represents helper flags for each "word" of T-SQL script.
// Those flags can be updated but in main scenario each flag is set once during
// creation of Script object.
//
// IsComment - bool flag for being inside a comment.
//
// IsMainKeyword - main T-SQL keywords like SELECT, GROUP BY, etc.
// full list of main keywords can be found in (read/KeywordsRegexp.go).
//
// LineNumber - number of line containing this word.
//
// CharIdStart - Start index of this word in the RawContent
//
// CharIdEnd - End index of this word in the RawContent
//
// SelectList - Flag for being element of SELECT list of columns and expresions
//
type WordFlag struct {
	IsComment     bool
	IsMainKeyword MainKeyword
	LineNumber    int
	LineIndentLvl int
	CharIdStart   int
	CharIdEnd     int
	SelectList    SelectColList
}

// MainKeyword represents struct for keyword detection. If "Is" is true than
// that word is a keyword. Keyword says which one. The Keyword filed is very
// useful to not parse something like "  \t\n  Group   \t \n   by" into
// "GROUP BY" phrase. Fields WordIdStart and WordIdEnd are the boundries of the
// keyword including whitespaces.
type MainKeyword struct {
	Is          bool
	Keyword     string
	WordIdStart int
	WordIdEnd   int
}

// SelectColList is one of word flags which represents being in SELECT list of
// column names and expressions. Everthing between SELECT or SELECT TOP [n] and
// FROM or INTO is SelectColList.
type SelectColList struct {
	Is          bool
	WordIdStart int
	WordIdEnd   int
}

// ToSQL method convertes RawScript into Script object.
func ToSQL(rs read.RawScript) SQL {
	sFlags := ScriptFlags{}
	script := SQL{rs.Name, rs.FullPath, rs.Content,
		stringF.SplitWithSep(rs.Content), &sFlags}

	script.InitFlags()
	script.MarkMainKeywords()
	script.MarkLineNumbers()
	script.MarkLineIndentLvl()
	script.MarkComments()

	return script
}

// Replace method replaces words and its flags from word number wIdFrom until
// wIdTo from SQL script. Replacement can be any size. Words and flags to be
// replaced must have the same length. This method replaces words, flags and
// updates RowContent based on new words. It should be usually used to
// formatting the script.
func (s *SQL) Replace(wIdFrom, wIdTo int, newWords []string,
	newFlags ScriptFlags) {
	maxSize := len(s.Words) + len(newWords)
	words := make([]string, 0, maxSize)
	flags := make(ScriptFlags, 0, maxSize)
	newWordsAdded := false

	for wId, word := range s.Words {
		if !newWordsAdded && wId >= wIdFrom && wId <= wIdTo {
			for nwId, newWord := range newWords {
				words = append(words, newWord)
				flags = append(flags, newFlags[nwId])
			}
			newWordsAdded = true
		}
		if newWordsAdded && wId >= wIdFrom && wId <= wIdTo {
			continue
		}

		words = append(words, word)
		flags = append(flags, (*s.Flags)[wId])
	}

	s.Words = words
	s.Flags = &flags
	s.RawContent = strings.Join(words, "")
}

// String method returns Script in form of all words and its flags. This
// function is rather for development and debugging then for production use.
func (s SQL) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("SQL Script: %s\n", s.FullPath))
	sb.WriteString("Content (words) with description flags:\n\n")

	for wId, word := range s.Words {
		flags := (*s.Flags)[wId]
		sb.WriteString(fmt.Sprintf("[%d] {%s} ", wId, word))
		sb.WriteString("{")

		if flags.IsComment {
			sb.WriteString("comment, ")
		}
		if flags.IsMainKeyword.Is {
			sb.WriteString(fmt.Sprintf("keyword (%s(%d-%d)), ",
				flags.IsMainKeyword.Keyword, flags.IsMainKeyword.WordIdStart,
				flags.IsMainKeyword.WordIdEnd))
		}

		if flags.SelectList.Is {
			sb.WriteString("SELECT colList, ")
		}

		fStr := fmt.Sprintf("#Line=%d, Indent=%d, (%d, %d)", flags.LineNumber,
			flags.LineIndentLvl, flags.CharIdStart, flags.CharIdEnd)
		sb.WriteString(fStr)
		sb.WriteString("}\n")
	}

	return sb.String()
}
