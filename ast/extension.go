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

var _ TypeSystemExtension = (*ScalarTypeExtension)(nil)
var _ TypeSystemExtension = (*ObjectTypeExtension)(nil)
var _ TypeSystemExtension = (*InterfaceTypeExtension)(nil)
var _ TypeSystemExtension = (*UnionTypeExtension)(nil)
var _ TypeSystemExtension = (*EnumTypeExtension)(nil)
var _ TypeSystemExtension = (*InputObjectTypeExtension)(nil)

type ScalarTypeExtension struct {
}

func (e *ScalarTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	return TypeSystemExtensionKindScalar
}

type ObjectTypeExtension struct {
}

func (e *ObjectTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	return TypeSystemExtensionKindObject
}

type InterfaceTypeExtension struct {
	Name                string
	ImplementInterfaces []string
	Directive           []Directive
	FieldsDefinition    []FieldDefinition
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
}

func (e *InputObjectTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	return TypeSystemExtensionKindInputObject
}
