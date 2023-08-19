package ast

type TypeSystemExtensionDocument struct {
	SchemaDefinitions    []SchemaDefinition
	TypeDefinitions      map[string]TypeDefinition
	DirectiveDefinitions []DirectiveDefinition
}

type RootOperationTypeDefinition struct {
	Type string
}
