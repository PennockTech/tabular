#!/bin/sh
#
# Relies upon: <https://github.com/wadey/gocovmerge>
#
# Based upon mmindenhall's solution in <https://github.com/golang/go/issues/6909>
#

TOP="go.pennock.tech/tabular"

progname="$(basename "$0")"
trace() { printf >&2 "%s: %s\n" "$progname" "$*" ; }

trace "removing old c*.out files"
# We used to use c.partial.out in each directory, prior to Go 1.10
# introducing coverprofiles across multiple packages
find . -name c\*.out -execdir rm -v {} \;

trace "generating new coverage.out"
go test -cover -covermode=count -coverprofile=coverage.out -coverpkg ./... ./...

trace "suggestions:"
echo "  go tool cover -func=coverage.out | less"
echo "  go tool cover -html=coverage.out"
