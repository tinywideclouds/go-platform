// Package transport provides Go-native wrappers for Protobuf message types.
package notification

import (
	"fmt"

	smv1 "github.com/tinywideclouds/gen-platform/go/types/notification/v1"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"
)

// Re-export the Protobuf types with convenient aliases.
// Any project importing this 'transport' package can now use these types
// without needing to import the long Protobuf package path.
type NotificationRequestPb = smv1.NotificationRequestPb
type DeviceTokenPb = smv1.DeviceTokenPb
type NotificationRequestPbContent = smv1.NotificationRequestPb_Content

// DeviceToken represents a push notification token for a user's device.
// This is the Go-native counterpart to the DeviceTokenPb message.
type DeviceToken struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

// NotificationContent holds the user-facing content of a push notification.
type NotificationContent struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Sound string `json:"sound"`
}

// NotificationRequest is the Go-native representation of a push notification job.
// It uses idiomatic Go types like urn.URN.
type NotificationRequest struct {
	RecipientID urn.URN             `json:"recipientId"`
	Tokens      []DeviceToken       `json:"tokens"`
	Content     NotificationContent `json:"content"`
	DataPayload map[string]string   `json:"dataPayload"`
}

// NotificationRequestToProto converts a native NotificationRequest struct to its
// Protobuf representation.
func NotificationRequestToProto(nativeReq *NotificationRequest) *NotificationRequestPb {
	if nativeReq == nil {
		return nil
	}

	protoTokens := make([]*DeviceTokenPb, len(nativeReq.Tokens))
	for i, token := range nativeReq.Tokens {
		protoTokens[i] = &DeviceTokenPb{
			Token:    token.Token,
			Platform: token.Platform,
		}
	}

	return &NotificationRequestPb{
		RecipientId: nativeReq.RecipientID.String(),
		Tokens:      protoTokens,
		Content: &smv1.NotificationRequestPb_Content{
			Title: nativeReq.Content.Title,
			Body:  nativeReq.Content.Body,
			Sound: nativeReq.Content.Sound,
		},
		DataPayload: nativeReq.DataPayload,
	}
}

// NotificationRequestFromProto converts a Protobuf NotificationRequestPb message to its
// Go-native representation, parsing URNs and handling potential errors.
func NotificationRequestFromProto(protoReq *NotificationRequestPb) (*NotificationRequest, error) {
	if protoReq == nil {
		return nil, nil
	}

	recipientURN, err := urn.Parse(protoReq.GetRecipientId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse recipient URN: %w", err)
	}

	nativeTokens := make([]DeviceToken, len(protoReq.GetTokens()))
	for i, protoToken := range protoReq.GetTokens() {
		nativeTokens[i] = DeviceToken{
			Token:    protoToken.GetToken(),
			Platform: protoToken.GetPlatform(),
		}
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
		Tokens:      nativeTokens,
		Content:     nativeContent,
		DataPayload: protoReq.GetDataPayload(),
	}, nil
}
