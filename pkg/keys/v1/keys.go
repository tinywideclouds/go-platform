package keys

import (
	// --- NEW IMPORTS ---
	"google.golang.org/protobuf/encoding/protojson"
	// ---
	keysv1 "github.com/tinywideclouds/gen-platform/src/types/key/v1"
)

var (
	// protojsonMarshalOptions tells protojson to use camelCase (json_name)
	// which is the standard for JSON, instead of the default snake_case (proto_name).
	protojsonMarshalOptions = &protojson.MarshalOptions{
		UseProtoNames:   false, // <-- THIS IS THE FIX. false = use camelCase.
		EmitUnpopulated: false, // Don't emit empty fields
	}

	// protojsonUnmarshalOptions tells protojson to ignore fields
	// in the JSON that aren't in our proto definition.
	protojsonUnmarshalOptions = &protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

type PublicKeys struct {
	// The raw SPKI bytes of the public encryption key (RSA-OAEP).
	EncKey []byte `json:"enc_key,omitempty"`
	// The raw SPKI bytes of the public signing key (RSA-PSS).
	// This is mission-critical for the "Sealed Sender" model.
	SigKey []byte `json:"sig_key,omitempty"`
}

// ToProto converts the idiomatic Go struct into its Protobuf representation.
func ToProto(native *PublicKeys) *keysv1.PublicKeysPb {
	if native == nil {
		return nil
	}
	return &keysv1.PublicKeysPb{
		EncKey: native.EncKey,
		SigKey: native.SigKey,
	}
}

// FromProto converts the Protobuf representation into the idiomatic Go struct.
func FromProto(proto *keysv1.PublicKeysPb) (*PublicKeys, error) {
	if proto == nil {
		return nil, nil
	}

	return &PublicKeys{
		EncKey: proto.EncKey,
		SigKey: proto.SigKey,
	}, nil
}

// --- NEW JSON METHODS ---

// MarshalJSON implements the json.Marshaler interface.
// It marshals the *Protobuf* representation, not the Go struct directly.
func (pk *PublicKeys) MarshalJSON() ([]byte, error) {
	// 1. Convert native Go struct to Protobuf struct
	protoPb := ToProto(pk)

	// 2. Marshal the Protobuf struct using protojson
	// This correctly handles byte arrays as base64, etc.
	return protojsonMarshalOptions.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It unmarshals into the *Protobuf* representation first.
func (pk *PublicKeys) UnmarshalJSON(data []byte) error {
	// 1. Unmarshal into a new Protobuf struct
	var protoPb keysv1.PublicKeysPb

	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}

	// 2. Convert from Protobuf to native Go struct
	native, err := FromProto(&protoPb)
	if err != nil {
		return err
	}

	// 3. Assign the fields to the receiver pointer
	// (Handle nil case from FromProto)
	if native != nil {
		*pk = *native
	} else {
		*pk = PublicKeys{} // or set to nil if the receiver was a pointer-to-pointer
	}

	return nil
}
