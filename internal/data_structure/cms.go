package data_structure

import (
	"github.com/spaolacci/murmur3"
	"math"
)

// Implementation of Count-Min Sketch data structure
// https://quanghoang.substack.com/p/count-min-sketch

const Log10PointFive = -0.30102999566

type CMS struct {
	width      uint32
	depth      uint32
	totalCount uint64
	counter    []uint32
}

func CreateCMS(w uint32, d uint32) *CMS {
	cms := &CMS{
		width:      w,
		depth:      d,
		totalCount: 0,
	}
	cms.counter = make([]uint32, d*w)
	return cms
}

// CalcCMSDim calculates the dimension of CMS when we want the error is at most errRate,
// with certainty of (1-errProb)
func CalcCMSDim(errRate float64, errProb float64) (uint32, uint32) {
	w := uint32(math.Ceil(2.0 / errRate))
	d := uint32(math.Ceil(math.Log10(errProb) / Log10PointFive))
	return w, d
}

func (c *CMS) calcHash(item string, seed uint32) uint32 {
	hasher := murmur3.New32WithSeed(seed)
	hasher.Write([]byte(item))
	return hasher.Sum32()
}

func (c *CMS) IncrBy(item string, value uint32) uint32 {
	var i, id, hash uint32
	var minCount uint32 = math.MaxUint32

	for i = 0; i < c.depth; i++ {
		hash = c.calcHash(item, i)
		id = (hash % c.width) + i*c.width
		if math.MaxUint32-c.counter[id] < value {
			c.counter[id] = math.MaxUint32
		} else {
			c.counter[id] += value
		}

		if c.counter[id] < minCount {
			minCount = c.counter[id]
		}
	}
	c.totalCount += uint64(value)
	return minCount
}

func (c *CMS) Count(item string) uint32 {
	var minCount uint32 = math.MaxUint32
	var i, id, hash uint32

	for i = 0; i < c.depth; i++ {
		hash = c.calcHash(item, i)
		id = (hash % c.width) + i*c.width
		if c.counter[id] < minCount {
			minCount = c.counter[id]
		}
	}
	return minCount
}
