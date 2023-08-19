package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

// https://spec.graphql.org/October2021/#sec-Objects
// TODO: parse description
// TODO: parse directives
func ParseTypeObjectDefinition(l *LexerWrapper) (d *ast.ObjectTypeDefinition, err error) {
	d = &ast.ObjectTypeDefinition{
		FieldDefinitions: make([]*ast.FieldDefinition, 0),
	}

	if err = l.SkipKeyword("type"); err != nil {
		return nil, err
	}

	if err = l.PeekAndMustBe([]gogqllexer.Kind{gogqllexer.Name}, func(t gogqllexer.Token, advanceLexer func()) error {
		defer advanceLexer()
		d.Name = t.Value
		return nil
	}); err != nil {
		return nil, err
	}

	// parse implements interface
	if err = l.PeekAndMayBe([]gogqllexer.Kind{gogqllexer.Name}, func(t gogqllexer.Token, advanceLexer func()) error {
		if err = l.SkipKeyword("implements"); err != nil {
			return err
		}

		interfaces := make([]string, 0)

		// implements at least one interface
		l.SkipIf(gogqllexer.Amp)
		if err = l.PeekAndMustBe(
			[]gogqllexer.Kind{gogqllexer.Name},
			func(t gogqllexer.Token, advanceLexer func()) error {
				defer advanceLexer()

				interfaces = append(interfaces, t.Value)
				return nil
			},
		); err != nil {
			return err
		}

		// read more interfaces
		for {
			if skip := l.SkipIf(gogqllexer.Amp); !skip {
				break
			}

			if err = l.PeekAndMustBe(
				[]gogqllexer.Kind{gogqllexer.Name},
				func(t gogqllexer.Token, advanceLexer func()) error {
					defer advanceLexer()

					interfaces = append(interfaces, t.Value)
					return nil
				},
			); err != nil {
				return err
			}
		}

		if len(interfaces) > 0 {
			d.Interfaces = interfaces
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if err = l.PeekAndMustBe(
		[]gogqllexer.Kind{gogqllexer.BraceL},
		func(t gogqllexer.Token, advanceLexer func()) error {
			l.NextToken()

			for {
				t = l.PeekToken()
				if t.Kind == gogqllexer.BraceR {
					l.NextToken()
					break
				}
				if t.Kind == gogqllexer.EOF {
					return fmt.Errorf("unexpected token %+v", t)
				}

				fieldDefinition, err := parseFieldDefinition(l)
				if err != nil {
					return err
				}

				d.FieldDefinitions = append(d.FieldDefinitions, fieldDefinition)
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return d, nil
}
