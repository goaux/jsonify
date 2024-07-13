# jsonify
Package jsonify provides utility functions for JSON encoding of various types, including protobuf messages and standard Go types.

[![Go Reference](https://pkg.go.dev/badge/github.com/goaux/jsonify.svg)](https://pkg.go.dev/github.com/goaux/jsonify)
[![Go Report Card](https://goreportcard.com/badge/github.com/goaux/jsonify)](https://goreportcard.com/report/github.com/goaux/jsonify)

It uses a custom [jsoniter][1] configuration for improved performance and consistent output.

[1]: https://pkg.go.dev/github.com/json-iterator/go

## Features

- Fast JSON encoding using jsoniter
- Support for protobuf messages
- Consistent output with sorted map keys
- Easy-to-use API with both error-returning and panic-on-error versions

## Installation

To install jsonify, use `go get`:

    go get github.com/goaux/jsonify

## Usage

Here are some examples of how to use jsonify:

    package main

    import (
        "fmt"
        "github.com/goaux/jsonify"
    )

    func main() {
        // Encoding to []byte
        data := map[string]interface{}{
            "name": "John Doe",
            "age":  30,
        }

        bytes, err := jsonify.Bytes(data)
        if err != nil {
            panic(err)
        }
        fmt.Printf("Bytes: %s\n", bytes)

        // Encoding to string
        str, err := jsonify.String(data)
        if err != nil {
            panic(err)
        }
        fmt.Printf("String: %s\n", str)

        // Using MustBytes and MustString (panics on error)
        fmt.Printf("MustBytes: %s\n", jsonify.MustBytes(data))
        fmt.Printf("MustString: %s\n", jsonify.MustString(data))
    }

## API

- `Bytes(v any) ([]byte, error)`: Encodes the given value as JSON and returns it as a byte slice.
- `MustBytes(v any) []byte`: Similar to Bytes but panics if an error occurs during encoding.
- `String(v any) (string, error)`: Encodes the given value as JSON and returns it as a string.
- `MustString(v any) string`: Similar to String but panics if an error occurs during encoding.
