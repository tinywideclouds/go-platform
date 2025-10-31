package keys

import (
	keysv1 "github.com/tinywideclouds/gen-platform/go/types/keys/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// --- Marshal/Unmarshal Options (Unchanged) ---
var (
	protojsonMarshalOptions = &protojson.MarshalOptions{
		UseProtoNames:   false, // Use camelCase
		EmitUnpopulated: false,
	}
	protojsonUnmarshalOptions = &protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

type PublicKeys struct {
	EncKey []byte `json:"encKey,omitempty"`
	SigKey []byte `json:"sigKey,omitempty"`
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

// --- JSON METHODS ---

// MarshalJSON implements the json.Marshaler interface.
//
// REFACTOR: This now has a VALUE RECEIVER (no *).
// This means both PublicKeys and *PublicKeys satisfy the interface,
// making our API robust and removing the "fragility".
func (pk PublicKeys) MarshalJSON() ([]byte, error) {
	// 1. Convert native Go struct to Protobuf struct
	// Note: We pass a pointer to ToProto
	protoPb := ToProto(&pk)

	// 2. Marshal using our camelCase options
	return protojsonMarshalOptions.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This remains a POINTER RECEIVER (*pk), which is correct
// because it needs to modify the struct it's called on.
func (pk *PublicKeys) UnmarshalJSON(data []byte) error {
	var protoPb keysv1.PublicKeysPb

	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}

	native, err := FromProto(&protoPb)
	if err != nil {
		return err
	}

	if native != nil {
		*pk = *native
	} else {
		*pk = PublicKeys{}
	}
	return nil
}
