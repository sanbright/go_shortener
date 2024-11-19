// Package noexitanalyzer пакет анализатора
package noexitanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// NoExitAnalyzer проверяет использование os.Exit в main.
var NoExitAnalyzer = &analysis.Analyzer{
	Name: "noexit",
	Doc:  "checks for os.Exit calls within main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if pass.Pkg.Name() == "main" && pass.Fset.File(file.Pos()).Name() == "main.go" {
			ast.Inspect(file, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
					ast.Inspect(fn.Body, func(n ast.Node) bool {
						if call, ok := n.(*ast.CallExpr); ok {
							if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
								switch pkgIdent, ok := sel.X.(*ast.Ident); {
								case pkgIdent.Name == "os" && ok && sel.Sel.Name == "Exit":
									pass.Reportf(call.Pos(), "os.Exit called within main function2d")
								}
							}
						}
						return true
					})
				}
				return true
			})
		}
	}

	return nil, nil
}
