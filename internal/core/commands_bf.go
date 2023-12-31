package core

import (
	"errors"
	"fmt"
	"memkv/internal/constant"
	"memkv/internal/data_structure"
	"strconv"
)

func cmdBFRESERVE(args []string) []byte {
	if !(len(args) == 3 || len(args) == 5) {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'BF.RESERVE' command"), false)
	}
	key := args[0]
	errRate, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return Encode(errors.New(fmt.Sprintf("error rate must be a floating point number %s", args[1])), false)
	}
	capacity, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return Encode(errors.New(fmt.Sprintf("capacity must be an integer number %s", args[2])), false)
	}
	var growthRate uint64 = data_structure.BfDefaultExpansion
	if len(args) == 5 {
		if args[3] != "EXPANSION" {
			return Encode(errors.New("(error) 4th param must be EXPANSION for 'BF.RESERVE' command"), false)
		}
		growthRate, err = strconv.ParseUint(args[2], 10, 32)
		if err != nil {
			return Encode(errors.New(fmt.Sprintf("growthRate must be an integer number %s", args[2])), false)
		}
		if growthRate < 1 {
			return Encode(errors.New(fmt.Sprintf("growthRate should be greater or equal to 1 %d", growthRate)), false)
		}
	}
	_, exist := sbStore[key]
	if exist {
		return Encode(errors.New(fmt.Sprintf("Bloom filter with key '%s' already exist", key)), false)
	}
	sbStore[key] = data_structure.CreateSBChain(capacity, errRate, growthRate)
	return constant.RespOk
}

func cmdBFINFO(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'BF.INFO' command"), false)
	}
	key := args[0]
	sb, exist := sbStore[key]
	if !exist {
		return Encode(errors.New(fmt.Sprintf("Bloom filter with key '%s' does not exist", key)), false)
	}
	var res []string
	res = append(res, "Capacity", fmt.Sprintf("%d", sb.GetCapacity()),
		"Size", fmt.Sprintf("%d", sb.GetMemUsage()),
		"Number of filters", fmt.Sprintf("%d", sb.GetFilterNumber()),
		"Number of items inserted", fmt.Sprintf("%d", sb.GetSize()),
		"Expansion rate", fmt.Sprintf("%d", sb.GetGrowthFactor()))

	return Encode(res, false)
}

func cmdBFMADD(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'BF.MADD' command"), false)
	}
	key := args[0]
	sb, exist := sbStore[key]
	var err error
	if !exist {
		sb = data_structure.CreateSBChain(data_structure.BfDefaultInitCapacity,
			data_structure.BfDefaultErrRate,
			data_structure.BfDefaultExpansion)
		sbStore[key] = sb
	}
	var res []string
	for i := 1; i < len(args); i++ {
		item := args[i]
		err = sb.Add(item)
		if err != nil {
			res = append(res, "ERR problem inserting into filter")
		} else {
			res = append(res, "1")
		}
	}
	return Encode(res, false)
}

func cmdBFEXISTS(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'BF.EXISTS' command"), false)
	}
	key, item := args[0], args[1]
	sb, exist := sbStore[key]
	if !exist {
		return constant.RespZero
	}
	if !sb.Exist(item) {
		return constant.RespZero
	}
	return constant.RespOne
}

func cmdBFMEXISTS(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'BF.MEXISTS' command"), false)
	}
	key := args[0]
	sb, exist := sbStore[key]
	var res []string
	for i := 1; i < len(args); i++ {
		if !exist {
			res = append(res, "0")
			continue
		}
		item := args[i]
		if !sb.Exist(item) {
			res = append(res, "0")
			continue
		}
		res = append(res, "1")
	}
	return Encode(res, false)
}
