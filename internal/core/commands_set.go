package core

import (
	"errors"
	"memkv/internal/data_structure"
	"strconv"
)

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
