package tfhe

func Nand(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 - in0.lwe_sample.A[i] - in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

func Or(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 + in0.lwe_sample.A[i] + in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

func And(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(-1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 + in0.lwe_sample.A[i] + in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

func Nor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(-1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 - in0.lwe_sample.A[i] - in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

func Xor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 4)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 + 2*in0.lwe_sample.A[i] + 2*in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

func Xnor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(-1, 4)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 - 2*in0.lwe_sample.A[i] - 2*in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

func Not(out, in *Ctxt) {
	for i := 0; i <= in.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = -in.lwe_sample.A[i]
	}
}

func Copy(out, in *Ctxt) {
	for i := 0; i <= in.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = in.lwe_sample.A[i]
	}
}
