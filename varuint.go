// Written in 2012 by Dmitry Chestnykh.
//
// To the extent possible under law, the author have dedicated all copyright
// and related and neighboring rights to this software to the public domain
// worldwide. This software is distributed without any warranty.
// http://creativecommons.org/publicdomain/zero/1.0/

// Package varuint implements SQLite4-like variable unsigned integer encoding
// as described in http://www.sqlite.org/src4/doc/trunk/www/varint.wiki.
//
// Unlike varint from encoding/binary package, this encoding uses fewer bytes
// for smaller values, and the number of encoded bytes can be determined by
// looking at the first byte. Varuint also preserves numeric and lexicographical
// ordering.
package varuint

// Maximum number of bytes required to encode uint64.
const MaxUint64Len = 9

// PutUint64 encodes a uint64 into b and returns the number of bytes written.
// Buffer must have enough space to encode the number.
func PutUint64(b []byte, v uint64) int {
	if v <= 240 {
		b[0] = byte(v)
		return 1
	}
	if v <= 2287 {
		b[0] = byte((v-240)/256 + 241)
		b[1] = byte((v - 240) % 256)
		return 2
	}
	if v <= 67823 {
		b[0] = 249
		b[1] = byte((v - 2288) / 256)
		b[2] = byte((v - 2288) % 256)
		return 3
	}
	if v <= 1<<24-1 {
		b[0] = 250
		b[1] = byte(v >> 16)
		b[2] = byte(v >> 8)
		b[3] = byte(v)
		return 4
	}
	if v <= 1<<31-1 {
		b[0] = 251
		b[1] = byte(v >> 24)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 8)
		b[4] = byte(v)
		return 5
	}
	if v <= 1<<40-1 {
		b[0] = 252
		b[1] = byte(v >> 32)
		b[2] = byte(v >> 24)
		b[3] = byte(v >> 16)
		b[4] = byte(v >> 8)
		b[5] = byte(v)
		return 6
	}
	if v <= 1<<48-1 {
		b[0] = 253
		b[1] = byte(v >> 40)
		b[2] = byte(v >> 32)
		b[3] = byte(v >> 24)
		b[4] = byte(v >> 16)
		b[5] = byte(v >> 8)
		b[6] = byte(v)
		return 7
	}
	if v <= 1<<56-1 {
		b[0] = 254
		b[1] = byte(v >> 48)
		b[2] = byte(v >> 40)
		b[3] = byte(v >> 32)
		b[4] = byte(v >> 24)
		b[5] = byte(v >> 16)
		b[6] = byte(v >> 8)
		b[7] = byte(v)
		return 8
	}
	b[0] = 255
	b[1] = byte(v >> 56)
	b[2] = byte(v >> 48)
	b[3] = byte(v >> 40)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 24)
	b[6] = byte(v >> 16)
	b[7] = byte(v >> 8)
	b[8] = byte(v)
	return 9
}

// Uint64 returns a decoded uint64 from b and the number of bytes read.
// If the buffer is too small, returns n < 0, where -n is the expected
// length of buffer.
func Uint64(b []byte) (v uint64, n int) {
	if len(b) < 1 {
		return 0, -1
	}
	b0 := b[0]
	if b0 <= 240 {
		v = uint64(b0)
		n = 1
		return
	}
	if b0 <= 248 {
		if len(b) < 2 {
			return 0, -2
		}
		v = uint64(b0-241)*256 + uint64(b[1]) + 240
		n = 2
		return
	}
	if b0 <= 249 {
		if len(b) < 3 {
			return 0, -3
		}
		v = 2288 + (256 * uint64(b[1])) + uint64(b[2])
		n = 3
		return
	}
	if b0 == 250 {
		if len(b) < 4 {
			return 0, -4
		}
		v = uint64(b[1])<<16 | uint64(b[2])<<8 | uint64(b[3])
		n = 4
		return
	}
	if b0 == 251 {
		if len(b) < 5 {
			return 0, -5
		}
		v = uint64(b[1])<<24 | uint64(b[2])<<16 | uint64(b[3])<<8 | uint64(b[4])
		n = 5
		return
	}
	if b0 == 252 {
		if len(b) < 6 {
			return 0, -6
		}
		v = uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5])
		n = 6
		return
	}
	if b0 == 253 {
		if len(b) < 7 {
			return 0, -7
		}
		v = uint64(b[1])<<40 | uint64(b[2])<<32 | uint64(b[3])<<24 | uint64(b[4])<<16 | uint64(b[5])<<8 | uint64(b[6])
		n = 7
		return
	}
	if b0 == 254 {
		if len(b) < 8 {
			return 0, -8
		}
		v = uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 | uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
		n = 8
		return
	}
	if len(b) < 9 {
		return 0, -9
	}
	v = uint64(b[1])<<56 | uint64(b[2])<<48 | uint64(b[3])<<40 | uint64(b[4])<<32 | uint64(b[5])<<24 | uint64(b[6])<<16 | uint64(b[7])<<8 | uint64(b[8])
	n = 9
	return
}
