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
const MercatorMax float64 = 20037726.37

// 52-bits gives us accuracy down to 0.6m
const GeoMaxStep uint8 = 26

type GeohashBits struct {
	Step uint8
	Bits uint64
}

type GeohashCircularSearchQuery struct {
	long        float64
	lat         float64
	radiusMeter float64
}

type GeohashRange struct {
	MinLat  float64
	MaxLat  float64
	MinLong float64
	MaxLong float64
}

type GeohashNeighbors struct {
	North     GeohashBits
	East      GeohashBits
	West      GeohashBits
	South     GeohashBits
	NorthEast GeohashBits
	SouthEast GeohashBits
	NorthWest GeohashBits
	SouthWest GeohashBits
}

type GeohashArea struct {
	hash   GeohashBits
	grange GeohashRange
}

/*
________________
|    |    |    |
|    |    |    |
----------------
|    |    |    |
|    |    |    |
----------------
|    |    |    |
|    |    |    |
----------------
*/
type GeohashRadius struct {
	hash      GeohashBits
	area      GeohashArea
	neighbors GeohashNeighbors
}

type GeoPoint struct {
	long   float64
	lat    float64
	dist   float64 // distance to searching point
	member string
	score  float64
}

var GeohashCoordRange = GeohashRange{
	MinLat:  GeoLatMin,
	MaxLat:  GeoLatMax,
	MinLong: GeoLongMin,
	MaxLong: GeoLongMax,
}

var GeohashStandardRange = GeohashRange{
	MinLat:  -90,
	MaxLat:  90,
	MinLong: -180,
	MaxLong: 180,
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
	res.Bits = uint64(GeohashAlign52Bits(*res))
	return res, nil
}

func GeohashDecode(geohashRange GeohashRange, hash GeohashBits) GeohashArea {
	var step = hash.Step
	latBits, longBits := Deinterleave(hash.Bits)
	latScale := geohashRange.MaxLat - geohashRange.MinLat
	longScale := geohashRange.MaxLong - geohashRange.MinLong
	exp2Step := 1 << step
	res := GeohashArea{
		hash: hash,
		grange: GeohashRange{
			MinLat:  geohashRange.MinLat + (float64(latBits)/float64(exp2Step))*latScale,
			MaxLat:  geohashRange.MinLat + (float64(latBits+1)/float64(exp2Step))*latScale,
			MinLong: geohashRange.MinLong + (float64(longBits)/float64(exp2Step))*longScale,
			MaxLong: geohashRange.MinLong + (float64(longBits+1)/float64(exp2Step))*longScale,
		},
	}
	return res
}

func GeohashDecodeAreaToLongLat(geohashRange GeohashRange, hash GeohashBits) (long float64, lat float64) {
	area := GeohashDecode(geohashRange, hash)
	// result is the center of the rectangle
	lat = (area.grange.MinLat + area.grange.MaxLat) / 2
	long = (area.grange.MinLong + area.grange.MaxLong) / 2
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

/* Move left/right 1 step */
func GeohashMoveX(hash *GeohashBits, d int8) {
	if d == 0 {
		return
	}

	x := hash.Bits & 0xaaaaaaaaaaaaaaaa
	y := hash.Bits & 0x5555555555555555
	var zz uint64 = 0x5555555555555555 >> (64 - hash.Step*2)

	if d > 0 {
		x = x + (zz + 1)
	} else {
		x = x | zz
		x = x - (zz + 1)
	}
	x &= 0xaaaaaaaaaaaaaaaa >> (64 - hash.Step*2)
	hash.Bits = x | y
}

/* Move up/down 1 step */
func GeohashMoveY(hash *GeohashBits, d int8) {
	if d == 0 {
		return
	}

	x := hash.Bits & 0xaaaaaaaaaaaaaaaa
	y := hash.Bits & 0x5555555555555555
	var zz uint64 = 0xaaaaaaaaaaaaaaaa >> (64 - hash.Step*2)

	if d > 0 {
		y = y + (zz + 1)
	} else {
		y = y | zz
		y = y - (zz + 1)
	}
	y &= 0x5555555555555555 >> (64 - hash.Step*2)
	hash.Bits = x | y
}

func (hash GeohashBits) GetNeighbors() GeohashNeighbors {
	ret := GeohashNeighbors{
		North:     hash,
		East:      hash,
		West:      hash,
		South:     hash,
		NorthEast: hash,
		SouthEast: hash,
		NorthWest: hash,
		SouthWest: hash,
	}

	GeohashMoveX(&ret.East, 1)
	GeohashMoveX(&ret.West, -1)
	GeohashMoveY(&ret.North, 1)
	GeohashMoveY(&ret.South, -1)

	GeohashMoveX(&ret.NorthWest, -1)
	GeohashMoveY(&ret.NorthWest, 1)

	GeohashMoveX(&ret.NorthEast, 1)
	GeohashMoveY(&ret.NorthEast, 1)

	GeohashMoveX(&ret.SouthEast, 1)
	GeohashMoveY(&ret.SouthEast, -1)

	GeohashMoveX(&ret.SouthWest, -1)
	GeohashMoveY(&ret.SouthWest, -1)

	return ret
}

/*
Calculate a set of areas (center + 8 neighbors) that are able to cover a range query
*/
func GeohashCalculateSearchingAreas(q GeohashCircularSearchQuery) (*GeohashRadius, error) {
	steps := GeohashEstimateStepsByRadius(q.radiusMeter)
	centerHash, err := GeohashEncode(GeohashCoordRange, q.long, q.lat, steps)
	if err != nil {
		return nil, err
	}
	neighbors := centerHash.GetNeighbors()
	areas := GeohashDecode(GeohashCoordRange, *centerHash)
	ret := GeohashRadius{
		hash:      *centerHash,
		area:      areas,
		neighbors: neighbors,
	}
	return &ret, nil
}

/*
Search all points inside area covered by 'hash' that is within searching distance to (lon, lat) point
*/
func GeohashGetMemberInsideBox(zset ZSet, q GeohashCircularSearchQuery, hash GeohashBits) []GeoPoint {
	mi, ma := GeohashGetScoreLimit(hash)
	// [min, max)
	zrange := ZRange{
		min:   float64(mi),
		max:   float64(ma),
		minex: false,
		maxex: true,
	}
	x := zset.zskiplist.FindFirstInRange(zrange)
	if x == nil {
		return []GeoPoint{}
	}
	var ret []GeoPoint
	for x != nil {
		if !zrange.ValueLteMax(x.score) {
			break
		}
		score := x.score
		long, lat := GeohashDecodeAreaToLongLat(GeohashCoordRange, GeohashBits{
			Step: GeoMaxStep,
			Bits: uint64(score),
		})
		dist := GeohashGetDistance(long, lat, q.long, q.lat)
		if dist <= q.radiusMeter {
			ret = append(ret, GeoPoint{
				long:   long,
				lat:    lat,
				dist:   dist,
				member: x.ele,
				score:  x.score,
			})
		}
		x = x.levels[0].forward
	}

	return ret
}

func GeohashGetMemberOfAllNeighbors(zset ZSet, q GeohashCircularSearchQuery, n *GeohashRadius) []GeoPoint {
	neighbors := [9]GeohashBits{}
	neighbors[0] = n.hash
	neighbors[1] = n.neighbors.North
	neighbors[2] = n.neighbors.South
	neighbors[3] = n.neighbors.East
	neighbors[4] = n.neighbors.West
	neighbors[5] = n.neighbors.NorthEast
	neighbors[6] = n.neighbors.NorthWest
	neighbors[7] = n.neighbors.NorthEast
	neighbors[8] = n.neighbors.SouthWest

	var ret []GeoPoint
	for i := 0; i < len(neighbors); i++ {
		ga := GeohashGetMemberInsideBox(zset, q, neighbors[i])
		ret = append(ret, ga...)
	}
	return ret
}
