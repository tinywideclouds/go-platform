package name

import (
	userv1 "github.com/tinywideclouds/gen-platform/src/types/user/v1"
	// --- NEW IMPORTS ---
	"google.golang.org/protobuf/encoding/protojson"
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

type User struct {
	// Updated JSON tags to camelCase
	Alias string `json:"alias,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
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

// --- JSON METHODS ---

// MarshalJSON implements the json.Marshaler interface.
//
// REFACTOR: This now has a VALUE RECEIVER (no *).
// This makes the marshaling robust and linter-friendly.
func (u User) MarshalJSON() ([]byte, error) {
	// 1. Convert native Go struct to Protobuf struct
	// Note: We pass a pointer to ToProto
	protoPb := ToProto(&u)

	// 2. Marshal using our camelCase options
	return protojsonMarshalOptions.Marshal(protoPb)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This remains a POINTER RECEIVER (*u), which is correct
// because it needs to modify the struct it's called on.
func (u *User) UnmarshalJSON(data []byte) error {
	var protoPb userv1.UserPb
	// Use our UnmarshalOptions
	if err := protojsonUnmarshalOptions.Unmarshal(data, &protoPb); err != nil {
		return err
	}

	native, err := FromProto(&protoPb)
	if err != nil {
		return err
	}

	if native != nil {
		*u = *native
	} else {
		*u = User{}
	}

	return nil
}
