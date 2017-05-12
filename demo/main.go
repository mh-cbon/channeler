package main

import "fmt"

//go:generate channeler Tomate:MyTomate *Tomate:MyTomatePointer

func main() {
	x := NewMyTomate()
	fmt.Println(
		x.Name("world"),
	)
	y := NewMyTomatePointer("s")
	fmt.Println(
		y.Name("world"),
	)
}

// Tomate is a vegetable.
type Tomate struct {
	name string
}

// Hello world!
func (t *Tomate) Hello() { fmt.Println(" world!") }

// Good bye!
func (t Tomate) Good() { fmt.Println(" bye!") }

// Name it!
func (t Tomate) Name(it string) string { return fmt.Sprintf("Hello %v!\n", it) }

// NewTomate is a contrstuctor
func NewTomate(n string) *Tomate {
	return &Tomate{
		name: n,
	}
}
