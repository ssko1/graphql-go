package types

import "github.com/graph-gophers/graphql-go/errors"

type Document struct {
	Operations OperationList
	Fragments  FragmentList
}

type Operation struct {
	Type       OperationType
	Name       Ident
	Vars       ArgumentsDefinition
	Selections []Selection
	Directives DirectiveList
	Loc        errors.Location
}

type OperationType string

type InputValueList []*InputValue

func (l InputValueList) Get(name string) *InputValue {
	for _, v := range l {
		if v.Name.Name == name {
			return v
		}
	}
	return nil
}

type Selection interface {
	isSelection()
}

type InlineFragment struct {
	Fragment
	Directives DirectiveList
	Loc        errors.Location
}

type Fragment struct {
	On         TypeName
	Selections []Selection
}

type FragmentDecl struct {
	Fragment
	Name       Ident
	Directives DirectiveList
	Loc        errors.Location
}

type FragmentSpread struct {
	Name       Ident
	Directives DirectiveList
	Loc        errors.Location
}

func (QueryField) isSelection()     {}
func (InlineFragment) isSelection() {}
func (FragmentSpread) isSelection() {}

type OperationList []*Operation

func (l OperationList) Get(name string) *Operation {
	for _, f := range l {
		if f.Name.Name == name {
			return f
		}
	}
	return nil
}

type FragmentList []*FragmentDecl

func (l FragmentList) Get(name string) *FragmentDecl {
	for _, f := range l {
		if f.Name.Name == name {
			return f
		}
	}
	return nil
}
