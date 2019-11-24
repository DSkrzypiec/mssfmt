package scanner

import (
	"mssfmt/token"
	"testing"
)

// Test for Scan method on short T-SQL script.
func TestFirstScan(t *testing.T) {
	var scanner Scanner
	scanner.Init("ShortEx", shortScriptEx())
	expTokens, expLits := shortScriptTokens()

	toks := make([]token.Token, 0, 50)
	lits := make([]string, 0, 50)

	for {
		tok, lit := scanner.Scan()
		if tok == token.EOF {
			break
		}
		toks = append(toks, tok)
		lits = append(lits, lit)
	}

	for id, tok := range expTokens {
		if tok != toks[id] {
			t.Errorf("[Token] [Id = %d] Expected <%s>, got <%s>", id, tok, toks[id])
		}

		if expLits[id] != lits[id] {
			t.Errorf("[Literal] [Id = %d] Expected <%s><%v>, got <%s><%v>",
				id, expLits[id], []byte(expLits[id]), lits[id], []byte(lits[id]))
		}
	}

}

func shortScriptEx() []byte {
	src := []byte(`--comment
;With tmp As (
	SELECT count(*) as cnt, sum(x) as sumx   from SampleTable 
	group   by z 
	order by z
)
Select Avg(sumx) from tmp`)
	return src
}

func shortScriptTokens() ([]token.Token, []string) {
	toks := make([]token.Token, 0, 50)
	lits := make([]string, 0, 50)

	toks = append(toks, token.COMMENT)
	toks = append(toks, token.SEMICOLON)
	toks = append(toks, token.WITH)
	toks = append(toks, token.IDENT)
	toks = append(toks, token.AS)
	toks = append(toks, token.LPAREN)
	toks = append(toks, token.SELECT)
	toks = append(toks, token.COUNT)
	toks = append(toks, token.LPAREN)
	toks = append(toks, token.MUL)
	toks = append(toks, token.RPAREN)
	toks = append(toks, token.AS)
	toks = append(toks, token.IDENT)
	toks = append(toks, token.COMMA)
	toks = append(toks, token.SUM)
	toks = append(toks, token.LPAREN)
	toks = append(toks, token.IDENT)
	toks = append(toks, token.RPAREN)
	toks = append(toks, token.AS)
	toks = append(toks, token.IDENT)
	toks = append(toks, token.FROM)
	toks = append(toks, token.IDENT)
	toks = append(toks, token.GROUPBY)
	toks = append(toks, token.IDENT)
	toks = append(toks, token.ORDERBY)
	toks = append(toks, token.IDENT)
	toks = append(toks, token.RPAREN)
	toks = append(toks, token.SELECT)
	toks = append(toks, token.AVG)
	toks = append(toks, token.LPAREN)
	toks = append(toks, token.IDENT)
	toks = append(toks, token.RPAREN)
	toks = append(toks, token.FROM)
	toks = append(toks, token.IDENT)

	lits = append(lits, "--comment")
	lits = append(lits, ";")
	lits = append(lits, "With")
	lits = append(lits, "tmp")
	lits = append(lits, "As")
	lits = append(lits, "(")
	lits = append(lits, "SELECT")
	lits = append(lits, "count")
	lits = append(lits, "(")
	lits = append(lits, "*")
	lits = append(lits, ")")
	lits = append(lits, "as")
	lits = append(lits, "cnt")
	lits = append(lits, ",")
	lits = append(lits, "sum")
	lits = append(lits, "(")
	lits = append(lits, "x")
	lits = append(lits, ")")
	lits = append(lits, "as")
	lits = append(lits, "sumx")
	lits = append(lits, "from")
	lits = append(lits, "SampleTable")
	lits = append(lits, "group by")
	lits = append(lits, "z")
	lits = append(lits, "order by")
	lits = append(lits, "z")
	lits = append(lits, ")")
	lits = append(lits, "Select")
	lits = append(lits, "Avg")
	lits = append(lits, "(")
	lits = append(lits, "sumx")
	lits = append(lits, ")")
	lits = append(lits, "from")
	lits = append(lits, "tmp")

	return toks, lits
}

// Test for scanning block comments in T-SQL.
func TestScanBlockComment(t *testing.T) {
	src1 := []byte("/*comment1*/")
	src2 := []byte("/**/")
	src3 := []byte("/*comment /*level1*/ level0 */")
	src4 := []byte("/* /*/*/* x\n */*/*/ */")
	var s1, s2, s3, s4 Scanner
	s1.Init("s1", src1)
	s2.Init("s2", src2)
	s3.Init("s3", src3)
	s4.Init("s4", src4)

	comments := make([]string, 4)
	comments[0] = s1.scanBlockComment()
	comments[1] = s2.scanBlockComment()
	comments[2] = s3.scanBlockComment()
	comments[3] = s4.scanBlockComment()

	exp := make([]string, 4)
	exp[0] = "/*comment1*/"
	exp[1] = "/**/"
	exp[2] = "/*comment /*level1*/ level0 */"
	exp[3] = "/* /*/*/* x\n */*/*/ */"

	for i := 0; i < 4; i++ {
		if comments[i] != exp[i] {
			t.Errorf("Expected [%s], got [%s]", exp[i], comments[i])
		}
	}
}

// Test for scanning line comment in T-SQL.
func TestScanLineComment(t *testing.T) {
	src1 := []byte("--first comment\n SELECT")
	src2 := []byte("--\n SELECT")
	src3 := []byte("--   Another -- comment\r\n SELECT")
	var s1, s2, s3 Scanner
	s1.Init("s1", src1)
	s2.Init("s2", src2)
	s3.Init("s3", src3)

	comments := make([]string, 3)
	comments[0] = s1.scanLineComment()
	comments[1] = s2.scanLineComment()
	comments[2] = s3.scanLineComment()

	exp := make([]string, 3)
	exp[0] = "--first comment"
	exp[1] = "--"
	exp[2] = "--   Another -- comment"

	for i := 0; i < 3; i++ {
		if comments[i] != exp[i] {
			t.Errorf("Expected [%s], got [%s]", exp[i], comments[i])
		}
	}
}

// Test for scanning numbers in T-SQL.
func TestScanNumber(t *testing.T) {
	src := []byte("42 512 3.1459 -2.1234 0 2.1324e-5")
	var s Scanner
	s.Init("s", src)

	toks := make([]token.Token, 0, 6)
	nums := make([]string, 0, 6)
	expToks := []token.Token{token.INT, token.INT, token.FLOAT, token.FLOAT,
		token.INT, token.FLOAT}
	expNums := []string{"42", "512", "3.1459", "-2.1234", "0", "2.1324e-5"}

	for i := 0; i < 6; i++ {
		t, n := s.scanNumber()
		toks = append(toks, t)
		nums = append(nums, n)
		s.skipWhitespace()
	}

	if len(toks) != len(expToks) {
		t.Errorf("Expected %d tokens, got %d tokens", len(expToks), len(toks))
	}
	if len(nums) != len(expNums) {
		t.Errorf("Expected %d numbers, got %d numbers", len(expNums), len(nums))
	}

	for i := 0; i < 6; i++ {
		if toks[i] != expToks[i] {
			t.Errorf("For token %d, expected %d, got %d", i, expToks[i], toks[i])
		}
		if nums[i] != expNums[i] {
			t.Errorf("For token %d, expected number %s, got %s", i, expNums[i], nums[i])
		}
	}

}

// Test for scanning SQL strings.
func TestScanSQLString(t *testing.T) {
	src := []byte("'''Test' 'stringValue'  'with '' escape'   'Cox''' ''")
	var s Scanner
	s.Init("s", src)

	string0 := s.scanSQLString()
	s.skipWhitespace()
	string1 := s.scanSQLString()
	s.skipWhitespace()
	string2 := s.scanSQLString()
	s.skipWhitespace()
	string3 := s.scanSQLString()
	s.skipWhitespace()
	string4 := s.scanSQLString()

	if string0 != "'''Test'" {
		t.Errorf("Expected '''Test', got: %s", string0)
	}
	if string1 != "'stringValue'" {
		t.Errorf("Expected 'stringValue', got: %s", string1)
	}
	if string2 != "'with '' escape'" {
		t.Errorf("Expected 'with '' escape', got: %s", string2)
	}
	if string3 != "'Cox'''" {
		t.Errorf("Expected 'Cos''', got: %s", string3)
	}
	if string4 != "''" {
		t.Errorf("Expected '', got: %s", string4)
	}
}

// Test for scanning several identifiers.
func TestIdentifierMany(t *testing.T) {
	src1 := []byte("       GOSIA    \t DamiansTable\n x.Name  #tempTable")
	src2 := []byte("#tempTable   \n @variableName")
	var s, s2 Scanner
	s.Init("s", src1)
	s2.Init("s2", src2)

	s.skipWhitespace()
	firstId := s.scanIdentifier()
	s.skipWhitespace()
	secondId := s.scanIdentifier()
	s.skipWhitespace()
	thirdId := s.scanIdentifier()

	if firstId != "GOSIA" {
		t.Errorf("Expected <GOSIA>, got: <%s>", firstId)
	}
	if secondId != "DamiansTable" {
		t.Errorf("Expected <DamiansTable>, got: <%s>", secondId)
	}
	if thirdId != "x" {
		t.Errorf("Expected <x>, got: <%s>", thirdId)
	}

	tempId := s2.scanIdentifier()
	s2.skipWhitespace()
	varId := s2.scanIdentifier()

	if hashSign != rune('#') {
		t.Errorf("Hash sign const is wrong!")
	}
	if tempId != "#tempTable" {
		t.Errorf("Expected <#tempTable>, got: <%s>", tempId)
	}
	if varId != "@variableName" {
		t.Errorf("Expected <@variableName>, got: <%s>", varId)
	}
}

// Tests for scanning delimited identifiers - [SELECT AS ColName] or
// "Illegal Var % Name!".
func TestIdentifierDelimited(t *testing.T) {
	src1 := []byte("[ illegal var __name]")
	src2 := []byte(`"Another illegal<!>"`)
	var s1, s2 Scanner
	s1.Init("s1", src1)
	s2.Init("s2", src2)

	id1 := s1.scanIdentifier()
	id2 := s2.scanIdentifier()
	const eid1 = "[ illegal var __name]"
	const eid2 = `"Another illegal<!>"`

	if id1 != eid1 {
		t.Errorf("Expected identifier <%s>, got <%s>", eid1, id1)
	}
	if id2 != eid2 {
		t.Errorf("Expected identifier <%s>, got <%s>", eid2, id2)
	}
}

// Tests for scanning regular identifiers in T-SQL.
func TestIdentifierRegular(t *testing.T) {
	src1 := []byte("VarNa#m3")
	src2 := []byte("Ciąg ")
	var s1, s2 Scanner
	s1.Init("s1", src1)
	s2.Init("s2", src2)

	id1 := s1.scanIdentifier()
	id2 := s2.scanIdentifier()
	const eid1 = "VarNa#m3"
	const eid2 = "Ciąg"

	if id1 != eid1 {
		t.Errorf("Expected identifier [%s], got [%s]", eid1, id1)
	}
	if id2 != eid2 {
		t.Errorf("Expected identifier [%s], got [%s]", eid2, id2)
	}

}

func TestSkipWhitespace(t *testing.T) {
	src1 := []byte(" X")
	src2 := []byte(" \t  \nX")
	src3 := []byte("XY")

	ss := make([]Scanner, 3)

	ss[0].Init("s1", src1)
	ss[1].Init("s2", src2)
	ss[2].Init("s3", src3)

	expOffset := []int{1, 5, 0}
	expChar := []byte("XXX")

	for i, s := range ss {
		s.skipWhitespace()

		if s.offset != expOffset[i] {
			t.Errorf("Expected offset: %d, got: %d \n", expOffset[i], s.offset)
		}
		if byte(s.char) != expChar[i] {
			t.Errorf("Expected char: %q, got: %q \n", s.char, expChar[i])
		}
	}

}

func TestPeek(t *testing.T) {
	src := []byte("SELECT")
	var s Scanner
	s.Init("s1", src)

	if string(s.peek()) != "E" {
		t.Errorf("Expected 'E', got: '%s'", string(s.peek()))
	}

	s.next()
	if string(s.peek()) != "L" {
		t.Errorf("Expected 'L', got: '%s'", string(s.peek()))
	}
}
