package core

import "time"

func hasExpired(obj *Obj) bool {
	exp, exist := expires[obj]
	if !exist {
		return false
	}
	return exp <= uint64(time.Now().UnixMilli())
}

func getExpiry(obj *Obj) (uint64, bool) {
	exp, exist := expires[obj]
	return exp, exist
}
