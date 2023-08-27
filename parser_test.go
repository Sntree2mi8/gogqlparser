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

enum HogeEnum {
	FOO
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

	if _, err = validator.ValidateTypeSystemExtensionDocument(d); err != nil {
		t.Fatal(err)
	}
}
