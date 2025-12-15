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
		// Populating the new buckets
		FCMTokens: []string{"fcm-token-1", "fcm-token-2"},
		WebSubscriptions: []notification.WebPushSubscription{
			{
				Endpoint: "https://fcm.googleapis.com/fcm/send/eR5...",
				Keys: struct {
					P256dh string `json:"p256dh"`
					Auth   string `json:"auth"`
				}{
					P256dh: "base64key...",
					Auth:   "base64auth...",
				},
			},
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

	t.Run("ToProto Strips Tokens", func(t *testing.T) {
		// 1. Convert native Go struct to Protobuf message
		protoReq := notification.NotificationRequestToProto(nativeReq)

		// Assert Recipient/Content are preserved
		require.Equal(t, "urn:contacts:user:recipient-456", protoReq.GetRecipientId())
		require.Equal(t, "New Message", protoReq.GetContent().GetTitle())

		// Assert Tokens are NOT in the proto (because the proto definition doesn't have them)
		// We can't even check `protoReq.GetTokens()` because the method doesn't exist anymore!
		// This compiles confirms the field is gone from the Proto.
	})

	t.Run("FromProto Returns Empty Buckets", func(t *testing.T) {
		// 1. Create a Proto (simulating incoming Pub/Sub message)
		protoReq := notification.NotificationRequestToProto(nativeReq)

		// 2. Convert back to Go
		convertedNative, err := notification.NotificationRequestFromProto(protoReq)
		require.NoError(t, err)

		// Assert Identity matches
		require.Equal(t, nativeReq.RecipientID, convertedNative.RecipientID)
		require.Equal(t, nativeReq.Content, convertedNative.Content)

		// Assert Buckets are Empty (expected behavior)
		require.Nil(t, convertedNative.FCMTokens)
		require.Nil(t, convertedNative.WebSubscriptions)
	})

	t.Run("Handling Nil Inputs", func(t *testing.T) {
		require.Nil(t, notification.NotificationRequestToProto(nil))

		converted, err := notification.NotificationRequestFromProto(nil)
		require.NoError(t, err)
		require.Nil(t, converted)
	})
}

// --- JSON Facade Tests ---
func TestNotificationRequest_JSON(t *testing.T) {
	nativeStruct := newTestRequest(t)

	t.Run("MarshalJSON Includes Tokens", func(t *testing.T) {
		// Even though Proto doesn't have tokens, our internal JSON logging/storage MUST have them.
		jsonBytes, err := json.Marshal(nativeStruct)
		require.NoError(t, err)

		jsonStr := string(jsonBytes)
		// Check for buckets
		assert.Contains(t, jsonStr, `"fcmTokens":["fcm-token-1","fcm-token-2"]`)
		assert.Contains(t, jsonStr, `"webSubscriptions"`)
		assert.Contains(t, jsonStr, `"endpoint":"https://fcm.googleapis.com/fcm/send/eR5..."`)
	})

	t.Run("UnmarshalJSON Restores Full Struct", func(t *testing.T) {
		jsonInput := `{
			"recipientId": "urn:contacts:user:recipient-456",
			"fcmTokens": ["fcm-token-A"],
			"webSubscriptions": [{
				"endpoint": "https://example.com",
				"keys": { "p256dh": "key", "auth": "auth" }
			}],
			"content": {
				"title": "New Message"
			}
		}`

		var resultStruct notification.NotificationRequest
		err := json.Unmarshal([]byte(jsonInput), &resultStruct)

		require.NoError(t, err)
		assert.Equal(t, "fcm-token-A", resultStruct.FCMTokens[0])
		assert.Equal(t, "https://example.com", resultStruct.WebSubscriptions[0].Endpoint)
	})
}
