package format

import (
	"mssfmt/script"
)

// Formatter is called any object that formats a SQL Script.
type Formatter interface {
	Format(*script.SQL)
}

// ApplyFormats applies all format actions to be performed on give .sql script.
func ApplyFormats(script *script.SQL, formats []Formatter) {
	for _, f := range formats {
		f.Format(script)
	}
}
