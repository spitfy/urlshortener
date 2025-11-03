package linter

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "check_panic",
	Doc:      "Checks for prohibited function calls panic or exit",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Собираем все функции main в пакете
	var mainFuncs []*ast.FuncDecl
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
				mainFuncs = append(mainFuncs, fn)
			}
		}
	}

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	i.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		// Проверка на panic - всегда запрещена
		if isIdent(call.Fun, "panic") {
			pass.Reportf(call.Pos(), "prohibited use of panic()")
			return
		}

		// Проверка на log.Fatal и os.Exit
		if isLogFatalOrOsExit(call.Fun) {
			// Разрешаем только если находимся в функции main пакета main
			if pass.Pkg.Name() != "main" || !isInsideMainFunc(mainFuncs, call) {
				pass.Reportf(call.Pos(), "prohibited use of log.Fatal or os.Exit outside main package main function")
			}
		}
	})

	return nil, nil
}

func isIdent(expr ast.Expr, name string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == name
}

func isLogFatalOrOsExit(expr ast.Expr) bool {
	switch x := expr.(type) {
	case *ast.SelectorExpr:
		if isIdent(x.X, "log") && x.Sel.Name == "Fatal" {
			return true
		}
		if isIdent(x.X, "os") && x.Sel.Name == "Exit" {
			return true
		}
	}
	return false
}

func isInsideMainFunc(mainFuncs []*ast.FuncDecl, node ast.Node) bool {
	for _, mainFunc := range mainFuncs {
		if mainFunc.Body != nil {
			// Проверяем, находится ли вызов внутри тела функции main
			if node.Pos() >= mainFunc.Body.Pos() && node.Pos() <= mainFunc.Body.End() {
				return true
			}
		}
	}
	return false
}
