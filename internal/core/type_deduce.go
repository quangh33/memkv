package core

import (
	"memkv/internal/constant"
	"strconv"
)

func deduceTypeString(v string) (uint8, uint8) {
	oType := constant.OBJ_TYPE_STRING
	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return oType, constant.OBJ_ENCODING_INT
	}
	return oType, constant.OBJ_ENCODING_RAW
}
