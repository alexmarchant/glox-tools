package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(1)
	}
	outputDir := args[0]

	// Expressions
	defineAst(outputDir, "Expr", map[string][]string{
		"BinaryExpr":   []string{"Left Expr", "Operator *Token", "Right Expr"},
		"GroupingExpr": []string{"Expression Expr"},
		"LiteralExpr":  []string{"Value interface{}"},
		"UnaryExpr":    []string{"Operator *Token", "Right Expr"},
	})

	// Statements
	defineAst(outputDir, "Stmt", map[string][]string{
		"ExpressionStmt": []string{"Expression Expr"},
		"PrintStmt":      []string{"Expression Expr"},
	})
}

func defineAst(outputDir string, baseName string, types map[string][]string) {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"
	f, err := os.Create(path)
	if err != nil {
		panic("Couldn't open file for writing.")
	}
	w := bufio.NewWriter(f)

	// Header
	w.WriteString("package main\n")
	w.WriteString("\n")

	// Type enum
	w.WriteString(fmt.Sprintf("type %sType int\n", baseName))
	w.WriteString("const (\n")
	enumTypeDeclared := false
	for typeName := range types {
		if !enumTypeDeclared {
			w.WriteString(fmt.Sprintf("\t%sType%s %sType = iota\n", baseName, typeName, baseName))
			enumTypeDeclared = true
		} else {
			w.WriteString(fmt.Sprintf("\t%sType%s\n", baseName, typeName))
		}
	}
	w.WriteString(")\n")
	w.WriteString("\n")

	// Type Interface
	w.WriteString(fmt.Sprintf("type %s interface {\n", baseName))
	w.WriteString(fmt.Sprintf("\t%sType() %sType\n", baseName, baseName))
	w.WriteString(fmt.Sprintf("\tAccept(%sVisitor) (interface{}, error)\n", baseName))
	w.WriteString("}\n")
	w.WriteString("\n")

	// Visitor Interfaces
	defineVisitor(w, baseName, types)

	// Types
	for typeName, fields := range types {
		defineType(w, baseName, typeName, fields)
	}

	w.Flush()
}

func defineType(w *bufio.Writer, baseName string, typeName string, fields []string) {
	// Struct
	w.WriteString(fmt.Sprintf("type %s struct {\n", typeName))
	for _, field := range fields {
		w.WriteString(fmt.Sprintf("\t%s\n", field))
	}
	w.WriteString("}\n")
	w.WriteString("\n")

	// Type interface method
	w.WriteString(fmt.Sprintf("func (t *%s) %sType() %sType {\n", typeName, baseName, baseName))
	w.WriteString(fmt.Sprintf("\treturn %sType%s\n", baseName, typeName))
	w.WriteString("}\n")
	w.WriteString("\n")

	// Visitor accept interface method
	w.WriteString(fmt.Sprintf("func (t *%s) Accept(visitor %sVisitor) (interface{}, error) {\n", typeName, baseName))
	w.WriteString(fmt.Sprintf("\treturn visitor.Visit%s(t)\n", typeName))
	w.WriteString("}\n")
	w.WriteString("\n")
}

func defineVisitor(w *bufio.Writer, baseName string, types map[string][]string) {
	w.WriteString(fmt.Sprintf("type %sVisitor interface {\n", baseName))
	for typeName := range types {
		w.WriteString(fmt.Sprintf("\tVisit%s(*%s) (interface{}, error)\n", typeName, typeName))
	}
	w.WriteString("}\n")
	w.WriteString("\n")
}
