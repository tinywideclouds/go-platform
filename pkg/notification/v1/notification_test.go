package notification_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinywideclouds/go-platform/pkg/net/v1"
	"github.com/tinywideclouds/go-platform/pkg/notification/v1"
)

func TestNotificationRequestConversions(t *testing.T) {
	recipientURN, err := urn.New("contacts", "user", "recipient-456")
	require.NoError(t, err)

	nativeReq := &notification.NotificationRequest{
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

	t.Run("ToProto and FromProto Symmetry", func(t *testing.T) {
		// 1. Convert native Go struct to Protobuf message
		protoReq := notification.NotificationRequestToProto(nativeReq)

		// Assert that the Protobuf message has the correct values
		require.Equal(t, "urn:contacts:user:recipient-456", protoReq.GetRecipientId())
		require.Len(t, protoReq.GetTokens(), 2)
		require.Equal(t, "token-1", protoReq.GetTokens()[0].GetToken())
		require.Equal(t, "apns", protoReq.GetTokens()[0].GetPlatform())
		require.Equal(t, "New Message", protoReq.GetContent().GetTitle())
		require.Equal(t, "msg-789", protoReq.GetDataPayload()["message_id"])

		// 2. Convert the Protobuf message back to the native Go struct
		convertedNative, err := notification.NotificationRequestFromProto(protoReq)
		require.NoError(t, err)

		// Assert that the round-trip conversion results in the original struct
		require.Equal(t, nativeReq, convertedNative)
	})

	t.Run("Handling Nil Inputs", func(t *testing.T) {
		// ToProto should handle a nil input gracefully
		require.Nil(t, notification.NotificationRequestToProto(nil))

		// FromProto should also handle a nil input gracefully
		converted, err := notification.NotificationRequestFromProto(nil)
		require.NoError(t, err)
		require.Nil(t, converted)
	})
}
