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

var keyValueStore map[string]*Obj

// map from expired obj to its expire time
var keyValueExpireStore map[*Obj]uint64

var zsetStore map[string]*data_structure.ZSet

func init() {
	keyValueStore = make(map[string]*Obj)
	keyValueExpireStore = make(map[*Obj]uint64)
	zsetStore = make(map[string]*data_structure.ZSet)
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
	keyValueExpireStore[obj] = uint64(time.Now().UnixMilli()) + uint64(ttlMs)
}

func Get(k string) *Obj {
	v := keyValueStore[k]
	if v != nil {
		if hasExpired(v) {
			Del(k)
			return nil
		}
	}
	return v
}

func Put(k string, obj *Obj) {
	if len(keyValueStore) >= config.KeyNumberLimit {
		evict()
	}
	keyValueStore[k] = obj
}

func Del(k string) bool {
	if obj, exist := keyValueStore[k]; exist {
		delete(keyValueStore, k)
		delete(keyValueExpireStore, obj)
		return true
	}
	return false
}
