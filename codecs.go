package katalis

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"math"
)

var (
	UintCodec   = uintCodec{}
	Uint64Codec = uint64Codec{}
	Uint32Codec = uint32Codec{}
	Uint16Codec = uint16Codec{}

	IntCodec   = intCodec{}
	Int64Codec = int64Codec{}
	Int32Codec = int32Codec{}
	Int16Codec = int16Codec{}

	Float64Codec = float64Codec{}
	Float32Codec = float32Codec{}

	BytesCodec  = bytesCodec{}
	StringCodec = stringCodec{}
)

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
	return UintCodec.Encode(uint(i))
}

func (ic intCodec) Decode(b []byte) (int, error) {
	i, err := UintCodec.Decode(b)
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

type GobCodec[T any] struct{}

func (pc GobCodec[T]) Encode(a T) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(a)
	return buf.Bytes(), err
}

func (pc GobCodec[T]) Decode(b []byte) (t T, err error) {
	dec := gob.NewDecoder(bytes.NewReader(b))
	err = dec.Decode(&t)
	return
}
