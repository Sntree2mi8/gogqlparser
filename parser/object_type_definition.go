package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

// https://spec.graphql.org/October2021/#sec-Objects
func (p *parser) ParseObjectTypeDefinition(description string) (d *ast.ObjectTypeDefinition, err error) {
	d = &ast.ObjectTypeDefinition{
		Description: description,
	}

	if err = p.SkipKeyword("type"); err != nil {
		return nil, err
	}

	if d.Name, err = p.ReadNameValue(); err != nil {
		return nil, err
	}

	if p.CheckKeyword("implements") {
		if d.Interfaces, err = p.parseImplementsInterfaces(); err != nil {
			return nil, err
		}
	}

	if p.CheckKind(gogqllexer.At) {
		if d.Directives, err = p.parseDirectives(); err != nil {
			return nil, err
		}
	}

	if err = p.PeekAndMustBe(
		[]gogqllexer.Kind{gogqllexer.BraceL},
		func(t gogqllexer.Token, advanceLexer func()) error {
			if d.FieldDefinitions, err = p.parseFieldsDefinition(); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return d, nil
}

// ParseObjectTypeExtension parse object type extension.
// "extend" keyword must be consumed before calling this function.
//
// Reference: https://spec.graphql.org/October2021/#sec-Object-Extensions
func (p *parser) ParseObjectTypeExtension() (def *ast.ObjectTypeExtension, err error) {
	def = &ast.ObjectTypeExtension{}

	if err = p.SkipKeyword("type"); err != nil {
		return nil, err
	}

	if def.Name, err = p.ReadNameValue(); err != nil {
		return nil, err
	}

	var canOmitFields bool
	if p.CheckKeyword("implements") {
		if def.ImplementInterfaces, err = p.parseImplementsInterfaces(); err != nil {
			return nil, err
		}

		canOmitFields = true
	}

	if p.CheckKind(gogqllexer.At) {
		if def.Directives, err = p.parseDirectives(); err != nil {
			return nil, err
		}

		canOmitFields = true
	}

	if p.CheckKind(gogqllexer.BraceL) {
		if def.FieldsDefinition, err = p.parseFieldsDefinition(); err != nil {
			return nil, err
		}
	} else if !canOmitFields {
		return nil, fmt.Errorf("unexpected token. expected interface implementation or directive or fields definition")
	}

	return def, nil
}
