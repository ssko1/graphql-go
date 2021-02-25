package parser

import (
	"github.com/graph-gophers/graphql-go/internal/common"
	"github.com/graph-gophers/graphql-go/internal/schema"
	"github.com/graph-gophers/graphql-go/types"
)

// the code in this file is needed as a bridge into existing Twitch code. We should delete it soon.

// ParseDirectory is unimplemented. This code is copied-ish from twitch's linter.
// I have not yet explored whether we should carry the `grouByFile` approach
// forward to this new codebase or if there is a better way. The requirements of ParseDirectory are to:
// - recursively descend and collect all .graphql files
// - validate that these files compose to form a coherent and unambiguous schema (no dupes)
// - (new feature): build a map of source files to files in the parse tree so that error messages are helpful to coders
// func ParseDirectory(path string) (*types.Schema, error) {
// 	s := &types.Schema{}
// 	err := filepath.Walk(path, func(name string, f os.FileInfo, err error) error {
// 		if err != nil || f.IsDir() || !strings.HasSuffix(name, ".graphql") {
// 			// Early-return if:
// 			// - already have an error
// 			// - the current file is a Directory
// 			// - the current file does not have the .grapqhl file extension.
// 			return err
// 		}

// 		l, err := common.NewFileScanner(name)
// 		if err != nil {
// 			return err
// 		}

// 		schema.ParseSchema(s, l)

// 		return l.Close()
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// Resolve is an unfortunate name. This block of code is really checking that
// all the files we just slurped up create an unambiguous schema.
// 	// if err := s.Resolve(); err != nil {
// 	// 	return nil, err
// 	// }

// 	return s, nil
// 	// return s.groupByFile()
// }

func Parse(schemaString string) (*types.Schema, error) {
	s := &types.Schema{}
	l := common.NewLexer(schemaString, false)
	schema.ParseSchema(s, l)
	// TODO where are the errors
	return s, nil
}
