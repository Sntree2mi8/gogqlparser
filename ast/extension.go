package ast

type TypeSystemExtensionKind int

const (
	TypeSystemExtensionKindScalar TypeSystemExtensionKind = iota
	TypeSystemExtensionKindObject
	TypeSystemExtensionKindInterface
	TypeSystemExtensionKindUnion
	TypeSystemExtensionKindEnum
	TypeSystemExtensionKindInputObject
)

type TypeSystemExtension interface {
	TypeSystemExtensionKind() TypeSystemExtensionKind
}

type ScalarTypeExtension struct {
	Name       string
	Directives []Directive
}

func (e *ScalarTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	return TypeSystemExtensionKindScalar
}

type ObjectTypeExtension struct {
	Name                string
	Directives          []Directive
	FieldsDefinition    []*FieldDefinition
	ImplementInterfaces []string
}

func (e *ObjectTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	return TypeSystemExtensionKindObject
}

type InterfaceTypeExtension struct {
	Name                string
	ImplementInterfaces []string
	Directives          []Directive
	FieldsDefinition    []*FieldDefinition
}

func (e *InterfaceTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	return TypeSystemExtensionKindInterface
}

type UnionTypeExtension struct {
	Name        string
	Directives  []Directive
	MemberTypes []Type
}

func (e *UnionTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	return TypeSystemExtensionKindUnion
}

type EnumTypeExtension struct {
	Name       string
	Directives []Directive
	EnumValue  []EnumValueDefinition
}

func (e *EnumTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	return TypeSystemExtensionKindEnum
}

type InputObjectTypeExtension struct {
	Name                  string
	Directives            []Directive
	InputsFieldDefinition []InputValueDefinition
}

func (e *InputObjectTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	return TypeSystemExtensionKindInputObject
}

type SchemaExtension struct {
	Directives   []Directive
	Query        *RootOperationTypeDefinition
	Mutation     *RootOperationTypeDefinition
	Subscription *RootOperationTypeDefinition
}
