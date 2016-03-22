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
		// experimental
		res = append(res, 76, 10)
		res = append(res, u64tobyte(siphash.Hash(0, 0, toen))...) // Checksum

	} else {
		// Not using checksum
		res = append(res, 108, 10)
		// return nil, DeltaError("Not implemented")
	}
	// Start adding data
	res = append(res, toen...)
	// "encrypt"
	var old0, old1 byte
	for i := range res {
		k := i - 4
		if *checksum {
			k -= 8
		}
		old1 = old0
		old0 = res[i]
		if k >= 0 {
			fmt.Printf("%02X->", res[i])
			res[i] = byte((int(res[i]) + int(paarr[k%8])) % 256)
			fmt.Printf("%02X:", res[i])
			if k == 0 {
				fmt.Print(";\n")
			}
		} else {
			fmt.Printf("%02X;", res[i])
		}
		if k > 0 {
			res[i] = byte((int(res[i]) + int(old1)) % 256)
			fmt.Printf("->%02X(%02X);\n", res[i], old1)
		}
	}
	// fin "encrypt"
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
		b[7-i] = byte((masked >> mov) % 256)
	}
	return
}
