package parser

import (
	"mssfmt/ast"
	"mssfmt/scanner"
	"testing"
)

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

	if p1.selectDistinct() != (ast.DistinctType{false, false}) {
		t.Errorf("Wrong parsed case without DISTINCT and ALL")
	}
	if p2.selectDistinct() != (ast.DistinctType{false, true}) {
		t.Errorf("Wrong parsed case with DISTINCT")
	}
	if p3.selectDistinct() != (ast.DistinctType{true, false}) {
		t.Errorf("Wrong parsed case with ALL")
	}
}
