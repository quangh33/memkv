package data_structure

import (
	"github.com/spaolacci/murmur3"
	"math"
)

const Ln2 float64 = 0.693147180559945
const Ln2Square float64 = 0.480453013918201
const ABigSeed uint32 = 0x9747b28c

type Bloom struct {
	Hashes      int
	Entries     uint64
	Error       float64
	bitPerEntry float64
	bf          []uint8
	bits        uint64 // size of bf in bit
	bytes       uint64 // size of bf in byte
}

type HashValue struct {
	a uint64
	b uint64
}

func calcBpe(err float64) float64 {
	num := math.Log(err)
	return math.Abs(-(num / Ln2Square))
}

/*
http://en.wikipedia.org/wiki/Bloom_filter
- Optimal number of bits is: bits = (entries * ln(error)) / ln(2)^2
- bitPerEntry = bits/entries
- Optimal number of hash functions is: hashes = bitPerEntry * ln(2)
*/
func CreateBloomFilter(entries uint64, errorRate float64) *Bloom {
	bloom := Bloom{
		Entries: entries,
		Error:   errorRate,
	}
	bloom.bitPerEntry = calcBpe(errorRate)
	bits := uint64(float64(entries) * bloom.bitPerEntry)
	if bits%64 != 0 {
		bloom.bytes = ((bits / 64) + 1) * 8
	} else {
		bloom.bytes = bits / 8
	}
	bloom.bits = bloom.bytes * 8
	bloom.Hashes = int(math.Ceil(Ln2 * bloom.bitPerEntry))
	bloom.bf = make([]uint8, bloom.bytes)
	return &bloom
}

func (b *Bloom) CalcHash(entry string) HashValue {
	hasher := murmur3.New128WithSeed(ABigSeed)
	hasher.Write([]byte(entry))
	x, y := hasher.Sum128()
	return HashValue{
		a: x,
		b: y,
	}
}

func (b *Bloom) Add(entry string) {
	var hash, bytePos uint64
	initHash := b.CalcHash(entry)
	for i := 0; i < b.Hashes; i++ {
		hash = (initHash.a + initHash.b*uint64(i)) % b.bits
		bytePos = hash >> 3 // div 8
		b.bf[bytePos] |= 1 << (hash % 8)
	}
}

func (b *Bloom) Exist(entry string) bool {
	var hash, bytePos uint64
	initHash := b.CalcHash(entry)
	for i := 0; i < b.Hashes; i++ {
		hash = (initHash.a + initHash.b*uint64(i)) % b.bits
		bytePos = hash >> 3 // div 8
		if (b.bf[bytePos] & (1 << (hash % 8))) == 0 {
			return false
		}
	}
	return true
}

func (b *Bloom) AddHash(initHash HashValue) {
	var hash, bytePos uint64
	for i := 0; i < b.Hashes; i++ {
		hash = (initHash.a + initHash.b*uint64(i)) % b.bits
		bytePos = hash >> 3 // div 8
		b.bf[bytePos] |= 1 << (hash % 8)
	}
}

func (b *Bloom) ExistHash(initHash HashValue) bool {
	var hash, bytePos uint64
	for i := 0; i < b.Hashes; i++ {
		hash = (initHash.a + initHash.b*uint64(i)) % b.bits
		bytePos = hash >> 3 // div 8
		if (b.bf[bytePos] & (1 << (hash % 8))) == 0 {
			return false
		}
	}
	return true
}
