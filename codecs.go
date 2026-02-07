package katalis

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"math"
)

// Predefined codecs for common Go types. These codecs handle encoding and
// decoding of primitive types to/from byte slices for database storage.
var (
	UintCodec   = uintCodec{}
	Uint16Codec = uint16Codec{}
	Uint32Codec = uint32Codec{}
	Uint64Codec = uint64Codec{}

	IntCodec   = intCodec{}
	Int16Codec = int16Codec{}
	Int32Codec = int32Codec{}
	Int64Codec = int64Codec{}

	Float64Codec = float64Codec{}
	Float32Codec = float32Codec{}

	BytesCodec  = bytesCodec{}
	StringCodec = stringCodec{}
)

// codecFor returns the appropriate codec for type T. For primitive types, it
// returns the corresponding predefined codec. For custom types, it returns a
// GobCodec.
func codecFor[T any]() Codec[T] {
	switch any(zero[T]()).(type) {
	case string:
		return any(StringCodec).(Codec[T])
	case []byte:
		return any(BytesCodec).(Codec[T])
	case uint:
		return any(UintCodec).(Codec[T])
	case uint16:
		return any(Uint16Codec).(Codec[T])
	case uint32:
		return any(Uint32Codec).(Codec[T])
	case uint64:
		return any(Uint64Codec).(Codec[T])
	case int:
		return any(IntCodec).(Codec[T])
	case int16:
		return any(Int16Codec).(Codec[T])
	case int32:
		return any(Int32Codec).(Codec[T])
	case int64:
		return any(Int64Codec).(Codec[T])
	case float32:
		return any(Float32Codec).(Codec[T])
	case float64:
		return any(Float64Codec).(Codec[T])
	default:
		return GobCodec[T]{}
	}
}

type uint64Codec struct{}

func (uc uint64Codec) Encode(i uint64) ([]byte, error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b, nil
}

func (uc uint64Codec) Decode(b []byte) (uint64, error) {
	return binary.BigEndian.Uint64(b), nil
}

type uint32Codec struct{}

func (uc uint32Codec) Encode(i uint32) ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b, nil
}

func (uc uint32Codec) Decode(b []byte) (uint32, error) {
	return binary.BigEndian.Uint32(b), nil
}

type uint16Codec struct{}

func (uc uint16Codec) Encode(i uint16) ([]byte, error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return b, nil
}

func (uc uint16Codec) Decode(b []byte) (uint16, error) {
	return binary.BigEndian.Uint16(b), nil
}

type uintCodec struct{}

func (uc uintCodec) Encode(i uint) ([]byte, error) {
	return Uint32Codec.Encode(uint32(i))
}

func (uc uintCodec) Decode(b []byte) (uint, error) {
	u32, err := Uint32Codec.Decode(b)
	return uint(u32), err
}

type int64Codec struct{}

func (ic int64Codec) Encode(i int64) ([]byte, error) {
	return Uint64Codec.Encode(uint64(i))
}

func (ic int64Codec) Decode(b []byte) (int64, error) {
	i, err := Uint64Codec.Decode(b)
	return int64(i), err
}

type int32Codec struct{}

func (ic int32Codec) Encode(i int32) ([]byte, error) {
	return Uint32Codec.Encode(uint32(i))
}

func (ic int32Codec) Decode(b []byte) (int32, error) {
	i, err := Uint32Codec.Decode(b)
	return int32(i), err
}

type int16Codec struct{}

func (ic int16Codec) Encode(i int16) ([]byte, error) {
	return Uint16Codec.Encode(uint16(i))
}

func (ic int16Codec) Decode(b []byte) (int16, error) {
	i, err := Uint16Codec.Decode(b)
	return int16(i), err
}

type intCodec struct{}

func (ic intCodec) Encode(i int) ([]byte, error) {
	return Int32Codec.Encode(int32(i))
}

func (ic intCodec) Decode(b []byte) (int, error) {
	i, err := Int32Codec.Decode(b)
	return int(i), err
}

type float64Codec struct{}

func (f64c float64Codec) Encode(f float64) ([]byte, error) {
	return Uint64Codec.Encode(math.Float64bits(f))
}

func (f32c float64Codec) Decode(b []byte) (float64, error) {
	i, err := Uint64Codec.Decode(b)
	return math.Float64frombits(i), err
}

type float32Codec struct{}

func (f32c float32Codec) Encode(f float32) ([]byte, error) {
	return Uint32Codec.Encode(math.Float32bits(f))
}

func (f32c float32Codec) Decode(b []byte) (float32, error) {
	i, err := Uint32Codec.Decode(b)
	return math.Float32frombits(i), err
}

type stringCodec struct{}

func (sc stringCodec) Encode(s string) ([]byte, error) {
	return []byte(s), nil
}

func (sc stringCodec) Decode(b []byte) (string, error) {
	return string(b), nil
}

type bytesCodec struct{}

func (sc bytesCodec) Encode(b []byte) ([]byte, error) {
	return b, nil
}

func (sc bytesCodec) Decode(b []byte) ([]byte, error) {
	return b, nil
}

// Gob returns a GobCodec for type T. The optional variadic parameter allows
// type inference from a value.
func Gob[T any](_ ...T) (g GobCodec[T]) { return }

// GobCodec is a generic codec that uses Go's gob encoding to serialize values.
// It works with any type that can be encoded by the encoding/gob package.
type GobCodec[T any] struct{}

// Encode serializes the value using gob encoding.
func (pc GobCodec[T]) Encode(a T) ([]byte, error) {
	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(a)
	return buf.Bytes(), err
}

// Decode deserializes the value using gob decoding.
func (pc GobCodec[T]) Decode(b []byte) (t T, err error) {
	return t, gob.NewDecoder(bytes.NewReader(b)).Decode(&t)
}
