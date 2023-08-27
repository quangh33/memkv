package core

import (
	"errors"
	"fmt"
	"io"
	"memkv/internal/constant"
	"strconv"
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
	default:
		return errors.New(fmt.Sprintf("command not found: %s", cmd.Cmd))
	}
	_, err := c.Write(res)
	return err
}
