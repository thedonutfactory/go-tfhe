package tfhe

type Operations interface {
	// comparison
	CompareBit(a, b, lsbCarry, tmp *LweSample) (result *LweSample)
	Equals(a, b []*LweSample, nbBits int) (result []*LweSample)
	Minimum(a, b []*LweSample, nbBits int) (result []*LweSample)
	Maximum(a, b []*LweSample, nbBits int) (result []*LweSample)
	Gte(a, b []*LweSample, nbBits int) (result []*LweSample)

	// arithmetic
	Add(a, b []*LweSample, nbBits int) (result []*LweSample)
	Sub(a, b []*LweSample, nbBits int) (result []*LweSample)
	Mul(a, b []*LweSample, nbBits int) (result []*LweSample)
	Div(a, b []*LweSample, nbBits int) (result []*LweSample)
	Pow(a []*LweSample, n, nbBits int) (result []*LweSample)

	// bitwise shift
	ShiftLeft(a []*LweSample, positions, nbBits int) (result []*LweSample)
	ShiftRight(a []*LweSample, positions, nbBits int) (result []*LweSample)
	UshiftLeft(a []*LweSample, positions, nbBits int) (result []*LweSample)
	UshiftRight(a []*LweSample, positions, nbBits int) (result []*LweSample)

	// bitwise operations
	Nand(a, b []*LweSample, nbBits int) (result []*LweSample)
	Or(a, b []*LweSample, nbBits int) (result []*LweSample)
	And(a, b []*LweSample, nbBits int) (result []*LweSample)
	Xor(a, b []*LweSample, nbBits int) (result []*LweSample)
	Xnor(a, b []*LweSample, nbBits int) (result []*LweSample)
	Not(a []*LweSample, nbBits int) (result []*LweSample)
	Nor(a, b []*LweSample, nbBits int) (result []*LweSample)
	AndNY(a, b []*LweSample, nbBits int) (result []*LweSample)
	AndYN(a, b []*LweSample, nbBits int) (result []*LweSample)
	OrNY(a, b []*LweSample, nbBits int) (result []*LweSample)
	OrYN(a, b []*LweSample, nbBits int) (result []*LweSample)
	Mux(a, b, c []*LweSample, nbBits int) (result []*LweSample)

	// misc
	Copy(a, b []*LweSample, nbBits int) (result []*LweSample)
	Constant(value int64, nbBits int) (result []*LweSample)
}

type CipheredOperations struct {
	bk *TFheGateBootstrappingCloudKeySet
}

// elementary full comparator gate that is used to compare the i-th bit:
//   input: ai and bi the i-th bit of a and b
//          lsb_carry: the result of the comparison on the lowest bits
//   algo: if (a==b) return lsb_carry else return b
func (ops *CipheredOperations) CompareBit(a, b, lsbCarry, tmp *LweSample) (result *LweSample) {
	result = NewLweSample(ops.bk.params.InOutParams)
	BootsXNOR(tmp, a, b, ops.bk)
	BootsMUX(result, tmp, lsbCarry, a, ops.bk)
	return result
}

// Returns a == b
func (ops *CipheredOperations) Equals(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	BootsCONSTANT(result[0], 1, ops.bk)
	for i := 0; i < nbBits; i++ {
		BootsXNOR(tmps[0], a[i], b[i], ops.bk)
		BootsAND(result[0], result[0], tmps[0], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) isNegative(result, a []*LweSample, nbBits int) {
	BootsCOPY(result[0], a[nbBits-1], ops.bk)
}

// this function compares two multibit words, and puts the min in result
func (ops *CipheredOperations) Minimum(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	aGreater := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	minimumSameSign := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	minimumOneNegative := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	oneNegative := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	negativoA := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	negativoB := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	ops.isNegative(negativoA, a, nbBits)
	ops.isNegative(negativoB, b, nbBits)

	BootsXOR(oneNegative[0], negativoA[0], negativoB[0], ops.bk)

	// a > b = soloOneNegative & is_negative(b)
	BootsAND(aGreater[0], oneNegative[0], negativoB[0], ops.bk)
	for i := 0; i < nbBits; i++ {
		BootsMUX(minimumOneNegative[i], aGreater[0], b[i], a[i], ops.bk)
	}

	//initialize the carry to 0
	BootsCONSTANT(tmps[0], 0, ops.bk)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		tmps[0] = ops.CompareBit(a[i], b[i], tmps[0], tmps[1])
	}

	//tmps[0] is the result of the comparaison: 0 if a is larger, 1 if b is larger
	//select the max and copy it to the result
	for i := 0; i < nbBits; i++ {
		BootsMUX(minimumSameSign[i], tmps[0], b[i], a[i], ops.bk)
	}

	// Result depending on whether we compare the same sign or not
	for i := 0; i < nbBits; i++ {
		BootsMUX(result[i], oneNegative[0], minimumOneNegative[i], minimumSameSign[i], ops.bk)
	}
	return result
}

// this function compares two multibit words, and puts the min in result
func (ops *CipheredOperations) Maximum2(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	aGreater := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	minimumSameSign := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	minimumOneNegative := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	oneNegative := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	negativoA := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	negativoB := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	ops.isNegative(negativoA, a, nbBits)
	ops.isNegative(negativoB, b, nbBits)

	BootsXOR(oneNegative[0], negativoA[0], negativoB[0], ops.bk)

	// a > b = soloOneNegative & is_negative(b)
	BootsAND(aGreater[0], oneNegative[0], negativoB[0], ops.bk)
	for i := 0; i < nbBits; i++ {
		BootsMUX(minimumOneNegative[i], aGreater[0], b[i], a[i], ops.bk)
	}

	//initialize the carry to 0
	BootsCONSTANT(tmps[0], 0, ops.bk)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		tmps[0] = ops.CompareBit(a[i], b[i], tmps[0], tmps[1])
	}

	//tmps[0] is the result of the comparaison: 0 if a is larger, 1 if b is larger
	//select the max and copy it to the result
	for i := 0; i < nbBits; i++ {
		BootsMUX(minimumSameSign[i], tmps[0], a[i], b[i], ops.bk)
	}

	// Result depending on whether we compare the same sign or not
	for i := 0; i < nbBits; i++ {
		BootsMUX(result[i], oneNegative[0], minimumOneNegative[i], minimumSameSign[i], ops.bk)
	}
	return result
}

// this function compares two multibit words, and puts the min in result
func (ops *CipheredOperations) Maximum(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	aGreater := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	minimumSameSign := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	minimumOneNegative := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	oneNegative := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	negativoA := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	negativoB := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	ops.isNegative(negativoA, a, nbBits)
	ops.isNegative(negativoB, b, nbBits)

	BootsXOR(oneNegative[0], negativoA[0], negativoB[0], ops.bk)

	// a > b = soloOneNegative & is_negative(b)
	BootsAND(aGreater[0], oneNegative[0], negativoB[0], ops.bk)
	for i := 0; i < nbBits; i++ {
		BootsMUX(minimumOneNegative[i], aGreater[0], b[i], a[i], ops.bk)
	}

	//initialize the carry to 0
	BootsCONSTANT(tmps[0], 0, ops.bk)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		tmps[0] = ops.CompareBit(a[i], b[i], tmps[0], tmps[1])
	}

	//tmps[0] is the result of the comparaison: 0 if a is larger, 1 if b is larger
	//select the max and copy it to the result
	for i := 0; i < nbBits; i++ {
		BootsMUX(minimumSameSign[i], tmps[0], b[i], a[i], ops.bk)
	}

	// Todo - same as in minimum, but returning the opposite
	for i := 0; i < nbBits; i++ {
		//BootsMUX(result[i], oneNegative[0], minimumOneNegative[i], minimumMismoSigno[i], ops.bk)
		BootsMUX(result[i], oneNegative[0], minimumSameSign[i], minimumOneNegative[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) addBit(result, carry_out, a, b, carry_in *LweSample) {
	s1 := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	c1 := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	c2 := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	BootsCONSTANT(s1[0], 0, ops.bk)
	BootsCONSTANT(c1[0], 0, ops.bk)
	BootsCONSTANT(c2[0], 0, ops.bk)

	BootsXOR(s1[0], a, b, ops.bk)
	BootsXOR(result, s1[0], carry_in, ops.bk)

	BootsAND(c1[0], s1[0], carry_in, ops.bk)
	BootsAND(c2[0], a, b, ops.bk)
	BootsOR(carry_out, c1[0], c2[0], ops.bk)

}

// return -a
func (ops *CipheredOperations) negative(result, a []*LweSample, nbBits int) {

	ha_changed := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	not_x := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	for i := 0; i < 2; i++ {
		BootsCONSTANT(ha_changed[i], 0, ops.bk)
		BootsCONSTANT(not_x[i], 0, ops.bk)
	}

	for i := 0; i < nbBits; i++ {
		BootsNOT(not_x[0], a[i], ops.bk)
		BootsMUX(result[i], ha_changed[0], not_x[0], a[i], ops.bk)
		BootsOR(ha_changed[0], ha_changed[0], a[i], ops.bk)
	}

}

func (ops *CipheredOperations) Add(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	tmpsCarry := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	//initialize the carry to 0
	BootsCONSTANT(tmpsCarry[0], 0, ops.bk)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		ops.addBit(result[i], tmpsCarry[0], a[i], b[i], tmpsCarry[0])
	}
	return result
}

func (ops *CipheredOperations) Sub(a, b []*LweSample, nbBits int) (result []*LweSample) {
	res := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	ops.negative(res, b, nbBits)
	return ops.Add(a, res, nbBits)
}

// Unsigned multiply
func (ops *CipheredOperations) umul(result, a, b []*LweSample, nbBits int) {
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	aux2 := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	for i := 0; i < nbBits; i++ {
		BootsCONSTANT(aux[i], 0, ops.bk)
		BootsCONSTANT(aux2[i], 0, ops.bk)
	}

	// Multiply opA * opB
	for i := 0; i < nbBits/2; i++ {
		// Reset the auxs
		for j := 0; j < nbBits; j++ {
			BootsCONSTANT(aux[j], 0, ops.bk)
			BootsCONSTANT(aux2[j], 0, ops.bk)
		}

		for j := 0; j < (nbBits/2)+1; j++ {
			BootsAND(aux[j+i], a[i], b[j], ops.bk)
		}

		// add(aux2, aux, result, nbBits, bk);
		aux2 = ops.Add(aux, result, nbBits)
		//result = ops.Add(aux2, aux, nbBits)

		for j := 0; j < nbBits; j++ {
			BootsCOPY(result[j], aux2[j], ops.bk)
		}

	}
}

// multiply two ciphertexts and return the result
func (ops *CipheredOperations) Mul(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	aux2 := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	// Parameters to take into account negative numbers
	negatA := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	negatB := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	opA := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	opB := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	isNegativeA := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	isNegativeB := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	corrige := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	// Set number of bits so: nb(result) = nb(a)+nb(b)
	for i := 0; i < nbBits; i++ {
		BootsCONSTANT(aux[i], 0, ops.bk)
		BootsCONSTANT(aux2[i], 0, ops.bk)
		BootsCONSTANT(negatA[i], 0, ops.bk)
		BootsCONSTANT(negatB[i], 0, ops.bk)
		BootsCONSTANT(opA[i], 0, ops.bk)
		BootsCONSTANT(opB[i], 0, ops.bk)
		BootsCONSTANT(result[i], 0, ops.bk)
	}

	for i := 0; i < 2; i++ {
		BootsCONSTANT(isNegativeA[i], 0, ops.bk)
		BootsCONSTANT(isNegativeB[i], 0, ops.bk)
		BootsCONSTANT(corrige[i], 0, ops.bk)
	}

	// BEGIN SIGN LOGIC
	ops.negative(negatA, a, nbBits)
	ops.negative(negatB, b, nbBits)

	// Put the two numbers in positive
	opA = ops.Maximum(negatA, a, nbBits)
	opB = ops.Maximum(negatB, b, nbBits)

	// If only one of the two is negative, the result is negative
	ops.isNegative(isNegativeA, a, nbBits)
	ops.isNegative(isNegativeB, b, nbBits)
	BootsXOR(corrige[0], isNegativeA[0], isNegativeB[0], ops.bk)
	// END SIGN LOGIC

	ops.umul(result, opA, opB, nbBits)

	// BEGIN SIGN LOGIC
	// We determine whether to return the positive or negative result
	ops.negative(aux, result, nbBits)

	for i := 0; i < nbBits; i++ {
		BootsMUX(result[i], corrige[0], aux[i], result[i], ops.bk)
	}
	// END SIGN LOGIC
	return result
}

/*
 0 si a >= b
 Ignores the sign!
*/
func (ops *CipheredOperations) Gte(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	ops.gte(result, a, b, nbBits)
	return result
}

func (ops *CipheredOperations) gte(result, a, b []*LweSample, nbBits int) {
	eq := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	BootsCONSTANT(result[0], 0, ops.bk)
	for i := 0; i < nbBits; i++ {
		BootsXNOR(eq[0], a[i], b[i], ops.bk)
		BootsMUX(result[0], eq[0], result[0], a[i], ops.bk)
	}
}

// signed bit shift left
func (ops *CipheredOperations) ShiftLeft(a []*LweSample, positions, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	neg := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	is_neg := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	val := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	ops.negative(neg, a, nbBits)
	ops.isNegative(is_neg, a, nbBits)

	for i := 0; i < nbBits; i++ {
		BootsMUX(val[i], is_neg[0], neg[i], a[i], ops.bk)
	}

	for i := 0; i < nbBits; i++ {
		BootsCOPY(result[i], val[i], ops.bk)
	}

	for i := 0; i < positions; i++ {
		for j := 1; j < nbBits; j++ {
			BootsCOPY(aux[j], result[j-1], ops.bk)
		}

		BootsCONSTANT(aux[0], 0, ops.bk)

		for j := 0; j < nbBits; j++ {
			BootsCOPY(result[j], aux[j], ops.bk)
		}
	}

	ops.negative(aux, result, nbBits)
	for i := 0; i < nbBits; i++ {
		BootsMUX(result[i], is_neg[0], aux[i], result[i], ops.bk)
	}
	return result
}

// signed bit shift right
func (ops *CipheredOperations) ShiftRight(a []*LweSample, positions, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	is_neg := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	neg := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	val := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	ops.negative(neg, a, nbBits)
	ops.isNegative(is_neg, a, nbBits)

	for i := 0; i < nbBits; i++ {
		BootsMUX(val[i], is_neg[0], neg[i], a[i], ops.bk)
	}

	for i := 0; i < nbBits; i++ {
		BootsCOPY(result[i], val[i], ops.bk)
	}

	for i := 0; i < positions; i++ {

		for j := 0; j < nbBits-1; j++ {
			BootsCOPY(aux[j], result[j+1], ops.bk)
		}

		BootsCONSTANT(aux[nbBits-1], 0, ops.bk)

		for j := 0; j < nbBits; j++ {
			BootsCOPY(result[j], aux[j], ops.bk)
		}
	}

	ops.negative(aux, result, nbBits)
	for i := 0; i < nbBits; i++ {
		BootsMUX(result[i], is_neg[0], aux[i], result[i], ops.bk)
	}
	return result
}

// Unsigned shift left
func (ops *CipheredOperations) UshiftLeft(a []*LweSample, positions, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	for i := 0; i < nbBits; i++ {
		BootsCOPY(result[i], a[i], ops.bk)
	}

	for i := 0; i < positions; i++ {
		for j := 1; j < nbBits; j++ {
			BootsCOPY(aux[j], result[j-1], ops.bk)
		}

		BootsCONSTANT(aux[0], 0, ops.bk)

		for j := 0; j < nbBits; j++ {
			BootsCOPY(result[j], aux[j], ops.bk)
		}
	}
	return result
}

// unsigned shift right
func (ops *CipheredOperations) UshiftRight(a []*LweSample, positions, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	for i := 0; i < nbBits; i++ {
		BootsCOPY(result[i], a[i], ops.bk)
	}

	for i := 0; i < positions; i++ {
		for j := 0; j < nbBits-1; j++ {
			BootsCOPY(aux[j], result[j+1], ops.bk)
		}

		BootsCONSTANT(aux[nbBits-1], 0, ops.bk)

		for j := 0; j < nbBits; j++ {
			BootsCOPY(result[j], aux[j], ops.bk)
		}
	}
	return result
}

// Scaling from nb_bits to nb_bits_result
func (ops *CipheredOperations) urescale(result, a []*LweSample, nbBitsResult, nbBits int) {

	for i := 0; i < nbBitsResult; i++ {
		BootsCONSTANT(result[i], 0, ops.bk)
	}

	// determine if the sign should be taken into account
	bits := nbBits
	if nbBits > nbBitsResult {
		bits = nbBitsResult
	}
	for i := 0; i < bits; i++ {
		BootsCOPY(result[i], a[i], ops.bk)
	}
}

func (ops *CipheredOperations) rescale(result, a []*LweSample, nbBitsResult, nbBits int) {
	auxA := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	corrige := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	aux_res := NewGateBootstrappingCiphertextArray(nbBitsResult, ops.bk.params)
	aux_res_neg := NewGateBootstrappingCiphertextArray(nbBitsResult, ops.bk.params)

	ops.negative(auxA, a, nbBits)
	ops.isNegative(corrige, a, nbBits)
	// Trabajaremos con el positivo
	n := ops.Maximum(auxA, a, nbBits)

	ops.urescale(aux_res, n, nbBitsResult, nbBits)

	ops.negative(aux_res_neg, aux_res, nbBitsResult)
	for i := 0; i < nbBitsResult; i++ {
		BootsMUX(result[i], corrige[0], aux_res_neg[i], aux_res[i], ops.bk)
	}
}

func (ops *CipheredOperations) udiv(cociente, a, b []*LweSample, nbBits int) {
	gt := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	div_aux := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.params)
	div_aux2 := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.params)
	dividendo := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.params)
	divisor := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.params)
	remainder := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.params)

	ops.rescale(dividendo, a, 2*nbBits, nbBits)
	ops.rescale(div_aux, b, 2*nbBits, nbBits)
	divisor = ops.UshiftLeft(div_aux, nbBits-1, 2*nbBits)

	for i := 0; i < nbBits; i++ {
		// gt = dividend >= divisor
		ops.gte(gt, dividendo, divisor, 2*nbBits)

		BootsCOPY(cociente[nbBits-i-1], gt[0], ops.bk)

		// remainder = gt? sub(dividend, divisor) : remainder
		div_aux = ops.Sub(dividendo, divisor, 2*nbBits)
		// divisor = shiftr(divisor, 1)
		div_aux2 = ops.UshiftRight(divisor, 1, 2*nbBits)
		for j := 0; j < 2*nbBits; j++ {
			BootsMUX(remainder[j], gt[0], div_aux[j], dividendo[j], ops.bk)
			// dividendo = gt ? remainder : dividendo
			BootsMUX(dividendo[j], gt[0], remainder[j], dividendo[j], ops.bk)
			BootsCOPY(divisor[j], div_aux2[j], ops.bk)
		}
	}

}

func (ops *CipheredOperations) Div(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	aux2 := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	negatA := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	negatB := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	//opA := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	//opB := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	gt := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	bit := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	corrige := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	isNegativeA := NewGateBootstrappingCiphertextArray(2, ops.bk.params)
	isNegativeB := NewGateBootstrappingCiphertextArray(2, ops.bk.params)

	div_aux := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.params)
	div_aux2 := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.params)
	dividendo := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.params)
	divisor := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.params)
	cociente := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.params)
	resto := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.params)

	for i := 0; i < nbBits; i++ {
		BootsCONSTANT(aux[i], 0, ops.bk)
		BootsCONSTANT(aux2[i], 0, ops.bk)
		//BootsCONSTANT(opA[i], 0, ops.bk)
		//BootsCONSTANT(opB[i], 0, ops.bk)
	}

	for i := 0; i < 2*nbBits; i++ {
		BootsCONSTANT(dividendo[i], 0, ops.bk)
		BootsCONSTANT(div_aux[i], 0, ops.bk)
		BootsCONSTANT(div_aux2[i], 0, ops.bk)
		BootsCONSTANT(divisor[i], 0, ops.bk)
		BootsCONSTANT(cociente[i], 0, ops.bk)
		BootsCONSTANT(resto[i], 0, ops.bk)
	}

	for i := 0; i < 2; i++ {
		BootsCONSTANT(gt[i], 0, ops.bk)
		BootsCONSTANT(bit[i], 0, ops.bk)
	}

	// BEGIN LOGICAL SIGN
	ops.negative(negatA, a, nbBits)
	ops.negative(negatB, b, nbBits)

	// put the two numbers in positive
	opA := ops.Maximum(negatA, a, nbBits)
	opB := ops.Maximum(negatB, b, nbBits)

	// if only one of the two is negative, the result is negative
	ops.isNegative(isNegativeA, a, nbBits)
	ops.isNegative(isNegativeB, b, nbBits)
	BootsXOR(corrige[0], isNegativeA[0], isNegativeB[0], ops.bk)
	// END LOGICAL SIGN

	ops.udiv(result, opA, opB, nbBits)

	// BEGIN LOGICAL SIGN
	// determine whether to return the positive or negative result
	ops.negative(aux, result, nbBits)

	for i := 0; i < nbBits; i++ {
		BootsMUX(result[i], corrige[0], aux[i], result[i], ops.bk)
	}
	// END LOGICAL SIGN
	return result
}

func (ops *CipheredOperations) Pow(a []*LweSample, n, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	// aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	cero := NewGateBootstrappingCiphertextArray(1, ops.bk.params)
	BootsCONSTANT(cero[0], 0, ops.bk)

	// Initializing result
	for i := 0; i < nbBits; i++ {
		if n > 0 {
			BootsCOPY(result[i], a[i], ops.bk)
		} else {
			BootsCONSTANT(result[i], 0, ops.bk)
		}
	}

	if n <= 0 {
		BootsCONSTANT(result[0], 1, ops.bk)
	}

	for i := 0; i < n-1; i++ {
		aux := ops.Mul(result, a, nbBits)
		for j := 0; j < nbBits; j++ {
			BootsCOPY(result[j], aux[j], ops.bk)
		}
	}
	return result
}

// boolean operations wrappers
func (ops *CipheredOperations) Nand(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsNAND(result[i], a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Or(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsOR(result[i], a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) And(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsAND(result[i], a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Xor(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsXOR(result[i], a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Xnor(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsXNOR(result[i], a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Not(a []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsNOT(result[i], a[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Copy(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsCOPY(result[i], a[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Constant(value int64, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsCONSTANT(result[i], value, ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Nor(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsNOR(result[i], a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) AndNY(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsANDNY(result[i], a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) AndYN(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsANDYN(result[i], a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) OrNY(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsORNY(result[i], a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) OrYN(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsORYN(result[i], a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Mux(a, b, c []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	for i := 0; i < nbBits; i++ {
		BootsMUX(result[i], a[i], b[i], c[i], ops.bk)
	}
	return result
}
