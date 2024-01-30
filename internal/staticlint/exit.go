package staticlint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// ExitCheckAnalyzer - структура для кастомной проверки внутри mutichecker
var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check os.Exit calls inside main func of main pkg",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// нам нужны только файлы с package main
		if file.Name.Name != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			f, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// нам нужна только функция main
			if f.Name.Name != "main" {
				return false
			}

			ast.Inspect(f, func(node ast.Node) bool {
				if fun, ok := node.(*ast.SelectorExpr); ok {
					if x, ok := fun.X.(*ast.Ident); ok {
						if x.Name == "os" && fun.Sel.Name == "Exit" {
							pass.Reportf(node.Pos(), "os.Exit call inside main")
						}
					}
					return false
				}
				return true
			})
			return true
		})
	}

	return nil, nil
}
