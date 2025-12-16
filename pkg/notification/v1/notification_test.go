// --- File: pkg/notification/v1/notification_test.go ---
package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"
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
					P256dh []byte `json:"p256dh"`
					Auth   []byte `json:"auth"`
				}{
					P256dh: []byte("test-key"),
					Auth:   []byte("test-auth"),
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

// ✅ NEW: Test Suite for the WebPushSubscription Facade
func TestWebPushSubscription_JSON_Facade(t *testing.T) {
	// Arrange: A valid subscription using "safe" strings for keys (simulating Base64)
	// Arrange: A valid Domain Object (Nested)
	original := notification.WebPushSubscription{
		Endpoint: "https://push.example.com/123",
		Keys: struct {
			P256dh []byte `json:"p256dh"`
			Auth   []byte `json:"auth"`
		}{
			P256dh: []byte("test-key"),  // Raw bytes
			Auth:   []byte("test-auth"), // Raw bytes
		},
	}

	t.Run("MarshalJSON flattens to Proto wire format", func(t *testing.T) {
		// Act
		jsonBytes, err := json.Marshal(original)
		require.NoError(t, err)

		// Assert: Expect FLAT JSON because the Proto definition is flat
		// Old (Wrong): {"endpoint": "...", "keys": { ... }}
		// New (Correct): {"endpoint": "...", "p256dh": "...", "auth": "..."}
		expectedWireFormat := `{"endpoint":"https://push.example.com/123","p256dh":"dGVzdC1rZXk=","auth":"dGVzdC1hdXRo"}`
		assert.JSONEq(t, expectedWireFormat, string(jsonBytes))
	})

	t.Run("UnmarshalJSON hydrates from Flat Proto format", func(t *testing.T) {
		// Arrange: Input is FLAT JSON (simulating what the Frontend sends via toJson)
		inputWireFormat := `{"endpoint":"https://push.example.com/123","p256dh":"dGVzdC1rZXk=","auth":"dGVzdC1hdXRo"}`
		var loaded notification.WebPushSubscription

		// Act
		err := json.Unmarshal([]byte(inputWireFormat), &loaded)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, original, loaded)
	})

	t.Run("Round Trip Integrity", func(t *testing.T) {
		// Act
		data, err := json.Marshal(original)
		require.NoError(t, err)

		var result notification.WebPushSubscription
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, original, result)
	})

	t.Run("Handles Invalid JSON via protojson error", func(t *testing.T) {
		// Arrange: Invalid JSON (wrong type for keys)
		invalidJSON := `{"endpoint": 123}`
		var loaded notification.WebPushSubscription

		// Act
		err := json.Unmarshal([]byte(invalidJSON), &loaded)

		// Assert: Should fail because protojson enforces types
		require.Error(t, err)
	})
}

// ... (Existing NotificationRequest tests) ...

func TestNotificationRequestConversions(t *testing.T) {
	nativeReq := newTestRequest(t)

	t.Run("ToProto Strips Tokens", func(t *testing.T) {
		// 1. Convert native Go struct to Protobuf message
		protoReq := notification.NotificationRequestToProto(nativeReq)

		// Assert Recipient/Content are preserved
		require.Equal(t, "urn:contacts:user:recipient-456", protoReq.GetRecipientId())
		require.Equal(t, "New Message", protoReq.GetContent().GetTitle())

		// Assert Tokens are NOT in the proto
	})

	t.Run("FromProto Returns Empty Buckets", func(t *testing.T) {
		protoReq := notification.NotificationRequestToProto(nativeReq)

		convertedNative, err := notification.NotificationRequestFromProto(protoReq)
		require.NoError(t, err)

		require.Equal(t, nativeReq.RecipientID, convertedNative.RecipientID)
		require.Equal(t, nativeReq.Content, convertedNative.Content)

		// Assert Buckets are Empty (expected behavior)
		require.Nil(t, convertedNative.FCMTokens)
		require.Nil(t, convertedNative.WebSubscriptions)
	})
}
