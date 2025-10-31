/*
File: pkg/routing/queue_test.go
Description: NEW test file to validate the "QueuedMessage" facade.
*/
package routing_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// Import the native packages we are testing and using
	"github.com/tinywideclouds/go-platform/pkg/net/v1"
	"github.com/tinywideclouds/go-platform/pkg/routing/v1"
	"github.com/tinywideclouds/go-platform/pkg/secure/v1"

	// Import the generated proto packages to check against
	routingv1 "github.com/tinywideclouds/gen-platform/go/types/routing/v1"
	securev1 "github.com/tinywideclouds/gen-platform/go/types/secure/v1"
)

// Helper to create a valid native SecureEnvelope for tests
func newTestEnvelope(t *testing.T) *secure.SecureEnvelope {
	t.Helper()
	recipientURN, err := urn.Parse("urn:sm:user:recipient-bob")
	require.NoError(t, err)

	return &secure.SecureEnvelope{
		RecipientID:           recipientURN,
		EncryptedData:         []byte{1, 2, 3},
		EncryptedSymmetricKey: []byte{4, 5, 6},
		Signature:             []byte{7, 8, 9},
	}
}

func TestQueuedMessage_Proto_RoundTrip(t *testing.T) {
	t.Run("Valid QueuedMessage", func(t *testing.T) {
		// Arrange
		nativeMsg := &routing.QueuedMessage{
			ID:       uuid.NewString(),
			Envelope: newTestEnvelope(t),
		}

		// Act: Native -> Proto
		protoPb := routing.ToProto(nativeMsg)

		// Assert: Check proto struct
		require.NotNil(t, protoPb)
		assert.Equal(t, nativeMsg.ID, protoPb.Id)
		require.NotNil(t, protoPb.Envelope)
		assert.Equal(t, nativeMsg.Envelope.RecipientID.String(), protoPb.Envelope.RecipientId)
		assert.Equal(t, nativeMsg.Envelope.EncryptedData, protoPb.Envelope.EncryptedData)

		// Act: Proto -> Native
		roundTripMsg, err := routing.FromProto(protoPb)
		require.NoError(t, err)

		// Assert: Check round trip
		assert.Equal(t, nativeMsg, roundTripMsg)
	})

	t.Run("Nil values", func(t *testing.T) {
		// Native -> Proto
		assert.Nil(t, routing.ToProto(nil))

		// Proto -> Native
		native, err := routing.FromProto(nil)
		require.NoError(t, err)
		assert.Nil(t, native)
	})

	t.Run("FromProto with invalid nested envelope", func(t *testing.T) {
		// Arrange
		badProto := &routingv1.QueuedMessagePb{
			Id: "test-id",
			Envelope: &securev1.SecureEnvelopePb{
				RecipientId: "not-a-valid-urn", // This will cause secure.FromProto to fail
			},
		}

		// Act
		_, err := routing.FromProto(badProto)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse nested envelope")
	})
}

func TestQueuedMessageList_Proto_RoundTrip(t *testing.T) {
	// Arrange
	nativeList := &routing.QueuedMessageList{
		Messages: []*routing.QueuedMessage{
			{ID: uuid.NewString(), Envelope: newTestEnvelope(t)},
			{ID: uuid.NewString(), Envelope: newTestEnvelope(t)},
		},
	}

	// Act: Native -> Proto
	protoListPb := routing.ListToProto(nativeList)

	// Assert: Check proto
	require.NotNil(t, protoListPb)
	require.Len(t, protoListPb.Messages, 2)
	assert.Equal(t, nativeList.Messages[0].ID, protoListPb.Messages[0].Id)
	assert.Equal(t, nativeList.Messages[1].ID, protoListPb.Messages[1].Id)

	// Act: Proto -> Native
	roundTripList, err := routing.ListFromProto(protoListPb)
	require.NoError(t, err)

	// Assert: Check round trip
	assert.Equal(t, nativeList, roundTripList)
}

func TestQueuedMessage_JSON_RoundTrip(t *testing.T) {
	// Arrange
	nativeMsg := &routing.QueuedMessage{
		ID:       "test-queue-id-123",
		Envelope: newTestEnvelope(t),
	}

	// This is what the client API will serve/consume
	expectedJSON := `
	{
		"id": "test-queue-id-123",
		"envelope": {
			"recipientId": "urn:sm:user:recipient-bob",
			"encryptedData": "AQID",
			"encryptedSymmetricKey": "BAUG",
			"signature": "BwgJ"
		}
	}`

	// --- Test 1: Marshal (Go struct -> JSON) ---
	t.Run("MarshalJSON", func(t *testing.T) {
		jsonBytes, err := json.Marshal(nativeMsg)
		require.NoError(t, err)
		assert.JSONEq(t, expectedJSON, string(jsonBytes))
	})

	// --- Test 2: Unmarshal (JSON -> Go struct) ---
	t.Run("UnmarshalJSON", func(t *testing.T) {
		var resultStruct routing.QueuedMessage
		err := json.Unmarshal([]byte(expectedJSON), &resultStruct)
		require.NoError(t, err)
		assert.Equal(t, nativeMsg, &resultStruct)
	})
}
