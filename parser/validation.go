package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"strings"
)

// This is a place to put validation related functions. At the time of creation, I couldn't imagine how to divide it, so I'll put it here for now.
func ValidateTypeSystemExtensionDocument(doc *ast.TypeSystemExtensionDocument) error {
	return validateDirectiveDefinitions(doc)
}

func validateDirectiveDefinitions(doc *ast.TypeSystemExtensionDocument) error {
	// validate directive
	// 1. A directive definition must not contain the use of a directive which references itself directly.
	// TODO: implement

	// 2. A directive definition must not contain the use of a directive which references itself indirectly by referencing a Type or Directive which transitively includes a reference to this directive.
	// TODO: implement

	// 3. The directive must not have a name which begins with the characters "__" (two underscores).
	for _, dd := range doc.DirectiveDefinitions {
		if strings.HasPrefix(dd.Name, "__") {
			return fmt.Errorf("directive name must not begins with \"__\": %s", dd.Name)
		}

		// for each argument
		for _, ad := range dd.ArgumentsDefinition {
			// 4. The argument must not have a name which begins with the characters "__" (two underscores).
			if strings.HasPrefix(ad.Name, "__") {
				return fmt.Errorf("argument name must not begins with \"__\": %s", ad.Name)
			}

			// 5. The argument must accept a type where IsInputType(argumentType) returns true.
			// これを判断するためにはtype definitionsが必要になる
			// argumentに割り当てられた方をtype definitionsから探してそれがinputTypeである必要がある
			// TODO: implement
		}
	}

	return nil
}
