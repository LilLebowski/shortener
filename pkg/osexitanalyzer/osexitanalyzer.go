// Package osexitanalyzer search call os.Exit in main packages and report position
// Implement analysis.Analyzer type interface for multi-check
package osexitanalyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// Analyzer describes os.Exit analysis function and its options
var Analyzer = &analysis.Analyzer{
	Name: "osexitanalyzer",
	Doc:  "Don't allow os.Exit in main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				if x.Name.Name != "main" {
					return false
				}
			case *ast.SelectorExpr:
				if x.Sel.Name == "Exit" {
					pass.Reportf(x.Pos(), "expression has os.Exit call in main package")

				}
			}
			return true
		})
	}
	return nil, nil
}
