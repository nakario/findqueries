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

func resolve(query ssa.Value) ([]string, error) {
	switch q := query.(type) {
	case *ssa.Phi:
		ret := make([]string, 0)
		for _, edge := range q.Edges {
			resolved, err := resolve(edge)
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve a phi query")
			}
			ret = append(ret, resolved...)
		}
		return ret, nil
	case *ssa.Const:
		queryStr, _ := strconv.Unquote(q.Value.ExactString())
		return []string{queryStr}, nil
	case *ssa.BinOp:
		xs, err := resolve(q.X)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve a query: lhs of +")
		}
		ys, err := resolve(q.Y)
		ret := make([]string, 0, len(xs) * len(ys))
		for _, x := range xs {
			for _, y := range ys {
				ret = append(ret, x + y)
			}
		}
		return ret, errors.Wrap(err, "failed to resolve a query: rhs of +")
	case *ssa.Extract:
		switch tuple := q.Tuple.(type) {
		case *ssa.Call:
			callee := tuple.Call.StaticCallee()
			if callee == nil {
				return nil, errors.New("dynamic call is not supported")
			}
			return resolveFunc(callee, q.Index)
		case *ssa.TypeAssert:
			return nil, errors.Errorf("not implemented extract: %#v", tuple.X)
		case *ssa.Next:
			return nil, errors.Errorf("not implemented extract: %#v", tuple.Iter)
		case *ssa.UnOp:
			return nil, errors.Errorf("not implemented extract: %#v", tuple.X)
		case *ssa.Lookup:
			return nil, errors.Errorf("not implemented extract: %#v", tuple.X)
		default:
			return nil, errors.New("unexpected extract")
		}
	case *ssa.Call:
		callee := q.Call.StaticCallee()
		if callee == nil {
			return nil, errors.New("dynamic call is not supported")
		}
		return resolveFunc(callee, 0)
	}
	return nil, errors.New("failed to resolve a query: unsupported value")
}

func resolveFunc(fn *ssa.Function, index int) ([]string, error) {
	queries := make([]string, 0)
	for _, block := range fn.Blocks {
		if len(block.Instrs) == 0 {
			continue
		}
		switch last := block.Instrs[len(block.Instrs)-1].(type) {
		case *ssa.If:
			return nil, errors.Errorf("not implemented instr: %#v", last)
		case *ssa.Jump:
			return nil, errors.Errorf("not implemented instr: %#v", last)
		case *ssa.Return:
			queryStr, err := resolve(last.Results[index])
			if err != nil {
				return nil, errors.WithMessage(err, "failed to resolve return value")
			}
			queries = append(queries, queryStr...)
		case *ssa.Panic:
			return nil, errors.Errorf("not implemented instr: %#v", last)
		default:
			return nil, errors.New("unexpected last instruction in a function block")
		}
	}
	return queries, nil
}

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
						fmt.Fprintln(os.Stderr, pkg.Prog.Fset.Position(site.Pos()))
						return nil, nil, errors.Errorf("found an unexpected call")
					}
					args := site.Common().Args
					if site.Common().Signature().Recv() != nil {
						args = args[1:]
					}
					query := args[pos]
					possibleQueries, err := resolve(query)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						unresolved = append(unresolved, qi)
					} else {
						for _, pq := range possibleQueries {
							cp := qi
							cp.Query = pq
							queries = append(queries, cp)
						}
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
