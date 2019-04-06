Title: Things I learned in Golang recently
Date: 2019-04-09
Category: software

## Repeating the same argument to Printf

If we wanted to repeat the same argument to a call to `fmt.Printf()`, we can make use of "indexed" arguments.
That is, instead of writing `fmt.Printf("%s %s", "Hello", "Hello")`, we can write `fmt.Printf("%[1]s %[1]s", "Hello")`.
Learn about it in the [docs](https://golang.org/pkg/fmt/). 


