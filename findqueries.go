package main

import (
	"go/parser"
	"go/token"
	"log"
)

type queryInfo struct{
	query string
}

func findQueries(dir string) map[string][]queryInfo {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		log.Fatalf("failed to parse dir %s: %v", dir, err)
	}

	queries := make(map[string][]queryInfo)

	for pkgName, pkg := range pkgs {
		queries[pkgName] = searchPackage(pkg)
	}
	
	return queries
}
