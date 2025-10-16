package utils

import (
	"math"
	"math/rand"

	"github.com/thedonutfactory/go-tfhe/params"
)

// F64ToTorus converts a float64 to a Torus value
func F64ToTorus(d float64) params.Torus {
	torus := math.Mod(d, 1.0) * float64(uint64(1)<<32)
	return params.Torus(int64(torus))
}

// TorusToF64 converts a Torus value to a float64 in range [0, 1)
func TorusToF64(t params.Torus) float64 {
	return float64(t) / float64(uint64(1)<<32)
}

// F64ToTorusVec converts a slice of float64 to a slice of Torus values
func F64ToTorusVec(d []float64) []params.Torus {
	result := make([]params.Torus, len(d))
	for i, val := range d {
		result[i] = F64ToTorus(val)
	}
	return result
}

// GaussianTorus samples from a Gaussian distribution and adds to mu
func GaussianTorus(mu params.Torus, stddev float64, rng *rand.Rand) params.Torus {
	sample := rng.NormFloat64() * stddev
	return mu + F64ToTorus(sample)
}

// GaussianF64 samples from a Gaussian distribution with mean mu
func GaussianF64(mu float64, stddev float64, rng *rand.Rand) params.Torus {
	muTorus := F64ToTorus(mu)
	return GaussianTorus(muTorus, stddev, rng)
}

// GaussianF64Vec samples a vector from a Gaussian distribution
func GaussianF64Vec(mu []float64, stddev float64, rng *rand.Rand) []params.Torus {
	result := make([]params.Torus, len(mu))
	for i, m := range mu {
		result[i] = GaussianTorus(F64ToTorus(m), stddev, rng)
	}
	return result
}
