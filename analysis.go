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
				switch expr := n.(type) {
				case *ast.CallExpr:
					pos2expr[expr.Lparen] = expr
				case *ast.GoStmt:
					pos2expr[expr.Go] = expr.Call
				case *ast.DeferStmt:
					pos2expr[expr.Defer] = expr.Call
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
