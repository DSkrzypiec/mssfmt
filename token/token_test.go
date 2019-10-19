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
		test{strings.ToUpper("inner Join"), "INNER JOIN"},
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

func TestAggFuncNames(t *testing.T) {
	type test struct {
		inputIdent string
		wantIdent  string
	}

	tests := []test{
		test{"SUM", "SUM"},
		test{"MyFunctionName", "IDENT"},
		test{"AVG", "AVG"},
		test{strings.ToUpper("stdev"), "STDEV"},
		test{"STRING_AGG", "STRING_AGG"},
		test{"AnotherFuncName", "IDENT"},
	}

	for _, tt := range tests {
		got := AggFuncLookup(tt.inputIdent)
		if got.String() != tt.wantIdent {
			t.Errorf("Expected: [%s], got: [%s]", tt.wantIdent,
				got.String())
		}
	}
}
