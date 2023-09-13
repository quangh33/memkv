package core

import (
	"memkv/internal/data_structure"
)

var zsetStore map[string]*data_structure.ZSet
var setStore map[string]data_structure.Set
var dictStore *data_structure.Dict
var sbStore map[string]*data_structure.SBChain
var cmsStore map[string]*data_structure.CMS

func init() {
	zsetStore = make(map[string]*data_structure.ZSet)
	setStore = make(map[string]data_structure.Set)
	dictStore = data_structure.CreateDict()
	sbStore = make(map[string]*data_structure.SBChain)
	cmsStore = make(map[string]*data_structure.CMS)
}
