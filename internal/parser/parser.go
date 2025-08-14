package parser

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/packages"
)

type FieldInfo struct {
	Name string
	Tag  string
	Type string
}

type StructInfo struct {
	Name           string
	Package        string
	PackageName    string
	Fields         []FieldInfo
	PrimaryKey     string
	PrimaryKeyType string
	FilePath       string
}

func ScanStructs(root string) ([]StructInfo, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes,
		Dir:  root,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, err
	}

	var structs []StructInfo

	for _, pkg := range pkgs {
		for i, file := range pkg.Syntax {
			filePath := pkg.GoFiles[i]
			ast.Inspect(file, func(n ast.Node) bool {
				ts, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					return true
				}

				var fields []FieldInfo
				var pkField, pkType string

				for _, f := range st.Fields.List {
					tag := ""
					if f.Tag != nil {
						tag = f.Tag.Value
					}

					for _, name := range f.Names {
						fieldType := exprToString(f.Type)
						fields = append(fields, FieldInfo{
							Name: name.Name,
							Tag:  tag,
							Type: fieldType,
						})

						if strings.Contains(tag, "primaryKey") || strings.Contains(tag, "PK") {
							pkField = name.Name
							pkType = fieldType
						}
					}
				}

				structs = append(structs, StructInfo{
					Name:           ts.Name.Name,
					Package:        pkg.PkgPath,
					PackageName:    pkg.Name,
					Fields:         fields,
					PrimaryKey:     pkField,
					PrimaryKeyType: pkType,
					FilePath:       filePath,
				})
				return true
			})
		}
	}

	return structs, nil
}

func exprToString(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", exprToString(t.X), t.Sel.Name)
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	default:
		return ""
	}
}
