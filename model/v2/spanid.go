// Copyright (c) 2023 The Jaeger Authors.
// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"errors"
)

const spanIDSize = 8

var (
	errMarshalSpanID   = errors.New("marshal: invalid buffer length for SpanID")
	errUnmarshalSpanID = errors.New("unmarshal: invalid SpanID length")
)

// SpanID is a custom data type that is used for all span_id fields in OTLP
// Protobuf messages.
type SpanID [spanIDSize]byte

var _ gogoCustom = (*SpanID)(nil)

// Size returns the size of the data to serialize.
func (sid SpanID) Size() int {
	if sid.IsEmpty() {
		return 0
	}
	return spanIDSize
}

// IsEmpty returns true if id contains at least one non-zero byte.
func (sid SpanID) IsEmpty() bool {
	return sid == [spanIDSize]byte{}
}

// MarshalTo converts trace ID into a binary representation.
// Called by Protobuf serialization.
func (sid SpanID) MarshalTo(data []byte) (n int, err error) {
	if sid.IsEmpty() {
		return 0, nil
	}

	if len(data) < spanIDSize {
		return 0, errMarshalSpanID
	}

	return copy(data, sid[:]), nil
}

// Marshal implements gogoCustom.
func (sid SpanID) Marshal() ([]byte, error) {
	return sid[:], nil
}

// Unmarshal inflates this trace ID from binary representation. Called by Protobuf serialization.
func (sid *SpanID) Unmarshal(data []byte) error {
	if len(data) == 0 {
		*sid = [spanIDSize]byte{}
		return nil
	}

	if len(data) != spanIDSize {
		return errUnmarshalSpanID
	}

	copy(sid[:], data)
	return nil
}

// MarshalJSON converts SpanID into a hex string enclosed in quotes.
func (sid SpanID) MarshalJSON() ([]byte, error) {
	if sid.IsEmpty() {
		return []byte(`""`), nil
	}
	return marshalJSON(sid[:])
}

// UnmarshalJSON decodes SpanID from hex string, possibly enclosed in quotes.
// Called by Protobuf JSON deserialization.
func (sid *SpanID) UnmarshalJSON(data []byte) error {
	// in base64 encoding an 8-byte array is padded to 9 bytes
	buf := [spanIDSize + 1]byte{}
	if err := unmarshalJSON(buf[:], data); err != nil {
		return err
	}
	*sid = [spanIDSize]byte{}
	copy(sid[:], buf[:spanIDSize])
	return nil
}
