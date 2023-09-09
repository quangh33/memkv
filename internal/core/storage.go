package core

import (
	"time"

	"memkv/internal/config"
	"memkv/internal/data_structure"
)

type Obj struct {
	Value        interface{}
	TypeEncoding uint8
	// type    | encoding
	// [][][][]|[][][][]
}

var dictStore map[string]*Obj

// map from expired obj to its expiry time
var keyValueExpireStore map[*Obj]uint64

var zsetStore map[string]*data_structure.ZSet

var setStore map[string]data_structure.Set

func init() {
	dictStore = make(map[string]*Obj)
	keyValueExpireStore = make(map[*Obj]uint64)
	zsetStore = make(map[string]*data_structure.ZSet)
	setStore = make(map[string]data_structure.Set)
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
	v := dictStore[k]
	if v != nil {
		if hasExpired(v) {
			Del(k)
			return nil
		}
	}
	return v
}

func GetAndAssert(k string, t uint8) (*Obj, error) {
	v := dictStore[k]
	if v != nil {
		if hasExpired(v) {
			Del(k)
			return nil, nil
		}
	}
	if err := assertType(v.TypeEncoding, t); err != nil {
		return nil, err
	}
	return v, nil
}

func Put(k string, obj *Obj) {
	if len(dictStore) >= config.KeyNumberLimit {
		evict()
	}
	dictStore[k] = obj
}

func Del(k string) bool {
	if obj, exist := dictStore[k]; exist {
		delete(dictStore, k)
		delete(keyValueExpireStore, obj)
		return true
	}
	return false
}
