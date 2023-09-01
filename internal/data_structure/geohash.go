package data_structure

import (
	"errors"
	"fmt"
	"math"
)

// Limits from EPSG:900913 / EPSG:3785 / OSGEO:41001
const GeoLatMin float64 = -85.05112878
const GeoLatMax float64 = 85.05112878
const GeoLongMin float64 = -180
const GeoLongMax float64 = 180
const DR float64 = math.Pi / 180.0
const EarthRadiusInMeters float64 = 6372797.560856

// 52-bits gives us accuracy down to 0.6m
const GeoMaxStep uint8 = 26

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

var GeohashCoordRange = GeohashRange{
	MinLat:  GeoLatMin,
	MaxLat:  GeoLatMax,
	MinLong: GeoLongMin,
	MaxLong: GeoLongMax,
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
	exp2Step := float64(1 << GeoMaxStep)
	latOffset *= exp2Step
	longOffset *= exp2Step
	// lat is at even position, long is at odd position
	res.Bits = Interleave(uint32(latOffset), uint32(longOffset))
	return res, nil
}

func GeohashDecode(geohashRange GeohashRange, hash GeohashBits) (long float64, lat float64) {
	var step = hash.Step
	latBits, longBits := Deinterleave(hash.Bits)
	latBits = latBits << 1
	longBits = longBits << 1
	latScale := geohashRange.MaxLat - geohashRange.MinLat
	longScale := geohashRange.MaxLong - geohashRange.MinLong
	exp2Step := 1 << step
	latMin := geohashRange.MinLat + (float64(latBits)/float64(exp2Step))*latScale
	latMax := geohashRange.MinLat + (float64(latBits+1)/float64(exp2Step))*latScale
	longMin := geohashRange.MinLong + (float64(longBits)/float64(exp2Step))*longScale
	longMax := geohashRange.MinLong + (float64(longBits+1)/float64(exp2Step))*longScale

	// result is the center of the rectangle
	lat = (latMin + latMax) / 2
	long = (longMin + longMax) / 2
	if lat > GeoLatMax {
		lat = GeoLongMax
	}
	if lat < GeoLatMin {
		lat = GeoLatMin
	}
	if long > GeoLongMax {
		long = GeoLongMax
	}
	if long < GeoLongMin {
		long = GeoLongMin
	}
	return long, lat
}

func degToRad(angle float64) float64 {
	return angle * DR
}

func geohashGetLatDistance(lat1 float64, lat2 float64) float64 {
	return EarthRadiusInMeters * math.Abs(degToRad(lat2)-degToRad(lat1))
}

// Calculate distance using haversine great circle distance formula.
// Unit: meter
func GeohashGetDistance(lon1 float64, lat1 float64, lon2 float64, lat2 float64) float64 {
	lon1r := degToRad(lon1)
	lon2r := degToRad(lon2)
	v := math.Sin((lon2r - lon1r) / 2.0)
	if v == 0.0 {
		return geohashGetLatDistance(lat1, lat2)
	}
	lat1r := degToRad(lat1)
	lat2r := degToRad(lat2)
	u := math.Sin((lat2r - lat1r) / 2.0)
	a := u*u + math.Cos(lat1r)*math.Cos(lat2r)*v*v
	return 2.0 * EarthRadiusInMeters * math.Asin(math.Sqrt(a))
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

func squash(X uint64) uint32 {
	X &= 0x5555555555555555
	X = (X | (X >> 1)) & 0x3333333333333333
	X = (X | (X >> 2)) & 0x0f0f0f0f0f0f0f0f
	X = (X | (X >> 4)) & 0x00ff00ff00ff00ff
	X = (X | (X >> 8)) & 0x0000ffff0000ffff
	X = (X | (X >> 16)) & 0x00000000ffffffff
	return uint32(X)
}

// from https://graphics.stanford.edu/~seander/bithacks.html#InterleaveBMN
func Interleave(x uint32, y uint32) uint64 {
	return spread(x) | (spread(y) << 1)
}

// return even and odd bitlevels of X
func Deinterleave(x uint64) (uint32, uint32) {
	return squash(x), squash(x >> 1)
}

/*
Compute sorted set score [min, max) we should query to get all the elements inside
the specific are 'hash'.
*/
func GeohashGetScoreLimit(hash GeohashBits) (min GeoHashFix52Bits, max GeoHashFix52Bits) {
	min = GeohashAlign52Bits(hash)
	hash.Bits++
	max = GeohashAlign52Bits(hash)
	return
}
