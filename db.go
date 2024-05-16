package katalis

import (
	"errors"
	"io"
	"time"

	"github.com/akrylysov/pogreb"
)

type Codec[T any] interface {
	Encode(T) ([]byte, error)
	Decode([]byte) (T, error)
}

type DB[KT, VT any] struct {
	*pogreb.DB
	path     string
	keyCodec Codec[KT]
	valCodec Codec[VT]
}

type Options = pogreb.Options

var ErrIterationDone = pogreb.ErrIterationDone

// Open opens or creates a new DB. The DB must be closed after use, by calling
// Close method.
func Open[KT, VT any](path string, keyCodec Codec[KT], valCodec Codec[VT]) (db DB[KT, VT], err error) {
	pg, err := pogreb.Open(path, &pogreb.Options{
		BackgroundSyncInterval:       -1,
		BackgroundCompactionInterval: time.Hour * 24,
	})

	db = DB[KT, VT]{
		DB:       pg,
		path:     path,
		keyCodec: keyCodec,
		valCodec: valCodec,
	}
	return
}

// OpenOptions is like Open but accepts an Options struct.
func OpenOptions[KT, VT any](path string, keyCodec Codec[KT], valCodec Codec[VT], opts *Options) (db DB[KT, VT], err error) {
	pg, err := pogreb.Open(path, opts)

	db = DB[KT, VT]{
		DB:       pg,
		path:     path,
		keyCodec: keyCodec,
		valCodec: valCodec,
	}
	return
}

// Get returns the value for the given key stored in the DB or an empty value if
// the key doesn't exist.
func (db DB[KT, VT]) Get(key KT) (res VT, err error) {
	// Encode the key to []byte.
	kb, err := db.keyCodec.Encode(key)
	if err != nil {
		return res, err
	}

	// Fetch from the DB the []byte of the value.
	b, err := db.DB.Get(kb)
	if err != nil {
		return res, err
	}
	// Decode the value into its type.
	return db.valCodec.Decode(b)
}

// Put sets the value for the given key. It updates the value for the existing key.
func (db DB[KT, VT]) Put(key KT, val VT) error {
	// Encode the key to []byte.
	kb, err := db.keyCodec.Encode(key)
	if err != nil {
		return err
	}

	// Encode the value to []byte.
	vb, err := db.valCodec.Encode(val)
	if err != nil {
		return err
	}
	// Write in the DB the key and the value both as []byte.
	return db.DB.Put(kb, vb)
}

func (db DB[KT, VT]) Del(key KT) error {
	// Encode the key to []byte.
	kb, err := db.keyCodec.Encode(key)
	if err != nil {
		return err
	}
	// Delete from the DB the key-value pair.
	return db.DB.Delete(kb)
}

// Has returns true if the DB contains the given key.
func (db DB[KT, VT]) Has(key KT) (bool, error) {
	// Encode the key to []byte.
	kb, err := db.keyCodec.Encode(key)
	if err != nil {
		return false, err
	}
	return db.DB.Has(kb)
}

func (db DB[KT, VT]) Fold(fn func(key KT, val VT, err error) error) (err error) {
	iter := db.DB.Items()
	for !IsTerminate(err) {
		// Fetch the key-value pair from DB as []byte.
		kb, vb, e := iter.Next()
		if IsTerminate(e) {
			return nil
		}
		err = errors.Join(err, e)

		// Decode the key into its type.
		key, e := db.keyCodec.Decode(kb)
		err = errors.Join(err, e)

		// Decode the value into its type.
		val, e := db.valCodec.Decode(vb)
		err = errors.Join(err, e)

		// Call the user-provided function.
		err = fn(key, val, err)
	}
	return
}

func IsTerminate(err error) bool {
	return errors.Is(err, pogreb.ErrIterationDone) || errors.Is(err, io.EOF)
}
