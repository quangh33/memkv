package data_structure

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCMS(t *testing.T) {
	cms := CreateCMS(10, 20)
	assert.EqualValues(t, 200, len(cms.counter))
	assert.EqualValues(t, 10, cms.width)
	assert.EqualValues(t, 20, cms.depth)
	assert.EqualValues(t, 0, cms.totalCount)
}

func TestCalcCMSDim(t *testing.T) {
	w, d := CalcCMSDim(0.001, 0.001)
	assert.EqualValues(t, 2000, w)
	assert.EqualValues(t, 10, d)
}

func TestCMS_IncrBy(t *testing.T) {
	cms := CreateCMS(10, 20)
	cms.IncrBy("a", 10)
	assert.EqualValues(t, 10, cms.Count("a"))
	cms.IncrBy("a", 10)
	assert.EqualValues(t, 20, cms.Count("a"))
	cms.IncrBy("b", 30)
	assert.EqualValues(t, 30, cms.Count("b"))
}
