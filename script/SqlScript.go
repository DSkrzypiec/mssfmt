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
// IsMainKeyword - main T-SQL keywords like SELECT, GROUP BY, etc.
// full list of main keywords can be found in (read/KeywordsRegexp.go).
//
// LineNumber - number of line containing this word.
//
// CharIdStart - TODO: description
//
// CharIdEnd - TODO: description
type WordFlag struct {
	IsComment     bool
	IsMainKeyword MainKeyword
	LineNumber    int
	LineIndentLvl int
	CharIdStart   int // Start index of this word in the RawContent
	CharIdEnd     int // End index of this word in the RawContent
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

		fStr := fmt.Sprintf("#Line=%d, Indent=%d, (%d, %d)", flags.LineNumber,
			flags.LineIndentLvl, flags.CharIdStart, flags.CharIdEnd)
		sb.WriteString(fStr)
		sb.WriteString("}\n")
	}

	return sb.String()
}
