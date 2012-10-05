// Written in 2012 by Dmitry Chestnykh.
//
// To the extent possible under law, the author have dedicated all copyright
// and related and neighboring rights to this software to the public domain
// worldwide. This software is distributed without any warranty.
// http://creativecommons.org/publicdomain/zero/1.0/

package varuint

import "testing"

func testUint64(t *testing.T, x uint64) {
	var tmp [MaxUint64Len]byte
	n := PutUint64(tmp[:], x)
	y, w := Uint64(tmp[:n])
	if n != w {
		t.Errorf("different number of bytes: expected %d, got %d", n, w)
	}
	if x != y {
		t.Errorf("expected %d, got %d", x, y)
	}
}

var tests = []uint64{
	0,
	240,
	241,
	2287,
	2288,
	67823,
	67824,
	16777215,
	16777216,
	1<<24-1,
	1<<24,
	1<<32-1,
	1<<32,
	1<<64-1,
}

func TestVaruint(t *testing.T) {
	for i := uint64(0); i < 3000; i++ {
		testUint64(t, i)
	}

	for _, v := range tests {
		testUint64(t, v)
	}
}

func BenchmarkPutUint64(b *testing.B) {
	buf := make([]byte, MaxUint64Len)
	b.SetBytes(8)
	for i := 0; i < b.N; i++ {
		for j := uint(0); j < MaxUint64Len; j++ {
			PutUint64(buf, 1<<(j*7))
		}
	}
}

func BenchmarkUint64(b *testing.B) {
	buf := make([][]byte, len(tests))
	bytes := int64(0)
	for i, v := range tests {
		buf[i] = make([]byte, MaxUint64Len)
		bytes += int64(PutUint64(buf[i], v))
	}
	b.SetBytes(bytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range buf {
			Uint64(v)
		}
	}
}
