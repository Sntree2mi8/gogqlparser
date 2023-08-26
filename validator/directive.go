package validator

import (
	"fmt"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"slices"
	"strings"
)

func (v *validator) isInputType(t ast.Type) (bool, error) {
	if t.ListType != nil {
		return v.isInputType(*t.ListType)
	}

	td, ok := v.typeDefs[t.NamedType]
	if !ok {
		return false, fmt.Errorf("undefined type: %s", t.NamedType)
	}
	switch td.TypeDefinitionKind() {
	case ast.TypeDefinitionKindScalar, ast.TypeDefinitionKindEnum, ast.TypeDefinitionKindInputObject:
		return true, nil
	default:
		return false, nil
	}
}

// この関数はastに送っても良さそう
func getUnderlyingType(t ast.Type) ast.Type {
	if t.ListType != nil {
		return getUnderlyingType(*t.ListType)
	}

	return t
}

func (v *validator) checkSelfDirectiveReferenceInDirective(self ast.DirectiveDefinition, d ast.Directive) (hasReference bool) {
	if d.Name == self.Name {
		return true
	}

	for _, ad := range v.directiveDefs[d.Name].ArgumentsDefinition {
		for _, adDir := range ad.Directives {
			return v.checkSelfDirectiveReferenceInDirective(self, adDir)
		}
	}

	return false
}

func (v *validator) checkSelfDirectiveReferenceInType(self ast.DirectiveDefinition, t ast.Type) (hasReference bool) {
	typeDef := v.typeDefs[getUnderlyingType(t).NamedType]

	// 自己参照を判定するこの関数自体がdirectiveの文脈からしか呼ばれないのでScalar, Enum, InputObject以外は考慮しない
	switch typeDef.TypeDefinitionKind() {
	case ast.TypeDefinitionKindScalar:
		if slices.Contains(self.DirectiveLocations, ast.DirectiveLocationScalar) {
			for _, td := range typeDef.GetDirectives() {
				if v.checkSelfDirectiveReferenceInDirective(self, td) {
					return true
				}
			}
		}
		return false
	case ast.TypeDefinitionKindEnum:
		if slices.Contains(self.DirectiveLocations, ast.DirectiveLocationEnum) {
			for _, td := range typeDef.GetDirectives() {
				if v.checkSelfDirectiveReferenceInDirective(self, td) {
					return true
				}
			}
		}
		if slices.Contains(self.DirectiveLocations, ast.DirectiveLocationEnumValue) {
			enumTypeDef := typeDef.(*ast.EnumTypeDefinition)
			for _, ev := range enumTypeDef.EnumValue {
				for _, evd := range ev.Directives {
					if v.checkSelfDirectiveReferenceInDirective(self, evd) {
						return true
					}
				}
			}
		}
	case ast.TypeDefinitionKindInputObject:
		if slices.Contains(self.DirectiveLocations, ast.DirectiveLocationInputObject) {
			for _, td := range typeDef.GetDirectives() {
				if v.checkSelfDirectiveReferenceInDirective(self, td) {
					return true
				}
			}
		}
		if slices.Contains(self.DirectiveLocations, ast.DirectiveLocationInputFieldDefinition) {
			inputObjectTypeDef := typeDef.(*ast.InputObjectTypeDefinition)
			for _, ifd := range inputObjectTypeDef.InputFields {
				for _, ifdd := range ifd.Directives {
					if v.checkSelfDirectiveReferenceInDirective(self, ifdd) {
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
			if strings.HasPrefix(ad.Name, "__") {
				return fmt.Errorf("argument name must not begins with \"__\": %s", ad.Name)
			}

			inputType, err := v.isInputType(ad.Type)
			if err != nil {
				return err
			}
			if !inputType {
				return fmt.Errorf("argument %s must be input type (scalar, enum, input object)", ad.Name)
			}

			// typeにより間接的に自分自身をreferenceしていないかを確認する
			if v.checkSelfDirectiveReferenceInType(dd, ad.Type) {
				return fmt.Errorf("argument %s must not contain the use of a directive which references itself", ad.Name)
			}

			// argument definitionにより間接的に自分自身をreferenceしていないかを確認する
			if canUseOnArgumentDefinition {
				for _, adDir := range ad.Directives {
					if v.checkSelfDirectiveReferenceInDirective(dd, adDir) {
						return fmt.Errorf("directive %s must not contain the use of a directive which references itself", dd.Name)
					}
				}
			}
		}
	}

	return nil
}
