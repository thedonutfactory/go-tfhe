// Package num implements various utility functions regarding numeric types.
package num

import (
	"math"
	"math/bits"
)

// Unsigned represents the unsigned Integer type.
type Unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

// Integer represents the Integer type.
type Integer interface {
	Unsigned | int | int8 | int16 | int32 | int64
}

// Real represents the Integer and Float type.
type Real interface {
	Integer | float32 | float64
}

// Number represents Integer, Float, and Complex type.
type Number interface {
	Real | complex64 | complex128
}

// Abs returns the absolute value of x.
func Abs[T Real](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// MaxT returns the maximum possible value of type T in uint64.
func MaxT[T Integer]() uint64 {
	var z T
	switch any(z).(type) {
	case int:
		return math.MaxInt
	case uint, uintptr:
		return math.MaxUint
	case int8:
		return math.MaxInt8
	case uint8:
		return math.MaxUint8
	case int16:
		return math.MaxInt16
	case uint16:
		return math.MaxUint16
	case int32:
		return math.MaxInt32
	case uint32:
		return math.MaxUint32
	case int64:
		return math.MaxInt64
	case uint64:
		return math.MaxUint64
	}
	return math.MaxUint
}

// SizeT returns the bits required to express value of type T.
func SizeT[T Integer]() int {
	var z T
	switch any(z).(type) {
	case int, uint, uintptr:
		return bits.UintSize
	case int8, uint8:
		return 8
	case int16, uint16:
		return 16
	case int32, uint32:
		return 32
	case int64, uint64:
		return 64
	}
	return 64
}

// ByteSizeT returns the bytes required to express value of type T.
func ByteSizeT[T Integer]() int {
	return SizeT[T]() / 8
}

// MinT returns the minimum possible value of type T in int64.
func MinT[T Integer]() int64 {
	var z T
	switch any(z).(type) {
	case int:
		return math.MinInt
	case int8:
		return math.MinInt8
	case int16:
		return math.MinInt16
	case int32:
		return math.MinInt32
	case int64:
		return math.MinInt64
	}
	return 0
}

// IsSigned returns true if type T is a signed type.
func IsSigned[T Real]() bool {
	var z T
	return z-1 < 0
}

// IsPowerOfTwo returns whether x is a power of two.
// If x <= 0, it always returns false.
func IsPowerOfTwo[T Integer](x T) bool {
	return (x > 0) && (x&(x-1)) == 0
}

// Log2 returns floor(log2(x)). Panics if x <= 0.
func Log2[T Integer](x T) int {
	if x <= 0 {
		panic("non-positive log2 undefined")
	}

	return int(bits.Len64(uint64(x))) - 1
}

// DivRound returns round(x/y).
func DivRound[T Integer](x, y T) T {
	return T(math.Round(float64(x) / float64(y)))
}

// DivRoundBits is a bit-optimzed version of RoundRatio: it returns round(x/2^bits).
//
// It only produces correct results for 0 <= bits < sizeT.
// If bits < 0, it panics.
func DivRoundBits[T Integer](x T, bits int) T {
	return (x >> bits) + ((x<<1)>>bits)&1
}

// Min returns the smaller value between x and y.
func Min[T Real](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Max returns the larger value between x and y.
func Max[T Real](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// MaxN returns the largest number of x.
// If x is empty, it returns the zero value of T.
func MaxN[T Real](x ...T) T {
	var max T
	if len(x) == 0 {
		return max
	}

	max = x[0]
	for i := 1; i < len(x); i++ {
		if x[i] > max {
			max = x[i]
		}
	}
	return max
}

// MinN returns the smallest number of x.
// If x is empty, it returns the zero value of T.
func MinN[T Real](x ...T) T {
	var min T
	if len(x) == 0 {
		return min
	}

	min = x[0]
	for i := 1; i < len(x); i++ {
		if x[i] < min {
			min = x[i]
		}
	}
	return min
}

// Sqrt returns floor(sqrt2(x)). Panics if x < 0.
func Sqrt[T Integer](x T) T {
	if x < 0 {
		panic("negative sqrt undefined")
	}

	t := uint64(x)
	c := uint64(0)
	d := uint64(1 << 62)
	for d > t {
		d >>= 2
	}

	for d != 0 {
		if t >= c+d {
			t -= c + d
			c = (c >> 1) + d
		} else {
			c >>= 1
		}
		d >>= 2
	}

	return T(c)
}

// ModInverse returns the modular inverse of x modulo m.
// Output is always positive.
// Panics if m <= 0 or x and m are not coprime.
func ModInverse[T Integer](x, m T) T {
	if m <= 0 {
		panic("modulus not positive")
	}

	x %= m
	if x < 0 {
		x += m
	}

	a, b := x, m
	u, v := T(1), T(0)
	for b != 0 {
		q := a / b
		a, b = b, a-q*b
		u, v = v, u-q*v
	}

	if a != 1 {
		panic("modular inverse does not exist")
	}

	u %= m
	if u < 0 {
		u += m
	}
	return u
}

// ModExp returns x^y mod m.
// Output is always positive.
// Panics if m <= 0.
//
// If y < 0, it returns (x^-1)^(-y) mod m.
// Panics if x and m are not coprime in this case.
func ModExp[T Integer](x, y, m T) T {
	if m <= 0 {
		panic("modulus not positive")
	}

	switch {
	case y < 0:
		x = ModInverse(x, m)
		y = -y
	case y == 0:
		return 1
	}

	x %= m
	if x < 0 {
		x += m
	}

	res := T(1)
	for y > 0 {
		if y&1 == 1 {
			res = (res * x) % m
		}
		x = (x * x) % m
		y >>= 1
	}

	if res < 0 {
		res += m
	}
	return res
}
