package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"log"
)

// https://spec.graphql.org/October2021/#UnionMemberTypes
func (p *parser) parseUnionMemberTypes() (memberTypes []ast.Type, err error) {
	if err = p.Skip(gogqllexer.Equal); err != nil {
		return nil, err
	}

	p.SkipIf(gogqllexer.Pipe)

	for {
		var mt ast.Type
		if mt.NamedType, err = p.ReadNameValue(); err != nil {
			log.Println("here")
			return nil, err
		}

		memberTypes = append(memberTypes, mt)

		if !p.SkipIf(gogqllexer.Pipe) {
			break
		}
	}

	return memberTypes, nil
}

// https://spec.graphql.org/October2021/#sec-Unions
func (p *parser) ParseUnionTypeDefinition(description string) (def *ast.UnionTypeDefinition, err error) {
	def = &ast.UnionTypeDefinition{
		Description: description,
	}

	if err = p.SkipKeyword("union"); err != nil {
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

	if err = p.PeekAndMustBe(
		[]gogqllexer.Kind{gogqllexer.Equal},
		func(t gogqllexer.Token, advanceLexer func()) error {
			if def.MemberTypes, err = p.parseUnionMemberTypes(); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return def, nil
}

// ParseUnionTypeExtension parses union type extension.
// "extend" keyword must be consumed before calling this function.
//
// Reference: https://spec.graphql.org/October2021/#sec-Union-Extensions
func (p *parser) ParseUnionTypeExtension() (def *ast.UnionTypeExtension, err error) {
	def = &ast.UnionTypeExtension{}

	if err = p.SkipKeyword("union"); err != nil {
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

	if p.CheckKind(gogqllexer.Equal) {
		if def.MemberTypes, err = p.parseUnionMemberTypes(); err != nil {
			return nil, err
		}
	}

	return def, err
}
