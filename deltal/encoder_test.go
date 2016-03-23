package deltal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type cryptTest struct {
	In, Out, Pass string
	Check         bool
}

var (
	crypting = []cryptTest{
		{"encoder.go", "../test/encoder.go.pw.delta", "pw", true},
		{"encoder.go", "../test/encoder.go.delta", "", true},
		{"encoder.go", "../test/encoder.go.nocheck.delta", "", false},
		{"encoder.go", "../test/encoder.go.pw.nocheck.delta", "pw", false},
		{"old/preReaderDelta.txt", "../test/preReaderDelta.txt.pw.delta", "pw", true},
		{"old/preReaderDelta.txt", "../test/preReaderDelta.txt.delta", "", true},
		{"old/preReaderDelta.txt", "../test/preReaderDelta.txt.nocheck.delta", "", false},
		{"old/preReaderDelta.txt", "../test/preReaderDelta.txt.pw.nocheck.delta", "pw", false},
	}
)

func (c *cryptTest) String() string {
	return fmt.Sprintf("Encoder: %s <-> %s using PW: \"%s\" Checksum: %v", c.In, c.Out, c.Pass, c.Check)
}

func TestEncoder(t *testing.T) {
	os.Mkdir("../test", os.ModePerm)
	writtenFiles := make(map[string]([]string))
	for _, v := range crypting {
		fmt.Println(v.String())

		// Encryption
		filein, err := os.Open(v.In)
		check(err, t)
		encoder, err := NewEncoderReader(filein, v.Pass, v.Check)
		check(err, t)
		encryptedData, err := ioutil.ReadAll(encoder)
		check(err, t)
		ioutil.WriteFile(v.Out, encryptedData, os.ModePerm)

		// Decryption
		cryptin, err := os.Open(v.Out)
		check(err, t)
		decoder := NewDecoderStream(cryptin, v.Pass)
		decryptedData, err := ioutil.ReadAll(decoder)
		check(err, t)
		decryptOut := v.Out + ".testdec"
		writtenFiles[v.In] = append(writtenFiles[v.In], decryptOut)
		ioutil.WriteFile(decryptOut, decryptedData, os.ModePerm)

		// Testing Seeker
		fmt.Println(decoder)
		decoder.Seek(0, 0)
	}
	for source, decrypts := range writtenFiles {
		sourceData, _ := ioutil.ReadFile(source)
		for _, file := range decrypts {
			decrpytData, _ := ioutil.ReadFile(file)
			if bytes.Equal(sourceData, decrpytData) {
				fmt.Println("Successfull", source, file)
			} else {
				fmt.Println("Failed", source, file)
				t.Fail()
			}
		}
	}
}

func BenchmarkEncryptNoPW(b *testing.B) {
	bdf, fd := benchEncrypter("../test/bench.txt", "")
	for i := 0; i < b.N; i++ {
		bdf()
	}
	fd.Close()
}
func BenchmarkEncryptPW(b *testing.B) {
	bdf, fd := benchEncrypter("../test/bench.txt", "password")
	for i := 0; i < b.N; i++ {
		bdf()
	}
	fd.Close()
}
func BenchmarkDecryptNoPW(b *testing.B) {
	bdf, fd := benchDecrypter("../test/bench.txt", "")
	for i := 0; i < b.N; i++ {
		bdf()
	}
	fd.Close()
}
func BenchmarkDecryptPW(b *testing.B) {
	bdf, fd := benchDecrypter("../test/bench.txt", "password")
	for i := 0; i < b.N; i++ {
		bdf()
	}
	fd.Close()
}

func check(err error, t *testing.T) {
	if err != nil {
		t.Fail()
		fmt.Println(err)
	}
}
