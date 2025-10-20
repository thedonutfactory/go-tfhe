package vec_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thedonutfactory/go-tfhe/math/vec"
)

func TestVec(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	N := 687

	v0 := make([]uint32, N)
	v1 := make([]uint32, N)
	vOut := make([]uint32, N)
	vOutAVX := make([]uint32, N)

	for i := 0; i < N; i++ {
		v0[i] = r.Uint32()
		v1[i] = r.Uint32()
	}

	w0 := make([]uint64, N)
	w1 := make([]uint64, N)
	wOut := make([]uint64, N)
	wOutAVX := make([]uint64, N)

	for i := 0; i < N; i++ {
		w0[i] = r.Uint64()
		w1[i] = r.Uint64()
	}

	t.Run("AddAssign", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] = v0[i] + v1[i]
			wOut[i] = w0[i] + w1[i]
		}
		vec.AddAssign(v0, v1, vOutAVX)
		vec.AddAssign(w0, w1, wOutAVX)

		assert.Equal(t, vOut, vOutAVX)
		assert.Equal(t, wOut, wOutAVX)
	})

	t.Run("SubAssign", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] = v0[i] - v1[i]
			wOut[i] = w0[i] - w1[i]
		}
		vec.SubAssign(v0, v1, vOutAVX)
		vec.SubAssign(w0, w1, wOutAVX)

		assert.Equal(t, vOut, vOutAVX)
		assert.Equal(t, wOut, wOutAVX)
	})

	t.Run("ScalarMulAssign", func(t *testing.T) {
		cv := vOut[0]
		cw := wOut[0]
		for i := 0; i < N; i++ {
			vOut[i] = cv * v0[i]
			wOut[i] = cw * w0[i]
		}
		vec.ScalarMulAssign(v0, cv, vOutAVX)
		vec.ScalarMulAssign(w0, cw, wOutAVX)

		assert.Equal(t, vOut, vOutAVX)
		assert.Equal(t, wOut, wOutAVX)
	})

	t.Run("ScalarMulAddAssign", func(t *testing.T) {
		cv := vOut[0]
		cw := wOut[0]
		for i := 0; i < N; i++ {
			vOut[i] += cv * v0[i]
			wOut[i] += cw * w0[i]
		}
		vec.ScalarMulAddAssign(v0, cv, vOutAVX)
		vec.ScalarMulAddAssign(w0, cw, wOutAVX)

		assert.Equal(t, vOut, vOutAVX)
		assert.Equal(t, wOut, wOutAVX)
	})

	t.Run("ScalarMulSubAssign", func(t *testing.T) {
		cv := vOut[0]
		cw := wOut[0]
		for i := 0; i < N; i++ {
			vOut[i] -= cv * v0[i]
			wOut[i] -= cw * w0[i]
		}
		vec.ScalarMulSubAssign(v0, cv, vOutAVX)
		vec.ScalarMulSubAssign(w0, cw, wOutAVX)

		assert.Equal(t, vOut, vOutAVX)
		assert.Equal(t, wOut, wOutAVX)
	})

	t.Run("ElementWiseMulAssign", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] = v0[i] * v1[i]
			wOut[i] = w0[i] * w1[i]
		}
		vec.ElementWiseMulAssign(v0, v1, vOutAVX)
		vec.ElementWiseMulAssign(w0, w1, wOutAVX)

		assert.Equal(t, vOut, vOutAVX)
		assert.Equal(t, wOut, wOutAVX)
	})

	t.Run("ElementWiseMulAddAssign", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] += v0[i] * v1[i]
			wOut[i] += w0[i] * w1[i]
		}
		vec.ElementWiseMulAddAssign(v0, v1, vOutAVX)
		vec.ElementWiseMulAddAssign(w0, w1, wOutAVX)

		assert.Equal(t, vOut, vOutAVX)
		assert.Equal(t, wOut, wOutAVX)
	})

	t.Run("ElementWiseMulSubAssign", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] -= v0[i] * v1[i]
			wOut[i] -= w0[i] * w1[i]
		}
		vec.ElementWiseMulSubAssign(v0, v1, vOutAVX)
		vec.ElementWiseMulSubAssign(w0, w1, wOutAVX)

		assert.Equal(t, vOut, vOutAVX)
		assert.Equal(t, wOut, wOutAVX)
	})
}
