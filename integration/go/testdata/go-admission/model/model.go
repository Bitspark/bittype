package model

type Person struct {
	Name   string
	Active bool
}

type Optional[T any] struct {
	Present bool
	Value   T
}

type Page[T any] struct {
	Items []T
	Next  Optional[string]
}

type Wire struct {
	Send func([]byte)
}
