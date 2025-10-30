package keys

import (
	"encoding/json" // We use the standard 'json' lib to test the interface
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicKeys_JSON_RoundTrip(t *testing.T) {
	// Arrange
	nativeStruct := &PublicKeys{
		EncKey: []byte{1, 2, 3},
		SigKey: []byte{4, 5, 6},
	}

	// This is the JSON string that protojson *should* create
	// (bytes are correctly marshaled as base64 strings)
	expectedJSON := `{"encKey":"AQID","sigKey":"BAUG"}`

	// --- Test 1: Marshal (Go struct -> JSON) ---
	t.Run("MarshalJSON", func(t *testing.T) {
		// Act
		// We call the standard json.Marshal, which will
		// automatically find our MarshalJSON() method.
		jsonBytes, err := json.Marshal(nativeStruct)
		require.NoError(t, err)

		// Assert
		assert.JSONEq(t, expectedJSON, string(jsonBytes))
	})

	// --- Test 2: Unmarshal (JSON -> Go struct) ---
	t.Run("UnmarshalJSON", func(t *testing.T) {
		// Arrange
		var resultStruct PublicKeys
		jsonBytes := []byte(expectedJSON)

		// Act
		// We call the standard json.Unmarshal, which will
		// automatically find our UnmarshalJSON() method.
		err := json.Unmarshal(jsonBytes, &resultStruct)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, nativeStruct, &resultStruct)
	})

	// --- Test 3: Unmarshal with extra fields ---
	t.Run("UnmarshalJSON with unknown fields", func(t *testing.T) {
		// Arrange
		var resultStruct PublicKeys
		// This JSON has an extra field, which our Unmarshaler should ignore
		jsonWithExtra := `{"encKey":"AQID","sigKey":"BAUG","extraField":"should_be_ignored"}`

		// Act
		err := json.Unmarshal([]byte(jsonWithExtra), &resultStruct)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, nativeStruct, &resultStruct)
	})
}
