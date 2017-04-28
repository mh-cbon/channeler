package demo

import "fmt"

//go:generate channeler tomate_gen.go Tomate:MyTomate

// Tomate is a vegetable.
type Tomate struct{}

// Hello world!
func (t *Tomate) Hello() { fmt.Println(" world!") }

// Good bye!
func (t Tomate) Good() { fmt.Println(" bye!") }

// Name it!
func (t Tomate) Name(it string) string { return fmt.Sprintf("Name:%v\n", it) }
