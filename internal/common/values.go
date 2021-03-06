package common

import (
	"github.com/graph-gophers/graphql-go/types"
)

func ParseInputValue(l *Lexer) *types.InputValue {
	p := &types.InputValue{}
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

func ParseArguments(l *Lexer) types.ArgumentList {
	var args types.ArgumentList
	l.ConsumeToken('(')
	for l.Peek() != ')' {
		name := l.ConsumeIdentWithLoc()
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

func ParseArgumentsDef(l *Lexer) types.ArgumentsDefinition {
	var args types.ArgumentsDefinition
	l.ConsumeToken('(')
	for l.Peek() != ')' {
		l.ConsumeToken(':')
		value := ParseInputValue(l)
		args = append(args, value)
	}
	l.ConsumeToken(')')
	return args
}
