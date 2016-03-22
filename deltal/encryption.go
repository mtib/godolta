package deltal

import (
	"fmt"
	"github.com/dchest/siphash"
	"io/ioutil"
)

// DeltaError to be thrown when encryption fails
type DeltaError string

// TODO rewrite this as io.Writer / io.Reader

// Encrypt returns encrypted files bytes
func Encrypt(file, pass *string, checksum *bool) ([]byte, error) {
	toen, err := ioutil.ReadFile(*file)
	fmt.Println("checksum:", *checksum)
	paarr := u64tobyte(siphash.Hash(0, 0, []byte(*pass)))
	if *pass == "" {
		paarr = u64tobyte(0)
	}
	fmt.Println(paarr)
	if err != nil {
		panic(DeltaError("Cannot read file"))
	}
	var res []byte
	res = append(res, 206, 148)
	if *checksum {
		res = append(res, 76, 10)            // = L\n
		res = append(res, rustHash(toen)...) // == Checksum

	} else { // Not using Checksum
		res = append(res, 108, 10) // = l\n
	}
	for i := range toen { // This could be realized via io.Reader
		res = append(res, byte((int(toen[i])+int(paarr[i%8]))%256))
		if i > 0 {
			res[len(res)-1] = byte((int(res[len(res)-1]) + int(toen[i-1])) % 256)
		}
	} // Added all data & encrypted it
	return res, nil
}

func (d DeltaError) Error() string {
	return "Delta En/De-cryption error: " + string(d)
}

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
