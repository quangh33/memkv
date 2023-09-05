package data_structure

func GeohashAlign52Bits(hash GeohashBits) uint64 {
	return hash.Bits << (52 - hash.Step*2)
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
