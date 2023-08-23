package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func (p *parser) parseInputValueDefinition(description string) (def ast.InputValueDefinition, err error) {
	def.Description = description
	if def.Name, err = p.ReadNameValue(); err != nil {
		return def, err
	}

	if err = p.Skip(gogqllexer.Colon); err != nil {
		return def, err
	}

	def.Type, err = p.parseType()
	if err != nil {
		return def, err
	}

	if p.SkipIf(gogqllexer.Equal) {
		if err = p.PeekAndMustBe(
			[]gogqllexer.Kind{gogqllexer.Int, gogqllexer.Float, gogqllexer.String, gogqllexer.BlockString, gogqllexer.Name},
			func(t gogqllexer.Token, advanceLexer func()) error {
				defer advanceLexer()

				def.RawDefaultValue = t.Value
				return nil
			},
		); err != nil {
			return def, err
		}
	}

	if p.CheckKind(gogqllexer.At) {
		if def.Directives, err = p.parseDirectives(); err != nil {
			return def, err
		}
	}

	return def, nil
}

// https://spec.graphql.org/October2021/#InputFieldsDefinition
func (p *parser) parseInputFieldsDefinition() (defs []ast.InputValueDefinition, err error) {
	if err = p.Skip(gogqllexer.BraceL); err != nil {
		return nil, err
	}

	for {
		description, _ := p.ReadDescription()
		def, err := p.parseInputValueDefinition(description)
		if err != nil {
			return nil, err
		}

		defs = append(defs, def)

		if p.SkipIf(gogqllexer.BraceR) {
			break
		}
	}

	return defs, nil
}

// https://spec.graphql.org/October2021/#sec-Input-Objects
func (p *parser) ParseInputObjectTypeDefinition(description string) (def *ast.InputObjectTypeDefinition, err error) {
	def = &ast.InputObjectTypeDefinition{
		Description: description,
	}

	if err := p.SkipKeyword("input"); err != nil {
		return nil, err
	}

	if def.Name, err = p.ReadNameValue(); err != nil {
		return nil, err
	}

	if p.CheckKind(gogqllexer.At) {
		if def.Directives, err = p.parseDirectives(); err != nil {
			return nil, err
		}
	}

	if def.InputFields, err = p.parseInputFieldsDefinition(); err != nil {
		return nil, err
	}

	return def, nil
}

// ParseInputObjectTypeExtension parse input object type extension.
// "extend" keyword must be consumed before calling this function.
//
// Reference: https://spec.graphql.org/October2021/#sec-Input-Object-Extensions
func (p *parser) ParseInputObjectTypeExtension() (def *ast.InputObjectTypeExtension, err error) {
	def = &ast.InputObjectTypeExtension{}

	if err = p.SkipKeyword("input"); err != nil {
		return nil, err
	}

	if def.Name, err = p.ReadNameValue(); err != nil {
		return nil, err
	}

	var hasDirective bool
	if p.CheckKind(gogqllexer.At) {
		if def.Directives, err = p.parseDirectives(); err != nil {
			return nil, err
		}

		hasDirective = true
	}

	// field definition can be omitted only when there is one or more directives.
	if p.CheckKind(gogqllexer.BraceL) {
		if def.InputsFieldDefinition, err = p.parseInputFieldsDefinition(); err != nil {
			return nil, err
		}
	} else if !hasDirective {
		// TODO: fix error message
		// if next token is punctuator, there is no value to print.
		return nil, fmt.Errorf("expected '{' or '@' but got %s", p.PeekToken().Value)
	}

	return def, nil
}
