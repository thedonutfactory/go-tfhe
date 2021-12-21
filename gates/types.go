package gates

import (
	. "github.com/thedonutfactory/go-tfhe/core"
)

type Int1 = [1]*LweSample
type Int2 = [2]*LweSample
type Int4 = [4]*LweSample
type Int8 = [8]*LweSample
type Int16 = [16]*LweSample
type Int32 = [32]*LweSample
type Int64 = [64]*LweSample

type Int = []*LweSample

func (p *GateBootstrappingParameterSet) Int(size int) Int {
	return NewInt(size, p.InOutParams)
}

func NewInt(size int, params *LweParams) Int {
	var s Int = NewLweSampleArray(int32(size), params)
	for i := 0; i < size; i++ {
		s[i] = NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int1() Int1 {
	return NewInt1(p.InOutParams)
}

func NewInt1(params *LweParams) Int1 {
	var s Int1
	s[0] = NewLweSample(params)
	return s
}

func (p *GateBootstrappingParameterSet) Int2() Int2 {
	return NewInt2(p.InOutParams)
}

func NewInt2(params *LweParams) Int2 {
	var s Int2
	for i := 0; i < 2; i++ {
		s[i] = NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int4() Int4 {
	return NewInt4(p.InOutParams)
}
func NewInt4(params *LweParams) Int4 {
	var s Int4
	for i := 0; i < 4; i++ {
		s[i] = NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int8() Int8 {
	return NewInt8(p.InOutParams)
}
func NewInt8(params *LweParams) Int8 {
	var s Int8
	for i := 0; i < 8; i++ {
		s[i] = NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int16() Int16 {
	return NewInt16(p.InOutParams)
}
func NewInt16(params *LweParams) Int16 {
	var s Int16
	for i := 0; i < 16; i++ {
		s[i] = NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int32() Int32 {
	return NewInt32(p.InOutParams)
}
func NewInt32(params *LweParams) Int32 {
	var s Int32
	for i := 0; i < 32; i++ {
		s[i] = NewLweSample(params)
	}
	return s
}

func (p *GateBootstrappingParameterSet) Int64() Int64 {
	return NewInt64(p.InOutParams)
}
func NewInt64(params *LweParams) Int64 {
	var s Int64
	for i := 0; i < 64; i++ {
		s[i] = NewLweSample(params)
	}
	return s
}
