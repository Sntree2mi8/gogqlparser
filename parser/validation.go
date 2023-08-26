package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"slices"
	"strings"
)

func ValidateTypeSystemExtensionDocument(doc *ast.TypeSystemExtensionDocument) error {
	return newValidator(doc).validateDirectiveDefinitions()
}

type validator struct {
	doc *ast.TypeSystemExtensionDocument

	typeDefs      map[string]ast.TypeDefinition
	directiveDefs map[string]ast.DirectiveDefinition
}

func newValidator(doc *ast.TypeSystemExtensionDocument) *validator {
	typeDefs := make(map[string]ast.TypeDefinition, len(doc.TypeDefinitions))
	for _, def := range doc.TypeDefinitions {
		typeDefs[def.TypeName()] = def
	}

	directiveDefs := make(map[string]ast.DirectiveDefinition, len(doc.DirectiveDefinitions))
	for _, def := range doc.DirectiveDefinitions {
		directiveDefs[def.Name] = def
	}

	return &validator{
		doc:           doc,
		typeDefs:      typeDefs,
		directiveDefs: directiveDefs,
	}
}

func (v *validator) validateDirectiveDefinitions() error {
	var checkReferenceInDirective func(dd ast.DirectiveDefinition, d ast.Directive) (hasReference bool)
	checkReferenceInDirective = func(dd ast.DirectiveDefinition, d ast.Directive) (hasReference bool) {
		if d.Name == dd.Name {
			return true
		}

		for _, ad := range v.directiveDefs[d.Name].ArgumentsDefinition {
			for _, adDir := range ad.Directives {
				return checkReferenceInDirective(dd, adDir)
			}
		}

		return false
	}

	var checkReferenceInType func(dd ast.DirectiveDefinition, t ast.Type) (hasReference bool)
	checkReferenceInType = func(dd ast.DirectiveDefinition, t ast.Type) (hasReference bool) {
		typeDef := v.typeDefs[t.NamedType]

		switch typeDef.TypeDefinitionKind() {
		case ast.TypeDefinitionKindScalar:
			// TODO: 本当はscalarにもdirectiveをつけられるけど、まだ実装できていない
			return false
		case ast.TypeDefinitionKindEnum:
			if slices.Contains(dd.DirectiveLocations, ast.DirectiveLocationEnum) {
				for _, td := range typeDef.GetDirectives() {
					if checkReferenceInDirective(dd, td) {
						return true
					}
				}
			}
			if slices.Contains(dd.DirectiveLocations, ast.DirectiveLocationEnumValue) {
				enumTypeDef := typeDef.(*ast.EnumTypeDefinition)
				for _, ev := range enumTypeDef.EnumValue {
					for _, evd := range ev.Directives {
						if evd.Name == dd.Name {
							return true
						}

						if checkReferenceInDirective(dd, evd) {
							return true
						}
					}
				}
			}
		case ast.TypeDefinitionKindInputObject:
			if slices.Contains(dd.DirectiveLocations, ast.DirectiveLocationInputObject) {
				for _, td := range typeDef.GetDirectives() {
					if checkReferenceInDirective(dd, td) {
						return true
					}
				}
			}
			if slices.Contains(dd.DirectiveLocations, ast.DirectiveLocationInputFieldDefinition) {
				inputObjectTypeDef := typeDef.(*ast.InputObjectTypeDefinition)
				for _, ifd := range inputObjectTypeDef.InputFields {
					for _, ifdd := range ifd.Directives {
						if checkReferenceInDirective(dd, ifdd) {
							return true
						}
					}
				}
			}
		default:
			return false
		}

		return false
	}

	var isInputType func(t ast.Type) bool
	isInputType = func(t ast.Type) bool {
		if t.ListType != nil {
			return isInputType(*t.ListType)
		}

		td := v.typeDefs[t.NamedType]
		switch td.TypeDefinitionKind() {
		case ast.TypeDefinitionKindScalar, ast.TypeDefinitionKindEnum, ast.TypeDefinitionKindInputObject:
			return true
		default:
			return false
		}
	}

	for _, dd := range v.doc.DirectiveDefinitions {
		if strings.HasPrefix(dd.Name, "__") {
			return fmt.Errorf("directive name must not begins with \"__\": %s", dd.Name)
		}

		canUseOnArgumentDefinition := slices.Contains(dd.DirectiveLocations, ast.DirectiveLocationArgumentDefinition)
		for _, ad := range dd.ArgumentsDefinition {
			if canUseOnArgumentDefinition {
				for _, adDir := range ad.Directives {
					// argumentDefinitionに使用したdirectiveが直接参照もしくは間接的に参照をしていないかを確認する
					if checkReferenceInDirective(dd, adDir) {
						return fmt.Errorf("directive %s must not contain the use of a directive which references itself", dd.Name)
					}
				}
			}

			if strings.HasPrefix(ad.Name, "__") {
				return fmt.Errorf("argument name must not begins with \"__\": %s", ad.Name)
			}

			if !isInputType(ad.Type) {
				return fmt.Errorf("argument %s must accept a type where IsInputType(argumentType) returns true", ad.Name)
			}

			if checkReferenceInType(dd, ad.Type) {
				return fmt.Errorf("argument %s must not contain the use of a directive which references itself", ad.Name)
			}
		}
	}

	return nil
}
