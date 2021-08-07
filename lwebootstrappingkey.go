package tfhe

type LweBootstrappingKey struct {
	in_out_params  *LweParams       ///< paramètre de l'input et de l'output. key: s
	bk_params      *TGswParams      ///< params of the Gsw elems in bk. key: s"
	accum_params   *TLweParams      ///< params of the accum variable key: s"
	extract_params *LweParams       ///< params after extraction: key: s'
	bk             []*TGswSample    ///< the bootstrapping key (s.s")
	ks             *LweKeySwitchKey ///< the keyswitch key (s'.s)
}

/*
func NewLweBootstrappingKey(in_out_params *LweParams,
	bk_params *TGswParams,
	accum_params *TLweParams,
	extract_params *LweParams,
	bk *TGswSample,
	ks *LweKeySwitchKey) *LweBootstrappingKey {

	return &LweBootstrappingKey{
		in_out_params:  in_out_params,
		bk_params:      bk_params,
		accum_params:   accum_params,
		extract_params: extract_params,
		bk:             bk,
		ks:             ks,
	}
}
*/

func NewLweBootstrappingKey(ks_t, ks_basebit int32, in_out_params *LweParams,
	bk_params *TGswParams) *LweBootstrappingKey {

	accum_params := bk_params.TlweParams
	extract_params := &accum_params.extractedLweparams
	n := in_out_params.N
	N := extract_params.N

	bk := NewTGswSampleArray(n, bk_params)
	ks := NewLweKeySwitchKey(N, ks_t, ks_basebit, in_out_params)

	return &LweBootstrappingKey{
		in_out_params:  in_out_params,
		bk_params:      bk_params,
		accum_params:   accum_params,
		extract_params: extract_params,
		bk:             bk,
		ks:             ks,
	}
}

/*
EXPORT void init_LweBootstrappingKey(LweBootstrappingKey *obj, int32_t ks_t, int32_t ks_basebit, const LweParams *in_out_params,
	const TGswParams *bk_params) {
	const TLweParams *accum_params = bk_params.tlwe_params;
	const LweParams *extract_params = &accum_params.extracted_lweparams;
	const int32_t n = in_out_params.n;
	const int32_t N = extract_params.n;

	TGswSample *bk = new_TGswSample_array(n, bk_params);
	LweKeySwitchKey *ks = new_LweKeySwitchKey(N, ks_t, ks_basebit, in_out_params);

	new(obj) LweBootstrappingKey(in_out_params, bk_params, accum_params, extract_params, bk, ks);
}
*/
