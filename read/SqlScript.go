package read

import "github.com/DSkrzypiec/sqlfmt/stringF"

// Script is the main object representing T-SQL script.
// Field RawContent is a content of the script as a single string.
// TODO:Extend description
type Script struct {
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
	IsMainKeyword bool
	LineNumber    int
	LineIndentLvl int
	CharIdStart   int // Start index of this word in the RawContent
	CharIdEnd     int // End index of this word in the RawContent
}

// ToScript method convertes RawScript into Script object.
func (rs RawScript) ToScript() Script {
	sFlags := ScriptFlags{}
	script := Script{rs.Name, rs.FullPath, rs.Content,
		stringF.SplitWithSep(rs.Content), &sFlags}

	script.initFlags()
	script.MarkMainKeywords()
	script.MarkLineNumbers()
	script.MarkLineIndentLvl()
	script.MarkComments()

	return script
}
