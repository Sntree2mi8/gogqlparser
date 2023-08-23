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
	Directives  []Directive
	EnumValue   []EnumValueDefinition
}

type EnumValueDefinition struct {
	Description string
	Value       EnumValue
	Directives  []Directive
}

func (d *EnumTypeDefinition) TypeDefinitionKind() TypeDefinitionKind {
	return TypeDefinitionKindEnum
}

type InputObjectTypeDefinition struct {
	Description string
	Name        string
	Directives  []Directive
	InputFields []InputValueDefinition
}

func (d *InputObjectTypeDefinition) TypeDefinitionKind() TypeDefinitionKind {
	return TypeDefinitionKindInputObject
}
