package main

import (
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

type queryInfo struct{
	Query string `json:"query"`
	Caller string `json:"caller"`
	Expr string `json:"expr"`
	Pos string `json:"pos"`
}

func analyze(dir string, queryers []queryerInfo) (map[string][]queryInfo, error) {
	conf := &packages.Config{
		Dir: dir,
		Mode: packages.NeedName |
		packages.NeedFiles |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(conf)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to parse dir %s", dir)
	}

	queries := make(map[string][]queryInfo)

	for _, pkg := range pkgs {
		qs, err := findQueries(pkg, queryers)
		if err != nil {
			return nil, err
		}
		queries[pkg.Name] = qs
	}
	
	return queries, nil
}
