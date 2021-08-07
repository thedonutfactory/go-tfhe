package tfhe

type LweBootstrappingKey struct {
	inOutParams   *LweParams       ///< paramètre de l'input et de l'output. key: s
	bkParams      *TGswParams      ///< params of the Gsw elems in bk. key: s"
	accumParams   *TLweParams      ///< params of the accum variable key: s"
	extractParams *LweParams       ///< params after extraction: key: s'
	bk            []*TGswSample    ///< the bootstrapping key (s.s")
	ks            *LweKeySwitchKey ///< the keyswitch key (s'.s)
}

func NewLweBootstrappingKey(ksT, ksBasebit int32, inOutParams *LweParams,
	bkParams *TGswParams) *LweBootstrappingKey {

	accumParams := bkParams.TlweParams
	extractParams := &accumParams.extractedLweparams
	n := inOutParams.N
	N := extractParams.N

	bk := NewTGswSampleArray(n, bkParams)
	ks := NewLweKeySwitchKey(N, ksT, ksBasebit, inOutParams)

	return &LweBootstrappingKey{
		inOutParams:   inOutParams,
		bkParams:      bkParams,
		accumParams:   accumParams,
		extractParams: extractParams,
		bk:            bk,
		ks:            ks,
	}
}
