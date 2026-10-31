package fixtures

import "time"

// Generic declarations, applications and recursive indirection.
type Person struct { Name string }
type Page[T any] struct { Items []T; Next *string }
type Applied = Page[Person]
type Integer[T interface { ~int | ~int64 }] struct { Value T }
type Tagged struct { Name string `json:"name"`; A, B int; *Person }
type Shapes struct {
	Fixed [3]int
	Mapping map[string][]*Person
	Receive <-chan int
	Send chan<- int
	Both chan bool
	Pkg time.Time
	F func(a, b int) (result string)
	Variadic func(args ...string)
}
type Reader[T any] interface { Read(key string) (T, error) }
type Parent interface { Reader[Person] }
