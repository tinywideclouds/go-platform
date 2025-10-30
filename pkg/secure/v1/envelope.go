package secure

import (
	"fmt"

	// --- NEW IMPORTS ---
	"google.golang.org/protobuf/encoding/protojson"
	// ---
	smv1 "github.com/tinywideclouds/gen-platform/src/types/secure/v1"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"
)

// --- Marshal/Unmarshal Options ---
var (
	// protojsonMarshalOptions tells protojson to use camelCase (json_name)
	protojsonMarshalOptions = &protojson.MarshalOptions{
		UseProtoNames:   false, // Use camelCase
		EmitUnpopulated: false,
	}

	// protojsonUnmarshalOptions tells protojson to ignore unknown fields.
	protojsonUnmarshalOptions = &protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

type SecureEnvelopePb = smv1.SecureEnvelopePb
type SecureEnvelopeListPb = smv1.SecureEnvelopeListPb

// --- SecureEnvelope (Single) ---

// SecureEnvelope is the canonical, idiomatic Go struct for a message.
type SecureEnvelope struct {
	// Updated JSON tags to camelCase
	RecipientID           urn.URN `json:"recipientId"`
	EncryptedData         []byte  `json:"encryptedData,omitempty"`
	EncryptedSymmetricKey []byte  `json:"encryptedSymmetricKey,omitempty"`
	Signature             []byte  `json:"signature,omitempty"`
}

// ToProto converts the idiomatic Go struct into its Protobuf representation.
func ToProto(native *SecureEnvelope) *SecureEnvelopePb {
	if native == nil {
		return nil
	}
	return &SecureEnvelopePb{
		RecipientId:           native.RecipientID.String(),
		EncryptedData:         native.EncryptedData,
		EncryptedSymmetricKey: native.EncryptedSymmetricKey,
		Signature:             native.Signature,
	}
}

// FromProto converts the Protobuf representation into the idiomatic Go struct.
func FromProto(proto *SecureEnvelopePb) (*SecureEnvelope, error) {
	if proto == nil {
		return nil, nil
	}

	recipientID, err := urn.Parse(proto.RecipientId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recipient id: %w", err)
	}

	return &SecureEnvelope{
		RecipientID:           recipientID,
		EncryptedData:         proto.EncryptedData,
		EncryptedSymmetricKey: proto.EncryptedSymmetricKey,
		Signature:             proto.Signature,
	}, nil
}

// --- JSON METHODS (Single) ---

// MarshalJSON implements the json.Marshaler interface for SecureEnvelope.
//
// REFACTOR: This now has a VALUE RECEIVER (no *).
func (se SecureEnvelope) MarshalJSON() ([]byte, error) {
	protoPb := ToProto(&se)
	return protojsonMarshalOptions.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface for SecureEnvelope.
// This remains a POINTER RECEIVER (*se) to modify the struct.
func (se *SecureEnvelope) UnmarshalJSON(data []byte) error {
	var protoPb SecureEnvelopePb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}

	native, err := FromProto(&protoPb)
	if err != nil {
		return err
	}

	if native != nil {
		*se = *native
	} else {
		*se = SecureEnvelope{}
	}
	return nil
}

// --- SecureEnvelopeList (List) ---

// SecureEnvelopeList is the idiomatic Go struct for a list of envelopes.
type SecureEnvelopeList struct {
	Envelopes []*SecureEnvelope `json:"envelopes,omitempty"`
}

// ListToProto converts the idiomatic Go list into its Protobuf representation.
func ListToProto(native *SecureEnvelopeList) *SecureEnvelopeListPb {
	if native == nil {
		return nil
	}
	protoEnvelopes := make([]*SecureEnvelopePb, len(native.Envelopes))
	for i, env := range native.Envelopes {
		protoEnvelopes[i] = ToProto(env)
	}
	return &SecureEnvelopeListPb{
		Envelopes: protoEnvelopes,
	}
}

// ListFromProto converts the Protobuf list into the idiomatic Go struct.
func ListFromProto(proto *SecureEnvelopeListPb) (*SecureEnvelopeList, error) {
	if proto == nil {
		return nil, nil
	}
	nativeEnvelopes := make([]*SecureEnvelope, len(proto.Envelopes))
	var err error
	for i, pEnv := range proto.Envelopes {
		nativeEnvelopes[i], err = FromProto(pEnv)
		if err != nil {
			// Wrap the error with context about which envelope failed
			return nil, fmt.Errorf("failed to parse envelope at index %d: %w", i, err)
		}
	}
	return &SecureEnvelopeList{
		Envelopes: nativeEnvelopes,
	}, nil
}

// --- JSON METHODS (List) ---

// MarshalJSON implements the json.Marshaler interface for SecureEnvelopeList.
//
// REFACTOR: This now has a VALUE RECEIVER (no *).
func (sel SecureEnvelopeList) MarshalJSON() ([]byte, error) {
	protoPb := ListToProto(&sel)
	return protojsonMarshalOptions.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface for SecureEnvelopeList.
// This remains a POINTER RECEIVER (*sel) to modify the struct.
func (sel *SecureEnvelopeList) UnmarshalJSON(data []byte) error {
	var protoPb SecureEnvelopeListPb
	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}

	native, err := ListFromProto(&protoPb)
	if err != nil {
		return err
	}

	if native != nil {
		*sel = *native
	} else {
		*sel = SecureEnvelopeList{}
	}
	return nil
}
