[![Go Report Card](https://goreportcard.com/badge/github.com/fmpwizard/go-quilljs-delta)](https://goreportcard.com/report/github.com/fmpwizard/go-quilljs-delta)

# go-quilljs-delta

It's a port of QuillJS Delta's repo to Go

Try to find a balance between staying close to the original and being idiomatic Go

Refer to [the original QuillJS repo](https://github.com/quilljs/delta) for reference

[![GoDoc](https://godoc.org/github.com/fmpwizard/go-quilljs-delta/delta?status.svg)](https://godoc.org/github.com/fmpwizard/go-quilljs-delta/delta)

## Motivation

Have a backend that supports operational transformations written in Go, that also supports Quill's Delta format, which we 
use at work.

### Why not using nodejs?

Pretty sure this will run faster with concurrent users.