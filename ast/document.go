package ast

type TypeSystemExtensionDocument struct {
	SchemaDefinitions    []SchemaDefinition
	TypeDefinitions      []TypeDefinition
	DirectiveDefinitions []DirectiveDefinition

	SchemaExtensions     []SchemaExtension
	TypeSystemExtensions []TypeSystemExtension
}

func (d *TypeSystemExtensionDocument) Merge(others ...*TypeSystemExtensionDocument) (merged *TypeSystemExtensionDocument) {
	merged = &TypeSystemExtensionDocument{
		SchemaDefinitions:    d.SchemaDefinitions,
		TypeDefinitions:      d.TypeDefinitions,
		DirectiveDefinitions: d.DirectiveDefinitions,
	}

	for _, other := range others {
		merged.SchemaDefinitions = append(merged.SchemaDefinitions, other.SchemaDefinitions...)
		merged.TypeDefinitions = append(merged.TypeDefinitions, other.TypeDefinitions...)
		merged.DirectiveDefinitions = append(merged.DirectiveDefinitions, other.DirectiveDefinitions...)
	}

	return
}

type RootOperationTypeDefinition struct {
	Type string
}
