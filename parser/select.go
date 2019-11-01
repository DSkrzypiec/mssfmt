package parser

import (
	"mssfmt/ast"
	"mssfmt/token"
)

// Method SelectQuery parse SELECT query. This method assumes that token SELECT
// was already parsed.
func (p *Parser) SelectQuery() ast.SelectQuery {
	distType := p.selectDistinct()

	// TODO
	return ast.SelectQuery{distType, nil, nil, nil, nil, nil, nil, nil, nil}
}

// Method selectDistinct parses [ALL | DISTINCT | ] clause in SELECT query.
func (p *Parser) selectDistinct() ast.DistinctType {
	if p.word.Token == token.DISTINCT {
		p.next()
		return ast.DistinctType{false, true}
	}
	if p.word.Token == token.ALL {
		p.next()
		return ast.DistinctType{true, false}
	}
	return ast.DistinctType{false, false}
}
