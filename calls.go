package main

import (
	"sort"

	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/ssa"
)

func findCalls(pkg *ssa.Package) []call {
	er2ees := make(map[string][]string)
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

	calls := make([]call, 0, numEdges)
	for _, er := range ers {
		for _, ee := range er2ees[er] {
			calls = append(calls, call{er, ee})
		}
	}

	return calls
}
