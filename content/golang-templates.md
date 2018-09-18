Title: Simple arrays in Golang templates
Date: 2018-09-18
Category: golang
Status: Draft


While working on creating [kops templates](), I had to figure out

```
package main

import (
	"html/template"
	"log"
	"os"
)



func main() {

	var names = []string{"Tom", "Jill"}
	//create a new template with some name
	tmpl := template.New("test")

	//parse some content and generate a template
	tmpl, err := tmpl.Parse("{{range .}}Hello {{.}}\n{{end}}")
	if err != nil {
		log.Fatal("Error Parsing template: ", err)
		return
	}
	err1 := tmpl.Execute(os.Stdout, names)
	if err1 != nil {
		log.Fatal("Error executing template: ", err1)
		return
	}
}

```
