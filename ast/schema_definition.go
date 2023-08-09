package ast

type SchemaDefinition struct {
	Description                  string
	Directives                   Directives
	RootOperationTypeDefinitions []RootOperationTypeDefinition
}
