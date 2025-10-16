package trgsw

import (
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

// IdentityKeySwitchingAssign performs identity key switching and writes to output
// Zero-allocation version
func IdentityKeySwitchingAssign(src *tlwe.TLWELv1, keySwitchingKey []*tlwe.TLWELv0, output *tlwe.TLWELv0) {
	n := params.GetTRGSWLv1().N
	basebit := params.GetTRGSWLv1().BASEBIT
	base := 1 << basebit
	iksT := params.GetTRGSWLv1().IKS_T
	tlweLv0N := params.GetTLWELv0().N

	// Clear output
	for i := 0; i < len(output.P); i++ {
		output.P[i] = 0
	}
	output.P[tlweLv0N] = src.P[len(src.P)-1]

	precOffset := params.Torus(1 << (32 - (1 + basebit*iksT)))

	for i := 0; i < n; i++ {
		aBar := src.P[i] + precOffset
		for j := 0; j < iksT; j++ {
			k := (aBar >> (32 - (j+1)*basebit)) & params.Torus((1<<basebit)-1)
			if k != 0 {
				idx := (base * iksT * i) + (base * j) + int(k)
				for x := 0; x < len(output.P); x++ {
					output.P[x] -= keySwitchingKey[idx].P[x]
				}
			}
		}
	}
}
