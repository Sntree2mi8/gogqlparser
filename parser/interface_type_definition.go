package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

// https://spec.graphql.org/October2021/#sec-Interfaces
func (p *parser) ParseInterfaceTypeDefinition(description string) (d *ast.InterfaceTypeDefinition, err error) {
	d = &ast.InterfaceTypeDefinition{}

	d.Description = description

	if err = p.SkipKeyword("interface"); err != nil {
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

	return d, err
}

// ParseInterfaceTypeExtension parse interface type extension.
// "extend" keyword must be consumed before calling this function.
//
// Reference: https://spec.graphqp.org/October2021/#sec-Interface-Extensions
func (p *parser) ParseInterfaceTypeExtension() (def *ast.InterfaceTypeExtension, err error) {
	def = &ast.InterfaceTypeExtension{}

	if err = p.SkipKeyword("interface"); err != nil {
		return nil, err
	}

	if def.Name, err = p.ReadNameValue(); err != nil {
		return nil, err
	}

	if p.CheckKeyword("implements") {
		if def.ImplementInterfaces, err = p.parseImplementsInterfaces(); err != nil {
			return nil, err
		}
	}

	if p.CheckKind(gogqllexer.At) {
		if def.Directives, err = p.parseDirectives(); err != nil {
			return nil, err
		}
	}

	if p.CheckKind(gogqllexer.BraceL) {
		if def.FieldsDefinition, err = p.parseFieldsDefinition(); err != nil {
			return nil, err
		}
	}

	return def, nil
}
