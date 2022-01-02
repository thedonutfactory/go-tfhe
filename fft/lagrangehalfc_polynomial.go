package fft

import "github.com/thedonutfactory/go-tfhe/types"

type LagrangeHalfCPolynomial struct {
	Coefs []complex128
}

func NewLagrangeHalfCPolynomial(n int32) *LagrangeHalfCPolynomial {
	return &LagrangeHalfCPolynomial{Coefs: make([]complex128, n/2)}
}

func NewLagrangeHalfCPolynomialArray(size, n int32) []*LagrangeHalfCPolynomial {
	arr := make([]*LagrangeHalfCPolynomial, size)
	for i := int32(0); i < size; i++ {
		arr[i] = NewLagrangeHalfCPolynomial(n)
	}
	return arr
}

func (p *LagrangeHalfCPolynomial) Clear() {
	p.Coefs = make([]complex128, len(p.Coefs))
}

func (p *LagrangeHalfCPolynomial) SetTorusConstant(mu types.Torus32) {
	muc := complex(types.TorusToDouble(mu), 0.)
	for j := 0; j < len(p.Coefs); j++ {
		p.Coefs[j] = muc
	}
}

func (p *LagrangeHalfCPolynomial) AddTorusConstant(mu types.Torus32) {
	muc := complex(types.TorusToDouble(mu), 0.)
	for j := 0; j < len(p.Coefs); j++ {
		p.Coefs[j] += muc
	}
}

/** termwise multiplication in Lagrange space */
func (p *LagrangeHalfCPolynomial) Mul(a *LagrangeHalfCPolynomial, b *LagrangeHalfCPolynomial) {
	for j := 0; j < len(p.Coefs); j++ {
		p.Coefs[j] = a.Coefs[j] * b.Coefs[j]
	}
}

/** termwise multiplication and addTo in Lagrange space */
func (p *LagrangeHalfCPolynomial) AddMul(a *LagrangeHalfCPolynomial, b *LagrangeHalfCPolynomial) {
	for j := 0; j < len(p.Coefs); j++ {
		p.Coefs[j] += a.Coefs[j] * b.Coefs[j]
	}
}

/** termwise multiplication and addTo in Lagrange space */
func (p *LagrangeHalfCPolynomial) SubMul(a *LagrangeHalfCPolynomial, b *LagrangeHalfCPolynomial) {
	for j := 0; j < len(p.Coefs); j++ {
		p.Coefs[j] += a.Coefs[j] * b.Coefs[j]
	}
}

func (p *LagrangeHalfCPolynomial) AddTo(a *LagrangeHalfCPolynomial) {
	for j := 0; j < len(p.Coefs); j++ {
		p.Coefs[j] += a.Coefs[j]
	}
}
