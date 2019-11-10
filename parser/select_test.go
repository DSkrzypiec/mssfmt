package parser

import (
	"mssfmt/ast"
	"mssfmt/scanner"
	"mssfmt/token"
	"testing"
)

func TestParseSelectTop(t *testing.T) {
	// select is ommited, it assumes that it was aleardy scanned
	src1 := []byte("top 42 * from")
	src2 := []byte("TOP (0.1 + 0.3)   PErcent * from")
	src3 := []byte("Not-a-top-clause * from")
	src4 := []byte("Top 55 With TIES")
	var s1, s2, s3, s4 scanner.Scanner
	var p1, p2, p3, p4 Parser
	s1.Init("s1", src1)
	s2.Init("s2", src2)
	s3.Init("s3", src3)
	s4.Init("s4", src4)
	p1.Init("p1", ScanWords(s1))
	p2.Init("p2", ScanWords(s2))
	p3.Init("p3", ScanWords(s3))
	p4.Init("p4", ScanWords(s4))

	st := ast.SelectQuery{}

	// Case 1
	p1.selectTop(&st)
	num := st.Top.Expr[0]

	if num.Token != token.INT {
		t.Errorf("Expected token INT, got: %s", num.Token)
	}
	if num.Literal != "42" {
		t.Errorf("Expected literal [42], got: [%s]", num.Literal)
	}
	if st.Top.PercentParam || st.Top.WithTiesParam {
		t.Errorf("Doesn't expect [PERCENT] or [WITH TIES] tokens")
	}

	// Case 2
	p2.selectTop(&st)
	t1 := st.Top.Expr[1]
	t4 := st.Top.Expr[4]

	if t1.Token != token.FLOAT {
		t.Errorf("Expected token FLOAT, got: %s", t1.Token)
	}
	if t4.Token != token.RPAREN {
		t.Errorf("Expected token RPAREN ')', got: [%s]", t4.Token)
	}
	if t1.Literal != "0.1" {
		t.Errorf("Expected literal [0.1], got: [%s]", t1.Literal)
	}
	if !st.Top.PercentParam {
		t.Errorf("Expected [PERCENT] token but not found.")
	}

	// Case 3
	p3.selectTop(&st)
	if st.Top != nil {
		t.Errorf("Expected nil, got %v", st.Top)
	}

	// Case 4
	p4.selectTop(&st)
	if !st.Top.WithTiesParam {
		t.Errorf("Expected [WITH TIES] but not parsed.")
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
