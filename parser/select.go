package parser

import (
	"fmt"
	"mssfmt/ast"
	"mssfmt/token"
)

// Method SelectQuery parse SELECT query. This method assumes that token SELECT
// was already parsed.
func (p *Parser) SelectQuery() *ast.SelectQuery {
	selectTree := ast.SelectQuery{}

	p.selectDistinct(&selectTree)
	p.selectTop(&selectTree)
	p.selectColList(&selectTree)
	p.selectInto(&selectTree)
	p.selectFrom(&selectTree)
	// ...

	return &selectTree
}

// Method selectDistinct parses [ALL | DISTINCT | ] clause in SELECT query.
func (p *Parser) selectDistinct(selectTree *ast.SelectQuery) {
	if p.word.Token != token.DISTINCT && p.word.Token != token.ALL {
		return
	}

	if p.word.Token == token.DISTINCT {
		p.next()
		(*selectTree).DistinctType = &ast.DistinctType{false, true}
		return
	}
	if p.word.Token == token.ALL {
		p.next()
		(*selectTree).DistinctType = &ast.DistinctType{true, false}
		return
	}
	(*selectTree).DistinctType = nil
}

// Method selectTop parses TOP clause in SELECT query. TOP clause is of the form
// "TOP (expression) [PERCENT] [WITH TIES]". This method uses Parser to update
// ast.SelectQuery object. In case when TOP clause doesn't occur in SELECT query
// then nil is set and Parser doesn't move next().
func (p *Parser) selectTop(selectTree *ast.SelectQuery) {
	if p.word.Token != token.TOP {
		(*selectTree).Top = nil
		return
	}

	// p.word.Token == token.TOP
	p.next()
	top := ast.TopClause{}
	tok := p.word.Token

	if (tok == token.INT || tok == token.FLOAT) && tok != token.LPAREN {
		top.Expr = ast.Expression{p.word}
		p.next()
	}

	if tok == token.LPAREN {
		expr := make(ast.Expression, 0, 10)

		for {
			if p.word.Token == token.RPAREN {
				expr = append(expr, p.word)
				break
			}

			expr = append(expr, p.word)
			p.next()
		}
		top.Expr = expr
		p.next()
	}

	if p.word.Token == token.PERCENT {
		top.PercentParam = true
		p.next()
	}
	if p.word.Token == token.WITH && p.peek().Token == token.TIES {
		top.WithTiesParam = true
		p.next()
	}

	(*selectTree).Top = &top
}

// Method selectColList parses list of "columns" in SELECT query. At this point
// single "column" is treated as T-SQL expression. It's usually comma-separated
// list of column names but it can be any valid T-SQL expression which uses
// functions, operators, etc. Formal definition what can be placed in SELECT
// column list exists but mssfmt assumes that given T-SQL code is a valid code.
// Therefore validation of all those details isn't necessary.
func (p *Parser) selectColList(selectTree *ast.SelectQuery) {
	stopTokens := colListStopTokens()
	cols := make([]ast.Expression, 0, 10)
	currCol := make(ast.Expression, 0, 3)

	for {
		if _, stop := stopTokens[p.word.Token]; stop {
			cols = append(cols, currCol)
			break
		}
		if p.word.Token == token.COMMA {
			cols = append(cols, currCol)
			currCol = make(ast.Expression, 0, 3)
			p.next()
			continue
		}
		currCol = append(currCol, p.word)
		p.next()
	}

	(*selectTree).Columns = cols
}

// The following function returns dict of Tokens which should state that this is
// end of list of columns in SELECT query.
func colListStopTokens() map[token.Token]bool {
	colStop := make(map[token.Token]bool)
	colStop[token.FROM] = true
	colStop[token.INTO] = true
	colStop[token.SEMICOLON] = true
	colStop[token.SELECT] = true
	colStop[token.UPDATE] = true
	colStop[token.INSERT] = true
	colStop[token.TRUNCATE] = true
	colStop[token.GO] = true
	//colStop[token.LPAREN] = true

	return colStop
}

// Method selectInto parses INTO expression in SELECT query.
func (p *Parser) selectInto(selectTree *ast.SelectQuery) {
	if p.word.Token != token.INTO {
		(*selectTree).Into = nil
		return
	}

	p.next()
	if p.word.Token != token.IDENT {
		(*selectTree).Into = nil
		return
	}

	intoLit := p.word.Literal
	(*selectTree).Into = &intoLit
}

// Method for parsing FROM clause in SELECT query.
func (p *Parser) selectFrom(selectTree *ast.SelectQuery) {
	if p.word.Token != token.FROM {
		(*selectTree).From = nil
		return
	}

	p.next()

	// parsing subquery
	if p.word.Token == token.LPAREN {
		fmt.Println("TODO(parser.selectFrom): Parsing subquery not implemented")
		return
	}

	// parsing table name
	if p.word.Token == token.IDENT {
		p.tableName(selectTree)
		p.next()
	}

	// parsing all JOIN expressions
	for {
		if !p.word.Token.IsJoinType() {
			break
			// TODO
		}
		p.joinClause(selectTree)
		p.next()
	}
}

// Method tableName parses table or view name just after FROM keyword in SELECT
// query in case when token after FROM is an identifier.
func (p *Parser) tableName(selectTree *ast.SelectQuery) {
	tabName := ast.TableName{}
	tabName.Name = p.word.Literal

	p.next()
	tabName.ASKeyword = p.word.Token == token.AS
	if tabName.ASKeyword {
		p.next()
	}
	if p.word.Token == token.IDENT {
		alias := p.word.Literal
		tabName.Alias = &alias
		p.next()
	}

	if p.word.Token == token.TABLESAMPLE {
		// TODO: make routine for parsing table sample clause
		// just after implementing ast.FromClause.Joins parsing
	}

	if p.word.Token == token.WITH {
		// TODO: make routine for parsing table hints
		// just after implementing ast.FromClause.Joins parsing
	}

	from := ast.FromClause{&tabName, nil}
	(*selectTree).From = &from
}

// Method joinClause parses single JOIN expression in T-SQL query. This method
// is meant to be called until all of JOIN expressions from the SELECT are
// parsed. When method is called current token supposed to be SQL JOIN type
// keyword.
func (p *Parser) joinClause(selectTree *ast.SelectQuery) {
	// TODO
}
