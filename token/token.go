// Package token defines tokens for lexical tokens of T-SQL (query language of
// Microsoft SQL Server). This package is based on Go package "token" for
// parsing Go language.
package token

// Token is the set of lexical tokens on T-SQL.
// TODO: extend Token. For now (2019-10-13) Token is limited to basic SELECT
// queries.
type Token int

const (
	EOF = iota
	COMMENT

	literalBeg
	IDENT  // ColName, TableName, CTEName, ...
	INT    // 53421
	FLOAT  // 123.123
	STRING // 'Value'
	AS     // ColName AS cn
	literalEnd

	keywordBeg
	// T-SQL keywords
	SELECT
	TOP
	FROM
	WHERE
	GROUPBY  // GROUP BY
	ORDERBY  // ORDER BY
	LEFTJOIN // LEFT JOIN
	RIGHTJOIN
	FULLJOIN
	INNERJOIN
	CROSSJOIN
	HAVING
	INTO
	CASE
	WHEN
	THEN
	END
	keywordEnd

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
	LPAREN    // (
	LBRACK    // [
	LBRACE    // {
	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	COMMA     // ,
	PERIOD    // .
	SEMICOLON // ;
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

	SELECT:    "SELECT",
	TOP:       "TOP",
	FROM:      "FROM",
	WHERE:     "WHERE",
	GROUPBY:   "GROUP BY",
	ORDERBY:   "ORDER BY",
	LEFTJOIN:  "LEFT JOIN",
	RIGHTJOIN: "RIGHT JOIN",
	FULLJOIN:  "FULL JOIN",
	INNERJOIN: "INNER JOIN",
	CROSSJOIN: "CROSS JOIN",
	HAVING:    "HAVING",
	INTO:      "INTO",
	CASE:      "CASE",
	WHEN:      "WHEN",
	THEN:      "THEN",
	END:       "END",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	MOD: "%",
	AND: "AND",
	OR:  "OR",
	NOT: "NOT",

	ASSIGN:    "=",
	EQL:       "=",
	NEQ:       "!=",
	LSS:       "<",
	GTR:       ">",
	LEQ:       "<=",
	GEQ:       ">=",
	IN:        "IN",
	BETWEEN:   "BETWEEN",
	ANY:       "ANY",
	ALL:       "ALL",
	LIKE:      "LIKE",
	SOME:      "SOME",
	LPAREN:    "(",
	LBRACK:    "[",
	LBRACE:    "{",
	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	COMMA:     ",",
	PERIOD:    ".",
	SEMICOLON: ";",
}

// String method returns string of token t. For GROUPBY it would be "GROUP BY"
// and for "COMMA" it is ",".
func (t Token) String() string {
	if 0 <= t && t < Token(len(tokens)) {
		return tokens[t]
	}

	return "token_" + strconv.Itoa(int(t))
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
