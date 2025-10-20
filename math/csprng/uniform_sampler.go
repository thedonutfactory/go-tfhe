package csprng

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"

	"github.com/thedonutfactory/go-tfhe/math/num"
	"github.com/thedonutfactory/go-tfhe/math/poly"
)

// bufSize is the default buffer size of UniformSampler.
const bufSize = 8192

// UniformSampler samples values from uniform distribution.
// This uses AES-CTR as a underlying prng.
type UniformSampler[T num.Integer] struct {
	prng cipher.Stream

	buf [bufSize]byte
	ptr int

	byteSizeT int
	maxT      T
}

// NewUniformSampler creates a new UniformSampler.
//
// Panics when read from crypto/rand or AES initialization fails.
func NewUniformSampler[T num.Integer]() *UniformSampler[T] {
	var seed [32]byte
	if _, err := rand.Read(seed[:]); err != nil {
		panic(err)
	}

	return NewUniformSamplerWithSeed[T](seed[:])
}

// NewUniformSamplerWithSeed creates a new UniformSampler, with user supplied seed.
//
// Panics when AES initialization fails.
func NewUniformSamplerWithSeed[T num.Integer](seed []byte) *UniformSampler[T] {
	r := sha512.Sum384(seed)

	block, err := aes.NewCipher(r[:32])
	if err != nil {
		panic(err)
	}

	prng := cipher.NewCTR(block, r[32:])

	return &UniformSampler[T]{
		prng: prng,

		buf: [bufSize]byte{},
		ptr: bufSize,

		byteSizeT: num.ByteSizeT[T](),
		maxT:      T(num.MaxT[T]()),
	}
}

// Sample uniformly samples a random integer of type T.
func (s *UniformSampler[T]) Sample() T {
	if s.ptr == bufSize {
		s.prng.XORKeyStream(s.buf[:], s.buf[:])
		s.ptr = 0
	}

	var res uint64
	switch s.byteSizeT {
	case 1:
		res |= uint64(s.buf[s.ptr+0])
	case 2:
		res |= uint64(s.buf[s.ptr+0])
		res |= uint64(s.buf[s.ptr+1]) << 8
	case 4:
		res |= uint64(s.buf[s.ptr+0])
		res |= uint64(s.buf[s.ptr+1]) << 8
		res |= uint64(s.buf[s.ptr+2]) << 16
		res |= uint64(s.buf[s.ptr+3]) << 24
	case 8:
		res |= uint64(s.buf[s.ptr+0])
		res |= uint64(s.buf[s.ptr+1]) << 8
		res |= uint64(s.buf[s.ptr+2]) << 16
		res |= uint64(s.buf[s.ptr+3]) << 24
		res |= uint64(s.buf[s.ptr+4]) << 32
		res |= uint64(s.buf[s.ptr+5]) << 40
		res |= uint64(s.buf[s.ptr+6]) << 48
		res |= uint64(s.buf[s.ptr+7]) << 56
	}
	s.ptr += s.byteSizeT

	return T(res)
}

// SampleN uniformly samples a random integer of type T in [0, N).
func (s *UniformSampler[T]) SampleN(N T) T {
	bound := s.maxT - (s.maxT % N)
	for {
		res := s.Sample()
		if 0 <= res && res < bound {
			return res % N
		}
	}
}

// SampleVecAssign samples uniform values to vOut.
func (s *UniformSampler[T]) SampleVecAssign(vOut []T) {
	for i := range vOut {
		vOut[i] = s.Sample()
	}
}

// SamplePolyAssign samples uniform values to p.
func (s *UniformSampler[T]) SamplePolyAssign(pOut poly.Poly[T]) {
	s.SampleVecAssign(pOut.Coeffs)
}
