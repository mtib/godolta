package main

import (
	"flag"
	"fmt"
	"github.com/mtib/godolta/deltal"
)

var (
	help = `USAGE:
    godelta encrypt <file> [-o <file>]
    godelta decrypt <file> [-o <file>]`
	pass     = flag.String("p", "", "passphrase to use")
	outp     = flag.String("o", "", "output file")
	override = flag.Bool("y", false, "override existing file")
	check    = flag.Bool("c", true, "disables checksum feature")
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		fmt.Println(help)
		flag.PrintDefaults()
		return
	}
	mode := flag.Arg(0)
	file := flag.Arg(1)
	switch mode {
	case "encrypt":
		deltal.Encrypt(&file)
	case "decrypt":
		deltal.Decrypt(&file)
	default:
		fmt.Println(help)
		flag.PrintDefaults()
	}
}
