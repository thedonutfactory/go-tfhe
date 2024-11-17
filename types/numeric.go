package types

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// 2^32 as int32
const two32 int = 1 << 32

var two32Double float64 = math.Pow(2, 32)

// from float64 to Torus32 - float64 to int32 conversion
func DoubleToTorus(d float64) Torus32 {
	return Torus32(math.Round(math.Mod(d, 1) * math.Pow(2, 32)))
}

// from Torus32 to float64
func TorusToDouble(x Torus32) float64 {
	return float64(x) / math.Pow(2, 32)
}

// Used to approximate the phase to the nearest message possible in the message space
// The constant Msize will indicate on which message space we are working (how many messages possible)
//
// "travailler sur 63 bits au lieu de 64, car dans nos cas pratiques, c'est plus précis"
func ApproxPhase(phase Torus32, Msize int32) Torus32 {
	interv := ((uint64(1) << 63) / uint64(Msize)) * 2 // width of each interval
	halfInterval := interv / 2                        // begin of the first intervall
	var phase64 uint64 = (uint64(phase) << 32) + halfInterval
	//floor to the nearest multiples of interv
	phase64 -= phase64 % interv
	//rescale to Torus32
	return Torus32(phase64 >> 32)
}

func UniformFloat64Dist(min, max int32) float64 {
	dist := distuv.Uniform{
		Min: float64(min),
		Max: float64(max),
	}
	return dist.Rand()
}

func UniformTorus32Dist() Torus32 {
	dist := distuv.Uniform{
		Min: math.MinInt32,
		Max: math.MaxInt32,
	}
	return Torus32(dist.Rand())
}

func UniformUintDist() uint32 {
	dist := distuv.Uniform{
		Min: 0,
		Max: math.MaxUint32,
	}
	return uint32(dist.Rand())
}

func UniformInt32Dist(min, max int32) int32 {
	dist := distuv.Uniform{
		Min: float64(min),
		Max: float64(max),
	}
	return int32(dist.Rand())
}

// Gaussian sample centered in message, with standard deviation sigma
func Gaussian32(message Torus32, sigma float64) Torus32 {
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
func ModSwitchFromTorus32(phase Torus32, Msize int32) int32 {
	interv := ((uint64(1) << 63) / uint64(Msize)) * 2 // width of each intervall
	halfInterval := interv / 2                        // begin of the first intervall
	phase64 := (uint64(phase) << 32) + halfInterval
	//floor to the nearest multiples of interv
	return int32(phase64 / interv)
}

// Used to approximate the phase to the nearest message possible in the message space
// The constant Msize will indicate on which message space we are working (how many messages possible)
//
// "travailler sur 63 bits au lieu de 64, car dans nos cas pratiques, c'est plus précis"
func ModSwitchToTorus32(mu, Msize int32) Torus32 {
	interv := ((uint64(1) << 63) / uint64(Msize)) * 2 // width of each intervall
	phase64 := uint64(mu) * interv
	//floor to the nearest multiples of interv
	return Torus32(phase64 >> 32)
}

// this function return the absolute value of the (centered) fractional part of d
// i.e. the distance between d and its nearest integer
func Absfrac(d float64) float64 {
	return math.Abs(d - math.Round(d))
	//return abs(d-rint(d))
}

func Abs(n Torus32) Torus32 {
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
