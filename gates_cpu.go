package tfhe

func Nand(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i < in0.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = 0 - in0.lwe_sample_.A[i] - in1.lwe_sample_.A[i]
	}
	out.lwe_sample_.B += fix
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
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i < in0.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = 0 + in0.lwe_sample_.A[i] + in1.lwe_sample_.A[i]
	}
	out.lwe_sample_.B += fix
	Bootstrap(out.lwe_sample_, out.lwe_sample_, mu, pub_key.bk_, pub_key.ksk_)
}

func Nor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i < in0.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = 0 - in0.lwe_sample_.A[i] - in1.lwe_sample_.A[i]
	}
	out.lwe_sample_.B += fix
	Bootstrap(out.lwe_sample_, out.lwe_sample_, mu, pub_key.bk_, pub_key.ksk_)
}

func Xor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i < in0.lwe_sample_.N; i++ {
		out.lwe_sample_.A[i] = 0 + 2*in0.lwe_sample_.A[i] + 2*in1.lwe_sample_.A[i]
	}
	out.lwe_sample_.B += fix
	Bootstrap(out.lwe_sample_, out.lwe_sample_, mu, pub_key.bk_, pub_key.ksk_)
}

func Xnor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 8)
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
