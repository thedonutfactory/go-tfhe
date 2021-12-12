package tfhe

func ModSwitchFromTorus(phase Torus, space uint32) uint32 {
	interv := ((uint64(1) << 63) / uint64(space)) * 2
	half_interval := interv / 2
	phase64 := (uint64(phase) << 32) + half_interval
	return uint32(phase64 / interv)
}

func PolyMulPowX(out []Torus, in []Torus, exp, n int) {
	temp := make([]Torus, n) //  new Torus[n];
	Assert(exp >= 0 && exp <= 2*n)
	if exp == 2*n {
		for i := 0; i < n; i++ {
			temp[i] = in[i]
		}
	} else if exp < n {
		for i := 0; i < exp; i++ {
			temp[i] = -in[i-exp+n]
		}
		for i := exp; i < n; i++ {
			temp[i] = in[i-exp]
		}
	} else {
		exp -= n
		for i := 0; i < exp; i++ {
			temp[i] = in[i-exp+n]
		}
		for i := exp; i < n; i++ {
			temp[i] = -in[i-exp]
		}
	}
	for i := 0; i < n; i++ {
		out[i] = temp[i]
	}
}

func PolyMulAdd(out, in0 []Torus, in1 []int32, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			out[i] += in0[j] * in1[i-j]
		}
		for j := i + 1; j < n; j++ {
			out[i] -= in0[j] * in1[n+i-j]
		}
	}
}

func PolySub(out, in0, in1 []Torus, n int) {
	for i := 0; i < n; i++ {
		out[i] = in0[i] - in1[i]
	}
}

func PolyDecomp(out [][]int32, in []Torus, n, bits, l, mask, half, offset int32) {
	for j := int32(0); j < n; j++ {
		for i := int32(0); i < l; i++ {
			out[i][j] = (((in[j] + offset) >> (32 - (i+1)*bits)) & mask) - half
		}
	}
}

func LWESampleSub(out, in0, in1 *LWESample) {
	for i := 0; i <= out.N; i++ {
		out.A[i] = in0.A[i] - in1.A[i]
	}
}

func Bootstrap(out, in *LWESample, mu Torus, bk *BootstrappingKey, ksk *KeySwitchingKey) {
	lwe_n := ksk.N
	tlwe_n := bk.N
	n2 := uint32(2 * tlwe_n)
	k := bk.K
	bk_l := bk.L
	bk_bits := bk.W
	bk_mask := (1 << bk_bits) - 1
	bk_half := 1 << (bk_bits - 1)
	bk_offset := 0
	for i := 0; i < bk_l; i++ {
		bk_offset += 0x1 << (32 - (i+1)*bk_bits)
	}
	bk_offset *= bk_half
	kpl := (k + 1) * bk_l
	ksk_l := ksk.L
	ksk_bits := ksk.W
	ksk_mask := (0x1 << ksk_bits) - 1
	ksk_offset := 0x1 << (31 - ksk_l*ksk_bits)
	temp := make([]Torus, tlwe_n) // new Torus[tlwe_n];
	acc := NewTLWESample(bk.N, bk.K)
	//std::pair<void*, MemoryDeleter> pair = AllocatorCPU::New(acc->SizeMalloc());
	//acc->set_data((TLWESample::PointerType)pair.first);
	//MemoryDeleter acc_deleter = pair.second;
	//int** decomp = new int32_t*[kpl];
	decomp := make([][]int32, kpl)
	for i := 0; i < kpl; i++ {
		decomp[i] = make([]int32, tlwe_n) //new int32_t[tlwe_n];
	}

	bar_b := ModSwitchFromTorus(*in.B, n2)
	for i := 0; i < tlwe_n; i++ {
		temp[i] = mu
	}
	//memset(acc->data(), 0, acc->SizeData());
	PolyMulPowX(acc.b(), temp, int(n2-bar_b), tlwe_n)

	//var bar_a Torus

	for i := 0; i < lwe_n; i++ {
		bar_a := ModSwitchFromTorus(in.A[i], n2)
		for j := 0; j <= k; j++ {
			PolyMulPowX(temp, acc.ExtractPoly(j), int(bar_a), tlwe_n)
			PolySub(temp, temp, acc.ExtractPoly(j), tlwe_n)
			//PolyDecomp(decomp+j*bk_l, temp, int32(tlwe_n), int32(bk_bits), int32(bk_l), int32(bk_mask), int32(bk_half), int32(bk_offset))
			PolyDecomp(decomp[j*bk_l:], temp, int32(tlwe_n), int32(bk_bits), int32(bk_l), int32(bk_mask), int32(bk_half), int32(bk_offset))
		}
		for j := 0; j <= k; j++ {
			for p := 0; p < kpl; p++ {
				PolyMulAdd(acc.ExtractPoly(j),
					bk.ExtractTGSWSample(i).ExtractTLWESample(p).ExtractPoly(j),
					decomp[p], tlwe_n)
			}
		}
	}

	var coeff, digit int32
	//memset(out->data(), 0, out->SizeData());
	//out.A = make([]int32, out.N)
	out = NewLWESample(lwe_n)
	*out.B = acc.b()[0]
	for i := 0; i < k*tlwe_n; i++ {
		if i == 0 {
			coeff = acc.a(acc.K)[i]
		} else {
			coeff = -acc.a(acc.K)[k*tlwe_n-i]
		}
		coeff += int32(ksk_offset)
		for j := 0; j < ksk_l; j++ {
			digit = (coeff >> (32 - (j+1)*ksk_bits)) & int32(ksk_mask)
			if digit != 0 {
				ksk_entry := ksk.A[i][j][digit]
				//ksk_entry := ksk.ExtractLWESample(ksk.GetLWESampleIndex(i, j, int(digit)))
				LWESampleSub(out, out, ksk_entry)
			}
		}
	}

}
