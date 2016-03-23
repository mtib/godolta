package main

import (
	"flag"
	"fmt"
	"github.com/mtib/godolta/deltal"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var (
	help = `USAGE:
    godelta encrypt <file> [-o <file>]
    godelta decrypt <file> [-o <file>]`
	pass = flag.String("p", "", "passphrase to use")
	outp = flag.String("o", "", "output file")
	// override = flag.Bool("y", false, "override existing file")
	check = flag.Bool("c", true, "use checksum feature")
	force = flag.Bool("f", false, "force overwrite")
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
		filename := file + ".delta"
		if *outp != "" {
			filename = *outp
		}
		f, err := tryCreateFile(filename)
		if err != nil {
			fmt.Println("Could not create file", filename)
			return
		}
		encoder, err := deltal.NewEncoder(file, *pass, *check)
		if err != nil {
			fmt.Println("Could not create encoder\n", err)
			return
		}
		io.Copy(f, encoder)
	case "decrypt", "d":
		filename := removeDelta(file)
		if *outp != "" {
			filename = *outp
		}
		f, err := tryCreateFile(filename)
		if err != nil {
			fmt.Println("Could not create file", filename)
			return
		}
		src, err := os.Open(file)
		if err != nil {
			fmt.Println("Could not open file", file)
			return
		}
		decoder := deltal.NewDecoderStream(src, *pass)
		io.Copy(f, decoder)
		csum, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("Could not verify checksum")
			return
		}
		if decoder.UseChecksum && !decoder.Check(csum) {
			fmt.Println("Wrong checksum, proceed with caution")
		}
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

func tryCreateFile(file string) (*os.File, error) {
	f, err := os.Open(file)
	if err == nil && !*force {
		// File exists
		fmt.Println("File exists, to overwrite use -f")
		return nil, os.ErrExist
	}
	if err == nil && *force {
		os.Remove(f.Name())
		f, err = os.Create(file)
		return f, err
	}
	f, err = os.Create(file)
	if err != nil {
		fmt.Println("Could not create file:", file)
		return nil, err
	}
	return f, nil
}
