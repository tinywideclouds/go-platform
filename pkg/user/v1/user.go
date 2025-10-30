package name

import (
	userv1 "github.com/tinywideclouds/gen-platform/src/types/user/v1"
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
