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
		return
	}
}
