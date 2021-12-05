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
	Constant(value bool, nbBits int) (result []*LweSample)
}

type CipheredOperations struct {
	bk *PublicKey
}

// elementary full comparator gate that is used to compare the i-th bit:
//   input: ai and bi the i-th bit of a and b
//          lsb_carry: the result of the comparison on the lowest bits
//   algo: if (a==b) return lsb_carry else return b
func (ops *CipheredOperations) CompareBit(a, b, lsbCarry, tmp *LweSample) *LweSample {
	result := NewLweSample(ops.bk.Params.InOutParams)
	tmp = Xnor(a, b, ops.bk)
	result = Mux(tmp, lsbCarry, a, ops.bk)
	return result
}

// Returns a == b
func (ops *CipheredOperations) Equals(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	result[0] = Constant(true, ops.bk)
	for i := 0; i < nbBits; i++ {
		tmps[0] = Xnor(a[i], b[i], ops.bk)
		result[0] = And(result[0], tmps[0], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) isNegative(a []*LweSample, nbBits int) *LweSample {
	return Copy(a[nbBits-1], ops.bk)
}

// this function compares two multibit words, and puts the min in result
func (ops *CipheredOperations) Minimum(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	aGreater := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	minimumSameSign := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	minimumOneNegative := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	oneNegative := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	negativoA := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	negativoB := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	negativoA[0] = ops.isNegative(a, nbBits)
	negativoB[0] = ops.isNegative(b, nbBits)

	oneNegative[0] = Xor(negativoA[0], negativoB[0], ops.bk)

	// a > b = soloOneNegative & is_negative(b)
	aGreater[0] = And(oneNegative[0], negativoB[0], ops.bk)
	for i := 0; i < nbBits; i++ {
		minimumOneNegative[i] = Mux(aGreater[0], b[i], a[i], ops.bk)
	}

	//initialize the carry to 0
	tmps[0] = Constant(false, ops.bk)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		tmps[0] = ops.CompareBit(a[i], b[i], tmps[0], tmps[1])
	}

	//tmps[0] is the result of the comparaison: 0 if a is larger, 1 if b is larger
	//select the max and copy it to the result
	for i := 0; i < nbBits; i++ {
		minimumSameSign[i] = Mux(tmps[0], b[i], a[i], ops.bk)
	}

	// Result depending on whether we compare the same sign or not
	for i := 0; i < nbBits; i++ {
		result[i] = Mux(oneNegative[0], minimumOneNegative[i], minimumSameSign[i], ops.bk)
	}
	return result
}

// this function compares two multibit words, and puts the min in result
func (ops *CipheredOperations) Maximum2(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	aGreater := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	minimumSameSign := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	minimumOneNegative := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	oneNegative := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	negativoA := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	negativoB := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	negativoA[0] = ops.isNegative(a, nbBits)
	negativoB[0] = ops.isNegative(b, nbBits)

	oneNegative[0] = Xor(negativoA[0], negativoB[0], ops.bk)

	// a > b = soloOneNegative & is_negative(b)
	aGreater[0] = And(oneNegative[0], negativoB[0], ops.bk)
	for i := 0; i < nbBits; i++ {
		minimumOneNegative[i] = Mux(aGreater[0], b[i], a[i], ops.bk)
	}

	//initialize the carry to 0
	tmps[0] = Constant(false, ops.bk)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		tmps[0] = ops.CompareBit(a[i], b[i], tmps[0], tmps[1])
	}

	//tmps[0] is the result of the comparaison: 0 if a is larger, 1 if b is larger
	//select the max and copy it to the result
	for i := 0; i < nbBits; i++ {
		minimumSameSign[i] = Mux(tmps[0], a[i], b[i], ops.bk)
	}

	// Result depending on whether we compare the same sign or not
	for i := 0; i < nbBits; i++ {
		result[i] = Mux(oneNegative[0], minimumOneNegative[i], minimumSameSign[i], ops.bk)
	}
	return result
}

// this function compares two multibit words, and puts the min in result
func (ops *CipheredOperations) Maximum(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	aGreater := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	minimumSameSign := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	minimumOneNegative := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	oneNegative := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	negativoA := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	negativoB := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	negativoA[0] = ops.isNegative(a, nbBits)
	negativoB[0] = ops.isNegative(b, nbBits)

	oneNegative[0] = Xor(negativoA[0], negativoB[0], ops.bk)

	// a > b = soloOneNegative & is_negative(b)
	aGreater[0] = And(oneNegative[0], negativoB[0], ops.bk)
	for i := 0; i < nbBits; i++ {
		minimumOneNegative[i] = Mux(aGreater[0], b[i], a[i], ops.bk)
	}

	//initialize the carry to 0
	tmps[0] = Constant(false, ops.bk)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		tmps[0] = ops.CompareBit(a[i], b[i], tmps[0], tmps[1])
	}

	//tmps[0] is the result of the comparaison: 0 if a is larger, 1 if b is larger
	//select the max and copy it to the result
	for i := 0; i < nbBits; i++ {
		minimumSameSign[i] = Mux(tmps[0], b[i], a[i], ops.bk)
	}

	// Todo - same as in minimum, but returning the opposite
	for i := 0; i < nbBits; i++ {
		//BootsMUX(result[i], oneNegative[0], minimumOneNegative[i], minimumMismoSigno[i], ops.bk)
		result[i] = Mux(oneNegative[0], minimumSameSign[i], minimumOneNegative[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) addBit(result, carry_out, a, b, carry_in *LweSample) {
	s1 := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	c1 := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	c2 := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	s1[0] = Constant(false, ops.bk)
	c1[0] = Constant(false, ops.bk)
	c2[0] = Constant(false, ops.bk)

	s1[0] = Xor(a, b, ops.bk)
	result = Xor(s1[0], carry_in, ops.bk)

	c1[0] = And(s1[0], carry_in, ops.bk)
	c2[0] = And(a, b, ops.bk)
	carry_out = Or(c1[0], c2[0], ops.bk)
}

// return -a
func (ops *CipheredOperations) negative(result, a []*LweSample, nbBits int) {

	ha_changed := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	not_x := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	for i := 0; i < 2; i++ {
		ha_changed[i] = Constant(false, ops.bk)
		not_x[i] = Constant(false, ops.bk)
	}

	for i := 0; i < nbBits; i++ {
		not_x[0] = Not(a[i], ops.bk)
		result[i] = Mux(ha_changed[0], not_x[0], a[i], ops.bk)
		ha_changed[0] = Or(ha_changed[0], a[i], ops.bk)
	}

}

func (ops *CipheredOperations) Add(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	tmpsCarry := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	//initialize the carry to 0
	tmpsCarry[0] = Constant(false, ops.bk)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		ops.addBit(result[i], tmpsCarry[0], a[i], b[i], tmpsCarry[0])
	}
	return result
}

func (ops *CipheredOperations) Sub(a, b []*LweSample, nbBits int) (result []*LweSample) {
	res := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	ops.negative(res, b, nbBits)
	return ops.Add(a, res, nbBits)
}

// Unsigned multiply
func (ops *CipheredOperations) umul(result, a, b []*LweSample, nbBits int) {
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	aux2 := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)

	for i := 0; i < nbBits; i++ {
		aux[i] = Constant(false, ops.bk)
		aux2[i] = Constant(false, ops.bk)
	}

	// Multiply opA * opB
	for i := 0; i < nbBits/2; i++ {
		// Reset the auxs
		for j := 0; j < nbBits; j++ {
			aux[j] = Constant(false, ops.bk)
			aux2[j] = Constant(false, ops.bk)
		}

		for j := 0; j < (nbBits/2)+1; j++ {
			aux[j+i] = And(a[i], b[j], ops.bk)
		}

		// add(aux2, aux, result, nbBits, bk);
		aux2 = ops.Add(aux, result, nbBits)
		//result = ops.Add(aux2, aux, nbBits)

		for j := 0; j < nbBits; j++ {
			result[j] = Copy(aux2[j], ops.bk)
		}

	}
}

// multiply two ciphertexts and return the result
func (ops *CipheredOperations) Mul(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	aux2 := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)

	// Parameters to take into account negative numbers
	negatA := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	negatB := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)

	opA := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	opB := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)

	isNegativeA := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	isNegativeB := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	corrige := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	// Set number of bits so: nb(result) = nb(a)+nb(b)
	for i := 0; i < nbBits; i++ {
		aux[i] = Constant(false, ops.bk)
		aux2[i] = Constant(false, ops.bk)
		negatA[i] = Constant(false, ops.bk)
		negatB[i] = Constant(false, ops.bk)
		opA[i] = Constant(false, ops.bk)
		opB[i] = Constant(false, ops.bk)
		result[i] = Constant(false, ops.bk)
	}

	for i := 0; i < 2; i++ {
		isNegativeA[i] = Constant(false, ops.bk)
		isNegativeB[i] = Constant(false, ops.bk)
		corrige[i] = Constant(false, ops.bk)
	}

	// BEGIN SIGN LOGIC
	ops.negative(negatA, a, nbBits)
	ops.negative(negatB, b, nbBits)

	// Put the two numbers in positive
	opA = ops.Maximum(negatA, a, nbBits)
	opB = ops.Maximum(negatB, b, nbBits)

	// If only one of the two is negative, the result is negative
	isNegativeA[0] = ops.isNegative(a, nbBits)
	isNegativeB[0] = ops.isNegative(b, nbBits)
	corrige[0] = Xor(isNegativeA[0], isNegativeB[0], ops.bk)
	// END SIGN LOGIC

	ops.umul(result, opA, opB, nbBits)

	// BEGIN SIGN LOGIC
	// We determine whether to return the positive or negative result
	ops.negative(aux, result, nbBits)

	for i := 0; i < nbBits; i++ {
		result[i] = Mux(corrige[0], aux[i], result[i], ops.bk)
	}
	// END SIGN LOGIC
	return result
}

//0 si a >= b
//Ignores the sign!
func (ops *CipheredOperations) Gte(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	ops.gte(result, a, b, nbBits)
	return result
}

func (ops *CipheredOperations) gte(result, a, b []*LweSample, nbBits int) {
	eq := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	result[0] = Constant(false, ops.bk)
	for i := 0; i < nbBits; i++ {
		eq[0] = Xnor(a[i], b[i], ops.bk)
		result[0] = Mux(eq[0], result[0], a[i], ops.bk)
	}
}

// signed bit shift left
func (ops *CipheredOperations) ShiftLeft(a []*LweSample, positions, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)

	neg := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	is_neg := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	val := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)

	ops.negative(neg, a, nbBits)
	is_neg[0] = ops.isNegative(a, nbBits)

	for i := 0; i < nbBits; i++ {
		val[i] = Mux(is_neg[0], neg[i], a[i], ops.bk)
	}

	for i := 0; i < nbBits; i++ {
		result[i] = Copy(val[i], ops.bk)
	}

	for i := 0; i < positions; i++ {
		for j := 1; j < nbBits; j++ {
			aux[j] = Copy(result[j-1], ops.bk)
		}

		aux[0] = Constant(false, ops.bk)

		for j := 0; j < nbBits; j++ {
			result[j] = Copy(aux[j], ops.bk)
		}
	}

	ops.negative(aux, result, nbBits)
	for i := 0; i < nbBits; i++ {
		result[i] = Mux(is_neg[0], aux[i], result[i], ops.bk)
	}
	return result
}

// signed bit shift right
func (ops *CipheredOperations) ShiftRight(a []*LweSample, positions, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)

	is_neg := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	neg := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	val := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)

	ops.negative(neg, a, nbBits)
	is_neg[0] = ops.isNegative(a, nbBits)

	for i := 0; i < nbBits; i++ {
		val[i] = Mux(is_neg[0], neg[i], a[i], ops.bk)
	}

	for i := 0; i < nbBits; i++ {
		result[i] = Copy(val[i], ops.bk)
	}

	for i := 0; i < positions; i++ {

		for j := 0; j < nbBits-1; j++ {
			aux[j] = Copy(result[j+1], ops.bk)
		}

		aux[nbBits-1] = Constant(false, ops.bk)

		for j := 0; j < nbBits; j++ {
			result[j] = Copy(aux[j], ops.bk)
		}
	}

	ops.negative(aux, result, nbBits)
	for i := 0; i < nbBits; i++ {
		result[i] = Mux(is_neg[0], aux[i], result[i], ops.bk)
	}
	return result
}

// Unsigned shift left
func (ops *CipheredOperations) UshiftLeft(a []*LweSample, positions, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)

	for i := 0; i < nbBits; i++ {
		result[i] = Copy(a[i], ops.bk)
	}

	for i := 0; i < positions; i++ {
		for j := 1; j < nbBits; j++ {
			aux[j] = Copy(result[j-1], ops.bk)
		}

		aux[0] = Constant(false, ops.bk)

		for j := 0; j < nbBits; j++ {
			result[j] = Copy(aux[j], ops.bk)
		}
	}
	return result
}

// unsigned shift right
func (ops *CipheredOperations) UshiftRight(a []*LweSample, positions, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)

	for i := 0; i < nbBits; i++ {
		result[i] = Copy(a[i], ops.bk)
	}

	for i := 0; i < positions; i++ {
		for j := 0; j < nbBits-1; j++ {
			aux[j] = Copy(result[j+1], ops.bk)
		}

		aux[nbBits-1] = Constant(false, ops.bk)

		for j := 0; j < nbBits; j++ {
			result[j] = Copy(aux[j], ops.bk)
		}
	}
	return result
}

// Scaling from nb_bits to nb_bits_result
func (ops *CipheredOperations) urescale(result, a []*LweSample, nbBitsResult, nbBits int) {

	for i := 0; i < nbBitsResult; i++ {
		result[i] = Constant(false, ops.bk)
	}

	// determine if the sign should be taken into account
	bits := nbBits
	if nbBits > nbBitsResult {
		bits = nbBitsResult
	}
	for i := 0; i < bits; i++ {
		result[i] = Copy(a[i], ops.bk)
	}
}

func (ops *CipheredOperations) rescale(result, a []*LweSample, nbBitsResult, nbBits int) {
	auxA := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	corrige := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	aux_res := NewGateBootstrappingCiphertextArray(nbBitsResult, ops.bk.Params)
	aux_res_neg := NewGateBootstrappingCiphertextArray(nbBitsResult, ops.bk.Params)

	ops.negative(auxA, a, nbBits)
	corrige[0] = ops.isNegative(a, nbBits)
	// Trabajaremos con el positivo
	n := ops.Maximum(auxA, a, nbBits)

	ops.urescale(aux_res, n, nbBitsResult, nbBits)

	ops.negative(aux_res_neg, aux_res, nbBitsResult)
	for i := 0; i < nbBitsResult; i++ {
		result[i] = Mux(corrige[0], aux_res_neg[i], aux_res[i], ops.bk)
	}
}

func (ops *CipheredOperations) udiv(cociente, a, b []*LweSample, nbBits int) {
	gt := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	div_aux := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.Params)
	div_aux2 := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.Params)
	dividendo := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.Params)
	divisor := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.Params)
	remainder := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.Params)

	ops.rescale(dividendo, a, 2*nbBits, nbBits)
	ops.rescale(div_aux, b, 2*nbBits, nbBits)
	divisor = ops.UshiftLeft(div_aux, nbBits-1, 2*nbBits)

	for i := 0; i < nbBits; i++ {
		// gt = dividend >= divisor
		ops.gte(gt, dividendo, divisor, 2*nbBits)

		cociente[nbBits-i-1] = Copy(gt[0], ops.bk)

		// remainder = gt? sub(dividend, divisor) : remainder
		div_aux = ops.Sub(dividendo, divisor, 2*nbBits)
		// divisor = shiftr(divisor, 1)
		div_aux2 = ops.UshiftRight(divisor, 1, 2*nbBits)
		for j := 0; j < 2*nbBits; j++ {
			remainder[j] = Mux(gt[0], div_aux[j], dividendo[j], ops.bk)
			// dividendo = gt ? remainder : dividendo
			dividendo[j] = Mux(gt[0], remainder[j], dividendo[j], ops.bk)
			divisor[j] = Copy(div_aux2[j], ops.bk)
		}
	}

}

func (ops *CipheredOperations) Div(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)

	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	aux2 := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	negatA := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	negatB := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	//opA := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	//opB := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	gt := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	bit := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	corrige := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	isNegativeA := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)
	isNegativeB := NewGateBootstrappingCiphertextArray(2, ops.bk.Params)

	div_aux := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.Params)
	div_aux2 := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.Params)
	dividendo := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.Params)
	divisor := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.Params)
	cociente := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.Params)
	resto := NewGateBootstrappingCiphertextArray(2*nbBits, ops.bk.Params)

	for i := 0; i < nbBits; i++ {
		aux[i] = Constant(false, ops.bk)
		aux2[i] = Constant(false, ops.bk)
		//BootsCONSTANT(opA[i], 0, ops.bk)
		//BootsCONSTANT(opB[i], 0, ops.bk)
	}

	for i := 0; i < 2*nbBits; i++ {
		dividendo[i] = Constant(false, ops.bk)
		div_aux[i] = Constant(false, ops.bk)
		div_aux2[i] = Constant(false, ops.bk)
		divisor[i] = Constant(false, ops.bk)
		cociente[i] = Constant(false, ops.bk)
		resto[i] = Constant(false, ops.bk)
	}

	for i := 0; i < 2; i++ {
		gt[i] = Constant(false, ops.bk)
		bit[i] = Constant(false, ops.bk)
	}

	// BEGIN LOGICAL SIGN
	ops.negative(negatA, a, nbBits)
	ops.negative(negatB, b, nbBits)

	// put the two numbers in positive
	opA := ops.Maximum(negatA, a, nbBits)
	opB := ops.Maximum(negatB, b, nbBits)

	// if only one of the two is negative, the result is negative
	isNegativeA[0] = ops.isNegative(a, nbBits)
	isNegativeB[0] = ops.isNegative(b, nbBits)
	corrige[0] = Xor(isNegativeA[0], isNegativeB[0], ops.bk)
	// END LOGICAL SIGN

	ops.udiv(result, opA, opB, nbBits)

	// BEGIN LOGICAL SIGN
	// determine whether to return the positive or negative result
	ops.negative(aux, result, nbBits)

	for i := 0; i < nbBits; i++ {
		result[i] = Mux(corrige[0], aux[i], result[i], ops.bk)
	}
	// END LOGICAL SIGN
	return result
}

func (ops *CipheredOperations) Pow(a []*LweSample, n, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	// aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	cero := NewGateBootstrappingCiphertextArray(1, ops.bk.Params)
	cero[0] = Constant(false, ops.bk)

	// Initializing result
	for i := 0; i < nbBits; i++ {
		if n > 0 {
			result[i] = Copy(a[i], ops.bk)
		} else {
			result[i] = Constant(false, ops.bk)
		}
	}

	if n <= 0 {
		result[0] = Constant(false, ops.bk)
	}

	for i := 0; i < n-1; i++ {
		aux := ops.Mul(result, a, nbBits)
		for j := 0; j < nbBits; j++ {
			result[j] = Copy(aux[j], ops.bk)
		}
	}
	return result
}

// boolean operations wrappers
func (ops *CipheredOperations) Nand(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = Nand(a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Or(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = Or(a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) And(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = And(a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Xor(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = Xor(a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Xnor(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = Xnor(a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Not(a []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = Not(a[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Copy(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = Copy(a[i], ops.bk)
	}
	return result
}

	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {

func (ops *CipheredOperations) Nor(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = Nor(a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) AndNY(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = AndNY(a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) AndYN(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = AndYN(a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) OrNY(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = OrNY(a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) OrYN(a, b []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = OrYN(a[i], b[i], ops.bk)
	}
	return result
}

func (ops *CipheredOperations) Mux(a, b, c []*LweSample, nbBits int) (result []*LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.bk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = Mux(a[i], b[i], c[i], ops.bk)
	}
	return result
}
