// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// fnv 包实现了 FNV-1 和 FNV-1a，这是由 Glenn Fowler，Landon Curt Noll 和 Phong Vo 创建的非加密哈希函数。
// 参见：https://en.wikipedia.org/wiki/Fowler-Noll-Vo_hash_function。
//
// 此包返回的所有 hash.Hash 实现也实现了 encoding.BinaryMarshaler 和 encoding.BinaryUnmarshaler 来封装和取消封装 hash 的内部状态。
package fnv

import (
	"errors"
	"hash"
	"math/bits"
)

type (
	sum32   uint32
	sum32a  uint32
	sum64   uint64
	sum64a  uint64
	sum128  [2]uint64
	sum128a [2]uint64
)

const (
	offset32        = 2166136261
	offset64        = 14695981039346656037
	offset128Lower  = 0x62b821756295c58d
	offset128Higher = 0x6c62272e07bb0142
	prime32         = 16777619
	prime64         = 1099511628211
	prime128Lower   = 0x13b
	prime128Shift   = 24
)

// New32 返回一个新的 32 位的 FNV-1 hash.Hash。
// 它的 Sum 方法将按照 big-endian 字节顺序排列值。
func New32() hash.Hash32 {
	var s sum32 = offset32
	return &s
}

// New32a 返回一个新的 32 位的 FNV-1a hash.Hash。
// 它的 Sum 方法将按照 big-endian 字节顺序排列值。
func New32a() hash.Hash32 {
	var s sum32a = offset32
	return &s
}

// New64 返回一个新的 64 位的 FNV-1 hash.Hash。
// 它的 Sum 方法将按照 big-endian 字节顺序排列值。
func New64() hash.Hash64 {
	var s sum64 = offset64
	return &s
}

// New64a 返回一个新的 64 位的 FNV-1a hash.Hash。
// 它的 Sum 方法将按照 big-endian 字节顺序排列值。
func New64a() hash.Hash64 {
	var s sum64a = offset64
	return &s
}

// New128 返回一个新的 128 位的 FNV-1 hash.Hash。
// 它的 Sum 方法将按照 big-endian 字节顺序排列值。
func New128() hash.Hash {
	var s sum128
	s[0] = offset128Higher
	s[1] = offset128Lower
	return &s
}

// New128a 返回一个新的 128 位的 FNV-1a hash.Hash。
// 它的 Sum 方法将按照 big-endian 字节顺序排列值。
func New128a() hash.Hash {
	var s sum128a
	s[0] = offset128Higher
	s[1] = offset128Lower
	return &s
}

func (s *sum32) Reset()   { *s = offset32 }
func (s *sum32a) Reset()  { *s = offset32 }
func (s *sum64) Reset()   { *s = offset64 }
func (s *sum64a) Reset()  { *s = offset64 }
func (s *sum128) Reset()  { s[0] = offset128Higher; s[1] = offset128Lower }
func (s *sum128a) Reset() { s[0] = offset128Higher; s[1] = offset128Lower }

func (s *sum32) Sum32() uint32  { return uint32(*s) }
func (s *sum32a) Sum32() uint32 { return uint32(*s) }
func (s *sum64) Sum64() uint64  { return uint64(*s) }
func (s *sum64a) Sum64() uint64 { return uint64(*s) }

func (s *sum32) Write(data []byte) (int, error) {
	hash := *s
	for _, c := range data {
		hash *= prime32
		hash ^= sum32(c)
	}
	*s = hash
	return len(data), nil
}

func (s *sum32a) Write(data []byte) (int, error) {
	hash := *s
	for _, c := range data {
		hash ^= sum32a(c)
		hash *= prime32
	}
	*s = hash
	return len(data), nil
}

func (s *sum64) Write(data []byte) (int, error) {
	hash := *s
	for _, c := range data {
		hash *= prime64
		hash ^= sum64(c)
	}
	*s = hash
	return len(data), nil
}

func (s *sum64a) Write(data []byte) (int, error) {
	hash := *s
	for _, c := range data {
		hash ^= sum64a(c)
		hash *= prime64
	}
	*s = hash
	return len(data), nil
}

func (s *sum128) Write(data []byte) (int, error) {
	for _, c := range data {
		// Compute the multiplication
		s0, s1 := bits.Mul64(prime128Lower, s[1])
		s0 += s[1]<<prime128Shift + prime128Lower*s[0]
		// Update the values
		s[1] = s1
		s[0] = s0
		s[1] ^= uint64(c)
	}
	return len(data), nil
}

func (s *sum128a) Write(data []byte) (int, error) {
	for _, c := range data {
		s[1] ^= uint64(c)
		// Compute the multiplication
		s0, s1 := bits.Mul64(prime128Lower, s[1])
		s0 += s[1]<<prime128Shift + prime128Lower*s[0]
		// Update the values
		s[1] = s1
		s[0] = s0
	}
	return len(data), nil
}

func (s *sum32) Size() int   { return 4 }
func (s *sum32a) Size() int  { return 4 }
func (s *sum64) Size() int   { return 8 }
func (s *sum64a) Size() int  { return 8 }
func (s *sum128) Size() int  { return 16 }
func (s *sum128a) Size() int { return 16 }

func (s *sum32) BlockSize() int   { return 1 }
func (s *sum32a) BlockSize() int  { return 1 }
func (s *sum64) BlockSize() int   { return 1 }
func (s *sum64a) BlockSize() int  { return 1 }
func (s *sum128) BlockSize() int  { return 1 }
func (s *sum128a) BlockSize() int { return 1 }

func (s *sum32) Sum(in []byte) []byte {
	v := uint32(*s)
	return append(in, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (s *sum32a) Sum(in []byte) []byte {
	v := uint32(*s)
	return append(in, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (s *sum64) Sum(in []byte) []byte {
	v := uint64(*s)
	return append(in, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (s *sum64a) Sum(in []byte) []byte {
	v := uint64(*s)
	return append(in, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (s *sum128) Sum(in []byte) []byte {
	return append(in,
		byte(s[0]>>56), byte(s[0]>>48), byte(s[0]>>40), byte(s[0]>>32), byte(s[0]>>24), byte(s[0]>>16), byte(s[0]>>8), byte(s[0]),
		byte(s[1]>>56), byte(s[1]>>48), byte(s[1]>>40), byte(s[1]>>32), byte(s[1]>>24), byte(s[1]>>16), byte(s[1]>>8), byte(s[1]),
	)
}

func (s *sum128a) Sum(in []byte) []byte {
	return append(in,
		byte(s[0]>>56), byte(s[0]>>48), byte(s[0]>>40), byte(s[0]>>32), byte(s[0]>>24), byte(s[0]>>16), byte(s[0]>>8), byte(s[0]),
		byte(s[1]>>56), byte(s[1]>>48), byte(s[1]>>40), byte(s[1]>>32), byte(s[1]>>24), byte(s[1]>>16), byte(s[1]>>8), byte(s[1]),
	)
}

const (
	magic32          = "fnv\x01"
	magic32a         = "fnv\x02"
	magic64          = "fnv\x03"
	magic64a         = "fnv\x04"
	magic128         = "fnv\x05"
	magic128a        = "fnv\x06"
	marshaledSize32  = len(magic32) + 4
	marshaledSize64  = len(magic64) + 8
	marshaledSize128 = len(magic128) + 8*2
)

func (s *sum32) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, marshaledSize32)
	b = append(b, magic32...)
	b = appendUint32(b, uint32(*s))
	return b, nil
}

func (s *sum32a) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, marshaledSize32)
	b = append(b, magic32a...)
	b = appendUint32(b, uint32(*s))
	return b, nil
}

func (s *sum64) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, marshaledSize64)
	b = append(b, magic64...)
	b = appendUint64(b, uint64(*s))
	return b, nil

}

func (s *sum64a) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, marshaledSize64)
	b = append(b, magic64a...)
	b = appendUint64(b, uint64(*s))
	return b, nil
}

func (s *sum128) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, marshaledSize128)
	b = append(b, magic128...)
	b = appendUint64(b, s[0])
	b = appendUint64(b, s[1])
	return b, nil
}

func (s *sum128a) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, marshaledSize128)
	b = append(b, magic128a...)
	b = appendUint64(b, s[0])
	b = appendUint64(b, s[1])
	return b, nil
}

func (s *sum32) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic32) || string(b[:len(magic32)]) != magic32 {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize32 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	*s = sum32(readUint32(b[4:]))
	return nil
}

func (s *sum32a) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic32a) || string(b[:len(magic32a)]) != magic32a {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize32 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	*s = sum32a(readUint32(b[4:]))
	return nil
}

func (s *sum64) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic64) || string(b[:len(magic64)]) != magic64 {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize64 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	*s = sum64(readUint64(b[4:]))
	return nil
}

func (s *sum64a) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic64a) || string(b[:len(magic64a)]) != magic64a {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize64 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	*s = sum64a(readUint64(b[4:]))
	return nil
}

func (s *sum128) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic128) || string(b[:len(magic128)]) != magic128 {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize128 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	s[0] = readUint64(b[4:])
	s[1] = readUint64(b[12:])
	return nil
}

func (s *sum128a) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic128a) || string(b[:len(magic128a)]) != magic128a {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize128 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	s[0] = readUint64(b[4:])
	s[1] = readUint64(b[12:])
	return nil
}

func readUint32(b []byte) uint32 {
	_ = b[3]
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

func appendUint32(b []byte, x uint32) []byte {
	a := [4]byte{
		byte(x >> 24),
		byte(x >> 16),
		byte(x >> 8),
		byte(x),
	}
	return append(b, a[:]...)
}

func appendUint64(b []byte, x uint64) []byte {
	a := [8]byte{
		byte(x >> 56),
		byte(x >> 48),
		byte(x >> 40),
		byte(x >> 32),
		byte(x >> 24),
		byte(x >> 16),
		byte(x >> 8),
		byte(x),
	}
	return append(b, a[:]...)
}

func readUint64(b []byte) uint64 {
	_ = b[7]
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}
