package core

import (
	"errors"
	"fmt"
	"io"
)

func cmdPING(args []string) []byte {
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
		res = cmdPING(cmd.Args)
	case "SET":
		res = cmdSET(cmd.Args)
	case "GET":
		res = cmdGET(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	case "DEL":
		res = cmdDEL(cmd.Args)
	case "EXPIRE":
		res = cmdEXPIRE(cmd.Args)
	case "INCR":
		res = cmdINCR(cmd.Args)
	// Set
	case "SADD":
		res = cmdSADD(cmd.Args)
	case "SREM":
		res = cmdSREM(cmd.Args)
	case "SCARD":
		res = cmdSCARD(cmd.Args)
	case "SMEMBERS":
		res = cmdSMEMBERS(cmd.Args)
	case "SISMEMBER":
		res = cmdSISMEMBER(cmd.Args)
	case "SMISMEMBER":
		res = cmdSMISMEMBER(cmd.Args)
	case "SRAND":
		res = cmdSRAND(cmd.Args)
	case "SPOP":
		res = cmdSPOP(cmd.Args)
	// Sorted set
	case "ZADD":
		res = cmdZADD(cmd.Args)
	case "ZRANK":
		res = cmdZRANK(cmd.Args)
	case "ZREM":
		res = cmdZREM(cmd.Args)
	case "ZSCORE":
		res = cmdZSCORE(cmd.Args)
	case "ZCARD":
		res = cmdZCARD(cmd.Args)
	// Geo Hash
	case "GEOADD":
		res = cmdGEOADD(cmd.Args)
	case "GEODIST":
		res = cmdGEODIST(cmd.Args)
	case "GEOHASH":
		res = cmdGEOHASH(cmd.Args)
	case "GEOSEARCH":
		res = cmdGEOSEARCH(cmd.Args)
	case "GEOPOS":
		res = cmdGEOPOS(cmd.Args)
	// Bloom filter
	case "BF.RESERVE":
		res = cmdBFRESERVE(cmd.Args)
	case "BF.INFO":
		res = cmdBFINFO(cmd.Args)
	case "BF.MADD":
		res = cmdBFMADD(cmd.Args)
	case "BF.EXISTS":
		res = cmdBFEXISTS(cmd.Args)
	case "BF.MEXISTS":
		res = cmdBFMEXISTS(cmd.Args)
	default:
		return errors.New(fmt.Sprintf("command not found: %s", cmd.Cmd))
	}
	_, err := c.Write(res)
	return err
}
