// GoGOST -- Pure Go GOST cryptographic functions library
// Copyright (C) 2015-2016 Sergey Matveev <stargrave@stargrave.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package gost341194

import (
	"crypto/rand"
	"crypto/subtle"
	"hash"
	"testing"
	"testing/quick"

	"github.com/stargrave/gogost/gost28147"
)

func TestHashInterface(t *testing.T) {
	h := New(SboxDefault)
	var _ hash.Hash = h
}

func TestVectors(t *testing.T) {
	h := New(SboxDefault)

	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0x8d, 0x0f, 0x49, 0x49, 0x2c, 0x91, 0xf4, 0x5a,
		0x68, 0xff, 0x5c, 0x05, 0xd2, 0xc2, 0xb4, 0xab,
		0x78, 0x02, 0x7b, 0x9a, 0xab, 0x5c, 0xe3, 0xfe,
		0xff, 0x52, 0x67, 0xc4, 0x9c, 0xb9, 0x85, 0xce,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("a"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0xdd, 0x14, 0xf3, 0x62, 0xce, 0xfd, 0x49, 0xf8,
		0x73, 0xa5, 0xc6, 0x44, 0x43, 0x1b, 0x87, 0x21,
		0x9c, 0x34, 0x49, 0x66, 0x1f, 0x80, 0x8a, 0xc8,
		0xe9, 0x66, 0x7c, 0x36, 0x9e, 0x53, 0x2c, 0xd4,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("abc"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0x1d, 0xd5, 0xa4, 0x06, 0x7c, 0x49, 0x70, 0x3b,
		0x75, 0xbc, 0x75, 0xc9, 0x29, 0x0f, 0x5e, 0xcb,
		0xb5, 0xeb, 0x85, 0x22, 0x9e, 0x72, 0x77, 0xa2,
		0xb2, 0xb1, 0x4f, 0xc4, 0x48, 0x43, 0x13, 0xf3,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("message digest"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0x4d, 0x9a, 0x88, 0xa4, 0x16, 0xde, 0x2f, 0xdb,
		0x72, 0xde, 0x48, 0x3f, 0x27, 0x65, 0x2b, 0x58,
		0x69, 0x24, 0x3d, 0xec, 0x59, 0xbe, 0x0c, 0xb6,
		0x99, 0x2c, 0x8f, 0xb1, 0xec, 0x34, 0x44, 0xad,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	for i := 0; i < 128; i++ {
		h.Write([]byte("U"))
	}
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0xa4, 0x33, 0x57, 0xfe, 0xe8, 0xa9, 0x26, 0xd9,
		0x52, 0x2a, 0x06, 0x87, 0x0a, 0x66, 0x25, 0x1c,
		0x55, 0x3e, 0x27, 0x74, 0xa0, 0x85, 0x1d, 0x0c,
		0xef, 0x0c, 0x18, 0x25, 0xed, 0xa3, 0xa3, 0x53,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("The quick brown fox jumps over the lazy dog"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0x94, 0x42, 0x1f, 0x6d, 0x37, 0x0f, 0xa1, 0xd1,
		0x6b, 0xa7, 0xac, 0x5e, 0x31, 0x29, 0x65, 0x29,
		0xc9, 0x68, 0x04, 0x7d, 0xca, 0x9b, 0xf4, 0x25,
		0x8a, 0xc5, 0x9a, 0x0c, 0x41, 0xfa, 0xb7, 0x77,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("The quick brown fox jumps over the lazy cog"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0x45, 0xc4, 0xee, 0x4e, 0xe1, 0xd2, 0x50, 0x91,
		0x31, 0x21, 0x35, 0x54, 0x0d, 0x67, 0x02, 0xe6,
		0x67, 0x7f, 0x7a, 0x73, 0xb5, 0xda, 0x31, 0xe1,
		0x0b, 0x8b, 0xb7, 0xaa, 0xda, 0xc4, 0xeb, 0xa3,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("This is message, length=32 bytes"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0xfa, 0xff, 0x37, 0xa6, 0x15, 0xa8, 0x16, 0x69,
		0x1c, 0xff, 0x3e, 0xf8, 0xb6, 0x8c, 0xa2, 0x47,
		0xe0, 0x95, 0x25, 0xf3, 0x9f, 0x81, 0x19, 0x83,
		0x2e, 0xb8, 0x19, 0x75, 0xd3, 0x66, 0xc4, 0xb1,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("Suppose the original message has length = 50 bytes"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0x08, 0x52, 0xf5, 0x62, 0x3b, 0x89, 0xdd, 0x57,
		0xae, 0xb4, 0x78, 0x1f, 0xe5, 0x4d, 0xf1, 0x4e,
		0xea, 0xfb, 0xc1, 0x35, 0x06, 0x13, 0x76, 0x3a,
		0x0d, 0x77, 0x0a, 0xa6, 0x57, 0xba, 0x1a, 0x47,
	}) != 1 {
		t.Fail()
	}
}

func TestVectorsCryptoPro(t *testing.T) {
	h := New(&gost28147.GostR3411_94_CryptoProParamSet)

	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0xc0, 0x56, 0xd6, 0x4c, 0x23, 0x83, 0xc4, 0x4a,
		0x58, 0x13, 0x9c, 0x9b, 0x56, 0x01, 0x11, 0xac,
		0x13, 0x3e, 0x43, 0xfb, 0x84, 0x0f, 0x83, 0x87,
		0x14, 0x84, 0x0c, 0xa3, 0x3c, 0x5f, 0x1e, 0x98,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("a"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0x11, 0x30, 0x40, 0x2f, 0xcf, 0xaa, 0xf1, 0xef,
		0x3c, 0x13, 0xe3, 0x17, 0x3f, 0x10, 0x5a, 0x71,
		0x55, 0x80, 0xf7, 0xc9, 0x79, 0x00, 0xaf, 0x37,
		0xbf, 0x83, 0x21, 0x28, 0xdd, 0x52, 0x4c, 0xe7,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("abc"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0x2c, 0xd4, 0x2f, 0xf9, 0x86, 0x29, 0x3b, 0x16,
		0x7e, 0x99, 0x43, 0x81, 0xed, 0x59, 0x74, 0x74,
		0x14, 0xdd, 0x24, 0x95, 0x36, 0x77, 0x76, 0x2d,
		0x39, 0xd7, 0x18, 0xbf, 0x6d, 0x05, 0x85, 0xb2,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("message digest"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0xa0, 0x1b, 0x72, 0x29, 0x9b, 0xc3, 0x9a, 0x54,
		0x0f, 0xd6, 0x72, 0xa9, 0x9a, 0x72, 0xb4, 0xbd,
		0xfe, 0x74, 0x41, 0x73, 0x86, 0x98, 0x6e, 0xfa,
		0xeb, 0x01, 0xa4, 0x2a, 0xdd, 0x41, 0x60, 0xbc,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("The quick brown fox jumps over the lazy dog"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0x76, 0x0a, 0x83, 0x65, 0xd5, 0x70, 0x47, 0x6e,
		0x78, 0x72, 0x54, 0x76, 0x1b, 0xe7, 0x65, 0x67,
		0x74, 0x02, 0x1b, 0x1f, 0x3d, 0xe5, 0x6f, 0x58,
		0x8c, 0x50, 0x1a, 0x36, 0x4a, 0x29, 0x04, 0x90,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("This is message, length=32 bytes"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0xeb, 0x48, 0xde, 0x3e, 0x89, 0xe7, 0x1b, 0xcb,
		0x69, 0x5f, 0xc7, 0x52, 0xd6, 0x17, 0xfa, 0xe7,
		0x57, 0xf3, 0x4f, 0xa7, 0x7f, 0xa5, 0x8e, 0xe1,
		0x14, 0xc5, 0xbd, 0xb7, 0xf7, 0xc2, 0xef, 0x2c,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	h.Write([]byte("Suppose the original message has length = 50 bytes"))
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0x11, 0x50, 0xa6, 0x30, 0x31, 0xdc, 0x61, 0x1a,
		0x5f, 0x5e, 0x40, 0xd9, 0x31, 0x53, 0xf7, 0x4e,
		0xbd, 0xe8, 0x21, 0x6f, 0x67, 0x92, 0xc2, 0x5a,
		0x91, 0xcf, 0xca, 0xbc, 0x5c, 0x0c, 0x73, 0xc3,
	}) != 1 {
		t.Fail()
	}

	h.Reset()
	for i := 0; i < 128; i++ {
		h.Write([]byte{'U'})
	}
	if subtle.ConstantTimeCompare(h.Sum(nil), []byte{
		0xe8, 0xc4, 0x49, 0xf6, 0x08, 0x10, 0x4c, 0x51,
		0x27, 0x10, 0xcd, 0x37, 0xfd, 0xed, 0x92, 0x0d,
		0xf1, 0xe8, 0x6b, 0x21, 0x16, 0x23, 0xfa, 0x27,
		0xf4, 0xbb, 0x91, 0x46, 0x61, 0xc7, 0x4a, 0x1c,
	}) != 1 {
		t.Fail()
	}
}

func TestRandom(t *testing.T) {
	h := New(SboxDefault)
	f := func(data []byte) bool {
		h.Reset()
		h.Write(data)
		d1 := h.Sum(nil)
		h.Reset()
		for _, c := range data {
			h.Write([]byte{c})
		}
		d2 := h.Sum(nil)
		return subtle.ConstantTimeCompare(d1, d2) == 1
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func BenchmarkHash(b *testing.B) {
	h := New(SboxDefault)
	src := make([]byte, BlockSize+1)
	rand.Read(src)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(src)
		h.Sum(nil)
	}
}
