package deltal

import (
	"fmt"
	"github.com/dchest/siphash"
	"io/ioutil"
)

// DeltaError to be thrown when encryption fails
type DeltaError string

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
			h := siphash.New(append(toen[i:i+8], 0, 0, 0, 0, 0, 0, 0, 0))
			if err != nil {
				panic(err)
			}
			for _, v := range h.Sum(nil) {
				res = append(res, v)
			}
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
