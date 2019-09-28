package main

import (
	"go/parser"
	"go/token"
	"log"
)

func findQueries(dir string) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		log.Fatalf("failed to parse dir %s: %v", dir, err)
	}

	for pkgName, pkg := range pkgs {
		log.Println(pkgName, ":", pkg)
	}
}
