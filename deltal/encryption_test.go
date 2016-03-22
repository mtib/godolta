package deltal

import (
	"fmt"
	"testing"
)

func TestIntToArray(t *testing.T) {
	var num uint64
	num = 0x8899AABBCCDDEEFF
	fmt.Printf("\"%016X\" -> ", num)
	res := u64tobyte(num)
	for _, v := range res {
		fmt.Printf("%02X ", v)
	}
	fmt.Print("\n")
}

func TestHashFunction(t *testing.T) {
	ta := []string{"I am a teapot", "hello world"}
	for _, text := range ta {
		fmt.Printf("\"%s\" -> %016X\n", text, intHash([]byte(text)))
		res := rustHash([]byte(text))
		fmt.Printf("\"%s\" (%v) -> ", text, []byte(text))
		for _, v := range res {
			fmt.Printf("%02X ", v)
		}
		fmt.Print("\n")
	}
}
