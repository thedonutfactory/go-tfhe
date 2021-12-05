package tfhe

import "fmt"

type TGswParams struct {
	L          int32       ///< decomp length
	Bgbit      int32       ///< log_2(Bg)
	Bg         int32       ///< decomposition base (must be a power of 2)
	HalfBg     int32       ///< Bg/2
	MaskMod    uint32      ///< Bg-1
	TlweParams *TLweParams ///< Params of each row
	Kpl        int32       ///< number of rows = (k+1)*l
	H          []Torus32   ///< powers of Bgbit
	Offset     uint32      ///< offset = Bg/2 * (2^(32-Bgbit) + 2^(32-2*Bgbit) + ... + 2^(32-l*Bgbit))
}

type TGswKey struct {
	Params     *TGswParams     ///< the parameters
	TlweParams *TLweParams     ///< the tlwe params of each rows
	Key        []IntPolynomial ///< the key (array of k polynomials)
	TlweKey    TLweKey
}

type TGswSample struct {
	AllSample  []TLweSample    ///< TLweSample* allSample; (k+1)l TLwe Sample
	BlocSample [][]*TLweSample ///< optional access to the different size l blocks
	// double currentVariance
	K int32
	L int32
}

/*

func TGswClear(result *TGswSample, params *TGswParams) {
	kpl := params.kpl
	for p := int(0); p < kpl; p++ {
		tLweClear(&result.AllSample[p], params.tlweParams)
	}
}

// Result += H
func TGswAddH(result *TGswSample, params *TGswParams) {
	k := params.tlweParams.k
	l := params.l
	h := params.H
	// compute result += H
	for bloc := int32(0); bloc <= k; bloc++ {
		for i := int32(0); i < l; i++ {
			result.BlocSample[bloc][i].A[bloc].Coefs[0] += h[i]
		}
	}
}
*/

func NewTGswSampleArray(size int, params *TGswParams) (arr []*TGswSample) {
	arr = make([]*TGswSample, size)
	for i := int(0); i < size; i++ {
		arr[i] = NewTGswSample(params)
	}
	return
}

func NewTGswSample(params *TGswParams) *TGswSample {
	k := params.TlweParams.K
	l := params.L
	kpl := params.Kpl
	allSamples := make([]TLweSample, kpl)
	for i := range allSamples {
		allSamples[i] = *NewTLweSample(params.TlweParams)
	}
	var c int = 0
	blocSamples := make([][]*TLweSample, k+1)
	for i := range blocSamples {
		blocSamples[i] = make([]*TLweSample, l)
		for j := range blocSamples[i] {
			blocSamples[i][j] = &allSamples[c]
			c++
		}
	}
	return &TGswSample{
		AllSample:  allSamples,
		BlocSample: blocSamples,
		K:          k,
		L:          l,
	}
}

func NewTGswParams(l, Bgbit int, tlweParams *TLweParams) *TGswParams {
	var Bg int = 1 << Bgbit
	var halfBg int = Bg / 2
	h := make([]Torus, l)
	for i := int(0); i < l; i++ {
		kk := (32 - (i+1)*Bgbit)
		h[i] = 1 << kk // 1/(Bg^(i+1)) as a Torus
	}

	// offset = Bg/2 * (2^(32-Bgbit) + 2^(32-2*Bgbit) + ... + 2^(32-l*Bgbit))
	var temp1 int = 0
	for i := int(0); i < l; i++ {
		temp0 := int(1 << (32 - (i+1)*Bgbit))
		temp1 += temp0
	}
	offset := temp1 * halfBg

	return &TGswParams{
		Bg:         Bg,
		L:          l,
		Bgbit:      Bgbit,
		HalfBg:     halfBg,
		MaskMod:    uint32(Bg - 1),
		TlweParams: tlweParams,
		Kpl:        int((tlweParams.K + 1) * l),
		H:          h,
		Offset:     uint32(offset),
	}
}

func NewTGswKey(params *TGswParams) *TGswKey {
	tlweKey := *NewTLweKey(params.TlweParams)
	return &TGswKey{
		Params:     params,
		TlweParams: params.TlweParams,
		TlweKey:    tlweKey,
		Key:        tlweKey.Key,
	}
}

func (s *TGswSample) DebugTGswSample(params *TGswParams) {
	tabs(1, "TGswSample {")
	kpl := params.Kpl
	k := params.TlweParams.K
	l := params.L

	tabs(2, fmt.Sprintf("k: %d", k))
	tabs(2, fmt.Sprintf("l: %d", l))
	tabs(2, fmt.Sprintf("kpl: %d", kpl))

	tabs(2, "AllSample {")
	for i := int(0); i < kpl; i++ {
		s.AllSample[i].DebugTLweSample()
	}
	tabs(2, "}")

	for bloc := int(0); bloc <= k; bloc++ {
		tabs(2, "BlockSample {")
		for i := int(0); i < l; i++ {
			s.BlocSample[bloc][i].DebugTLweSample()
		}
		tabs(2, "}")
	}
	tabs(1, "}")
}

// TGsw
/** generate a tgsw key (in fact, a tlwe key) */
func TGswKeyGen(result *TGswKey) {
	TLweKeyGen(&result.TlweKey)
}

// support Functions for TGsw
// Result = 0
func TGswClear(result *TGswSample, params *TGswParams) {
	kpl := params.Kpl
	for p := int(0); p < kpl; p++ {
		TLweClear(&result.AllSample[p], params.TlweParams)
	}
}

// Result += H
func TGswAddH(result *TGswSample, params *TGswParams) {
	k := params.TlweParams.K
	l := params.L
	h := params.H
	// compute result += H
	for bloc := int32(0); bloc <= k; bloc++ {
		for i := int32(0); i < l; i++ {
			result.BlocSample[bloc][i].A[bloc].Coefs[0] += h[i]
		}
	}
}

// Result += mu*H
func TGswAddMuH(result *TGswSample, message *IntPolynomial, params *TGswParams) {
	k := params.TlweParams.K
	N := params.TlweParams.N
	l := params.L
	h := params.H
	mu := message.Coefs

	// compute result += H
	for bloc := int32(0); bloc <= k; bloc++ {
		for i := int32(0); i < l; i++ {
			target := result.BlocSample[bloc][i].A[bloc].Coefs
			hi := h[i]
			for j := int(0); j < N; j++ {
				target[j] += mu[j] * hi
			}
		}
	}
}

// Result += mu*H, mu integer
func TGswAddMuIntH(result *TGswSample, message int, params *TGswParams) {
	k := params.TlweParams.K
	l := params.L
	h := params.H

	// compute result += H
	for bloc := int32(0); bloc <= k; bloc++ {
		for i := int32(0); i < l; i++ {
			result.BlocSample[bloc][i].A[bloc].Coefs[0] += message * h[i]
		}
	}
}

// Result = tGsw(0)
func TGswEncryptZero(result *TGswSample, alpha double, key *TGswKey) {
	rlkey := &key.TlweKey
	kpl := key.Params.Kpl
	for p := int32(0); p < kpl; p++ {
		tLweSymEncryptZero(&result.AllSample[p], alpha, rlkey)
	}
}

//mult externe de X^{a_i} par bki
func TGswMulByXaiMinusOne(result *TGswSample, ai int, bk *TGswSample, params *TGswParams) {
	par := params.TlweParams
	kpl := params.Kpl
	for i := int(0); i < kpl; i++ {
		TLweMulByXaiMinusOne(&result.AllSample[i], ai, &bk.AllSample[i], par)
	}
}

//Update l'accumulateur ligne 5 de l'algo toujours
//void tGswTLweDecompH(IntPolynomial* result, const TLweSample* sample,const TGswParams* params)
//accum *= sample
func TGswExternMulToTLwe(accum *TLweSample, sample *TGswSample, params *TGswParams) {
	par := params.TlweParams
	N := par.N
	kpl := int(params.Kpl)
	//TODO: improve this new/delete
	dec := NewIntPolynomialArray(kpl, N)

	TGswTLweDecompH(dec, accum, params)
	TLweClear(accum, par)
	for i := 0; i < kpl; i++ {
		TLweAddMulRTo(accum, &dec[i], &sample.AllSample[i], par)
	}
}

/**
 * encrypts a poly message
 */
func TGswSymEncrypt(result *TGswSample, message *IntPolynomial, alpha double, key *TGswKey) {
	TGswEncryptZero(result, alpha, key)
	TGswAddMuH(result, message, key.Params)
}

/**
 * encrypts a constant message
 */
func TGswSymEncryptInt(result *TGswSample, message int, alpha double, key *TGswKey) {
	TGswEncryptZero(result, alpha, key)
	TGswAddMuIntH(result, message, key.Params)
}

/**
 * encrypts a message = 0 ou 1
 */
func TGswEncryptB(result *TGswSample, message int, alpha double, key *TGswKey) {
	TGswEncryptZero(result, alpha, key)
	if message == 1 {
		TGswAddH(result, key.Params)
	}
}

// à revoir
func TGswSymDecrypt(result *IntPolynomial, sample *TGswSample, key *TGswKey, Msize int) {
	params := key.Params
	rlweParams := params.TlweParams
	N := rlweParams.N
	l := params.L
	k := rlweParams.K
	testvec := NewTorusPolynomial(N)
	tmp := NewTorusPolynomial(N)
	decomp := NewIntPolynomialArray(int(l), N)

	indic := ModSwitchToTorus(1, int(Msize))
	torusPolynomialClear(testvec)
	testvec.Coefs[0] = indic
	TGswTorus32PolynomialDecompH(decomp, testvec, params)

	torusPolynomialClear(testvec)
	for i := int(0); i < l; i++ {
		for j := int(1); j < N; j++ {
			Assert(decomp[i].Coefs[j] == 0)
		}
		TLwePhase(tmp, sample.BlocSample[k][i], &key.TlweKey)
		TorusPolynomialAddMulR(testvec, &decomp[i], tmp)
	}
	for i := int32(0); i < N; i++ {
		result.Coefs[i] = ModSwitchFromTorus32(testvec.Coefs[i], int32(Msize))
	}
}

//fonction de decomposition
func TGswTLweDecompH(result []IntPolynomial, sample *TLweSample, params *TGswParams) {
	k := params.TlweParams.K
	l := params.L
	var j = 0
	for i := int(0); i <= k*l; i += l {
		TGswTorusPolynomialDecompH(result[i:i+l], &sample.A[j], params)
		j++
	}

}

func TorusPolynomialDecompHOld(result []IntPolynomial, sample *TorusPolynomial, params *TGswParams) {
	N := params.TlweParams.N
	l := params.L
	Bgbit := params.Bgbit
	maskMod := params.MaskMod
	halfBg := params.HalfBg
	offset := params.Offset

	for j := int32(0); j < N; j++ {
		temp0 := uint32(sample.Coefs[j]) + offset
		for p := int32(0); p < l; p++ {
			temp1 := (temp0 >> (32 - (p+1)*Bgbit)) & maskMod // doute
			result[p].Coefs[j] = int(temp1) - halfBg
		}
	}
}

func TGswTorusPolynomialDecompH(result []IntPolynomial, sample *TorusPolynomial, params *TGswParams) {
	N := params.TlweParams.N
	l := params.L
	Bgbit := params.Bgbit
	buf := []uint32{}
	for _, vNum := range sample.Coefs {
		buf = append(buf, uint32(vNum))
	}
	maskMod := params.MaskMod
	halfBg := params.HalfBg
	offset := params.Offset
	//First, add offset to everyone
	for j := int(0); j < N; j++ {
		buf[j] += offset
	}

	//then, do the decomposition (in parallel)
	for p := int(0); p < l; p++ {
		var decal int = 32 - (p+1)*Bgbit
		for j := int(0); j < N; j++ {
			var temp1 int = int((buf[j] >> uint(decal)) & maskMod)
			result[p].Coefs[j] = temp1 - halfBg
		}
	}
	//finally, remove offset from everyone
	for j := int(0); j < N; j++ {
		buf[j] -= offset
	}
}

//result = a*b
func TGswExternProduct(result *TLweSample, a *TGswSample, b *TLweSample, params *TGswParams) {
	parlwe := params.TlweParams
	N := parlwe.N
	kpl := params.Kpl
	dec := NewIntPolynomialArray(int(kpl), N)
	TGswTLweDecompH(dec, b, params)
	TLweClear(result, parlwe)
	for i := int(0); i < kpl; i++ {
		TLweAddMulRTo(result, &dec[i], &a.AllSample[i], parlwe)
	}
	result.CurrentVariance += b.CurrentVariance //todo + the error term?
}

/**
 * result = (0,mu)
 */
func TGswNoiselessTrivial(result *TGswSample, mu *IntPolynomial, params *TGswParams) {
	TGswClear(result, params)
	TGswAddMuH(result, mu, params)
}
