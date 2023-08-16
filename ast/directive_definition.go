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
	Directives      []Directive
}

type Directive struct {
	Name      string
	Arguments []Argument
}

type Argument struct {
	Name  string
	Value string
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
