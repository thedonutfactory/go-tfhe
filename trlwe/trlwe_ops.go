package trlwe

import (
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

// SampleExtractIndexAssign extracts a TLWE sample from TRLWE at index k and writes to output
// Zero-allocation version
func SampleExtractIndexAssign(trlwe *TRLWELv1, k int, output *tlwe.TLWELv1) {
	n := params.GetTRLWELv1().N

	for i := 0; i < n; i++ {
		if i <= k {
			output.P[i] = trlwe.A[k-i]
		} else {
			output.P[i] = ^params.Torus(0) - trlwe.A[n+k-i]
		}
	}
	output.SetB(trlwe.B[k])
}
