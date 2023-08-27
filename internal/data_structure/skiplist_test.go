package data_structure

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

func TestSkiplist_Insert(t *testing.T) {
	sl := CreateSkiplist()
	sl.Insert(10, "k1")
	sl.Insert(20, "k3")
	sl.Insert(40, "k4")
	sl.Insert(10, "k2")
	if sl.length != 4 {
		t.Fatal("len should be 4")
	}
	if sl.head.levels[0].forward.score != 10 {
		t.Fatal("1st node score should be 10")
	}
	if sl.head.levels[0].forward.ele != "k1" {
		t.Fatal("1st node should be k1")
	}
	if sl.head.levels[0].forward.levels[0].forward.score != 10 {
		t.Fatal("2nd node score should be 10")
	}
	if sl.head.levels[0].forward.levels[0].forward.ele != "k2" {
		t.Fatal("2nd node should be k2")
	}
	if sl.head.levels[0].forward.levels[0].forward.levels[0].forward.score != 20 {
		t.Fatal("3rd node score should be 20")
	}
	if sl.head.levels[0].forward.levels[0].forward.levels[0].forward.ele != "k3" {
		t.Fatal("3rd node score should be k3")
	}
	if sl.head.levels[0].forward.backward != nil {
		t.Fail()
	}
	if sl.tail != sl.head.levels[0].forward.levels[0].forward.levels[0].forward.levels[0].forward {
		t.Fatal("tail should be 4th node")
	}
	if sl.head.levels[0].span != 1 {
		t.Fail()
	}
	for i := sl.level; i < SkiplistMaxLevel; i++ {
		if sl.head.levels[i].span != 0 {
			t.Fail()
		}
	}
}
