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
	DISTINCT
	TOP
	PERCENT
	TIES
	AS
	FROM
	WHERE
	GROUPBY // GROUP BY
	ORDERBY // ORDER BY
	PARTITIONBY
	FORCEORDER
	JOIN
	ON
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
	UPDATE
	DELETE
	INSERT
	GO
	TRUNCATE

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
	RECOMPILE
	IN
	BETWEEN
	ANY
	ALL
	LIKE
	SOME
	AND
	OR
	NOT
	WITH
	OPTION
	keywordEnd

	operatorBeg
	// Operators in T-SQL
	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %

	ASSIGN       // =
	EQL          // =
	NEQ          // !=
	LSS          // <
	GTR          // >
	LEQ          // <=
	GEQ          // >=
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

	SELECT:      "SELECT",
	DISTINCT:    "DISTINCT",
	TOP:         "TOP",
	PERCENT:     "PERCENT",
	TIES:        "TIES",
	FROM:        "FROM",
	WHERE:       "WHERE",
	GROUPBY:     "GROUP BY",
	ORDERBY:     "ORDER BY",
	PARTITIONBY: "PARTITION BY",
	FORCEORDER:  "FORCE ORDER",
	JOIN:        "JOIN",
	ON:          "ON",
	LEFT:        "LEFT",
	RIGHT:       "RIGHT",
	FULL:        "FULL",
	INNER:       "INNER",
	CROSS:       "CROSS",
	HAVING:      "HAVING",
	INTO:        "INTO",
	CASE:        "CASE",
	WHEN:        "WHEN",
	THEN:        "THEN",
	END:         "END",
	CUBE:        "CUBE",
	ROLLUP:      "ROLLUP",
	WITH:        "WITH",
	OPTION:      "OPTION",
	UPDATE:      "UPDATE",
	DELETE:      "DELETE",
	INSERT:      "INSERT",
	GO:          "GO",
	TRUNCATE:    "TRUNCATE",

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
	RECOMPILE:             "RECOMPILE",
	AND:                   "AND",
	OR:                    "OR",
	NOT:                   "NOT",
	IN:                    "IN",
	BETWEEN:               "BETWEEN",
	ANY:                   "ANY",
	ALL:                   "ALL",
	LIKE:                  "LIKE",
	SOME:                  "SOME",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	MOD: "%",

	ASSIGN:       "=",
	EQL:          "=",
	NEQ:          "!=",
	LSS:          "<",
	GTR:          ">",
	LEQ:          "<=",
	GEQ:          ">=",
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

var keywords map[string]Token

// Map MultiwordKeywords stores multi-word keywords like "GROUP BY", "ORDER BY"
// etc. Key in this map is the first word of the keyword and value is an int
// which points how many words are contained in this keyword.
// For example in "GROUP BY": MultiwordKeywords["GROUP"] = 1 ("BY").
var MultiwordKeywords map[string]int

// Function init creates mainKeywords and keywords global package maps.
func init() {
	keywords = make(map[string]Token)
	for i := keywordBeg + 1; i < keywordEnd; i++ {
		keywords[tokens[i]] = Token(i)
	}

	MultiwordKeywords = make(map[string]int)
	MultiwordKeywords["GROUP"] = 1
	MultiwordKeywords["ORDER"] = 1
	MultiwordKeywords["PARTITION"] = 1
	MultiwordKeywords["FORCE"] = 1
	// MultiwordKeywords["WITH"] = 1 TODO: will it clash with regular WITH?
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

// IsLiteral returns true for tokens which are defined as literals.
func (t Token) IsLiteral() bool {
	return literalBeg < t && t < literalEnd
}

// IsKeyword returns true for tokens which are defined as keywords.
func (t Token) IsKeyword() bool {
	return keywordBeg < t && t < keywordEnd
}

// IsOperator returns true for tokens which are defined as operator.
func (t Token) IsOperator() bool {
	return operatorBeg < t && t < operatorEnd
}
