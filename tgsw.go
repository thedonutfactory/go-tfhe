package tfhe

import "fmt"

type TGswParams struct {
	l          int64       ///< decomp length
	Bgbit      int64       ///< log_2(Bg)
	Bg         int64       ///< decomposition base (must be a power of 2)
	halfBg     int64       ///< Bg/2
	maskMod    uint64      ///< Bg-1
	TlweParams *TLweParams ///< Params of each row
	Kpl        int64       ///< number of rows = (k+1)*l
	H          []Torus     ///< powers of Bgbit
	offset     uint64      ///< offset = Bg/2 * (2^(32-Bgbit) + 2^(32-2*Bgbit) + ... + 2^(32-l*Bgbit))
}

type TGswKey struct {
	params     *TGswParams     ///< the parameters
	tlweParams *TLweParams     ///< the tlwe params of each rows
	key        []IntPolynomial ///< the key (array of k polynomials)
	TlweKey    TLweKey
}

type TGswSample struct {
	AllSample  []TLweSample    ///< TLweSample* allSample; (k+1)l TLwe Sample
	BlocSample [][]*TLweSample ///< optional access to the different size l blocks
	// double currentVariance
	k int
	l int
}

/*

func TGswClear(result *TGswSample, params *TGswParams) {
	kpl := params.kpl
	for p := int64(0); p < kpl; p++ {
		tLweClear(&result.AllSample[p], params.tlweParams)
	}
}

// Result += H
func TGswAddH(result *TGswSample, params *TGswParams) {
	k := params.tlweParams.k
	l := params.l
	h := params.H
	// compute result += H
	for bloc := int64(0); bloc <= k; bloc++ {
		for i := int64(0); i < l; i++ {
			result.BlocSample[bloc][i].A[bloc].CoefsT[0] += h[i]
		}
	}
}
*/

func NewTGswSampleArray(size int, params *TGswParams) (arr []*TGswSample) {
	arr = make([]*TGswSample, size)
	for i := 0; i < size; i++ {
		arr[i] = NewTGswSample(params)
	}
	return
}

func NewTGswSample(params *TGswParams) *TGswSample {
	k := params.TlweParams.K
	l := params.l
	kpl := params.Kpl
	//h := params.H

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
	/*
		blocSamples := make([][]TLweSample, k+1)
		for i := range blocSamples {
			blocSamples[i] = make([]TLweSample, l)
			for j := range blocSamples[i] {
				blocSamples[i][j] = *NewTLweSample(params.tlweParams)
			}
		}
	*/

	return &TGswSample{
		AllSample:  allSamples,
		BlocSample: blocSamples,
		k:          k,
		l:          int(l),
	}
}

func NewTGswParams(l, Bgbit int64, tlweParams *TLweParams) *TGswParams {
	var Bg int64 = 1 << Bgbit
	var halfBg int64 = Bg / 2
	h := make([]Torus, l)
	for i := int64(0); i < l; i++ {
		kk := (64 - (i+1)*Bgbit)
		h[i] = 1 << kk // 1/(Bg^(i+1)) as a Torus
	}

	// offset = Bg/2 * (2^(32-Bgbit) + 2^(32-2*Bgbit) + ... + 2^(32-l*Bgbit))
	var temp1 int64 = 0
	for i := int64(0); i < l; i++ {
		temp0 := int64(1 << (64 - (i+1)*Bgbit))
		temp1 += temp0
	}
	offset := temp1 * halfBg

	return &TGswParams{
		Bg:         Bg,
		l:          l,
		Bgbit:      Bgbit,
		halfBg:     halfBg,
		maskMod:    uint64(Bg - 1),
		TlweParams: tlweParams,
		Kpl:        int64(tlweParams.K+1) * l,
		H:          h,
		offset:     uint64(offset),
	}
}

func NewTGswKey(params *TGswParams) *TGswKey {
	tlweKey := *NewTLweKey(params.TlweParams)
	return &TGswKey{
		params:     params,
		tlweParams: params.TlweParams,
		TlweKey:    tlweKey,
		key:        tlweKey.key,
	}
}

func (s *TGswSample) DebugTGswSample(params *TGswParams) {
	tabs(1, "TGswSample {")
	kpl := params.Kpl
	k := params.TlweParams.K
	l := params.l

	tabs(2, fmt.Sprintf("k: %d", k))
	tabs(2, fmt.Sprintf("l: %d", l))
	tabs(2, fmt.Sprintf("kpl: %d", kpl))

	tabs(2, "AllSample {")
	for i := int64(0); i < kpl; i++ {
		s.AllSample[i].DebugTLweSample()
	}
	tabs(2, "}")

	for bloc := 0; bloc <= k; bloc++ {
		tabs(2, "BlockSample {")
		for i := int64(0); i < l; i++ {
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
	for p := int64(0); p < kpl; p++ {
		TLweClear(&result.AllSample[p], params.TlweParams)
	}
}

// Result += H
func TGswAddH(result *TGswSample, params *TGswParams) {
	k := params.TlweParams.K
	l := params.l
	h := params.H
	// compute result += H
	for bloc := 0; bloc <= k; bloc++ {
		for i := int64(0); i < l; i++ {
			result.BlocSample[bloc][i].A[bloc].CoefsT[0] += h[i]
		}
	}
}

// Result += mu*H
func TGswAddMuH(result *TGswSample, message *IntPolynomial, params *TGswParams) {
	k := params.TlweParams.K
	N := params.TlweParams.N
	l := params.l
	h := params.H
	mu := message.Coefs

	// compute result += H
	for bloc := 0; bloc <= k; bloc++ {
		for i := int64(0); i < l; i++ {
			target := result.BlocSample[bloc][i].A[bloc].CoefsT
			hi := h[i]
			for j := 0; j < N; j++ {
				target[j] += mu[j] * hi
			}
		}
	}
}

// Result += mu*H, mu integer
func TGswAddMuIntH(result *TGswSample, message int64, params *TGswParams) {
	k := params.TlweParams.K
	l := params.l
	h := params.H

	// compute result += H
	for bloc := 0; bloc <= k; bloc++ {
		for i := 0; i < int(l); i++ {
			result.BlocSample[bloc][i].A[bloc].CoefsT[0] += message * h[i]
		}
	}
}

// Result = tGsw(0)
func TGswEncryptZero(result *TGswSample, alpha double, key *TGswKey) {
	rlkey := &key.TlweKey
	kpl := key.params.Kpl
	for p := int64(0); p < kpl; p++ {
		tLweSymEncryptZero(&result.AllSample[p], alpha, rlkey)
	}
}

//mult externe de X^{a_i} par bki
func TGswMulByXaiMinusOne(result *TGswSample, ai int64, bk *TGswSample, params *TGswParams) {
	par := params.TlweParams
	kpl := params.Kpl
	for i := int64(0); i < kpl; i++ {
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
	TGswAddMuH(result, message, key.params)
}

/**
 * encrypts a constant message
 */
func TGswSymEncryptInt(result *TGswSample, message int64, alpha double, key *TGswKey) {
	TGswEncryptZero(result, alpha, key)
	TGswAddMuIntH(result, message, key.params)
}

/**
 * encrypts a message = 0 ou 1
 */
func TGswEncryptB(result *TGswSample, message int64, alpha double, key *TGswKey) {
	TGswEncryptZero(result, alpha, key)
	if message == 1 {
		TGswAddH(result, key.params)
	}
}

// à revoir
func TGswSymDecrypt(result *IntPolynomial, sample *TGswSample, key *TGswKey, Msize int) {
	params := key.params
	rlweParams := params.TlweParams
	N := rlweParams.N
	l := params.l
	k := rlweParams.K
	testvec := NewTorusPolynomial(N)
	tmp := NewTorusPolynomial(N)
	decomp := NewIntPolynomialArray(int(l), N)

	indic := ModSwitchToTorus(1, Msize)
	torusPolynomialClear(testvec)
	testvec.CoefsT[0] = indic
	TGswTorusPolynomialDecompH(decomp, testvec, params)

	torusPolynomialClear(testvec)
	for i := int64(0); i < l; i++ {
		for j := 1; j < N; j++ {
			Assert(decomp[i].Coefs[j] == 0)
		}
		TLwePhase(tmp, sample.BlocSample[k][i], &key.TlweKey)
		TorusPolynomialAddMulR(testvec, &decomp[i], tmp)
	}
	for i := 0; i < N; i++ {
		result.Coefs[i] = ModSwitchFromTorus(testvec.CoefsT[i], Msize)
	}
}

//fonction de decomposition
func TGswTLweDecompH(result []IntPolynomial, sample *TLweSample, params *TGswParams) {
	k := params.TlweParams.K
	l := int(params.l)
	var j = 0
	for i := 0; i <= k*l; i += l {
		TGswTorusPolynomialDecompH(result[i:i+l], &sample.A[j], params)
		j++
	}

}

/*
func TorusPolynomialDecompHOld(result []IntPolynomial, sample *TorusPolynomial, params *TGswParams) {
	N := params.TlweParams.N
	l := params.l
	Bgbit := params.Bgbit
	maskMod := params.maskMod
	halfBg := params.halfBg
	offset := params.offset

	for j := 0; j < N; j++ {
		temp0 := uint64(sample.CoefsT[j]) + offset
		for p := int64(0); p < l; p++ {
			temp1 := (temp0 >> (64 - (p+1)*Bgbit)) & maskMod // doute
			result[p].Coefs[j] = int64(temp1) - halfBg
		}
	}
}
*/

func TGswTorusPolynomialDecompH(result []IntPolynomial, sample *TorusPolynomial, params *TGswParams) {
	N := params.TlweParams.N
	l := params.l
	Bgbit := params.Bgbit
	buf := []uint64{}
	for _, vNum := range sample.CoefsT {
		buf = append(buf, uint64(vNum))
	}
	maskMod := params.maskMod
	halfBg := params.halfBg
	offset := params.offset
	//First, add offset to everyone
	for j := 0; j < N; j++ {
		buf[j] += offset
	}

	//then, do the decomposition (in parallel)
	for p := int64(0); p < l; p++ {
		var decal int64 = 64 - (p+1)*Bgbit
		for j := 0; j < N; j++ {
			temp1 := int64((buf[j] >> uint64(decal)) & maskMod)
			result[p].Coefs[j] = temp1 - int64(halfBg)
		}
	}
	//finally, remove offset from everyone
	for j := 0; j < N; j++ {
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
	for i := int64(0); i < kpl; i++ {
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
