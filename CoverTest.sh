#!/bin/sh

progname="$(basename -s .sh "$0")"
trace() { printf >&2 "%s: %s\n" "$progname" "$*" ; }

# Remove this block 2018Q3 or thereafter:
trace "removing old c*.out files"
# We used to use c.partial.out in each directory, prior to Go 1.10
# introducing coverprofiles across multiple packages
find . -name c\*.out -execdir rm -v {} \;

trace "generating new coverage.out"
go test -cover -covermode=count -coverprofile=coverage.out -coverpkg ./... ./...

trace "suggestions:"
echo "  go tool cover -func=coverage.out | less"
echo "  go tool cover -html=coverage.out"
