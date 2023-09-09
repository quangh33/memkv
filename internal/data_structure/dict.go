package data_structure

import (
	"memkv/internal/config"
	"time"
)

type Obj struct {
	Value        interface{}
	TypeEncoding uint8
	// type    | encoding
	// [][][][]|[][][][]
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[*Obj]uint64
}

func CreateDict() *Dict {
	res := Dict{
		dictStore:        make(map[string]*Obj),
		expiredDictStore: make(map[*Obj]uint64),
	}
	return &res
}

func (d *Dict) NewObj(value interface{}, ttlMs int64, oType uint8, oEnc uint8) *Obj {
	obj := &Obj{
		Value:        value,
		TypeEncoding: oType | oEnc,
	}
	if ttlMs > 0 {
		d.SetExpiry(obj, ttlMs)
	}
	return obj
}

func (d *Dict) HasExpired(obj *Obj) bool {
	exp, exist := d.expiredDictStore[obj]
	if !exist {
		return false
	}
	return exp <= uint64(time.Now().UnixMilli())
}

func (d *Dict) GetExpiry(obj *Obj) (uint64, bool) {
	exp, exist := d.expiredDictStore[obj]
	return exp, exist
}

func (d *Dict) SetExpiry(obj *Obj, ttlMs int64) {
	d.expiredDictStore[obj] = uint64(time.Now().UnixMilli()) + uint64(ttlMs)
}

func (d *Dict) Get(k string) *Obj {
	v := d.dictStore[k]
	if v != nil {
		if d.HasExpired(v) {
			d.Del(k)
			return nil
		}
	}
	return v
}

func (d *Dict) Put(k string, obj *Obj) {
	if len(d.dictStore) >= config.KeyNumberLimit {
		d.evict()
	}
	d.dictStore[k] = obj
}

func (d *Dict) Del(k string) bool {
	if obj, exist := d.dictStore[k]; exist {
		delete(d.dictStore, k)
		delete(d.expiredDictStore, obj)
		return true
	}
	return false
}

func (d *Dict) evictFirst() {
	for k := range d.dictStore {
		d.Del(k)
		return
	}
}

func (d *Dict) evict() {
	switch config.EvictStrategy {
	case config.EvictFirst:
		d.evictFirst()
	default:
		d.evictFirst()
	}
}
