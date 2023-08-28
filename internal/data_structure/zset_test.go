package data_structure

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZSet_Add(t *testing.T) {
	zs := CreateZSet()
	ret, flagOut := zs.Add(10.0, "k1", ZAddInXX)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutNop, flagOut)

	ret, flagOut = zs.Add(10.0, "k1", 0)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutAdded, flagOut)
	v, ok := zs.dict["k1"]
	assert.True(t, ok)
	assert.EqualValues(t, 10, v)
	assert.EqualValues(t, "k1", zs.zskiplist.head.levels[0].forward.ele)
	assert.EqualValues(t, 10, zs.zskiplist.head.levels[0].forward.score)
	assert.EqualValues(t, 1, zs.zskiplist.length)

	ret, flagOut = zs.Add(20.0, "k1", ZAddInNX)
	assert.EqualValues(t, 1, ret)
	assert.EqualValues(t, ZAddOutNop, flagOut)
}
