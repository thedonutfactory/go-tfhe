package csprng

import (
	"github.com/thedonutfactory/go-tfhe/math/num"
	"github.com/thedonutfactory/go-tfhe/math/poly"
	"github.com/thedonutfactory/go-tfhe/math/vec"
)

// BinarySampler samples values from uniform and block binary distribution.
type BinarySampler[T num.Integer] struct {
	baseSampler *UniformSampler[uint64]
}

// NewBinarySampler creates a new BinarySampler.
//
// Panics when read from crypto/rand or AES initialization fails.
func NewBinarySampler[T num.Integer]() *BinarySampler[T] {
	return &BinarySampler[T]{
		baseSampler: NewUniformSampler[uint64](),
	}
}

// NewBinarySamplerWithSeed creates a new BinarySampler, with user supplied seed.
//
// Panics when AES initialization fails.
func NewBinarySamplerWithSeed[T num.Integer](seed []byte) *BinarySampler[T] {
	return &BinarySampler[T]{
		baseSampler: NewUniformSamplerWithSeed[uint64](seed),
	}
}

// Sample uniformly samples a random binary integer.
func (s *BinarySampler[T]) Sample() T {
	return T(s.baseSampler.Sample() & 1)
}

// SampleVecAssign samples uniform binary values to vOut.
func (s *BinarySampler[T]) SampleVecAssign(vOut []T) {
	var buf uint64
	for i := 0; i < len(vOut); i++ {
		if i&63 == 0 {
			buf = s.baseSampler.Sample()
		}
		vOut[i] = T(buf & 1)
		buf >>= 1
	}
}

// SamplePolyAssign samples uniform binary values to pOut.
func (s *BinarySampler[T]) SamplePolyAssign(pOut poly.Poly[T]) {
	s.SampleVecAssign(pOut.Coeffs)
}

// SampleBlockVecAssign samples block binary values to vOut.
func (s *BinarySampler[T]) SampleBlockVecAssign(blockSize int, vOut []T) {
	if len(vOut)%blockSize != 0 {
		panic("length not multiple of blocksize")
	}

	for i := 0; i < len(vOut); i += blockSize {
		vec.Fill(vOut[i:i+blockSize], 0)
		offset := int(s.baseSampler.SampleN(uint64(blockSize) + 1))
		if offset == blockSize {
			continue
		}
		vOut[i+offset] = 1
	}
}

// SampleBlockPolyAssign samples block binary values to pOut.
func (s *BinarySampler[T]) SampleBlockPolyAssign(blockSize int, pOut poly.Poly[T]) {
	s.SampleBlockVecAssign(blockSize, pOut.Coeffs)
}
