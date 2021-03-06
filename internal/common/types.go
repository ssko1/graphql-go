package common

import (
	"github.com/graph-gophers/graphql-go/types"
)

func ParseType(l *Lexer) types.Type {
	t := parseNullType(l)
	if l.Peek() == '!' {
		l.ConsumeToken('!')
		return &types.NonNull{OfType: t}
	}
	return t
}

func parseNullType(l *Lexer) types.Type {
	if l.Peek() == '[' {
		l.ConsumeToken('[')
		ofType := ParseType(l)
		l.ConsumeToken(']')
		return &types.List{OfType: ofType}
	}

	return &types.TypeName{Ident: l.ConsumeIdentWithLoc()}
}
