// REFACTOR: This file is updated to fix the field assignment bug in the New()
// constructor. It also adds public getter methods to allow safe inspection
// of the URN's components.
//
// ADDED: This file is now updated to gracefully handle empty strings
// and zero-value URNs, making ToProto/FromProto compatible.
//
// V2 REFACTOR: Added ToProto and FromProto methods to complete the
// native façade pattern.

package urn

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	// --- NEW IMPORT ---
	netv1 "github.com/tinywideclouds/gen-platform/go/types/net/v1"
)

const (
	// Scheme is the required scheme for all URNs in the system.
	Scheme = "urn"
	// SecureMessaging is the required namespace for all URNs in the system.
	SecureMessaging = "sm"
	urnParts        = 4
	urnDelimiter    = ":"
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
// Its fields are unexported to ensure that all instances are created via the
// validating New() constructor.
type URN struct {
	scheme     string
	namespace  string
	entityType string
	entityID   string
}

// New is the constructor for a URN. It validates that the provided namespace,
// entity type, and ID are not empty, ensuring no invalid URNs can be created.
func New(namespace, entityType, entityID string) (URN, error) {
	// REFACTOR: This constructor ensures all parts are valid before assignment.
	if namespace == "" || entityType == "" || entityID == "" {
		return URN{}, ErrInvalidFormat
	}
	// We only support the 'sm' namespace for new URNs.
	if namespace != SecureMessaging {
		return URN{}, fmt.Errorf("invalid namespace: expected 'sm', got '%s'", namespace)
	}

	return URN{
		scheme:     Scheme, // Hardcoded, as it's our standard.
		namespace:  namespace,
		entityType: entityType,
		entityID:   entityID,
	}, nil
}

// Parse converts a URN string into a validated URN struct.
//
// FIXED: This function now gracefully handles an empty string by returning
// a zero-value URN, which makes it compatible with ToProto/FromProto.
func Parse(s string) (URN, error) {
	// Handle empty string as a zero-value URN
	if s == "" {
		return URN{}, nil
	}

	parts := strings.Split(s, urnDelimiter)
	if len(parts) != urnParts {
		// --- Backward Compatibility for Legacy UserIDs ---
		// If it's not a URN, check if it's a legacy ID.
		// We assume a legacy ID has no colons.
		if len(parts) == 1 {
			// It's a legacy ID, auto-migrate it to the "user" type.
			return New(SecureMessaging, EntityTypeUser, s)
		}
		return URN{}, fmt.Errorf("%w: expected %d parts, got %d", ErrInvalidFormat, urnParts, len(parts))
	}

	if parts[0] != Scheme {
		return URN{}, fmt.Errorf("invalid scheme: expected 'urn', got '%s'", parts[0])
	}
	if parts[1] != SecureMessaging {
		return URN{}, fmt.Errorf("invalid namespace: expected 'sm', got '%s'", parts[1])
	}
	if parts[2] == "" {
		return URN{}, fmt.Errorf("entity type must not be empty")
	}
	if parts[3] == "" {
		return URN{}, fmt.Errorf("entity ID must not be empty")
	}

	return URN{
		scheme:     parts[0],
		namespace:  parts[1],
		entityType: parts[2],
		entityID:   parts[3],
	}, nil
}

// String implements the fmt.Stringer interface.
//
// FIXED: This now gracefully handles a zero-value URN by returning an
// empty string, which is the expected behavior for ToProto.
func (u URN) String() string {
	if u.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s:%s:%s:%s", u.scheme, u.namespace, u.entityType, u.entityID)
}

// --- Getters ---

// Namespace returns the namespace (e.g., "sm").
func (u URN) Namespace() string {
	return u.namespace
}

// EntityType returns the type of the entity (e.g., "user", "device").
func (u URN) EntityType() string {
	return u.entityType
}

// EntityID returns the unique identifier for the entity.
func (u URN) EntityID() string {
	return u.entityID
}

// IsZero returns true if the URN has not been initialized.
func (u URN) IsZero() bool {
	return u.scheme == "" && u.namespace == "" && u.entityType == "" && u.entityID == ""
}

// --- JSON Methods (Unchanged) ---

// MarshalJSON implements the json.Marshaler interface.
func (u URN) MarshalJSON() ([]byte, error) {
	if u.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(u.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
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

// --- NEW PROTO METHODS ---

// ToProto converts the idiomatic Go URN struct into its Protobuf representation.
// Note: The "scheme" is omitted as it's implied.
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

// FromProto converts the Protobuf representation into the idiomatic Go URN struct.
func FromProto(proto *netv1.UrnPb) (URN, error) {
	if proto == nil {
		return URN{}, nil
	}

	// Use the New() constructor to ensure validation logic is applied.
	native, err := New(proto.Namespace, proto.EntityType, proto.EntityId)
	if err != nil {
		return URN{}, fmt.Errorf("failed to convert proto to native URN: %w", err)
	}
	return native, nil
}
