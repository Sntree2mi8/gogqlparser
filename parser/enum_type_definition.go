package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

// https://spec.graphql.org/October2021/#EnumValuesDefinition
func (p *parser) parseEnumValuesDefinition() (enumValuesDef []ast.EnumValueDefinition, err error) {
	if err = p.Skip(gogqllexer.BraceL); err != nil {
		return nil, err
	}

	for {
		var enumValueDef ast.EnumValueDefinition
		enumValueDef.Description, _ = p.ReadDescription()

		if enumValueDef.Value.Value, err = p.ReadNameValue(); err != nil {
			return nil, err
		}

		if p.CheckKind(gogqllexer.At) {
			if enumValueDef.Directives, err = p.parseDirectives(); err != nil {
				return nil, err
			}
		}

		enumValuesDef = append(enumValuesDef, enumValueDef)

		if p.SkipIf(gogqllexer.BraceR) {
			break
		}
	}
	return enumValuesDef, nil
}

// https://spec.graphql.org/October2021/#sec-Enums
func (p *parser) ParseEnumTypeDefinition(description string) (def *ast.EnumTypeDefinition, err error) {
	def = &ast.EnumTypeDefinition{
		Description: description,
	}

	if err = p.SkipKeyword("enum"); err != nil {
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

	if def.EnumValue, err = p.parseEnumValuesDefinition(); err != nil {
		return nil, err
	}

	return def, nil
}

// https://spec.graphql.org/October2021/#sec-Enum-Extensions
// NOTION: consume "extend" keyword before call this function.
func (p *parser) ParseEnumTypeExtension() (def *ast.EnumTypeExtension, err error) {
	def = &ast.EnumTypeExtension{}

	if err := p.SkipKeyword("enum"); err != nil {
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

	if p.CheckKind(gogqllexer.BraceL) {
		if def.EnumValue, err = p.parseEnumValuesDefinition(); err != nil {
			return nil, err
		}
	}

	return def, nil
}
