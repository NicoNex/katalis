package katalis

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

type Uint64Codec struct{}

func (uc Uint64Codec) Encode(i uint64) ([]byte, error) {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b, nil
}

func (uc Uint64Codec) Decode(b []byte) (uint64, error) {
	return binary.BigEndian.Uint64(b), nil
}

type StringCodec struct{}

func (sc StringCodec) Encode(s string) ([]byte, error) {
	return []byte(s), nil
}

func (sc StringCodec) Decode(b []byte) (string, error) {
	return string(b), nil
}

type GobCodec[T any] struct{}

func (pc GobCodec[T]) Encode(a T) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(a)
	return buf.Bytes(), err
}

func (pc GobCodec[T]) Decode(b []byte) (p T, err error) {
	dec := gob.NewDecoder(bytes.NewReader(b))
	err = dec.Decode(&p)
	return
}
