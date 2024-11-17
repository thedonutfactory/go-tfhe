package types

import (
	"math"
	"testing"
)

func TestDoubleToTorus(t *testing.T) {
	tests := []struct {
		input    float64
		expected Torus32
	}{
		{0.0, 0},
		{-0.5, Torus32(-two32 / 2)},
		{0.25, Torus32(two32 / 4)},
		{-0.25, Torus32(-two32 / 4)},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := DoubleToTorus(tt.input)
			if result != tt.expected {
				t.Errorf("DoubleToTorus(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestApproxPhase(t *testing.T) {
	tests := []struct {
		phase    Torus32
		Msize    int32
		expected Torus32
	}{
		{123456789, 16, 0},
		{int32(two32 / 4), 4, Torus32(two32 / 4)},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := ApproxPhase(tt.phase, tt.Msize)
			if result != tt.expected {
				t.Errorf("ApproxPhase(%v, %v) = %v, expected %v", tt.phase, tt.Msize, result, tt.expected)
			}
		})
	}
}

func TestUniformFloat64Dist(t *testing.T) {
	result := UniformFloat64Dist(-5, 5)
	if result < -5 || result > 5 {
		t.Errorf("UniformFloat64Dist(-5, 5) generated out of range value: %v", result)
	}
}

func TestUniformTorus32Dist(t *testing.T) {
	result := UniformTorus32Dist()
	if result < math.MinInt32 || result > math.MaxInt32 {
		t.Errorf("UniformTorus32Dist() generated out of range value: %v", result)
	}
}

func TestModSwitchFromTorus32(t *testing.T) {
	tests := []struct {
		phase    Torus32
		Msize    int32
		expected int32
	}{
		{0, 4, 0},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := ModSwitchFromTorus32(tt.phase, tt.Msize)
			if result != tt.expected {
				t.Errorf("ModSwitchFromTorus32(%v, %v) = %v, expected %v", tt.phase, tt.Msize, result, tt.expected)
			}
		})
	}
}

func TestModSwitchToTorus32(t *testing.T) {
	tests := []struct {
		mu       int32
		Msize    int32
		expected Torus32
	}{
		{0, 4, 0},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := ModSwitchToTorus32(tt.mu, tt.Msize)
			if result != tt.expected {
				t.Errorf("ModSwitchToTorus32(%v, %v) = %v, expected %v", tt.mu, tt.Msize, result, tt.expected)
			}
		})
	}
}

func TestAbsfrac(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0.3, 0.3},
		{-0.7, 0.3}, // Closest integer is -1, so fractional part is 0.3
		{1.5, 0.5},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Absfrac(tt.input)
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Absfrac(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		input    Torus32
		expected Torus32
	}{
		{-10, 10},
		{10, 10},
		{0, 0},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Abs(tt.input)
			if result != tt.expected {
				t.Errorf("Abs(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAbsInt(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{-10, 10},
		{10, 10},
		{0, 0},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := AbsInt(tt.input)
			if result != tt.expected {
				t.Errorf("AbsInt(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
