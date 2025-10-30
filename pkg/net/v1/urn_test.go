// REFACTOR: This test is updated to validate the new getter methods and the
// bug fix in the New() constructor.
//
// ADDED: This test is now updated to validate the new behavior for
// zero-value URNs and empty strings in Parse(), String(), and UnmarshalJSON().

package urn_test

import (
	"encoding/json"
	"testing"

	"github.com/illmade-knight/go-secure-messaging/pkg/urn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewURN validates the behavior of the new constructor function.
func TestNewURN(t *testing.T) {
	t.Run("Valid URN", func(t *testing.T) {
		u, err := urn.New(urn.SecureMessaging, "user", "user-123")
		require.NoError(t, err)
		assert.Equal(t, "urn:sm:user:user-123", u.String())
		// REFACTOR: Test the new getter methods.
		assert.Equal(t, "user", u.EntityType())
		assert.Equal(t, "user-123", u.EntityID())
	})

	t.Run("Empty Namespace", func(t *testing.T) {
		_, err := urn.New("", "user", "user-123")
		require.Error(t, err)
		assert.ErrorIs(t, err, urn.ErrInvalidFormat)
	})

	t.Run("Empty Entity Type", func(t *testing.T) {
		_, err := urn.New(urn.SecureMessaging, "", "user-123")
		require.Error(t, err)
		assert.ErrorIs(t, err, urn.ErrInvalidFormat)
	})

	t.Run("Empty Entity ID", func(t *testing.T) {
		_, err := urn.New(urn.SecureMessaging, "user", "")
		require.Error(t, err)
		assert.ErrorIs(t, err, urn.ErrInvalidFormat)
	})
}

// TestString validates the String() method, including the zero-value fix.
func TestString(t *testing.T) {
	t.Run("Valid URN", func(t *testing.T) {
		u, err := urn.New(urn.SecureMessaging, "user", "user-123")
		require.NoError(t, err)
		assert.Equal(t, "urn:sm:user:user-123", u.String())
	})

	t.Run("Zero-value URN", func(t *testing.T) {
		var u urn.URN
		assert.Equal(t, "", u.String(), "A zero-value URN should serialize to an empty string")
		assert.True(t, u.IsZero())
	})
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedURN   string // We check the string representation
		expectErr     bool
		expectedErrIs error
	}{
		{
			name:        "Valid User URN",
			input:       "urn:sm:user:user-123",
			expectedURN: "urn:sm:user:user-123",
			expectErr:   false,
		},
		{
			name:        "Valid Device URN",
			input:       "urn:sm:device:uuid-abc-123",
			expectedURN: "urn:sm:device:uuid-abc-123",
			expectErr:   false,
		},
		{
			name:        "Parse Empty String (FIXED BEHAVIOR)",
			input:       "",
			expectedURN: "", // A zero-value URN's string rep is ""
			expectErr:   false,
		},
		{
			name:          "Invalid Scheme",
			input:         "foo:sm:user:user-123",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
		{
			name:          "Invalid string ':::' (too few parts)",
			input:         ":::",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
		// Parse now delegates to New, which checks all fields.
		{
			name:          "Empty Namespace",
			input:         "urn::user:user-123",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
		{
			name:          "Empty Entity Type",
			input:         "urn:sm::user-123",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
		{
			name:          "Empty Entity ID",
			input:         "urn:sm:user:",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsedURN, err := urn.Parse(tc.input)
			if tc.expectErr {
				require.Error(t, err)
				if tc.expectedErrIs != nil {
					assert.ErrorIs(t, err, tc.expectedErrIs)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedURN, parsedURN.String())
				if tc.input == "" {
					assert.True(t, parsedURN.IsZero())
				}
			}
		})
	}
}

func TestJSONMarshaling(t *testing.T) {
	u, err := urn.New(urn.SecureMessaging, "user", "user-123")
	require.NoError(t, err)
	expectedJSON := `"urn:sm:user:user-123"`

	jsonData, err := json.Marshal(u)
	require.NoError(t, err)
	assert.Equal(t, expectedJSON, string(jsonData))

	// Test marshaling a zero-value URN
	// This works because MarshalJSON has an explicit IsZero() check.
	var zeroURN urn.URN
	zeroJSON, err := json.Marshal(zeroURN)
	require.NoError(t, err)
	assert.Equal(t, "null", string(zeroJSON))
}

func TestJSONUnmarshaling(t *testing.T) {
	testCases := []struct {
		name        string
		jsonInput   string
		expectedURN string // String representation
		expectErr   bool
		expectZero  bool // Explicitly check for zero-value
	}{
		{
			name:        "Unmarshal Full URN",
			jsonInput:   `"urn:sm:user:user-123"`,
			expectedURN: "urn:sm:user:user-123",
			expectErr:   false,
		},
		{
			name:        "Unmarshal Legacy UserID (Backward Compatibility)",
			jsonInput:   `"legacy-user-456"`,
			expectedURN: "urn:sm:user:legacy-user-456",
			expectErr:   false,
		},
		{
			name:      "Unmarshal Invalid URN",
			jsonInput: `"urn:sm:user"`, // Too short
			expectErr: true,
		},
		{
			name:        "Unmarshal Empty String (FIXED BEHAVIOR)",
			jsonInput:   `""`,
			expectedURN: "",
			expectErr:   false,
			expectZero:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var u urn.URN
			err := json.Unmarshal([]byte(tc.jsonInput), &u)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedURN, u.String())
				if tc.expectZero {
					assert.True(t, u.IsZero())
				}
			}
		})
	}
}
