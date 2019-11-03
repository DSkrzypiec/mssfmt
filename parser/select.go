package parser

import (
	"mssfmt/ast"
	"mssfmt/token"
	"strings"
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

//
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
		numberLit := p.word.Literal
		top.NumberLit = &numberLit
		p.next()
	}

	if tok == token.LPAREN {
		expr := make([]string, 0, 10)

		for {
			if p.word.Token == token.RPAREN {
				expr = append(expr, ")")
				break
			}

			expr = append(expr, p.word.Literal)
			p.next()
		}
		exprJ := strings.Join(expr, " ")
		top.Expression = &exprJ
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
