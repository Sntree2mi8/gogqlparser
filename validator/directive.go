package validator

import (
	"fmt"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"log"
	"slices"
	"strings"
)

func (v *validator) isInputType(t ast.Type) bool {
	if t.ListType != nil {
		return v.isInputType(*t.ListType)
	}

	td, ok := v.typeDefs[t.NamedType]
	if !ok {
		log.Println("type not found: ", t.NamedType)
		return false
	}
	switch td.TypeDefinitionKind() {
	case ast.TypeDefinitionKindScalar, ast.TypeDefinitionKindEnum, ast.TypeDefinitionKindInputObject:
		return true
	default:
		return false
	}
}

func (v *validator) checkDirectiveReferenceInDirective(dd ast.DirectiveDefinition, d ast.Directive) (hasReference bool) {
	if d.Name == dd.Name {
		return true
	}

	for _, ad := range v.directiveDefs[d.Name].ArgumentsDefinition {
		for _, adDir := range ad.Directives {
			return v.checkDirectiveReferenceInDirective(dd, adDir)
		}
	}

	return false
}

func (v *validator) checkDirectiveReferenceInType(dd ast.DirectiveDefinition, t ast.Type) (hasReference bool) {
	typeDef := v.typeDefs[t.NamedType]

	switch typeDef.TypeDefinitionKind() {
	case ast.TypeDefinitionKindScalar:
		// TODO: 本当はscalarにもdirectiveをつけられるけど、まだ実装できていない
		return false
	case ast.TypeDefinitionKindEnum:
		if slices.Contains(dd.DirectiveLocations, ast.DirectiveLocationEnum) {
			for _, td := range typeDef.GetDirectives() {
				if v.checkDirectiveReferenceInDirective(dd, td) {
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

					if v.checkDirectiveReferenceInDirective(dd, evd) {
						return true
					}
				}
			}
		}
	case ast.TypeDefinitionKindInputObject:
		if slices.Contains(dd.DirectiveLocations, ast.DirectiveLocationInputObject) {
			for _, td := range typeDef.GetDirectives() {
				if v.checkDirectiveReferenceInDirective(dd, td) {
					return true
				}
			}
		}
		if slices.Contains(dd.DirectiveLocations, ast.DirectiveLocationInputFieldDefinition) {
			inputObjectTypeDef := typeDef.(*ast.InputObjectTypeDefinition)
			for _, ifd := range inputObjectTypeDef.InputFields {
				for _, ifdd := range ifd.Directives {
					if v.checkDirectiveReferenceInDirective(dd, ifdd) {
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

func (v *validator) validateDirectiveDefinitions() error {
	for _, dd := range v.doc.DirectiveDefinitions {
		if strings.HasPrefix(dd.Name, "__") {
			return fmt.Errorf("directive name must not begins with \"__\": %s", dd.Name)
		}

		canUseOnArgumentDefinition := slices.Contains(dd.DirectiveLocations, ast.DirectiveLocationArgumentDefinition)
		for _, ad := range dd.ArgumentsDefinition {
			if canUseOnArgumentDefinition {
				for _, adDir := range ad.Directives {
					// argumentDefinitionに使用したdirectiveが直接参照もしくは間接的に参照をしていないかを確認する
					if v.checkDirectiveReferenceInDirective(dd, adDir) {
						return fmt.Errorf("directive %s must not contain the use of a directive which references itself", dd.Name)
					}
				}
			}

			if strings.HasPrefix(ad.Name, "__") {
				return fmt.Errorf("argument name must not begins with \"__\": %s", ad.Name)
			}

			if !v.isInputType(ad.Type) {
				return fmt.Errorf("argument %s must accept a type where IsInputType(argumentType) returns true", ad.Name)
			}

			if v.checkDirectiveReferenceInType(dd, ad.Type) {
				return fmt.Errorf("argument %s must not contain the use of a directive which references itself", ad.Name)
			}
		}
	}

	return nil
}
