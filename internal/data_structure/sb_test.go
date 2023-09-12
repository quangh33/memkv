package data_structure

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSBChain(t *testing.T) {
	sb := CreateSBChain(10, 0.01, 2)
	assert.EqualValues(t, 1, len(sb.filters))
	assert.EqualValues(t, 0, sb.filters[0].size)
	assert.EqualValues(t, 10, sb.filters[0].bloom.Entries)
	assert.EqualValues(t, 0.01, sb.filters[0].bloom.Error)
	assert.EqualValues(t, 2, sb.growthFactor)
	assert.EqualValues(t, 0, sb.size)
}

func TestSBChain_AddLink(t *testing.T) {
	sb := CreateSBChain(10, 0.01, 2)
	sb.AddLink(20, 0.005)
	assert.EqualValues(t, 2, len(sb.filters))
	assert.EqualValues(t, 10, sb.filters[0].bloom.Entries)
	assert.EqualValues(t, 0.01, sb.filters[0].bloom.Error)
	assert.EqualValues(t, 20, sb.filters[1].bloom.Entries)
	assert.EqualValues(t, 0.005, sb.filters[1].bloom.Error)
}

func TestSBChain_Add(t *testing.T) {
	var err error

	sb := CreateSBChain(10, 0.01, 2)
	err = sb.Add("0")
	assert.Nil(t, err)
	assert.EqualValues(t, 1, sb.size)
	assert.EqualValues(t, 1, sb.filters[0].size)

	for i := 1; i < 50; i++ {
		err = sb.Add(fmt.Sprintf("%d", i))
		assert.Nil(t, err)
	}
	assert.EqualValues(t, 3, len(sb.filters))
	assert.EqualValues(t, 20, sb.filters[2].size)

	assert.EqualValues(t, 10, sb.filters[0].bloom.Entries)
	assert.EqualValues(t, 0.01, sb.filters[0].bloom.Error)
	assert.EqualValues(t, 20, sb.filters[1].bloom.Entries)
	assert.EqualValues(t, 0.01/2.0, sb.filters[1].bloom.Error)
	assert.EqualValues(t, 40, sb.filters[2].bloom.Entries)
	assert.EqualValues(t, 0.01/4.0, sb.filters[2].bloom.Error)

	for i := 0; i < 50; i++ {
		assert.True(t, sb.Exist(fmt.Sprintf("%d", i)))
	}
	assert.False(t, sb.Exist("50"))
}
