package common

import (
	"github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/types"
)

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

func (l InputValueList) Get(name string) *InputValue {
	for _, v := range l {
		if v.Name.Name == name {
			return v
		}
	}
	return nil
}

func ParseInputValue(l *Lexer) *InputValue {
	p := &InputValue{}
	p.Loc = l.Location()
	p.Desc = l.DescComment()
	p.Name = l.ConsumeIdentWithLoc()
	l.ConsumeToken(':')
	p.TypeLoc = l.Location()
	p.Type = ParseType(l)
	if l.Peek() == '=' {
		l.ConsumeToken('=')
		p.Default = ParseLiteral(l, true)
	}
	p.Directives = ParseDirectives(l)
	return p
}

func ParseInputValuePrime(l *Lexer) *types.InputValue {
	p := &types.InputValue{}
	p.Loc = l.Location()
	p.Desc = l.DescComment()
	p.Name = l.ConsumeIdentWithLocPrime()
	l.ConsumeToken(':')
	p.TypeLoc = l.Location()
	p.Type = ParseType(l)
	if l.Peek() == '=' {
		l.ConsumeToken('=')
		p.Default = ParseLiteral(l, true)
	}
	p.Directives = ParseDirectivesPrime(l)
	return p
}

type Argument struct {
	Name  Ident
	Value Literal
}

type ArgumentList []Argument

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

func ParseArguments(l *Lexer) ArgumentList {
	var args ArgumentList
	l.ConsumeToken('(')
	for l.Peek() != ')' {
		name := l.ConsumeIdentWithLoc()
		l.ConsumeToken(':')
		value := ParseLiteral(l, false)
		args = append(args, Argument{Name: name, Value: value})
	}
	l.ConsumeToken(')')
	return args
}

func ParseArgumentsPrime(l *Lexer) types.ArgumentList {
	var args types.ArgumentList
	l.ConsumeToken('(')
	for l.Peek() != ')' {
		name := l.ConsumeIdentWithLocPrime()
		l.ConsumeToken(':')
		value := ParseLiteral(l, false)
		// TODO this is not a pointer in the old version. Check for references to ArgumentList
		args = append(args, &types.Argument{
			Name:  name,
			Value: value,
		})
	}
	l.ConsumeToken(')')
	return args
}
