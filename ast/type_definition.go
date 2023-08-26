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
	TypeName() string
	GetDirectives() []Directive
}

type ScalarTypeDefinition struct {
	Name string
}

func (d *ScalarTypeDefinition) TypeDefinitionKind() TypeDefinitionKind {
	return TypeDefinitionKindScalar
}

func (d *ScalarTypeDefinition) TypeName() string {
	return d.Name
}

func (d *ScalarTypeDefinition) GetDirectives() []Directive {
	return nil
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

func (d *ObjectTypeDefinition) TypeName() string {
	return d.Name
}

func (d *ObjectTypeDefinition) GetDirectives() []Directive {
	return d.Directives
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

func (d *InterfaceTypeDefinition) TypeName() string {
	return d.Name
}

func (d *InterfaceTypeDefinition) GetDirectives() []Directive {
	return d.Directives
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

func (d *UnionTypeDefinition) TypeName() string {
	return d.Name
}

func (d *UnionTypeDefinition) GetDirectives() []Directive {
	return d.Directives
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

func (d *EnumTypeDefinition) TypeName() string {
	return d.Name
}

func (d *EnumTypeDefinition) GetDirectives() []Directive {
	return d.Directives
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

func (d *InputObjectTypeDefinition) TypeName() string {
	return d.Name
}

func (d *InputObjectTypeDefinition) GetDirectives() []Directive {
	return d.Directives
}
