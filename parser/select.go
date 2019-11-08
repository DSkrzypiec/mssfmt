package parser

import (
	"mssfmt/ast"
	"mssfmt/token"
)

// Method SelectQuery parse SELECT query. This method assumes that token SELECT
// was already parsed.
func (p *Parser) SelectQuery() *ast.SelectQuery {
	selectTree := ast.SelectQuery{}

	p.selectDistinct(&selectTree)
	p.selectTop(&selectTree)
	// TODO

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
		top.Expr = ast.Expression{ast.Word{p.word.Token, p.word.Literal}}
		p.next()
	}

	if tok == token.LPAREN {
		expr := make(ast.Expression, 0, 10)

		for {
			if p.word.Token == token.RPAREN {
				expr = append(expr, ast.Word{token.RPAREN, ")"})
				break
			}

			expr = append(expr, ast.Word{p.word.Token, p.word.Literal})
			p.next()
		}
		top.Expr = expr
		p.next()
	}

	if p.word.Token == token.PERCENT {
		top.PercentParam = true
		p.next()
	}
	if p.word.Token == token.WITHTIES {
		top.WithTiesParam = true
		p.next()
	}

	(*selectTree).Top = &top
}
