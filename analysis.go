package main

import (
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa/ssautil"
)

type result struct{
	Queries map[string][]queryInfo `json:"queries"`
	Calls map[string][]call `json:"calls"`
}

type queryInfo struct{
	Query string `json:"query"`
	Caller string `json:"caller"`
	Expr string `json:"expr"`
	Pos string `json:"pos"`
}

type call struct{
	Caller string `json:"caller"`
	Callee string `json:"callee"`
}

func analyze(dir string, queryers []queryerInfo) (*result, error) {
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

	_, ssaPkgs := ssautil.Packages(pkgs, 0)
	calls := make(map[string][]call)
	for _, pkg := range ssaPkgs {
		if pkg == nil { continue }
		cs := findCalls(pkg)
		calls[pkg.Pkg.Name()] = cs
	}
	
	return &result{queries, calls}, nil
}
