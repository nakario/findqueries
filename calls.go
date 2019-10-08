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

type queryResolver struct {
	pkg *types.Package
	buildersMap map[string]builderInfo
}

func newQueryResolver(pkg *types.Package, builders []builderInfo) *queryResolver {
	buildersMap := make(map[string]builderInfo)
	for _, bi := range builders {
		buildersMap[bi.FullName] = bi
	}
	return &queryResolver{pkg, buildersMap}
}

func (qr *queryResolver) resolve(query ssa.Value) ([]string, error) {
	switch q := query.(type) {
	case *ssa.Phi:
		ret := make([]string, 0)
		for _, edge := range q.Edges {
			resolved, err := qr.resolve(edge)
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve a phi edge")
			}
			ret = append(ret, resolved...)
		}
		if len(ret) == 0 {
			return nil, errors.New("there was no possible edge in a phi node")
		}
		return ret, nil
	case *ssa.Const:
		queryStr, _ := strconv.Unquote(q.Value.ExactString())
		return []string{queryStr}, nil
	case *ssa.BinOp:
		xs, err := qr.resolve(q.X)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve lhs of +")
		}
		if len(xs) == 0 {
			return nil, errors.New("couldn't find any queries from lhs of +")
		}
		ys, err := qr.resolve(q.Y)
		ret := make([]string, 0, len(xs) * len(ys))
		for _, x := range xs {
			for _, y := range ys {
				ret = append(ret, x + y)
			}
		}
		if len(ys) == 0 {
			return nil, errors.New("couldn't find any queries from rhs of +")
		}
		return ret, errors.Wrap(err, "failed to resolve rhs of +")
	case *ssa.Extract:
		switch tuple := q.Tuple.(type) {
		case *ssa.Call:
			callee := tuple.Call.StaticCallee()
			if callee == nil {
				return nil, errors.New("dynamic call is not supported")
			}
			if bi, ok := qr.buildersMap[callee.RelString(qr.pkg)]; ok {
				return qr.resolve(tuple.Call.Args[bi.ArgIndex])
			}
			return qr.resolveFunc(callee, q.Index)
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
		if bi, ok := qr.buildersMap[callee.RelString(qr.pkg)]; ok {
			return qr.resolve(q.Call.Args[bi.ArgIndex])
		}
		return qr.resolveFunc(callee, 0)
	}
	return nil, errors.New("failed to resolve a query: unsupported value")
}

func (qr *queryResolver) resolveFunc(fn *ssa.Function, index int) ([]string, error) {
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
			queryStr, err := qr.resolve(last.Results[index])
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
	if len(queries) == 0 {
		return nil, errors.New("couldn't find any queries in a function")
	}
	return queries, nil
}

func findCalls(pkg *ssa.Package, queryers []queryerInfo, builders []builderInfo, pos2expr map[token.Pos]*ast.CallExpr) ([]queryInfo, []call, error) {
	queryersMap := make(map[string]int)
	for _, qi := range queryers {
		queryersMap[qi.FullName] = qi.QueryPos
	}
	qr := newQueryResolver(pkg.Pkg, builders)
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
					possibleQueries, err := qr.resolve(query)
					if err != nil {
						qi.err = err
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
		fmt.Fprintln(os.Stderr, qi.err)
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
