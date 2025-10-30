// REFACTOR: This file is updated to fix the field assignment bug in the New()
// constructor. It also adds public getter methods to allow safe inspection
// of the URN's components.
//
// ADDED: This file is now updated to gracefully handle empty strings
// and zero-value URNs, making ToProto/FromProto compatible.

package urn

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	if namespace == "" {
		return URN{}, fmt.Errorf("%w: namespace cannot be empty", ErrInvalidFormat)
	}
	if entityType == "" {
		return URN{}, fmt.Errorf("%w: entity type cannot be empty", ErrInvalidFormat)
	}
	if entityID == "" {
		return URN{}, fmt.Errorf("%w: entity ID cannot be empty", ErrInvalidFormat)
	}
	// REFACTOR: Corrected the field assignments.
	return URN{
		scheme:     Scheme,
		namespace:  namespace,
		entityType: entityType,
		entityID:   entityID,
	}, nil
}

// Parse converts a raw string into a structured URN, validating its format.
//
// FIXED: Now correctly handles an empty string by returning a zero-value URN
// instead of an error. This is required for transport.FromProto.
func Parse(s string) (URN, error) {
	// Handle the zero-value case
	if s == "" {
		return URN{}, nil
	}

	parts := strings.Split(s, urnDelimiter)
	if len(parts) != urnParts {
		return URN{}, fmt.Errorf("%w: expected %d parts, but got %d", ErrInvalidFormat, urnParts, len(parts))
	}

	if parts[0] != Scheme {
		return URN{}, fmt.Errorf("%w: invalid scheme '%s', expected '%s'", ErrInvalidFormat, parts[0], Scheme)
	}

	// Delegate final validation to the constructor.
	return New(parts[1], parts[2], parts[3])
}

// String reassembles the URN into its canonical string representation.
//
// FIXED: A zero-value URN now serializes to an empty string ""
// instead of ":::". This is required for transport.ToProto.
func (u URN) String() string {
	if u.IsZero() {
		return ""
	}
	return strings.Join([]string{u.scheme, u.namespace, u.entityType, u.entityID}, urnDelimiter)
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

// MarshalJSON implements the json.Marshaler interface.
//
// FIXED: This now correctly handles the `u.String()` returning "" for a zero-value
// URN, by explicitly checking IsZero() first.
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

	// FIXED: This logic now works because Parse("") returns a valid,
	// zero-value URN instead of an error.
	if strings.HasPrefix(s, Scheme+urnDelimiter) {
		parsedURN, parseErr := Parse(s)
		if parseErr != nil {
			return parseErr
		}
		*u = parsedURN
		return nil
	}

	// FIXED: Handle the empty string case from JSON
	if s == "" {
		*u = URN{}
		return nil
	}

	if s != "" {
		legacyURN, err := New(SecureMessaging, EntityTypeUser, s)
		if err != nil {
			return err
		}
		*u = legacyURN
		return nil
	}

	return ErrInvalidFormat
}
