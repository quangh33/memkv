package core

import (
	"errors"
	"fmt"
	"io"
	"memkv/constant"
	"strconv"
	"time"
)

func evalSET(args []string, c io.ReadWriter) error {
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		return errors.New("(error) ERR wrong number of arguments for 'SET' command")
	}

	var key, value string
	var ttlMs int64 = -1

	key, value = args[0], args[1]
	if len(args) > 2 {
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return errors.New("(error) ERR value is not an integer or out of range")
		}
		ttlMs = ttlSec * 1000
	}

	Put(key, NewObj(value, ttlMs))
	c.Write(Encode("OK", true))
	return nil
}

func evalGET(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong number of arguments for 'GET' command")
	}

	key := args[0]
	obj := Get(key)
	if obj == nil {
		c.Write(constant.RESP_NIL)
		return nil
	}

	if obj.ExpireAt != constant.NO_EXPIRE && obj.ExpireAt <= time.Now().UnixMilli() {
		c.Write(constant.RESP_NIL)
		return nil
	}

	c.Write(Encode(obj.Value, false))
	return nil
}

func evalPING(args []string, c io.ReadWriter) error {
	var buf []byte

	if len(args) > 1 {
		return errors.New("ERR wrong number of arguments for 'PING' command")
	}

	if len(args) == 0 {
		buf = Encode("PONG", true)
	} else {
		buf = Encode(args[0], false)
	}

	_, err := c.Write(buf)
	return err
}

func evalTTL(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong number of arguments for 'TTL' command")
	}
	key := args[0]
	obj := Get(key)
	if obj == nil {
		c.Write(constant.TTL_KEY_NOT_EXIST)
		return nil
	}

	if obj.ExpireAt == constant.NO_EXPIRE {
		c.Write(constant.TTL_KEY_EXIST_NO_EXPIRE)
		return nil
	}

	remainMs := obj.ExpireAt - time.Now().UnixMilli()
	if remainMs < 0 {
		c.Write(constant.TTL_KEY_NOT_EXIST)
		return nil
	}

	c.Write(Encode(int64(remainMs/1000), false))
	return nil
}

func EvalAndResponse(cmd *MemKVCmd, c io.ReadWriter) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, c)
	case "SET":
		return evalSET(cmd.Args, c)
	case "GET":
		return evalGET(cmd.Args, c)
	case "TTL":
		return evalTTL(cmd.Args, c)
	}
	return errors.New(fmt.Sprintf("command not found: %s", cmd.Cmd))
}
