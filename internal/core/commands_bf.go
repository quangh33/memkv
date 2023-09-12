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
