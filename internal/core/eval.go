package core

import (
	"errors"
	"fmt"
	"io"
)

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
	case "GEOSEARCH":
		res = evalGEOSEARCH(cmd.Args)
	case "GEOPOS":
		res = evalGEOPOS(cmd.Args)
	default:
		return errors.New(fmt.Sprintf("command not found: %s", cmd.Cmd))
	}
	_, err := c.Write(res)
	return err
}
