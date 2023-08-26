package core_test

import (
	"memkv/core"
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
		value, _ := core.GeohashEncode(k[0], k[1], core.GEO_MAX_STEP)
		output := core.Base32Encode(value.Bits)
		// fmt.Println(output)
		if output != v {
			t.Fail()
		}
	}
}

func TestInterleave(t *testing.T) {
	cases := map[[2]uint32]uint64{
		[2]uint32{0b1111, 0b1010}: 0b11011101,
		[2]uint32{0b1, 0b0}:       0b1,
		[2]uint32{0b101, 0b111}:   0b111011,
	}

	for k, v := range cases {
		value := core.Interleave(k[0], k[1])
		if v != value {
			t.Fail()
		}
	}
}
