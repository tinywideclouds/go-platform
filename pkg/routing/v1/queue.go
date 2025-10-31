/*
File: pkg/routing/queue.go
Description: REFACTORED to add the native 'QueuedMessage' and 'QueuedMessageList'
structs and their "To/FromProto" facade functions.
*/
package routing

import (
	"fmt"

	// --- NEW: Protojson for JSON methods ---
	"google.golang.org/protobuf/encoding/protojson"

	// --- NEW: Platform imports for the facade ---
	routingv1 "github.com/tinywideclouds/gen-platform/go/types/routing/v1"
	"github.com/tinywideclouds/go-platform/pkg/secure/v1"
)

// ConnectionInfo holds details about a user's real-time connection.
type ConnectionInfo struct {
	ServerInstanceID string `json:"serverInstanceId"`
	ConnectedAt      int64  `json:"connectedAt"`
}

// DeviceToken represents a push notification token for a user's device.
type DeviceToken struct {
	Token    string `json:"token"`
	Platform string `json:"platform"` // e.g., "ios", "android"
}

// --- NEW: Protojson marshal/unmarshal options (copied from envelope.go) ---
var (
	protojsonMarshalOptions = &protojson.MarshalOptions{
		UseProtoNames:   false, // Use camelCase
		EmitUnpopulated: false,
	}
	protojsonUnmarshalOptions = &protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

// --- NEW: Protobuf type aliases ---
type QueuedMessagePb = routingv1.QueuedMessagePb
type QueuedMessageListPb = routingv1.QueuedMessageListPb

// --- NEW: QueuedMessage (Single) ---

// QueuedMessage is the canonical, idiomatic Go struct for the "wrapper" message.
// It holds the internal router-generated ID and the native SecureEnvelope.
type QueuedMessage struct {
	ID       string                 `json:"id"`
	Envelope *secure.SecureEnvelope `json:"envelope"`
}

// ToProto converts the idiomatic Go struct into its Protobuf representation.
func ToProto(native *QueuedMessage) *QueuedMessagePb {
	if native == nil {
		return nil
	}
	return &QueuedMessagePb{
		Id:       native.ID,
		Envelope: secure.ToProto(native.Envelope), // Calls the envelope's facade
	}
}

// FromProto converts the Protobuf representation into the idiomatic Go struct.
func FromProto(proto *QueuedMessagePb) (*QueuedMessage, error) {
	if proto == nil {
		return nil, nil
	}

	nativeEnvelope, err := secure.FromProto(proto.Envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse nested envelope from proto: %w", err)
	}

	return &QueuedMessage{
		ID:       proto.Id,
		Envelope: nativeEnvelope,
	}, nil
}

// --- NEW: JSON Methods (Single) ---

// MarshalJSON implements the json.Marshaler interface.
func (qm QueuedMessage) MarshalJSON() ([]byte, error) {
	protoPb := ToProto(&qm)
	return protojsonMarshalOptions.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (qm *QueuedMessage) UnmarshalJSON(data []byte) error {
	var protoPb QueuedMessagePb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}
	native, err := FromProto(&protoPb)
	if err != nil {
		return err
	}
	if native != nil {
		*qm = *native
	} else {
		*qm = QueuedMessage{}
	}
	return nil
}

// --- NEW: QueuedMessageList (List) ---

// QueuedMessageList is the idiomatic Go struct for a list of queued messages.
type QueuedMessageList struct {
	Messages []*QueuedMessage `json:"messages,omitempty"`
}

// ListToProto converts the idiomatic Go list into its Protobuf representation.
func ListToProto(native *QueuedMessageList) *QueuedMessageListPb {
	if native == nil {
		return nil
	}
	protoMessages := make([]*QueuedMessagePb, len(native.Messages))
	for i, msg := range native.Messages {
		protoMessages[i] = ToProto(msg)
	}
	return &QueuedMessageListPb{
		Messages: protoMessages,
	}
}

// ListFromProto converts the Protobuf list into the idiomatic Go struct.
func ListFromProto(proto *QueuedMessageListPb) (*QueuedMessageList, error) {
	if proto == nil {
		return nil, nil
	}
	nativeMessages := make([]*QueuedMessage, len(proto.Messages))
	var err error
	for i, pMsg := range proto.Messages {
		nativeMessages[i], err = FromProto(pMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse message at index %d: %w", i, err)
		}
	}
	return &QueuedMessageList{
		Messages: nativeMessages,
	}, nil
}

// --- NEW: JSON Methods (List) ---

// MarshalJSON implements the json.Marshaler interface.
func (qml QueuedMessageList) MarshalJSON() ([]byte, error) {
	protoPb := ListToProto(&qml)
	return protojsonMarshalOptions.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (qml *QueuedMessageList) UnmarshalJSON(data []byte) error {
	var protoPb QueuedMessageListPb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}
	native, err := ListFromProto(&protoPb)
	if err != nil {
		return err
	}
	if native != nil {
		*qml = *native
	} else {
		*qml = QueuedMessageList{}
	}
	return nil
}
