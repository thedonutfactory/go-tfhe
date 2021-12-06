package tfhe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaussian32(t *testing.T) {
	assert := assert.New(t)

	const MESSAGE1 int64 = 123456789
	const MESSAGE2 int64 = 987654321

	reps1 := gaussian32(MESSAGE1, 0)
	reps2 := gaussian32(MESSAGE2, 0)
	assert.Equal(MESSAGE1, reps1)
	assert.Equal(MESSAGE2, reps2)

	reps1 = gaussian32(MESSAGE1, 0.01)
	reps2 = gaussian32(MESSAGE2, 0.5)
	assert.NotEqual(MESSAGE1, reps1)
	assert.NotEqual(MESSAGE2, reps2)
	//assert.LessOrEqual(Abs(MESSAGE1-reps1), int64(80000000))
}

func TestConversion(t *testing.T) {
	assert := assert.New(t)
	// conversion from double to Torus
	// conversion from Torus to double
	assert.EqualValues(uint64(0), Dtot32(0))
	assert.EqualValues(uint64(1)<<63, Dtot32(0.5))
	assert.EqualValues(uint64(1)<<63, Dtot32(-0.5))
	assert.EqualValues(uint64(1)<<62, Dtot32(0.25))
	assert.EqualValues(uint64(0xC000000000000000), Dtot32(-0.25))
}

//  Used to approximate the phase to the nearest multiple of  1/Msize
func TestApproxPhase(t *testing.T) {
	assert := assert.New(t)
	for i := 2; i < 200; i++ {
		v := UniformTorusDist()
		w := approxPhase(v, i)
		dv := T32tod(v)
		dw := T32tod(w)
		fmt.Printf("%d, %f, %f, %f \n", i, dv, dw, float64(i)*dw)
		assert.LessOrEqual(absfrac(dv-dw), 1./(2.*float64(i))+1e-40)
		assert.LessOrEqual(absfrac(float64(i)*dw), float64(i)*1e-9)
	}
}

func TestModSwitchFromTorus(t *testing.T) {
	assert := assert.New(t)

	for i := 2; i < 200; i++ {
		v := UniformTorusDist()
		w := ModSwitchFromTorus(v, i)
		dv := T32tod(v)
		dw := double(w) / double(i)
		assert.LessOrEqual(absfrac(dv-dw), 1./(2.*float64(i))+1e-40)
	}
}

// converts mu/Msize to a Torus for mu in [0,Msize[
func TestModSwitchToTorus(t *testing.T) {
	assert := assert.New(t)
	for i := 2; i < 200; i++ {
		j := UniformTorusDist() % int64(i)
		v := ModSwitchToTorus(j, i)
		dv := T32tod(v)
		//printf("%d, %d, %lf, %lf\n", j, i, dv, double(j)/i);
		assert.LessOrEqual(absfrac(dv-double(j)/double(i)), 1./(2.*float64(i))+1e-40)
	}
}
