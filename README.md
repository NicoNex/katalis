# katalis
Katalis is a Go package that wraps the [Pogreb](https://github.com/akrylysov/pogreb) key-value store, enhancing it with Go generics to provide type-safe operations. This package allows developers to leverage the efficiency of Pogreb while ensuring both compile-time type safety for their data and ease of use.

## Features
- Type-Safe Operations: Use Go generics to enforce type safety for key-value pairs.
- High Performance: Built on top of Pogreb, known for its fast and efficient storage capabilities.
- Simple API: Easy-to-use interface that integrates seamlessly with Go's type system.

## Installation
```bash
go get github.com/NicoNex/katalis@latest
```

## Usage
Here's a quick example on how to use katalis to store a key of type `string` and a value of type `User`:
```go
package main

import (
  "fmt"

  "github.com/NicoNex/katalis"
)

type User struct {
  Name    string
  Age     int
  Country string
}

func main() {
  // Create a custom gob codec for the User struct using the provided katalis.GobCodec.
  var userCodec = katalis.GobCodec[User]{}

  // Get the instance to the existing or newly created DB passing katalis.StringCodec as the
  // codec for the key, and our custom userCodec for the value.
  db, err := katalis.Open("path/to/db", katalis.StringCodec, userCodec)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  // Easily store the string-User key value pair without dealing with the conversion to []byte.
  err := db.Put("alice-id", User{Name: "Alice", Age: 20, Country: "Italy"})
  if err != nil {
    panic(err)
  }

  // The variable alice here will be of type User.
  alice, err := db.Get("alice-id")
  if err != nil {
    panic(err)
  }

  fmt.Println(alice)
}
```
