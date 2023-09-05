package core

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"memkv/internal/constant"
	"memkv/internal/data_structure"
)

func evalSET(args []string) []byte {
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SET' command"), false)
	}

	var key, value string
	var ttlMs int64 = -1

	key, value = args[0], args[1]
	oType, oEnc := deduceTypeString(value)
	if len(args) > 2 {
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
		}
		ttlMs = ttlSec * 1000
	}

	Put(key, NewObj(value, ttlMs, oType, oEnc))
	return constant.RespOk
}

func evalGET(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command"), false)
	}

	key := args[0]
	obj := Get(key)
	if obj == nil {
		return constant.RespNil
	}

	if hasExpired(obj) {
		return constant.RespNil
	}

	return Encode(obj.Value, false)
}

func evalPING(args []string) []byte {
	var buf []byte

	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'PING' command"), false)
	}

	if len(args) == 0 {
		buf = Encode("PONG", true)
	} else {
		buf = Encode(args[0], false)
	}

	return buf
}

func evalTTL(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'TTL' command"), false)
	}
	key := args[0]
	obj := Get(key)
	if obj == nil {
		return constant.TtlKeyNotExist
	}

	exp, isExpirySet := getExpiry(obj)
	if !isExpirySet {
		return constant.TtlKeyExistNoExpire
	}

	remainMs := exp - uint64(time.Now().UnixMilli())
	if remainMs < 0 {
		return constant.TtlKeyNotExist
	}

	return Encode(int64(remainMs/1000), false)
}

func evalDEL(args []string) []byte {
	delCount := 0

	for _, key := range args {
		if ok := Del(key); ok {
			delCount++
		}
	}

	return Encode(delCount, false)
}

func evalEXPIRE(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'EXPIRE' command"), false)
	}
	key := args[0]
	ttlSec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
	}

	obj := Get(key)
	if obj == nil {
		return constant.RespZero
	}

	setExpiry(obj, ttlSec*1000)
	return constant.RespOne
}

func evalBGREWRITEAOF(args []string) []byte {
	DumpAllAOF()
	return constant.RespOk
}

func evalINCR(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'INCR' command"), false)
	}
	key := args[0]
	obj := Get(key)
	if obj == nil {
		obj = NewObj("0", constant.NoExpire, constant.ObjTypeString, constant.ObjEncodingInt)
		Put(key, obj)
	}

	if err := assertType(obj.TypeEncoding, constant.ObjTypeString); err != nil {
		return Encode(err, false)
	}

	if err := assertEncoding(obj.TypeEncoding, constant.ObjEncodingInt); err != nil {
		return Encode(err, false)
	}

	i, _ := strconv.ParseInt(obj.Value.(string), 10, 64)
	i++
	obj.Value = strconv.FormatInt(i, 10)

	return Encode(i, false)
}

/*
	Set
*/

func evalSADD(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SADD' command"), false)
	}
	key := args[0] // TODO: check key is used by other types or not
	set, exist := setStore[key]
	if !exist {
		set = data_structure.CreateSet(key)
		setStore[key] = set
	}
	count := set.Add(args[1:]...)
	return Encode(count, false)
}

func evalSREM(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SADD' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		set = data_structure.CreateSet(key)
		setStore[key] = set
	}
	count := set.Rem(args[1:]...)
	return Encode(count, false)
}

func evalSCARD(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SCARD' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(0, false)
	}
	return Encode(set.Size(), false)
}

func evalSMEMBERS(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SMEMBERS' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(make([]string, 0), false)
	}
	return Encode(set.Members(), false)
}

func evalSISMEMBER(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SISMEMBER' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(0, false)
	}
	return Encode(set.IsMember(args[1]), false)
}

func evalSMISMEMBER(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SMISMEMBER' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		res := make([]int, len(args)-1)
		return Encode(res, false)
	}
	return Encode(set.MIsMember(args[1:]...), false)
}

func evalSPOP(args []string) []byte {
	if len(args) > 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SPOP' command"), false)
	}
	key := args[0]
	hasCount := len(args) > 1
	count := 0
	if hasCount {
		n, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return Encode(errors.New("(error) Count must be int"), false)
		}
		count = int(n)
	}

	set, exist := setStore[key]
	if !exist {
		if !hasCount {
			return Encode(nil, false)
		}
		return Encode(make([]string, 0), false)
	}
	if !hasCount {
		return Encode(set.Pop(count)[0], false)
	}
	return Encode(set.Pop(count), false)
}

func evalSRAND(args []string) []byte {
	if len(args) > 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SRAND' command"), false)
	}
	key := args[0]
	hasCount := len(args) > 1
	count := 0
	if hasCount {
		n, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return Encode(errors.New("(error) Count must be int"), false)
		}
		count = int(n)
	}

	set, exist := setStore[key]
	if !exist {
		if !hasCount {
			return Encode(nil, false)
		}
		return Encode(make([]string, 0), false)
	}
	if !hasCount {
		return Encode(set.Rand(count)[0], false)
	}
	return Encode(set.Rand(count), false)
}

/*
	Sorted Set
*/

func evalZADD(args []string) []byte {
	if len(args) < 3 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZADD' command"), false)
	}
	key := args[0]
	scoreIndex := 1
	flags := 0
	for scoreIndex < len(args) {
		if strings.ToLower(args[scoreIndex]) == "nx" {
			flags |= data_structure.ZAddInNX
		} else if strings.ToLower(args[scoreIndex]) == "xx" {
			flags |= data_structure.ZAddInXX
		} else {
			break
		}
		scoreIndex++
	}
	nx := (flags & data_structure.ZAddInNX) != 0
	xx := (flags & data_structure.ZAddInXX) != 0
	if nx && xx {
		return Encode(errors.New("(error) Cannot have both NN and XX flag for 'ZADD' command"), false)
	}
	numScoreEleArgs := len(args) - scoreIndex
	if numScoreEleArgs%2 == 1 || numScoreEleArgs == 0 {
		return Encode(errors.New(fmt.Sprintf("(error) Wrong number of (score, member) arg: %d", numScoreEleArgs)), false)
	}

	zset, exist := zsetStore[key]
	if !exist {
		zset = data_structure.CreateZSet()
		zsetStore[key] = zset
	}

	count := 0
	for i := scoreIndex; i < len(args); i += 2 {
		ele := args[i+1]
		score, err := strconv.ParseFloat(args[i], 64)
		if err != nil {
			return Encode(errors.New("(error) Score must be floating point number"), false)
		}
		ret, outFlag := zset.Add(score, ele, flags)
		if ret != 1 {
			return Encode(errors.New("error when adding element"), false)
		}
		if outFlag != data_structure.ZAddOutNop {
			count++
		}
	}
	return Encode(count, false)
}

func evalZRANK(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZRANK' command"), false)
	}
	key, member := args[0], args[1]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespNil
	}
	rank, _ := zset.GetRank(member, false)
	return Encode(rank, false)
}

func evalZREM(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZREM' command"), false)
	}
	key := args[0]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespZero
	}
	deleted := 0
	for i := 1; i < len(args); i++ {
		ret := zset.Del(args[i])
		if ret == 1 {
			deleted++
		}
		if zset.Len() == 0 {
			delete(zsetStore, key)
			break
		}
	}
	return Encode(deleted, false)
}

func evalZSCORE(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZSCORE' command"), false)
	}
	key, member := args[0], args[1]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespNil
	}
	ret, score := zset.GetScore(member)
	if ret == 0 {
		return constant.RespNil
	}
	return Encode(fmt.Sprintf("%f", score), false)
}

func evalZCARD(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZCARD' command"), false)
	}
	key := args[0]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespZero
	}
	return Encode(zset.Len(), false)
}

/*
	Geo Hash
*/

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
		hash := Base32encoding.Encode(value.Bits)
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

func EvalAndResponse(cmd *MemKVCmd, c io.ReadWriter) error {
	var res []byte

	switch cmd.Cmd {
	case "PING":
		res = evalPING(cmd.Args)
	case "SET":
		res = evalSET(cmd.Args)
	case "GET":
		res = evalGET(cmd.Args)
	case "TTL":
		res = evalTTL(cmd.Args)
	case "DEL":
		res = evalDEL(cmd.Args)
	case "EXPIRE":
		res = evalEXPIRE(cmd.Args)
	case "BGREWRITEAOF":
		res = evalBGREWRITEAOF(cmd.Args)
	case "INCR":
		res = evalINCR(cmd.Args)
	// Set
	case "SADD":
		res = evalSADD(cmd.Args)
	case "SREM":
		res = evalSREM(cmd.Args)
	case "SCARD":
		res = evalSCARD(cmd.Args)
	case "SMEMBERS":
		res = evalSMEMBERS(cmd.Args)
	case "SISMEMBER":
		res = evalSISMEMBER(cmd.Args)
	case "SMISMEMBER":
		res = evalSMISMEMBER(cmd.Args)
	case "SRAND":
		res = evalSRAND(cmd.Args)
	case "SPOP":
		res = evalSPOP(cmd.Args)
	// Sorted set
	case "ZADD":
		res = evalZADD(cmd.Args)
	case "ZRANK":
		res = evalZRANK(cmd.Args)
	case "ZREM":
		res = evalZREM(cmd.Args)
	case "ZSCORE":
		res = evalZSCORE(cmd.Args)
	case "ZCARD":
		res = evalZCARD(cmd.Args)
	// Geo Hash
	case "GEOADD":
		res = evalGEOADD(cmd.Args)
	case "GEODIST":
		res = evalGEODIST(cmd.Args)
	case "GEOHASH":
		res = evalGEOHASH(cmd.Args)
	default:
		return errors.New(fmt.Sprintf("command not found: %s", cmd.Cmd))
	}
	_, err := c.Write(res)
	return err
}
