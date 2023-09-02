package data_structure

type GeoHashFix52Bits uint64

func GeohashAlign52Bits(hash GeohashBits) GeoHashFix52Bits {
	hash.Bits <<= 52 - hash.Step*2
	return GeoHashFix52Bits(hash.Bits)
}

func GeohashEstimateStepsByRadius(radiusMeters float64) uint8 {
	var step uint8 = 1
	for radiusMeters < MercatorMax {
		radiusMeters *= 2
		step++
	}

	step -= 2
	// TODO: handle edge case where we need to search a wider range near the poles
	return step
}
