package deltal

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const (
	testInput    = "encoder.go"
	testOutDir   = "../test/"
	testOutFile  = "encoder.go.pw.delta"
	testPassword = "pw"
)

func TestEncoder(t *testing.T) {
	encoder, err := NewEncoder(testInput, testPassword, true)
	check(err)
	data, err := ioutil.ReadAll(encoder)
	fmt.Println("Encrypt:", data[:12])
	fmt.Println("Passhash (e):", encoder.passarray)
	check(err)
	os.Mkdir(testOutDir, os.ModePerm)
	ioutil.WriteFile(testOutDir+testOutFile, data, os.ModePerm)

	file, err := os.Open(testOutDir + testOutFile)
	check(err)
	decoder := NewDecoderStream(file, testPassword)
	fmt.Println("Passhash (d):", encoder.passarray)
	data, err = ioutil.ReadAll(decoder)
	check(err)
	ioutil.WriteFile(testOutDir+"decrypted.txt", data, os.ModePerm)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
