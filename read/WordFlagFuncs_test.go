package read

import (
    "fmt"
	"testing"
)

func TestCommentRemarkSimple(t *testing.T) {
	sql := `--some inline comment select
	select top 10 * from tableName
	
	select top /*some comment*/ 10 * from tableName
	/*
		comment
				*/
	`
	rawS := RawScript{"x.sql", "./x.sql", sql}
	s := rawS.ToScript()

	commIds := []int{0, 1, 2, 3, 4, 5, 6, 28, 29, 30, 41, 42, 43, 44, 45}
	notComm := []int{8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}

	for _, wId := range commIds {
		if !(*s.Flags)[wId].IsComment {
			t.Errorf("Expected comment for {%s|id=%d}", s.Words[wId], wId)
		}
	}

	for _, wId := range notComm {
		if (*s.Flags)[wId].IsComment {
			t.Errorf("Expected not a comment for {%s|id=%d}", s.Words[wId], wId)
		}
	}
}

func TestIsMainKeyword(t *testing.T) {
	sql1 := "... Select     toP \t \t \t 1000 someColumnName"
	sql2 := "GROUP \n         \t    \t \n          bY"
	sql3 := "froM"

	rawS1 := RawScript{"x.sql", "./x.sql", sql1}
	rawS2 := RawScript{"x.sql", "./x.sql", sql2}
	rawS3 := RawScript{"x.sql", "./x.sql", sql3}
	s1 := rawS1.ToScript()
	s2 := rawS2.ToScript()
	s3 := rawS3.ToScript()

	for wId, w := range s1.Words {
		if (wId == 0 || wId == 18) &&
			(*s1.Flags)[wId].IsMainKeyword {
			t.Errorf("[sql1] %s - Should be not marked as MainKeyword. \n", w)
		}
		if wId > 0 && wId < 18 && !(*s1.Flags)[wId].IsMainKeyword {
			t.Errorf("[sql1] %s - Should be marked as MainKeyword. \n", w)
		}
	}

	// every word supposed to be MainKeyword
	for wId, w := range s2.Words {
		if !(*s2.Flags)[wId].IsMainKeyword {
			t.Errorf("[sql2] %s - Should be marked as MainKeyword \n", w)
		}
	}

	if !(*s3.Flags)[0].IsMainKeyword {
		t.Errorf("[sql3] Should be marked as MainKeyword: %s \n", s3.Words[0])
	}

	for wId, w := range s1.Words {
		fmt.Printf("%q - IsKey = %t \n", w, (*s1.Flags)[wId].IsMainKeyword)
	}
}

// Test for MarkCharIds method.
func TestCharIds(t *testing.T) {
	sql := `select top 10
*
from tableName t
where t.X = 1`

	rawS := RawScript{"x.sql", "./x.sql", sql}
	s := rawS.ToScript()
	flags := *s.Flags

	if flags[0].CharIdStart != 0 || flags[0].CharIdEnd != 5 {
		t.Errorf("Expected [0, 5]. Got: [%d, %d]. \n", flags[0].CharIdStart,
			flags[0].CharIdEnd)
	}
	if flags[3].CharIdStart != 10 || flags[3].CharIdEnd != 10 {
		t.Errorf("Expected [10, 10]. Got: [%d, %d]. \n", flags[3].CharIdStart,
			flags[3].CharIdEnd)
	}

	if flags[10].CharIdStart != 21 || flags[10].CharIdEnd != 29 {
		t.Errorf("Expected [21, 29]. Got: [%d, %d]. \n", flags[10].CharIdStart,
			flags[10].CharIdEnd)
	}
}

func TestAllFlags(t *testing.T) {
	sql := `-- some comment in select T-SQL script
SELECT
	t.X,
	COUNT(*)  	AS Cnt
From
	tableName t
where
	t.SomeTableName = 1
GROUP BY
	t.X


;WITH cte AS (
	select top 10 * from tableName tn
)
select top 3 * from cte
`

	rawS := RawScript{"x.sql", "./x.sql", sql}
	s := rawS.ToScript()
	fmt.Println(sql)

	for wId, w := range s.Words {
		f := (*s.Flags)[wId]
		fmt.Printf("[%d] %q {Comm=%t, MainK=%t, #Line=%d, LineInd=%d, Char={%d, %d}}\n",
			wId, w, f.IsComment, f.IsMainKeyword, f.LineNumber, f.LineIndentLvl,
			f.CharIdStart, f.CharIdEnd)
	}
}
