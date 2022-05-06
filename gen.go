//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"text/template"

	_ "embed"
)

// Inspired by sort/gen_sort_variants.go
type Variant struct {
	// Name is the variant name: should be unique among variants.
	Name string

	// Path is the file path into which the generator will emit the code for this
	// variant.
	Path string

	// Imports is the imports needed for this package.
	Imports string

	// FuncSuffix is appended to all function names in this variant's code. All
	// suffixes should be unique within a package.
	FuncSuffix string

	StructPrefix    string
	StructPrefixLow string
	StructSuffix    string

	ExtraFileds string

	// TypeParam is the optional type parameter for the function.
	TypeParam string

	// Funcs is a map of functions used from within the template. The following
	// functions are expected to exist:
	//
	//    Less (name, i, j):
	//      emits a comparison expression that checks if the value `name` at
	//      index `i` is smaller than at index `j`.
	//
	//    Swap (name, i, j):
	//      emits a statement that performs a data swap between elements `i` and
	//      `j` of the value `name`.
	Funcs template.FuncMap
}

func main() {
	base := &Variant{
		Name:            "ordered",
		Path:            "ordered.go",
		Imports:         "\"sync\"\n\"sync/atomic\"\n\"unsafe\"\n",
		FuncSuffix:      "",
		TypeParam:       "[T ordered]",
		StructPrefix:    "Ordered",
		StructPrefixLow: "ordered",
		StructSuffix:    "",
		Funcs: template.FuncMap{
			"Less": func(i, j string) string {
				return fmt.Sprintf("(%s < %s)", i, j)
			},
			"Equal": func(i, j string) string {
				return fmt.Sprintf("%s == %s", i, j)
			},
		},
	}
	generate(base)
	base.Name += "Desc"
	base.FuncSuffix += "Desc"
	base.StructSuffix += "Desc"
	base.Path = "ordereddesc.go"
	base.Funcs = template.FuncMap{
		"Less": func(i, j string) string {
			return fmt.Sprintf("(%s > %s)", i, j)
		},
		"Equal": func(i, j string) string {
			return fmt.Sprintf("%s == %s", i, j)
		},
	}
	generate(base)

	basefunc := &Variant{
		Name:            "func",
		Path:            "func.go",
		Imports:         "\"sync\"\n\"sync/atomic\"\n\"unsafe\"\n",
		FuncSuffix:      "Func",
		TypeParam:       "[T any]",
		ExtraFileds:     "\nless func(a,b T)bool\n",
		StructPrefix:    "Func",
		StructPrefixLow: "func",
		StructSuffix:    "",
		Funcs: template.FuncMap{
			"Less": func(i, j string) string {
				return fmt.Sprintf("s.less(%s,%s)", i, j)
			},
			"Equal": func(i, j string) string {
				return fmt.Sprintf("!s.less(%s,%s)", j, i)
			},
		},
	}
	generate(basefunc)
}

// generate generates the code for variant `v` into a file named by `v.Path`.
func generate(v *Variant) {
	// Parse templateCode anew for each variant because Parse requires Funcs to be
	// registered, and it helps type-check the funcs.
	tmpl, err := template.New("gen").Funcs(v.Funcs).Parse(templateCode)
	if err != nil {
		log.Fatal("template Parse:", err)
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, v)
	if err != nil {
		log.Fatal("template Execute:", err)
	}

	formatted, err := format.Source(out.Bytes())
	if err != nil {
		log.Fatal("format:", err)
	}

	if err := os.WriteFile(v.Path, formatted, 0644); err != nil {
		log.Fatal("WriteFile:", err)
	}
}

//go:embed skipset.tpl
var templateCode string
