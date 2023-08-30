package core

import (
	"memkv/internal/config"
	"memkv/internal/data_structure"
	"time"
)

type Obj struct {
	Value        interface{}
	TypeEncoding uint8
	// type    | encoding
	// [][][][]|[][][][]
}

var store map[string]*Obj

// map from expired obj to its expire time
var expires map[*Obj]uint64

var zsetMap map[string]*data_structure.ZSet

func init() {
	store = make(map[string]*Obj)
	expires = make(map[*Obj]uint64)
	zsetMap = make(map[string]*data_structure.ZSet)
}

func NewObj(value interface{}, ttlMs int64, oType uint8, oEnc uint8) *Obj {
	obj := &Obj{
		Value:        value,
		TypeEncoding: oType | oEnc,
	}
	if ttlMs > 0 {
		setExpiry(obj, ttlMs)
	}
	return obj
}

func setExpiry(obj *Obj, ttlMs int64) {
	expires[obj] = uint64(time.Now().UnixMilli()) + uint64(ttlMs)
}

func Get(k string) *Obj {
	v := store[k]
	if v != nil {
		if hasExpired(v) {
			Del(k)
			return nil
		}
	}
	return v
}

func Put(k string, obj *Obj) {
	if len(store) >= config.KeyNumberLimit {
		evict()
	}
	store[k] = obj
}

func Del(k string) bool {
	if obj, exist := store[k]; exist {
		delete(store, k)
		delete(expires, obj)
		return true
	}
	return false
}
