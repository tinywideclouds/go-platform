// libs/platform/pkg/net/v1/urn.go

package urn

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	netv1 "github.com/tinywideclouds/gen-platform/go/types/net/v1"
)

const (
	// Scheme is the required scheme for all URNs in the system.
	Scheme = "urn"

	// --- Namespaces ---
	// SecureMessaging is the legacy/core namespace ("sm").
	SecureMessaging = "sm"
	// AuthNamespace is for federated identities ("auth").
	AuthNamespace = "auth"
	// LookupNamespace is for database lookup keys ("lookup").
	LookupNamespace = "lookup"

	urnParts     = 4
	urnDelimiter = ":"

	// EntityTypeUser is a standard entity type for users.
	EntityTypeUser = "user"
	// EntityTypeGroup is a standard entity type for groups.
	EntityTypeGroup = "group"
)

var (
	// ErrInvalidFormat is returned when a string or components do not conform
	// to the expected URN structure.
	ErrInvalidFormat = errors.New("invalid URN format")
)

// URN represents a parsed, validated Uniform Resource Name.
type URN struct {
	scheme     string
	namespace  string
	entityType string
	entityID   string
}

// New is the constructor for a URN. It validates that the provided namespace,
// entity type, and ID are not empty, and checks against allowed namespaces.
func New(namespace, entityType, entityID string) (URN, error) {
	if namespace == "" || entityType == "" || entityID == "" {
		return URN{}, ErrInvalidFormat
	}

	// FIX: Allow 'sm', 'auth', or 'lookup' namespaces.
	switch namespace {
	case SecureMessaging, AuthNamespace, LookupNamespace:
		// Valid
	default:
		return URN{}, fmt.Errorf("invalid namespace: expected 'sm', 'auth', or 'lookup', got '%s'", namespace)
	}

	return URN{
		scheme:     Scheme,
		namespace:  namespace,
		entityType: entityType,
		entityID:   entityID,
	}, nil
}

// Parse converts a URN string into a validated URN struct.
func Parse(s string) (URN, error) {
	// Handle empty string as a zero-value URN
	if s == "" {
		return URN{}, nil
	}

	parts := strings.Split(s, urnDelimiter)
	if len(parts) != urnParts {
		// --- Backward Compatibility for Legacy UserIDs ---
		if len(parts) == 1 {
			return New(SecureMessaging, EntityTypeUser, s)
		}
		return URN{}, fmt.Errorf("%w: expected %d parts, got %d", ErrInvalidFormat, urnParts, len(parts))
	}

	if parts[0] != Scheme {
		return URN{}, fmt.Errorf("invalid scheme: expected 'urn', got '%s'", parts[0])
	}

	// Use the New constructor to validate the namespace and parts
	return New(parts[1], parts[2], parts[3])
}

// String implements the fmt.Stringer interface.
func (u URN) String() string {
	if u.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s:%s:%s:%s", u.scheme, u.namespace, u.entityType, u.entityID)
}

// --- Getters ---

func (u URN) Namespace() string {
	return u.namespace
}

func (u URN) EntityType() string {
	return u.entityType
}

func (u URN) EntityID() string {
	return u.entityID
}

func (u URN) IsZero() bool {
	return u.scheme == "" && u.namespace == "" && u.entityType == "" && u.entityID == ""
}

// --- JSON Methods ---

func (u URN) MarshalJSON() ([]byte, error) {
	if u.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(u.String())
}

func (u *URN) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("URN should be a string, but got %s: %w", string(data), err)
	}

	if s == "" {
		*u = URN{}
		return nil
	}

	parsedURN, parseErr := Parse(s)
	if parseErr != nil {
		return parseErr
	}
	*u = parsedURN
	return nil
}

// --- Proto Methods ---

func ToProto(native URN) *netv1.UrnPb {
	if native.IsZero() {
		return nil
	}
	return &netv1.UrnPb{
		Namespace:  native.Namespace(),
		EntityType: native.EntityType(),
		EntityId:   native.EntityID(),
	}
}

func FromProto(proto *netv1.UrnPb) (URN, error) {
	if proto == nil {
		return URN{}, nil
	}
	native, err := New(proto.Namespace, proto.EntityType, proto.EntityId)
	if err != nil {
		return URN{}, fmt.Errorf("failed to convert proto to native URN: %w", err)
	}
	return native, nil
}
