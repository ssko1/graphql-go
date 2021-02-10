package types

import (
	"github.com/graph-gophers/graphql-go/errors"
	"strings"
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

type InputValueList []*InputValue

type DirectiveDecl struct {
	Name string
	Desc string
	Locs []string
	Args InputValueList
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
}

// FieldDefinition is a list of an Object's Fields.
//
// http://facebook.github.io/graphql/draft/#FieldsDefinition
type FieldDefinition []*Field

type ArgumentsDefinition []*InputValue

// Field is a conceptual function which yields values.
// http://facebook.github.io/graphql/draft/#FieldDefinition
type Field struct {
	Name       string
	Args       ArgumentsDefinition
	Type       Type
	Directives DirectiveList
	Desc       string
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

type ArgumentList []Argument

type DirectiveList []*Directive

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
