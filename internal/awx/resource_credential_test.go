package awx

import (
	"testing"
)

func Test_getSecretFieldsFromInputs(t *testing.T) {
	tests := []struct {
		name     string
		inputs   interface{}
		expected map[string]struct{}
	}{
		{
			name:     "nil inputs",
			inputs:   nil,
			expected: map[string]struct{}{},
		},
		{
			name:     "inputs is not a map",
			inputs:   "not a map",
			expected: map[string]struct{}{},
		},
		{
			name:     "no fields key",
			inputs:   map[string]interface{}{"required": []interface{}{"foo"}},
			expected: map[string]struct{}{},
		},
		{
			name: "fields with no secret fields",
			inputs: map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{"id": "username", "label": "Username", "type": "string"},
				},
			},
			expected: map[string]struct{}{},
		},
		{
			name: "fields with secret fields",
			inputs: map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{"id": "username", "label": "Username", "type": "string"},
					map[string]interface{}{"id": "password", "label": "Password", "type": "string", "secret": true},
					map[string]interface{}{"id": "token", "label": "Token", "type": "string", "secret": true},
				},
			},
			expected: map[string]struct{}{
				"password": {},
				"token":    {},
			},
		},
		{
			name: "secret explicitly false",
			inputs: map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{"id": "password", "label": "Password", "secret": false},
				},
			},
			expected: map[string]struct{}{},
		},
		{
			name: "mixed secret and non-secret fields",
			inputs: map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{"id": "host", "label": "Host"},
					map[string]interface{}{"id": "api_key", "label": "API Key", "secret": true},
					map[string]interface{}{"id": "port", "label": "Port"},
				},
			},
			expected: map[string]struct{}{
				"api_key": {},
			},
		},
		{
			name: "field without id is skipped",
			inputs: map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{"label": "No ID", "secret": true},
				},
			},
			expected: map[string]struct{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSecretFieldsFromInputs(tt.inputs)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d secret fields, got %d", len(tt.expected), len(result))
				return
			}
			for k := range tt.expected {
				if _, ok := result[k]; !ok {
					t.Errorf("expected secret field %q not found in result", k)
				}
			}
		})
	}
}

func Test_sanitizeEncryptedInputs(t *testing.T) {
	tests := []struct {
		name         string
		apiInputs    map[string]interface{}
		stateInputs  map[string]interface{}
		secretFields map[string]struct{}
		expected     map[string]interface{}
	}{
		{
			name: "no secret fields - no changes",
			apiInputs: map[string]interface{}{
				"username": "admin",
				"host":     "example.com",
			},
			stateInputs:  map[string]interface{}{},
			secretFields: map[string]struct{}{},
			expected: map[string]interface{}{
				"username": "admin",
				"host":     "example.com",
			},
		},
		{
			name: "encrypted field replaced with state value",
			apiInputs: map[string]interface{}{
				"username": "admin",
				"password": "$encrypted$",
			},
			stateInputs: map[string]interface{}{
				"username": "admin",
				"password": "my-secret-password",
			},
			secretFields: map[string]struct{}{
				"password": {},
			},
			expected: map[string]interface{}{
				"username": "admin",
				"password": "my-secret-password",
			},
		},
		{
			name: "multiple encrypted fields replaced",
			apiInputs: map[string]interface{}{
				"host":     "example.com",
				"password": "$encrypted$",
				"token":    "$encrypted$",
			},
			stateInputs: map[string]interface{}{
				"host":     "example.com",
				"password": "secret1",
				"token":    "secret2",
			},
			secretFields: map[string]struct{}{
				"password": {},
				"token":    {},
			},
			expected: map[string]interface{}{
				"host":     "example.com",
				"password": "secret1",
				"token":    "secret2",
			},
		},
		{
			name: "secret field not in API response - no change",
			apiInputs: map[string]interface{}{
				"username": "admin",
			},
			stateInputs: map[string]interface{}{
				"username": "admin",
				"password": "secret",
			},
			secretFields: map[string]struct{}{
				"password": {},
			},
			expected: map[string]interface{}{
				"username": "admin",
			},
		},
		{
			name: "secret field not encrypted - keep API value",
			apiInputs: map[string]interface{}{
				"password": "plaintext-value",
			},
			stateInputs: map[string]interface{}{
				"password": "old-value",
			},
			secretFields: map[string]struct{}{
				"password": {},
			},
			expected: map[string]interface{}{
				"password": "plaintext-value",
			},
		},
		{
			name: "nil state inputs - encrypted field kept as is",
			apiInputs: map[string]interface{}{
				"password": "$encrypted$",
			},
			stateInputs: nil,
			secretFields: map[string]struct{}{
				"password": {},
			},
			expected: map[string]interface{}{
				"password": "$encrypted$",
			},
		},
		{
			name: "nil secretFields (fallback) - all encrypted fields replaced",
			apiInputs: map[string]interface{}{
				"username": "admin",
				"password": "$encrypted$",
				"token":    "$encrypted$",
			},
			stateInputs: map[string]interface{}{
				"username": "admin",
				"password": "secret-pw",
				"token":    "secret-tok",
			},
			secretFields: nil,
			expected: map[string]interface{}{
				"username": "admin",
				"password": "secret-pw",
				"token":    "secret-tok",
			},
		},
		{
			name: "nil secretFields - non-encrypted values preserved",
			apiInputs: map[string]interface{}{
				"username": "admin",
				"password": "$encrypted$",
			},
			stateInputs: map[string]interface{}{
				"username": "old-admin",
				"password": "secret-pw",
			},
			secretFields: nil,
			expected: map[string]interface{}{
				"username": "admin",
				"password": "secret-pw",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeEncryptedInputs(tt.apiInputs, tt.stateInputs, tt.secretFields)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d fields, got %d", len(tt.expected), len(result))
				return
			}
			for k, expectedV := range tt.expected {
				gotV, ok := result[k]
				if !ok {
					t.Errorf("expected key %q not found in result", k)
					continue
				}
				if gotV != expectedV {
					t.Errorf("field %q: expected %v, got %v", k, expectedV, gotV)
				}
			}
		})
	}
}
