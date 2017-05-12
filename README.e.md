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
# Create a channeled version of Tomate to MyTomate to stdout
channeler - demo/Tomate:ChanTomate
# Create a channeled version of Tomate to MyTomate to gen_test/chantomate.go
channeler demo/Tomate:gen_test/ChanTomate
```
# API example

Following example demonstates a program using it to generate a channeled version of a type.

#### > {{cat "demo/main.go" | color "go"}}

Following code is the generated implementation of `Tomate` type.

#### > {{cat "demo/mytomate.go" | color "go"}}


# Recipes

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```

# History

[CHANGELOG](CHANGELOG.md)
