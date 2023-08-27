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
	assert.EqualValues(t, 1, sl.head.levels[0].span)
	for i := sl.level; i < SkiplistMaxLevel; i++ {
		assert.EqualValues(t, 0, sl.head.levels[i].span)
	}
}
