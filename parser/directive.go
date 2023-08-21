package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func parseDirectiveArguments(l *LexerWrapper) (args []ast.Argument, err error) {
	if err = l.Skip(gogqllexer.ParenL); err != nil {
		return nil, err
	}

	for {
		arg := ast.Argument{}

		if arg.Name, err = l.ReadNameValue(); err != nil {
			return nil, err
		}
		if l.Skip(gogqllexer.Colon) != nil {
			return nil, err
		}
		if err = l.PeekAndMustBe(
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

		if l.SkipIf(gogqllexer.ParenR) {
			break
		}
	}

	return args, err
}

func parseDirectives(l *LexerWrapper) (directives []ast.Directive, err error) {
	for {
		var d ast.Directive

		d, err = parseDirective(l)
		if err != nil {
			return nil, err
		}

		directives = append(directives, d)

		if !l.CheckKind(gogqllexer.At) {
			break
		}
	}

	return directives, nil
}

func parseDirective(l *LexerWrapper) (d ast.Directive, err error) {
	if err = l.Skip(gogqllexer.At); err != nil {
		return d, err
	}

	if d.Name, err = l.ReadNameValue(); err != nil {
		return d, err
	}

	if l.CheckKind(gogqllexer.ParenL) {
		if d.Arguments, err = parseDirectiveArguments(l); err != nil {
			return d, err
		}
	}

	return d, err
}

func ParseDirectiveDefinition(l *LexerWrapper, description string) (def *ast.DirectiveDefinition, err error) {
	l.NextToken()
	t := l.NextToken()
	if t.Kind != gogqllexer.At {
		return nil, fmt.Errorf("unexpected token %+v", t)
	}
	t = l.NextToken()
	if t.Kind != gogqllexer.Name {
		return nil, fmt.Errorf("unexpected token %+v", t)
	}

	directiveName := t.Value
	isRepeatable := false

	t = l.PeekToken()
	var inputValueDefinitions []ast.InputValueDefinition
	if t.Kind == gogqllexer.ParenL {
		inputValueDefinitions, err = ParseArgumentsDefinition(l)
		if err != nil {
			return nil, err
		}
	}

	t = l.PeekToken()
	if t.Kind != gogqllexer.Name {
		return nil, fmt.Errorf("unexpected token %+v", t)
	}

	if t.Value == "repeatable" {
		isRepeatable = true
		l.NextToken()
	}

	t = l.NextToken()
	if t.Kind != gogqllexer.Name {
		return nil, fmt.Errorf("unexpected token %+v", t)
	} else if t.Value != "on" {
		return nil, fmt.Errorf("unexpected token %+v", t)
	}

	t = l.NextToken()
	if t.Kind != gogqllexer.Name {
		return nil, fmt.Errorf("unexpected token %+v", t)
	}
	dLocations := make([]ast.DirectiveLocation, 0)
	dl := ParsDirectiveLocation(t.Value)
	if dl == ast.DirectiveLocationUnknown {
		return nil, fmt.Errorf("unexpected token %+v", t)
	}
	dLocations = append(dLocations, dl)

	for {
		t = l.PeekToken()
		if t.Kind != gogqllexer.Pipe {
			break
		} else {
			l.NextToken()
		}

		t = l.NextToken()
		if t.Kind != gogqllexer.Name {
			return nil, fmt.Errorf("unexpected token %+v", t)
		}
		dl := ParsDirectiveLocation(t.Value)
		if dl == ast.DirectiveLocationUnknown {
			return nil, fmt.Errorf("unexpected token %+v", t)
		}
		dLocations = append(dLocations, dl)
	}

	def = &ast.DirectiveDefinition{
		Description:         description,
		Name:                directiveName,
		ArgumentsDefinition: inputValueDefinitions,
		IsRepeatable:        isRepeatable,
		DirectiveLocations:  dLocations,
	}

	return def, nil
}
