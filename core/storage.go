package core

import (
	"memkv/constant"
	"time"
)

type Obj struct {
	Value    interface{}
	ExpireAt int64
}

var store map[string]*Obj

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, ttlMs int64) *Obj {
	var expireAt int64 = constant.NO_EXPIRE
	if ttlMs > 0 {
		expireAt = time.Now().UnixMilli() + ttlMs
	}

	return &Obj{
		Value:    value,
		ExpireAt: expireAt,
	}
}

func Get(k string) *Obj {
	v := store[k]
	if v != nil {
		if v.ExpireAt != constant.NO_EXPIRE && v.ExpireAt <= time.Now().UnixMilli() {
			Del(k)
			return nil
		}
	}
	return v
}

func Put(k string, obj *Obj) {
	store[k] = obj
}

func Del(k string) bool {
	if _, exist := store[k]; exist {
		delete(store, k)
		return true
	}
	return false
}
