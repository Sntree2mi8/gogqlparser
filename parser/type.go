package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func parseType(l *LexerWrapper) (t ast.Type, err error) {
	if l.SkipIf(gogqllexer.BracketL) {
		listType, err := parseType(l)
		if err != nil {
			return t, err
		}

		t.ListType = &listType
		if err = l.Skip(gogqllexer.BracketR); err != nil {
			return t, err
		}
	} else {
		if t.NamedType, err = l.ReadNameValue(); err != nil {
			return t, err
		}
	}

	if l.SkipIf(gogqllexer.Bang) {
		t.NotNull = true
	}

	return t, nil
}
