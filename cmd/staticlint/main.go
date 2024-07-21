// Static analytic service. Include static analytic packages:
// - golang.org/x/tools/go/analysis/passes
// - all SA classes staticcheck.io
// - Go-critic and nilerr linters
// - OsExitAnalyzer to check os.Exit calls in main packages
package main

import (
	"fmt"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"strings"

	"github.com/LilLebowski/shortener/pkg/osexitanalyzer"
	goc "github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/gostaticanalysis/nilerr"
)

// StaticCheckConfig describe config data
type StaticCheckConfig struct {
	Staticcheck []string
	Stylecheck  []string
}

func main() {
	cfg := StaticCheckConfig{
		Staticcheck: []string{"SA"},
		Stylecheck:  []string{"ST1000", "ST1005"},
	}

	checks := prepareChecks(cfg)

	fmt.Println("Run static checks:\n", checks)

	multichecker.Main(
		checks...,
	)
}

func prepareChecks(cfg StaticCheckConfig) []*analysis.Analyzer {
	checks := []*analysis.Analyzer{
		osexitanalyzer.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		goc.Analyzer,
		nilerr.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		for _, sc := range cfg.Staticcheck {
			if strings.HasPrefix(v.Name, sc) {
				checks = append(checks, v)
			}
		}
	}

	for _, v := range stylecheck.Analyzers {
		for _, sc := range cfg.Stylecheck {
			if strings.HasPrefix(v.Name, sc) {
				checks = append(checks, v)
			}
		}
	}
	return checks
}
