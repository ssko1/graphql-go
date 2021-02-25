package parser

import (
	"github.com/graph-gophers/graphql-go/internal/common"
	"github.com/graph-gophers/graphql-go/internal/schema"
	"github.com/graph-gophers/graphql-go/types"
	"os"
	"path/filepath"
	"strings"
)

// the code in this file is needed as a bridge into existing Twitch code. We should delete it soon.

func ParseDirectory(path string) (*types.Schema, error) {
	s := &types.Schema{}
	err := filepath.Walk(path, func(name string, f os.FileInfo, err error) error {
		if err != nil || f.IsDir() || !strings.HasSuffix(name, ".graphql") {
			// Early-return if:
			// - already have an error
			// - the current file is a Directory
			// - the current file does not have the .grapqhl file extension.
			return err
		}

		l, err := common.NewFileScanner(name)
		if err != nil {
			return err
		}

		// TODO this ain't right
		schema.ParseSchema(s, l)

		if err := s.parse(l); err != nil {
			return err
		}

		return l.Close()
	})
	if err != nil {
		return nil, err
	}

	if err := s.Resolve(); err != nil {
		return nil, err
	}

	return s, nil
	// return s.groupByFile()
}

func Parse(schemaString string) (*types.Schema, error) {
}
