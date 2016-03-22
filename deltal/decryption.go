package deltal

import (
	"io/ioutil"
)

// Decrypt returns encrypted files bytes
func Decrypt(file, pass *string) ([]byte, error) {
	cdata, err := ioutil.ReadFile(*file)
	if err != nil {
		panic(DeltaError("Could not read file"))
	}
	return []byte(cdata), err
}
