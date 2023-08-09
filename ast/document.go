package ast

type TypeSystemExtensionDocument struct {
	SchemaDefinitions    []SchemaDefinition
	TypeDefinitions      []TypeDefinition
	DirectiveDefinitions []DirectiveDefinition
}

// 以下、未整理

type RootOperationTypeKind int

const (
	OperationTypeQuery RootOperationTypeKind = iota
	OperationTypeMutation
	OperationTypeSubscription
)

type RootOperationTypeDefinition struct {
	OperationType RootOperationTypeKind
	Type          string
}
