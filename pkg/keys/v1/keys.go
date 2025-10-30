package keys

import (
	keysv1 "github.com/tinywideclouds/gen-platform/src/types/key/v1"
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
