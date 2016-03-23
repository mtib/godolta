package deltal

import (
	"io"
	"io/ioutil"
	"os"
)

// Encoder of delta-l
type Encoder struct {
	Counter     uint64
	passarray   []byte
	FileReader  io.ReadSeeker
	Checksum    []byte
	UseChecksum bool
	headerpos   int
	last        uint8
}

// NewEncoderReader creates an Encoder
func NewEncoderReader(reader io.ReadSeeker, password string, checksum bool) (*Encoder, error) {
	full, err := ioutil.ReadAll(reader)
	reader.Seek(0, 0)

	var pa []byte
	if password != "" {
		pa = rustHash([]byte(password))
	} else {
		pa = u64tobyte(0)
	}
	return &Encoder{passarray: pa, FileReader: reader, UseChecksum: checksum, Checksum: rustHash(full)}, err
}

// NewEncoder returns pointer to decoder reader interface
func NewEncoder(file, password string, checksum bool) (*Encoder, error) {
	f, _ := os.Open(file)
	return NewEncoderReader(f, password, checksum)
}

func (d *Encoder) Read(b []byte) (n int, err error) {
	// Writing Checksum
	for (d.UseChecksum && d.headerpos < 12) || (!d.UseChecksum && d.headerpos < 4) {
		if d.headerpos > 3 {
			if d.UseChecksum {
				b[n] = d.Checksum[d.headerpos-4]
			} else {
				return
			}
		} else {
			b[n] = header(d.headerpos, d.UseChecksum)
		}
		d.headerpos++
		n++
		if n == len(b) {
			return
		}
	}

	// Encrypting File
	b2 := make([]byte, len(b)-n)
	n2, err := d.FileReader.Read(b2)
	for k := range b2 {
		s := b2[k]
		b[k+n] = uint8(b2[k] + d.last + d.passarray[d.Counter%uint64(len(d.passarray))])
		d.Counter++
		d.last = s
	}
	n += n2
	return
}

// FastEncrpyt is an easy-call function for benchmarking
func FastEncrpyt(file, pass string) {
	filein, _ := os.Open(file)
	encoder, _ := NewEncoderReader(filein, pass, true)
	encryptedData, _ := ioutil.ReadAll(encoder)
	ioutil.WriteFile(file+".delta", encryptedData, os.ModePerm)
}

var (
	checkHeader    = []byte{206, 148, 76, 10}
	nonCheckHeader = []byte{206, 148, 108, 10}
)

func header(pos int, check bool) uint8 {
	if check {
		return checkHeader[pos]
	}
	return nonCheckHeader[pos]
}
