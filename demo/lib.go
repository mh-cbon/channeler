package demo

import "fmt"

//go:generate channeler tomate_gen.go Tomate:MyTomate *Tomate:MyTomatePointer

// Tomate is a vegetable.
type Tomate struct {
	name string
}

// Hello world!
func (t *Tomate) Hello() { fmt.Println(" world!") }

// Good bye!
func (t Tomate) Good() { fmt.Println(" bye!") }

// Name it!
func (t Tomate) Name(it string) string { return fmt.Sprintf("Name:%v\n", it) }

// NewTomate isa contrstuctor
func NewTomate(n string) *Tomate {
	return &Tomate{
		name: n,
	}
}
