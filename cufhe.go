package tfhe

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

const KPtxtSpace = 2

type Param struct {
	lwe_n_                    int
	tlwe_n_                   int
	tlwe_k_                   int
	tgsw_decomp_bits_         int
	tgsw_decomp_size_         int
	keyswitching_decomp_bits_ int
	keyswitching_decomp_size_ int
	lwe_noise_                float64
	tlwe_noise_               float64
}

func NewParam() *Param {
	return &Param{
		lwe_n_:                    500,
		tlwe_n_:                   1024,
		tlwe_k_:                   1,
		tgsw_decomp_bits_:         10,
		tgsw_decomp_size_:         2,
		keyswitching_decomp_bits_: 2,
		keyswitching_decomp_size_: 8,
		lwe_noise_:                math.Pow(2.0, -15),
		tlwe_noise_:               9.e-9,
	}
}

func BuildParam(lwe_n, tlwe_n, tlwe_k, tgsw_decomp_bits,
	tgsw_decomp_size, keyswitching_decomp_bits,
	keyswitching_decomp_size int,
	lwe_noise, tlwe_noise float64) *Param {
	return &Param{
		lwe_n_:                    lwe_n,
		tlwe_n_:                   tlwe_n,
		tlwe_k_:                   tlwe_k,
		tgsw_decomp_bits_:         tgsw_decomp_bits,
		tgsw_decomp_size_:         tgsw_decomp_size,
		keyswitching_decomp_bits_: keyswitching_decomp_bits,
		keyswitching_decomp_size_: keyswitching_decomp_size,
		lwe_noise_:                lwe_noise,
		tlwe_noise_:               tlwe_noise,
	}
}

func GetDefaultParam() *Param {
	return &Param{
		500,
		1024,
		1,
		10,
		2,
		2,
		8,
		math.Pow(2.0, -15),
		9.e-9,
	}
}

/**
* Private Key.
* Necessary for encryption/decryption and public key generation.
 */
type PriKey struct {
	lwe_key_  *LWEKey
	tlwe_key_ *TLWEKey
}

func NewPriKey(isAlias bool) *PriKey {
	param := GetDefaultParam()
	return &PriKey{
		lwe_key_:  NewLWEKey(param.lwe_n_),
		tlwe_key_: NewTLWEKey(param.tlwe_n_, param.tlwe_k_),
	}
}

/**
* Public Key.
* Necessary for a server to perform homomorphic evaluation.
 */
type PubKey struct {
	bk_  *BootstrappingKey
	ksk_ *KeySwitchingKey
}

func NewPubKey(isAlias bool) *PubKey {
	param := GetDefaultParam()
	return &PubKey{
		bk_: NewBootstrappingKey(param.tlwe_n_, param.tlwe_k_,
			param.tgsw_decomp_size_,
			param.tgsw_decomp_bits_, param.lwe_n_),
		ksk_: NewKeySwitchingKey(param.lwe_n_,
			param.keyswitching_decomp_size_,
			param.keyswitching_decomp_bits_,
			param.tlwe_n_*param.tlwe_k_),
	}
}

/** Ciphertext. */
type Ctxt struct {
	lwe_sample_        *LWESample
	lwe_sample_device_ *LWESample
}

/** Plaintext is in {0, 1}. */
type Ptxt struct {
	Message    uint32
	kPtxtSpace uint32
}

func NewPtxt() *Ptxt {
	return &Ptxt{
		kPtxtSpace: 2,
	}
}

func NewPtxtArray(n int) (arr []*Ptxt) {
	arr = make([]*Ptxt, n)
	for i := 0; i < n; i++ {
		arr[i] = NewPtxt()
	}
	return
}

func (ptxt *Ptxt) opEquals(message uint32) {
	ptxt.Message = message % ptxt.kPtxtSpace
}

func (ptxt *Ptxt) set(message uint32) {
	ptxt.Message = message % ptxt.kPtxtSpace
}

func (ptxt *Ptxt) get() uint32 {
	return ptxt.Message
}

////////////////////////////////////////////////////////////////////////////////
/*
std::default_random_engine generator; // @todo Set Seed!

void RandomGeneratorSetSeed(uint32_t* values, int32_t size = 1) {
	std::seed_seq seeds(values, values + size);
	generator.seed(seeds);
}
*/

func SDFromBound(noise_bound float64) float64 {
	return noise_bound * math.Sqrt(2.0/math.Pi)
}

// Conversions go back to -0.5~0.5 if Torus is int32!
func TorusFromDouble(d float64) Torus {
	return Torus(math.Round(math.Mod(d, 1) * math.Pow(2, 32)))
}

func DoubleFromTorus(t Torus) float64 {
	return float64(t) / math.Pow(2, 32)
}

func ApproxPhase(phase Torus, msg_space int32) Torus {
	//interv := (uint64(1) << 63) / msg_space * 2
	interv := ((uint64(1) << 63) / uint64(msg_space)) * 2 // width of each interval
	half_interval := interv / 2
	phase64 := (uint64(phase) << 32) + half_interval
	//floor to the nearest multiples of interv
	phase64 -= phase64 % interv
	//rescale to torus32
	return Torus(phase64 >> 32)
}

func PolyMulAddBinary(r, a []Torus, b []Binary, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n-i; j++ {
			r[i+j] += a[i] & int32(-b[j])
		}
		for j := n - i; j < n; j++ {
			r[i+j-n] -= a[i] & int32(-b[j])
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
func LWEKeyGen(param *Param) *LWEKey {
	dist := distuv.Uniform{
		Min: 0,
		Max: 1,
	}

	key := NewLWEKey(param.lwe_n_)
	for i := range key.Key {
		// key.data()[i] = Torus32(math.Round(dist.Rand()))
		key.Key[i] = uint32(math.Round(dist.Rand()))
	}
	return key
}

/*
func TLweKeyGen(result *TLweKey) {
	N := result.params.N
	k := result.params.K
	dist := distuv.Uniform{
		Min: 0,
		Max: 1,
	}
	for i := int32(0); i < k; i++ {
		for j := int32(0); j < N; j++ {
			result.key[i].Coefs[j] = Torus32(math.Round(dist.Rand()))
		}
	}
}
*/

func TLWEKeyGen(param *Param) *TLWEKey {
	key := NewTLWEKey(param.tlwe_n_, param.tlwe_k_)
	dist := distuv.Uniform{
		Min: 0,
		Max: 1,
	}
	for i := 0; i < key.NumPolys(); i++ {
		for j := 0; j < key.N; j++ {
			key.ExtractPoly(i)[j] = uint32(math.Round(dist.Rand()))
		}
	}
	return key
}

func LWEEncrypt(ct *LWESample, pt Torus, key *LWEKey) {
	//noise_bound := GetDefaultParam().lwe_noise_
	var M int32 = 8
	noise_bound := 1.0 / (10.0 * float64(M))

	//d1 := NormalDist(0.0, SDFromBound(noise_bound))
	d1 := NormalDist(0.0, noise_bound)
	ct.B = pt + TorusFromDouble(d1.Rand())
	d2 := UniformDist(math.MinInt32, math.MaxInt32)
	for i := 0; i < key.N; i++ {
		ct.A[i] = Torus(d2.Rand())
		ct.B += ct.A[i] * int32(key.Key[i])
	}
}

func LWEEncryptExternalNoise(ct *LWESample, pt Torus, key *LWEKey, noise float64) {
	ct.B = pt + int32(noise)
	d1 := UniformDist(math.MinInt32, math.MaxInt32)
	for i := 0; i < key.N; i++ {
		ct.A[i] = Torus(d1.Rand())
		ct.B += ct.A[i] * int32(key.Key[i])
	}
}

func LwePhase(sample *LWESample, key *LWEKey) Torus {
	var axs Torus = 0
	for i := 0; i < int(key.N); i++ {
		axs += sample.A[i] * int32(key.Key[i])
	}
	return sample.B - axs
}

func LWEDecrypt(ct *LWESample, key *LWEKey, space int32) Torus {
	Assert(ct.N == key.N)
	err := ct.B

	/*
		for i := 0; i < ct.N; i++ {
			err -= ct.A[i] * int32(key.Key[i])
		}
		return err
	*/

	LwePhase(ct, key)
	return ApproxPhase(err, space)
}

func KeySwitchingKeyGen(lwe_key_to, lwe_key_from *LWEKey) *KeySwitchingKey {
	param := GetDefaultParam()
	key := NewKeySwitchingKey(param.lwe_n_,
		param.keyswitching_decomp_size_,
		param.keyswitching_decomp_bits_,
		param.tlwe_n_*param.tlwe_k_)

	var mu Torus
	var temp uint32
	// lwe_sample := NewLWESample(param.lwe_n_)

	//base := int32(1 << basebit)
	// n=1024, t=8, base=2
	//sizeks := n * t * (base - 1)

	total := key.M * key.L * (0x1 << key.W)

	//total := key.NumLWESamples()
	noise := make([]float64, total)
	var err float64 = 0
	nd := NormalDist(0.0, SDFromBound(param.lwe_noise_))
	for i := 0; i < total; i++ {
		noise[i] = nd.Rand()
		err += noise[i]
	}
	err /= float64(total)
	for i := 0; i < total; i++ {
		noise[i] -= err
	}

	var index uint32 = 0
	for i := 0; i < key.M; i++ {
		temp = lwe_key_from.Key[i]
		for j := 0; j < key.L; j++ {
			looper := (0x1 << key.W)
			/*
				for k := 0; k < looper; k++ {
					lwe_sample := key.ExtractLWESample(key.GetLWESampleIndex(i, j, k))
					mu = Torus((temp * uint32(k)) * (1 << (32 - (j+1)*key.W)))
					LWEEncryptExternalNoise(lwe_sample, mu, lwe_key_to, noise[index])
					index++
				}
			*/

			for k := 0; k < looper; k++ {
				mu = Torus((temp * uint32(k)) * (1 << (32 - (j+1)*key.W)))
				LWEEncryptExternalNoise(key.A[i][j][k], mu, lwe_key_to, noise[index])
				index++
			}
		}
	}
	return key
}

/*
func tLweSymEncryptZero(ct *TLWESample, key *TLWEKey) {
	N := key.N
	k := key.K

	for j := 0; j < N; j++ {
		result.B().CoefsT[j] = gaussian32(0, alpha)
	}

	for i := int32(0); i < k; i++ {
		torusPolynomialUniform(&result.A[i])
		TorusPolynomialAddMulR(result.B(), &key.key[i], &result.A[i])
	}

	result.CurrentVariance = alpha * alpha
}
*/

func TLWEEncryptZero(ct *TLWESample, key *TLWEKey) {
	noise_bound := GetDefaultParam().tlwe_noise_
	dist_b := NormalDist(0.0, SDFromBound(noise_bound))
	for i := 0; i < key.N; i++ {
		ct.b()[i] = TorusFromDouble(dist_b.Rand())
	}

	dist_a := UniformDist(math.MinInt32, math.MaxInt32)
	for i := 0; i < key.K; i++ {
		for j := 0; j < key.N; j++ {
			ct.a(i)[j] = TorusFromDouble(dist_a.Rand())
		}
		// PolyMulAddBinary(ct.b(), ct.a(i), key.data(), key.N)
		PolyMulAddBinary(ct.b(), ct.a(i), key.A[i], key.N)
	}
}

func TGSWEncryptBinary(ct *TGSWSample, pt Binary, key *TGSWKey) {
	//param := GetDefaultParam()
	l := ct.L
	k := ct.K
	w := ct.W
	//tlwe_sample := NewTLWESample(param.tlwe_n_, param.tlwe_k_)
	for i := 0; i < ct.NumTLWESamples(); i++ {
		tlwe_sample := ct.ExtractTLWESample(i)
		TLWEEncryptZero(tlwe_sample, key)
	}
	for i := 0; i < l; i++ {
		mu := Torus(pt << (32 - w*(i+1)))
		for j := 0; j < k; j++ {
			tlwe_sample := ct.ExtractTLWESample(j*l + i)
			tlwe_sample.a(j)[0] += mu
		}
		tlwe_sample := ct.ExtractTLWESample(k*l + i)
		tlwe_sample.b()[0] += mu
	}
}

func BootstrappingKeyGen(
	lwe_key *LWEKey,
	tgsw_key *TGSWKey) *BootstrappingKey {

	param := GetDefaultParam()
	//tgsw_sample := NewTGSWSample(0, 0, 0, 0)
	key := NewBootstrappingKey(param.tlwe_n_, param.tlwe_k_,
		param.tgsw_decomp_size_,
		param.tgsw_decomp_bits_, param.lwe_n_)
	for i := 0; i < lwe_key.N; i++ {
		tgsw_sample := key.ExtractTGSWSample(i)
		TGSWEncryptBinary(tgsw_sample, lwe_key.Key[i], tgsw_key)
	}
	return key
}

////////////////////////////////////////////////////////////////////////////////

/*
void SetSeed(uint32_t seed) {
	srand(seed);
	RandomGeneratorSetSeed(&seed, 1);
}
*/

func PubKeyGen(pri_key *PriKey) *PubKey {
	/*
		param := GetDefaultParam()
		return &PubKey{
			bk_: NewBootstrappingKey(param.tlwe_n_, param.tlwe_k_,
				param.tgsw_decomp_size_,
				param.tgsw_decomp_bits_, param.lwe_n_),
			ksk_: NewKeySwitchingKey(param.lwe_n_,
				param.keyswitching_decomp_size_,
				param.keyswitching_decomp_bits_,
				param.tlwe_n_*param.tlwe_k_),
		}
	*/

	lwe_key_extract := pri_key.tlwe_key_.ExtractLWEKey()
	return &PubKey{
		bk_: BootstrappingKeyGen(pri_key.lwe_key_, pri_key.tlwe_key_),
		//lwe_key_extract := NewLWEKey()
		ksk_: KeySwitchingKeyGen(pri_key.lwe_key_,
			lwe_key_extract),
	}
}

func PriKeyGen() *PriKey {
	//return NewPriKey(false)
	param := GetDefaultParam()
	return &PriKey{
		lwe_key_:  LWEKeyGen(param),
		tlwe_key_: TLWEKeyGen(param),
	}
}

func KeyGen() (*PubKey, *PriKey) {
	pri_key := PriKeyGen()
	pub_key := PubKeyGen(pri_key)
	return pub_key, pri_key
}

func Encrypt(ctxt *Ctxt, ptxt *Ptxt, pri_key *PriKey) {
	//assert(Find(ptxt, kPtxtSet));
	//  Torus mu = TorusFromDouble((double)1.0 * ptxt.message_ / Ptxt::kPtxtSpace);
	one := ModSwitchToTorus(1, 8)
	//mu := ptxt.message_ ? one : -one;
	var mu = one
	if ptxt.Message == 0 {
		mu = -one
	}

	LWEEncrypt(ctxt.lwe_sample_, mu, pri_key.lwe_key_)
}

func Decrypt(ptxt *Ptxt, ctxt *Ctxt, pri_key *PriKey) {
	mu := LWEDecrypt(ctxt.lwe_sample_, pri_key.lwe_key_, KPtxtSpace)
	//  ptxt.message_ = (uint32_t)int32_t((Ptxt::kPtxtSpace) * DoubleFromTorus(mu));
	//  ptxt.message_ %= ptxt.kPtxtSpace;

	if mu > 0 {
		ptxt.Message = 1
	} else {
		ptxt.Message = 0
	}

	//ptxt.message_ = mu > 0 ? 1 : 0;
}
