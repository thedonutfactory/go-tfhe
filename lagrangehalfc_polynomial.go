package tfhe

//"github.com/takatoh/fft"

type LagrangeHalfCPolynomial struct {
	Coefs []complex128
}

func NewLagrangeHalfCPolynomial(n int) *LagrangeHalfCPolynomial {
	//Assert(n == 1024)
	return &LagrangeHalfCPolynomial{Coefs: make([]complex128, n/2)}
}

func NewLagrangeHalfCPolynomialArray(size, n int32) []*LagrangeHalfCPolynomial {
	arr := make([]*LagrangeHalfCPolynomial, size)
	for i := int32(0); i < size; i++ {
		arr[i] = NewLagrangeHalfCPolynomial(n)
	}
	return arr
}

func LagrangeHalfCPolynomialClear(p *LagrangeHalfCPolynomial) {
	p.Coefs = make([]complex128, len(p.Coefs))
}

func LagrangeHalfCPolynomialSetTorusConstant(result *LagrangeHalfCPolynomial, mu Torus32) {
	muc := complex(TorusToDouble(mu), 0.)
	for j := 0; j < len(result.Coefs); j++ {
		result.Coefs[j] = muc
	}
}

func LagrangeHalfCPolynomialAddTorusConstant(result *LagrangeHalfCPolynomial, mu Torus32) {
	muc := complex(TorusToDouble(mu), 0.)
	for j := 0; j < len(result.Coefs); j++ {
		if j < 10 {
			//fmt.Printf("Before: (%f, %f)\n", real(result.coefsC[j]), imag(result.coefsC[j]))
			//fmt.Printf("Add Mu: (%f, %f)\n", real(muc), imag(muc))
			result.Coefs[j] += muc
			//fmt.Printf("Result: (%f, %f)\n", real(result.coefsC[j]), imag(result.coefsC[j]))
			//fmt.Println()
		} else {
			result.Coefs[j] += muc
		}
	}
}

/** termwise multiplication in Lagrange space */
func LagrangeHalfCPolynomialMul(result *LagrangeHalfCPolynomial, a *LagrangeHalfCPolynomial, b *LagrangeHalfCPolynomial) {
	for j := 0; j < len(result.Coefs); j++ {
		result.Coefs[j] = a.Coefs[j] * b.Coefs[j]
	}
}

/** termwise multiplication and addTo in Lagrange space */
func LagrangeHalfCPolynomialAddMul(accum *LagrangeHalfCPolynomial, a *LagrangeHalfCPolynomial, b *LagrangeHalfCPolynomial) {
	for j := 0; j < len(accum.Coefs); j++ {
		accum.Coefs[j] += a.Coefs[j] * b.Coefs[j]
	}
}

/** termwise multiplication and addTo in Lagrange space */
func LagrangeHalfCPolynomialSubMul(accum *LagrangeHalfCPolynomial, a *LagrangeHalfCPolynomial, b *LagrangeHalfCPolynomial) {
	for j := 0; j < len(accum.Coefs); j++ {
		accum.Coefs[j] += a.Coefs[j] * b.Coefs[j]
	}
}

func LagrangeHalfCPolynomialAddTo(accum *LagrangeHalfCPolynomial, a *LagrangeHalfCPolynomial) {
	for j := 0; j < len(accum.Coefs); j++ {
		accum.Coefs[j] += a.Coefs[j]
	}
}
