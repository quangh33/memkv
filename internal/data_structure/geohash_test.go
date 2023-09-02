package data_structure_test

import (
	"github.com/stretchr/testify/assert"
	"math"
	"memkv/internal/core"
	"memkv/internal/data_structure"
	"testing"
)

func TestGeohashEncode(t *testing.T) {
	cases := map[[2]float64]string{
		// {long, lat} => geohash
		[2]float64{13.361389, 38.115556}:  "sqc8b49rny0",
		[2]float64{15.087269, 37.502669}:  "sqdtr74hyu0",
		[2]float64{100, 80}:               "ynpp5e9cbb0",
		[2]float64{40, -20}:               "kukqnpp5e90",
		[2]float64{-20, 39}:               "ewcvbgsrqn0",
		[2]float64{0, -50}:                "hp0581b0bh0",
		[2]float64{12.345678, -20.654321}: "kk2f2zvzg50",
		[2]float64{0, 0}:                  "s0000000000",
		[2]float64{180, 85}:               "bp0581b0bh0",
		[2]float64{-180, -85}:             "00bh2n0p050",
	}

	for k, v := range cases {
		value, _ := data_structure.GeohashEncode(data_structure.GeohashStandardRange, k[0], k[1], data_structure.GeoMaxStep)
		output := core.Base32encoding.Encode(value.Bits)
		assert.EqualValues(t, v, output)
	}
}

func TestGeohashDecode(t *testing.T) {
	cases := map[string][2]float64{
		"sqc8b49rny0": {13.361389, 38.115556},
		"sqdtr74hyu0": {15.087269, 37.502669},
		"ynpp5e9cbb0": {100, 80},
		"kukqnpp5e90": {40, -20},
		"ewcvbgsrqn0": {-20, 39},
		"hp0581b0bh0": {0, -50},
		"kk2f2zvzg50": {12.345678, -20.654321},
		"s0000000000": {0, 0},
		"bp0581b0bh0": {180, 85},
		"00bh2n0p050": {-180, -85},
	}

	normalGeoRange := data_structure.GeohashRange{
		MinLat:  -90,
		MaxLat:  90,
		MinLong: -180,
		MaxLong: 180,
	}
	for hash, expected := range cases {
		geohashBits := data_structure.GeohashBits{
			Step: data_structure.GeoMaxStep,
			Bits: core.Base32encoding.Decode(hash) << 2,
			// need to shift-left 2 because base32 decode returns a 50bits value
		}
		long, lat := data_structure.GeohashDecodeAreaToLongLat(normalGeoRange, geohashBits)
		assert.LessOrEqual(t, data_structure.GeohashGetDistance(long, lat, expected[0], expected[1]), 1.0)
	}
}

func TestInterleave(t *testing.T) {
	cases := map[[2]uint32]uint64{
		[2]uint32{0b1111, 0b1010}: 0b11011101,
		[2]uint32{0b1, 0b0}:       0b1,
		[2]uint32{0b101, 0b111}:   0b111011,
	}

	for k, v := range cases {
		value := data_structure.Interleave(k[0], k[1])
		assert.EqualValues(t, v, value)
	}
}

func TestDeinterleave(t *testing.T) {
	cases := map[uint64][2]uint32{
		0b11011101: {0b1111, 0b1010},
		0b1:        {0b1, 0b0},
		0b111011:   {0b101, 0b111},
	}

	for k, v := range cases {
		even, odd := data_structure.Deinterleave(k)
		assert.EqualValues(t, v[0], even)
		assert.EqualValues(t, v[1], odd)
	}
}

func TestBase32Decode(t *testing.T) {
	cases := []uint64{
		0b1001011010100101011010100101011010100101011010100101,
		0b110111100010111101101010011111100010111101101010011,
		0b1000010101000000010101000000010101000000010101000000,
		0b10101000000010101000000010101000000010101,
		0b1100010110010110100001010001000100110111101001111011,
		0b1101111010101101011010110101011010100101011010100101,
	}

	for _, x := range cases {
		s := core.Base32encoding.Encode(x)
		decode := core.Base32encoding.Decode(s)
		assert.EqualValues(t, x>>2, decode)
	}
}

func TestGeohashGetDistance(t *testing.T) {
	cases := map[[4]float64]float64{
		[4]float64{20, 10, 40, 30}:        3041460.716138,
		[4]float64{10, 85, 10, -85}:       18908471,
		[4]float64{180, 20, 179.9999, 20}: 10.451853,
	}

	for points, dis := range cases {
		output := data_structure.GeohashGetDistance(points[0], points[1], points[2], points[3])
		assert.LessOrEqual(t, math.Abs(output-dis), 1e-5)
	}
}

func TestGeohashMoveX(t *testing.T) {
	hash := data_structure.GeohashBits{
		Step: 2,
		Bits: 0b1001,
	}
	data_structure.GeohashMoveX(&hash, 1)
	assert.EqualValues(t, 0b1011, hash.Bits)
	hash = data_structure.GeohashBits{
		Step: 2,
		Bits: 0b1001,
	}
	data_structure.GeohashMoveX(&hash, -1)
	assert.EqualValues(t, 0b11, hash.Bits)

	hash = data_structure.GeohashBits{
		Step: 3,
		Bits: 0b101010,
	}
	data_structure.GeohashMoveX(&hash, -1)
	assert.EqualValues(t, 0b101000, hash.Bits)

	hash = data_structure.GeohashBits{
		Step: 3,
		Bits: 0b101010,
	}
	data_structure.GeohashMoveX(&hash, 1)
	assert.EqualValues(t, 0b0, hash.Bits)

	hash = data_structure.GeohashBits{
		Step: 3,
		Bits: 0b101111,
	}
	data_structure.GeohashMoveX(&hash, 1)
	assert.EqualValues(t, 0b101, hash.Bits)

	hash = data_structure.GeohashBits{
		Step: 2,
		Bits: 0b1011,
	}
	data_structure.GeohashMoveX(&hash, -1)
	assert.EqualValues(t, 0b1001, hash.Bits)

	hash = data_structure.GeohashBits{
		Step: 2,
		Bits: 0b1011,
	}
	data_structure.GeohashMoveX(&hash, 1)
	assert.EqualValues(t, 0b1, hash.Bits)
}

func TestGeohashMoveY(t *testing.T) {
	hash := data_structure.GeohashBits{
		Step: 2,
		Bits: 0b1001,
	}
	data_structure.GeohashMoveY(&hash, 1)
	assert.EqualValues(t, 0b1100, hash.Bits)
	hash = data_structure.GeohashBits{
		Step: 2,
		Bits: 0b1001,
	}
	data_structure.GeohashMoveY(&hash, -1)
	assert.EqualValues(t, 0b1000, hash.Bits)

	hash = data_structure.GeohashBits{
		Step: 3,
		Bits: 0b101010,
	}
	data_structure.GeohashMoveY(&hash, -1)
	assert.EqualValues(t, 0b111111, hash.Bits)

	hash = data_structure.GeohashBits{
		Step: 3,
		Bits: 0b101010,
	}
	data_structure.GeohashMoveY(&hash, 1)
	assert.EqualValues(t, 0b101011, hash.Bits)

	hash = data_structure.GeohashBits{
		Step: 3,
		Bits: 0b101111,
	}
	data_structure.GeohashMoveY(&hash, 1)
	assert.EqualValues(t, 0b111010, hash.Bits)

	hash = data_structure.GeohashBits{
		Step: 2,
		Bits: 0b1011,
	}
	data_structure.GeohashMoveY(&hash, -1)
	assert.EqualValues(t, 0b1010, hash.Bits)

	hash = data_structure.GeohashBits{
		Step: 2,
		Bits: 0b1011,
	}
	data_structure.GeohashMoveY(&hash, 1)
	assert.EqualValues(t, 0b1110, hash.Bits)
}

func TestGeohashBits_GetNeighbors(t *testing.T) {
	hash := data_structure.GeohashBits{
		Step: 2,
		Bits: 0b1001,
	}
	ret := hash.GetNeighbors()
	assert.EqualValues(t, 0b1100, ret.North.Bits)
	assert.EqualValues(t, 0b1011, ret.East.Bits)
	assert.EqualValues(t, 0b1000, ret.South.Bits)
	assert.EqualValues(t, 0b0011, ret.West.Bits)
	assert.EqualValues(t, 0b1110, ret.NorthEast.Bits)
	assert.EqualValues(t, 0b1010, ret.SouthEast.Bits)
	assert.EqualValues(t, 0b0010, ret.SouthWest.Bits)
	assert.EqualValues(t, 0b0110, ret.NorthWest.Bits)
}
