package docs

import (
	"fmt"
	"go/ast"
	"strings"
)

func exprName(e ast.Expr) string {
	switch s := e.(type) {
	case *ast.StarExpr:
		return "*" + exprName(s.X)
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", s.X, s.Sel)
	case *ast.CallExpr:
		name := exprName(s.Fun)
		args := exprNameList(s.Args)

		return fmt.Sprintf("%s(%s)", name, strings.Join(args, ","))
	case *ast.Ident:
		return s.Name
	case *ast.BasicLit:
		return s.Value
	default:
		// log.Println("->", reflect.TypeOf(s))
	}

	return ""
}

func callName(e ast.Expr) string {
	if c, ok := e.(*ast.CallExpr); ok {
		return strings.Join(exprNameList(c.Args), ",")
	}

	return ""
}

func argName(e ast.Expr) string {
	switch b := e.(type) {
	case *ast.BasicLit:
		return b.Value[1 : len(b.Value)-1]
	}
	return ""
}

func exprNameList(arr []ast.Expr) []string {
	names := make([]string, len(arr))

	for i, e := range arr {
		names[i] = exprName(e)
	}

	return names
}
