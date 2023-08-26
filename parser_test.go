package gogqlparser

import (
	"github.com/Sntree2mi8/gogqlparser/ast"
	"github.com/Sntree2mi8/gogqlparser/validator"
	"testing"
)

var schema = `
"""
This is schema
"""
schema @schema_directive @schema_directive {
	query: Query
}

type Query {
	hello: String
}

scalar Hoge

directive @schema_directive repeatable on SCHEMA
directive @arg_directive(arg1: String) on INPUT_OBJECT
input InputObjectForArgDirective @arg_directive(arg1: "arg1") {
	field1: String
}
`

func TestParser_ParseTypeSystem(t *testing.T) {
	p := New()
	d, err := p.parseTypeSystemDocument(&ast.Source{
		Name: "schema.graphql",
		Body: schema,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = validator.ValidateTypeSystemExtensionDocument(d); err != nil {
		t.Fatal(err)
	}
}
