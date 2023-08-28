package data_structure

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSkiplist(t *testing.T) {
	sl := CreateSkiplist()
	assert.NotNil(t, sl)
	assert.EqualValues(t, 0, sl.length)
	assert.EqualValues(t, 1, sl.level)
	assert.EqualValues(t, SkiplistMaxLevel, len(sl.head.levels))
	assert.Nil(t, sl.tail)
}

func TestSkiplist_Insert(t *testing.T) {
	sl := CreateSkiplist()
	sl.Insert(10, "k1")
	sl.Insert(20, "k3")
	sl.Insert(40, "k4")
	sl.Insert(10, "k2")

	assert.EqualValues(t, 4, sl.length)
	assert.EqualValues(t, 10, sl.head.levels[0].forward.score)
	assert.EqualValues(t, "k1", sl.head.levels[0].forward.ele)
	assert.EqualValues(t, 10, sl.head.levels[0].forward.levels[0].forward.score)
	assert.EqualValues(t, "k2", sl.head.levels[0].forward.levels[0].forward.ele)
	assert.EqualValues(t, 20, sl.head.levels[0].forward.levels[0].forward.levels[0].forward.score)
	assert.EqualValues(t, "k3", sl.head.levels[0].forward.levels[0].forward.levels[0].forward.ele)
	assert.Nil(t, sl.head.levels[0].forward.backward)
	assert.Equal(t, sl.head.levels[0].forward.levels[0].forward.levels[0].forward.levels[0].forward, sl.tail)
	assert.Equal(t, sl.head.levels[0].forward.levels[0].forward.levels[0].forward, sl.tail.backward)
	assert.EqualValues(t, 1, sl.head.levels[0].span)
	for i := sl.level; i < SkiplistMaxLevel; i++ {
		assert.EqualValues(t, 0, sl.head.levels[i].span)
		assert.Nil(t, sl.head.levels[i].forward)
	}
}

func TestSkiplist_Delete(t *testing.T) {
	sl := CreateSkiplist()
	sl.Insert(10, "k1")
	sl.Insert(20, "k3")
	sl.Insert(40, "k4")
	sl.Insert(10, "k2")

	res := sl.Delete(10, "k5")
	assert.EqualValues(t, 0, res)
	res = sl.Delete(30, "k5")
	assert.EqualValues(t, 0, res)

	res = sl.Delete(20, "k3")
	assert.EqualValues(t, 1, res)
	assert.EqualValues(t, 3, sl.length)
	node1 := sl.head.levels[0].forward
	node2 := node1.levels[0].forward
	node3 := node2.levels[0].forward
	assert.EqualValues(t, 10, node1.score)
	assert.EqualValues(t, "k1", node1.ele)
	assert.EqualValues(t, 10, node2.score)
	assert.EqualValues(t, "k2", node2.ele)
	assert.EqualValues(t, 40, node3.score)
	assert.EqualValues(t, "k4", node3.ele)
	assert.EqualValues(t, 10, node3.backward.score)
	assert.EqualValues(t, 10, sl.tail.backward.score)
	assert.Equal(t, sl.tail, node3)
	for i := sl.level; i < SkiplistMaxLevel; i++ {
		assert.Nil(t, sl.head.levels[i].forward)
	}

	res = sl.Delete(10, "k1")
	assert.EqualValues(t, 1, res)
	assert.EqualValues(t, 2, sl.length)
	node1 = sl.head.levels[0].forward
	node2 = node1.levels[0].forward
	assert.EqualValues(t, 10, node1.score)
	assert.EqualValues(t, "k2", node1.ele)
	assert.EqualValues(t, 40, node2.score)
	assert.EqualValues(t, "k4", node2.ele)
	assert.Equal(t, sl.tail, node2)
	for i := sl.level; i < SkiplistMaxLevel; i++ {
		assert.Nil(t, sl.head.levels[i].forward)
	}
}
