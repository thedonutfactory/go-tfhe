package gates

import (
	"github.com/thedonutfactory/go-tfhe/core"
)

type Operations interface {
	// comparison
	CompareBit(a, b, lsbCarry, tmp *core.LweSample) (result *core.LweSample)
	Equals(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Minimum(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Maximum(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Gte(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)

	// arithmetic
	Add(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Sub(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Mul(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Div(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Pow(a []*core.LweSample, n, nbBits int) (result []*core.LweSample)

	// bitwise shift
	ShiftLeft(a []*core.LweSample, positions, nbBits int) (result []*core.LweSample)
	ShiftRight(a []*core.LweSample, positions, nbBits int) (result []*core.LweSample)
	UshiftLeft(a []*core.LweSample, positions, nbBits int) (result []*core.LweSample)
	UshiftRight(a []*core.LweSample, positions, nbBits int) (result []*core.LweSample)

	// bitwise operations
	Nand(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Or(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	And(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Xor(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Xnor(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Not(a []*core.LweSample, nbBits int) (result []*core.LweSample)
	Nor(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	AndNY(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	AndYN(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	OrNY(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	OrYN(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Mux(a, b, c []*core.LweSample, nbBits int) (result []*core.LweSample)

	// misc
	Copy(a, b []*core.LweSample, nbBits int) (result []*core.LweSample)
	Constant(value bool, nbBits int) (result []*core.LweSample)
}

type CipheredOperations struct {
	Pk *PublicKey
}

// elementary full comparator gate that is used to compare the i-th bit:
//   input: ai ops.bk.And bi the i-th bit of a ops.bk.And b
//          lsb_carry: the result of the comparison on the lowest bits
//   algo: if (a==b) return lsb_carry else return b
func (ops *CipheredOperations) CompareBit(a, b, lsbCarry, tmp *core.LweSample) *core.LweSample {
	result := core.NewLweSample(ops.Pk.Params.InOutParams)
	tmp = ops.Pk.Xnor(a, b)
	result = ops.Pk.Mux(tmp, lsbCarry, a)
	return result
}

// Returns a == b
func (ops *CipheredOperations) Equals(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	result[0] = ops.Pk.Constant(true)
	for i := 0; i < nbBits; i++ {
		tmps[0] = ops.Pk.Xnor(a[i], b[i])
		result[0] = ops.Pk.And(result[0], tmps[0])
	}
	return result
}

func (ops *CipheredOperations) isNegative(a []*core.LweSample, nbBits int) *core.LweSample {
	return ops.Pk.Copy(a[nbBits-1])
}

// this function compares two multibit words, and puts the min in result
func (ops *CipheredOperations) Minimum(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	aGreater := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	minimumSameSign := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	minimumOneNegative := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	oneNegative := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	negativoA := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	negativoB := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	negativoA[0] = ops.isNegative(a, nbBits)
	negativoB[0] = ops.isNegative(b, nbBits)

	oneNegative[0] = ops.Pk.Xor(negativoA[0], negativoB[0])

	// a > b = soloOneNegative & is_negative(b)
	aGreater[0] = ops.Pk.And(oneNegative[0], negativoB[0])
	for i := 0; i < nbBits; i++ {
		minimumOneNegative[i] = ops.Pk.Mux(aGreater[0], b[i], a[i])
	}

	//initialize the carry to 0
	tmps[0] = ops.Pk.Constant(false)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		tmps[0] = ops.CompareBit(a[i], b[i], tmps[0], tmps[1])
	}

	//tmps[0] is the result of the comparaison: 0 if a is larger, 1 if b is larger
	//select the max and ops.bk.Copy it to the result
	for i := 0; i < nbBits; i++ {
		minimumSameSign[i] = ops.Pk.Mux(tmps[0], b[i], a[i])
	}

	// Result depending on whether we compare the same sign or not
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Mux(oneNegative[0], minimumOneNegative[i], minimumSameSign[i])
	}
	return result
}

// this function compares two multibit words, and puts the min in result
func (ops *CipheredOperations) Maximum2(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	aGreater := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	minimumSameSign := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	minimumOneNegative := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	oneNegative := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	negativoA := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	negativoB := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	negativoA[0] = ops.isNegative(a, nbBits)
	negativoB[0] = ops.isNegative(b, nbBits)

	oneNegative[0] = ops.Pk.Xor(negativoA[0], negativoB[0])

	// a > b = soloOneNegative & is_negative(b)
	aGreater[0] = ops.Pk.And(oneNegative[0], negativoB[0])
	for i := 0; i < nbBits; i++ {
		minimumOneNegative[i] = ops.Pk.Mux(aGreater[0], b[i], a[i])
	}

	//initialize the carry to 0
	tmps[0] = ops.Pk.Constant(false)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		tmps[0] = ops.CompareBit(a[i], b[i], tmps[0], tmps[1])
	}

	//tmps[0] is the result of the comparaison: 0 if a is larger, 1 if b is larger
	//select the max and ops.bk.Copy it to the result
	for i := 0; i < nbBits; i++ {
		minimumSameSign[i] = ops.Pk.Mux(tmps[0], a[i], b[i])
	}

	// Result depending on whether we compare the same sign or not
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Mux(oneNegative[0], minimumOneNegative[i], minimumSameSign[i])
	}
	return result
}

// this function compares two multibit words, and puts the min in result
func (ops *CipheredOperations) Maximum(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	tmps := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	aGreater := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	minimumSameSign := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	minimumOneNegative := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	oneNegative := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	negativoA := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	negativoB := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	negativoA[0] = ops.isNegative(a, nbBits)
	negativoB[0] = ops.isNegative(b, nbBits)

	oneNegative[0] = ops.Pk.Xor(negativoA[0], negativoB[0])

	// a > b = soloOneNegative & is_negative(b)
	aGreater[0] = ops.Pk.And(oneNegative[0], negativoB[0])
	for i := 0; i < nbBits; i++ {
		minimumOneNegative[i] = ops.Pk.Mux(aGreater[0], b[i], a[i])
	}

	//initialize the carry to 0
	tmps[0] = ops.Pk.Constant(false)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		tmps[0] = ops.CompareBit(a[i], b[i], tmps[0], tmps[1])
	}

	//tmps[0] is the result of the comparaison: 0 if a is larger, 1 if b is larger
	//select the max and ops.bk.Copy it to the result
	for i := 0; i < nbBits; i++ {
		minimumSameSign[i] = ops.Pk.Mux(tmps[0], b[i], a[i])
	}

	// Todo - same as in minimum, but returning the opposite
	for i := 0; i < nbBits; i++ {
		//BootsMUX(result[i], oneNegative[0], minimumOneNegative[i], minimumMismoSigno[i])
		result[i] = ops.Pk.Mux(oneNegative[0], minimumSameSign[i], minimumOneNegative[i])
	}
	return result
}

func (ops *CipheredOperations) addBit(a, b, carry_in *core.LweSample) (result, carry_out *core.LweSample) {

	s1 := ops.Constant(false, 2)
	c1 := ops.Constant(false, 2)
	c2 := ops.Constant(false, 2)

	//s1[0] = ops.Pk.Constant(false)
	//c1[0] = ops.Pk.Constant(false)
	//c2[0] = ops.Pk.Constant(false)

	s1[0] = ops.Pk.Xor(a, b)
	result = ops.Pk.Xor(s1[0], carry_in)

	c1[0] = ops.Pk.And(s1[0], carry_in)
	c2[0] = ops.Pk.And(a, b)
	carry_out = ops.Pk.Or(c1[0], c2[0])

	return result, carry_out
}

// return -a
func (ops *CipheredOperations) negative(result, a []*core.LweSample, nbBits int) {

	ha_changed := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	not_x := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	for i := 0; i < 2; i++ {
		ha_changed[i] = ops.Pk.Constant(false)
		not_x[i] = ops.Pk.Constant(false)
	}

	for i := 0; i < nbBits; i++ {
		not_x[0] = ops.Pk.Not(a[i])
		result[i] = ops.Pk.Mux(ha_changed[0], not_x[0], a[i])
		ha_changed[0] = ops.Pk.Or(ha_changed[0], a[i])
	}

}

func (ops *CipheredOperations) Add(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	tmpsCarry := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	//initialize the carry to 0
	//tmpsCarry := ops.Constant(false, nbBits)
	tmpsCarry[0] = ops.Pk.Constant(false)

	//run the elementary comparator gate n times
	for i := 0; i < nbBits; i++ {
		result[i], tmpsCarry[0] = ops.addBit(a[i], b[i], tmpsCarry[0])
	}
	return result
}

func (ops *CipheredOperations) Sub(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	res := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	ops.negative(res, b, nbBits)
	return ops.Add(a, res, nbBits)
}

// Unsigned multiply
func (ops *CipheredOperations) umul(result, a, b []*core.LweSample, nbBits int) {
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	aux2 := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)

	for i := 0; i < nbBits; i++ {
		aux[i] = ops.Pk.Constant(false)
		aux2[i] = ops.Pk.Constant(false)
	}

	// Multiply opA * opB
	for i := 0; i < nbBits/2; i++ {
		// Reset the auxs
		for j := 0; j < nbBits; j++ {
			aux[j] = ops.Pk.Constant(false)
			aux2[j] = ops.Pk.Constant(false)
		}

		for j := 0; j < (nbBits/2)+1; j++ {
			aux[j+i] = ops.Pk.And(a[i], b[j])
		}

		// add(aux2, aux, result, nbBits, bk);
		aux2 = ops.Add(aux, result, nbBits)
		//result = ops.Add(aux2, aux, nbBits)

		for j := 0; j < nbBits; j++ {
			result[j] = ops.Pk.Copy(aux2[j])
		}

	}
}

// multiply two ciphertexts and return the result
func (ops *CipheredOperations) Mul(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	aux2 := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)

	// Parameters to take into account negative numbers
	negatA := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	negatB := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)

	opA := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	opB := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)

	isNegativeA := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	isNegativeB := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	corrige := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	// Set number of bits so: nb(result) = nb(a)+nb(b)
	for i := 0; i < nbBits; i++ {
		aux[i] = ops.Pk.Constant(false)
		aux2[i] = ops.Pk.Constant(false)
		negatA[i] = ops.Pk.Constant(false)
		negatB[i] = ops.Pk.Constant(false)
		opA[i] = ops.Pk.Constant(false)
		opB[i] = ops.Pk.Constant(false)
		result[i] = ops.Pk.Constant(false)
	}

	for i := 0; i < 2; i++ {
		isNegativeA[i] = ops.Pk.Constant(false)
		isNegativeB[i] = ops.Pk.Constant(false)
		corrige[i] = ops.Pk.Constant(false)
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
	corrige[0] = ops.Pk.Xor(isNegativeA[0], isNegativeB[0])
	// END SIGN LOGIC

	ops.umul(result, opA, opB, nbBits)

	// BEGIN SIGN LOGIC
	// We determine whether to return the positive or negative result
	ops.negative(aux, result, nbBits)

	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Mux(corrige[0], aux[i], result[i])
	}
	// END SIGN LOGIC
	return result
}

//0 si a >= b
//Ignores the sign!
func (ops *CipheredOperations) Gte(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	ops.gte(result, a, b, nbBits)
	return result
}

func (ops *CipheredOperations) gte(result, a, b []*core.LweSample, nbBits int) {
	eq := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	result[0] = ops.Pk.Constant(false)
	for i := 0; i < nbBits; i++ {
		eq[0] = ops.Pk.Xnor(a[i], b[i])
		result[0] = ops.Pk.Mux(eq[0], result[0], a[i])
	}
}

// signed bit shift left
func (ops *CipheredOperations) ShiftLeft(a []*core.LweSample, positions, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)

	neg := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	is_neg := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	val := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)

	ops.negative(neg, a, nbBits)
	is_neg[0] = ops.isNegative(a, nbBits)

	for i := 0; i < nbBits; i++ {
		val[i] = ops.Pk.Mux(is_neg[0], neg[i], a[i])
	}

	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Copy(val[i])
	}

	for i := 0; i < positions; i++ {
		for j := 1; j < nbBits; j++ {
			aux[j] = ops.Pk.Copy(result[j-1])
		}

		aux[0] = ops.Pk.Constant(false)

		for j := 0; j < nbBits; j++ {
			result[j] = ops.Pk.Copy(aux[j])
		}
	}

	ops.negative(aux, result, nbBits)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Mux(is_neg[0], aux[i], result[i])
	}
	return result
}

// signed bit shift right
func (ops *CipheredOperations) ShiftRight(a []*core.LweSample, positions, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)

	is_neg := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	neg := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	val := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)

	ops.negative(neg, a, nbBits)
	is_neg[0] = ops.isNegative(a, nbBits)

	for i := 0; i < nbBits; i++ {
		val[i] = ops.Pk.Mux(is_neg[0], neg[i], a[i])
	}

	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Copy(val[i])
	}

	for i := 0; i < positions; i++ {

		for j := 0; j < nbBits-1; j++ {
			aux[j] = ops.Pk.Copy(result[j+1])
		}

		aux[nbBits-1] = ops.Pk.Constant(false)

		for j := 0; j < nbBits; j++ {
			result[j] = ops.Pk.Copy(aux[j])
		}
	}

	ops.negative(aux, result, nbBits)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Mux(is_neg[0], aux[i], result[i])
	}
	return result
}

// Unsigned shift left
func (ops *CipheredOperations) UshiftLeft(a []*core.LweSample, positions, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)

	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Copy(a[i])
	}

	for i := 0; i < positions; i++ {
		for j := 1; j < nbBits; j++ {
			aux[j] = ops.Pk.Copy(result[j-1])
		}

		aux[0] = ops.Pk.Constant(false)

		for j := 0; j < nbBits; j++ {
			result[j] = ops.Pk.Copy(aux[j])
		}
	}
	return result
}

// unsigned shift right
func (ops *CipheredOperations) UshiftRight(a []*core.LweSample, positions, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)

	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Copy(a[i])
	}

	for i := 0; i < positions; i++ {
		for j := 0; j < nbBits-1; j++ {
			aux[j] = ops.Pk.Copy(result[j+1])
		}

		aux[nbBits-1] = ops.Pk.Constant(false)

		for j := 0; j < nbBits; j++ {
			result[j] = ops.Pk.Copy(aux[j])
		}
	}
	return result
}

// Scaling from nb_bits to nb_bits_result
func (ops *CipheredOperations) urescale(result, a []*core.LweSample, nbBitsResult, nbBits int) {

	for i := 0; i < nbBitsResult; i++ {
		result[i] = ops.Pk.Constant(false)
	}

	// determine if the sign should be taken into account
	bits := nbBits
	if nbBits > nbBitsResult {
		bits = nbBitsResult
	}
	for i := 0; i < bits; i++ {
		result[i] = ops.Pk.Copy(a[i])
	}
}

func (ops *CipheredOperations) rescale(result, a []*core.LweSample, nbBitsResult, nbBits int) {
	auxA := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	corrige := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	aux_res := NewGateBootstrappingCiphertextArray(nbBitsResult, ops.Pk.Params)
	aux_res_neg := NewGateBootstrappingCiphertextArray(nbBitsResult, ops.Pk.Params)

	ops.negative(auxA, a, nbBits)
	corrige[0] = ops.isNegative(a, nbBits)
	// Trabajaremos con el positivo
	n := ops.Maximum(auxA, a, nbBits)

	ops.urescale(aux_res, n, nbBitsResult, nbBits)

	ops.negative(aux_res_neg, aux_res, nbBitsResult)
	for i := 0; i < nbBitsResult; i++ {
		result[i] = ops.Pk.Mux(corrige[0], aux_res_neg[i], aux_res[i])
	}
}

func (ops *CipheredOperations) udiv(cociente, a, b []*core.LweSample, nbBits int) {
	gt := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	div_aux := NewGateBootstrappingCiphertextArray(2*nbBits, ops.Pk.Params)
	div_aux2 := NewGateBootstrappingCiphertextArray(2*nbBits, ops.Pk.Params)
	dividendo := NewGateBootstrappingCiphertextArray(2*nbBits, ops.Pk.Params)
	divisor := NewGateBootstrappingCiphertextArray(2*nbBits, ops.Pk.Params)
	remainder := NewGateBootstrappingCiphertextArray(2*nbBits, ops.Pk.Params)

	ops.rescale(dividendo, a, 2*nbBits, nbBits)
	ops.rescale(div_aux, b, 2*nbBits, nbBits)
	divisor = ops.UshiftLeft(div_aux, nbBits-1, 2*nbBits)

	for i := 0; i < nbBits; i++ {
		// gt = dividend >= divisor
		ops.gte(gt, dividendo, divisor, 2*nbBits)

		cociente[nbBits-i-1] = ops.Pk.Copy(gt[0])

		// remainder = gt? sub(dividend, divisor) : remainder
		div_aux = ops.Sub(dividendo, divisor, 2*nbBits)
		// divisor = shiftr(divisor, 1)
		div_aux2 = ops.UshiftRight(divisor, 1, 2*nbBits)
		for j := 0; j < 2*nbBits; j++ {
			remainder[j] = ops.Pk.Mux(gt[0], div_aux[j], dividendo[j])
			// dividendo = gt ? remainder : dividendo
			dividendo[j] = ops.Pk.Mux(gt[0], remainder[j], dividendo[j])
			divisor[j] = ops.Pk.Copy(div_aux2[j])
		}
	}

}

func (ops *CipheredOperations) Div(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)

	aux := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	aux2 := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	negatA := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	negatB := NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	//opA := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)
	//opB := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	gt := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	bit := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	corrige := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	isNegativeA := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)
	isNegativeB := NewGateBootstrappingCiphertextArray(2, ops.Pk.Params)

	div_aux := NewGateBootstrappingCiphertextArray(2*nbBits, ops.Pk.Params)
	div_aux2 := NewGateBootstrappingCiphertextArray(2*nbBits, ops.Pk.Params)
	dividendo := NewGateBootstrappingCiphertextArray(2*nbBits, ops.Pk.Params)
	divisor := NewGateBootstrappingCiphertextArray(2*nbBits, ops.Pk.Params)
	cociente := NewGateBootstrappingCiphertextArray(2*nbBits, ops.Pk.Params)
	resto := NewGateBootstrappingCiphertextArray(2*nbBits, ops.Pk.Params)

	for i := 0; i < nbBits; i++ {
		aux[i] = ops.Pk.Constant(false)
		aux2[i] = ops.Pk.Constant(false)
		//BootsCONSTANT(opA[i], 0)
		//BootsCONSTANT(opB[i], 0)
	}

	for i := 0; i < 2*nbBits; i++ {
		dividendo[i] = ops.Pk.Constant(false)
		div_aux[i] = ops.Pk.Constant(false)
		div_aux2[i] = ops.Pk.Constant(false)
		divisor[i] = ops.Pk.Constant(false)
		cociente[i] = ops.Pk.Constant(false)
		resto[i] = ops.Pk.Constant(false)
	}

	for i := 0; i < 2; i++ {
		gt[i] = ops.Pk.Constant(false)
		bit[i] = ops.Pk.Constant(false)
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
	corrige[0] = ops.Pk.Xor(isNegativeA[0], isNegativeB[0])
	// END LOGICAL SIGN

	ops.udiv(result, opA, opB, nbBits)

	// BEGIN LOGICAL SIGN
	// determine whether to return the positive or negative result
	ops.negative(aux, result, nbBits)

	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Mux(corrige[0], aux[i], result[i])
	}
	// END LOGICAL SIGN
	return result
}

func (ops *CipheredOperations) Pow(a []*core.LweSample, n, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	// aux := NewGateBootstrappingCiphertextArray(nbBits, ops.bk.params)

	cero := NewGateBootstrappingCiphertextArray(1, ops.Pk.Params)
	cero[0] = ops.Pk.Constant(false)

	// Initializing result
	for i := 0; i < nbBits; i++ {
		if n > 0 {
			result[i] = ops.Pk.Copy(a[i])
		} else {
			result[i] = ops.Pk.Constant(false)
		}
	}

	if n <= 0 {
		result[0] = ops.Pk.Constant(false)
	}

	for i := 0; i < n-1; i++ {
		aux := ops.Mul(result, a, nbBits)
		for j := 0; j < nbBits; j++ {
			result[j] = ops.Pk.Copy(aux[j])
		}
	}
	return result
}

// boolean operations wrappers
func (ops *CipheredOperations) Nand(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Nand(a[i], b[i])
	}
	return result
}

func (ops *CipheredOperations) Or(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Or(a[i], b[i])
	}
	return result
}

func (ops *CipheredOperations) And(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.And(a[i], b[i])
	}
	return result
}

func (ops *CipheredOperations) Xor(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Xor(a[i], b[i])
	}
	return result
}

func (ops *CipheredOperations) Xnor(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Xnor(a[i], b[i])
	}
	return result
}

func (ops *CipheredOperations) Not(a []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Not(a[i])
	}
	return result
}

func (ops *CipheredOperations) Copy(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Copy(a[i])
	}
	return result
}

func (ops *CipheredOperations) Constant(value bool, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Constant(value)
	}
	return result
}

func (ops *CipheredOperations) Nor(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Nor(a[i], b[i])
	}
	return result
}

func (ops *CipheredOperations) AndNY(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.AndNY(a[i], b[i])
	}
	return result
}

func (ops *CipheredOperations) AndYN(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.AndYN(a[i], b[i])
	}
	return result
}

func (ops *CipheredOperations) OrNY(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.OrNY(a[i], b[i])
	}
	return result
}

func (ops *CipheredOperations) OrYN(a, b []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.OrYN(a[i], b[i])
	}
	return result
}

func (ops *CipheredOperations) Mux(a, b, c []*core.LweSample, nbBits int) (result []*core.LweSample) {
	result = NewGateBootstrappingCiphertextArray(nbBits, ops.Pk.Params)
	for i := 0; i < nbBits; i++ {
		result[i] = ops.Pk.Mux(a[i], b[i], c[i])
	}
	return result
}
