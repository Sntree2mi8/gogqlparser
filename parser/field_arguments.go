package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func ParseArgumentsDefinition(l *LexerWrapper) (defs []ast.InputValueDefinition, err error) {
	if err = l.Skip(gogqllexer.ParenL); err != nil {
		return nil, err
	}

	for {
		var inputValDescription string
		inputValDescription, _ = l.ReadDescription()

		if l.CheckKind(gogqllexer.Name) {
			ivd, err := parseInputValueDefinition(l, inputValDescription)
			if err != nil {
				return nil, err
			}
			defs = append(defs, ivd)
		} else {
			return nil, fmt.Errorf("unexpected token %+v", l.PeekToken())
		}

		if l.SkipIf(gogqllexer.ParenR) {
			break
		}
	}

	if len(defs) == 0 {
		return nil, fmt.Errorf("empty arguments definition")
	}

	return defs, nil
}
