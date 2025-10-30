package name

import (
	"google.golang.org/protobuf/encoding/protojson"

	userv1 "github.com/tinywideclouds/gen-platform/src/types/user/v1"
)

var (
	// protojsonMarshalOptions tells protojson to use camelCase (json_name)
	protojsonMarshalOptions = &protojson.MarshalOptions{
		UseProtoNames:   false, // <-- THIS IS THE FIX. false = use camelCase.
		EmitUnpopulated: false,
	}

	// protojsonUnmarshalOptions tells protojson to ignore unknown fields.
	protojsonUnmarshalOptions = &protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

type User struct {
	Alias string
	Name  string
	Email string
}

// ToProto converts the idiomatic Go struct into its Protobuf representation.
func ToProto(native *User) *userv1.UserPb {
	if native == nil {
		return nil
	}
	return &userv1.UserPb{
		Alias: native.Alias,
		Name:  native.Name,
		Email: native.Email,
	}
}

// FromProto converts the Protobuf representation into the idiomatic Go struct.
func FromProto(proto *userv1.UserPb) (*User, error) {
	if proto == nil {
		return nil, nil
	}

	return &User{
		Alias: proto.Alias,
		Name:  proto.Name,
		Email: proto.Email,
	}, nil
}

// --- NEW JSON METHODS ---

// MarshalJSON implements the json.Marshaler interface.
// It marshals the *Protobuf* representation, not the Go struct directly.
func (u *User) MarshalJSON() ([]byte, error) {
	// 1. Convert native Go struct to Protobuf struct
	protoPb := ToProto(u)

	// 2. Marshal the Protobuf struct using protojson
	return protojsonMarshalOptions.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It unmarshals into the *Protobuf* representation first.
func (u *User) UnmarshalJSON(data []byte) error {
	// 1. Unmarshal into a new Protobuf struct
	var protoPb userv1.UserPb

	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}

	// 2. Convert from Protobuf to native Go struct
	native, err := FromProto(&protoPb)
	if err != nil {
		return err
	}

	// 3. Assign the fields to the receiver pointer
	if native != nil {
		*u = *native
	} else {
		*u = User{}
	}

	return nil
}
