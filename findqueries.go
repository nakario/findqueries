package main

import (
	"go/parser"
	"go/token"

	"github.com/pkg/errors"
)

type queryInfo struct{
	query string
}

func findQueries(dir string) (map[string][]queryInfo, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to parse dir %s", dir)
	}

	queries := make(map[string][]queryInfo)

	for pkgName, pkg := range pkgs {
		qs, err := searchPackage(pkg, dir, fset)
		if err != nil {
			return nil, err
		}
		queries[pkgName] = qs
	}
	
	return queries, nil
}
