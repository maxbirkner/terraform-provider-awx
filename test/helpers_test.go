package awx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
