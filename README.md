tabular
=======

[![Continuous Integration](https://secure.travis-ci.org/PennockTech/tabular.svg?branch=main)](http://travis-ci.org/PennockTech/tabular)
[![Documentation](https://godoc.org/go.pennock.tech/tabular?status.svg)](https://godoc.org/go.pennock.tech/tabular)
[![Coverage Status](https://coveralls.io/repos/github/PennockTech/tabular/badge.svg)](https://coveralls.io/github/PennockTech/tabular)
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

---

### Projects using Tabular

* [character](https://github.com/philpennock/character) — Unicode character
  lookups and manipulations
