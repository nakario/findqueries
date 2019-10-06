package main

import (
	"go/ast"
	"go/token"
	"go/types"

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

	pkg2pos2expr := make(map[*types.Package]map[token.Pos]*ast.CallExpr)
	for _, p := range pkgs {
		pos2expr := make(map[token.Pos]*ast.CallExpr)
		for _, file := range p.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				if ce, ok := n.(*ast.CallExpr); ok {
					pos2expr[ce.Lparen] = ce
				}
				return true
			})
		}
		pkg2pos2expr[p.Types] = pos2expr
	}

	queries := make(map[string][]queryInfo)
	calls := make(map[string][]call)

	_, ssaPkgs := ssautil.Packages(pkgs, 0)
	for _, pkg := range ssaPkgs {
		if pkg == nil { continue }
		qs, cs, err := findCalls(pkg, queryers, pkg2pos2expr[pkg.Pkg])
		if err != nil {
			return nil, err
		}
		queries[pkg.Pkg.Name()] = qs
		calls[pkg.Pkg.Name()] = cs
	}
	
	return &result{queries, calls}, nil
}
