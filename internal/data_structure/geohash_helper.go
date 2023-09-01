package data_structure

type GeoHashFix52Bits uint64

func GeohashAlign52Bits(hash GeohashBits) GeoHashFix52Bits {
	bits := hash.Bits
	bits <<= 52 - hash.Step*2
	return GeoHashFix52Bits(bits)
}
