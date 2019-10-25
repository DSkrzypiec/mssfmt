// Package token defines tokens for lexical tokens of T-SQL (query language of
// Microsoft SQL Server). This package is based on Go package "token" for
// parsing Go language.
package token

import "strconv"

// Token is the set of lexical tokens on T-SQL.
// TODO: extend Token. For now (2019-10-13) Token is limited to basic SELECT
// queries.
type Token int

const (
	EOF = iota
	ILLEGAL
	COMMENT

	literalBeg
	IDENT  // ColName, TableName, CTEName, ...
	INT    // 53421
	FLOAT  // 123.123, 4321.123e-3
	STRING // 'Value'
	literalEnd

	keywordBeg
	// T-SQL keywords
	SELECT
	TOP
	AS
	FROM
	WHERE
	GROUPBY // GROUP BY
	ORDERBY // ORDER BY
	JOIN    //
	LEFT
	RIGHT
	FULL
	INNER
	CROSS
	HAVING
	INTO
	CASE
	WHEN
	THEN
	END
	CUBE
	ROLLUP
	keywordEnd

	aggFuncsBeg
	APPROX_COUNT_DISTINCT
	AVG
	CHECKSUM_AGG
	COUNT
	COUNT_BIG
	GROUPING
	GROUPING_ID
	MAX
	MIN
	STDEV
	STDEVP
	STRING_AGG
	SUM
	VAR
	VARP
	aggFuncsEnd

	operatorBeg
	// Operators in T-SQL
	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %
	AND
	OR
	NOT

	ASSIGN // =
	EQL    // =
	NEQ    // !=
	LSS    // <
	GTR    // >
	LEQ    // <=
	GEQ    // >=
	IN
	BETWEEN
	ANY
	ALL
	LIKE
	SOME
	LPAREN       // (
	LBRACK       // [
	LBRACE       // {
	RPAREN       // )
	RBRACK       // ]
	RBRACE       // }
	COMMA        // ,
	PERIOD       // .
	SEMICOLON    // ;
	SINGLEQUOTE  // '
	DOUBLEQUOTES // "
	operatorEnd
)

var tokens = [...]string{
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",
	AS:     "AS",

	SELECT:  "SELECT",
	TOP:     "TOP",
	FROM:    "FROM",
	WHERE:   "WHERE",
	GROUPBY: "GROUP BY",
	ORDERBY: "ORDER BY",
	JOIN:    "JOIN",
	LEFT:    "LEFT",
	RIGHT:   "RIGHT",
	FULL:    "FULL",
	INNER:   "INNER",
	CROSS:   "CROSS",
	HAVING:  "HAVING",
	INTO:    "INTO",
	CASE:    "CASE",
	WHEN:    "WHEN",
	THEN:    "THEN",
	END:     "END",
	CUBE:    "CUBE",
	ROLLUP:  "ROLLUP",

	APPROX_COUNT_DISTINCT: "APPROX_COUNT_DISTINCT",
	AVG:                   "AVG",
	CHECKSUM_AGG:          "CHECKSUM_AGG",
	COUNT:                 "COUNT",
	COUNT_BIG:             "COUNT_BIG",
	GROUPING:              "GROUPING",
	GROUPING_ID:           "GROUPING_ID",
	MAX:                   "MAX",
	MIN:                   "MIN",
	STDEV:                 "STDEV",
	STDEVP:                "STDEVP",
	STRING_AGG:            "STRING_AGG",
	SUM:                   "SUM",
	VAR:                   "VAR",
	VARP:                  "VARP",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	MOD: "%",
	AND: "AND",
	OR:  "OR",
	NOT: "NOT",

	ASSIGN:       "=",
	EQL:          "=",
	NEQ:          "!=",
	LSS:          "<",
	GTR:          ">",
	LEQ:          "<=",
	GEQ:          ">=",
	IN:           "IN",
	BETWEEN:      "BETWEEN",
	ANY:          "ANY",
	ALL:          "ALL",
	LIKE:         "LIKE",
	SOME:         "SOME",
	LPAREN:       "(",
	LBRACK:       "[",
	LBRACE:       "{",
	RPAREN:       ")",
	RBRACK:       "]",
	RBRACE:       "}",
	COMMA:        ",",
	PERIOD:       ".",
	SEMICOLON:    ";",
	SINGLEQUOTE:  "'",
	DOUBLEQUOTES: `"`,
}

// String method returns string of token t. For GROUPBY it would be "GROUP BY"
// and for "COMMA" it is ",".
func (t Token) String() string {
	if 0 <= t && t < Token(len(tokens)) {
		return tokens[t]
	}

	return "token_" + strconv.Itoa(int(t))
}

var aggFuncNames map[string]Token
var keywords map[string]Token

// Function init creates mainKeywords and keywords global package maps.
func init() {
	keywords = make(map[string]Token)
	for i := keywordBeg + 1; i < keywordEnd; i++ {
		keywords[tokens[i]] = Token(i)
	}

	aggFuncNames = make(map[string]Token)
	for i := aggFuncsBeg + 1; i < aggFuncsEnd; i++ {
		aggFuncNames[tokens[i]] = Token(i)
	}
}

// Ranges for precedence ranges for operators in T-SQL.
const (
	LowestPrec  = 0
	HighestOrec = 9
)

// Precedence method returns value of precedence of T-SQL operators, from lowest
// to highest precedence. If given token oper isn't operator than LowestPrec is
// returned.
func (oper Token) Precedence() int {
	switch oper {
	case ASSIGN:
		return 1
	case ALL, ANY, BETWEEN, IN, LIKE, OR, SOME:
		return 2
	case AND:
		return 3
	case NOT:
		return 4
	case EQL, LSS, GTR, LEQ, GEQ, NEQ:
		return 5
	case ADD, SUB:
		return 6
	case MUL, DIV, MOD:
		return 7
	}
	return LowestPrec
}

// KeywordLookup search for keyword Token based on given identifier.
func KeywordLookup(identifier string) Token {
	tok, isKeyword := keywords[identifier]
	if isKeyword {
		return tok
	}
	return IDENT
}

// AggFuncLookup search for aggregation function name Token based on given
// identifier.
func AggFuncLookup(identifier string) Token {
	tok, isAggFuncName := aggFuncNames[identifier]
	if isAggFuncName {
		return tok
	}
	return IDENT
}

// IsLiteral returns true for tokens which are defined as literals.
func (t Token) IsLiteral() bool {
	return literalBeg < t && t < literalEnd
}

// IsKeyword returns true for tokens which are defined as keywords.
func (t Token) IsKeyword() bool {
	return keywordBeg < t && t < keywordEnd
}

// IsAggFunc returns true for tokens which are defined as aggregation functions.
func (t Token) IsAggFunc() bool {
	return aggFuncsBeg < t && t < aggFuncsEnd
}

// IsOperator returns true for tokens which are defined as operator.
func (t Token) IsOperator() bool {
	return operatorBeg < t && t < operatorEnd
}
