package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

// https://spec.graphql.org/October2021/#FieldsDefinition
func (p *parser) parseFieldsDefinition() (d []*ast.FieldDefinition, err error) {
	if err = p.Skip(gogqllexer.BraceL); err != nil {
		return nil, err
	}

	for {
		var fieldDefinition *ast.FieldDefinition
		fieldDefinition, err = p.parseFieldDefinition()
		if err != nil {
			return nil, err
		}

		d = append(d, fieldDefinition)

		if p.SkipIf(gogqllexer.BraceR) {
			break
		}
	}

	return d, nil
}

// https://spec.graphql.org/October2021/#FieldDefinition
func (p *parser) parseFieldDefinition() (d *ast.FieldDefinition, err error) {
	d = &ast.FieldDefinition{}

	d.Description, _ = p.ReadDescription()

	if d.Name, err = p.ReadNameValue(); err != nil {
		return nil, err
	}

	if p.CheckKind(gogqllexer.ParenL) {
		if d.ArgumentDefinition, err = p.ParseArgumentsDefinition(); err != nil {
			return nil, err
		}
	}

	if err = p.Skip(gogqllexer.Colon); err != nil {
		return nil, err
	}

	if err = p.PeekAndMustBe(
		[]gogqllexer.Kind{gogqllexer.Name},
		func(t gogqllexer.Token, advanceLexer func()) error {
			if d.Type, err = p.parseType(); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	if p.CheckKind(gogqllexer.At) {
		if d.Directives, err = p.parseDirectives(); err != nil {
			return nil, err
		}
	}

	return d, err
}
