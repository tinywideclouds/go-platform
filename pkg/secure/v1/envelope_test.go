/*
File: envelope_test.go
Description: REFACTORED to add the missing "Proto Round Trip" test
to validate the ToProto/FromProto facade functions.
*/
package secure_test

import (
	"encoding/json" // We use the standard 'json' lib to test the interface
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// --- Import the native packages we are testing ---
	"github.com/tinywideclouds/go-platform/pkg/net/v1"
	"github.com/tinywideclouds/go-platform/pkg/secure/v1"
)

// Helper to create a valid native SecureEnvelope for tests
func newTestEnvelope(t *testing.T) *secure.SecureEnvelope {
	t.Helper()
	recipientURN, err := urn.Parse("urn:contacts:user:recipient-bob")
	require.NoError(t, err)

	return &secure.SecureEnvelope{
		RecipientID:           recipientURN,
		Priority:              0,
		EncryptedData:         []byte{1, 2, 3},
		EncryptedSymmetricKey: []byte{4, 5, 6},
		Signature:             []byte{7, 8, 9},
	}
}

// --- NEW TEST ---
func TestSecureEnvelope_Proto_RoundTrip(t *testing.T) {
	t.Run("Valid Envelope", func(t *testing.T) {
		// Arrange
		nativeEnv := newTestEnvelope(t)

		// Act: Native -> Proto
		protoPb := secure.ToProto(nativeEnv)

		// Assert: Check proto struct
		require.NotNil(t, protoPb)
		assert.Equal(t, "urn:contacts:user:recipient-bob", protoPb.RecipientId)
		assert.Equal(t, []byte{1, 2, 3}, protoPb.EncryptedData)

		// Act: Proto -> Native
		roundTripEnv, err := secure.FromProto(protoPb)
		require.NoError(t, err)

		// Assert: Check round trip
		assert.Equal(t, nativeEnv, roundTripEnv)
	})

	t.Run("Nil values", func(t *testing.T) {
		// Native -> Proto
		assert.Nil(t, secure.ToProto(nil))

		// Proto -> Native
		native, err := secure.FromProto(nil)
		require.NoError(t, err)
		assert.Nil(t, native)
	})

}

// --- EXISTING TEST (UNCHANGED, but verified) ---
func TestSecureEnvelope_JSON_RoundTrip(t *testing.T) {
	// Arrange
	nativeStruct := newTestEnvelope(t)

	// This now expects camelCase
	expectedJSON := `{
		"recipientId": "urn:contacts:user:recipient-bob",
		"encryptedData": "AQID",
		"encryptedSymmetricKey": "BAUG",
		"priority": 0,
		"signature": "BwgJ"
	}`

	// --- Test 1: Marshal (Go struct -> JSON) ---
	t.Run("MarshalJSON", func(t *testing.T) {
		// Act
		jsonBytes, err := json.Marshal(nativeStruct)
		require.NoError(t, err)

		// Assert
		assert.JSONEq(t, expectedJSON, string(jsonBytes))
	})

	// --- Test 2: Unmarshal (JSON -> Go struct) ---
	t.Run("UnmarshalJSON", func(t *testing.T) {
		// Act
		var resultStruct secure.SecureEnvelope
		err := json.Unmarshal([]byte(expectedJSON), &resultStruct)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, nativeStruct, &resultStruct)
	})

	// --- Test 3: Unmarshal (from snake_case) ---
	t.Run("UnmarshalJSON from snake_case", func(t *testing.T) {
		// Arrange
		// The protojson unmarshaler should handle both
		// camelCase (recipientId) and snake_case (recipient_id).
		jsonWithSnakeCase := `{
			"recipient_id": "urn:contacts:user:recipient-bob",
			"encrypted_data": "AQID",
			"encrypted_symmetric_key": "BAUG",
			"signature": "BwgJ"
		}`

		// Act
		var resultStruct secure.SecureEnvelope
		err := json.Unmarshal([]byte(jsonWithSnakeCase), &resultStruct)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, nativeStruct, &resultStruct)
	})
}

func TestSecureEnvelopeList_JSON_RoundTrip(t *testing.T) {
	// Arrange
	recipientURN, err := urn.Parse("urn:contacts:user:recipient-bob")
	require.NoError(t, err)

	nativeList := &secure.SecureEnvelopeList{
		Envelopes: []*secure.SecureEnvelope{
			{
				RecipientID:           recipientURN,
				Priority:              0,
				EncryptedData:         []byte{1, 2, 3},
				EncryptedSymmetricKey: []byte{4, 5, 6},
				Signature:             []byte{7, 8, 9},
			},
		},
	}

	// REFACTORED: This now expects camelCase
	expectedListJSON := `{
		"envelopes": [
			{
				"recipientId": "urn:contacts:user:recipient-bob",
				"encryptedData": "AQID",
				"priority": 0,
				"encryptedSymmetricKey": "BAUG",
				"signature": "BwgJ"
			}
		]
	}`

	// --- Test 1: Marshal (Go struct -> JSON) ---
	t.Run("MarshalJSON List", func(t *testing.T) {
		// Act
		jsonBytes, err := json.Marshal(nativeList)
		require.NoError(t, err)

		// Assert
		assert.JSONEq(t, expectedListJSON, string(jsonBytes))
	})

	// --- Test 2: Unmarshal (JSON -> Go struct) ---
	t.Run("UnmarshalJSON List", func(t *testing.T) {
		// Act
		var resultList secure.SecureEnvelopeList
		err := json.Unmarshal([]byte(expectedListJSON), &resultList)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, nativeList, &resultList)
	})
}
