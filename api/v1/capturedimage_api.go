package v1

import "github.com/lib/pq"

func ConverElevationsFromDB(elevations pq.Float64Array) []float32 {
	res := make([]float32, len(elevations))
	for i, v := range elevations {
		res[i] = float32(v)
	}
	return res
}

func ConverElevationsToDB(elevations []float32) pq.Float64Array {
	res := make(pq.Float64Array, len(elevations))
	for i, v := range elevations {
		res[i] = float64(v)
	}
	return res
}
