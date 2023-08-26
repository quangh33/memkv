package core

import (
	"errors"
	"fmt"
)

const GEO_LAT_MIN float64 = -90
const GEO_LAT_MAX float64 = 90
const GEO_LONG_MIN float64 = -180
const GEO_LONG_MAX float64 = 180
const GEO_ALPHABET string = "0123456789bcdefghjkmnpqrstuvwxyz"
const GEO_MAX_STEP uint8 = 26 // 52-bits gives us accuracy down to 0.6m

type GeohashBits struct {
	Step uint8
	Bits uint64
}

type GeohashRange struct {
	MinLat  float64
	MaxLat  float64
	MinLong float64
	MaxLong float64
}

func Base32Encode(x uint64) string {
	b := [11]byte{}
	for i := 0; i < 11; i++ {
		shift := 52 - (i+1)*5
		if shift <= 0 {
			b[i] = GEO_ALPHABET[0]
			break
		}
		idx := (x >> shift) & 0b11111
		b[i] = GEO_ALPHABET[idx]
	}
	return string(b[:])
}

func GeohashEncode(geohashRange GeohashRange, long float64, lat float64, step uint8) (*GeohashBits, error) {
	if long > geohashRange.MaxLong || long < geohashRange.MinLong ||
		lat > geohashRange.MaxLat || lat < geohashRange.MinLat {
		return nil, errors.New(fmt.Sprintf("invalid coord: %f, %f", long, lat))
	}

	res := &GeohashBits{
		Step: step,
		Bits: 0,
	}

	latOffset := (lat - geohashRange.MinLat) / (geohashRange.MaxLat - geohashRange.MinLat)
	longOffset := (long - geohashRange.MinLong) / (geohashRange.MaxLong - geohashRange.MinLong)
	exp2Step := float64(1 << GEO_MAX_STEP)
	latOffset *= exp2Step
	longOffset *= exp2Step
	res.Bits = Interleave(uint32(latOffset), uint32(longOffset))
	return res, nil
}

func spread(x uint32) uint64 {
	X := uint64(x)
	X = (X | (X << 16)) & 0x0000ffff0000ffff
	X = (X | (X << 8)) & 0x00ff00ff00ff00ff
	X = (X | (X << 4)) & 0x0f0f0f0f0f0f0f0f
	X = (X | (X << 2)) & 0x3333333333333333
	X = (X | (X << 1)) & 0x5555555555555555
	return X
}

// from https://graphics.stanford.edu/~seander/bithacks.html#InterleaveBMN
func Interleave(x uint32, y uint32) uint64 {
	return spread(x) | (spread(y) << 1)
}
