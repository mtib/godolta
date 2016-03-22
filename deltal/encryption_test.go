package deltal

import (
	"fmt"
	"testing"
)

func TestIntToArray(t *testing.T) {
	var num uint64
	num = 0x8899AABBCCDDEEFF
	fmt.Println(u64tobyte(num))
}
