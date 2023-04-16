package preprocessor

import "strings"

func ReplaceString(s string, replaceIdx map[string]string) string {
	if strings.HasPrefix(s, "{{") && strings.HasSuffix(s, "}}") {
		key := s[2 : len(s)-2]
		if val, ok := replaceIdx[key]; ok {
			return val
		}
	}

	return s
}

/*

package main

import (
	"fmt"
	"os"
	"text/template"
)

func main() {
	fmt.Println("Hello World")
	type Inventory struct {
		Material string
		Count    uint
	}
	foobar := map[string]string{}
	foobar["Count"] = "42"
	foobar["Material"] = "Polyester"
	//sweaters := Inventory{"wool", 17}
	tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, foobar)
	if err != nil {
		panic(err)
	}
}

*/
