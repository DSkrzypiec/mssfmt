package parser

import (
	"mssfmt/ast"
	"mssfmt/scanner"
	"mssfmt/token"
	"testing"
)

// Test for parsing table or view name.
func TestParseFromTableName(t *testing.T) {
	src1 := []byte("From tableName AS tn")
	src2 := []byte("FROM tableName tn")
	src3 := []byte("FROM #tempTableName")
	var s1, s2, s3 scanner.Scanner
	var p1, p2, p3 Parser
	s1.Init("s1", src1)
	s2.Init("s2", src2)
	s3.Init("s3", src3)
	p1.Init("p1", ScanWords(s1))
	p2.Init("p2", ScanWords(s2))
	p3.Init("p3", ScanWords(s3))

	st := ast.SelectQuery{}
	var tabName ast.TableName

	// Case 1
	p1.selectFrom(&st)
	tabName = *st.From.TableOrViewName

	if tabName.Name != "tableName" {
		t.Errorf("Expected 'tableName', got: '%s'", tabName.Name)
	}
	if !tabName.ASKeyword {
		t.Errorf("Expected 'AS' after tableName, in [%s]", string(src1))
	}
	if tabName.Alias == nil {
		t.Errorf("Expected non-nil table alias - 'tn'")
	}
	if *tabName.Alias != "tn" {
		t.Errorf("Expected table alias to be 'tn', got: '%s'", *tabName.Alias)
	}

	// Case 2
	p2.selectFrom(&st)
	tabName = *st.From.TableOrViewName

	if tabName.Name != "tableName" {
		t.Errorf("Expected 'tableName', got: '%s'", tabName.Name)
	}
	if tabName.ASKeyword {
		t.Errorf("Unexpected 'AS' token in [%s]", string(src2))
	}
	if tabName.Alias == nil {
		t.Errorf("Expected non-nil table alias - 'tn'")
	}
	if *tabName.Alias != "tn" {
		t.Errorf("Expected table alias to be 'tn', got: '%s'", *tabName.Alias)
	}

	// Case 3
	//p3.selectFrom(&st)
	//tabName = *p3.From.TableOrViewName
}

// Test for select INTO clause.
func TestParseSelectInto(t *testing.T) {
	src1 := []byte("inTo tableName")
	src2 := []byte("INTo  /* comment  */  #tempTable")
	src3 := []byte("not-a-into-clause")
	var s1, s2, s3 scanner.Scanner
	var p1, p2, p3 Parser
	s1.Init("s1", src1)
	s2.Init("s2", src2)
	s3.Init("s3", src3)
	p1.Init("p1", ScanWords(s1))
	p2.Init("p2", ScanWords(s2))
	p3.Init("p3", ScanWords(s3))

	st := ast.SelectQuery{}

	// Case 1
	p1.selectInto(&st)
	into1 := st.Into

	if into1 == nil {
		t.Errorf("Expected non-nil INTO clause. Got [%v]", *into1)
	}
	if *into1 != "tableName" {
		t.Errorf("Expected [INTO tableName], got [INTO %s]", *into1)
	}

	// Case 2
	p2.selectInto(&st)
	into2 := st.Into

	if into2 == nil {
		t.Errorf("Expected non-nil INTO clause.")
	}
	if *into2 != "#tempTable" {
		t.Errorf("Expected [INTO #tempTable], got [INTO %s]", *into2)
	}

	// Case 3
	p3.selectInto(&st)
	into3 := st.Into

	if into3 != nil {
		t.Errorf("Expected nil INTO clause, got something: [%s]", *into3)
	}
}

// Test for parsing SELECT column list. Tests assumes that all tokens before the
// column list was already parsed.
func TestParseSelectColList(t *testing.T) {
	ps := prepareColListParsers()
	sts := make([]ast.SelectQuery, len(ps))

	for id, p := range ps {
		p.selectColList(&sts[id])
	}

	// Case "* from"
	cols := sts[0].Columns
	if len(cols) != 1 {
		t.Errorf("Expected length = 1, got: %d", len(cols))
	}
	if len(cols[0]) != 1 {
		t.Errorf("Expected num of token = 1, got: %d", len(cols[0]))
	}
	lit := cols[0][0].Literal
	if lit != "*" {
		t.Errorf("Expected literal '*', got: '%s'", lit)
	}

	// Case "5 AS X;"
	cols1 := sts[1].Columns
	if len(cols1) != 1 {
		t.Errorf("Expected length = 1, got: %d", len(cols1))
	}
	if len(cols1[0]) != 3 {
		t.Errorf("Expected num of token = 3, got: %d", len(cols1[0]))
	}
	if cols1[0][0].Token != token.INT || cols1[0][1].Token != token.AS ||
		cols1[0][2].Literal != "X" || cols1[0][2].Token != token.IDENT {
		t.Errorf("Expected '5 AS X', got: '%s %s %s'", cols1[0][0].Literal,
			cols1[0][1].Token, cols1[0][2].Literal)
	}

	// Case "x.Col1, y.Col2, y.Col3, x.Col1 + y.Col2 FROM ..."
	cols2 := sts[2].Columns
	if len(cols2) != 4 {
		t.Errorf("Expected 4 columns, got: %d", len(cols2))
	}
	if len(cols2[0]) != 3 || len(cols2[1]) != 3 || len(cols2[2]) != 3 ||
		len(cols2[3]) != 7 {
		t.Errorf("Expected #tokens (3, 3, 3, 7) in columns respectively, got: (%d, %d, %d, %d)",
			len(cols2[0]), len(cols2[1]), len(cols2[2]), len(cols2[3]))
	}
	if cols2[1][2].Literal != "Col2" {
		t.Errorf("In second column expected third literal as 'Col2', got '%s'",
			cols2[1][2].Literal)
	}
	if cols2[3][3].Token != token.ADD {
		t.Errorf("In forth column expected forth token as '+', got '%s'",
			cols2[3][3].Token)
	}

	// Case "1 + 5 AS X \n SELECT ..."
	cols3 := sts[3].Columns
	if len(cols3) != 1 {
		t.Errorf("Expected single column, got: %d", len(cols3))
	}
	if len(cols3[0]) != 5 {
		t.Errorf("Expected num of token = 5, got: %d (%v)", len(cols3[0]),
			cols3[0])
	}

	if cols3[0][0].Literal != "1" || cols3[0][1].Token != token.ADD ||
		cols3[0][2].Literal != "5" || cols3[0][3].Token != token.AS ||
		cols3[0][4].Literal != "X" {
		t.Errorf("Expected '1 + 5 AS X', got '%s %s %s %s %s'",
			cols3[0][0].Literal, cols3[0][1].Literal, cols3[0][2].Literal,
			cols3[0][3].Literal, cols3[0][4].Literal)
	}

	// Case "Col1, /* comment */ Col2 --comment \n , 42 AS Y Into"
	cols4 := sts[4].Columns
	if len(cols4) != 3 {
		t.Errorf("Expected 3 columns, got: %d", len(cols4))
	}
	if cols4[0][0].Literal != "Col1" {
		t.Errorf("Expected first column 'Col1', got '%s'", cols4[0][0].Literal)
	}
	if cols4[1][0].Literal != "Col2" {
		t.Errorf("Expected second column 'Col2', got '%s'", cols4[0][1].Literal)
	}
	if len(cols4[2]) != 3 {
		t.Errorf("Expected 3 tokens in third column, got: %d (%v)",
			len(cols4[2]), cols4[2])
	}
	if cols4[2][0].Literal != "42" || cols4[2][1].Token != token.AS ||
		cols4[2][2].Literal != "Y" {
		t.Errorf("Expected '42 AS Y', got '%s %s %s'", cols4[2][0].Literal,
			cols4[2][1].Token, cols4[2][2].Literal)
	}
}

// Prepares mock SQL codes in form of Parser object for testing parsing SELECT
// column list.
func prepareColListParsers() []Parser {
	const n = 5
	src := make([][]byte, n)
	ss := make([]scanner.Scanner, n)
	ps := make([]Parser, n)

	src[0] = []byte("* from")
	src[1] = []byte("5 AS X;")
	src[2] = []byte("x.Col1, y.Col2, y.Col3, x.Col1 + y.Col2 FROM ...")
	src[3] = []byte("1 + 5 AS X \n SELECT ...")
	src[4] = []byte("Col1, /* comment */ Col2 --comment \n , 42 AS Y Into")

	for i := 0; i < n; i++ {
		ss[i].Init("s", src[i])
		ps[i].Init("p", ScanWords(ss[i]))
	}

	return ps
}

// Test for parsing TOP clause in SELECT query. Tests assumes that SELECT
// keyword was already parsed.
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
