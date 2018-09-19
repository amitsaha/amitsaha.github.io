Title: Simple arrays in Golang templates
Date: 2018-09-18
Category: golang
Status: Draft


While working on creating a template file for a Golang project, I wanted to better understand how to work
with arrays in Golang templates as available via the `html/template` package. 

Let's consider our first program:

```
package main

import (
	"html/template"
	"log"
	"os"
)

func main() {

	var names = []string{"Tabby", "Jill"}

	tmpl := template.New("test")

	tmpl, err := tmpl.Parse("Array contents: {{.}}")
	if err != nil {
		log.Fatal("Error Parsing template: ", err)
		return
	}
	err1 := tmpl.Execute(os.Stdout, names)
	if err1 != nil {
		log.Fatal("Error executing template: ", err1)

	}
}



```

When we run the above program [Go playground link](https://play.golang.org/p/St0g-6_G8_1), the output we get is:

```
Array contents: [Tabby Jill]
```

There are three main stages in the above program:

- Create a new `Template` object: `tmpl := template.New("test")`
- Parse a template string: `tmpl, err := tmpl.Parse("Array contents: {{.}}")`
- Execute the template: `err1 := tmpl.Execute(os.Stdout, names)` passing data in the `names` variable

Anything within `{{ }}` inside the template string is where we do _something_ with the data that we pass in 
when executing the template. The `.` (dot) refers to the data that is passed in. In the above example, the 
entire array contents of `names` is the value of `.`. Hence, the output has the entire array including the surrounding
`[]`.


## Learn more

- [Golang documentation on template](https://golang.org/pkg/text/template/)
