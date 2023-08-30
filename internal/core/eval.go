package core

import (
	"errors"
	"fmt"
	"io"
	"memkv/internal/constant"
	"memkv/internal/data_structure"
	"strconv"
	"strings"
	"time"
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
			return Encode(errors.New("Error when adding element"), false)
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
	case "ZADD":
		res = evalZADD(cmd.Args)
	case "ZRANK":
		res = evalZRANK(cmd.Args)
	case "ZREM":
		res = evalZREM(cmd.Args)
	case "ZSCORE":
		res = evalZSCORE(cmd.Args)
	default:
		return errors.New(fmt.Sprintf("command not found: %s", cmd.Cmd))
	}
	_, err := c.Write(res)
	return err
}
