package gates

import (
	"github.com/thedonutfactory/go-tfhe/core"
)

type Int1 = [1]*core.LweSample
type Int2 = [2]*core.LweSample
type Int4 = [4]*core.LweSample
type Int8 = [8]*core.LweSample
type Int16 = [16]*core.LweSample
type Int32 = [32]*core.LweSample
type Int64 = [64]*core.LweSample

type Int = []*core.LweSample

type Ctxt = []*core.LweSample
type Ptxt = []bool

func (p *GateBootstrappingParameterSet) Int(size int) Int {
	return NewInt(size, p.InOutParams)
}

func NewInt(size int, params *core.LweParams) Int {
	var s Ctxt = core.NewLweSampleArray(int32(size), params)
	for i := 0; i < size; i++ {
		s[i] = core.NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Ctxt(size int) Ctxt {
	return NewInt(size, p.InOutParams)
}

func NewCtxt(size int, params *core.LweParams) Ctxt {
	var s Ctxt = core.NewLweSampleArray(int32(size), params)
	for i := 0; i < size; i++ {
		s[i] = core.NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int1() Int1 {
	return NewInt1(p.InOutParams)
}

func NewInt1(params *core.LweParams) Int1 {
	var s Int1
	s[0] = core.NewLweSample(params)
	return s
}

func (p *GateBootstrappingParameterSet) Int2() Int2 {
	return NewInt2(p.InOutParams)
}

func NewInt2(params *core.LweParams) Int2 {
	var s Int2
	for i := 0; i < 2; i++ {
		s[i] = core.NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int4() Int4 {
	return NewInt4(p.InOutParams)
}
func NewInt4(params *core.LweParams) Int4 {
	var s Int4
	for i := 0; i < 4; i++ {
		s[i] = core.NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int8() Int8 {
	return NewInt8(p.InOutParams)
}
func NewInt8(params *core.LweParams) Int8 {
	var s Int8
	for i := 0; i < 8; i++ {
		s[i] = core.NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int16() Int16 {
	return NewInt16(p.InOutParams)
}
func NewInt16(params *core.LweParams) Int16 {
	var s Int16
	for i := 0; i < 16; i++ {
		s[i] = core.NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int32() Int32 {
	return NewInt32(p.InOutParams)
}
func NewInt32(params *core.LweParams) Int32 {
	var s Int32
	for i := 0; i < 32; i++ {
		s[i] = core.NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int64() Int64 {
	return NewInt64(p.InOutParams)
}
func NewInt64(params *core.LweParams) Int64 {
	var s Int64
	for i := 0; i < 64; i++ {
		s[i] = core.NewLweSample(params)
	}
	return s
}
