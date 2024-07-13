package jsonify_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"unicode"

	"github.com/goaux/jsonify"
	"google.golang.org/protobuf/types/known/structpb"
)

func ExampleBytes() {
	b, _ := jsonify.Bytes(map[string]interface{}{"A": true, "B": "<b>"})
	fmt.Printf(">%s<\n", b)
	// Output:
	// >{"A":true,"B":"<b>"}<
}

func ExampleMustBytes() {
	fmt.Printf(">%s<\n", jsonify.MustBytes(map[string]any{"A": true, "B": "<b>"}))
	// output:
	// >{"A":true,"B":"<b>"}<
}

func ExampleString() {
	s, _ := jsonify.String(map[string]interface{}{"A": true, "B": "<b>"})
	fmt.Printf(">%s<\n", s)
	// Output:
	// >{"A":true,"B":"<b>"}<
}

func ExampleMustString() {
	fmt.Println(">" + jsonify.MustString(map[string]any{"A": true, "B": "<b>"}) + "<")
	// output:
	// >{"A":true,"B":"<b>"}<
}

func BenchmarkMustString(b *testing.B) {
	target := struct {
		Hello bool    `json:"hello"`
		World float64 `json:"world"`
	}{
		Hello: true,
		World: 42.195,
	}

	b.Run("jsoniter", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jsonify.MustString(target)
		}
	})
	b.Run("standard", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			j, err := json.Marshal(target)
			if err != nil {
				panic(err)
			}
			bytes.TrimRightFunc(j, unicode.IsSpace)
		}
	})
}

func TestBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []byte
		wantErr  bool
	}{
		{
			name:     "simple map",
			input:    map[string]interface{}{"A": true, "B": "<b>"},
			expected: []byte(`{"A":true,"B":"<b>"}`),
		},
		{
			name:     "json.RawMessage",
			input:    json.RawMessage(`{"raw":"message"}`),
			expected: []byte(`{"raw":"message"}`),
		},
		{
			name:     "nested structure",
			input:    map[string]interface{}{"A": 1, "B": map[string]interface{}{"C": "nested"}},
			expected: []byte(`{"A":1,"B":{"C":"nested"}}`),
		},
		{
			name:     "slice",
			input:    []int{1, 2, 3},
			expected: []byte(`[1,2,3]`),
		},
		{
			name:    "channel (invalid JSON type)",
			input:   make(chan int),
			wantErr: true,
		},
		{
			name:     "nil value",
			input:    nil,
			expected: []byte(`null`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonify.Bytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Bytes() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "simple map",
			input:    map[string]interface{}{"A": true, "B": "<b>"},
			expected: `{"A":true,"B":"<b>"}`,
		},
		{
			name:     "json.RawMessage",
			input:    json.RawMessage(`{"raw":"message"}`),
			expected: `{"raw":"message"}`,
		},
		{
			name:     "nested structure",
			input:    map[string]interface{}{"A": 1, "B": map[string]interface{}{"C": "nested"}},
			expected: `{"A":1,"B":{"C":"nested"}}`,
		},
		{
			name:    "channel (invalid JSON type)",
			input:   make(chan int),
			wantErr: true,
		},
		{
			name:     "nil value",
			input:    nil,
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonify.String(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMustBytes(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		input := map[string]interface{}{"A": true, "B": "<b>"}
		expected := []byte(`{"A":true,"B":"<b>"}`)
		result := jsonify.MustBytes(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("MustBytes() = %s, want %s", result, expected)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustBytes() did not panic")
			}
		}()
		jsonify.MustBytes(make(chan int))
	})
}

func TestMustString(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		input := map[string]interface{}{"A": true, "B": "<b>"}
		expected := `{"A":true,"B":"<b>"}`
		result := jsonify.MustString(input)
		if result != expected {
			t.Errorf("MustString() = %v, want %v", result, expected)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustString() did not panic")
			}
		}()
		jsonify.MustString(make(chan int))
	})
}

func TestProtobufMessage(t *testing.T) {
	pbMsg, err := structpb.NewStruct(map[string]any{
		"foo": "bar",
	})
	if err != nil {
		panic(err)
	}

	t.Run("Bytes with protobuf", func(t *testing.T) {
		got, err := jsonify.Bytes(pbMsg)
		if err != nil {
			t.Fatalf("Bytes() error = %v", err)
		}
		expected := []byte(`{"foo":"bar"}`)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Bytes() = %s, want %s", got, expected)
		}
	})

	t.Run("String with protobuf", func(t *testing.T) {
		got, err := jsonify.String(pbMsg)
		if err != nil {
			t.Fatalf("String() error = %v", err)
		}
		expected := `{"foo":"bar"}`
		if got != expected {
			t.Errorf("String() = %v, want %v", got, expected)
		}
	})
}
