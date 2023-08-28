package data_structure

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZSet_Add_NoOps(t *testing.T) {
	zs := CreateZSet()
	ret, flagOut := zs.Add(10.0, "k1", ZAddInXX)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutNop, flagOut)

	ret, flagOut = zs.Add(10.0, "k1", 0)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutAdded, flagOut)

	ret, flagOut = zs.Add(20.0, "k1", ZAddInNX)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutNop, flagOut)

	ret, flagOut = zs.Add(100.0, "", ZAddInNX)
	assert.EqualValues(t, 0, ret)
	assert.EqualValues(t, ZAddOutNop, flagOut)
}

func TestZSet_Add_AddNew(t *testing.T) {
	zs := CreateZSet()
	ret, flagOut := zs.Add(10.0, "k1", 0)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutAdded, flagOut)
	v, ok := zs.dict["k1"]
	assert.True(t, ok)
	assert.EqualValues(t, 10.0, v)
	assert.EqualValues(t, "k1", zs.zskiplist.head.levels[0].forward.ele)
	assert.EqualValues(t, 10, zs.zskiplist.head.levels[0].forward.score)
	assert.EqualValues(t, 1, zs.zskiplist.length)

	ret, flagOut = zs.Add(20.0, "k2", 0)
	v, ok = zs.dict["k2"]
	assert.EqualValues(t, 1, ret)
	assert.True(t, ok)
	assert.EqualValues(t, 20, v)
	assert.EqualValues(t, "k2", zs.zskiplist.tail.ele)
	assert.EqualValues(t, 20, zs.zskiplist.tail.score)
	assert.EqualValues(t, 2, zs.zskiplist.length)
}

func TestZSet_Add_UpdateExist(t *testing.T) {
	zs := CreateZSet()
	ret, flagOut := zs.Add(10.0, "k1", 0)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutAdded, flagOut)
	v, ok := zs.dict["k1"]
	assert.True(t, ok)
	assert.EqualValues(t, 10.0, v)
	assert.EqualValues(t, "k1", zs.zskiplist.head.levels[0].forward.ele)
	assert.EqualValues(t, 10, zs.zskiplist.head.levels[0].forward.score)
	assert.EqualValues(t, 1, zs.zskiplist.length)

	ret, flagOut = zs.Add(5.0, "k1", 0)
	v, ok = zs.dict["k1"]
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutUpdated, flagOut)
	assert.True(t, ok)
	assert.EqualValues(t, 5, v)
	assert.EqualValues(t, "k1", zs.zskiplist.head.levels[0].forward.ele)
	assert.EqualValues(t, 5, zs.zskiplist.head.levels[0].forward.score)
	assert.EqualValues(t, 1, zs.zskiplist.length)
}

func TestZSet_Add_AddDuplicateEleAndScore(t *testing.T) {
	zs := CreateZSet()
	ret, flagOut := zs.Add(10.0, "k1", 0)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutAdded, flagOut)
	v, ok := zs.dict["k1"]
	assert.True(t, ok)
	assert.EqualValues(t, 10.0, v)
	assert.EqualValues(t, "k1", zs.zskiplist.head.levels[0].forward.ele)
	assert.EqualValues(t, 10, zs.zskiplist.head.levels[0].forward.score)
	assert.EqualValues(t, 1, zs.zskiplist.length)

	ret, flagOut = zs.Add(10.0, "k1", 0)
	v, ok = zs.dict["k1"]
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutNop, flagOut)
	assert.True(t, ok)
}

func TestZSet_Del(t *testing.T) {
	zs := CreateZSet()
	ret, flagOut := zs.Add(20.0, "k2", 0)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutAdded, flagOut)
	ret, flagOut = zs.Add(10.0, "k1", 0)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutAdded, flagOut)
	ret, flagOut = zs.Add(30.0, "k3", 0)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutAdded, flagOut)

	assert.EqualValues(t, 3, zs.zskiplist.length)
	zs.Del("k1")
	assert.EqualValues(t, 2, zs.zskiplist.length)
	assert.EqualValues(t, 20.0, zs.zskiplist.head.levels[0].forward.score)
	assert.EqualValues(t, "k2", zs.zskiplist.head.levels[0].forward.ele)
	zs.Del("k3")
	assert.EqualValues(t, 1, zs.zskiplist.length)
	assert.EqualValues(t, 20.0, zs.zskiplist.head.levels[0].forward.score)
	assert.EqualValues(t, "k2", zs.zskiplist.head.levels[0].forward.ele)
	zs.Del("k2")
	assert.EqualValues(t, 0, zs.zskiplist.length)
	assert.Nil(t, zs.zskiplist.head.levels[0].forward)
	assert.Nil(t, zs.zskiplist.tail)
}

func TestZSet_GetRank(t *testing.T) {
	zs := CreateZSet()
	_, flagOut := zs.Add(20.0, "k2", 0)
	assert.EqualValues(t, ZAddOutAdded, flagOut)
	_, flagOut = zs.Add(40.0, "k4", 0)
	assert.EqualValues(t, ZAddOutAdded, flagOut)
	_, flagOut = zs.Add(10.0, "k1", 0)
	assert.EqualValues(t, ZAddOutAdded, flagOut)
	_, flagOut = zs.Add(15.0, "k2", 0)
	assert.EqualValues(t, ZAddOutUpdated, flagOut)
	_, flagOut = zs.Add(30.0, "k3", 0)
	assert.EqualValues(t, ZAddOutAdded, flagOut)

	/* { (k1, 10), (k2, 15), (k3, 30), (k4, 40) } */
	rank, score := zs.GetRank("k1", false)
	assert.EqualValues(t, 0, rank)
	assert.EqualValues(t, 10.0, score)

	rank, score = zs.GetRank("k2", false)
	assert.EqualValues(t, 1, rank)
	assert.EqualValues(t, 15.0, score)

	rank, score = zs.GetRank("k3", false)
	assert.EqualValues(t, 2, rank)
	assert.EqualValues(t, 30.0, score)

	rank, score = zs.GetRank("k4", false)
	assert.EqualValues(t, 3, rank)
	assert.EqualValues(t, 40.0, score)
}
