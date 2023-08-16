package ast

type SchemaDefinition struct {
	Description                  string
	Directives                   []Directive
	RootOperationTypeDefinitions []RootOperationTypeDefinition
}
