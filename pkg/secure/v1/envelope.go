package secure

import (
	"fmt"

	// --- NEW IMPORTS ---
	"google.golang.org/protobuf/encoding/protojson"
	// ---
	smv1 "github.com/tinywideclouds/gen-platform/src/types/secure/v1"
	urn "github.com/tinywideclouds/go-platform/pkg/net/v1"
)

type SecureEnvelopePb = smv1.SecureEnvelopePb
type SecureEnvelopeListPb = smv1.SecureEnvelopeListPb

// --- SecureEnvelope (Single) ---

// SecureEnvelope is the canonical, idiomatic Go struct for a message.
type SecureEnvelope struct {
	RecipientID           urn.URN
	EncryptedData         []byte
	EncryptedSymmetricKey []byte
	Signature             []byte
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

// --- NEW JSON METHODS (Single) ---

// MarshalJSON implements the json.Marshaler interface.
func (se *SecureEnvelope) MarshalJSON() ([]byte, error) {
	protoPb := ToProto(se)
	return protojson.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (se *SecureEnvelope) UnmarshalJSON(data []byte) error {
	var protoPb SecureEnvelopePb
	unmarshalOpts := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	if err := unmarshalOpts.Unmarshal(data, &protoPb); err != nil {
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
	Envelopes []*SecureEnvelope
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
	for i, pEnv := range proto.Envelopes {
		native, err := FromProto(pEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse envelope at index %d: %w", i, err)
		}
		nativeEnvelopes[i] = native
	}
	return &SecureEnvelopeList{
		Envelopes: nativeEnvelopes,
	}, nil
}

// --- NEW JSON METHODS (List) ---

// MarshalJSON implements the json.Marshaler interface for the list.
func (sel *SecureEnvelopeList) MarshalJSON() ([]byte, error) {
	protoPb := ListToProto(sel)
	return protojson.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface for the list.
func (sel *SecureEnvelopeList) UnmarshalJSON(data []byte) error {
	var protoPb SecureEnvelopeListPb
	unmarshalOpts := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	if err := unmarshalOpts.Unmarshal(data, &protoPb); err != nil {
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
