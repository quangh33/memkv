package core

import (
	"testing"
)

func TestCreateSkiplist(t *testing.T) {
	sl := CreateSkiplist()
	if sl == nil {
		t.Fail()
	}
	if sl.length != 0 {
		t.Fail()
	}
	if sl.level != 1 {
		t.Fail()
	}
	if len(sl.head.levels) != SkiplistMaxLevel {
		t.Fail()
	}
	if sl.tail != nil {
		t.Fail()
	}
}
