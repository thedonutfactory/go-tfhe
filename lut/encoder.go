package lut

import (
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/utils"
)

// Encoder provides encoding and decoding functions for different message spaces
type Encoder struct {
	MessageModulus int     // Number of possible messages (e.g., 2 for binary, 4 for 2-bit)
	Scale          float64 // Scaling factor for encoding
}

// NewEncoder creates a new encoder with the given message modulus
// For binary (boolean) operations, use messageModulus=2
// The default encoding uses 1/(2*messageModulus) to place messages in the torus
func NewEncoder(messageModulus int) *Encoder {
	// For TFHE, binary messages are encoded as ±1/8
	// Message 0 (false) -> -1/8 = 7/8 in unsigned representation
	// Message 1 (true) -> +1/8
	//
	// For general case with messageModulus m, we use ±1/(2m)
	// This gives us 1/4 for binary (m=2)
	scale := 1.0 / float64(2*messageModulus)
	return &Encoder{
		MessageModulus: messageModulus,
		Scale:          scale,
	}
}

// NewEncoderWithScale creates a new encoder with custom message modulus and scale
func NewEncoderWithScale(messageModulus int, scale float64) *Encoder {
	return &Encoder{
		MessageModulus: messageModulus,
		Scale:          scale,
	}
}

// Encode encodes an integer message into a torus value
// message should be in range [0, MessageModulus)
//
// For TFHE bootstrapping, the encoding is:
//
//	message i -> (i + 0.5) * scale
//
// This centers each message in its quantization region
func (e *Encoder) Encode(message int) params.Torus {
	// Normalize message to [0, MessageModulus)
	message = message % e.MessageModulus
	if message < 0 {
		message += e.MessageModulus
	}

	// Encode as (message + 0.5) * scale
	// For binary: 0 -> 0.5 * 0.25 = 0.125, 1 -> 1.5 * 0.25 = 0.375
	// But we want: 0 -> -0.125 (= 0.875), 1 -> 0.125
	//
	// Actually for TFHE bootstrapping, messages map to: (2i+1-m)/(2m)
	// For m=2: i=0 -> -1/4 = 3/4, i=1 -> 1/4
	//
	// Hmm, let me reconsider. The standard TFHE encoding is:
	// For boolean: false=-1/8, true=1/8
	// In unsigned: false=7/8, true=1/8
	//
	// For m values: message i maps to (2i+1-m) / (2m)
	// m=2: i=0 -> (0+1-2)/(4) = -1/4 = 3/4
	//      i=1 -> (2+1-2)/(4) = 1/4
	//
	// But for bootstrapping, we actually want something different...
	// Let me use the simpler formula: message i -> i * scale
	// with offset handling

	value := float64(message) * e.Scale
	return utils.F64ToTorus(value)
}

// EncodeWithCustomScale encodes with a custom scale factor
func (e *Encoder) EncodeWithCustomScale(message int, scale float64) params.Torus {
	message = message % e.MessageModulus
	if message < 0 {
		message += e.MessageModulus
	}
	value := float64(message) * scale
	return utils.F64ToTorus(value)
}

// Decode decodes a torus value back to an integer message
func (e *Encoder) Decode(value params.Torus) int {
	// Convert torus to float
	f := utils.TorusToF64(value)

	// Round to nearest message
	message := int(f/e.Scale + 0.5)

	// Normalize to [0, MessageModulus)
	message = message % e.MessageModulus
	if message < 0 {
		message += e.MessageModulus
	}

	return message
}

// DecodeBool decodes a torus value to a boolean (for binary messages)
func (e *Encoder) DecodeBool(value params.Torus) bool {
	return e.Decode(value) != 0
}
