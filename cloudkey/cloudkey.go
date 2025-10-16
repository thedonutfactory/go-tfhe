package cloudkey

import (
	"sync"

	"github.com/thedonutfactory/go-tfhe/key"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/poly"
	"github.com/thedonutfactory/go-tfhe/tlwe"
	"github.com/thedonutfactory/go-tfhe/trgsw"
	"github.com/thedonutfactory/go-tfhe/trlwe"
	"github.com/thedonutfactory/go-tfhe/utils"
)

// CloudKey contains the public evaluation keys
type CloudKey struct {
	DecompositionOffset params.Torus
	BlindRotateTestvec  *trlwe.TRLWELv1
	KeySwitchingKey     []*tlwe.TLWELv0
	BootstrappingKey    []*trgsw.TRGSWLv1FFT
}

// NewCloudKey generates a new cloud key from a secret key
func NewCloudKey(secretKey *key.SecretKey) *CloudKey {
	return &CloudKey{
		DecompositionOffset: genDecompositionOffset(),
		BlindRotateTestvec:  genTestvec(),
		KeySwitchingKey:     genKeySwitchingKey(secretKey),
		BootstrappingKey:    genBootstrappingKey(secretKey),
	}
}

// NewCloudKeyNoKSK creates a cloud key without key switching key (for testing)
func NewCloudKeyNoKSK() *CloudKey {
	base := 1 << params.GetTRGSWLv1().BASEBIT
	iksT := params.GetTRGSWLv1().IKS_T
	n := params.GetTRGSWLv1().N
	lv0N := params.GetTLWELv0().N

	ksk := make([]*tlwe.TLWELv0, base*iksT*n)
	for i := range ksk {
		ksk[i] = tlwe.NewTLWELv0()
	}

	polyEval := poly.NewEvaluator(n)
	bsk := make([]*trgsw.TRGSWLv1FFT, lv0N)
	for i := range bsk {
		bsk[i] = trgsw.NewTRGSWLv1FFTDummy(polyEval)
	}

	return &CloudKey{
		DecompositionOffset: genDecompositionOffset(),
		BlindRotateTestvec:  genTestvec(),
		KeySwitchingKey:     ksk,
		BootstrappingKey:    bsk,
	}
}

// genDecompositionOffset generates the decomposition offset
func genDecompositionOffset() params.Torus {
	var offset params.Torus
	l := params.GetTRGSWLv1().L
	bg := params.GetTRGSWLv1().BG
	bgbit := params.GetTRGSWLv1().BGBIT

	for i := 0; i < l; i++ {
		offset += params.Torus(bg/2) * params.Torus(1<<(32-((i+1)*int(bgbit))))
	}

	return offset
}

// genTestvec generates the test vector for blind rotation
func genTestvec() *trlwe.TRLWELv1 {
	n := params.GetTRGSWLv1().N
	testvec := trlwe.NewTRLWELv1()
	bTorus := utils.F64ToTorus(0.125)

	for i := 0; i < n; i++ {
		testvec.A[i] = 0
		testvec.B[i] = bTorus
	}

	return testvec
}

// genKeySwitchingKey generates the key switching key
func genKeySwitchingKey(secretKey *key.SecretKey) []*tlwe.TLWELv0 {
	basebit := params.GetTRGSWLv1().BASEBIT
	iksT := params.GetTRGSWLv1().IKS_T
	base := 1 << basebit
	n := params.GetTRGSWLv1().N

	result := make([]*tlwe.TLWELv0, base*iksT*n)
	for i := range result {
		result[i] = tlwe.NewTLWELv0()
	}

	for i := 0; i < n; i++ {
		for j := 0; j < iksT; j++ {
			for k := 0; k < base; k++ {
				if k == 0 {
					continue
				}
				shift := uint((j + 1) * basebit)
				p := (float64(k) * float64(secretKey.KeyLv1[i])) / float64(uint64(1)<<shift)
				idx := (base * iksT * i) + (base * j) + k
				result[idx] = tlwe.NewTLWELv0().EncryptF64(p, params.KSKAlpha(), secretKey.KeyLv0)
			}
		}
	}

	return result
}

// genBootstrappingKey generates the bootstrapping key (parallelized)
func genBootstrappingKey(secretKey *key.SecretKey) []*trgsw.TRGSWLv1FFT {
	lv0N := params.GetTLWELv0().N
	result := make([]*trgsw.TRGSWLv1FFT, lv0N)

	var wg sync.WaitGroup
	for i := 0; i < lv0N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			polyEval := poly.NewEvaluator(params.GetTRGSWLv1().N)
			trgswCipher := trgsw.NewTRGSWLv1().EncryptTorus(
				secretKey.KeyLv0[idx],
				params.BSKAlpha(),
				secretKey.KeyLv1,
				polyEval,
			)
			result[idx] = trgsw.NewTRGSWLv1FFT(trgswCipher, polyEval)
		}(i)
	}
	wg.Wait()

	return result
}
