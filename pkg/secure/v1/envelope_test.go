package secure

import (
	"encoding/json" // We use the standard 'json' lib to test the interface
	"testing"

	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecureEnvelope_JSON_RoundTrip(t *testing.T) {
	// Arrange
	recipientURN, err := urn.Parse("urn:sm:user:recipient-bob")
	require.NoError(t, err)

	nativeStruct := &SecureEnvelope{
		RecipientID:           recipientURN,
		EncryptedData:         []byte{1, 2, 3},
		EncryptedSymmetricKey: []byte{4, 5, 6},
		Signature:             []byte{7, 8, 9},
	}

	// This is the JSON string that protojson *should* create
	expectedJSON := `{
		"recipientId": "urn:sm:user:recipient-bob",
		"encryptedData": "AQID",
		"encryptedSymmetricKey": "BAUG",
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
		// Arrange
		var resultStruct SecureEnvelope
		jsonBytes := []byte(expectedJSON)

		// Act
		err := json.Unmarshal(jsonBytes, &resultStruct)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, nativeStruct, &resultStruct)
	})

	// --- Test 3: Unmarshal with extra fields ---
	t.Run("UnmarshalJSON with unknown fields", func(t *testing.T) {
		// Arrange
		var resultStruct SecureEnvelope
		jsonWithExtra := `{
			"recipientId": "urn:sm:user:recipient-bob",
			"encryptedData": "AQID",
			"encryptedSymmetricKey": "BAUG",
			"signature": "BwgJ",
			"unknownField": "should-be-ignored"
		}`

		// Act
		err := json.Unmarshal([]byte(jsonWithExtra), &resultStruct)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, nativeStruct, &resultStruct)
	})
}

func TestSecureEnvelopeList_JSON_RoundTrip(t *testing.T) {
	// Arrange
	recipientURN, err := urn.Parse("urn:sm:user:recipient-bob")
	require.NoError(t, err)

	nativeList := &SecureEnvelopeList{
		Envelopes: []*SecureEnvelope{
			{
				RecipientID:           recipientURN,
				EncryptedData:         []byte{1, 2, 3},
				EncryptedSymmetricKey: []byte{4, 5, 6},
				Signature:             []byte{7, 8, 9},
			},
		},
	}

	// This is the JSON string that protojson *should* create
	expectedListJSON := `{
		"envelopes": [
			{
				"recipientId": "urn:sm:user:recipient-bob",
				"encryptedData": "AQID",
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
		// Arrange
		var resultList SecureEnvelopeList
		jsonBytes := []byte(expectedListJSON)

		// Act
		err := json.Unmarshal(jsonBytes, &resultList)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, nativeList, &resultList)
	})
}
