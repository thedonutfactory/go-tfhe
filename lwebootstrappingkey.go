package tfhe

type LweBootstrappingKey struct {
	in_out_params  *LweParams       ///< paramètre de l'input et de l'output. key: s
	bk_params      *TGswParams      ///< params of the Gsw elems in bk. key: s"
	accum_params   *TLweParams      ///< params of the accum variable key: s"
	extract_params *LweParams       ///< params after extraction: key: s'
	bk             *TGswSample      ///< the bootstrapping key (s->s")
	ks             *LweKeySwitchKey ///< the keyswitch key (s'->s)
}

