package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thedonutfactory/go-tfhe/types"
)

func TestGaussian32(t *testing.T) {
	assert := assert.New(t)

	const MESSAGE1 types.Torus32 = 123456789
	const MESSAGE2 types.Torus32 = 987654321

	reps1 := types.Gaussian32(MESSAGE1, 0)
	reps2 := types.Gaussian32(MESSAGE2, 0)
	assert.Equal(MESSAGE1, reps1)
	assert.Equal(MESSAGE2, reps2)

	reps1 = types.Gaussian32(MESSAGE1, 0.01)
	reps2 = types.Gaussian32(MESSAGE2, 0.5)
	assert.NotEqual(MESSAGE1, reps1)
	assert.NotEqual(MESSAGE2, reps2)
	assert.LessOrEqual(types.Abs(MESSAGE1-reps1), int32(80000000))
}

func TestConversion(t *testing.T) {
	assert := assert.New(t)
	// conversion from float64 to Torus32
	// conversion from Torus32 to float64
	assert.EqualValues(int32(0), types.DoubleToTorus(0))
	assert.EqualValues(1<<31, types.DoubleToTorus(-0.5))
	assert.EqualValues(1<<30, types.DoubleToTorus(0.25))
	assert.EqualValues(0xC0000000, types.DoubleToTorus(-0.25))
}

// Used to approximate the phase to the nearest multiple of  1/Msize
func TestApproxPhase(t *testing.T) {
	assert := assert.New(t)
	for i := int32(2); i < 200; i++ {
		v := types.UniformTorus32Dist()
		w := types.ApproxPhase(v, i)
		dv := types.TorusToDouble(v)
		dw := types.TorusToDouble(w)
		// fmt.Printf("%d, %f, %f, %f \n", i, dv, dw, float64(i)*dw)
		assert.LessOrEqual(types.Absfrac(dv-dw), 1./(2.*float64(i))+1e-40)
		assert.LessOrEqual(types.Absfrac(float64(i)*dw), float64(i)*1e-9)
	}
}

func TestModSwitchFromTorus32(t *testing.T) {
	assert := assert.New(t)

	for i := 2; i < 200; i++ {
		v := types.UniformTorus32Dist()
		w := types.ModSwitchFromTorus32(v, int32(i))
		dv := types.TorusToDouble(v)
		dw := float64(w) / float64(i)
		assert.LessOrEqual(types.Absfrac(dv-dw), 1./(2.*float64(i))+1e-40)
	}
}

// converts mu/Msize to a Torus32 for mu in [0,Msize[
func TestModSwitchToTorus32(t *testing.T) {
	assert := assert.New(t)
	for i := int32(2); i < 200; i++ {
		j := types.UniformTorus32Dist() % i
		v := types.ModSwitchToTorus32(j, i)
		dv := types.TorusToDouble(v)
		//printf("%d, %d, %lf, %lf\n", j, i, dv, float64(j)/i);
		assert.LessOrEqual(types.Absfrac(dv-float64(j)/float64(i)), 1./(2.*float64(i))+1e-40)
	}
}
