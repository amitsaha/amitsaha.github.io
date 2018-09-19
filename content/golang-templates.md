Title: Examples of consuming data in Golang templates
Date: 2018-09-18
Category: golang
Status: Draft

While working on creating a template file for a Golang project, I wanted to better understand how to work
with data in Golang templates as available via the `html/template` package. In this post, I discuss
some common use cases that may arise.

## Access a variable

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
`[]`. This also means that `names` could have been of another type - a struct for [example](https://play.golang.org/p/vAmNzNFg8LR)
and our program above would have worked:

```
package main

import (
	"html/template"
	"log"
	"os"
)

type Test struct {
	name string
}

func main() {

	t := Test{name: "Tabby"}

	tmpl := template.New("test")

	//parse some content and generate a template
	tmpl, err := tmpl.Parse("Variable contents: {{.}}")
	if err != nil {
		log.Fatal("Error Parsing template: ", err)
		return
	}
	err1 := tmpl.Execute(os.Stdout, t)
	if err1 != nil {
		log.Fatal("Error executing template: ", err1)

	}
}

```

The output now is:

```
Variable contents: {Tabby}
```

## Accessing structure members

Now, let's consider that our structure has multiple members and we want to access the individual members in our
template. Here's how we can do so (Golang playground)[https://play.golang.org/p/8BSiYJ_7Mfd]:

```
package main

import (
	"html/template"
	"log"
	"os"
)

type Person struct {
	Name string
	Age  int
}

func main() {

	p := Person{Name: "Tabby", Age: 21}

	tmpl := template.New("test")

	//parse some content and generate a template
	tmpl, err := tmpl.Parse("{{.Name}} is {{.Age}} years old")
	if err != nil {
		log.Fatal("Error Parsing template: ", err)
		return
	}
	err1 := tmpl.Execute(os.Stdout, p)
	if err1 != nil {
		log.Fatal("Error executing template: ", err1)

	}
}
```

The `dot` operator referes to the structure object, `p` and then inside the template, we just specify the
field name, like so, `.<Field>`. The output will be:

```
Tabby is 21 years old
```

## Do something with array elements

## Arrays of structure objects

## Values from files


## Learn more

- [Golang documentation on template](https://golang.org/pkg/text/template/)
