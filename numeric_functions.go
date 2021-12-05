package tfhe

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

const two32 int64 = int64(1) << 32 // 2^32
var two32Double double = math.Pow(2, 32)

// from double to Torus32 - float64 to int32 conversion
func DoubleToTorus(d double) Torus32 {
	return Torus32(math.Round(math.Mod(d, 1) * math.Pow(2, 32)))
}

// from Torus32 to double
func TorusToDouble(x Torus32) double {
	return double(x) / math.Pow(2, 32)
}

// Used to approximate the phase to the nearest message possible in the message space
// The constant Msize will indicate on which message space we are working (how many messages possible)
//
// "travailler sur 63 bits au lieu de 64, car dans nos cas pratiques, c'est plus précis"
func approxPhase(phase Torus, Msize int) Torus {
	interv := ((uint64(1) << 63) / uint64(Msize)) * 2 // width of each interval
	halfInterval := interv / 2                        // begin of the first intervall
	var phase64 uint64 = (uint64(phase) << 32) + halfInterval
	//floor to the nearest multiples of interv
	phase64 -= phase64 % interv
	//rescale to Torus
	return Torus(phase64 >> 32)
}

func UniformFloat64Dist(min, max int) float64 {
	dist := distuv.Uniform{
		Min: float64(min),
		Max: float64(max),
	}
	return dist.Rand()
}

func UniformTorusDist() Torus {
	dist := distuv.Uniform{
		Min: math.MinInt,
		Max: math.MaxInt,
	}
	return Torus(dist.Rand())
}

func UniformUintDist() uint {
	dist := distuv.Uniform{
		Min: 0,
		Max: math.MaxUint,
	}
	return uint(dist.Rand())
}

func UniformintDist(min, max int) int {
	dist := distuv.Uniform{
		Min: float64(min),
		Max: float64(max),
	}
	return int(dist.Rand())
}

// Gaussian sample centered in message, with standard deviation sigma
func gaussian32(message Torus, sigma double) Torus {
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
	return message + DoubleToTorus(dist.Rand())
}

// Used to approximate the phase to the nearest message possible in the message space
// The constant Msize will indicate on which message space we are working (how many messages possible)
//
// "travailler sur 63 bits au lieu de 64, car dans nos cas pratiques, c'est plus précis"
func ModSwitchFromTorus(phase Torus, Msize int) int {
	interv := ((uint64(1) << 63) / uint64(Msize)) * 2 // width of each intervall
	halfInterval := interv / 2                        // begin of the first intervall
	phase64 := (uint64(phase) << 32) + halfInterval
	//floor to the nearest multiples of interv
	return int(phase64 / interv)
}

// Used to approximate the phase to the nearest message possible in the message space
// The constant Msize will indicate on which message space we are working (how many messages possible)
//
// "travailler sur 63 bits au lieu de 64, car dans nos cas pratiques, c'est plus précis"
func ModSwitchToTorus(mu, Msize int) Torus {
	interv := ((uint64(1) << 63) / uint64(Msize)) * 2 // width of each intervall
	phase64 := uint64(mu) * interv
	//floor to the nearest multiples of interv
	return Torus(phase64 >> 32)
}

//this function return the absolute value of the (centered) fractional part of d
//i.e. the distance between d and its nearest integer
func absfrac(d double) double {
	return math.Abs(d - math.Round(d))
	//return abs(d-rint(d))
}

func Abs(n Torus) Torus {
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
