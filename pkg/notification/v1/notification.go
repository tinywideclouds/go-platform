// --- File: pkg/notification/v1/notification.go ---
package notification

import (
	"fmt"

	nv1 "github.com/tinywideclouds/gen-platform/go/types/notification/v1"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// Re-export the Protobuf types.
type NotificationRequestPb = nv1.NotificationRequestPb
type WebPushSubscriptionPb = nv1.WebPushSubscriptionPb

// --- Domain Structs ---

type WebPushSubscription struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh []byte `json:"p256dh"`
		Auth   []byte `json:"auth"`
	} `json:"keys"`
}

type NotificationContent struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Sound string `json:"sound"`
}

// ... (NotificationRequest remains unchanged) ...
type NotificationRequest struct {
	RecipientID      urn.URN               `json:"recipientId"`
	FCMTokens        []string              `json:"fcmTokens"`
	WebSubscriptions []WebPushSubscription `json:"webSubscriptions"`
	Content          NotificationContent   `json:"content"`
	DataPayload      map[string]string     `json:"dataPayload"`
}

// --- FACADE PATTERN IMPLEMENTATION ---

// UnmarshalJSON implements the json.Unmarshaler interface.
// It uses protojson to parse the wire format strictly, then maps it to the domain struct.
func (w *WebPushSubscription) UnmarshalJSON(data []byte) error {
	var pb nv1.WebPushSubscriptionPb

	// 1. Use protojson to parse the wire format (handling Base64, etc.)
	// We use DiscardUnknown to allow forward compatibility.
	opts := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err := opts.Unmarshal(data, &pb); err != nil {
		return err
	}

	// 2. Map Proto -> Domain
	w.Endpoint = pb.GetEndpoint()
	w.Keys.P256dh = pb.GetP256Dh()
	w.Keys.Auth = pb.GetAuth()

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
// It maps the domain struct to the Proto, then uses protojson to generate the wire format.
func (w WebPushSubscription) MarshalJSON() ([]byte, error) {
	// 1. Map Domain -> Proto
	pb := &nv1.WebPushSubscriptionPb{
		Endpoint: w.Endpoint,
		P256Dh:   w.Keys.P256dh,
		Auth:     w.Keys.Auth,
	}

	// 2. Use protojson to generate JSON
	opts := protojson.MarshalOptions{UseProtoNames: false, EmitUnpopulated: false}
	return opts.Marshal(pb)
}

// ... (Existing NotificationRequestToProto / FromProto functions remain unchanged) ...
func NotificationRequestToProto(nativeReq *NotificationRequest) *NotificationRequestPb {
	if nativeReq == nil {
		return nil
	}
	return &NotificationRequestPb{
		RecipientId: nativeReq.RecipientID.String(),
		Content: &nv1.NotificationRequestPb_Content{
			Title: nativeReq.Content.Title,
			Body:  nativeReq.Content.Body,
			Sound: nativeReq.Content.Sound,
		},
		DataPayload: nativeReq.DataPayload,
	}
}

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
		RecipientID:      recipientURN,
		FCMTokens:        nil,
		WebSubscriptions: nil,
		Content:          nativeContent,
		DataPayload:      protoReq.GetDataPayload(),
	}, nil
}
