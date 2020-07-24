package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func main() {
	filePath := os.Getenv("GOFILE")
	pkgName := os.Getenv("GOPACKAGE")

	if len(filePath) == 0 || len(pkgName) == 0 {
		tools.Die("Required environment variable did not set, must call by `go generate`")
	}
	fmt.Printf("GOFILE=%s\n", filePath)
	fmt.Printf("GOPACKAGE=%s\n", pkgName)

	fileNameBase := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	resultFile := filepath.Join(filepath.Dir(filePath), fileNameBase+".getters"+filepath.Ext(filePath))

	fmt.Printf(" * resultFile: %s\n", resultFile)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)

	if err != nil {
		tools.Die("Failed parse input file: %s", err.Error())
	}

	genContent := "package " + file.Name.String() + "\n\n"

	for _, node := range file.Decls {
		if decl, ok := node.(*ast.GenDecl); ok {
			if decl.Tok != token.TYPE {
				continue
			}
			typeSpec, ok := decl.Specs[0].(*ast.TypeSpec)
			if !ok {
				continue
			}
			structName := typeSpec.Name.String()
			structSpec, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			genInterface := fmt.Sprintf("type readonly%s interface {\n", tools.Ucfirst(structName))
			genGetter := ""
			for _, field := range structSpec.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				varName := field.Names[0].String()
				var varType string = ""
				// := field.Type.(*ast.Ident).String()
				if id, ok := field.Type.(*ast.Ident); ok {
					varType = id.String()
				} else if id, ok := field.Type.(*ast.ArrayType); ok {
					varType = id.Elt.(*ast.Ident).String()
					if id.Len == nil {
						varType = "[]" + varType
					} else {
						varType = "[" + id.Len.(*ast.BasicLit).Value + "]" + varType
					}
				} else {
					fmt.Println("Unknown type:")
					spew.Dump(field)
					os.Exit(1)
				}

				genInterface += fmt.Sprintf("\tGet%s() %s\n", varName, varType)
				genGetter += fmt.Sprintf("func (s %s) Get%s() %s {\n\treturn s.%s;\n}\n\n", structName, varName, varType, varName)
			}

			genContent += genInterface + "}\n\n" + genGetter + "\n"
		}
	}

	err = ioutil.WriteFile(resultFile, []byte(genContent), os.FileMode(0666))

	if err != nil {
		tools.Die("Failed write output file: %s", err.Error())
	}
}
