package script

import (
	"mssfmt/read"
	"testing"
)

func TestReplace(t *testing.T) {
	script := mockScriptShort()
	// 	fmt.Println(script)

	newWords := []string{"SELECT", " "}
	f1 := WordFlag{false, MainKeyword{true, "select", 0, 1}, 1, 0, 0, 5,
		SelectColList{false, 0, 0}}
	f2 := WordFlag{false, MainKeyword{true, "select", 0, 1}, 1, 0, 6, 6,
		SelectColList{false, 0, 0}}
	newFlags := []WordFlag{f1, f2}
	script.Replace(0, 3, newWords, newFlags)

	const newContent = `SELECT * FROM   tableName`

	if script.RawContent != newContent {
		t.Errorf("Expected [%s], got [%s]", newContent, script.RawContent)
	}

	if (*script.Flags)[2].IsMainKeyword.Is {
		t.Errorf("2nd word [%s] shouldn't be a main keyword", script.Words[2])
	}

	newWords2 := []string{" "}
	f3 := WordFlag{false, MainKeyword{true, "from", 3, 5}, 1, 0, 15, 15,
		SelectColList{false, 0, 0}}
	script.Replace(5, 7, newWords2, []WordFlag{f3})

	const newContent2 = `SELECT * FROM tableName`

	if script.RawContent != newContent2 {
		t.Errorf("Expected [%s], got [%s]", newContent2, script.RawContent)
	}

	if (*script.Flags)[6].IsMainKeyword.Is {
		t.Errorf("6th word [%s] shouldn't be a main keyword. Flags: %v",
			script.Words[6], (*script.Flags)[6])
	}
}

// Prepares SQL script based on short query.
func mockScriptShort() SQL {
	content := `Select   * FROM   tableName`
	rawScript := read.RawScript{"short.sql", "./short.sql", content}
	return ToSQL(rawScript)
}
