package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

// https://spec.graphql.org/October2021/#FieldDefinition
func parseFieldDefinition(l *LexerWrapper) (d *ast.FieldDefinition, err error) {
	d = &ast.FieldDefinition{}

	if err = l.PeekAndMayBe(
		[]gogqllexer.Kind{gogqllexer.String, gogqllexer.BlockString},
		func(t gogqllexer.Token, advanceLexer func()) error {
			defer advanceLexer()

			d.Description = t.Value
			return nil
		},
	); err != nil {
		return nil, err
	}

	if err = l.PeekAndMustBe(
		[]gogqllexer.Kind{gogqllexer.Name},
		func(t gogqllexer.Token, advanceLexer func()) error {
			defer advanceLexer()

			d.Name = t.Value
			return nil
		},
	); err != nil {
		return nil, err
	}

	if err = l.PeekAndMayBe(
		[]gogqllexer.Kind{gogqllexer.ParenL},
		func(t gogqllexer.Token, advanceLexer func()) error {
			if d.ArgumentDefinition, err = ParseArgumentsDefinition(l); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	if err = l.Skip(gogqllexer.Colon); err != nil {
		return nil, err
	}

	if err = l.PeekAndMustBe(
		[]gogqllexer.Kind{gogqllexer.Name},
		func(t gogqllexer.Token, advanceLexer func()) error {
			if d.Type, err = parseType(l); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	if err = l.PeekAndMayBe(
		[]gogqllexer.Kind{gogqllexer.At},
		func(t gogqllexer.Token, advanceLexer func()) error {
			for {
				t := l.PeekToken()
				if t.Kind == gogqllexer.At {
					directive, err := parseDirective(l)
					if err != nil {
						return err
					}
					d.Directives = append(d.Directives, directive)
					continue
				}
				break
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return d, err
}
