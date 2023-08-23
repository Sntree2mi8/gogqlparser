package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func (p *parser) ParseArgumentsDefinition() (defs []ast.InputValueDefinition, err error) {
	if err = p.Skip(gogqllexer.ParenL); err != nil {
		return nil, err
	}

	for {
		var inputValDescription string
		inputValDescription, _ = p.ReadDescription()

		if p.CheckKind(gogqllexer.Name) {
			ivd, err := p.parseInputValueDefinition(inputValDescription)
			if err != nil {
				return nil, err
			}
			defs = append(defs, ivd)
		} else {
			return nil, fmt.Errorf("unexpected token %+v", p.PeekToken())
		}

		if p.SkipIf(gogqllexer.ParenR) {
			break
		}
	}

	if len(defs) == 0 {
		return nil, fmt.Errorf("empty arguments definition")
	}

	return defs, nil
}
