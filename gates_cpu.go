package tfhe

// Homomorphic bootstrapped NAND gate
// Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
func Nand(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 - in0.lwe_sample.A[i] - in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

// Homomorphic bootstrapped OR gate
// Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
func Or(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 + in0.lwe_sample.A[i] + in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

// Homomorphic bootstrapped AND gate
// Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
func And(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(-1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 + in0.lwe_sample.A[i] + in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

// Homomorphic bootstrapped NOR gate
// Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
func Nor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(-1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 - in0.lwe_sample.A[i] - in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

// Homomorphic bootstrapped XOR gate
// Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
func Xor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 4)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 + 2*in0.lwe_sample.A[i] + 2*in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

// Homomorphic bootstrapped XNOR gate
// Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
func Xnor(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(-1, 4)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 - 2*in0.lwe_sample.A[i] - 2*in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

// Homomorphic bootstrapped NOT gate (doesn't need to be bootstrapped)
// Takes in input 1 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
func Not(out, in *Ctxt) {
	for i := 0; i <= in.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = -in.lwe_sample.A[i]
	}
	*out.lwe_sample.B = -*in.lwe_sample.B
}

// Homomorphic bootstrapped COPY gate (doesn't need to be bootstrapped)
// Takes in input 1 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
func Copy(out, in *Ctxt) {
	for i := 0; i <= in.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = in.lwe_sample.A[i]
	}
	*out.lwe_sample.B = *in.lwe_sample.B
}

// Homomorphic Trivial Constant gate (doesn't need to be bootstrapped)
// Takes a boolean value)
// Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
func Constant(out *Ctxt, value int32) {
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i <= out.lwe_sample.N; i++ {
		out.lwe_sample.A[i] += 0
	}
	*out.lwe_sample.B = fix
}

// Homomorphic bootstrapped AndNY Gate: not(a) and b
// Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
func AndNY(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(-1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 - in0.lwe_sample.A[i] + in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

// Homomorphic bootstrapped AndYN Gate: a and not(b)
// Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
func AndYN(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(-1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 + in0.lwe_sample.A[i] - in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

// Homomorphic bootstrapped OrNY Gate: not(a) or b
// Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
func OrNY(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 - in0.lwe_sample.A[i] + in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}

// Homomorphic bootstrapped OrYN Gate: a or not(b)
// Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
// Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
func OrYN(out, in0, in1 *Ctxt, pub_key *PubKey) {
	mu := ModSwitchToTorus(1, 8)
	fix := ModSwitchToTorus(1, 8)
	for i := 0; i <= in0.lwe_sample.N; i++ {
		out.lwe_sample.A[i] = 0 + in0.lwe_sample.A[i] - in1.lwe_sample.A[i]
	}
	*out.lwe_sample.B += fix
	Bootstrap(out.lwe_sample, out.lwe_sample, mu, pub_key.Bk, pub_key.Ksk)
}
