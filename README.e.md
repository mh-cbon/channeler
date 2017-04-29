---
License: MIT
LicenseFile: LICENSE
LicenseColor: yellow
---
# {{.Name}}

{{template "badge/travis" .}} {{template "badge/appveyor" .}} {{template "badge/goreport" .}} {{template "badge/godoc" .}} {{template "license/shields" .}}

{{pkgdoc}}

# {{toc 5}}

# Install
{{template "glide/install" .}}

## Usage

#### $ {{exec "channeler" "-help" | color "sh"}}

## Cli examples

```sh
# Create a channeled version os Tomate to MyTomate
channeler tomate_gen.go Tomate:MyTomate
```
# API example

Following example demonstates a program using it to generate a channeled version of a type.

#### > {{cat "demo/lib.go" | color "go"}}

Following code is the generated implementation of `Tomate` type.

#### > {{cat "demo/tomate_gen.go" | color "go"}}


# Recipes

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```

# History

[CHANGELOG](CHANGELOG.md)
