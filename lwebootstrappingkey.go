package tfhe

type LweBootstrappingKey struct {
	InOutParams   *LweParams       ///< paramètre de l'input et de l'output. key: s
	BkParams      *TGswParams      ///< params of the Gsw elems in bk. key: s"
	AccumParams   *TLweParams      ///< params of the accum variable key: s"
	ExtractParams *LweParams       ///< params after extraction: key: s'
	Bk            []*TGswSample    ///< the bootstrapping key (s.s")
	Ks            *LweKeySwitchKey ///< the keyswitch key (s'.s)
}

func NewLweBootstrappingKey(ksT, ksBasebit int32, inOutParams *LweParams,
	bkParams *TGswParams) *LweBootstrappingKey {

	accumParams := bkParams.TlweParams
	extractParams := &accumParams.ExtractedLweparams
	n := inOutParams.N
	N := extractParams.N

	bk := NewTGswSampleArray(n, bkParams)
	ks := NewLweKeySwitchKey(N, ksT, ksBasebit, inOutParams)

	return &LweBootstrappingKey{
		InOutParams:   inOutParams,
		BkParams:      bkParams,
		AccumParams:   accumParams,
		ExtractParams: extractParams,
		Bk:            bk,
		Ks:            ks,
	}
}
