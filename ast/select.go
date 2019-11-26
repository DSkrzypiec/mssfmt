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

// -- Syntax for SQL Server and Azure SQL Database
//
// [ FROM { <table_source> } [ ,...n ] ]
// <table_source> ::=
// {
//     table_or_view_name [ [ AS ] table_alias ]
//         [ <tablesample_clause> ]
//         [ WITH ( < table_hint > [ [ , ]...n ] ) ]
//     | rowset_function [ [ AS ] table_alias ]
//         [ ( bulk_column_alias [ ,...n ] ) ]
//     | user_defined_function [ [ AS ] table_alias ]
//     | OPENXML <openxml_clause>
//     | derived_table [ [ AS ] table_alias ] [ ( column_alias [ ,...n ] ) ]
//     | <joined_table>
//     | <pivoted_table>
//     | <unpivoted_table>
//     | @variable [ [ AS ] table_alias ]
//     | @variable.function_call ( expression [ ,...n ] )
//         [ [ AS ] table_alias ] [ (column_alias [ ,...n ] ) ]
//     | FOR SYSTEM_TIME <system_time>
// }
// <tablesample_clause> ::=
//     TABLESAMPLE [SYSTEM] ( sample_number [ PERCENT | ROWS ] )
//         [ REPEATABLE ( repeat_seed ) ]
//
// <joined_table> ::=
// {
//     <table_source> <join_type> <table_source> ON <search_condition>
//     | <table_source> CROSS JOIN <table_source>
//     | left_table_source { CROSS | OUTER } APPLY right_table_source
//     | [ ( ] <joined_table> [ ) ]
// }
// <join_type> ::=
//     [ { INNER | { { LEFT | RIGHT | FULL } [ OUTER ] } } [ <join_hint> ] ]
//     JOIN
//
// <pivoted_table> ::=
//     table_source PIVOT <pivot_clause> [ [ AS ] table_alias ]
//
// <pivot_clause> ::=
//         ( aggregate_function ( value_column [ [ , ]...n ])
//         FOR pivot_column
//         IN ( <column_list> )
//     )
//
// <unpivoted_table> ::=
//     table_source UNPIVOT <unpivot_clause> [ [ AS ] table_alias ]
//
// <unpivot_clause> ::=
//     ( value_column FOR pivot_column IN ( <column_list> ) )
//
// <column_list> ::=
//     column_name [ ,...n ]
//
// <system_time> ::=
// {
//        AS OF <date_time>
//     |  FROM <start_date_time> TO <end_date_time>
//     |  BETWEEN <start_date_time> AND <end_date_time>
//     |  CONTAINED IN (<start_date_time> , <end_date_time>)
//     |  ALL
// }
//
//     <date_time>::=
//         <date_time_literal> | @date_time_variable
//
//     <start_date_time>::=
//         <date_time_literal> | @date_time_variable
//
//     <end_date_time>::=
//         <date_time_literal> | @date_time_variable
type FromClause struct {
	// At this point (2019-11-26) FromClause handels only simple FROM caluse
	// like tableName and joins
	TableOrViewName *TableName
	Joins           *[]SQLJoin
}

// TableName represents single table or view name with some properties like
// alias, sample clause and table hints.
type TableName struct {
	Name      string
	ASKeyword bool
	Alias     *string
	Sample    *tableSampleClause
	Hints     *tableHints
}

type tableSampleClause struct {
	// TODO
}

type tableHints struct {
	// TODO
}

// SQLJoin represents T-SQL JOIN expressions.
type SQLJoin struct {
	LeftTableName  string
	Type           SQLJoinType
	Hints          SQLJoinHints
	RightTableName string
	Condition      Expression
}

// SQLJoinType is an enum for JOIN types in T-SQL.
type SQLJoinType int

const (
	INNER SQLJoinType = iota
	LEFT
	RIGHT
	FULL
	CROSS
	LEFTOUTER
	RIGHTOUTER
	FULLOUTER
)

// SQLJoinHints contains
type SQLJoinHints struct{}

type WhereClause struct{}
type GroupByClause struct{}
type HavingClause struct{}
type SelectOptions struct{}
