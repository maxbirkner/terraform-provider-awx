package awx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestNormalizeJsonOk(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		ok       bool
	}{
		{
			name:     "Valid JSON",
			input:    `{"key": "value"}`,
			expected: `{"key":"value"}`,
			ok:       true,
		},
		{
			name:     "Invalid JSON",
			input:    `{"key": "value"`,
			expected: "Error parsing JSON: unexpected end of JSON input",
			ok:       false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
			ok:       true,
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: "",
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := normalizeJsonOk(tt.input)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeYamlOk(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		ok       bool
	}{
		{
			name:     "Valid YAML",
			input:    "key: value",
			expected: "key: value\n",
			ok:       true,
		},
		{
			name:     "Valid quoted YAML",
			input:    `"key": value`,
			expected: "key: value\n",
			ok:       true,
		},
		{
			name:     "Valid YAML (bool)",
			input:    `"key": true`,
			expected: "key: true\n",
			ok:       true,
		},
		{
			name:     "Valid YAML (special characters)",
			input:    `"key": "&a, b, c"`,
			expected: "key: '&a, b, c'\n",
			ok:       true,
		},
		{
			name:     "Invalid YAML",
			input:    "key: [a, b, c",
			expected: "Error parsing YAML: yaml: line 1: did not find expected ',' or ']'",
			ok:       false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
			ok:       true,
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: "",
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := normalizeYamlOk(tt.input)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeJsonYaml(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "Valid JSON",
			input:    `{"key": "value"}`,
			expected: `{"key":"value"}`,
		},
		{
			name:     "Valid YAML",
			input:    "key: value",
			expected: "key: value\n",
		},
		{
			name:     "Valid quoted YAML",
			input:    `"key": "value"`,
			expected: "key: value\n",
		},
		{
			name:     "Invalid JSON",
			input:    `{"key": "value"`,
			expected: "{\"key\": \"value\"",
		},
		{
			name:     "Invalid YAML",
			input:    "key: [value",
			expected: "key: [value",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeJsonYaml(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
