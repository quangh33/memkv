package core

import (
	"errors"
	"fmt"
	"memkv/internal/constant"
	"memkv/internal/data_structure"
	"memkv/internal/util"
	"strconv"
	"strings"
)

func evalGEOADD(args []string) []byte {
	if len(args) < 4 || len(args)%3 != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GEOADD' command"), false)
	}

	key := args[0]
	zaddArgs := []string{key}
	for i := 1; i < len(args); i += 3 {
		lon, err := strconv.ParseFloat(args[i], 64)
		if err != nil {
			return Encode(errors.New(fmt.Sprintf("lon value must be a floating point number %s\n", args[i])), false)
		}
		lat, err := strconv.ParseFloat(args[i+1], 64)
		if err != nil {
			return Encode(errors.New(fmt.Sprintf("lat value must be a floating point number %s\n", args[i+1])), false)
		}
		member := args[i+2]
		hash, err := data_structure.GeohashEncode(data_structure.GeohashCoordRange, lon, lat, data_structure.GeoMaxStep)
		if err != nil {
			return Encode(err, false)
		}
		bits := data_structure.GeohashAlign52Bits(*hash)
		zaddArgs = append(zaddArgs, fmt.Sprintf("%d", bits))
		zaddArgs = append(zaddArgs, member)
	}
	return evalZADD(zaddArgs)
}

/*
The distance is computed assuming that the Earth is a perfect sphere, so errors up to 0.5% are possible in edge cases.
*/
func evalGEODIST(args []string) []byte {
	if !(len(args) == 3 || len(args) == 4) {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GEODIST' command"), false)
	}
	key, mem1, mem2 := args[0], args[1], args[2]
	var unit float64 = 1
	if len(args) == 4 {
		args[3] = strings.ToLower(args[3])
		if args[3] == "km" {
			unit = 1000
		} else if args[3] == "ft" {
			unit = 0.3048
		} else if args[3] == "mi" {
			unit = 1609.34
		} else {
			return Encode(errors.New("unsupported unit provided. please use M, KM, FT, MI"), false)
		}
	}

	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespNil
	}
	err, score1 := zset.GetScore(mem1)
	if err != 0 {
		return constant.RespNil
	}
	err, score2 := zset.GetScore(mem2)
	if err != 0 {
		return constant.RespNil
	}
	score1GeohashBit := data_structure.GeohashBits{
		Step: data_structure.GeoMaxStep,
		Bits: uint64(score1),
	}
	lon1, lat1 := data_structure.GeohashDecodeAreaToLongLat(data_structure.GeohashCoordRange, score1GeohashBit)
	score2GeohashBit := data_structure.GeohashBits{
		Step: data_structure.GeoMaxStep,
		Bits: uint64(score2),
	}
	lon2, lat2 := data_structure.GeohashDecodeAreaToLongLat(data_structure.GeohashCoordRange, score2GeohashBit)
	dist := data_structure.GeohashGetDistance(lon1, lat1, lon2, lat2) / unit
	return Encode(fmt.Sprintf("%f", dist), false)
}

func evalGEOHASH(args []string) []byte {
	if len(args) < 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GEOHASH' command"), false)
	}
	if len(args) == 1 {
		return constant.RespEmptyArray
	}
	key := args[0]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespNil
	}
	var res []string
	for i := 1; i < len(args); i++ {
		member := args[i]
		err, score := zset.GetScore(member)
		if err != 0 {
			res = append(res, "")
			continue
		}
		scoreGeohashBit := data_structure.GeohashBits{
			Step: data_structure.GeoMaxStep,
			Bits: uint64(score),
		}
		lon, lat := data_structure.GeohashDecodeAreaToLongLat(data_structure.GeohashCoordRange, scoreGeohashBit)
		/* The internal format we use for geocoding is a bit different
		 * than the standard, since we use as initial latitude range
		 * -85,85, while the normal geohashing algorithm uses -90,90.
		 * So we have to decode our position and re-encode using the
		 * standard ranges in order to output a valid geohash string.
		 */
		value, _ := data_structure.GeohashEncode(data_structure.GeohashStandardRange, lon, lat, data_structure.GeoMaxStep)
		value.Bits = data_structure.GeohashAlign52Bits(*value)
		hash := util.Base32encoding.Encode(value.Bits)
		res = append(res, hash)
	}
	return Encode(res, false)
}

/*
GEOSEARCH key [FROMMEMBER member] [FROMLONLAT long lat] radius
TODO: support more options like Redis
*/
func evalGEOSEARCH(args []string) []byte {
	if !(len(args) == 4 || len(args) == 5) {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GEOHASH' command"), false)
	}
	var err error
	var ga []data_structure.GeoPoint
	var res []string
	var member string
	var long, lat float64
	fromMember := false

	key, radiusMeter := args[0], args[len(args)-1]
	if args[1] == "FROMMEMBER" {
		member = args[2]
		fromMember = true
	} else if args[1] == "FROMLONLAT" {
		long, err = strconv.ParseFloat(args[2], 64)
		if err != nil {
			return Encode(errors.New("(error) longitude must be a floating point number"), false)
		}
		lat, err = strconv.ParseFloat(args[3], 64)
		if err != nil {
			return Encode(errors.New("(error) latitude must be a floating point number"), false)
		}
	} else {
		return Encode(errors.New("(error) 2nd param must be FROMMEMBER or FROMLONLAT"), false)
	}
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespEmptyArray
	}
	q := data_structure.GeohashCircularSearchQuery{}
	q.RadiusMeter, err = strconv.ParseFloat(radiusMeter, 64)
	if err != nil {
		return Encode(errors.New("(error) radius must be a floating point number"), false)
	}
	if fromMember {
		memberExist, score := zset.GetScore(member)
		if memberExist < 0 {
			return Encode(errors.New("(error) could not decode requested zset member"), false)
		}
		hash := data_structure.GeohashBits{
			Step: data_structure.GeoMaxStep,
			Bits: uint64(score),
		}
		q.Long, q.Lat = data_structure.GeohashDecodeAreaToLongLat(data_structure.GeohashCoordRange, hash)
	} else {
		q.Long, q.Lat = long, lat
	}

	geohashRadius, err := data_structure.GeohashCalculateSearchingAreas(q)
	if err != nil {
		return Encode(err, false)
	}
	ga = data_structure.GeohashGetMemberOfAllNeighbors(*zset, q, geohashRadius)
	for _, g := range ga {
		res = append(res, g.Member)
	}
	return Encode(res, false)
}

func evalGEOPOS(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GEOPOS' command"), false)
	}
	key := args[0]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespNil
	}
	var res [][]string
	for i := 1; i < len(args); i++ {
		member := args[i]
		memberExist, score := zset.GetScore(member)
		if memberExist < 0 {
			res = append(res, []string{})
			continue
		}
		hash := data_structure.GeohashBits{
			Step: data_structure.GeoMaxStep,
			Bits: uint64(score),
		}
		long, lat := data_structure.GeohashDecodeAreaToLongLat(data_structure.GeohashCoordRange, hash)
		res = append(res, []string{fmt.Sprintf("%f", long), fmt.Sprintf("%f", lat)})
	}
	return Encode(res, false)
}
