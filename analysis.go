package findqueries

import (
	"go/types"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

type Result struct {
	Name       string      `json:"name"`
	Queries    []queryInfo `json:"queries"`
	Calls      []call      `json:"calls"`
	unresolved []queryInfo
}

type queryInfo struct {
	Query     string `json:"query"`
	Caller    string `json:"caller"`
	Expr      string `json:"expr"`
	Pos       string `json:"pos"`
	calleeObj types.Object
	err       error
}

type call struct {
	Caller string `json:"caller"`
	Callee string `json:"callee"`
}

func analyzePackage(pkg *ssa.Package, er ExprResolver, queriers []querierInfo, builders []builderInfo) (*Result, error) {
	if pkg == nil {
		return nil, errors.New("nil package found")
	}
	queries, unresolved, calls, err := findQueries(pkg, queriers, builders, er)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to analyze a package")
	}
	return &Result{
		Name:       pkg.Pkg.Name(),
		Queries:    queries,
		unresolved: unresolved,
		Calls:      calls,
	}, nil
}
