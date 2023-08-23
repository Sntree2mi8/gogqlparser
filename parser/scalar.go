package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

// ParseScalarTypeExtension parse scalar type extension.
// "extend" keyword must be consumed before calling this function.
//
// Reference: https://spec.graphql.org/October2021/#sec-Scalar-Extensions
func (p *parser) ParseScalarTypeExtension() (def *ast.ScalarTypeExtension, err error) {
	def = &ast.ScalarTypeExtension{}

	if err = p.SkipKeyword("scalar"); err != nil {
		return nil, err
	}

	if def.Name, err = p.ReadNameValue(); err != nil {
		return nil, err
	}

	if p.CheckKind(gogqllexer.At) {
		if def.Directives, err = p.parseDirectives(); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("scalar type extension needs at least one directive")
	}

	return def, nil
}
