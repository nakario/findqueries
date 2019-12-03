package findqueries

import (
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

type Result struct {
	Queries map[string][]queryInfo `json:"queries"`
	Calls   map[string][]call      `json:"calls"`
}

type queryInfo struct {
	Query  string `json:"query"`
	Caller string `json:"caller"`
	Expr   string `json:"expr"`
	Pos    string `json:"pos"`
	err    error
}

type call struct {
	Caller string `json:"caller"`
	Callee string `json:"callee"`
}

func analyzePackage(pkg *ssa.Package, er ExprResolver, queriers []querierInfo, builders []builderInfo) (*Result, error) {
	if pkg == nil {
		return nil, errors.New("nil package found")
	}
	qs, cs, err := findQueries(pkg, queriers, builders, er)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to analyze a package")
	}
	queries := make(map[string][]queryInfo)
	calls := make(map[string][]call)
	queries[pkg.Pkg.Name()] = qs
	calls[pkg.Pkg.Name()] = cs
	return &Result{
		Queries: queries,
		Calls:   calls,
	}, nil
}
