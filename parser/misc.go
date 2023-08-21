package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func parseType(l *LexerWrapper) (t ast.Type, err error) {
	var token gogqllexer.Token

	token = l.NextToken()
	switch token.Kind {
	case gogqllexer.Name:
		t.NamedType = token.Value
		token = l.PeekToken()
		if token.Kind == gogqllexer.Bang {
			t.NotNull = true
			l.NextToken()
		}
	case gogqllexer.BracketL:
		t.ListType = &ast.Type{}

		elmType, err := parseType(l)
		if err != nil {
			return t, err
		}
		t.ListType = &elmType

		token = l.NextToken()
		if token.Kind != gogqllexer.BracketR {
			return t, fmt.Errorf("unexpected token %+v", token)
		}

		token = l.PeekToken()
		if token.Kind == gogqllexer.Bang {
			t.NotNull = true
			l.NextToken()
		}
	default:
		return t, fmt.Errorf("unexpected token %+v", token)
	}

	return t, nil
}
