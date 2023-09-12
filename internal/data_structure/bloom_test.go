package data_structure

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBloom_Exist(t *testing.T) {
	b := CreateBloomFilter(10, 0.01)
	b.Add("a")
	b.Add("b")
	assert.EqualValues(t, 10, b.Entries)
	assert.EqualValues(t, 0.01, b.Error)
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

func TestBloom_AddHash(t *testing.T) {
	b := CreateBloomFilter(10, 0.01)
	hash := b.CalcHash("abcdef")
	b.AddHash(hash)
	assert.True(t, b.ExistHash(hash))
	assert.True(t, b.Exist("abcdef"))
}
