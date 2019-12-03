// The findqueries package defines an Analyzer that find SQL queries
// and all functions calling such queries inside themselves.

package findqueries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `find SQL queries and functions calling them internally`

var Analyzer = &analysis.Analyzer{
	Name: "findqueries",
	Doc: Doc,
	Run: run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		buildssa.Analyzer,
	},
	ResultType: reflect.TypeOf((*Result)(nil)),
	FactTypes: []analysis.Fact{},
}

var (
	queriersInfoPath string // -queriers flag
	buildersInfoPath string // -builders flag
)

func init() {
	Analyzer.Flags.StringVar(&queriersInfoPath, "queriers", defaultQuerierInfoPath(), "path to queriers.json")
	Analyzer.Flags.StringVar(&buildersInfoPath, "builders", defaultBuilderInfoPath(), "path to builders.json")
}

type ExprResolver interface {
	ResolveFrom(pos token.Pos) ast.Expr
}

type exprResolver struct {
	pos2expr map[token.Pos]ast.Expr
}

func (er exprResolver) ResolveFrom(pos token.Pos) ast.Expr {
	return er.pos2expr[pos]
}

func NewExprResolver(inspctr *inspector.Inspector) ExprResolver {
	pos2expr := make(map[token.Pos]ast.Expr)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.GoStmt)(nil),
		(*ast.DeferStmt)(nil),
	}
	inspctr.Preorder(nodeFilter, func(n ast.Node) {
		switch expr := n.(type) {
		case *ast.CallExpr:
			pos2expr[expr.Lparen] = expr
		case *ast.GoStmt:
			pos2expr[expr.Go] = expr.Call
		case *ast.DeferStmt:
			pos2expr[expr.Defer] = expr.Call
		}
	})
	return exprResolver{pos2expr: pos2expr}
}

func run(pass *analysis.Pass) (interface{}, error) {
	withMessage := func(err error) error {
		return errors.WithMessagef(err, "failed to analyze pass %s", pass.String())
	}
	queryers, err := loadQuerierInfo(queriersInfoPath)
	if err != nil {
		return nil, withMessage(err)
	}

	builders, err := loadBuilderInfo(buildersInfoPath)
	if err != nil {
		return nil, withMessage(err)
	}
	ssa := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	inspctr := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	er := NewExprResolver(inspctr)

	result, err := analyzePackage(ssa.Pkg, er, queryers, builders)
	if err != nil {
		return nil, withMessage(err)
	}

	b, err := json.Marshal(result)
	if err != nil {
		return nil, withMessage(errors.Wrap(err, "failed to jsonify result"))
	}
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, b); err != nil {
		return nil, withMessage(errors.Wrap(err, "failed to compact json"))
	}
	fmt.Println(buf)

	return result, nil
}
