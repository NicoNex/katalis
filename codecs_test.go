package katalis_test

import (
	"math"
	"testing"

	"github.com/NicoNex/katalis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringCodec(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"simple", "hello"},
		{"unicode", "こんにちは"},
		{"special chars", "!@#$%^&*()"},
		{"long string", string(make([]byte, 10000))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.StringCodec.Encode(tt.input)
			require.NoError(t, err)

			decoded, err := katalis.StringCodec.Decode(encoded)
			require.NoError(t, err)

			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestBytesCodec(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"nil", nil},
		{"empty", []byte{}},
		{"simple", []byte{1, 2, 3, 4, 5}},
		{"zeros", []byte{0, 0, 0}},
		{"large", make([]byte, 10000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.BytesCodec.Encode(tt.input)
			require.NoError(t, err)

			decoded, err := katalis.BytesCodec.Decode(encoded)
			require.NoError(t, err)

			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestUint64Codec(t *testing.T) {
	tests := []struct {
		name  string
		input uint64
	}{
		{"zero", 0},
		{"small", 42},
		{"max", math.MaxUint64},
		{"powers of 2", 1 << 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.Uint64Codec.Encode(tt.input)
			require.NoError(t, err)
			assert.Len(t, encoded, 8)

			decoded, err := katalis.Uint64Codec.Decode(encoded)
			require.NoError(t, err)

			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestUint32Codec(t *testing.T) {
	tests := []struct {
		name  string
		input uint32
	}{
		{"zero", 0},
		{"small", 42},
		{"max", math.MaxUint32},
		{"powers of 2", 1 << 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.Uint32Codec.Encode(tt.input)
			require.NoError(t, err)
			assert.Len(t, encoded, 4)

			decoded, err := katalis.Uint32Codec.Decode(encoded)
			require.NoError(t, err)

			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestUint16Codec(t *testing.T) {
	tests := []struct {
		name  string
		input uint16
	}{
		{"zero", 0},
		{"small", 42},
		{"max", math.MaxUint16},
		{"powers of 2", 1 << 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.Uint16Codec.Encode(tt.input)
			require.NoError(t, err)
			assert.Len(t, encoded, 2)

			decoded, err := katalis.Uint16Codec.Decode(encoded)
			require.NoError(t, err)

			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestUintCodec(t *testing.T) {
	tests := []struct {
		name  string
		input uint
	}{
		{"zero", 0},
		{"small", 42},
		{"large", 1 << 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.UintCodec.Encode(tt.input)
			require.NoError(t, err)
			assert.Len(t, encoded, 4) // Uses uint32 under the hood

			decoded, err := katalis.UintCodec.Decode(encoded)
			require.NoError(t, err)

			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestInt64Codec(t *testing.T) {
	tests := []struct {
		name  string
		input int64
	}{
		{"zero", 0},
		{"positive", 42},
		{"negative", -42},
		{"max", math.MaxInt64},
		{"min", math.MinInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.Int64Codec.Encode(tt.input)
			require.NoError(t, err)
			assert.Len(t, encoded, 8)

			decoded, err := katalis.Int64Codec.Decode(encoded)
			require.NoError(t, err)

			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestInt32Codec(t *testing.T) {
	tests := []struct {
		name  string
		input int32
	}{
		{"zero", 0},
		{"positive", 42},
		{"negative", -42},
		{"max", math.MaxInt32},
		{"min", math.MinInt32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.Int32Codec.Encode(tt.input)
			require.NoError(t, err)
			assert.Len(t, encoded, 4)

			decoded, err := katalis.Int32Codec.Decode(encoded)
			require.NoError(t, err)

			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestInt16Codec(t *testing.T) {
	tests := []struct {
		name  string
		input int16
	}{
		{"zero", 0},
		{"positive", 42},
		{"negative", -42},
		{"max", math.MaxInt16},
		{"min", math.MinInt16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.Int16Codec.Encode(tt.input)
			require.NoError(t, err)
			assert.Len(t, encoded, 2)

			decoded, err := katalis.Int16Codec.Decode(encoded)
			require.NoError(t, err)

			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestIntCodec(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{"zero", 0},
		{"positive", 42},
		{"negative", -42},
		{"large positive", 1 << 30},
		{"large negative", -(1 << 30)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.IntCodec.Encode(tt.input)
			require.NoError(t, err)
			assert.Len(t, encoded, 4) // Uses uint32 under the hood

			decoded, err := katalis.IntCodec.Decode(encoded)
			require.NoError(t, err)

			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestFloat64Codec(t *testing.T) {
	tests := []struct {
		name  string
		input float64
	}{
		{"zero", 0.0},
		{"positive", 3.14159},
		{"negative", -2.71828},
		{"max", math.MaxFloat64},
		{"smallest positive", math.SmallestNonzeroFloat64},
		{"infinity", math.Inf(1)},
		{"negative infinity", math.Inf(-1)},
		{"NaN", math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.Float64Codec.Encode(tt.input)
			require.NoError(t, err)
			assert.Len(t, encoded, 8)

			decoded, err := katalis.Float64Codec.Decode(encoded)
			require.NoError(t, err)

			if math.IsNaN(tt.input) {
				assert.True(t, math.IsNaN(decoded))
			} else {
				assert.Equal(t, tt.input, decoded)
			}
		})
	}
}

func TestFloat32Codec(t *testing.T) {
	tests := []struct {
		name  string
		input float32
	}{
		{"zero", 0.0},
		{"positive", 3.14},
		{"negative", -2.71},
		{"max", math.MaxFloat32},
		{"smallest positive", math.SmallestNonzeroFloat32},
		{"infinity", float32(math.Inf(1))},
		{"negative infinity", float32(math.Inf(-1))},
		{"NaN", float32(math.NaN())},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := katalis.Float32Codec.Encode(tt.input)
			require.NoError(t, err)
			assert.Len(t, encoded, 4)

			decoded, err := katalis.Float32Codec.Decode(encoded)
			require.NoError(t, err)

			if math.IsNaN(float64(tt.input)) {
				assert.True(t, math.IsNaN(float64(decoded)))
			} else {
				assert.Equal(t, tt.input, decoded)
			}
		})
	}
}

func TestGobCodecStruct(t *testing.T) {
	type Address struct {
		Street string
		City   string
		ZIP    int
	}

	type Person struct {
		Name    string
		Age     int
		Address Address
		Tags    []string
	}

	person := Person{
		Name: "Bob",
		Age:  25,
		Address: Address{
			Street: "123 Main St",
			City:   "Springfield",
			ZIP:    12345,
		},
		Tags: []string{"developer", "golang"},
	}

	codec := katalis.Gob[Person]()

	encoded, err := codec.Encode(person)
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := codec.Decode(encoded)
	require.NoError(t, err)
	assert.Equal(t, person, decoded)
}

func TestGobCodecMap(t *testing.T) {
	input := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	codec := katalis.Gob[map[string]int]()

	encoded, err := codec.Encode(input)
	require.NoError(t, err)

	decoded, err := codec.Decode(encoded)
	require.NoError(t, err)

	assert.Equal(t, input, decoded)
}

func TestGobCodecSlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}

	codec := katalis.Gob[[]int]()

	encoded, err := codec.Encode(input)
	require.NoError(t, err)

	decoded, err := codec.Decode(encoded)
	require.NoError(t, err)

	assert.Equal(t, input, decoded)
}

func TestGobCodecWithTypeInference(t *testing.T) {
	type Data struct {
		ID    int
		Value string
	}

	data := Data{ID: 1, Value: "test"}
	codec := katalis.Gob(data)

	encoded, err := codec.Encode(data)
	require.NoError(t, err)

	decoded, err := codec.Decode(encoded)
	require.NoError(t, err)

	assert.Equal(t, data, decoded)
}

func TestGobCodecEmpty(t *testing.T) {
	type Empty struct{}

	codec := katalis.Gob[Empty]()

	encoded, err := codec.Encode(Empty{})
	require.NoError(t, err)

	decoded, err := codec.Decode(encoded)
	require.NoError(t, err)

	assert.Equal(t, Empty{}, decoded)
}

func TestCodecForPrimitives(t *testing.T) {
	// This tests the codecFor function indirectly by testing it
	// works correctly for different types

	// Note: codecFor is unexported, but we can test it through
	// a hypothetical OpenAuto function if it exists

	// For now, we'll just verify all codecs work as expected
	t.Run("all primitives", func(t *testing.T) {
		// String
		{
			b, err := katalis.StringCodec.Encode("test")
			require.NoError(t, err)
			s, err := katalis.StringCodec.Decode(b)
			require.NoError(t, err)
			assert.Equal(t, "test", s)
		}

		// Bytes
		{
			b, err := katalis.BytesCodec.Encode([]byte{1, 2, 3})
			require.NoError(t, err)
			bytes, err := katalis.BytesCodec.Decode(b)
			require.NoError(t, err)
			assert.Equal(t, []byte{1, 2, 3}, bytes)
		}

		// All integer types work
		{
			b, err := katalis.IntCodec.Encode(42)
			require.NoError(t, err)
			i, err := katalis.IntCodec.Decode(b)
			require.NoError(t, err)
			assert.Equal(t, 42, i)
		}
	})
}
