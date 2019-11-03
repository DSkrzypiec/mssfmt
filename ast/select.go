package ast

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
	Columns      []Column
	Into         *IntoClause
	From         *FromClause
	Where        *WhereClause
	GroupBy      *GroupByClause
	Having       *HavingClause
	Options      *SelectOptions
}

// DistinctType represents [ALL | DISTINCT | ] clause in SELECT query.
// Both All and Distinct can be false but they both cannot be true at the same
// time - this is invalid query.
type DistinctType struct {
	All      bool
	Distinct bool
}

// TopClaues represents SELECT TOP clause. // TODO
type TopClause struct {
	NumberLit     *string
	Expression    *string
	PercentParam  bool
	WithTiesParam bool
}

type Column struct{}
type IntoClause struct{}
type FromClause struct{}
type WhereClause struct{}
type GroupByClause struct{}
type HavingClause struct{}
type SelectOptions struct{}
