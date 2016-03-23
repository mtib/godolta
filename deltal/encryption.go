package deltal

import (
	"github.com/dchest/siphash"
)

func u64tobyte(num uint64) (b []byte) {
	var i uint
	b = make([]byte, 8)
	for i = 0; i < 8; i++ {
		mov := i * uint(8)
		mask := uint64(0xFF) << mov
		masked := num & mask
		b[i] = byte((masked >> mov) % 256)
	}
	return
}

func rustHash(b []byte) []byte {
	return u64tobyte(intHash(b))
}

func intHash(b []byte) uint64 {
	return siphash.Hash(uint64(0), uint64(0), b)
}
