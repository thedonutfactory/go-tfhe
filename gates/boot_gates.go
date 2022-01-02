package gates

import (
	"github.com/thedonutfactory/go-tfhe/core"
	"github.com/thedonutfactory/go-tfhe/types"
)

/*
 * Homomorphic bootstrapped NAND gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) Nand(ca, cb *core.LweSample) *core.LweSample {
	MU := types.ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := core.NewLweSample(inOutParams)

	//compute: (0,1/8) - ca - cb
	NandConst := types.ModSwitchToTorus32(1, 8)
	core.LweNoiselessTrivial(tempResult, NandConst, inOutParams)
	core.LweSubTo(tempResult, ca, inOutParams)
	core.LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//return tfheBootstrap(bk.Bkw.Bk, MU, tempResult)
	return core.TfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
	// tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped OR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) Or(ca, cb *core.LweSample) *core.LweSample {
	MU := types.ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := core.NewLweSample(inOutParams)

	//compute: (0,1/8) + ca + cb
	OrConst := types.ModSwitchToTorus32(1, 8)
	core.LweNoiselessTrivial(tempResult, OrConst, inOutParams)
	core.LweAddTo(tempResult, ca, inOutParams)
	core.LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return core.TfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped AND gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) And(ca, cb *core.LweSample) *core.LweSample {
	MU := types.ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := core.NewLweSample(inOutParams)

	//compute: (0,-1/8) + ca + cb
	AndConst := types.ModSwitchToTorus32(-1, 8)
	core.LweNoiselessTrivial(tempResult, AndConst, inOutParams)
	core.LweAddTo(tempResult, ca, inOutParams)
	core.LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return core.TfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped XOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) Xor(ca, cb *core.LweSample) *core.LweSample {
	MU := types.ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := core.NewLweSample(inOutParams)

	//compute: (0,1/4) + 2*(ca + cb)
	XorConst := types.ModSwitchToTorus32(1, 4)
	core.LweNoiselessTrivial(tempResult, XorConst, inOutParams)
	core.LweAddMulTo(tempResult, 2, ca, inOutParams)
	core.LweAddMulTo(tempResult, 2, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return core.TfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped XNOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) Xnor(ca, cb *core.LweSample) *core.LweSample {
	MU := types.ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := core.NewLweSample(inOutParams)

	//compute: (0,-1/4) + 2*(-ca-cb)
	XnorConst := types.ModSwitchToTorus32(-1, 4)
	core.LweNoiselessTrivial(tempResult, XnorConst, inOutParams)
	core.LweSubMulTo(tempResult, 2, ca, inOutParams)
	core.LweSubMulTo(tempResult, 2, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return core.TfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped NOT gate (doesn't need to be bootstrapped)
 * Takes in input 1 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) Not(ca *core.LweSample) *core.LweSample {
	inOutParams := bk.Params.InOutParams
	result := core.NewLweSample(inOutParams)
	core.LweNegate(result, ca, inOutParams)
	return result
}

/*
 * Homomorphic bootstrapped COPY gate (doesn't need to be bootstrapped)
 * Takes in input 1 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) Copy(ca *core.LweSample) *core.LweSample {
	inOutParams := bk.Params.InOutParams
	result := core.NewLweSample(inOutParams)
	core.LweCopy(result, ca, inOutParams)
	return result
}

/*
 * Homomorphic Trivial Constant gate (doesn't need to be bootstrapped)
 * Takes a boolean value)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) Constant(value bool) *core.LweSample {
	inOutParams := bk.Params.InOutParams
	MU := types.ModSwitchToTorus32(1, 8)
	var muValue = -MU
	if value {
		muValue = MU
	}
	result := core.NewLweSample(inOutParams)
	core.LweNoiselessTrivial(result, muValue, inOutParams)
	return result
}

/*
 * Homomorphic bootstrapped NOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) Nor(ca, cb *core.LweSample) *core.LweSample {
	MU := types.ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := core.NewLweSample(inOutParams)

	//compute: (0,-1/8) - ca - cb
	NorConst := types.ModSwitchToTorus32(-1, 8)
	core.LweNoiselessTrivial(tempResult, NorConst, inOutParams)
	core.LweSubTo(tempResult, ca, inOutParams)
	core.LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return core.TfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped AndNY Gate: not(a) and b
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) AndNY(ca, cb *core.LweSample) *core.LweSample {
	MU := types.ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := core.NewLweSample(inOutParams)

	//compute: (0,-1/8) - ca + cb
	AndNYConst := types.ModSwitchToTorus32(-1, 8)
	core.LweNoiselessTrivial(tempResult, AndNYConst, inOutParams)
	core.LweSubTo(tempResult, ca, inOutParams)
	core.LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return core.TfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped AndYN Gate: a and not(b)
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) AndYN(ca, cb *core.LweSample) *core.LweSample {
	MU := types.ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := core.NewLweSample(inOutParams)

	//compute: (0,-1/8) + ca - cb
	AndYNConst := types.ModSwitchToTorus32(-1, 8)
	core.LweNoiselessTrivial(tempResult, AndYNConst, inOutParams)
	core.LweAddTo(tempResult, ca, inOutParams)
	core.LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return core.TfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped OrNY Gate: not(a) or b
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) OrNY(ca, cb *core.LweSample) *core.LweSample {
	MU := types.ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := core.NewLweSample(inOutParams)

	//compute: (0,1/8) - ca + cb
	OrNYConst := types.ModSwitchToTorus32(1, 8)
	core.LweNoiselessTrivial(tempResult, OrNYConst, inOutParams)
	core.LweSubTo(tempResult, ca, inOutParams)
	core.LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return core.TfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped OrYN Gate: a or not(b)
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) OrYN(ca, cb *core.LweSample) *core.LweSample {
	MU := types.ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := core.NewLweSample(inOutParams)

	//compute: (0,1/8) + ca - cb
	OrYNConst := types.ModSwitchToTorus32(1, 8)
	core.LweNoiselessTrivial(tempResult, OrYNConst, inOutParams)
	core.LweAddTo(tempResult, ca, inOutParams)
	core.LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return core.TfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped Mux(a,b,c) = a?b:c = a*b + not(a)*c
 * Takes in input 3 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func (bk *PublicKey) Mux(a, b, c *core.LweSample) *core.LweSample {
	MU := types.ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	extractedParams := &bk.Params.TgswParams.TlweParams.ExtractedLweparams

	tempResult := core.NewLweSample(inOutParams)
	tempResult1 := core.NewLweSample(extractedParams)
	u1 := core.NewLweSample(extractedParams)
	u2 := core.NewLweSample(extractedParams)

	//compute "AND(a,b)": (0,-1/8) + a + b
	AndConst := types.ModSwitchToTorus32(-1, 8)
	core.LweNoiselessTrivial(tempResult, AndConst, inOutParams)
	core.LweAddTo(tempResult, a, inOutParams)
	core.LweAddTo(tempResult, b, inOutParams)
	// Bootstrap without KeySwitch
	// tfheBootstrapWoKS_FFT(u1, bk.bkFFT, MU, tempResult);
	core.TfheBootstrapWoKSFFT(u1, bk.Bkw.BkFFT, MU, tempResult)

	//compute "AND(not(a),c)": (0,-1/8) - a + c
	core.LweNoiselessTrivial(tempResult, AndConst, inOutParams)
	core.LweSubTo(tempResult, a, inOutParams)
	core.LweAddTo(tempResult, c, inOutParams)
	// Bootstrap without KeySwitch
	//tfheBootstrapWoKS_FFT(u2, bk.bkFFT, MU, tempResult);
	core.TfheBootstrapWoKSFFT(u2, bk.Bkw.BkFFT, MU, tempResult)

	// Add u1=u1+u2
	MuxConst := types.ModSwitchToTorus32(1, 8)
	core.LweNoiselessTrivial(tempResult1, MuxConst, extractedParams)
	core.LweAddTo(tempResult1, u1, extractedParams)
	core.LweAddTo(tempResult1, u2, extractedParams)
	// Key switching
	//core.LweKeySwitch(result, bk.bkFFT.ks, tempResult1)
	return core.LweKeySwitch(bk.Bkw.BkFFT.Ks, tempResult1)
}
