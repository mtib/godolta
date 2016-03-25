package deltal

import (
	"bytes"
	"compress/gzip"
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
	return &Encoder{passarray: stringToArray(password), FileReader: reader, UseChecksum: checksum, Checksum: rustHash(full)}, err
}

// NewEncoder returns pointer to decoder reader interface
func NewEncoder(file, password string, checksum bool) (*Encoder, error) {
	f, _ := os.Open(file)
	return NewEncoderReader(f, password, checksum)
}

type compressedByte struct {
	data   []byte
	offset int
}

func (b *compressedByte) Read(p []byte) (n int, err error) {
	delta := len(b.data) - b.offset
	if len(p) >= delta {
		for k, v := range b.data[b.offset:] {
			p[k] = v
		}
		b.offset = len(b.data)
		return delta, io.EOF
	}
	for k, v := range b.data[b.offset : b.offset+len(p)] {
		p[k] = v
	}
	b.offset += len(p)
	return len(p), nil
}

// Seek makes this a io.Seeker
func (b *compressedByte) Seek(offset int64, whence int) (int64, error) {
	var err error
	switch whence {
	case 0:
		b.offset = int(offset)
	case 2:
		b.offset = len(b.data)
		fallthrough
	case 1:
		b.offset += int(offset)
	}
	if int(b.offset) > len(b.data) {
		err = io.EOF
	}
	return int64(b.offset), err
}

// NewCompressedEncoderReader gzips the file before encryption
func NewCompressedEncoderReader(reader io.ReadSeeker, password string, checksum bool) (*Encoder, error) {
	var b bytes.Buffer
	gzwriter := gzip.NewWriter(&b)
	io.Copy(gzwriter, reader)
	gzwriter.Flush()
	gzwriter.Close()
	buf, err := ioutil.ReadAll(&b)
	return &Encoder{passarray: stringToArray(password), FileReader: &compressedByte{buf, 0}, UseChecksum: checksum, Checksum: rustHash(buf)}, err
}

func stringToArray(p string) (pa []byte) {
	if p != "" {
		pa = rustHash([]byte(p))
	} else {
		pa = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}
	return
}

// Reset resets the encoder, synonym to Seek(0,0)
func (d *Encoder) Reset() {
	d.Counter = 0
	d.headerpos = 0
	d.FileReader.Seek(0, 0)
	d.last = 0
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

// FastEncrypt is an easy-call function
func FastEncrypt(file, pass string) {
	filein, _ := os.Open(file)
	encoder, _ := NewEncoderReader(filein, pass, true)
	encryptedData, _ := ioutil.ReadAll(encoder)
	ioutil.WriteFile(file+".delta", encryptedData, os.ModePerm)
}

func benchEncrypter(file, pass string) (func(), *os.File) {
	filein, _ := os.Open(file)
	encoder, _ := NewEncoderReader(filein, pass, true)
	return func() {
		ioutil.ReadAll(encoder)
		encoder.Reset()
	}, filein
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
