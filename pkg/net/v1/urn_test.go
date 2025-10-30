// REFACTOR: This test is updated to validate the new getter methods and the
// bug fix in the New() constructor.
//
// ADDED: This test is now updated to validate the new behavior for
// zero-value URNs and empty strings in Parse(), String(), and UnmarshalJSON().
//
// V2 REFACTOR: Added Proto round trip tests.

package urn_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"

	// --- NEW IMPORT ---
	netv1 "github.com/tinywideclouds/gen-platform/src/types/net/v1"
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

	t.Run("Invalid Namespace", func(t *testing.T) {
		_, err := urn.New("other", "user", "user-123")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid namespace: expected 'sm'")
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

// TestParseURN validates the string parsing logic.
func TestParseURN(t *testing.T) {
	t.Run("Valid URN", func(t *testing.T) {
		u, err := urn.Parse("urn:sm:user:user-123")
		require.NoError(t, err)
		assert.Equal(t, "urn:sm:user:user-123", u.String())
		assert.Equal(t, "user", u.EntityType())
		assert.Equal(t, "user-123", u.EntityID())
	})

	t.Run("Invalid Scheme", func(t *testing.T) {
		_, err := urn.Parse("http:sm:user:user-123")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid scheme")
	})

	t.Run("Invalid Namespace", func(t *testing.T) {
		_, err := urn.Parse("urn:other:user:user-123")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid namespace")
	})

	t.Run("Invalid Format - Too Few Parts", func(t *testing.T) {
		_, err := urn.Parse("urn:sm:user")
		require.Error(t, err)
		assert.ErrorIs(t, err, urn.ErrInvalidFormat)
	})

	t.Run("Invalid Format - Empty Entity", func(t *testing.T) {
		_, err := urn.Parse("urn:sm::user-123")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "entity type must not be empty")
	})

	t.Run("Backward Compatibility - Legacy UserID", func(t *testing.T) {
		u, err := urn.Parse("legacy-user-456")
		require.NoError(t, err)
		assert.Equal(t, "urn:sm:user:legacy-user-456", u.String())
		assert.Equal(t, "user", u.EntityType())
		assert.Equal(t, "legacy-user-456", u.EntityID())
	})

	// FIXED: Test for zero-value/empty string behavior
	t.Run("Parse Empty String", func(t *testing.T) {
		u, err := urn.Parse("")
		require.NoError(t, err)
		assert.True(t, u.IsZero())
		assert.Equal(t, "", u.String())
	})
}

// TestStringer verifies the String() method behavior.
func TestStringer(t *testing.T) {
	u, err := urn.New(urn.SecureMessaging, "group", "abc-def")
	require.NoError(t, err)
	assert.Equal(t, "urn:sm:group:abc-def", u.String())

	// FIXED: Test zero-value behavior
	var zeroURN urn.URN
	assert.Equal(t, "", zeroURN.String())
}

// TestJSONMarshaling verifies the custom JSON marshaler.
func TestJSONMarshaling(t *testing.T) {
	u, err := urn.New(urn.SecureMessaging, "user", "user-123")
	require.NoError(t, err)

	// Wrap in a struct for a realistic test
	data := struct {
		UserURN urn.URN `json:"userUrn"`
	}{UserURN: u}

	expectedJSON := `{"userUrn":"urn:sm:user:user-123"}`
	jsonData, err := json.Marshal(data)
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
			expectZero:  true, // Should be a zero URN
		},
		{
			name:       "Unmarshal null (FIXED BEHAVIOR)",
			jsonInput:  `null`,
			expectErr:  false,
			expectZero: true, // Should be a zero URN
		},
		{
			name:      "Unmarshal Invalid Type (Number)",
			jsonInput: `123`,
			expectErr: true,
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

// --- NEW PROTO TEST ---

func TestURN_Proto_RoundTrip(t *testing.T) {
	t.Run("Valid URN", func(t *testing.T) {
		// Arrange
		nativeURN, err := urn.New(urn.SecureMessaging, "user", "proto-test-123")
		require.NoError(t, err)

		// Act: Native -> Proto
		protoPb := urn.ToProto(nativeURN)

		// Assert: Check proto struct
		require.NotNil(t, protoPb)
		assert.Equal(t, "sm", protoPb.Namespace)
		assert.Equal(t, "user", protoPb.EntityType)
		assert.Equal(t, "proto-test-123", protoPb.EntityId)

		// Act: Proto -> Native
		roundTripURN, err := urn.FromProto(protoPb)
		require.NoError(t, err)

		// Assert: Check round trip
		assert.Equal(t, nativeURN, roundTripURN)
	})

	t.Run("Zero URN", func(t *testing.T) {
		// Arrange
		var zeroURN urn.URN

		// Act: Native -> Proto
		protoPb := urn.ToProto(zeroURN)

		// Assert
		assert.Nil(t, protoPb)

		// Act: Proto -> Native
		roundTripURN, err := urn.FromProto(protoPb)
		require.NoError(t, err)

		// Assert
		assert.True(t, roundTripURN.IsZero())
	})

	t.Run("Invalid Proto -> Native (validation fail)", func(t *testing.T) {
		// Arrange
		invalidProto := &netv1.UrnPb{
			Namespace:  "invalid-namespace", // This should fail our New() validator
			EntityType: "user",
			EntityId:   "test",
		}

		// Act
		_, err := urn.FromProto(invalidProto)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid namespace")
	})
}
