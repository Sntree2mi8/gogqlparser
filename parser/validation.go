package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"strings"
)

// This is a place to put validation related functions. At the time of creation, I couldn't imagine how to divide it, so I'll put it here for now.
func ValidateTypeSystemExtensionDocument(doc *ast.TypeSystemExtensionDocument) error {
	return validateDirectiveDefinitions(doc.DirectiveDefinitions)
}

func validateDirectiveDefinitions(dds []ast.DirectiveDefinition) error {
	// validate directive
	// 1. A directive definition must not contain the use of a directive which references itself directly.
	// 2. A directive definition must not contain the use of a directive which references itself indirectly by referencing a Type or Directive which transitively includes a reference to this directive.
	// 3. The directive must not have a name which begins with the characters "__" (two underscores).
	for _, dd := range dds {
		if strings.HasPrefix(dd.Name, "__") {
			return fmt.Errorf("directive name must not begin with \"__\": %s", dd.Name)
		}
	}
	// for each argument
	// 4. The argument must not have a name which begins with the characters "__" (two underscores).
	// 5. The argument must accept a type where IsInputType(argumentType) returns true.

	return nil
}
