package ast

type TypeSystemExtensionDocument struct {
	SchemaDefinitions    []SchemaDefinition
	TypeDefinitions      []TypeDefinition
	DirectiveDefinitions []DirectiveDefinition
}

type RootOperationTypeDefinition struct {
	Type string
}
