package core

import (
	"sync/atomic"
	"time"

	"golang.org/x/exp/rand"
)

// Uniform is a random number generator that generates draws from a uniform
// distribution.
type Uniform struct {
	min int32
	max int32
}

// NewUniform constructs a new Uniform generator with the given
// parameters. Returns an error if the parameters are outside the accepted
// range.
func NewUniform(min, max int32) *Uniform {
	return &Uniform{min: min, max: max}
}

// IncMax increments max.
func (g *Uniform) IncMax(delta int32) {
	atomic.AddInt32(&g.max, delta)
}

// Uint64 returns a random Uint64 between min and max, drawn from a uniform
// distribution.
func (g *Uniform) Int32() int32 {
	rng := rand.New(rand.NewSource(uint64(time.Now().Nanosecond())))
	max := atomic.LoadInt32(&g.max)
	return rng.Int31n(max-g.min+1) + g.min // rng.Uint64n(max-g.min+1) + g.min
}
