package main

import (
	"embed"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

//go:embed tmpl/*.tmpl
var tmplFS embed.FS

func main() {
	if err := execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute() error {
	var src string
	flag.StringVar(&src, "source", "", "input file name")
	flag.Parse()

	dir, base := filepath.Split(src)
	dst := filepath.Join(
		dir,
		strings.TrimSuffix(base, filepath.Ext(base))+".csv"+filepath.Ext(base),
	)

	fsrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	fdst, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fdst.Close()

	return run(fsrc, fdst)
}

func run(src io.Reader, dst io.Writer) error {
	f, err := parser.ParseFile(token.NewFileSet(), "", src, 0)
	if err != nil {
		return err
	}

	data := Data{Package: f.Name.Name}
	ast.Walk(&data, f)

	tmpl := template.Must(template.ParseFS(tmplFS, "tmpl/*.tmpl"))
	return tmpl.ExecuteTemplate(dst, "csv.go.tmpl", data)
}

type Data struct {
	Package string
	Structs []Struct
}

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Names []string
	Type  string
	Tag   string
}

func (data *Data) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch t := n.(type) {
	case *ast.TypeSpec:
		st, ok := t.Type.(*ast.StructType)
		if !ok {
			return nil
		}

		s := Struct{Name: t.Name.Name}
		var isCSV bool
		for _, f := range st.Fields.List {
			var names []string
			for _, name := range f.Names {
				names = append(names, name.Name)
			}

			var tag reflect.StructTag
			if f.Tag != nil {
				tag = reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
			}
			if v, ok := tag.Lookup("csv"); ok {
				isCSV = true
				s.Fields = append(s.Fields, Field{
					Names: names,
					Type:  fmt.Sprint(f.Type),
					Tag:   v,
				})
			}
		}
		if isCSV {
			data.Structs = append(data.Structs, s)
		}
	}

	return data
}
