# channeler

[![travis Status](https://travis-ci.org/mh-cbon/channeler.svg?branch=master)](https://travis-ci.org/mh-cbon/channeler) [![Appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/channeler?branch=master&svg=true)](https://ci.appveyor.com/projects/mh-cbon/channeler) [![Go Report Card](https://goreportcard.com/badge/github.com/mh-cbon/channeler)](https://goreportcard.com/report/github.com/mh-cbon/channeler) [![GoDoc](https://godoc.org/github.com/mh-cbon/channeler?status.svg)](http://godoc.org/github.com/mh-cbon/channeler) [![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Package channeler generates synced type using channels.


# TOC
- [Install](#install)
  - [Usage](#usage)
    - [$ channeler -help](#-channeler--help)
  - [Cli examples](#cli-examples)
- [API example](#api-example)
  - [> demo/main.go](#-demomaingo)
  - [> demo/mytomate.go](#-demomytomatego)
- [Recipes](#recipes)
  - [Release the project](#release-the-project)
- [History](#history)

# Install
```sh
mkdir -p $GOPATH/src/github.com/mh-cbon/channeler
cd $GOPATH/src/github.com/mh-cbon/channeler
git clone https://github.com/mh-cbon/channeler.git .
glide install
go install
```

## Usage

#### $ channeler -help
```sh
channeler 0.0.0

Usage

  channeler [-p name] [...types]

  types:  A list of types such as src:dst.
          A type is defined by its package path and its type name,
          [pkgpath/]name
          If the Package path is empty, it is set to the package name being generated.
          Name can be a valid type identifier such as TypeName, *TypeName, []TypeName 
  -p:     The name of the package output.
```

## Cli examples

```sh
# Create a channeled version of Tomate to MyTomate to stdout
channeler - demo/Tomate:ChanTomate
# Create a channeled version of Tomate to MyTomate to gen_test/chantomate.go
channeler demo/Tomate:gen_test/ChanTomate
```
# API example

Following example demonstates a program using it to generate a channeled version of a type.

#### > demo/main.go
```go
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
```

Following code is the generated implementation of `Tomate` type.

#### > demo/mytomate.go
```go
package main

// file generated by
// github.com/mh-cbon/channeler
// do not edit

import (
	"encoding/json"
)

// MyTomate is channeled.
type MyTomate struct {
	embed Tomate
	ops   chan func()
	stop  chan bool
	tick  chan bool
}

// NewMyTomate constructs a channeled version of Tomate
func NewMyTomate() *MyTomate {
	ret := &MyTomate{
		ops:  make(chan func()),
		tick: make(chan bool),
		stop: make(chan bool),
	}
	go ret.Start()
	return ret
}

// Hello is channeled
func (t *MyTomate) Hello() {
	t.ops <- func() {
		t.embed.Hello()
	}
	<-t.tick
}

// Good is channeled
func (t *MyTomate) Good() {
	t.ops <- func() {
		t.embed.Good()
	}
	<-t.tick
}

// Name is channeled
func (t *MyTomate) Name(it string) string {
	var retVar0 string
	t.ops <- func() {
		retVar0 = t.embed.Name(it)
	}
	<-t.tick
	return retVar0
}

// Start the main loop
func (t *MyTomate) Start() {
	for {
		select {
		case op := <-t.ops:
			op()
			t.tick <- true
		case <-t.stop:
			return
		}
	}
}

// Stop the main loop
func (t *MyTomate) Stop() {
	t.stop <- true
}

//UnmarshalJSON JSON unserializes MyTomate
func (t *MyTomate) UnmarshalJSON(b []byte) error {
	var embed Tomate
	var err error
	t.ops <- func() {
		err = json.Unmarshal(b, &embed)
		if err == nil {
			t.embed = embed
		}
	}
	<-t.tick
	return err
}

//MarshalJSON JSON serializes MyTomate
func (t *MyTomate) MarshalJSON() ([]byte, error) {
	var ret []byte
	var err error
	t.ops <- func() {
		ret, err = json.Marshal(t.embed)
	}
	<-t.tick
	return ret, err
}
```


# Recipes

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```

# History

[CHANGELOG](CHANGELOG.md)
