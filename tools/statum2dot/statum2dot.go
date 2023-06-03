package main

import (
	"flag"
	"fmt"
	"go/ast"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
)

func main() {
	flag.Parse()

	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	// Print the names of the source files
	// for each package listed on the command line.
	for _, pkg := range pkgs {
		fmt.Println(pkg.ID, pkg.GoFiles)
		extractPkg(pkg)
	}
}

func extractPkg(pkg *packages.Package) {
	for _, f := range pkg.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.AssignStmt:
				evalAssignStmt(v)
			}
			return true
		})
	}
}

func evalAssignStmt(v *ast.AssignStmt) {
	for _, rh := range v.Rhs {
		evalNewStateMachineConfig(rh)
	}
}

func evalNewStateMachineConfig(expr ast.Expr) {
	log.Printf("evalNewStateMachineConfig: %v\n", expr)
	callExpr, ok := expr.(*ast.CallExpr)

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)

	for _, arg := range callExpr.Args {

	}

	ie, ok := expr.(*ast.IndexExpr)
	if !ok || ie == nil {
		return
	}
	se, ok := ie.X.(*ast.SelectorExpr)
	//se, ok := expr.(*ast.SelectorExpr)
	if !ok || se == nil {
		return
	}
	log.Println("evalNewStateMachineConfig", se)

	i, ok := se.X.(*ast.Ident)
	if !ok || i == nil {
		return
	}
	if i.Name == "statum" {
		log.Println("big o", i)
	}

}
