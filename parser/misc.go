package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func ParseArgumentsDefinition(l *LexerWrapper) (defs []ast.InputValueDefinition, err error) {
	t := l.NextToken()
	if t.Kind != gogqllexer.ParenL {
		return defs, err
	}

	for {
		inputValDescription := ""

		t = l.PeekToken()
		switch t.Kind {
		case gogqllexer.String, gogqllexer.BlockString:
			inputValDescription = t.Value
			l.NextToken()
		case gogqllexer.Name:
			ivd, err := parseInputValueDefinition(l, inputValDescription)
			if err != nil {
				return defs, err
			}
			defs = append(defs, ivd)
		case gogqllexer.ParenR:
			if len(defs) == 0 {
				return defs, fmt.Errorf("unexpected token %+v", t)
			}
			l.NextToken()
			return defs, nil
		default:
			return defs, fmt.Errorf("unexpected token %+v", t)
		}
	}
}

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
