package core

type GeoHashFix52Bits uint64

func geohashAlign52Bits(hash GeohashBits) GeoHashFix52Bits {
	bits := hash.Bits
	bits <<= 52 - hash.Step*2
	return GeoHashFix52Bits(bits)
}
