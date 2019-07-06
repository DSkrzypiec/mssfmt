package format

import (
	"strings"

	"mssfmt/read"
	"mssfmt/stringF"
)

// Keywords struct is a type for satisfying Formatter interface for formatting
// T-SQL keywords. It's field ToUpper is configuration info.
type Keywords struct {
	ToUpper bool
}

// Method Format of object Keywords is a function for satisfying Formatter
// interface. This method is a root function for formating T-SQL keywords.
// It search keywords ignoring comments, whitespaces and some special
// characters.
func (k Keywords) Format(script *read.Script) {
	for wordId, word := range script.Words {
		isKeyword, formatted := isKeyword(word, k.ToUpper)
		if isKeyword && !(*script.Flags)[wordId].IsComment {
			script.Words[wordId] = formatted
		}
	}
	script.RawContent = strings.Join(script.Words, "")
}

// Function isKeyword checks if given "word" is a T-SQL keyword and if it is
// then also returns it's formatted version.
func isKeyword(word string, toUpper bool) (bool, string) {
	wordAndNoise := stringF.SplitWithSep(word) // TODO: check if this is redundant after using SplitWithSep in constr.
	isKey := false

	for id, w := range wordAndNoise {
		_, isKWord := reservedKeywords[strings.ToLower(w)]
		if isKWord {
			isKey = true
			if toUpper {
				wordAndNoise[id] = strings.ToUpper(w)
			} else {
				wordAndNoise[id] = strings.ToLower(w)
			}
		}
	}
	return isKey, strings.Join(wordAndNoise, "")
}

// reservedKeywords is a dictionary of reserved T-SQL key words.
var reservedKeywords map[string]bool = map[string]bool{
	"add": true, "external": true, "procedure": true,
	"all": true, "fetch": true, "public": true,
	"alter": true, "file": true, "raiserror": true,
	"and": true, "fillfactor": true, "read": true,
	"any": true, "for": true, "readtext": true,
	"as": true, "foreign": true, "reconfigure": true,
	"asc": true, "freetext": true, "references": true,
	"authorization": true, "freetexttable": true, "replication": true,
	"backup": true, "from": true, "restore": true,
	"begin": true, "full": true, "restrict": true,
	"between": true, "function": true, "return": true,
	"break": true, "goto": true, "revert": true,
	"browse": true, "grant": true, "revoke": true,
	"bulk": true, "group": true, "right": true,
	"by": true, "having": true, "rollback": true,
	"cascade": true, "holdlock": true, "rowcount": true,
	"case": true, "identity": true, "rowguidcol": true,
	"check": true, "identity_insert": true, "rule": true,
	"checkpoint": true, "identitycol": true, "save": true,
	"close": true, "if": true, "schema": true,
	"clustered": true, "in": true, "securityaudit": true,
	"coalesce": true, "index": true, "select": true,
	"collate": true, "inner": true, "semantickeyphrasetable": true,
	"column": true, "insert": true, "semanticsimilaritydetailstable": true,
	"commit": true, "intersect": true, "semanticsimilaritytable": true,
	"compute": true, "into": true, "session_user": true,
	"constraint": true, "is": true, "set": true,
	"contains": true, "join": true, "setuser": true,
	"containstable": true, "key": true, "shutdown": true,
	"continue": true, "kill": true, "some": true,
	"convert": true, "left": true, "statistics": true,
	"create": true, "like": true, "system_user": true,
	"cross": true, "lineno": true, "table": true,
	"current": true, "load": true, "tablesample": true,
	"current_date": true, "merge": true, "textsize": true,
	"current_time": true, "national": true, "then": true,
	"current_timestamp": true, "nocheck": true, "to": true,
	"current_user": true, "nonclustered": true, "top": true,
	"cursor": true, "not": true, "tran": true,
	"database": true, "null": true, "transaction": true,
	"dbcc": true, "nullif": true, "trigger": true,
	"deallocate": true, "of": true, "truncate": true,
	"declare": true, "off": true, "try_convert": true,
	"default": true, "offsets": true, "tsequal": true,
	"delete": true, "on": true, "union": true,
	"deny": true, "open": true, "unique": true,
	"desc": true, "opendatasource": true, "unpivot": true,
	"disk": true, "openquery": true, "update": true,
	"distinct": true, "openrowset": true, "updatetext": true,
	"distributed": true, "openxml": true, "use": true,
	"double": true, "option": true, "user": true,
	"drop": true, "or": true, "values": true,
	"dump": true, "order": true, "varying": true,
	"else": true, "outer": true, "view": true,
	"end": true, "over": true, "waitfor": true,
	"errlvl": true, "percent": true, "when": true,
	"escape": true, "pivot": true, "where": true,
	"except": true, "plan": true, "while": true,
	"exec": true, "precision": true, "with": true,
	"execute": true, "primary": true, "within group": true,
	"exists": true, "print": true, "writetext": true,
	"exit": true, "proc": true, "sum": true,
	"avg": true, "partition": true, "min": true,
	"max": true, "count": true, "recompile": true,
}
