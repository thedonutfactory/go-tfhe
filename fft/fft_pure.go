//go:build !rust
// +build !rust

package fft

// This file is used when NOT building with the rust tag
// It uses the pure Go implementation (go-dsp/fft)

// The existing fft.go implementation is already pure Go,
// so this file just ensures compatibility with the build tag system.

// When building without -tags rust, the existing FFTProcessor
// implementation in fft.go will be used.

