package tfhe

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

const two32 int64 = int64(1) << 32 // 2^32
var two32Double double = math.Pow(2, 32)

// from double to Torus - float64 to int64 conversion
func Dtot32(d double) Torus {
	return Torus(math.Round(math.Mod(d, 1) * math.Pow(2, 64)))
}

// from Torus to double
func T32tod(x Torus) double {
	return double(x) / math.Pow(2, 64)
}

// Used to approximate the phase to the nearest message possible in the message space
// The constant Msize will indicate on which message space we are working (how many messages possible)
//
// "travailler sur 63 bits au lieu de 64, car dans nos cas pratiques, c'est plus précis"
func approxPhase(phase int64, Msize int) int64 {
	interv := ((uint64(1) << 63) / uint64(Msize)) * 2 // width of each interval
	halfInterval := interv / 2                        // begin of the first intervall
	// var phase64 uint64 = (uint64(phase) << 32) + halfInterval
	var phase64 uint64 = uint64(phase) + halfInterval
	//floor to the nearest multiples of interv
	phase64 -= phase64 % interv
	//rescale to Torus
	return int64(phase64)
}

func UniformFloat64Dist(min, max int64) float64 {
	dist := distuv.Uniform{
		Min: float64(min),
		Max: float64(max),
	}
	return dist.Rand()
}

func UniformTorusDist() int64 {
	dist := distuv.Uniform{
		Min: math.MinInt64,
		Max: math.MaxInt64,
	}
	return int64(dist.Rand())
}

func UniformUintDist() uint64 {
	dist := distuv.Uniform{
		Min: 0,
		Max: math.MaxUint64,
	}
	return uint64(dist.Rand())
}

func Uniformint64Dist(min, max int64) int64 {
	dist := distuv.Uniform{
		Min: float64(min),
		Max: float64(max),
	}
	return int64(dist.Rand())
}

// Gaussian sample centered in message, with standard deviation sigma
func gaussian32(message int64, sigma double) Torus {
	// Create a standard normal (mean = 0, stdev = 1)
	dist := distuv.Normal{
		Mu:    0,     // Mean of the normal distribution
		Sigma: sigma, // Standard deviation of the normal distribution
	}
	/*
		    z := make([]float64, n)

			// Generate some random numbers from standard normal distribution
			for i := range z {
				z[i] = dist.Rand()
			}
	*/
	return message + Dtot32(dist.Rand())
}

// Used to approximate the phase to the nearest message possible in the message space
// The constant Msize will indicate on which message space we are working (how many messages possible)
//
// "travailler sur 63 bits au lieu de 64, car dans nos cas pratiques, c'est plus précis"
func ModSwitchFromTorus(phase Torus, Msize int) int64 {
	interv := ((uint64(1) << 63) / uint64(Msize)) * 2 // width of each intervall
	halfInterval := interv / 2                        // begin of the first intervall
	//phase64 := (uint64(phase) << 32) + halfInterval
	phase64 := uint64(phase) + halfInterval
	//floor to the nearest multiples of interv
	return int64(phase64 / interv)
}

// Used to approximate the phase to the nearest message possible in the message space
// The constant Msize will indicate on which message space we are working (how many messages possible)
//
// "travailler sur 63 bits au lieu de 64, car dans nos cas pratiques, c'est plus précis"
func ModSwitchToTorus(mu int64, Msize int) Torus {
	interv := ((uint64(1) << 63) / uint64(Msize)) * 2 // width of each intervall
	phase64 := uint64(mu) * interv
	//floor to the nearest multiples of interv
	//return int64(phase64 >> 32)
	return Torus(phase64)
}

//this function return the absolute value of the (centered) fractional part of d
//i.e. the distance between d and its nearest integer
func absfrac(d double) double {
	return math.Abs(d - math.Round(d))
	//return abs(d-rint(d))
}

func Abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
