package poly

import (
	"math/cmplx"
	"math/rand"
	"testing"

	"github.com/thedonutfactory/go-tfhe/math/vec"
)

func TestVecCmplxAssembly(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	N := 1 << 10
	eps := 1e-10

	v0 := make([]complex128, N)
	v1 := make([]complex128, N)
	for i := 0; i < N; i++ {
		v0[i] = complex(r.Float64(), r.Float64())
		v1[i] = complex(r.Float64(), r.Float64())
	}
	v0Float4 := vec.CmplxToFloat4(v0)
	v1Float4 := vec.CmplxToFloat4(v1)

	vOut := make([]complex128, N)
	vOutAVX2 := make([]complex128, N)
	vOutAVX2Float4 := make([]float64, 2*N)

	t.Run("Add", func(t *testing.T) {
		vec.AddAssign(v0, v1, vOut)
		addCmplxAssign(v0Float4, v1Float4, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("Add: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})

	t.Run("Sub", func(t *testing.T) {
		vec.SubAssign(v0, v1, vOut)
		subCmplxAssign(v0Float4, v1Float4, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("Sub: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})

	t.Run("Neg", func(t *testing.T) {
		vec.NegAssign(v0, vOut)
		negCmplxAssign(v0Float4, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("Neg: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})

	t.Run("FloatMul", func(t *testing.T) {
		c := r.Float64()
		vec.ScalarMulAssign(v0, complex(c, 0), vOut)
		floatMulCmplxAssign(v0Float4, c, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("FloatMul: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})

	t.Run("FloatMulAdd", func(t *testing.T) {
		vec.Fill(vOut, 0)
		vec.Fill(vOutAVX2Float4, 0)

		c := r.Float64()
		vec.ScalarMulAddAssign(v0, complex(c, 0), vOut)
		floatMulAddCmplxAssign(v0Float4, c, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("FloatMulAdd: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})

	t.Run("FloatMulSub", func(t *testing.T) {
		vec.Fill(vOut, 0)
		vec.Fill(vOutAVX2Float4, 0)

		c := r.Float64()
		vec.ScalarMulSubAssign(v0, complex(c, 0), vOut)
		floatMulSubCmplxAssign(v0Float4, c, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("FloatMulSub: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})

	t.Run("CmplxMul", func(t *testing.T) {
		c := complex(r.Float64(), r.Float64())
		vec.ScalarMulAssign(v0, c, vOut)
		cmplxMulCmplxAssign(v0Float4, c, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("CmplxMul: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})

	t.Run("CmplxMulAdd", func(t *testing.T) {
		vec.Fill(vOut, 0)
		vec.Fill(vOutAVX2Float4, 0)

		c := complex(r.Float64(), r.Float64())
		vec.ScalarMulAddAssign(v0, c, vOut)
		cmplxMulAddCmplxAssign(v0Float4, c, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("CmplxMulAdd: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})

	t.Run("CmplxMulSub", func(t *testing.T) {
		vec.Fill(vOut, 0)
		vec.Fill(vOutAVX2Float4, 0)

		c := complex(r.Float64(), r.Float64())
		vec.ScalarMulSubAssign(v0, c, vOut)
		cmplxMulSubCmplxAssign(v0Float4, c, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("CmplxMulSub: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})

	t.Run("Mul", func(t *testing.T) {
		vec.ElementWiseMulAssign(v0, v1, vOut)
		elementWiseMulCmplxAssign(v0Float4, v1Float4, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("Mul: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})

	t.Run("MulAdd", func(t *testing.T) {
		vec.Fill(vOut, 0)
		vec.Fill(vOutAVX2Float4, 0)

		vec.ElementWiseMulAddAssign(v0, v1, vOut)
		elementWiseMulAddCmplxAssign(v0Float4, v1Float4, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("MulAdd: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})

	t.Run("MulSub", func(t *testing.T) {
		vec.Fill(vOut, 0)
		vec.Fill(vOutAVX2Float4, 0)

		vec.ElementWiseMulSubAssign(v0, v1, vOut)
		elementWiseMulSubCmplxAssign(v0Float4, v1Float4, vOutAVX2Float4)
		vec.Float4ToCmplxAssign(vOutAVX2Float4, vOutAVX2)
		for i := 0; i < N; i++ {
			if cmplx.Abs(vOut[i]-vOutAVX2[i]) > eps {
				t.Fatalf("MulSub: %v != %v", vOut[i], vOutAVX2[i])
			}
		}
	})
}
