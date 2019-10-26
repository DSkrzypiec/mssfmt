package token

import (
	"strings"
	"testing"
)

func TestKeywordLookup(t *testing.T) {
	type test struct {
		inputIdent string
		wantIdent  string
	}

	tests := []test{
		test{"SELECT", "SELECT"},
		test{"Select", "IDENT"},
		test{"CASE", "CASE"},
		test{strings.ToUpper("Join"), "JOIN"},
		test{"tableName", "IDENT"},
		test{"[another name]", "IDENT"},
		test{"INTO", "INTO"},
	}

	for _, tt := range tests {
		got := KeywordLookup(tt.inputIdent)
		if got.String() != tt.wantIdent {
			t.Errorf("Expected: [%s], got: [%s]", tt.wantIdent,
				got.String())
		}
	}
}
