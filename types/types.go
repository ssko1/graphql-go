package types

import (
	"strings"

	"github.com/graph-gophers/graphql-go/errors"
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
	EntryPoints map[string]NamedType

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

	EntryPointNames map[string]string
	Objects         []*Object
	Unions          []*Union
	Enums           []*Enum
	Extensions      []*Extension
}

func (s *Schema) Resolve(name string) Type {
	return s.Types[name]
}

// Extension type defines a GraphQL type extension.
// Schemas, Objects, Inputs and Scalars can be extended.
//
// https://facebook.github.io/graphql/draft/#sec-Type-System-Extensions
type Extension struct {
	Type       NamedType
	Directives DirectiveList
}

// EnumValuesDefinition types are unique values that may be serialized as a string: the name of the
// represented value.
//
// http://facebook.github.io/graphql/draft/#EnumValueDefinition
type EnumValuesDefinition struct {
	Name       string
	Directives DirectiveList
	Desc       string
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
	Directives DirectiveList
}

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

// http://facebook.github.io/graphql/draft/#InputValueDefinition
type InputValue struct {
	Name       Ident
	Type       Type
	Default    Literal
	Desc       string
	Directives DirectiveList
	Loc        errors.Location
	TypeLoc    errors.Location
}

// Replaced with ArgumentsDefinition
// type InputValueList []*InputValue

type DirectiveDecl struct {
	Name string
	Desc string
	Locs []string
	Args ArgumentsDefinition
}

// NamedType represents a type with a name.
//
// http://facebook.github.io/graphql/draft/#NamedType
type NamedType interface {
	Type
	TypeName() string
	Description() string
}

type Object struct {
	Name       string
	Interfaces []*Interface
	Fields     FieldDefinition
	Desc       string
	Directives DirectiveList

	InterfaceNames []string
}

// FieldDefinition is a list of an Object's Fields.
//
// http://facebook.github.io/graphql/draft/#FieldsDefinition
type FieldDefinition []*Field

func (l FieldDefinition) Get(name string) *Field {
	for _, f := range l {
		if f.Name.Name == name {
			return f
		}
	}
	return nil
}

func (l FieldDefinition) Names() []string {
	names := make([]string, len(l))
	for i, f := range l {
		names[i] = f.Name.Name
	}
	return names
}

type ArgumentsDefinition []*InputValue

func (a ArgumentsDefinition) Get(name string) *InputValue {
	for _, inputValue := range a {
		if inputValue.Name.Name == name {
			return inputValue
		}
	}
	return nil
}

type Field struct {
	Alias           Ident
	Name            Ident
	Arguments       ArgumentsDefinition
	Type            Type
	Directives      DirectiveList
	Desc            string
	Selections      []Selection
	SelectionSetLoc errors.Location
}

// QueryField represents a field used in a query. It is the realization of a `Field`
type QueryField struct {
	Alias           Ident
	Name            Ident
	Arguments       ArgumentList
	Directives      DirectiveList
	Selections      []Selection
	SelectionSetLoc errors.Location
}

// -type Field struct {
// 	-	Alias           common.Ident
// 	-	Name            common.Ident
// 	-	Arguments       common.ArgumentList
// 	-	Directives      common.DirectiveList
// 	-	Selections      []Selection
// 	-	SelectionSetLoc errors.Location
// 	-}

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
	Directives    DirectiveList
}

type Directive struct {
	Name Ident
	Args ArgumentList
}

type Argument struct {
	Name  Ident
	Value Literal
}

type ArgumentList []*Argument

func (l ArgumentList) Get(name string) (Literal, bool) {
	for _, arg := range l {
		if arg.Name.Name == name {
			return arg.Value, true
		}
	}
	return nil, false
}

func (l ArgumentList) MustGet(name string) Literal {
	value, ok := l.Get(name)
	if !ok {
		panic("argument not found")
	}
	return value
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

type Ident struct {
	Name string
	Loc  errors.Location
}

type Type interface {
	Kind() string
	String() string
}

type List struct {
	OfType Type
}

type NonNull struct {
	OfType Type
}

type TypeName struct {
	Ident
}

func (*List) Kind() string     { return "LIST" }
func (*NonNull) Kind() string  { return "NON_NULL" }
func (*TypeName) Kind() string { panic("TypeName needs to be resolved to actual type") }

func (t *List) String() string    { return "[" + t.OfType.String() + "]" }
func (t *NonNull) String() string { return t.OfType.String() + "!" }
func (*TypeName) String() string  { panic("TypeName needs to be resolved to actual type") }

type Resolver func(name string) Type

func ResolveType(t Type, resolver Resolver) (Type, *errors.QueryError) {
	switch t := t.(type) {
	case *List:
		ofType, err := ResolveType(t.OfType, resolver)
		if err != nil {
			return nil, err
		}
		return &List{OfType: ofType}, nil
	case *NonNull:
		ofType, err := ResolveType(t.OfType, resolver)
		if err != nil {
			return nil, err
		}
		return &NonNull{OfType: ofType}, nil
	case *TypeName:
		refT := resolver(t.Name)
		if refT == nil {
			err := errors.Errorf("Unknown type %q.", t.Name)
			err.Rule = "KnownTypeNames"
			err.Locations = []errors.Location{t.Loc}
			return nil, err
		}
		return refT, nil
	default:
		return t, nil
	}
}

// TODO I don't think this belongs in the types package
func (lit *BasicLit) Value(vars map[string]interface{}) interface{} {
	// switch lit.Type {
	// case scanner.Int:
	// 	value, err := strconv.ParseInt(lit.Text, 10, 32)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return int32(value)

	// case scanner.Float:
	// 	value, err := strconv.ParseFloat(lit.Text, 64)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return value

	// case scanner.String:
	// 	value, err := strconv.Unquote(lit.Text)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return value

	// case scanner.Ident:
	// 	switch lit.Text {
	// 	case "true":
	// 		return true
	// 	case "false":
	// 		return false
	// 	default:
	// 		return lit.Text
	// 	}

	panic("literal value not implemented in static context")
}

func (lit *BasicLit) String() string {
	return lit.Text
}

func (lit *BasicLit) Location() errors.Location {
	return lit.Loc
}

type ListLit struct {
	Entries []Literal
	Loc     errors.Location
}

func (lit *ListLit) Value(vars map[string]interface{}) interface{} {
	entries := make([]interface{}, len(lit.Entries))
	for i, entry := range lit.Entries {
		entries[i] = entry.Value(vars)
	}
	return entries
}

func (lit *ListLit) String() string {
	entries := make([]string, len(lit.Entries))
	for i, entry := range lit.Entries {
		entries[i] = entry.String()
	}
	return "[" + strings.Join(entries, ", ") + "]"
}

func (lit *ListLit) Location() errors.Location {
	return lit.Loc
}

type ObjectLit struct {
	Fields []*ObjectLitField
	Loc    errors.Location
}

type ObjectLitField struct {
	Name  Ident
	Value Literal
}

func (lit *ObjectLit) Value(vars map[string]interface{}) interface{} {
	fields := make(map[string]interface{}, len(lit.Fields))
	for _, f := range lit.Fields {
		fields[f.Name.Name] = f.Value.Value(vars)
	}
	return fields
}

func (lit *ObjectLit) String() string {
	entries := make([]string, 0, len(lit.Fields))
	for _, f := range lit.Fields {
		entries = append(entries, f.Name.Name+": "+f.Value.String())
	}
	return "{" + strings.Join(entries, ", ") + "}"
}

func (lit *ObjectLit) Location() errors.Location {
	return lit.Loc
}

type NullLit struct {
	Loc errors.Location
}

func (lit *NullLit) Value(vars map[string]interface{}) interface{} {
	return nil
}

func (lit *NullLit) String() string {
	return "null"
}

func (lit *NullLit) Location() errors.Location {
	return lit.Loc
}

type Variable struct {
	Name string
	Loc  errors.Location
}

func (v Variable) Value(vars map[string]interface{}) interface{} {
	return vars[v.Name]
}

func (v Variable) String() string {
	return "$" + v.Name
}

func (v *Variable) Location() errors.Location {
	return v.Loc
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
	Directives DirectiveList
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
	Directives    DirectiveList
	TypeNames     []string
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
	Values     ArgumentsDefinition
	Directives DirectiveList
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
