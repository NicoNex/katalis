package katalis

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"math"
)

var (
	Uint64Codec  = uint64Codec{}
	Int64Codec   = int64Codec{}
	Float64Codec = float64Codec{}
	StringCodec  = stringCodec{}
	BytesCodec   = bytesCodec{}
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

type int64Codec struct{}

func (ic int64Codec) Encode(i int64) ([]byte, error) {
	return Uint64Codec.Encode(uint64(i))
}

func (ic int64Codec) Decode(b []byte) (int64, error) {
	i, err := Uint64Codec.Decode(b)
	return int64(i), err
}

type float64Codec struct{}

func (f64c float64Codec) Encode(f float64) ([]byte, error) {
	return Uint64Codec.Encode(math.Float64bits(f))
}

func (f64c float64Codec) Decode(b []byte) (float64, error) {
	i, err := Uint64Codec.Decode(b)
	return math.Float64frombits(i), err
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
