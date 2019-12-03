package findqueries

import (
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

type Result struct {
	Queries    map[string][]queryInfo `json:"queries"`
	Unresolved map[string][]queryInfo `json:"unresolved"`
	Calls      map[string][]call      `json:"calls"`
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
	qs, us, cs, err := findQueries(pkg, queriers, builders, er)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to analyze a package")
	}
	queries := make(map[string][]queryInfo)
	unresolved := make(map[string][]queryInfo)
	calls := make(map[string][]call)
	queries[pkg.Pkg.Name()] = qs
	unresolved[pkg.Pkg.Name()] = us
	calls[pkg.Pkg.Name()] = cs
	return &Result{
		Queries:    queries,
		Unresolved: unresolved,
		Calls:      calls,
	}, nil
}
