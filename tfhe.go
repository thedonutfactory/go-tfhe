package tfhe

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

const KPtxtSpace = 2

type Param struct {
	lwe_n                    int
	tlwe_n                   int
	tlwe_k                   int
	tgsw_decomp_bits         int
	tgsw_decomp_size         int
	keyswitching_decomp_bits int
	keyswitching_decomp_size int
	lwe_noise                float64
	tlwe_noise               float64
}

func GetDefaultParam() *Param {
	return &Param{
		lwe_n:                    500,
		tlwe_n:                   1024,
		tlwe_k:                   1,
		tgsw_decomp_bits:         10,
		tgsw_decomp_size:         2,
		keyswitching_decomp_bits: 2,
		keyswitching_decomp_size: 8,
		lwe_noise:                math.Pow(2.0, -15),
		tlwe_noise:               9.e-9,
	}
}

/**
* Private Key.
* Necessary for encryption/decryption and public key generation.
 */
type PriKey struct {
	Lwe_key  *LWEKey
	Tlwe_key *TLWEKey
}

func NewPriKey(isAlias bool) *PriKey {
	param := GetDefaultParam()
	return &PriKey{
		Lwe_key:  NewLWEKey(param.lwe_n),
		Tlwe_key: NewTLWEKey(param.tlwe_n, param.tlwe_k),
	}
}

/**
* Public Key.
* Necessary for a server to perform homomorphic evaluation.
 */
type PubKey struct {
	Bk  *BootstrappingKey
	Ksk *KeySwitchingKey
}

func NewPubKey(isAlias bool) *PubKey {
	param := GetDefaultParam()
	return &PubKey{
		Bk: NewBootstrappingKey(param.tlwe_n, param.tlwe_k,
			param.tgsw_decomp_size,
			param.tgsw_decomp_bits, param.lwe_n),
		Ksk: NewKeySwitchingKey(param.lwe_n,
			param.keyswitching_decomp_size,
			param.keyswitching_decomp_bits,
			param.tlwe_n*param.tlwe_k),
	}
}

/** Ciphertext. */
type Ctxt struct {
	lwe_sample        *LWESample
	lwe_sample_device *LWESample
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

func NewPtxtInit(message uint32) *Ptxt {
	return &Ptxt{
		kPtxtSpace: 2,
		Message:    message,
	}
}

func NewPtxtArray(n int) (arr []*Ptxt) {
	arr = make([]*Ptxt, n)
	for i := 0; i < n; i++ {
		arr[i] = NewPtxt()
	}
	return
}

func SDFromBound(noise_bound float64) float64 {
	return noise_bound * math.Sqrt(2.0/math.Pi)
}

// Conversions go back to -0.5~0.5 if Torus is int32!
func TorusFromDouble(d float64) Torus {
	return Torus(math.Round(math.Mod(d, 1) * float64(int64(1)<<32)))
}

func DoubleFromTorus(t Torus) float64 {
	return float64(t) / math.Pow(2, 32)
}

func ApproxPhase(phase Torus, msg_space int32) Torus {
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

func LWEKeyGen(param *Param) *LWEKey {
	dist := distuv.Uniform{
		Min: 0,
		Max: 1,
	}

	key := NewLWEKey(param.lwe_n)
	for i := range key.Key {
		key.Key[i] = uint32(math.Round(dist.Rand()))
	}
	return key
}

func TLWEKeyGen(param *Param) *TLWEKey {
	key := NewTLWEKey(param.tlwe_n, param.tlwe_k)
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
	noise_bound := GetDefaultParam().lwe_noise
	d1 := NormalDist(0.0, SDFromBound(noise_bound))
	*ct.B = pt + TorusFromDouble(d1.Rand())
	d2 := UniformDist(math.MinInt32, math.MaxInt32)
	for i := 0; i < key.N; i++ {
		ct.A[i] = Torus(d2.Rand())
		*ct.B += ct.A[i] * int32(key.Key[i])
	}
}

func LWEEncryptExternalNoise(ct *LWESample, pt Torus, key *LWEKey, noise float64) {
	*ct.B = pt + TorusFromDouble(noise)
	d1 := UniformDist(math.MinInt32, math.MaxInt32)
	for i := 0; i < key.N; i++ {
		ct.A[i] = Torus(d1.Rand())
		*ct.B += ct.A[i] * int32(key.Key[i])
	}
}

func LWEDecrypt(ct *LWESample, key *LWEKey, space int32) Torus {
	Assert(ct.N == key.N)
	err := *ct.B

	for i := 0; i < ct.N; i++ {
		err -= ct.A[i] * int32(key.Key[i])
	}
	return err
}

func KeySwitchingKeyGen(lwe_key_to, lwe_key_from *LWEKey) *KeySwitchingKey {
	param := GetDefaultParam()
	key := NewKeySwitchingKey(param.lwe_n,
		param.keyswitching_decomp_size,
		param.keyswitching_decomp_bits,
		param.tlwe_n*param.tlwe_k)

	var mu Torus
	var temp uint32
	total := key.NumLWESamples()
	noise := make([]float64, total)
	var err float64 = 0
	for i := 0; i < total; i++ {
		nd := NormalDist(0.0, SDFromBound(param.lwe_noise))
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
			for k := 0; k < (0x1 << key.W); k++ {
				lwe_sample := key.A[i][j][k]
				mu = Torus((temp * uint32(k)) * (1 << (32 - (j+1)*key.W)))
				LWEEncryptExternalNoise(lwe_sample, mu, lwe_key_to, noise[index])
				index++
			}
		}
	}
	return key
}

func TLWEEncryptZero(ct *TLWESample, key *TLWEKey) {
	noise_bound := GetDefaultParam().tlwe_noise
	dist_b := NormalDist(0.0, SDFromBound(noise_bound))
	for i := 0; i < key.N; i++ {
		ct.b()[i] = TorusFromDouble(dist_b.Rand())
	}

	dist_a := UniformDist(math.MinInt32, math.MaxInt32)
	for i := 0; i < key.K; i++ {
		for j := 0; j < key.N; j++ {
			ct.a(i)[j] = Torus(dist_a.Rand())
		}
		PolyMulAddBinary(ct.b(), ct.A[i], key.A[0], key.N)
	}
}

func TGSWEncryptBinary(ct *TGSWSample, pt Binary, key *TGSWKey) {
	l := ct.L
	k := ct.K
	w := ct.W
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
	key := NewBootstrappingKey(param.tlwe_n, param.tlwe_k,
		param.tgsw_decomp_size,
		param.tgsw_decomp_bits, param.lwe_n)
	for i := 0; i < lwe_key.N; i++ {
		tgsw_sample := key.ExtractTGSWSample(i)
		TGSWEncryptBinary(tgsw_sample, lwe_key.Key[i], tgsw_key)
	}
	return key
}

func PubKeyGen(pri_key *PriKey) *PubKey {
	lwe_key_extract := pri_key.Tlwe_key.ExtractLWEKey()
	return &PubKey{
		Bk: BootstrappingKeyGen(pri_key.Lwe_key, pri_key.Tlwe_key),
		Ksk: KeySwitchingKeyGen(pri_key.Lwe_key,
			lwe_key_extract),
	}
}

func PriKeyGen() *PriKey {
	param := GetDefaultParam()
	return &PriKey{
		Lwe_key:  LWEKeyGen(param),
		Tlwe_key: TLWEKeyGen(param),
	}
}

func KeyGen() (*PubKey, *PriKey) {
	pri_key := PriKeyGen()
	pub_key := PubKeyGen(pri_key)
	return pub_key, pri_key
}

func Encrypt(ctxt *Ctxt, ptxt *Ptxt, pri_key *PriKey) {
	one := ModSwitchToTorus(1, 8)
	var mu = one
	if ptxt.Message == 0 {
		mu = -one
	}
	LWEEncrypt(ctxt.lwe_sample, mu, pri_key.Lwe_key)
}

func Decrypt(ptxt *Ptxt, ctxt *Ctxt, pri_key *PriKey) {
	mu := LWEDecrypt(ctxt.lwe_sample, pri_key.Lwe_key, KPtxtSpace)
	if mu > 0 {
		ptxt.Message = 1
	} else {
		ptxt.Message = 0
	}
}
