package tfhe

/** result = (0,mu) */
func LweNoiselessTrivial(result *LWESample, mu Torus) {
	p := GetDefaultParam()
	for i := 0; i < p.lwe_n_; i++ {
		result.A[i] = 0
	}
	result.B = mu
}

/** result = result - sample */
func LweSubTo(result *LWESample, sample *LWESample) {
	p := GetDefaultParam()
	for i := 0; i < p.lwe_n_; i++ {
		result.A[i] -= sample.A[i]
	}
	result.B -= sample.B
}

func Nand(out, ca, cb *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	//fix := ModSwitchToTorus(1, 8)

	p := GetDefaultParam()
	tempResult := NewLWESample(p.lwe_n_)

	//compute: (0,1/8) - ca - cb
	NandConst := ModSwitchToTorus(1, 8)
	LweNoiselessTrivial(tempResult, NandConst)
	LweSubTo(tempResult, ca.lwe_sample_)
	LweSubTo(tempResult, cb.lwe_sample_)

	//out.lwe_sample_.B += fix
	Bootstrap(out.lwe_sample_, out.lwe_sample_, mu, pub_key.bk_, pub_key.ksk_)
}

func Or(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i < in0.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = 0 + in0.lwe_sample_.A[i] + in1.lwe_sample_.A[i]
	}
	out.lwe_sample_.B += fix
	Bootstrap(out.lwe_sample_, out.lwe_sample_, mu, pub_key.bk_, pub_key.ksk_)
}

func And(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(-1, 8)
	for i := 0; i < in0.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = 0 + in0.lwe_sample_.A[i] + in1.lwe_sample_.A[i]
	}
	out.lwe_sample_.B += fix
	Bootstrap(out.lwe_sample_, out.lwe_sample_, mu, pub_key.bk_, pub_key.ksk_)
}

func Nor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(-1, 8)
	for i := 0; i < in0.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = 0 - in0.lwe_sample_.A[i] - in1.lwe_sample_.A[i]
	}
	out.lwe_sample_.B += fix
	Bootstrap(out.lwe_sample_, out.lwe_sample_, mu, pub_key.bk_, pub_key.ksk_)
}

func Xor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 4)
	for i := 0; i < in0.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = 0 + 2*in0.lwe_sample_.A[i] + 2*in1.lwe_sample_.A[i]
	}
	for i := 0; i < in0.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = 0 + 2*in0.lwe_sample_.A[i] + 2*in1.lwe_sample_.A[i]
	}
	out.lwe_sample_.B += fix
	Bootstrap(out.lwe_sample_, out.lwe_sample_, mu, pub_key.bk_, pub_key.ksk_)
}

func Xnor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(-1, 4)
	for i := 0; i < in0.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = 0 - 2*in0.lwe_sample_.A[i] - 2*in1.lwe_sample_.A[i]
	}
	out.lwe_sample_.B += fix
	Bootstrap(out.lwe_sample_, out.lwe_sample_, mu, pub_key.bk_, pub_key.ksk_)
}

func Not(out, in *Ctxt) {
	for i := 0; i < in.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = -in.lwe_sample_.A[i]
	}
}

func Copy(out, in *Ctxt) {
	for i := 0; i < in.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = in.lwe_sample_.A[i]
	}
}
