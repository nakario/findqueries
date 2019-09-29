package main

import (
	"sort"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/callgraph/cha"
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
		er2ees := make(map[string][]string)
		if pkg == nil { continue }
		pkg.Build()
		cg := cha.CallGraph(pkg.Prog)
		for fn, node := range cg.Nodes {
			if fn == nil || node == nil { continue }
			if node.Func.Package() != pkg {
				continue
			}
			for _, e := range node.Out {
				caller := e.Caller.Func.RelString(pkg.Pkg)
				callee := e.Callee.Func.RelString(pkg.Pkg)
				er2ees[caller] = append(er2ees[caller], callee)
			}
		}

		ers := make([]string, 0, len(er2ees))
		numEdges := 0
		for er, ees := range er2ees {
			ers = append(ers, er)
			sort.Strings(ees)
			numEdges += len(ees)
		}
	
		cs := make([]call, 0, numEdges)
		for _, er := range ers {
			for _, ee := range er2ees[er] {
				cs = append(cs, call{er, ee})
			}
		}
		calls[pkg.Pkg.Name()] = cs
	}
	
	return &result{queries, calls}, nil
}
