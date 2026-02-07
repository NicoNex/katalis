package katalis

import (
	"errors"
	"io"
	"iter"
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
	return OpenOptions(
		path,
		keyCodec,
		valCodec,
		&pogreb.Options{
			BackgroundSyncInterval:       -1,
			BackgroundCompactionInterval: time.Hour * 24,
		},
	)
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

// Del deletes the value for the given key from the DB.
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

// Fold iterates over all keys in the database calling the function `fn` for
// each key. If the function returns an error, no further keys are processed
// and the error returned.
func (db DB[KT, VT]) Fold(fn func(key KT, val VT, err error) error) (err error) {
	iter := db.DB.Items()
	for err == nil {
		// Fetch the key-value pair from DB as []byte.
		kb, vb, e := iter.Next()
		if isTerminate(e) {
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

// Items returns an iterator over all key-value pairs in the database. Decode
// errors are silently skipped, allowing iteration to continue. Use AllItems if
// you need to handle errors explicitly.
func (db DB[KT, VT]) Items() iter.Seq2[KT, VT] {
	items := db.DB.Items()

	return func(yield func(KT, VT) bool) {
		for {
			var key KT
			var val VT

			kb, vb, err := items.Next()
			if isTerminate(err) {
				return
			}

			if err == nil {
				key, err = db.keyCodec.Decode(kb)
			}
			if err == nil {
				val, err = db.valCodec.Decode(vb)
			}

			// Skip entries with errors
			if err != nil {
				continue
			}

			if !yield(key, val) {
				return
			}
		}
	}
}

// Entry represents a key-value pair from the database. It is used by AllItems
// to return both the key and value together with potential errors during iteration.
type Entry[KT, VT any] struct {
	Key   KT
	Value VT
}

// AllItems returns an iterator over all key-value pairs in the database with
// error reporting. Unlike Items, decode errors are yielded to the caller rather
// than terminating iteration.
func (db DB[KT, VT]) AllItems() iter.Seq2[Entry[KT, VT], error] {
	return func(yield func(Entry[KT, VT], error) bool) {
		iter := db.DB.Items()
		for {
			var entry Entry[KT, VT]

			kb, vb, err := iter.Next()
			if isTerminate(err) {
				return
			}

			if err == nil {
				entry.Key, err = db.keyCodec.Decode(kb)
			}
			if err == nil {
				entry.Value, err = db.valCodec.Decode(vb)
			}

			if !yield(entry, err) {
				return
			}
		}
	}
}

func isTerminate(err error) bool {
	return errors.Is(err, pogreb.ErrIterationDone) || errors.Is(err, io.EOF)
}
