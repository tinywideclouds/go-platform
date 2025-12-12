// --- File: pkg/notification/v1/notification_test.go ---
package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinywideclouds/go-platform/pkg/net/v1"
	"github.com/tinywideclouds/go-platform/pkg/notification/v1"
)

// newTestRequest creates a populated NotificationRequest for testing.
func newTestRequest(t *testing.T) *notification.NotificationRequest {
	t.Helper()
	recipientURN, err := urn.New("contacts", "user", "recipient-456")
	require.NoError(t, err)

	return &notification.NotificationRequest{
		RecipientID: recipientURN,
		Tokens: []notification.DeviceToken{
			{Token: "token-1", Platform: "apns"},
			{Token: "token-2", Platform: "fcm"},
		},
		Content: notification.NotificationContent{
			Title: "New Message",
			Body:  "You have a new secure message.",
			Sound: "default",
		},
		DataPayload: map[string]string{
			"message_id": "msg-789",
		},
	}
}

func TestNotificationRequestConversions(t *testing.T) {
	nativeReq := newTestRequest(t)

	t.Run("ToProto and FromProto Symmetry", func(t *testing.T) {
		// 1. Convert native Go struct to Protobuf message
		protoReq := notification.NotificationRequestToProto(nativeReq)

		// Assert that the Protobuf message has the correct values
		require.Equal(t, "urn:contacts:user:recipient-456", protoReq.GetRecipientId())
		require.Len(t, protoReq.GetTokens(), 2)
		require.Equal(t, "token-1", protoReq.GetTokens()[0].GetToken())
		require.Equal(t, "New Message", protoReq.GetContent().GetTitle())

		// 2. Convert the Protobuf message back to the native Go struct
		convertedNative, err := notification.NotificationRequestFromProto(protoReq)
		require.NoError(t, err)

		// Assert that the round-trip conversion results in the original struct
		require.Equal(t, nativeReq, convertedNative)
	})

	t.Run("Handling Nil Inputs", func(t *testing.T) {
		require.Nil(t, notification.NotificationRequestToProto(nil))

		converted, err := notification.NotificationRequestFromProto(nil)
		require.NoError(t, err)
		require.Nil(t, converted)
	})
}

// --- NEW TEST: JSON Facade Round Trip ---
func TestNotificationRequest_JSON_RoundTrip(t *testing.T) {
	nativeStruct := newTestRequest(t)

	// We expect camelCase keys because protojson defaults are overridden by our variables
	// AND the proto definition typically uses camelCase for JSON mapping.
	expectedJSONSubstring := `"recipientId":"urn:contacts:user:recipient-456"`

	// --- Test 1: Marshal (Go struct -> JSON) ---
	t.Run("MarshalJSON", func(t *testing.T) {
		// Act: Standard JSON marshal should trigger our custom MarshalJSON
		jsonBytes, err := json.Marshal(nativeStruct)
		require.NoError(t, err)

		jsonStr := string(jsonBytes)
		assert.Contains(t, jsonStr, expectedJSONSubstring)
		assert.Contains(t, jsonStr, `"title":"New Message"`)
	})

	// --- Test 2: Unmarshal (JSON -> Go struct) ---
	t.Run("UnmarshalJSON", func(t *testing.T) {
		// Arrange: Create JSON manually
		jsonInput := `{
			"recipientId": "urn:contacts:user:recipient-456",
			"tokens": [
				{"token": "token-1", "platform": "apns"},
				{"token": "token-2", "platform": "fcm"}
			],
			"content": {
				"title": "New Message",
				"body": "You have a new secure message.",
				"sound": "default"
			},
			"dataPayload": {
				"message_id": "msg-789"
			}
		}`

		// Act
		var resultStruct notification.NotificationRequest
		err := json.Unmarshal([]byte(jsonInput), &resultStruct)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, nativeStruct, &resultStruct)
	})

	// --- Test 3: Unmarshal (from snake_case) ---
	t.Run("UnmarshalJSON from snake_case", func(t *testing.T) {
		// Arrange: Protojson should handle snake_case inputs gracefully
		jsonWithSnakeCase := `{
			"recipient_id": "urn:contacts:user:recipient-456",
			"tokens": [
				{"token": "token-1", "platform": "apns"}
			],
			"content": {
				"title": "New Message"
			},
			"data_payload": {
				"message_id": "msg-789"
			}
		}`

		// Act
		var resultStruct notification.NotificationRequest
		err := json.Unmarshal([]byte(jsonWithSnakeCase), &resultStruct)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "urn:contacts:user:recipient-456", resultStruct.RecipientID.String())
		assert.Equal(t, "msg-789", resultStruct.DataPayload["message_id"])
	})
}
