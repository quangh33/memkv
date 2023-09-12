package data_structure

// Implementation of Scalable Bloom Filter data structure
// https://gsd.di.uminho.pt/members/cbm/ps/dbloom.pdf

const ErrorTighteningRatio = 0.5

type SBLink struct {
	bloom *Bloom
	size  uint64 // number of items in the link
}

func (sbl *SBLink) AddHash(hash HashValue) {
	sbl.bloom.AddHash(hash)
	sbl.size++
}

// SBChain A chain of bloom filters
type SBChain struct {
	filters      []SBLink
	size         uint64 // total number of items in all filters
	growthFactor uint32 // growth factor of filter's size
}

func CreateSBChain(initSize uint64, errorRate float64, growthFactor uint32) *SBChain {
	if initSize == 0 || errorRate == 0 || errorRate >= 1 {
		return nil
	}
	sb := &SBChain{
		size:         0,
		growthFactor: growthFactor,
		filters:      []SBLink{},
	}
	sb.AddLink(initSize, errorRate)
	return sb
}

func (sb *SBChain) AddLink(size uint64, errorRate float64) {
	newLink := SBLink{
		size: 0,
	}
	newLink.bloom = CreateBloomFilter(size, errorRate)
	sb.filters = append(sb.filters, newLink)
}

func (sb *SBChain) Add(item string) error {
	hash := sb.filters[0].bloom.CalcHash(item)
	if sb.existHash(hash) {
		return nil
	}
	curFilter := &sb.filters[len(sb.filters)-1]
	if curFilter.size >= curFilter.bloom.Entries {
		newErrorRate := curFilter.bloom.Error * ErrorTighteningRatio
		newSize := curFilter.bloom.Entries * uint64(sb.growthFactor)
		sb.AddLink(newSize, newErrorRate)
		curFilter = &sb.filters[len(sb.filters)-1]
	}
	curFilter.AddHash(hash)
	sb.size++
	return nil
}

func (sb *SBChain) existHash(hash HashValue) bool {
	for i := len(sb.filters) - 1; i >= 0; i-- {
		if sb.filters[i].bloom.ExistHash(hash) {
			return true
		}
	}
	return false
}

func (sb *SBChain) Exist(item string) bool {
	hash := sb.filters[0].bloom.CalcHash(item)
	return sb.existHash(hash)
}
