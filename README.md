tabular
=======

[![Continuous Integration](https://github.com/PennockTech/tabular/actions/workflows/pushes.yaml/badge.svg)](https://github.com/PennockTech/tabular/actions/workflows/pushes.yaml)
[![Documentation](https://godoc.org/go.pennock.tech/tabular?status.svg)](https://godoc.org/go.pennock.tech/tabular)
[![Coverage Status](https://coveralls.io/repos/github/PennockTech/tabular/badge.svg?branch=main)](https://coveralls.io/github/PennockTech/tabular?branch=main)
[![Current Tag](https://img.shields.io/github/tag/PennockTech/tabular.svg)](https://github.com/PennockTech/tabular/releases)
[![Issues](https://img.shields.io/github/issues/PennockTech/tabular.svg)](https://github.com/PennockTech/tabular/issues)
[![Repo Size](https://img.shields.io/github/repo-size/PennockTech/tabular.svg)](https://github.com/PennockTech/tabular)

The `tabular` package provides a Golang library for storing data in a table
consisting of rows and columns.  Sub-packages provide for rendering such a
table as a terminal box-table (line-drawing with UTF-8 box-drawing in various
styles, or ASCII), as HTML, or as CSV data.

The core data model is designed to be extensible and powerful, letting such
a table be embedded in various more sophisticated models.  (Eg, core of a
spreadsheet).  Cells in the table contain arbitrary data and possess metadata
in the form of "properties", modelled after the `context` package's `Context`.

A table can be created from the base package and then populated, before being
passed to any of the renderers, or a table can be directly created using a
sub-package, such that you _probably_ won't need to import the base package
directly.

An overview guide to the codebase can be found in
[the Overview.md](Overview.md)

[The usage documentation is in Godoc format](https://godoc.org/go.pennock.tech/tabular)

See the [examples](examples/) for a gentler introduction.

This package should be installable in the usual `go get` manner.

This software is under a [MIT-type license](LICENSE.txt).

When embedding into another tool, Go Modules support is able to report on the
version numbers of all dependencies with `go version -m $cmdname`; to support
your own version reporting framework, `go.pennock.tech/tabular.Versions()`
returns a slice of strings (including API versions and the value of the
`LinkerSpecifiedVersion` top-level variable).

This package uses [semantic versioning](https://semver.org/).  
Note that Go only supports the most recent two minor versions of the language;
for the purposes of semver, we do not consider it a breaking change to add a
dependency upon a language or standard library feature supported by all
currently-supported releases of Go.

---

### Projects using Tabular

* [character](https://github.com/philpennock/character) — Unicode character
  lookups and manipulations
