package parser

import (
	"fmt"
	"mssfmt/ast"
	"mssfmt/scanner"
	"testing"
)

func TestParseSelectTop(t *testing.T) {
	// select is ommited, it assumes that it was aleardy scanned
	src1 := []byte("top 10 * from")
	src2 := []byte("TOP (0.1 + 0.3)   PErcent * from")
	src3 := []byte("Not-a-top-clause * from")
	var s1, s2, s3 scanner.Scanner
	var p1, p2, p3 Parser
	s1.Init("s1", src1)
	s2.Init("s2", src2)
	s3.Init("s3", src3)
	p1.Init("p1", ScanWords(s1))
	p2.Init("p2", ScanWords(s2))
	p3.Init("p3", ScanWords(s3))

	st := ast.SelectQuery{}
	p1.selectTop(&st)
	fmt.Println(st.Top)
	fmt.Printf("Num: %s | Perc: %t | Ties: %t\n",
		*(st.Top.NumberLit), st.Top.PercentParam, st.Top.WithTiesParam)

	p2.selectTop(&st)
	fmt.Printf("Expr: %s | Perc: %t | Ties: %t\n",
		*(st.Top.Expression), st.Top.PercentParam, st.Top.WithTiesParam)

	p3.selectTop(&st)
	if st.Top != nil {
		t.Errorf("Expected nil, got %v", st.Top)
	}
}

func TestParseSelectDistinct(t *testing.T) {
	src1 := []byte("select X from Y")
	src2 := []byte("select distinct X from Y")
	src3 := []byte("select ALL X from Y")
	var s1, s2, s3 scanner.Scanner
	var p1, p2, p3 Parser
	s1.Init("s1", src1)
	s2.Init("s2", src2)
	s3.Init("s3", src3)
	p1.Init("p1", ScanWords(s1))
	p2.Init("p2", ScanWords(s2))
	p3.Init("p3", ScanWords(s3))

	// Skip "select" as already parsed
	p1.next()
	p2.next()
	p3.next()
	st := ast.SelectQuery{}
	p1.selectDistinct(&st)

	if st.DistinctType != nil {
		t.Errorf("Wrong parsed case without DISTINCT and ALL")
	}

	p2.selectDistinct(&st)
	if *st.DistinctType != (ast.DistinctType{false, true}) {
		t.Errorf("Wrong parsed case with DISTINCT")
	}

	p3.selectDistinct(&st)
	if *st.DistinctType != (ast.DistinctType{true, false}) {
		t.Errorf("Wrong parsed case with ALL")
	}
}
