// --- File: pkg/notification/v1/notification.go ---
// Package notification provides Go-native wrappers for Protobuf message types.
package notification

import (
	"encoding/json"
	"fmt"

	nv1 "github.com/tinywideclouds/gen-platform/go/types/notification/v1"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"
)

// Re-export the Protobuf types.
type NotificationRequestPb = nv1.NotificationRequestPb
type WebPushSubscriptionPb = nv1.WebPushSubscriptionPb
type NotificationRequestPbContent = nv1.NotificationRequestPb_Content

// WebPushSubscription represents a push notification token for a user's browser.
// This matches the W3C standard and the new Proto definition.
type WebPushSubscription struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

// NotificationContent holds the user-facing content of a push notification.
type NotificationContent struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Sound string `json:"sound"`
}

// NotificationRequest is the Go-native representation of a push notification job.
// It holds the "Fan-Out" buckets which are populated by the TokenStore lookup.
type NotificationRequest struct {
	RecipientID urn.URN `json:"recipientId"`

	// Bucket A: Mobile (Simple Strings for Firebase)
	FCMTokens []string `json:"fcmTokens"`

	// Bucket B: Web (Complex Objects for VAPID)
	WebSubscriptions []WebPushSubscription `json:"webSubscriptions"`

	Content     NotificationContent `json:"content"`
	DataPayload map[string]string   `json:"dataPayload"`
}

// NotificationRequestToProto converts a native NotificationRequest struct to its
// Protobuf representation.
// NOTE: This conversion is strictly for the Wire Protocol (Pub/Sub).
// Tokens are NOT included in the Proto, so they are not mapped here.
func NotificationRequestToProto(nativeReq *NotificationRequest) *NotificationRequestPb {
	if nativeReq == nil {
		return nil
	}

	return &NotificationRequestPb{
		RecipientId: nativeReq.RecipientID.String(),
		// Tokens are intentionally omitted as they are not part of the Request Proto
		Content: &nv1.NotificationRequestPb_Content{
			Title: nativeReq.Content.Title,
			Body:  nativeReq.Content.Body,
			Sound: nativeReq.Content.Sound,
		},
		DataPayload: nativeReq.DataPayload,
	}
}

// NotificationRequestFromProto converts a Protobuf NotificationRequestPb message to its
// Go-native representation.
// NOTE: The resulting struct will have empty Token buckets, meant to be populated later.
func NotificationRequestFromProto(protoReq *NotificationRequestPb) (*NotificationRequest, error) {
	if protoReq == nil {
		return nil, nil
	}

	recipientURN, err := urn.Parse(protoReq.GetRecipientId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse recipient URN: %w", err)
	}

	var nativeContent NotificationContent
	if protoReq.GetContent() != nil {
		nativeContent = NotificationContent{
			Title: protoReq.GetContent().GetTitle(),
			Body:  protoReq.GetContent().GetBody(),
			Sound: protoReq.GetContent().GetSound(),
		}
	}

	return &NotificationRequest{
		RecipientID: recipientURN,
		// Token buckets initialized as nil/empty
		FCMTokens:        nil,
		WebSubscriptions: nil,
		Content:          nativeContent,
		DataPayload:      protoReq.GetDataPayload(),
	}, nil
}

// --- JSON METHODS ---

// MarshalJSON implements the json.Marshaler interface.
// We use default JSON marshaling for the struct because the Proto no longer
// contains all the fields (tokens). If we used protojson here, we would lose
// the tokens in our logs/storage.
//
// However, if strict Proto compliance is required for the *wire* fields,
// we generally trust the struct tags defined above.
func (nr NotificationRequest) MarshalJSON() ([]byte, error) {
	// We deliberately DO NOT use protojson here for the full struct,
	// because NotificationRequest contains data (Tokens) that the Proto definition lacks.
	// We use standard Go JSON marshaling via a type alias to avoid infinite recursion.
	type Alias NotificationRequest
	return jsonMarshal(Alias(nr))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (nr *NotificationRequest) UnmarshalJSON(data []byte) error {
	type Alias NotificationRequest
	aux := &Alias{}
	if err := jsonUnmarshal(data, aux); err != nil {
		return err
	}
	*nr = NotificationRequest(*aux)
	return nil
}

func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func jsonUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
