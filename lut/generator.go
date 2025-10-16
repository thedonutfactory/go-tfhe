package lut

import (
	"math"

	"github.com/thedonutfactory/go-tfhe/params"
)

// Generator creates lookup tables from functions for programmable bootstrapping
type Generator struct {
	Encoder         *Encoder
	PolyDegree      int
	LookUpTableSize int // For binary: equals PolyDegree (not 2*PolyDegree!)
}

// NewGenerator creates a new LUT generator
func NewGenerator(messageModulus int) *Generator {
	polyDegree := params.GetTRGSWLv1().N
	// CRITICAL: For standard TFHE, lookUpTableSize = polyDegree (polyExtendFactor = 1)
	// Only for extended configurations is lookUpTableSize > polyDegree
	lookUpTableSize := polyDegree

	return &Generator{
		Encoder:         NewEncoder(messageModulus),
		PolyDegree:      polyDegree,
		LookUpTableSize: lookUpTableSize,
	}
}

// NewGeneratorWithScale creates a new LUT generator with custom scale
func NewGeneratorWithScale(messageModulus int, scale float64) *Generator {
	polyDegree := params.GetTRGSWLv1().N
	return &Generator{
		Encoder:         NewEncoderWithScale(messageModulus, scale),
		PolyDegree:      polyDegree,
		LookUpTableSize: polyDegree, // Standard: lookUpTableSize = polyDegree
	}
}

// GenLookUpTable generates a lookup table from a function f: int -> int
func (g *Generator) GenLookUpTable(f func(int) int) *LookUpTable {
	lut := NewLookUpTable()
	g.GenLookUpTableAssign(f, lut)
	return lut
}

// GenLookUpTableAssign generates a lookup table and writes to lutOut
//
// Algorithm from tfhe-go reference implementation (bootstrap_lut.go:111-132)
// For standard TFHE with polyExtendFactor=1 (lookUpTableSize = polyDegree):
// 1. Create lutRaw[lookUpTableSize]
// 2. For each message x, fill range with encoded f(x)
// 3. Rotate by offset
// 4. Negate tail
// 5. Store in polynomial
func (g *Generator) GenLookUpTableAssign(f func(int) int, lutOut *LookUpTable) {
	messageModulus := g.Encoder.MessageModulus

	// Create raw LUT buffer (size = lookUpTableSize, which equals N for standard TFHE)
	lutRaw := make([]params.Torus, g.LookUpTableSize)

	// Fill each message's range with encoded output
	for x := 0; x < messageModulus; x++ {
		start := divRound(x*g.LookUpTableSize, messageModulus)
		end := divRound((x+1)*g.LookUpTableSize, messageModulus)

		// Apply function to message index
		y := f(x)

		// Encode the output: message * scale
		encodedY := g.Encoder.Encode(y)

		// Fill range
		for xx := start; xx < end; xx++ {
			lutRaw[xx] = encodedY
		}
	}

	// Rotate by offset
	offset := divRound(g.LookUpTableSize, 2*messageModulus)

	// Apply rotation
	rotated := make([]params.Torus, g.LookUpTableSize)
	for i := 0; i < g.LookUpTableSize; i++ {
		srcIdx := (i + offset) % g.LookUpTableSize
		rotated[i] = lutRaw[srcIdx]
	}

	// Negate tail portion
	for i := g.LookUpTableSize - offset; i < g.LookUpTableSize; i++ {
		rotated[i] = -rotated[i]
	}

	// Store in polynomial
	// For polyExtendFactor=1: just copy all lookUpTableSize coefficients
	for i := 0; i < g.LookUpTableSize; i++ {
		lutOut.Poly.B[i] = rotated[i]
		lutOut.Poly.A[i] = 0
	}
}

// GenLookUpTableFull generates a lookup table from a function f: int -> Torus
func (g *Generator) GenLookUpTableFull(f func(int) params.Torus) *LookUpTable {
	lut := NewLookUpTable()
	g.GenLookUpTableFullAssign(f, lut)
	return lut
}

// GenLookUpTableFullAssign generates a lookup table with full control
func (g *Generator) GenLookUpTableFullAssign(f func(int) params.Torus, lutOut *LookUpTable) {
	messageModulus := g.Encoder.MessageModulus

	lutRaw := make([]params.Torus, g.LookUpTableSize)

	for x := 0; x < messageModulus; x++ {
		start := divRound(x*g.LookUpTableSize, messageModulus)
		end := divRound((x+1)*g.LookUpTableSize, messageModulus)

		y := f(x)

		for i := start; i < end; i++ {
			lutRaw[i] = y
		}
	}

	offset := divRound(g.LookUpTableSize, 2*messageModulus)
	rotated := make([]params.Torus, g.LookUpTableSize)
	for i := 0; i < g.LookUpTableSize; i++ {
		srcIdx := (i + offset) % g.LookUpTableSize
		rotated[i] = lutRaw[srcIdx]
	}

	for i := g.LookUpTableSize - offset; i < g.LookUpTableSize; i++ {
		rotated[i] = -rotated[i]
	}

	for i := 0; i < g.LookUpTableSize; i++ {
		lutOut.Poly.B[i] = rotated[i]
		lutOut.Poly.A[i] = 0
	}
}

// GenLookUpTableCustom generates a lookup table with custom message modulus and scale
func (g *Generator) GenLookUpTableCustom(f func(int) int, messageModulus int, scale float64) *LookUpTable {
	lut := NewLookUpTable()

	oldEncoder := g.Encoder
	g.Encoder = NewEncoderWithScale(messageModulus, scale)

	g.GenLookUpTableAssign(f, lut)

	g.Encoder = oldEncoder

	return lut
}

// ModSwitch switches the modulus of x from Torus (2^32) to lookUpTableSize
// For standard TFHE with lookUpTableSize=N: result in [0, N)
func (g *Generator) ModSwitch(x params.Torus) int {
	scaled := float64(x) / float64(uint64(1)<<32) * float64(g.LookUpTableSize)
	result := int(math.Round(scaled)) % g.LookUpTableSize

	if result < 0 {
		result += g.LookUpTableSize
	}

	return result
}

// divRound performs integer division with rounding
func divRound(a, b int) int {
	return (a + b/2) / b
}
