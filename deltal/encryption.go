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
func Encrypt(file, pass *string, checksum bool) ([]byte, error) {
	toen, err := ioutil.ReadFile(*file)
	if err != nil {
		panic(DeltaError("Cannot read file"))
	}
	var res []byte
	res = append(res, 206, 148)
	if checksum {
		// experimental
		res = append(res, 76, 10)
		for i := 0; i < len(toen); i += 8 {
			hnum := siphash.Hash(0, 0, toen[i:i+8])
			fmt.Println(hnum)
			// if err != nil {
			// 	panic(err)
			// }
		}
	} else {
		return nil, DeltaError("Not implemented")
	}
	fmt.Println(res)
	return res, nil
}

func (d DeltaError) Error() string {
	return "Delta En/De-cryption error: " + string(d)
}

func u64tobyte(num uint64) (b [8]byte) {
	var i uint
	fmt.Printf("%016X\n", num)
	for i = 0; i < 8; i++ {
		mov := i * uint(8)
		mask := uint64(0xFF) << mov
		masked := num & mask
		fmt.Printf("%016X -> ", masked)
		b[7-i] = byte((masked >> mov) % 256)
		fmt.Printf("%02X\n", b[7-i])
	}
	return
}
