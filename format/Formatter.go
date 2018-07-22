package format

import "github.com/DSkrzypiec/mssfmt/read"

// Formatter is called any object that formats a SQL Script.
type Formatter interface {
	Format(*read.Script)
}

// ApplyFormats applies all format actions to be performed on give .sql script.
func ApplyFormats(script *read.Script, formats []Formatter) {
	for _, f := range formats {
		f.Format(script)
	}
}
