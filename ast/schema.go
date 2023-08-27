package ast

type Schema struct {
	Description string

	Query        *RootOperationTypeDefinition
	Mutation     *RootOperationTypeDefinition
	Subscription *RootOperationTypeDefinition

	Types      map[string]TypeDefinition
	Directives map[string]*DirectiveDefinition
}
