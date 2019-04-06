Title: Things I learned in Golang recently
Date: 2019-04-09
Category: software

## Repeating the same argument to Printf

If we wanted to repeat the same argument to a call to `fmt.Printf()`, we can make use of "indexed" arguments.
That is, instead of writing `fmt.Printf("%s %s", "Hello", "Hello")`, we can write `fmt.Printf("%[1]s %[1]s", "Hello")`.
Learn about it in the [docs](https://golang.org/pkg/fmt/). 

## Multi-line strings

Things are hassle free on the multi-line strings front:

```
package main

import (
	"fmt"
)

func main() {

	s := `Multi line strings
are easy. So are strings with "double quotes"
and 'single quotes'`
	fmt.Print(s)

}

```


## Maps with values as maps

First, we define the map:

```
var clusters = make(map[int]map[string]int)
```


Then, we assign a value which is another map, that we create:

```
clusters[clusterNum] = map[string]int{
				"a": 1,
				"b": 1,
}

```

We can then modify the map defined as the value, like so:

```
clusters[1]["a"] = 5
```

## Check if a key is present in a map

Here we can make use of the multiple return value from the map query statement:

```
// The ok variable will be true if key present
// else false
if _, ok := flatMap[tblName]; !ok {
     // not present
} else {
    // present
}
```

