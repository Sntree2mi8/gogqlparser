package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func (p *parser) parseDirectiveArguments() (args []ast.Argument, err error) {
	if err = p.Skip(gogqllexer.ParenL); err != nil {
		return nil, err
	}

	for {
		arg := ast.Argument{}

		if arg.Name, err = p.ReadNameValue(); err != nil {
			return nil, err
		}
		if p.Skip(gogqllexer.Colon) != nil {
			return nil, err
		}
		if err = p.PeekAndMustBe(
			[]gogqllexer.Kind{gogqllexer.Int, gogqllexer.Float, gogqllexer.Name, gogqllexer.String, gogqllexer.BlockString},
			func(t gogqllexer.Token, advanceLexer func()) error {
				defer advanceLexer()
				arg.Value = t.Value
				return nil
			},
		); err != nil {
			return nil, err
		}

		args = append(args, arg)

		if p.SkipIf(gogqllexer.ParenR) {
			break
		}
	}

	return args, err
}

func (p *parser) parseDirectives() (directives []ast.Directive, err error) {
	for {
		var d ast.Directive

		d, err = p.parseDirective()
		if err != nil {
			return nil, err
		}

		directives = append(directives, d)

		if !p.CheckKind(gogqllexer.At) {
			break
		}
	}

	return directives, nil
}

func (p *parser) parseDirective() (d ast.Directive, err error) {
	if err = p.Skip(gogqllexer.At); err != nil {
		return d, err
	}

	if d.Name, err = p.ReadNameValue(); err != nil {
		return d, err
	}

	if p.CheckKind(gogqllexer.ParenL) {
		if d.Arguments, err = p.parseDirectiveArguments(); err != nil {
			return d, err
		}
	}

	return d, err
}

func parseDirectiveLocation(v string) ast.DirectiveLocation {
	switch v {
	case "QUERY":
		return ast.DirectiveLocationQuery
	case "MUTATION":
		return ast.DirectiveLocationMutation
	case "SUBSCRIPTION":
		return ast.DirectiveLocationSubscription
	case "FIELD":
		return ast.DirectiveLocationField
	case "FRAGMENT_DEFINITION":
		return ast.DirectiveLocationFragmentDefinition
	case "FRAGMENT_SPREAD":
		return ast.DirectiveLocationFragmentSpread
	case "INLINE_FRAGMENT":
		return ast.DirectiveLocationInlineFragment
	case "VARIABLE_DEFINITION":
		return ast.DirectiveLocationVariableDefinition

	case "SCHEMA":
		return ast.DirectiveLocationSchema
	case "SCALAR":
		return ast.DirectiveLocationScalar
	case "OBJECT":
		return ast.DirectiveLocationObject
	case "FIELD_DEFINITION":
		return ast.DirectiveLocationFieldDefinition
	case "ARGUMENT_DEFINITION":
		return ast.DirectiveLocationArgumentDefinition
	case "INTERFACE":
		return ast.DirectiveLocationInterface
	case "UNION":
		return ast.DirectiveLocationUnion
	case "ENUM":
		return ast.DirectiveLocationEnum
	case "ENUM_VALUE":
		return ast.DirectiveLocationEnumValue
	case "INPUT_OBJECT":
		return ast.DirectiveLocationInputObject
	case "INPUT_FIELD_DEFINITION":
		return ast.DirectiveLocationInputFieldDefinition
	default:
		return ast.DirectiveLocationUnknown
	}
}

func (p *parser) parseDirectiveLocations() (locs []ast.DirectiveLocation, err error) {
	p.SkipIf(gogqllexer.Pipe)

	for {
		var locationValue string
		if locationValue, err = p.ReadNameValue(); err != nil {
			return nil, err
		}
		loc := parseDirectiveLocation(locationValue)
		if loc == ast.DirectiveLocationUnknown {
			return nil, fmt.Errorf("unexpected token %+v", locationValue)
		}

		locs = append(locs, loc)

		if p.SkipIf(gogqllexer.Pipe) {
			continue
		}
		break
	}

	return locs, nil
}

// ParseDirectiveDefinition parses a directive definition
//
// Reference: https://spec.graphqp.org/October2021/#sec-Type-System.Directives
func (p *parser) ParseDirectiveDefinition(description string) (def *ast.DirectiveDefinition, err error) {
	def = &ast.DirectiveDefinition{
		Description: description,
	}

	if err = p.SkipKeyword("directive"); err != nil {
		return nil, err
	}
	if err = p.Skip(gogqllexer.At); err != nil {
		return nil, err
	}
	if def.Name, err = p.ReadNameValue(); err != nil {
		return nil, err
	}

	if p.CheckKind(gogqllexer.ParenL) {
		if def.ArgumentsDefinition, err = p.ParseArgumentsDefinition(); err != nil {
			return nil, err
		}
	}

	if p.SkipKeywordIf("repeatable") {
		def.IsRepeatable = true
	}

	if err = p.SkipKeyword("on"); err != nil {
		return nil, err
	}

	def.DirectiveLocations, err = p.parseDirectiveLocations()
	if err != nil {
		return nil, err
	}

	return def, nil
}
