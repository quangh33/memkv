package core

import (
	"memkv/internal/data_structure"
)

type Obj struct {
	Value        interface{}
	TypeEncoding uint8
	// type    | encoding
	// [][][][]|[][][][]
}

var zsetStore map[string]*data_structure.ZSet

var setStore map[string]data_structure.Set

var dict *data_structure.Dict

func init() {
	zsetStore = make(map[string]*data_structure.ZSet)
	setStore = make(map[string]data_structure.Set)
	dict = data_structure.CreateDict()
}
