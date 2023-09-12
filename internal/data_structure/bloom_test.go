package data_structure

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBloom_Exist(t *testing.T) {
	b := CreateBloomFilter(10, 0.01)
	b.Add("a")
	b.Add("b")
	assert.True(t, b.Exist("a"))
	assert.True(t, b.Exist("b"))
	assert.False(t, b.Exist("c"))
	assert.False(t, b.Exist("d"))
}

func TestBloom_CalcHash(t *testing.T) {
	b := CreateBloomFilter(10, 0.01)
	x := b.CalcHash("abcdef")
	y := b.CalcHash("abcdef")
	assert.EqualValues(t, x.a, y.a)
	assert.EqualValues(t, x.b, y.b)
}
