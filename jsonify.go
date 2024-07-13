// Package jsonify provides utility functions for JSON encoding of various
// types, including protobuf messages and standard Go types.
//
// It uses a custom [jsoniter] configuration for improved performance and
// consistent output.
//
// custom [jsoniter] configuration:
//
//	var config = jsoniter.Config{
//		SortMapKeys:            true,
//		ValidateJsonRawMessage: true,
//	}.Froze()
//
// This configuration is similar to [jsoniter.ConfigCompatibleWithStandardLibrary].
// The only difference is that EscapeHTML is set to false.
package jsonify

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var config = jsoniter.Config{
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
}.Froze()

// Bytes encodes the given value as JSON and returns it as a byte slice.
//
// It handles [json.RawMessage], [proto.Message], and other types differently.
// For [json.RawMessage], it returns the raw bytes.
// For [proto.Message], it uses [protojson] for marshaling.
// For other types, it uses a custom [jsoniter] configuration.
func Bytes(v any) ([]byte, error) {
	switch v := v.(type) {
	case json.RawMessage:
		return []byte(v), nil
	case proto.Message:
		return protojson.Marshal(v)
	}
	return config.Marshal(v)
}

// MustBytes is similar to [Bytes] but panics if an error occurs during encoding.
//
// It's useful when you're certain that the encoding will succeed.
func MustBytes(v any) []byte {
	b, err := Bytes(v)
	if err != nil {
		panic(err)
	}
	return b
}

// String encodes the given value as JSON and returns it as a string.
//
// It handles [json.RawMessage], [proto.Message], and other types differently.
// For [json.RawMessage], it returns the raw message as a string.
// For [proto.Message], it uses [protojson] for marshaling.
// For other types, it uses a custom [jsoniter] configuration.
func String(v any) (string, error) {
	switch v := v.(type) {
	case json.RawMessage:
		return string(v), nil
	case proto.Message:
		b, err := protojson.Marshal(v)
		return string(b), err
	}
	return config.MarshalToString(v)
}

// MustString is similar to [String] but panics if an error occurs during encoding.
//
// It's useful when you're certain that the encoding will succeed.
func MustString(v any) string {
	s, err := String(v)
	if err != nil {
		panic(err)
	}
	return s
}
