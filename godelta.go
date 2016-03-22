package main

import (
	"flag"
	"fmt"
	"github.com/mtib/godolta/deltal"
	"io/ioutil"
	"os"
	"strings"
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
	case "encrypt", "e":
		encrypt, err := deltal.Encrypt(&file, pass, *check)
		if err != nil {
			panic(err)
		}
		filename := file + ".delta"
		if *outp != "" {
			filename = *outp
		}
		ioutil.WriteFile(filename, encrypt, os.ModePerm)
	case "decrypt", "d":
		decrypt, err := deltal.Decrypt(&file, pass, *check)
		if err != nil {
			panic(err)
		}
		filename := removeDelta(file)
		if *outp != "" {
			filename = *outp
		}
		ioutil.WriteFile(filename, decrypt, os.ModePerm)
	default:
		fmt.Println(help)
		flag.PrintDefaults()
	}
}

func removeDelta(file string) string {
	if strings.HasSuffix(file, ".delta") {
		return file[:len(file)-6]
	}
	return file
}
