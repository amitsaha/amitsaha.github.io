package main

import "fmt"

type myiface interface {
	func1(r myiface)
	func2()
}

// mytype implements myiface
type mytype struct{}

func (t *mytype) func1(r myiface) {
	fmt.Printf("func1 of mytype called")
	r.func2()
}

func (t *mytype) func2() {
	fmt.Printf("func2 of mytype called")
}

func myfunc(t myiface) {
	t.func1(t)
}

// mytype1 implements myiface
// also embeds mytype
type mytype1 struct {
	mytype
}

func (t *mytype1) func1(t1 myiface) {
	fmt.Printf("func1 of mytype1 called")
	t.mytype.func1(t1)
}

func (t *mytype1) func2() {
	fmt.Printf("func2 of mytype1 called")
}

func main() {
	t := &mytype1{}
	myfunc(t)
}
