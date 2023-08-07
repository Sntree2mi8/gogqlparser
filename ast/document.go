package ast

type SchemaDefinition struct {
	Description                  string
	Directives                   Directives
	RootOperationTypeDefinitions []RootOperationTypeDefinition
}

type DirectiveDefinition struct{}

type Directives []Directive

type Directive struct {
	Name      string
	Arguments Arguments
}

type Arguments []Argument

type Argument struct {
	Name  string
	Value string
}

type OperationType int

const (
	OperationTypeQuery OperationType = iota
	OperationTypeMutation
	OperationTypeSubscription
)

type Type struct {
	Name string
}

type RootOperationTypeDefinition struct {
	OperationType OperationType
	Type          Type
}
