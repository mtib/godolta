package deltal

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncoder(t *testing.T) {
	encoder, err := NewEncoder("encoder.go", "pw", true)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(encoder)
	fmt.Println(data[:12])
	if err != nil {
		panic(err)
	}
	os.Mkdir("../test", os.ModePerm)
	ioutil.WriteFile("../test/encoder.go.pw.delta", data, os.ModePerm)
}
