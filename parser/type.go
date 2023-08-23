package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func (p *parser) parseType() (t ast.Type, err error) {
	if p.SkipIf(gogqllexer.BracketL) {
		listType, err := p.parseType()
		if err != nil {
			return t, err
		}

		t.ListType = &listType
		if err = p.Skip(gogqllexer.BracketR); err != nil {
			return t, err
		}
	} else {
		if t.NamedType, err = p.ReadNameValue(); err != nil {
			return t, err
		}
	}

	if p.SkipIf(gogqllexer.Bang) {
		t.NotNull = true
	}

	return t, nil
}
