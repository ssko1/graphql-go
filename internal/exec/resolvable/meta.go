package resolvable

import (
	"fmt"
	"reflect"

	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/graph-gophers/graphql-go/types"
)

// Meta defines the details of the metadata schema for introspection.
type Meta struct {
	FieldSchema   Field
	FieldType     Field
	FieldTypename Field
	Schema        *Object
	Type          *Object
}

func newMeta(s *types.Schema) *Meta {
	var err error
	b := newBuilder(s)

	metaSchema := s.Types["__Schema"].(*types.Object)
	so, err := b.makeObjectExec(metaSchema.Name, metaSchema.Fields, nil, false, reflect.TypeOf(&introspection.Schema{}))
	if err != nil {
		panic(err)
	}

	metaType := s.Types["__Type"].(*types.Object)
	t, err := b.makeObjectExec(metaType.Name, metaType.Fields, nil, false, reflect.TypeOf(&introspection.Type{}))
	if err != nil {
		panic(err)
	}

	if err := b.finish(); err != nil {
		panic(err)
	}

	fieldTypename := Field{
		Field: types.Field{
			Name: types.Ident{
				Name: "__typename",
			},
			Type: &types.NonNull{OfType: s.Types["String"]},
		},
		TraceLabel: fmt.Sprintf("GraphQL field: __typename"),
	}

	fieldSchema := Field{
		Field: types.Field{
			Name: types.Ident{
				Name: "__schema",
			},
			Type: s.Types["__Schema"],
		},
		TraceLabel: fmt.Sprintf("GraphQL field: __schema"),
	}

	fieldType := Field{
		Field: types.Field{
			Name: types.Ident{
				Name: "__type",
			},
			Type: s.Types["__Type"],
		},
		TraceLabel: fmt.Sprintf("GraphQL field: __type"),
	}

	return &Meta{
		FieldSchema:   fieldSchema,
		FieldTypename: fieldTypename,
		FieldType:     fieldType,
		Schema:        so,
		Type:          t,
	}
}
