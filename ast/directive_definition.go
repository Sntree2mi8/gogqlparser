package ast

type Type struct {
	NamedType string
	ListType  *Type
	NotNull   bool
}

type DirectiveDefinition struct {
	Description         string
	Name                string
	ArgumentsDefinition []InputValueDefinition
	IsRepeatable        bool
	DirectiveLocations  []DirectiveLocation
}

type InputValueDefinition struct {
	Description     string
	Name            string
	Type            Type
	RawDefaultValue string
	Directives      Directives
}

type DirectiveLocation int

const (
	DirectiveLocationUnknown DirectiveLocation = iota

	// type system directive location
	DirectiveLocationSchema
	DirectiveLocationScalar
	DirectiveLocationObject
	DirectiveLocationFieldDefinition
	DirectiveLocationArgumentDefinition
	DirectiveLocationInterface
	DirectiveLocationUnion
	DirectiveLocationEnum
	DirectiveLocationEnumValue
	DirectiveLocationInputObject
	DirectiveLocationInputFieldDefinition

	// executable directive location
	DirectiveLocationQuery
	DirectiveLocationMutation
	DirectiveLocationSubscription
	DirectiveLocationField
	DirectiveLocationFragmentDefinition
	DirectiveLocationFragmentSpread
	DirectiveLocationInlineFragment
	DirectiveLocationVariableDefinition
)

// 以下はdefinitionではない
type Directives []Directive

type Directive struct {
	Name      string
	Arguments Arguments
}

type Arguments []Argument

type Argument struct {
	Name  string
	Value string
}
