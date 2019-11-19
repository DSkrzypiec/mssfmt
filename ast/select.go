package ast

import "mssfmt/token"

// SelectQuery represents AST for SELECT query. TODO...
// From SQL Server 2019 documentation:
//
//	SELECT [ ALL | DISTINCT ]
//	    [TOP ( expression ) [PERCENT] [ WITH TIES ] ]
//	    < select_list >
//	    [ INTO new_table ]
//	    [ FROM { <table_source> } [ ,...n ] ]
//	    [ WHERE <search_condition> ]
//	    [ <GROUP BY> ]
//	    [ HAVING < search_condition > ]
//
type SelectQuery struct {
	DistinctType *DistinctType
	Top          *TopClause
	Columns      []Expression
	Into         *string
	From         *FromClause
	Where        *WhereClause
	GroupBy      *GroupByClause
	Having       *HavingClause
	Options      *SelectOptions
}

// Word represents single "word" in SQL script. It's a pair of token.Token and
// corresponding literal.
type Word struct {
	Token   token.Token
	Literal string
}

// T-SQL expression is just a slice of Words.
type Expression []Word

// DistinctType represents [ALL | DISTINCT | ] clause in SELECT query.
// Both All and Distinct can be false but they both cannot be true at the same
// time - this is invalid query.
type DistinctType struct {
	All      bool
	Distinct bool
}

// TopClaues represents SELECT TOP clause. Expr is an T-SQL expression. It might
// be any INT, FLOAT or any expression which gives INT or FLOAT.
type TopClause struct {
	Expr          Expression
	PercentParam  bool
	WithTiesParam bool
}

type FromClause struct{}
type WhereClause struct{}
type GroupByClause struct{}
type HavingClause struct{}
type SelectOptions struct{}
