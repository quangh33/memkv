package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvalGEOADD(t *testing.T) {
	res, err := Decode(evalGEOADD([]string{"vn", "10", "20", "p1"}))
	assert.Nil(t, err)
	assert.EqualValues(t, res, 1)

	res, err = Decode(evalGEOADD([]string{"vn", "10", "20", "p1"}))
	assert.Nil(t, err)
	assert.EqualValues(t, res, 0)

	res, err = Decode(evalGEOADD([]string{"vn", "-10", "20", "p1"}))
	assert.Nil(t, err)
	assert.EqualValues(t, res, 1)

	res, err = Decode(evalGEOADD([]string{"vn", "-10", "20", "p2", "-1", "2", "p3"}))
	assert.Nil(t, err)
	assert.EqualValues(t, res, 2)

	zset, exist := zsetStore["vn"]
	assert.True(t, exist)
	assert.EqualValues(t, 3, zset.Len())
}
