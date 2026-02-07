package katalis_test

import (
	"path/filepath"
	"testing"

	"github.com/NicoNex/katalis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.StringCodec)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	assert.DirExists(t, dbPath)
}

func TestPutAndGet(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.StringCodec)
	require.NoError(t, err)
	defer db.Close()

	// Put a value
	err = db.Put("key1", "value1")
	require.NoError(t, err)

	// Get the value
	val, err := db.Get("key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)
}

func TestGetNonExistent(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.StringCodec)
	require.NoError(t, err)
	defer db.Close()

	// Get non-existent key
	val, err := db.Get("nonexistent")
	require.NoError(t, err)
	assert.Equal(t, "", val) // Should be zero value
}

func TestDel(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.IntCodec)
	require.NoError(t, err)
	defer db.Close()

	// Put and verify
	err = db.Put("key1", 42)
	require.NoError(t, err)

	has, err := db.Has("key1")
	require.NoError(t, err)
	assert.True(t, has)

	// Delete
	err = db.Del("key1")
	require.NoError(t, err)

	// Verify deletion
	has, err = db.Has("key1")
	require.NoError(t, err)
	assert.False(t, has)
}

func TestHas(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.StringCodec)
	require.NoError(t, err)
	defer db.Close()

	// Check non-existent key
	has, err := db.Has("key1")
	require.NoError(t, err)
	assert.False(t, has)

	// Put a value
	err = db.Put("key1", "value1")
	require.NoError(t, err)

	// Check existing key
	has, err = db.Has("key1")
	require.NoError(t, err)
	assert.True(t, has)
}

func TestMultipleKeysAndValues(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.IntCodec)
	require.NoError(t, err)
	defer db.Close()

	// Put multiple values
	testData := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}

	for k, v := range testData {
		err := db.Put(k, v)
		require.NoError(t, err)
	}

	// Verify all values
	for k, expected := range testData {
		val, err := db.Get(k)
		require.NoError(t, err)
		assert.Equal(t, expected, val)
	}
}

func TestUpdateValue(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.StringCodec)
	require.NoError(t, err)
	defer db.Close()

	// Put initial value
	err = db.Put("key1", "value1")
	require.NoError(t, err)

	// Update value
	err = db.Put("key1", "value2")
	require.NoError(t, err)

	// Verify updated value
	val, err := db.Get("key1")
	require.NoError(t, err)
	assert.Equal(t, "value2", val)
}

func TestFold(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.IntCodec)
	require.NoError(t, err)
	defer db.Close()

	// Put test data
	testData := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	for k, v := range testData {
		err := db.Put(k, v)
		require.NoError(t, err)
	}

	// Fold over all items
	collected := make(map[string]int)
	err = db.Fold(func(key string, val int, iterErr error) error {
		require.NoError(t, iterErr)
		collected[key] = val
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, testData, collected)
}

func TestItems(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.IntCodec)
	require.NoError(t, err)
	defer db.Close()

	// Put test data
	testData := map[string]int{
		"x": 10,
		"y": 20,
		"z": 30,
	}

	for k, v := range testData {
		err := db.Put(k, v)
		require.NoError(t, err)
	}

	// Iterate using Items
	collected := make(map[string]int)
	for key, val := range db.Items() {
		collected[key] = val
	}

	assert.Equal(t, testData, collected)
}

func TestAllItems(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.IntCodec)
	require.NoError(t, err)
	defer db.Close()

	// Put test data
	testData := map[string]int{
		"p": 100,
		"q": 200,
		"r": 300,
	}

	for k, v := range testData {
		err := db.Put(k, v)
		require.NoError(t, err)
	}

	// Iterate using AllItems
	collected := make(map[string]int)
	for entry, iterErr := range db.AllItems() {
		require.NoError(t, iterErr)
		collected[entry.Key] = entry.Value
	}

	assert.Equal(t, testData, collected)
}

func TestAllItemsWithEmptyDB(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.StringCodec)
	require.NoError(t, err)
	defer db.Close()

	// Iterate over empty database
	count := 0
	for _, iterErr := range db.AllItems() {
		require.NoError(t, iterErr)
		count++
	}

	assert.Equal(t, 0, count)
}

func TestItemsEarlyExit(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.IntCodec, katalis.StringCodec)
	require.NoError(t, err)
	defer db.Close()

	// Put multiple values
	for i := 0; i < 10; i++ {
		err := db.Put(i, "value")
		require.NoError(t, err)
	}

	// Exit early from iteration
	count := 0
	for range db.Items() {
		count++
		if count >= 3 {
			break
		}
	}

	assert.Equal(t, 3, count)
}

func TestDifferentCodecs(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name     string
		key      any
		value    any
		keyCodec any
		valCodec any
	}{
		{"uint64", uint64(123), uint64(456), katalis.Uint64Codec, katalis.Uint64Codec},
		{"int64", int64(-123), int64(-456), katalis.Int64Codec, katalis.Int64Codec},
		{"float64", float64(3.14), float64(2.71), katalis.Float64Codec, katalis.Float64Codec},
		{"bytes", []byte("key"), []byte("value"), katalis.BytesCodec, katalis.BytesCodec},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbPath := filepath.Join(dir, tt.name+".db")

			switch kc := tt.keyCodec.(type) {
			case katalis.Codec[uint64]:
				db, err := katalis.Open(dbPath, kc, tt.valCodec.(katalis.Codec[uint64]))
				require.NoError(t, err)
				defer db.Close()

				err = db.Put(tt.key.(uint64), tt.value.(uint64))
				require.NoError(t, err)

				val, err := db.Get(tt.key.(uint64))
				require.NoError(t, err)
				assert.Equal(t, tt.value, val)

			case katalis.Codec[int64]:
				db, err := katalis.Open(dbPath, kc, tt.valCodec.(katalis.Codec[int64]))
				require.NoError(t, err)
				defer db.Close()

				err = db.Put(tt.key.(int64), tt.value.(int64))
				require.NoError(t, err)

				val, err := db.Get(tt.key.(int64))
				require.NoError(t, err)
				assert.Equal(t, tt.value, val)

			case katalis.Codec[float64]:
				db, err := katalis.Open(dbPath, kc, tt.valCodec.(katalis.Codec[float64]))
				require.NoError(t, err)
				defer db.Close()

				err = db.Put(tt.key.(float64), tt.value.(float64))
				require.NoError(t, err)

				val, err := db.Get(tt.key.(float64))
				require.NoError(t, err)
				assert.Equal(t, tt.value, val)

			case katalis.Codec[[]byte]:
				db, err := katalis.Open(dbPath, kc, tt.valCodec.(katalis.Codec[[]byte]))
				require.NoError(t, err)
				defer db.Close()

				err = db.Put(tt.key.([]byte), tt.value.([]byte))
				require.NoError(t, err)

				val, err := db.Get(tt.key.([]byte))
				require.NoError(t, err)
				assert.Equal(t, tt.value, val)
			}
		})
	}
}

func TestGobCodec(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	type Person struct {
		Name string
		Age  int
	}

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.Gob[Person]())
	require.NoError(t, err)
	defer db.Close()

	person := Person{Name: "Alice", Age: 30}
	err = db.Put("person1", person)
	require.NoError(t, err)

	retrieved, err := db.Get("person1")
	require.NoError(t, err)
	assert.Equal(t, person, retrieved)
}

func TestReopenDatabase(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	// Create and write to database
	{
		db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.StringCodec)
		require.NoError(t, err)

		err = db.Put("persist", "value")
		require.NoError(t, err)

		err = db.Close()
		require.NoError(t, err)
	}

	// Reopen and verify data persists
	{
		db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.StringCodec)
		require.NoError(t, err)
		defer db.Close()

		val, err := db.Get("persist")
		require.NoError(t, err)
		assert.Equal(t, "value", val)
	}
}

func TestConcurrentReads(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := katalis.Open(dbPath, katalis.StringCodec, katalis.IntCodec)
	require.NoError(t, err)
	defer db.Close()

	// Populate database
	for i := 0; i < 100; i++ {
		err := db.Put("key", i)
		require.NoError(t, err)
	}

	// Concurrent reads should work
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := db.Get("key")
			assert.NoError(t, err)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
