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

func TestSkiplist_GetRank(t *testing.T) {
	sl := CreateSkiplist()
	sl.Insert(10, "k1")
	sl.Insert(20, "k3")
	sl.Insert(50, "k5")
	sl.Insert(40, "k4")
	sl.Insert(10, "k2")
	sl.Insert(50, "k6")

	assert.EqualValues(t, 1, sl.GetRank(10, "k1"))
	assert.EqualValues(t, 2, sl.GetRank(10, "k2"))
	assert.EqualValues(t, 3, sl.GetRank(20, "k3"))
	assert.EqualValues(t, 4, sl.GetRank(40, "k4"))
	assert.EqualValues(t, 5, sl.GetRank(50, "k5"))
	assert.EqualValues(t, 6, sl.GetRank(50, "k6"))
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

func TestZRange_ValueGteMin(t *testing.T) {
	zr := ZRange{
		min:   10,
		max:   20,
		minex: false,
		maxex: false,
	}
	// [10, 20]
	assert.True(t, zr.ValueGteMin(10))
	assert.False(t, zr.ValueGteMin(9.99))

	zr = ZRange{
		min:   10,
		max:   20,
		minex: true,
		maxex: false,
	}
	// (10, 20]
	assert.False(t, zr.ValueGteMin(10))
	assert.True(t, zr.ValueGteMin(10.1))
}

func TestZRange_ValueLteMax(t *testing.T) {
	zr := ZRange{
		min:   10,
		max:   20,
		minex: false,
		maxex: false,
	}
	// [10, 20]
	assert.True(t, zr.ValueLteMax(20))
	assert.False(t, zr.ValueLteMax(20.1))

	zr = ZRange{
		min:   10,
		max:   20,
		minex: false,
		maxex: true,
	}
	// [10, 20)
	assert.False(t, zr.ValueLteMax(20))
	assert.True(t, zr.ValueLteMax(19.99))
}

func TestSkiplist_InRange(t *testing.T) {
	sl := CreateSkiplist()
	sl.Insert(10, "k1")
	sl.Insert(20, "k2")
	sl.Insert(30, "k3")
	sl.Insert(40, "k4")
	// 10->20->30->40
	zr := ZRange{
		min:   5,
		max:   15,
		minex: false,
		maxex: false,
	}
	// [5, 15]
	assert.True(t, sl.InRange(zr))
	zr.min = 5
	zr.max = 10
	// [5, 10]
	assert.True(t, sl.InRange(zr))
	zr.maxex = true
	// [5, 10)
	assert.False(t, sl.InRange(zr))
	zr.min = 40
	zr.max = 41
	zr.minex = false
	zr.maxex = false
	// [40, 41]
	assert.True(t, sl.InRange(zr))
	zr.minex = true
	// (40, 41]
	assert.False(t, sl.InRange(zr))
	zr = ZRange{
		min:   5,
		max:   50,
		minex: false,
		maxex: false,
	}
	// [5, 50]
	assert.True(t, sl.InRange(zr))
	zr = ZRange{
		min:   30,
		max:   20,
		minex: false,
		maxex: false,
	}
	assert.False(t, sl.InRange(zr))

	zr = ZRange{
		min:   20,
		max:   20,
		minex: true,
		maxex: false,
	}
	assert.False(t, sl.InRange(zr))

	zr = ZRange{
		min:   20,
		max:   20,
		minex: false,
		maxex: true,
	}
	assert.False(t, sl.InRange(zr))
}

func TestSkiplist_FindFirstInRange(t *testing.T) {
	sl := CreateSkiplist()
	sl.Insert(10, "k1")
	sl.Insert(20, "k2")
	sl.Insert(30, "k3")
	sl.Insert(40, "k4")
	sl.Insert(50, "k4")

	zr := ZRange{
		min:   1,
		max:   9,
		minex: false,
		maxex: false,
	}
	// [1, 9]
	assert.Nil(t, sl.FindFirstInRange(zr))
	zr = ZRange{
		min:   1,
		max:   15,
		minex: false,
		maxex: false,
	}
	// [1, 15]
	assert.Equal(t, sl.head.levels[0].forward, sl.FindFirstInRange(zr))
	zr = ZRange{
		min:   40,
		max:   100,
		minex: false,
		maxex: false,
	}
	// [40, 100]
	assert.EqualValues(t, 40, sl.FindFirstInRange(zr).score)
	zr = ZRange{
		min:   40,
		max:   100,
		minex: true,
		maxex: false,
	}
	// (40, 100]
	assert.EqualValues(t, 50, sl.FindFirstInRange(zr).score)
	zr = ZRange{
		min:   50,
		max:   100,
		minex: true,
		maxex: false,
	}
	// (50, 100]
	assert.Nil(t, sl.FindFirstInRange(zr))
}
