# Katalis

[![Go Reference](https://pkg.go.dev/badge/github.com/NicoNex/katalis.svg)](https://pkg.go.dev/github.com/NicoNex/katalis)
[![Go Report Card](https://goreportcard.com/badge/github.com/NicoNex/katalis)](https://goreportcard.com/report/github.com/NicoNex/katalis)
[![codecov](https://codecov.io/gh/NicoNex/katalis/branch/master/graph/badge.svg)](https://codecov.io/gh/NicoNex/katalis)

Katalis is a type-safe, generic wrapper around the [Pogreb](https://github.com/akrylysov/pogreb) embedded key-value store for Go. It leverages Go generics to provide compile-time type safety while maintaining the performance and simplicity of Pogreb.

## Why Katalis?

Traditional key-value stores in Go require manual serialization/deserialization of keys and values to `[]byte`. This approach is:
- **Error-prone**: Type mismatches discovered only at runtime
- **Verbose**: Boilerplate code for encoding/decoding
- **Unsafe**: Easy to store/retrieve wrong types

Katalis solves these issues through:
- **Type Safety**: Compile-time guarantees for key-value types
- **Zero Boilerplate**: Automatic encoding/decoding via codec system
- **Flexibility**: Pluggable codecs for any type
- **Modern Go**: Built-in support for Go 1.23+ iterators

## Features

- **Type-Safe Operations**: Use Go generics to enforce type safety for key-value pairs
- **Flexible Codec System**: Predefined codecs for primitives, custom codecs for complex types
- **High Performance**: Built on Pogreb's fast, embedded storage engine
- **Iterator Support**: Native support for Go 1.23+ range-over-function pattern
- **Error Handling**: Multiple iteration strategies for different error handling needs
- **Simple API**: Clean, idiomatic Go interface

## Installation

```bash
go get github.com/NicoNex/katalis@latest
```

Requires Go 1.23 or later.

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/NicoNex/katalis"
)

func main() {
	// Open a type-safe database with string keys and int values
	db, err := katalis.Open("mydb", katalis.StringCodec, katalis.IntCodec)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Put and Get operations are type-safe
	db.Put("age", 42)          // Compile-time type checking
	age, _ := db.Get("age")    // Returns int, not interface{}

	fmt.Println(age) // Output: 42
}
```

## Design Philosophy

### Codec System

Katalis uses a **codec pattern** to handle serialization transparently. Each type needs a codec that implements:

```go
type Codec[T any] interface {
    Encode(T) ([]byte, error)
    Decode([]byte) (T, error)
}
```

#### Predefined Codecs

Katalis provides codecs for all Go primitive types:

```go
// Integer types
katalis.IntCodec, katalis.Int8Codec, katalis.Int16Codec, katalis.Int32Codec, katalis.Int64Codec
katalis.UintCodec, katalis.Uint8Codec, katalis.Uint16Codec, katalis.Uint32Codec, katalis.Uint64Codec

// Floating point
katalis.Float32Codec, katalis.Float64Codec

// Strings and bytes
katalis.StringCodec, katalis.BytesCodec
```

#### Custom Types with Gob

For structs and complex types, use `Gob`:

```go
type User struct {
	Name string
	Age  int
}

// Option 1: Explicit type parameter
db, _ := katalis.Open("users.db", katalis.StringCodec, katalis.Gob[User]())

// Option 2: Type inference
var user User
db, _ := katalis.Open("users.db", katalis.StringCodec, katalis.Gob(user))
```

### Iteration Methods

Katalis provides three iteration strategies, each with different tradeoffs:

#### 1. `Items()` - Simple, Error-Skipping Iteration

Use when you want clean, simple iteration that automatically skips corrupted entries:

```go
for key, value := range db.Items() {
	fmt.Printf("%s: %d\n", key, value)
	// Automatically skips entries that fail to decode
}
```

**When to use**: Default choice for most cases, when you want simplicity and resilience without explicit error handling.

#### 2. `AllItems()` - Explicit Error Handling

Use when you need to know about and handle decode errors explicitly:

```go
for entry, err := range db.AllItems() {
	if err != nil {
		log.Printf("Error decoding entry: %v", err)
		continue // or handle differently
	}
	fmt.Printf("%s: %d\n", entry.Key, entry.Value)
}
```

**When to use**: Production monitoring, debugging, when you need to log/report errors, or take specific action on failures.

#### 3. `Fold()` - Callback-Based Iteration

Use when you need custom error handling or accumulation logic:

```go
total := 0
err := db.Fold(func(key string, val int, err error) error {
	if err != nil {
		return err // Stop on error
	}
	total += val
	return nil
})
```

**When to use**: Aggregations, reductions, custom control flow.

## Usage Examples

### Basic CRUD Operations

```go
db, _ := katalis.Open("db", katalis.StringCodec, katalis.IntCodec)
defer db.Close()

// Put
db.Put("score", 100)

// Get
score, err := db.Get("score")
if err != nil {
	// Handle error
}

// Has
exists, _ := db.Has("score")

// Delete
db.Del("score")
```

### Working with Complex Types

```go
type Article struct {
	Title   string
	Author  string
	Content string
	Tags    []string
}

db, _ := katalis.Open(
	"articles.db",
	katalis.StringCodec,      // Keys are strings
	katalis.Gob[Article](),   // Values are Articles
)
defer db.Close()

article := Article{
	Title:   "Understanding Katalis",
	Author:  "Alice",
	Content: "...",
	Tags:    []string{"go", "database"},
}

// Store
db.Put("article-1", article)

// Retrieve with full type safety
retrieved, _ := db.Get("article-1")
fmt.Println(retrieved.Title) // No type assertions needed!
```

### Iteration Examples

```go
// Simple iteration (automatically skips errors)
count := 0
for key, val := range db.Items() {
	count++
}

// Iteration with explicit error handling
for entry, err := range db.AllItems() {
	if err != nil {
		log.Printf("Corrupted entry: %v", err)
		continue
	}
	process(entry.Key, entry.Value)
}

// Early exit
for key, val := range db.Items() {
	if key == "target" {
		break // Stop iteration
	}
}
```

### Error Handling Patterns

```go
// Pattern 1: Strict error handling
for entry, err := range db.AllItems() {
	if err != nil {
		return fmt.Errorf("database corrupted: %w", err)
	}
	// Process entry.Key and entry.Value
}

// Pattern 2: Graceful degradation
var failed []string
for entry, err := range db.AllItems() {
	if err != nil {
		failed = append(failed, "unknown")
		continue
	}
	// Process valid entries
}

// Pattern 3: Fold with accumulation
err := db.Fold(func(key string, val int, err error) error {
	if err != nil {
		return err // Propagate error
	}
	// Process and accumulate
	return nil
})
```

## API Reference

### Database Operations

```go
// Open database with codecs
Open[KT, VT any](path string, keyCodec Codec[KT], valCodec Codec[VT]) (DB[KT, VT], error)
OpenOptions[KT, VT any](path string, keyCodec Codec[KT], valCodec Codec[VT], opts *Options) (DB[KT, VT], error)

// CRUD operations
Put(key KT, val VT) error
Get(key KT) (VT, error)
Has(key KT) (bool, error)
Del(key KT) error

// Iteration
Items() iter.Seq2[KT, VT]                            // Simple, skips errors
AllItems() iter.Seq2[Entry[KT, VT], error]           // With error handling
Fold(fn func(key KT, val VT, err error) error) error // Callback-based

// Close
Close() error
```

### Entry Type

```go
type Entry[KT, VT any] struct {
	Key   KT
	Value VT
}
```

Used by `AllItems()` to return key-value pairs with error information.

## Performance Considerations

- **Codec Choice**: Primitive codecs (Int, String) are faster than Gob for simple types
- **Iteration Method**: `Items()` and `AllItems()` have similar performance; use `Items()` for simplicity, `AllItems()` when you need error details
- **Database Options**: Tune Pogreb options via `OpenOptions` for specific workloads

## Contributing

Contributions are welcome! Please ensure:
- Tests pass: `go test ./...`
- Code is formatted: `go fmt ./...`
- Changes are documented

## License

See [LICENSE](LICENSE) file for details.

## Acknowledgments

Built on top of [Pogreb](https://github.com/akrylysov/pogreb) by Constantine Peresypkin.
