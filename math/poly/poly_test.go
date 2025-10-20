package poly_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/thedonutfactory/go-tfhe/math/poly"
)

var (
	LogN = []int{9, 10, 11, 12, 13, 14, 15}
)

func BenchmarkOps(b *testing.B) {
	r := rand.New(rand.NewSource(0))

	for _, logN := range LogN {
		b.Run(fmt.Sprintf("logN=%v", logN), func(b *testing.B) {
			N := 1 << logN

			pev := poly.NewEvaluator[uint64](N)

			p0 := pev.NewPoly()
			p1 := pev.NewPoly()
			pOut := pev.NewPoly()

			for i := 0; i < pev.Degree(); i++ {
				p0.Coeffs[i] = r.Uint64()
				p1.Coeffs[i] = r.Uint64()
			}

			fp := pev.ToFourierPoly(p0)

			b.Run("op=Add", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					pev.AddPolyAssign(p0, p1, pOut)
				}
			})

			b.Run("op=Sub", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					pev.SubPolyAssign(p0, p1, pOut)
				}
			})

			b.Run("op=BinaryMul", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					pev.ShortFourierPolyMulPolyAssign(p0, fp, pOut)
				}
			})

			b.Run("op=Mul", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					pev.MulPolyAssign(p0, p1, pOut)
				}
			})
		})
	}
}

func BenchmarkFourierOps(b *testing.B) {
	r := rand.New(rand.NewSource(0))

	for _, logN := range LogN {
		b.Run(fmt.Sprintf("logN=%v", logN), func(b *testing.B) {
			N := 1 << logN

			pev := poly.NewEvaluator[uint64](N)

			fp0 := pev.NewFourierPoly()
			fp1 := pev.NewFourierPoly()
			fpOut := pev.NewFourierPoly()

			for i := 0; i < pev.Degree(); i++ {
				fp0.Coeffs[i] = (2*r.Float64() - 1.0) * math.Exp(63)
				fp1.Coeffs[i] = (2*r.Float64() - 1.0) * math.Exp(63)
			}

			b.Run("op=Add", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					pev.AddFourierPolyAssign(fp0, fp1, fpOut)
				}
			})

			b.Run("op=Sub", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					pev.SubFourierPolyAssign(fp0, fp1, fpOut)
				}
			})

			b.Run("op=Mul", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					pev.MulFourierPolyAssign(fp0, fp1, fpOut)
				}
			})
		})
	}
}

func BenchmarkFourierTransform(b *testing.B) {
	r := rand.New(rand.NewSource(0))

	for _, logN := range LogN {
		b.Run(fmt.Sprintf("logN=%v", logN), func(b *testing.B) {
			N := 1 << logN

			pev := poly.NewEvaluator[uint64](N)

			p := pev.NewPoly()
			fp := pev.NewFourierPoly()

			for i := 0; i < pev.Degree(); i++ {
				p.Coeffs[i] = r.Uint64()
			}

			b.Run("op=ToFourierPoly", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					pev.ToFourierPolyAssign(p, fp)
				}
			})

			x := N / 3
			b.Run("op=MonomialToFourierPoly", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					pev.MonomialToFourierPolyAssign(x, fp)
				}
			})

			b.Run("op=ToPoly", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					pev.ToPolyAssignUnsafe(fp, p)
				}
			})
		})
	}
}
