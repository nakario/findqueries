// The findqueries command runs the findqueries analyzer.
package main

import (
	"github.com/nakario/findqueries"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(findqueries.Analyzer) }
