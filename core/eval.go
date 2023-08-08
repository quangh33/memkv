package core

import (
	"errors"
	"fmt"
	"io"
	"memkv/constant"
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
	if len(args) > 2 {
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
		}
		ttlMs = ttlSec * 1000
	}

	Put(key, NewObj(value, ttlMs))
	return constant.RESP_OK
}

func evalGET(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command"), false)
	}

	key := args[0]
	obj := Get(key)
	if obj == nil {
		return constant.RESP_NIL
	}

	if obj.ExpireAt != constant.NO_EXPIRE && obj.ExpireAt <= time.Now().UnixMilli() {
		return constant.RESP_NIL
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
		return constant.TTL_KEY_NOT_EXIST
	}

	if obj.ExpireAt == constant.NO_EXPIRE {
		return constant.TTL_KEY_EXIST_NO_EXPIRE
	}

	remainMs := obj.ExpireAt - time.Now().UnixMilli()
	if remainMs < 0 {
		return constant.TTL_KEY_NOT_EXIST
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
		return constant.RESP_ZERO
	}

	obj.ExpireAt = time.Now().UnixMilli() + ttlSec*1000
	return constant.RESP_ONE
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
	default:
		return errors.New(fmt.Sprintf("command not found: %s", cmd.Cmd))
	}
	_, err := c.Write(res)
	return err
}
