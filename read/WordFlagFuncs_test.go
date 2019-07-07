package read

import (
	"testing"
    "reflect"
)

func TestCommentRemarkSimple(t *testing.T) {
	rawS := RawScript{"x.sql", "./x.sql", sqlComment}
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
	sql1 := "SELECT TOP 10 x.Name from dbo.Table x"
	sql2 := "  group \n         \t    \t \n          bY varX"
	sql3 := "SELECT y, count(*) from TableName Where x = 1 grouP by y"

	rawS1 := RawScript{"x.sql", "./x.sql", sql1}
	rawS2 := RawScript{"x.sql", "./x.sql", sql2}
    rawS3 := RawScript{"x.sql", "./x.sql", sql3}
	s1 := rawS1.ToScript()
	s2 := rawS2.ToScript()
	s3 := rawS3.ToScript()

	for wId, w := range s1.Words {
        f := (*s1.Flags)[wId]
        if (wId != 6 && wId < 10) && !f.IsMainKeyword {
            t.Errorf("Supposed to be a keyword: %s\n", w)
        }
	}

    for i := 10; i < 13; i++ {
        if (*s1.Flags)[i].IsMainKeyword {
            t.Errorf("%s shouldn't be a keyword.\n", s1.Words[6])
        }
    }

    for i := 0; i < len(s2.Words) - 1; i++ {
        f := (*s2.Flags)[i]
        if !f.IsMainKeyword {
            t.Errorf("Supposed to be a keyword: %s\n", s2.Words[i])
        }
    }

    if (*s2.Flags)[len(s2.Words) - 1].IsMainKeyword {
        t.Errorf("%s shouldn't be a keyword.\n", s2.Words[len(s2.Words)-1])
    }

    expected3Key := []int{0, 1, 9, 10, 11, 13, 14, 15, 21, 22, 23, 24, 25}
    real3Key := make([]int, 0, 25)

    for wId, _ := range s3.Words {
        if (*s3.Flags)[wId].IsMainKeyword {
            real3Key = append(real3Key, wId)
        }
    }

    if len(expected3Key) != len(real3Key) {
        t.Errorf("Expected: %d keywords, got: %d \n", len(expected3Key), 
            len(real3Key))
    }

    if !reflect.DeepEqual(expected3Key, real3Key) {
        t.Errorf("Expected keyword ids: [%v], got: [%v]\n", expected3Key, 
            real3Key)
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

const sqlComment = `--some inline comment select
	select top 10 * from tableName
	
	select top /*some comment*/ 10 * from tableName
	/*
		comment
				*/
	`

const sqlFull = `-- some comment in select T-SQL script
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
