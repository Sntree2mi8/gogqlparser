package ast

// https://spec.graphql.org/October2021/#TypeDefinition

type TypeDefinitionKind int

const (
	TypeDefinitionKindScalar TypeDefinitionKind = iota
	TypeDefinitionKindObject
	TypeDefinitionKindInterface
	TypeDefinitionKindUnion
	TypeDefinitionKindEnum
	TypeDefinitionKindInputObject
)

type TypeDefinition interface {
	TypeDefinitionKind() TypeDefinitionKind
}

var _ TypeDefinition = (*ScalarTypeDefinition)(nil)
var _ TypeDefinition = (*ObjectTypeDefinition)(nil)
var _ TypeDefinition = (*InterfaceTypeDefinition)(nil)
var _ TypeDefinition = (*UnionTypeDefinition)(nil)
var _ TypeDefinition = (*EnumTypeDefinition)(nil)
var _ TypeDefinition = (*InputObjectTypeDefinition)(nil)

type ScalarTypeDefinition struct {
	Description string
	Name        string
}

func (d *ScalarTypeDefinition) TypeDefinitionKind() TypeDefinitionKind {
	return TypeDefinitionKindScalar
}

type FieldDefinition struct {
	Description        string
	Name               string
	ArgumentDefinition []InputValueDefinition
	Type               Type
	Directives         []Directive
}

type ObjectTypeDefinition struct {
	Description      string
	Name             string
	Directives       []Directive
	FieldDefinitions []*FieldDefinition
	Interfaces       []string
}

func (d *ObjectTypeDefinition) TypeDefinitionKind() TypeDefinitionKind {
	return TypeDefinitionKindObject
}

type InterfaceTypeDefinition struct {
	Description      string
	Name             string
	Directives       []Directive
	FieldDefinitions []*FieldDefinition
	Interfaces       []string
}

func (d *InterfaceTypeDefinition) TypeDefinitionKind() TypeDefinitionKind {
	return TypeDefinitionKindInterface
}

type UnionTypeDefinition struct {
	Description string
	Name        string
	Directives  []Directive
	MemberTypes []Type
}

func (d *UnionTypeDefinition) TypeDefinitionKind() TypeDefinitionKind {
	return TypeDefinitionKindUnion
}

type EnumTypeDefinition struct {
	Description string
	Name        string
}

func (d *EnumTypeDefinition) TypeDefinitionKind() TypeDefinitionKind {
	return TypeDefinitionKindEnum
}

type InputObjectTypeDefinition struct {
	Description string
	Name        string
}

func (d *InputObjectTypeDefinition) TypeDefinitionKind() TypeDefinitionKind {
	return TypeDefinitionKindInputObject
}
