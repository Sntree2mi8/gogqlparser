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

func (s *ScalarTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	//TODO implement me
	panic("implement me")
}

type ObjectTypeExtension struct {
}

func (o *ObjectTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	//TODO implement me
	panic("implement me")
}

type InterfaceTypeExtension struct {
}

func (i *InterfaceTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	//TODO implement me
	panic("implement me")
}

type UnionTypeExtension struct {
}

func (u *UnionTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	//TODO implement me
	panic("implement me")
}

type EnumTypeExtension struct {
	Name       string
	Directives []Directive
	EnumValue  []EnumValueDefinition
}

func (e *EnumTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	//TODO implement me
	panic("implement me")
}

type InputObjectTypeExtension struct {
}

func (i *InputObjectTypeExtension) TypeSystemExtensionKind() TypeSystemExtensionKind {
	//TODO implement me
	panic("implement me")
}
