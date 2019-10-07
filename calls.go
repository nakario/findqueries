package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/ssa"
)

func findCalls(pkg *ssa.Package, queryers []queryerInfo, pos2expr map[token.Pos]*ast.CallExpr) ([]queryInfo, []call, error) {
	queryersMap := make(map[string]int)
	for _, qi := range queryers {
		queryersMap[qi.FullName] = qi.QueryPos
	}
	queries := make([]queryInfo, 0)
	unresolved := make([]queryInfo, 0)
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

			if pos, ok := queryersMap[callee]; ok {
				site := e.Site
				if site != nil {
					if site.Pos() == token.NoPos {
						return nil, nil, errors.New("unexpectedly a call doesn't have its position")
					}
					qi := queryInfo{Caller: caller, Pos: pkg.Prog.Fset.Position(site.Pos()).String()}
					if expr, ok := pos2expr[site.Pos()]; ok {
						qi.Expr = types.ExprString(expr)
					} else {
						return nil, nil, errors.New("found an unexpected call")
					}
					args := site.Common().Args
					if site.Common().Signature().Recv() != nil {
						args = args[1:]
					}
					query := args[pos]
					if query, ok := query.(*ssa.Const); ok {
						query, _ := strconv.Unquote(query.Value.ExactString())
						qi.Query = query
						queries = append(queries, qi)
					} else {
						unresolved = append(unresolved, qi)
					}
				} else {
					fmt.Println("unexpected synthetic or intrinsic call")
				}
			}
			er2ees[caller] = append(er2ees[caller], callee)
		}
	}

	if len(unresolved) > 0 {
		fmt.Fprintln(os.Stderr, "UNRESOLVED")
	}
	for _, qi := range unresolved {
		fmt.Fprintln(os.Stderr, qi.Pos, qi.Expr)
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

	return queries, calls, nil
}