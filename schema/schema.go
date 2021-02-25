package schema

import (
	"github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/internal/common"
	"github.com/graph-gophers/graphql-go/types"
)

// Schema represents a GraphQL service's collective type system capabilities.
// A schema is defined in terms of the types and directives it supports as well as the root
// operation types for each kind of operation: `query`, `mutation`, and `subscription`.
//
// For a more formal definition, read the relevant section in the specification:
//
// http://facebook.github.io/graphql/draft/#sec-Schema
type Schema struct {
	// EntryPoints determines the place in the type system where `query`, `mutation`, and
	// `subscription` operations begin.
	//
	// http://facebook.github.io/graphql/draft/#sec-Root-Operation-Types
	//
	RootOperationTypes map[string]NamedType

	// Types are the fundamental unit of any GraphQL schema.
	// There are six kinds of named types, and two wrapping types.
	//
	// http://facebook.github.io/graphql/draft/#sec-Types
	Types map[string]NamedType

	// Directives are used to annotate various parts of a GraphQL document as an indicator that they
	// should be evaluated differently by a validator, executor, or client tool such as a code
	// generator.
	//
	// http://facebook.github.io/graphql/draft/#sec-Type-System.Directives
	Directives map[string]*DirectiveDecl

	UseFieldResolvers bool
}

// Resolve a named type in the schema by its name.
func (s *Schema) Resolve(name string) common.Type {
	return s.Types[name]
}

// NamedType represents a type with a name.
//
// http://facebook.github.io/graphql/draft/#NamedType
type NamedType interface {
	common.Type
	TypeName() string
	Description() string
}

// Scalar types represent primitive leaf values (e.g. a string or an integer) in a GraphQL type
// system.
//
// GraphQL responses take the form of a hierarchical tree; the leaves on these trees are GraphQL
// scalars.
//
// http://facebook.github.io/graphql/draft/#sec-Scalars
type Scalar struct {
	Name       string
	Desc       string
	Directives common.DirectiveList
}

// Object types represent a list of named fields, each of which yield a value of a specific type.
//
// GraphQL queries are hierarchical and composed, describing a tree of information.
// While Scalar types describe the leaf values of these hierarchical types, Objects describe the
// intermediate levels.
//
// http://facebook.github.io/graphql/draft/#sec-Objects
type Object struct {
	Name       string
	Interfaces []*Interface
	Fields     FieldDefinition
	Desc       string
	Directives common.DirectiveList
}

// Interface types represent a list of named fields and their arguments.
//
// GraphQL objects can then implement these interfaces which requires that the object type will
// define all fields defined by those interfaces.
//
// http://facebook.github.io/graphql/draft/#sec-Interfaces
type Interface struct {
	Name          string
	PossibleTypes []*Object
	Fields        FieldDefinition
	Desc          string
	Directives    common.DirectiveList
}

// Union types represent objects that could be one of a list of GraphQL object types, but provides no
// guaranteed fields between those types.
//
// They also differ from interfaces in that object types declare what interfaces they implement, but
// are not aware of what unions contain them.
//
// http://facebook.github.io/graphql/draft/#sec-Unions
type Union struct {
	Name          string
	PossibleTypes []*Object // NOTE: the spec refers to this as `UnionMemberTypes`.
	Desc          string
	Directives    common.DirectiveList
}

// Enum types describe a set of possible values.
//
// Like scalar types, Enum types also represent leaf values in a GraphQL type system.
//
// http://facebook.github.io/graphql/draft/#sec-Enums
type Enum struct {
	Name       string
	Values     []*EnumValuesDefinition // NOTE: the spec refers to this as `EnumValuesDefinition`.
	Desc       string
	Directives common.DirectiveList
}

// EnumValuesDefinition types are unique values that may be serialized as a string: the name of the
// represented value.
//
// http://facebook.github.io/graphql/draft/#EnumValueDefinition
type EnumValuesDefinition struct {
	Name       string
	Directives common.DirectiveList
	Desc       string
}

// InputObject types define a set of input fields; the input fields are either scalars, enums, or
// other input objects.
//
// This allows arguments to accept arbitrarily complex structs.
//
// http://facebook.github.io/graphql/draft/#sec-Input-Objects
type InputObject struct {
	Name       string
	Desc       string
	Values     types.InputValueList
	Directives common.DirectiveList
}

// Extension type defines a GraphQL type extension.
// Schemas, Objects, Inputs and Scalars can be extended.
//
// https://facebook.github.io/graphql/draft/#sec-Type-System-Extensions
type Extension struct {
	Type       NamedType
	Directives common.DirectiveList
}

// FieldDefinition is a list of an Object's Fields.
//
// http://facebook.github.io/graphql/draft/#FieldsDefinition
type FieldDefinition []*Field

// Get iterates over the field list, returning a pointer-to-Field when the field name matches the
// provided `name` argument.
// Returns nil when no field was found by that name.
func (l FieldDefinition) Get(name string) *Field {
	for _, f := range l {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// Names returns a string slice of the field names in the FieldDefinition.
func (l FieldDefinition) Names() []string {
	names := make([]string, len(l))
	for i, f := range l {
		names[i] = f.Name
	}
	return names
}

// http://facebook.github.io/graphql/draft/#sec-Type-System.Directives
type DirectiveDecl struct {
	Name string
	Desc string
	Locs []string
	Args types.InputValueList
}

func (*Scalar) Kind() string      { return "SCALAR" }
func (*Object) Kind() string      { return "OBJECT" }
func (*Interface) Kind() string   { return "INTERFACE" }
func (*Union) Kind() string       { return "UNION" }
func (*Enum) Kind() string        { return "ENUM" }
func (*InputObject) Kind() string { return "INPUT_OBJECT" }

func (t *Scalar) String() string      { return t.Name }
func (t *Object) String() string      { return t.Name }
func (t *Interface) String() string   { return t.Name }
func (t *Union) String() string       { return t.Name }
func (t *Enum) String() string        { return t.Name }
func (t *InputObject) String() string { return t.Name }

func (t *Scalar) TypeName() string      { return t.Name }
func (t *Object) TypeName() string      { return t.Name }
func (t *Interface) TypeName() string   { return t.Name }
func (t *Union) TypeName() string       { return t.Name }
func (t *Enum) TypeName() string        { return t.Name }
func (t *InputObject) TypeName() string { return t.Name }

func (t *Scalar) Description() string      { return t.Desc }
func (t *Object) Description() string      { return t.Desc }
func (t *Interface) Description() string   { return t.Desc }
func (t *Union) Description() string       { return t.Desc }
func (t *Enum) Description() string        { return t.Desc }
func (t *InputObject) Description() string { return t.Desc }

type Argument struct {
	Name  Ident
	Value Literal
}

type ArgumentList []Argument

type Literal interface {
	Value(vars map[string]interface{}) interface{}
	String() string
	Location() errors.Location
}

type BasicLit struct {
	Type rune
	Text string
	Loc  errors.Location
}

type Ident struct {
	Name string
	Loc  errors.Location
}

type InputValue struct {
	Name       Ident
	Type       Type
	Default    Literal
	Desc       string
	Directives DirectiveList
	Loc        errors.Location
	TypeLoc    errors.Location
}

type ArgumentsDefinition []*InputValue

func (l ArgumentsDefinition) Get(name string) *InputValue {
	for _, v := range l {
		if v.Name.Name == name {
			return v
		}
	}
	return nil
}

type Type interface {
	Kind() string
	String() string
}

type Directive struct {
	Name Ident
	Args ArgumentList
}

type DirectiveList []*Directive

func (l DirectiveList) Get(name string) *Directive {
	for _, d := range l {
		if d.Name.Name == name {
			return d
		}
	}
	return nil
}

// Field is a conceptual function which yields values.
// http://facebook.github.io/graphql/draft/#FieldDefinition
type Field struct {
	Name       string
	Args       ArgumentsDefinition
	Type       Type
	Directives DirectiveList
	Desc       string
}

func (f *Field) Description() *string {
	if f.Desc == "" {
		return nil
	}
	return &f.Desc
}
