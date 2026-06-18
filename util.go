package main

import (
	"encoding/binary"
	"math"
	"strings"
)

const hexChars = "0123456789abcdef"

// encode s with unicode escape sequences and add quotes
func Quote(s []byte, quote byte) string {
	var output strings.Builder
	output.Grow(len(s) * 2)
	// output.WriteByte(quote)
	for i := range len(s) {
		char := s[i]
		switch {
		case char == quote:
			output.WriteByte('\\')
			output.WriteByte(quote)
		case char == '\\':
			output.WriteByte('\\')
			output.WriteByte('\\')
		case char >= 0x20 && char < 0x7f:
			output.WriteByte(char)
		case char == 0x09:
			output.WriteByte('\\')
			output.WriteByte('t')
		case char == 0x0a:
			output.WriteByte('\\')
			output.WriteByte('n')
		case char == 0x0d:
			output.WriteByte('\\')
			output.WriteByte('r')
		default:
			output.WriteByte('\\')
			output.WriteByte('x')
			output.WriteByte(hexChars[char>>4])
			output.WriteByte(hexChars[char&0xf])
		}
	}
	// output.WriteByte(byte(quote))
	return output.String()
}

func encodeUShort(i int) []byte {
	higher := 0
	if i > 255 {
		higher = i << 8
	}
	lower := i & 0xff
	return []byte{byte(lower), byte(higher)}
}

func encodeFloat(f64 float64) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, math.Float32bits(float32(f64)))
	return bytes
}

func tex(f float64) float64 {
	return (f * textureSpaceFactor) + 0.5
}
