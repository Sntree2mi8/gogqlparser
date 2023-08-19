package ast

type SchemaDefinition struct {
	Description  string
	Directives   []Directive
	Query        *RootOperationTypeDefinition
	Mutation     *RootOperationTypeDefinition
	Subscription *RootOperationTypeDefinition
}
