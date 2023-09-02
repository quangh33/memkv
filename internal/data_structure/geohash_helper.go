package data_structure

type GeoHashFix52Bits uint64

func GeohashAlign52Bits(hash GeohashBits) GeoHashFix52Bits {
	hash.Bits <<= 52 - hash.Step*2
	return GeoHashFix52Bits(hash.Bits)
}
