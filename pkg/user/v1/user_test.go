package name

import (
	"encoding/json" // We use the standard 'json' lib to test the interface
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_JSON_RoundTrip(t *testing.T) {
	// Arrange
	nativeStruct := &User{
		Alias: "Testy",
		Name:  "Test McTester",
		Email: "test@example.com",
	}

	// REFACTORED: This now expects camelCase, which matches
	// the test that was failing in the handler.
	expectedJSON := `{"alias":"Testy","name":"Test McTester","email":"test@example.com"}`

	// --- Test 1: Marshal (Go struct -> JSON) ---
	t.Run("MarshalJSON", func(t *testing.T) {
		// Act
		jsonBytes, err := json.Marshal(nativeStruct)
		require.NoError(t, err)

		// Assert
		// This assertion will now fail if protojson reverts to snake_case
		assert.JSONEq(t, expectedJSON, string(jsonBytes))
	})

	// --- Test 2: Unmarshal (JSON -> Go struct) ---
	t.Run("UnmarshalJSON", func(t *testing.T) {
		// Arrange
		var resultStruct User
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
		var resultStruct User
		jsonWithExtra := `{"alias":"Testy","name":"Test McTester","email":"test@example.com","unknown":"field"}`

		// Act
		err := json.Unmarshal([]byte(jsonWithExtra), &resultStruct)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, nativeStruct, &resultStruct)
	})

	// --- Test 4: Unmarshal with snake_case (for robustness) ---
	t.Run("UnmarshalJSON with snake_case", func(t *testing.T) {
		// Arrange
		var resultStruct User
		// Our protojson unmarshaler should handle both camelCase and snake_case
		jsonWithSnakeCase := `{"alias":"Testy","name":"Test McTester","email":"test@example.com"}`

		// Act
		err := json.Unmarshal([]byte(jsonWithSnakeCase), &resultStruct)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, nativeStruct, &resultStruct)
	})
}
