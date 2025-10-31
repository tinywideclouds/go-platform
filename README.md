# Go Platform Library (go-platform)
This repository provides the canonical, idiomatic Go façades for our platform's shared data types. It acts as a crucial abstraction layer between our Go microservices and the Protobuf-generated code (gen-platform).
The core purpose of this library is to provide robust, safe, and easy-to-use native Go structs.
It solves two primary problems:
Separation of Concerns: Our services should not be tightly coupled to generated *Pb structs. This library provides clean, native Go types (like urn.URN, keys.PublicKeys) for services to use.
Robust JSON Serialization: This library guarantees that all our Go types serialize to and from JSON as camelCase (the API standard), not snake_case (the Protobuf standard).

## The Façade Pattern

This library follows a strict façade pattern. Each native Go struct (e.g., keys.PublicKeys) is a "façade" over its corresponding Protobuf struct (e.g., keysv1.PublicKeysPb).
Each façade provides a standard set of methods:
ToProto(native *MyStruct) *MyStructPb: Converts the native Go struct to its Protobuf representation.
FromProto(proto *MyStructPb) (*MyStruct, error): Converts the Protobuf struct back to the native Go struct, performing validation if necessary (like urn.Parse).
func (s MyStruct) MarshalJSON() ([]byte, error): Implements the json.Marshaler interface.
func (s *MyStruct) UnmarshalJSON(data []byte) error: Implements the json.Unmarshaler interface.

## The JSON Marshaling Contract (Robustness)

This is the most critical feature of the library. We enforce a specific pattern to prevent bugs and ensure consistency.
MarshalJSON has a VALUE receiver:
````
// Note: (pk PublicKeys), not (pk *PublicKeys)
func (pk PublicKeys) MarshalJSON() ([]byte, error) {
// ...
}
````
This is intentional. By using a value receiver, both a PublicKeys and a *PublicKeys satisfy the json.Marshaler interface. This makes our API handlers robust: json.Encode(myStruct) will work correctly whether myStruct is a pointer or a value.
UnmarshalJSON has a POINTER receiver:
Go
````
func (pk *PublicKeys) UnmarshalJSON(data []byte) error {
// ...
}
````
This is required by Go to modify the struct the method is called on.
protojson is used internally:
All MarshalJSON methods use protojson.MarshalOptions{ UseProtoNames: false } to guarantee camelCase output.

##Usage Examples


In an HTTP Handler

You should always use the standard encoding/json library. Our façade methods will be called automatically.

Go

````
// handler.go
import (
"encoding/json"
"net/http"
"github.com/tinywideclouds/go-platform/pkg/keys/v1"
"github.com/tinywideclouds/go-platform/pkg/net/v1"
)

func GetPublicKeysHandler(w http.ResponseWriter, r *http.Request) {
// ...
entityURN, err := urn.Parse(r.PathValue("entityURN"))
// ...

    // 1. The store returns the *native* struct
    retrievedKeys, err := myStore.GetPublicKeys(r.Context(), entityURN)
    if err != nil {
        // ...
        return
    }

    // 2. Use standard json.Encode.
    // This will *automatically* call our (pk PublicKeys) MarshalJSON() method,
    // which correctly uses protojson and camelCase.
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(retrievedKeys)
    
    // This is robust. It would also work if retrievedKeys was a pointer:
    // json.NewEncoder(w).Encode(&retrievedKeys)
}

func StorePublicKeysHandler(w http.ResponseWriter, r *http.Request) {
// ...

    // 1. Use standard json.Decode.
    // This will *automatically* call our (pk *PublicKeys) UnmarshalJSON() method.
    var keysToStore keys.PublicKeys
    if err := json.NewDecoder(r.Body).Decode(&keysToStore); err != nil {
        // ...
        return
    }
    
    // 2. The struct is now populated and validated, ready for use.
    if err := myStore.StorePublicKeys(r.Context(), entityURN, keysToStore); err != nil {
        // ...
        return
    }
    
    w.WriteHeader(http.StatusCreated)
}
````


Internal Service Logic

When passing data between internal services (e.g., to a gRPC client), use the ToProto and FromProto helpers.

Go

````
// my_service.go
import (
"github.com/tinywideclouds/go-platform/pkg/keys/v1"
"github.com/tinywideclouds/go-platform/pkg/net/v1"
keysv1 "github.com/tinywideclouds/gen-platform/src/types/key/v1"
)

func (s *MyService) DoSomething(nativeURN urn.URN) (*keys.PublicKeys, error) {
// 1. Convert native URN to a *UrnPb for the gRPC request
protoURN := urn.ToProto(nativeURN)

    // 2. Make the gRPC call
    protoKeys, err := s.grpcClient.GetKeys(ctx, &keysv1.GetKeysRequest{Urn: protoURN})
    if err != nil {
        return nil, err
    }

    // 3. Convert the *PublicKeysPb response back to a native struct
    nativeKeys, err := keys.FromProto(protoKeys)
    if err != nil {
        return nil, err
    }
    
    return nativeKeys, nil
}
````


### Available Packages

pkg/net/v1: Provides the smart urn.URN struct, which handles parsing, validation, and string formatting.
pkg/keys/v1: Provides the keys.PublicKeys struct used for the "Sealed Sender" model.
pkg/secure/v1: Provides the secure.SecureEnvelope and secure.SecureEnvelopeList façades for the E2EE wrapper.
pkg/name/v1: Provides the name.User struct for user profile information.

### Contributing

When adding a new façade, you must adhere to this pattern:
Add the ToProto and FromProto functions.
Add the (s MyStruct) MarshalJSON() (value receiver) and (s *MyStruct) UnmarshalJSON(...) (pointer receiver) methods.
Ensure the JSON methods use the protojson library with the correct camelCase options.
Add a _test.go file validating both the Proto and JSON round-trip functionality.
